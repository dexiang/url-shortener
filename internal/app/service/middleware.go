package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Middleware func(URLShortenerService) URLShortenerService

type loggingMiddleware struct {
	logger log.Logger
	next   URLShortenerService
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next URLShortenerService) URLShortenerService {
		return loggingMiddleware{level.Info(logger), next}
	}
}

func (mv loggingMiddleware) ShortenURL(ctx context.Context, url string, expireAt time.Time) (id string, err error) {
	defer func() {
		mv.logger.Log("method", "ShortenURL", "url", url, "expireAt", expireAt, "err", err)
	}()

	return mv.next.ShortenURL(ctx, url, expireAt)
}

func (mv loggingMiddleware) GetOriginalURL(ctx context.Context, id string) (url string, err error) {
	defer func() {
		mv.logger.Log("method", "GetOriginalURL", "id", id, "err", err)
	}()

	return mv.next.GetOriginalURL(ctx, id)
}
