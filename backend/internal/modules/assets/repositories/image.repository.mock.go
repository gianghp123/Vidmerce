package repositories

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/stretchr/testify/mock"
)

type MockImageRepository struct {
	mock.Mock
}

func NewMockImageRepository() *MockImageRepository {
	return &MockImageRepository{}
}

func (m *MockImageRepository) FindByAssetID(ctx context.Context, assetID string) ([]models.ImageEntity, error) {
	args := m.Called(ctx, assetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ImageEntity), args.Error(1)
}

func (m *MockImageRepository) FindByAssetIDAndOrder(ctx context.Context, assetID string, order int) (*models.ImageEntity, error) {
	args := m.Called(ctx, assetID, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ImageEntity), args.Error(1)
}

func (m *MockImageRepository) Delete(ctx context.Context, assetID string, imageID string) error {
	args := m.Called(ctx, assetID, imageID)
	return args.Error(0)
}
