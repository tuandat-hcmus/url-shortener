package service

import (
	"github.com/google/uuid"
	"github.com/tuandat-hcmus/url-shortener/internal/model"
)

//go:generate mockery --name=HealthCheck --filename healthcheck.go --outpkg healthcheck_svc_mock
type HealthCheck interface {
	Check() model.HealthStatus
}

type healthCheck struct {
	serviceName string
	instanceID  string
}

func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	if instanceID == "" {
		instanceID = uuid.New().String()
	}
	return &healthCheck{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

func (h *healthCheck) Check() model.HealthStatus {
	return model.HealthStatus{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}
}
