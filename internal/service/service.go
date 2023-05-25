package service

import (
	"context"
	"time"
)

type URLShortenerService interface {
	ShortenURL(ctx context.Context, url string, expireAt time.Time) (string, error)
	GetOriginalURL(ctx context.Context, url string) (string, error)
}

type tinyURLService struct{}

func (s tinyURLService) ShortenURL(ctx context.Context, url string, expireAt time.Time) (string, error) {
	return "<ID>", nil
}

func (s tinyURLService) GetOriginalURL(ctx context.Context, url string) (string, error) {
	return "<original_url>", nil
}

func New() URLShortenerService {
	return tinyURLService{}
}
