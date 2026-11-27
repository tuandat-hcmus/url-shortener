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

	"github.com/tuandat-hcmus/url-shortener/internal/service"
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

func TestRedirect(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name             string
		setupMockService func(ctx context.Context) *svc_mock.ShortenURL
		setupTestRequest func(ctx *gin.Context)
		expectedStatus   int
		expectedUrl      string
	}{
		{
			name: "successful redirect",
			setupMockService: func(ctx context.Context) *svc_mock.ShortenURL {
				mockSvc := svc_mock.NewShortenURL(t)
				mockSvc.On("GetUrlFromCode", ctx, "abc123").Return("http://example.com", nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/v1/links/redirect/abc123",
					nil,
				)
				ctx.Params = gin.Params{{Key: "code", Value: "abc123"}}
			},
			expectedStatus: http.StatusMovedPermanently,
			expectedUrl:    "http://example.com",
		},
		{
			name: "empty code parameter",
			setupMockService: func(ctx context.Context) *svc_mock.ShortenURL {
				return svc_mock.NewShortenURL(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/v1/links/redirect",
					nil,
				)
			},
			expectedStatus: http.StatusBadRequest,
			expectedUrl:    "",
		},
		{
			name: "url not found",
			setupMockService: func(ctx context.Context) *svc_mock.ShortenURL {
				mockSvc := svc_mock.NewShortenURL(t)
				mockSvc.On("GetUrlFromCode", ctx, "abc123").Return("", service.ErrUrlNotFound)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/v1/links/redirect/abc123",
					nil,
				)
				ctx.Params = gin.Params{{Key: "code", Value: "abc123"}}
			},
			expectedStatus: http.StatusNotFound,
			expectedUrl:    "",
		},
		{
			name: "internal server error",
			setupMockService: func(ctx context.Context) *svc_mock.ShortenURL {
				mockSvc := svc_mock.NewShortenURL(t)
				mockSvc.On("GetUrlFromCode", ctx, "abc123").Return("", errors.New("some internal error"))
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/v1/links/redirect/abc123",
					nil,
				)
				ctx.Params = gin.Params{{Key: "code", Value: "abc123"}}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedUrl:    "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			tc.setupTestRequest(ginCtx)
			svcMock := tc.setupMockService(ginCtx)
			handler := NewShortenURL(svcMock)
			handler.Redirect(ginCtx)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedUrl != "" {
				assert.Equal(t, tc.expectedUrl, rec.Header().Get("Location"))
			}
		})
	}
}
