package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/stretchr/testify/mock"
)

type MockAssetRepository struct {
	mock.Mock
}

func NewMockAssetRepository() *MockAssetRepository {
	return &MockAssetRepository{}
}

func (m *MockAssetRepository) DBClient() *dynamodb.Client {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dynamodb.Client)
}

func (m *MockAssetRepository) FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.AssetEntity], error) {
	args := m.Called(ctx, limit, lastKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.PaginatedResult[models.AssetEntity]), args.Error(1)
}

func (m *MockAssetRepository) FindByID(ctx context.Context, id string) (*models.AssetEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssetEntity), args.Error(1)
}

func (m *MockAssetRepository) Create(ctx context.Context, asset models.AssetEntity) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *MockAssetRepository) UpdateAssetStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockAssetRepository) IncrementImageCount(ctx context.Context, id string) (int, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int), args.Error(1)
}

func (m *MockAssetRepository) TransactWriteItems(ctx context.Context, items ...interface{}) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}
