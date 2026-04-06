package repositories

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
)

type VideoRepository interface {
	FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.VideoMetadataEntity], error)
	FindByID(ctx context.Context, id string) (*models.VideoMetadataEntity, error)
	FindInteractiveByVideoID(ctx context.Context, videoID string) (*models.VideoInteractiveEntity, error)
}

type videoRepository struct {
	dbClient *dynamodb.Client
}

func NewVideoRepository(dbClient *dynamodb.Client) VideoRepository {
	return &videoRepository{dbClient: dbClient}
}

func (r *videoRepository) FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.VideoMetadataEntity], error) {
	exclusiveStartKey, err := core.DecodeCursor(lastKey)
	if err != nil {
		return nil, err
	}

	keyCond := expression.Key("GSI1PK").Equal(expression.Value("ENTITY#VIDEO"))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(core.TableName),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit)),
		ScanIndexForward:          aws.Bool(false),
	}
	if exclusiveStartKey != nil {
		input.ExclusiveStartKey = exclusiveStartKey
	}

	resp, err := r.dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var items []models.VideoMetadataEntity
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &items); err != nil {
		return nil, err
	}

	// Strip "VIDEO#" prefix from PK for each item
	for i := range items {
		items[i].PK = strings.TrimPrefix(items[i].PK, "VIDEO#")
	}

	nextCursor, err := core.EncodeCursor(resp.LastEvaluatedKey)
	if err != nil {
		return nil, err
	}

	hasMore := resp.LastEvaluatedKey != nil
	return &response.PaginatedResult[models.VideoMetadataEntity]{
		Data: items,
		Meta: response.NewCursorMeta(limit, nextCursor, hasMore),
	}, nil
}

func (r *videoRepository) FindByID(ctx context.Context, id string) (*models.VideoMetadataEntity, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": "VIDEO#" + id,
		"SK": "METADATA",
	})
	if err != nil {
		return nil, err
	}

	resp, err := r.dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(core.TableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}

	if resp.Item == nil {
		return nil, nil
	}

	var item models.VideoMetadataEntity
	if err := attributevalue.UnmarshalMap(resp.Item, &item); err != nil {
		return nil, err
	}
	// Strip "VIDEO#" prefix from PK to return raw video ID
	item.PK = strings.TrimPrefix(item.PK, "VIDEO#")
	return &item, nil
}

func (r *videoRepository) FindInteractiveByVideoID(ctx context.Context, videoID string) (*models.VideoInteractiveEntity, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": "VIDEO#" + videoID,
		"SK": "INTERACTIVE",
	})
	if err != nil {
		return nil, err
	}

	resp, err := r.dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(core.TableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}

	if resp.Item == nil {
		return nil, nil
	}

	var item models.VideoInteractiveEntity
	if err := attributevalue.UnmarshalMap(resp.Item, &item); err != nil {
		return nil, err
	}
	// Strip "VIDEO#" prefix from PK to return raw video ID
	item.PK = strings.TrimPrefix(item.PK, "VIDEO#")
	return &item, nil
}
