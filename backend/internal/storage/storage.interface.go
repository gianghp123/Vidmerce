package storage

import (
	"context"
	"time"
)

type Storage interface {
	GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expire time.Duration) (string, error)
	GeneratePresignedGetURL(ctx context.Context, key string, expire time.Duration) (string, error)
	ObjectExists(ctx context.Context, key string) (bool, error)
	DeleteObject(ctx context.Context, key string) error
}
