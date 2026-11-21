package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuandat-hcmus/url-shortener/internal/service"
	"time"
)

type ShortenURL interface {
	CreateShortURL(c *gin.Context)
}

type shortenURL struct {
	svc service.ShortenURL
}

// NewShortenURL creates a new instance of ShortenURL handler
func NewShortenURL(svc service.ShortenURL) ShortenURL {
	return &shortenURL{svc: svc}
}

type shortenURLRequest struct {
	URL     string `json:"url"`
	ExpTime int64  `json:"exp"`
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
		c.AbortWithStatusJSON((http.StatusInternalServerError), gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, &shortenURLResponse{
		Code:    code,
		Message: "Short URL created successfully",
	})
}
