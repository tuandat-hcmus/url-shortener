package service

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
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
				urlStorageMock.On("GetURL", ctx, "abc123").Return("", nil)
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
				urlStorageMock.On("GetURL", ctx, "def456").Return("", nil)
				urlStorageMock.On("StoreURL", ctx, "def456", "http://example.com", 1*time.Minute).Return(nil)
				return urlStorageMock, codeGenMock
			},
			expectedCode: "def456",
			expectedErr:  nil,
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
