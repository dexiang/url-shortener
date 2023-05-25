package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/dexiang/url-shortener/internal/endpoints"
	"github.com/dexiang/url-shortener/internal/service"
	"github.com/dexiang/url-shortener/internal/transports"
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
	service := service.New(logger)
	endpoints := endpoints.New(service, logger)
	handler := transports.NewHTTPHandler(ctx, endpoints)

	http.ListenAndServe(":80", handler)
}
