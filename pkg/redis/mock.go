package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"testing"
)

// InitMockRedis initializes a mock Redis client for testing purposes
func InitMockRedis(t *testing.T) *redis.Client {
	mock := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{
		Addr: mock.Addr(),
	})
}
