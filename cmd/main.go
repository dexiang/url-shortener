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

func main() {

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
	handler := transports.NewHTTPHandler(ctx, endpoints)

	logger.Log("Listening on port 8080")
	http.ListenAndServe(":8080", handler)
}
