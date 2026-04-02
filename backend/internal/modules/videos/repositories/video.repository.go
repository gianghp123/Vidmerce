package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gianghp123/Vidmerce/backend/services/internal/database/models"
)

// VideoRepository defines data access for videos
type VideoRepository interface {
	FindAll(ctx context.Context, limit int, lastKey string) ([]models.VideoMetadataEntity, string, bool, error)
	FindByID(ctx context.Context, id string) (*models.VideoMetadataEntity, error)
	FindInteractiveByVideoID(ctx context.Context, videoID string) (*models.VideoInteractiveEntity, error)
}

type videoRepository struct {
	dbClient *dynamodb.Client
}

func NewVideoRepository(dbClient *dynamodb.Client) VideoRepository {
	return &videoRepository{dbClient: dbClient}
}

func (r *videoRepository) FindAll(ctx context.Context, limit int, lastKey string) ([]models.VideoMetadataEntity, string, bool, error) {
	// TODO: implement DynamoDB scan/query with cursor
	return []models.VideoMetadataEntity{}, "", false, nil
}

func (r *videoRepository) FindByID(ctx context.Context, id string) (*models.VideoMetadataEntity, error) {
	// TODO: implement DynamoDB getItem
	return nil, nil
}

func (r *videoRepository) FindInteractiveByVideoID(ctx context.Context, videoID string) (*models.VideoInteractiveEntity, error) {
	// TODO: implement DynamoDB getItem for interactive data
	return nil, nil
}
