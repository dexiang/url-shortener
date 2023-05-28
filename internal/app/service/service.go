package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"

	"github.com/dexiang/url-shortener/internal/app/model"
)

var ErrIDNotFound = errors.New("cannot found url id")

type URLShortenerService interface {
	ShortenURL(ctx context.Context, url string, expireAt time.Time) (string, error)
	GetOriginalURL(ctx context.Context, id string) (string, error)
}

type tinyURLService struct {
	repo model.URLRepository
}

func (s tinyURLService) ShortenURL(ctx context.Context, url string, expireAt time.Time) (string, error) {
	uuid := uuid.New()
	key := uuid.String()
	s.repo.Set(ctx, key, url)
	return key, nil
}

func (s tinyURLService) GetOriginalURL(ctx context.Context, id string) (string, error) {
	original_url, err := s.repo.GetByID(ctx, id)
	if err == model.ErrIDNotFound {
		return "", ErrIDNotFound
	}
	return original_url, err
}

func NewTinyURLService(repo model.URLRepository) URLShortenerService {
	return tinyURLService{
		repo: repo,
	}
}

func New(repo model.URLRepository, logger log.Logger) URLShortenerService {
	var svc URLShortenerService
	{
		svc = NewTinyURLService(repo)
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}
