package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// URLStorage defines the interface for storing and retrieving URLs
//
//go:generate mockery --name URLStorage --filename urlstorage.go --outpkg repo_mock
var NoKeyFound = errors.New("No key found")

type URLStorage interface {
	StoreURL(ctx context.Context, key, url string, expTime time.Duration) error
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

// NewURLStorage creates a new instance of URLStorage
func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{c: c}
}

// StoreURL stores a URL with the given key and expiration time
func (s *urlStorage) StoreURL(ctx context.Context, key, url string, expTime time.Duration) error {
	return s.c.Set(ctx, key, url, expTime).Err()
}

// GetURL retrieves a URL by its key
func (s *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	url, err := s.c.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", NoKeyFound
		}
		return "", err
	}
	return url, nil
}
