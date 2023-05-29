package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"

	"github.com/dexiang/url-shortener/internal/app/model"
)

const maxExpirationTime = 3 * 365 * 24 * time.Hour // 3 Year

var ErrIDNotFound = errors.New("cannot found url id")
var ErrMaxExpirationTimeExceeded = errors.New("the expiration time cannot exceed three years")
var ErrTimeExpired = errors.New("cannot set expired time")

type URLShortenerService interface {
	ShortenURL(ctx context.Context, url string, expireAt time.Time) (string, error)
	GetOriginalURL(ctx context.Context, id string) (string, error)
}

type tinyURLService struct {
	repo model.URLRepository
}

func (s tinyURLService) ShortenURL(ctx context.Context, url string, expireAt time.Time) (string, error) {

	current := time.Now()
	diff := expireAt.Sub(current)

	if diff > maxExpirationTime {
		return "", ErrMaxExpirationTimeExceeded
	}
	if diff < 0 {
		return "", ErrTimeExpired
	}

	var uid string
	for uid == "" {
		key := uuid.New().String()
		if isExist, _ := s.repo.Exists(ctx, key); !isExist {
			uid = key
		}
	}

	s.repo.Set(ctx, uid, url, diff)

	return uid, nil
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
