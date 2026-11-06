package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuandat-hcmus/url-shortener/internal/api"
	"github.com/tuandat-hcmus/url-shortener/internal/model"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name             string
		setupTestRequest func() (*httptest.ResponseRecorder, *http.Request)
		expectedStatus   int
	}{
		{
			name: "case 1: normal health check request",
			setupTestRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				req := httptest.NewRequest("GET", "/health-check", nil)
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "case 2: invalid method",
			setupTestRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				req := httptest.NewRequest("POST", "/health-check", nil)
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedStatus: http.StatusNotFound,
		},
	}
	cfg, err := api.NewConfig()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testApiEngine := api.NewEngine(cfg)
			rec, req := tc.setupTestRequest()
			testApiEngine.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedStatus != http.StatusOK {
				return
			}
			err := json.Unmarshal(rec.Body.Bytes(), &model.HealthStatus{})
			assert.NoError(t, err)
		})
	}
}
