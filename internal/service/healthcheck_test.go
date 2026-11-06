package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		serviceName string
		instanceID  string
		expectedInstanceID string
	}{
		{
			name: "case 1: instanceID is provided",
			serviceName: "test-service",
			instanceID:  "test-instance-id",
			expectedInstanceID: "test-instance-id",
		}, 
		{
			name: "case 2: instanceID is empty",
			serviceName: "test-service",
			instanceID:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := NewHealthCheck(
				tc.serviceName,
				tc.instanceID,
			)
			status := svc.Check()
			assert.Equal(t, "OK", status.Message)
			assert.Equal(t, tc.serviceName, status.ServiceName)
			if tc.expectedInstanceID != "" {
				assert.Equal(t, tc.expectedInstanceID, status.InstanceID)
				return 
			}
			assert.NotEmpty(t, status.InstanceID)
			_, err := uuid.Parse(status.InstanceID)
			assert.NoError(t, err)
		})
	}
}