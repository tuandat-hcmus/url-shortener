package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tuandat-hcmus/url-shortener/internal/handler"
	"github.com/tuandat-hcmus/url-shortener/internal/service"
	"fmt"
	"net/http"
)

type Engine interface {
	// Start starts the API engine and listens for incoming requests.
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	// Add any necessary fields for your API engine here.
	app *gin.Engine
	cfg *Config
}

func NewEngine(cfg *Config) Engine {
	e := &engine{
		app: gin.Default(),
		cfg: cfg,
	}
	e.initRoutes()
	return e
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%d", e.cfg.Port))
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

func (e *engine) initRoutes() {
	healCheckSvc := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID)
	healthCheckHandler := handler.NewHealthCheck(healCheckSvc)
	e.app.GET("/health-check", healthCheckHandler.Check)
}
