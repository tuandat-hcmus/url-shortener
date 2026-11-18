package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "github.com/tuandat-hcmus/url-shortener/docs"
	"github.com/tuandat-hcmus/url-shortener/internal/handler"
	"github.com/tuandat-hcmus/url-shortener/internal/repository"
	"github.com/tuandat-hcmus/url-shortener/internal/service"
	"github.com/tuandat-hcmus/url-shortener/pkg/utils"
	"net/http"
)

type Engine interface {
	// Start starts the API engine and listens for incoming requests.
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	// Add any necessary fields for your API engine here.
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
}

func NewEngine(cfg *Config, redisClient *redis.Client) Engine {
	e := &engine{
		app:         gin.Default(),
		cfg:         cfg,
		redisClient: redisClient,
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
	urlStorageRepo := repository.NewURLStorage(e.redisClient)
	codeGen := utils.NewCodeGen()
	shortenURLSvc := service.NewShortenURL(urlStorageRepo, codeGen, e.cfg.CodeLength)
	shortenURLHandler := handler.NewShortenURL(shortenURLSvc)
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.app.GET("/health-check", healthCheckHandler.Check)
	v1Routes := e.app.Group("/v1")
	{
		v1Routes.POST("/links/shorten", shortenURLHandler.CreateShortURL)
	}
}
