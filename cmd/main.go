package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/redis/go-redis/v9"

	"github.com/dexiang/url-shortener/internal/app/endpoints"
	"github.com/dexiang/url-shortener/internal/app/model"
	"github.com/dexiang/url-shortener/internal/app/service"
	"github.com/dexiang/url-shortener/internal/app/transports"
)

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		os.Setenv(env, fallback)
		return fallback
	}
	return e
}

func main() {

	var (
		ServiceHost = envString("SERVICE_HOST", "localhost")
		HttpPort    = envString("SERVICE_PORT", "80")
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	ctx := context.Background()

	db := redis.NewClient(&redis.Options{
		Addr:     "redis-master:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := db.Ping(ctx).Result()
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	repository := model.New(db, logger)
	service := service.New(repository, logger)
	endpoints := endpoints.New(service, logger)
	handler := transports.NewHTTPHandler(ctx, endpoints, logger)

	logger.Log("Listening ", ServiceHost, " on port ", HttpPort)
	http.ListenAndServe(":"+HttpPort, handler)
}
