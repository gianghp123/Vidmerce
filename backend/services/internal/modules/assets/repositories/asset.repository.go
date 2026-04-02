package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gianghp123/Vidmerce/backend/services/internal/database/models"
)

type AssetRepository interface {
	FindAll(ctx context.Context, limit int, lastKey string) ([]models.AssetEntity, string, bool, error)
	FindByID(ctx context.Context, id string) (*models.AssetEntity, error)
	Create(ctx context.Context, asset models.AssetEntity) error
}

type assetRepository struct {
	dbClient *dynamodb.Client
}

func NewAssetRepository(dbClient *dynamodb.Client) AssetRepository {
	return &assetRepository{dbClient: dbClient}
}

func (r *assetRepository) FindAll(ctx context.Context, limit int, lastKey string) ([]models.AssetEntity, string, bool, error) {
	return []models.AssetEntity{}, "", false, nil
}

func (r *assetRepository) FindByID(ctx context.Context, id string) (*models.AssetEntity, error) {
	return nil, nil
}

func (r *assetRepository) Create(ctx context.Context, asset models.AssetEntity) error {
	return nil
}
