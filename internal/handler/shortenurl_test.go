package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	svc_mock "github.com/tuandat-hcmus/url-shortener/internal/service/mocks"
)

// TestCreateShortURL tests the CreateShortURL handler
func TestCreateShortURL(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *svc_mock.ShortenURL
		expectedStatus   int
		expectedBody     string
	}{
		{
			name: "successful creation",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"http://example.com","exp":60}`)),
				)
			},
			setupMockService: func(ctx context.Context) *svc_mock.ShortenURL {
				svcMock := svc_mock.NewShortenURL(t)
				svcMock.On("CreateShortURL", ctx, "http://example.com", 1*time.Minute).Return("abc123", nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":"abc123","message":"Short URL created successfully"}`,
		},
		{
			name: "invalid request body",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"http://example.com"`)), // Invalid JSON
				)
			},
			setupMockService: func(ctx context.Context) *svc_mock.ShortenURL {
				return svc_mock.NewShortenURL(t) // No need to set up any expectations for this test case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid Request"}`,
		},
		{
			name: "service error",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"http://example.com","exp":60}`)),
				)
			},
			setupMockService: func(ctx context.Context) *svc_mock.ShortenURL {
				svcMock := svc_mock.NewShortenURL(t)
				svcMock.On("CreateShortURL", ctx, "http://example.com", 1*time.Minute).Return("", errors.New(""))
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Internal Server Error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ginCtx)
			svcMock := tc.setupMockService(ginCtx)
			handler := NewShortenURL(svcMock)
			handler.CreateShortURL(ginCtx)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedBody, rec.Body.String())
		})
	}
}
