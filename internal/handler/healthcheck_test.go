package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tuandat-hcmus/url-shortener/internal/model"
	healthcheck_svc_mock "github.com/tuandat-hcmus/url-shortener/internal/service/mocks"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		setupMockSvc func(ctx context.Context) *healthcheck_svc_mock.HealthCheck
		expectedStatus int 
		expectedBody model.HealthStatus
	}{
		{
			name: "case 1: health check is OK",
			setupMockSvc: func(ctx context.Context) *healthcheck_svc_mock.HealthCheck {
				mockSvc := healthcheck_svc_mock.NewHealthCheck(t)
				mockSvc.On("Check").Return(model.HealthStatus{
					Message:     "OK",
					ServiceName: "test-service",
					InstanceID:  "test-instance-id",
				})
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedBody: model.HealthStatus{
				Message:     "OK",
				ServiceName: "test-service",
				InstanceID:  "test-instance-id",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			mockSvc := tc.setupMockSvc(ginCtx)
			healthCheckHandler := NewHealthCheck(mockSvc)
			healthCheckHandler.Check(ginCtx)
			var status model.HealthStatus
			err := json.Unmarshal(rec.Body.Bytes(), &status)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedBody, status)
		})
	}
}