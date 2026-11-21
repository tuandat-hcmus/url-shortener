package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/tuandat-hcmus/url-shortener/internal/api"
	redisPkg "github.com/tuandat-hcmus/url-shortener/pkg/redis"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestShortenURL tests the shorten URL endpoint
func TestShortenURL(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name             string
		setupTestRequest func() (rec *httptest.ResponseRecorder, req *http.Request)
		setupMock        func(ctx context.Context) *redis.Client
		expectedStatus   int
		verifyFunc       func(ctx context.Context, c *redis.Client, code string)
	}{
		{
			name: "successful shorten URL creation",
			setupTestRequest: func() (rec *httptest.ResponseRecorder, req *http.Request) {
				rec = httptest.NewRecorder()
				req = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"http://example.com","exp":60}`)),
				)
				return
			},
			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			expectedStatus: http.StatusOK,
			verifyFunc: func(ctx context.Context, c *redis.Client, code string) {
				url, err := c.Get(ctx, code).Result()
				assert.NoError(t, err)
				assert.Equal(t, "http://example.com", url)
			},
		},
	}
	cfg, err := api.NewConfig()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRedis := tc.setupMock(t.Context())
			testApiEngine := api.NewEngine(cfg, mockRedis)
			rec, req := tc.setupTestRequest()
			testApiEngine.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.verifyFunc != nil {
				var response struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				tc.verifyFunc(t.Context(), mockRedis, response.Code)
			}
		})
	}
}
