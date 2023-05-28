package model

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/go-kit/kit/log"
)

var ErrIDNotFound = errors.New("cannot found url id")

type URLRepository interface {
	Set(context.Context, string, string) error
	GetByID(context.Context, string) (string, error)
}

type urlRepository struct {
	db *redis.Client
}

func (m urlRepository) Set(ctx context.Context, id string, url string) error {
	err := m.db.Set(ctx, id, url, 0).Err()
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
