package model

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-kit/kit/log"
)

var ErrIDNotFound = errors.New("cannot found url id")

type URLRepository interface {
	Exists(context.Context, string) (bool, error)
	Set(context.Context, string, string, time.Duration) error
	GetByID(context.Context, string) (string, error)
}

type urlRepository struct {
	db *redis.Client
}

func (m urlRepository) Exists(ctx context.Context, id string) (bool, error) {
	n, err := m.db.Exists(ctx, id).Result()
	return n > 0, err
}

func (m urlRepository) Set(ctx context.Context, id string, url string, expiration time.Duration) error {
	err := m.db.SetEx(ctx, id, url, expiration).Err()
	return err
}

func (m urlRepository) GetByID(ctx context.Context, id string) (string, error) {
	originalUrl, err := m.db.Get(ctx, id).Result()
	if err == redis.Nil {
		return "", ErrIDNotFound
	}
	return originalUrl, err
}

func NewURLRepository(db *redis.Client) URLRepository {
	return urlRepository{
		db: db,
	}
}

func New(db *redis.Client, logger log.Logger) URLRepository {
	var model URLRepository
	{
		model = NewURLRepository(db)
	}
	return model
}
