package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
)

type AssetRepository interface {
	FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.AssetEntity], error)
	FindByID(ctx context.Context, id string) (*models.AssetEntity, error)
	Create(ctx context.Context, asset models.AssetEntity) error
	UpdateAssetStatus(ctx context.Context, id string, status string) error
}

type assetRepository struct {
	dbClient *dynamodb.Client
}

func NewAssetRepository(dbClient *dynamodb.Client) AssetRepository {
	return &assetRepository{dbClient: dbClient}
}

func (r *assetRepository) FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.AssetEntity], error) {
	exclusiveStartKey, err := core.DecodeCursor(lastKey)
	if err != nil {
		return nil, err
	}

	keyCond := expression.Key("GSI1PK").Equal(expression.Value("ENTITY#ASSET"))
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

	var items []models.AssetEntity
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &items); err != nil {
		return nil, err
	}

	nextCursor, err := core.EncodeCursor(resp.LastEvaluatedKey)
	if err != nil {
		return nil, err
	}

	hasMore := resp.LastEvaluatedKey != nil
	return &response.PaginatedResult[models.AssetEntity]{
		Data: items,
		Meta: response.NewCursorMeta(limit, nextCursor, hasMore),
	}, nil
}

func (r *assetRepository) FindByID(ctx context.Context, id string) (*models.AssetEntity, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": "ASSET#" + id,
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

	var item models.AssetEntity
	if err := attributevalue.UnmarshalMap(resp.Item, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *assetRepository) Create(ctx context.Context, asset models.AssetEntity) error {
	item, err := attributevalue.MarshalMap(asset)
	if err != nil {
		return err
	}

	resp, err := r.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(core.TableName),
		Item:      item,
	})
	if err != nil {
		return err
	}

	if resp.Attributes != nil {
		if err := attributevalue.UnmarshalMap(resp.Attributes, &asset); err != nil {
			return err
		}
	}

	return nil
}

func (r *assetRepository) UpdateAssetStatus(ctx context.Context, id string, status string) error {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": "ASSET#" + id,
		"SK": "METADATA",
	})
	if err != nil {
		return err
	}

	_, err = r.dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(core.TableName),
		Key:              key,
		UpdateExpression: aws.String("SET #status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: status},
		},
	})
	return err
}
