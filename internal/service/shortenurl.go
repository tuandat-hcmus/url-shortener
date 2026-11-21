package service

import (
	"context"
	"github.com/tuandat-hcmus/url-shortener/internal/repository"
	"github.com/tuandat-hcmus/url-shortener/pkg/utils"
	"time"
)

// ShortenURL defines the interface for creating short URLs
//
//go:generate mockery --name ShortenURL --filename shortenurl.go --outpkg svc_mock
type ShortenURL interface {
	CreateShortURL(ctx context.Context, url string, expTime time.Duration) (string, error)
}

type shortenURL struct {
	r          repository.URLStorage
	c          utils.CodeGen
	codeLength int
}

// NewShortenURL creates a new instance of ShortenURL service
func NewShortenURL(r repository.URLStorage, c utils.CodeGen, codeLength int) ShortenURL {
	return &shortenURL{r: r, c: c, codeLength: codeLength}
}

// CreateShortURL creates a short URL for the given long URL
func (s *shortenURL) CreateShortURL(ctx context.Context, url string, expTime time.Duration) (string, error) {
	code, err := s.c.GenerateCode(s.codeLength)
	if err != nil {
		return "", err
	}
	_, err = s.r.GetURL(ctx, code)
	if err == nil {
		return s.CreateShortURL(ctx, url, expTime)
	} else {
		if err != repository.NoKeyFound {
			return "", err
		}
	}
	err = s.r.StoreURL(ctx, code, url, expTime)
	if err != nil {
		return "", err
	}
	return code, nil
}
