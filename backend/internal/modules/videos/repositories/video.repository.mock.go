package repositories

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/stretchr/testify/mock"
)

type MockVideoRepository struct {
	mock.Mock
}

func NewMockVideoRepository() *MockVideoRepository {
	return &MockVideoRepository{}
}

func (m *MockVideoRepository) FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.VideoMetadataEntity], error) {
	args := m.Called(ctx, limit, lastKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.PaginatedResult[models.VideoMetadataEntity]), args.Error(1)
}

func (m *MockVideoRepository) FindByID(ctx context.Context, id string) (*models.VideoMetadataEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VideoMetadataEntity), args.Error(1)
}

func (m *MockVideoRepository) FindInteractiveByVideoID(ctx context.Context, videoID string) (*models.VideoInteractiveEntity, error) {
	args := m.Called(ctx, videoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VideoInteractiveEntity), args.Error(1)
}
