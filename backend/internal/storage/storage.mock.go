package storage

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (m *MockStorage) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expire time.Duration) (string, error) {
	args := m.Called(ctx, key, contentType, expire)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GeneratePresignedGetURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	args := m.Called(ctx, key, expire)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) ObjectExists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}
