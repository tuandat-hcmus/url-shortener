package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tuandat-hcmus/url-shortener/internal/repository"
	repo_mock "github.com/tuandat-hcmus/url-shortener/internal/repository/mocks"
	utils_mock "github.com/tuandat-hcmus/url-shortener/pkg/utils/mocks"
	"testing"
	"time"
)

func TestCreateShortURL(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		setupMock    func(ctx context.Context) (*repo_mock.URLStorage, *utils_mock.CodeGen)
		expectedCode string
		expectedErr  error
	}{
		{
			name: "successful creation",
			setupMock: func(ctx context.Context) (*repo_mock.URLStorage, *utils_mock.CodeGen) {
				urlStorageMock := repo_mock.NewURLStorage(t)
				codeGenMock := utils_mock.NewCodeGen(t)
				codeGenMock.On("GenerateCode", 6).Return("abc123", nil)
				urlStorageMock.On("GetURL", ctx, "abc123").Return("", repository.NoKeyFound)
				urlStorageMock.On("StoreURL", ctx, "abc123", "http://example.com", 1*time.Minute).Return(nil)
				return urlStorageMock, codeGenMock
			},
			expectedCode: "abc123",
			expectedErr:  nil,
		},
		{
			name: "code generation fails",
			setupMock: func(ctx context.Context) (*repo_mock.URLStorage, *utils_mock.CodeGen) {
				urlStorageMock := repo_mock.NewURLStorage(t)
				codeGenMock := utils_mock.NewCodeGen(t)
				codeGenMock.On("GenerateCode", 6).Return("", fmt.Errorf("code generation failed"))
				return urlStorageMock, codeGenMock
			},
			expectedCode: "",
			expectedErr:  fmt.Errorf("code generation failed"),
		},
		{
			name: "code already exists",
			setupMock: func(ctx context.Context) (*repo_mock.URLStorage, *utils_mock.CodeGen) {
				urlStorageMock := repo_mock.NewURLStorage(t)
				codeGenMock := utils_mock.NewCodeGen(t)
				codeGenMock.On("GenerateCode", 6).Return("abc123", nil).Once()
				codeGenMock.On("GenerateCode", 6).Return("def456", nil).Once()
				urlStorageMock.On("GetURL", ctx, "abc123").Return("http://example.com", nil)
				urlStorageMock.On("GetURL", ctx, "def456").Return("", repository.NoKeyFound)
				urlStorageMock.On("StoreURL", ctx, "def456", "http://example.com", 1*time.Minute).Return(nil)
				return urlStorageMock, codeGenMock
			},
			expectedCode: "def456",
			expectedErr:  nil,
		},
		{
			name: "storage error on store",
			setupMock: func(ctx context.Context) (*repo_mock.URLStorage, *utils_mock.CodeGen) {
				urlStorageMock := repo_mock.NewURLStorage(t)
				codeGenMock := utils_mock.NewCodeGen(t)
				codeGenMock.On("GenerateCode", 6).Return("abc123", nil)
				urlStorageMock.On("GetURL", ctx, "abc123").Return("", repository.NoKeyFound)
				urlStorageMock.On("StoreURL", ctx, "abc123", "http://example.com", 1*time.Minute).Return(errors.New("storage error"))
				return urlStorageMock, codeGenMock
			},
			expectedCode: "",
			expectedErr:  errors.New("storage error"),
		},
		{
			name: "storage error on get",
			setupMock: func(ctx context.Context) (*repo_mock.URLStorage, *utils_mock.CodeGen) {
				urlStorageMock := repo_mock.NewURLStorage(t)
				codeGenMock := utils_mock.NewCodeGen(t)
				codeGenMock.On("GenerateCode", 6).Return("abc123", nil)
				urlStorageMock.On("GetURL", ctx, "abc123").Return("", errors.New("storage error"))
				return urlStorageMock, codeGenMock
			},
			expectedCode: "",
			expectedErr:  errors.New("storage error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			urlStorageMock, codeGenMock := tc.setupMock(ctx)
			svc := NewShortenURL(urlStorageMock, codeGenMock, 6)
			res, err := svc.CreateShortURL(ctx, "http://example.com", 1*time.Minute)
			assert.Equal(t, tc.expectedCode, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetUrlFromCode(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		setupMock   func(ctx context.Context) *repo_mock.URLStorage
		expectedUrl string
		expectedErr error
	}{
		{
			name: "successful retrieval",
			setupMock: func(ctx context.Context) *repo_mock.URLStorage {
				urlStorageMock := repo_mock.NewURLStorage(t)
				urlStorageMock.On("GetURL", ctx, "abc123").Return("http://example.com", nil)
				return urlStorageMock
			},
			expectedUrl: "http://example.com",
			expectedErr: nil,
		},
		{
			name: "url not found",
			setupMock: func(ctx context.Context) *repo_mock.URLStorage {
				urlStorageMock := repo_mock.NewURLStorage(t)
				urlStorageMock.On("GetURL", ctx, "abc123").Return("", repository.NoKeyFound)
				return urlStorageMock
			},
			expectedUrl: "",
			expectedErr: ErrUrlNotFound,
		},
		{
			name: "storage error",
			setupMock: func(ctx context.Context) *repo_mock.URLStorage {
				urlStorageMock := repo_mock.NewURLStorage(t)
				urlStorageMock.On("GetURL", ctx, "abc123").Return("", errors.New("storage error"))
				return urlStorageMock
			},
			expectedUrl: "",
			expectedErr: errors.New("storage error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			urlStorageMock := tc.setupMock(ctx)
			svc := NewShortenURL(urlStorageMock, nil, 6)
			res, err := svc.GetUrlFromCode(ctx, "abc123")
			assert.Equal(t, tc.expectedUrl, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
