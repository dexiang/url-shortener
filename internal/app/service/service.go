package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type URLShortenerService interface {
	ShortenURL(ctx context.Context, url string, expireAt time.Time) (string, error)
	GetOriginalURL(ctx context.Context, url string) (string, error)
}

type tinyURLService struct{}

func (s tinyURLService) ShortenURL(ctx context.Context, url string, expireAt time.Time) (string, error) {
	return "<ID>", nil
}

func (s tinyURLService) GetOriginalURL(ctx context.Context, id string) (string, error) {
	return "<original_url>", nil
}

func NewTinyURLService() URLShortenerService {
	return tinyURLService{}
}

func New(logger log.Logger) URLShortenerService {
	var svc URLShortenerService
	{
		svc = NewTinyURLService()
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}
