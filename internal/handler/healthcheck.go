package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tuandat-hcmus/url-shortener/internal/service"
	"net/http"
)

type HealthCheck interface {
	Check(c *gin.Context)
}

type healthCheck struct {
	svc service.HealthCheck
}

func NewHealthCheck(svc service.HealthCheck) HealthCheck {
	return &healthCheck{
		svc: svc,
	}
}

func (h *healthCheck) Check(c *gin.Context) {
	status := h.svc.Check()
	c.JSON(http.StatusOK, status)
}
