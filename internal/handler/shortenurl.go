package handler

import (
	"net/http"

	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tuandat-hcmus/url-shortener/internal/service"
)

type ShortenURL interface {
	CreateShortURL(c *gin.Context)
	Redirect(c *gin.Context)
}

type shortenURL struct {
	svc service.ShortenURL
}

// NewShortenURL creates a new instance of ShortenURL handler
func NewShortenURL(svc service.ShortenURL) ShortenURL {
	return &shortenURL{svc: svc}
}

type shortenURLRequest struct {
	URL     string `json:"url" binding:"required,url"`
	ExpTime int64  `json:"exp" binding:"required,gt=0,lte=86400"` // Expiration time in seconds
}

type shortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CreateShortURL
// @Summary Create Short URL
// @Description Create a short URL for the given long URL
// @Tags links
// @Produces application/json
// @Param request body shortenURLRequest true "Request body"
// @Success 200 {object} shortenURLResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/shorten [post]
func (h *shortenURL) CreateShortURL(c *gin.Context) {
	input := &shortenURLRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON((http.StatusBadRequest), gin.H{"error": "Invalid Request"})
		return
	}

	code, err := h.svc.CreateShortURL(c, input.URL, time.Duration(input.ExpTime)*time.Second)
	if err != nil {
		log.Error().Err(err).Str("url", input.URL).Int64("exp", input.ExpTime).Msg("Generate short URL failed - ShortendURL Service")
		c.AbortWithStatusJSON((http.StatusInternalServerError), gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, &shortenURLResponse{
		Code:    code,
		Message: "Short URL created successfully",
	})
}

// Redirect
// @Summary Redirect to Original URL
// @Description Redirect to the original URL for the given short code
// @Tags links
// @Param code path string true "Short code"
// @Success 301
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/redirect/{code} [get]
func (h *shortenURL) Redirect(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.AbortWithStatusJSON((http.StatusBadRequest), gin.H{"error": "Invalid Input"})
		return
	}
	url, err := h.svc.GetUrlFromCode(c, code)
	if err != nil {
		if errors.Is(err, service.ErrUrlNotFound) {
			c.AbortWithStatusJSON((http.StatusNotFound), gin.H{"error": "Code Not Exists"})
			return
		}
		log.Error().Err(err).Str("code", code).Msg("Get URL from code failed - ShortendURL Service")
		c.AbortWithStatusJSON((http.StatusInternalServerError), gin.H{"error": "Internal Server Error"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, url)
}
