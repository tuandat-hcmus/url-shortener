package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/tuandat-hcmus/url-shortener/pkg/redis"
	"testing"
	"time"
)

// TestURLStorage_StoreURL tests the StoreURL method of the URLStorage interface
func TestURLStorage_StoreURL(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		setupMock   func(ctx context.Context) *redis.Client
		verifyFunc  func(ctx context.Context, c *redis.Client)
		expectedErr error
	}{
		{
			name: "normal case: store URL successfully",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			verifyFunc: func(ctx context.Context, c *redis.Client) {
				res, err := c.Get(ctx, "test-key").Result()
				assert.NoError(t, err)
				assert.Equal(t, "http://example.com", res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mock := tc.setupMock(ctx)
			testRepo := NewURLStorage(mock)
			err := testRepo.StoreURL(ctx, "test-key", "http://example.com", 1*time.Minute)
			assert.Equal(t, tc.expectedErr, err)
			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock)
			}
		})
	}
}

// TestURLStorage_GetURL tests the GetURL method of the URLStorage interface
func TestURLStorage_GetURL(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		setupMock   func(ctx context.Context) *redis.Client
		setupFunc   func(ctx context.Context, c *redis.Client)
		expectedURL string
		expectedErr error
	}{
		{
			name: "normal case: get URL successfully",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			setupFunc: func(ctx context.Context, c *redis.Client) {
				c.Set(ctx, "test-key", "http://example.com", 1*time.Minute)
			},
			expectedURL: "http://example.com",
			expectedErr: nil,
		},
		{
			name: "error case: key does not exist",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			setupFunc:   nil,
			expectedURL: "",
			expectedErr: redis.Nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mock := tc.setupMock(ctx)
			if tc.setupFunc != nil {
				tc.setupFunc(ctx, mock)
			}
			testRepo := NewURLStorage(mock)
			url, err := testRepo.GetURL(ctx, "test-key")
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedURL, url)
		})
	}
}
