package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
)

type ImageRepository interface {
	FindByAssetID(ctx context.Context, assetID string) ([]models.ImageEntity, error)
	FindByAssetIDAndOrder(ctx context.Context, assetID string, order int) (*models.ImageEntity, error)
	Create(ctx context.Context, images []models.ImageEntity) error
	UpdateStatus(ctx context.Context, assetID string, order int, status string) error
	Delete(ctx context.Context, assetID string, order int) error
}

type imageRepository struct {
	dbClient *dynamodb.Client
}

func NewImageRepository(dbClient *dynamodb.Client) ImageRepository {
	return &imageRepository{dbClient: dbClient}
}

func (r *imageRepository) FindByAssetID(ctx context.Context, assetID string) ([]models.ImageEntity, error) {
	keyCond := expression.Key("PK").Equal(expression.Value("ASSET#" + assetID)).
		And(expression.Key("SK").BeginsWith("IMAGE#"))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(core.TableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(true),
	}

	resp, err := r.dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var images []models.ImageEntity
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &images); err != nil {
		return nil, err
	}

	// Strip "ASSET#" prefix from PK for each image
	for i := range images {
		images[i].PK = strings.TrimPrefix(images[i].PK, "ASSET#")
	}

	return images, nil
}

func (r *imageRepository) Create(ctx context.Context, images []models.ImageEntity) error {
	for _, img := range images {
		item, err := attributevalue.MarshalMap(img)
		if err != nil {
			return err
		}

		_, err = r.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName:           aws.String(core.TableName),
			Item:                item,
			ConditionExpression: aws.String("attribute_not_exists(PK)"),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *imageRepository) UpdateStatus(ctx context.Context, assetID string, order int, status string) error {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": "ASSET#" + assetID,
		"SK": fmt.Sprintf("IMAGE#%d", order),
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

func (r *imageRepository) FindByAssetIDAndOrder(ctx context.Context, assetID string, order int) (*models.ImageEntity, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": "ASSET#" + assetID,
		"SK": fmt.Sprintf("IMAGE#%d", order),
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

	var img models.ImageEntity
	if err := attributevalue.UnmarshalMap(resp.Item, &img); err != nil {
		return nil, err
	}
	// Strip "ASSET#" prefix from PK to return raw asset ID
	img.PK = strings.TrimPrefix(img.PK, "ASSET#")
	return &img, nil
}

func (r *imageRepository) Delete(ctx context.Context, assetID string, order int) error {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": "ASSET#" + assetID,
		"SK": fmt.Sprintf("IMAGE#%d", order),
	})
	if err != nil {
		return err
	}

	_, err = r.dbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(core.TableName),
		Key:       key,
	})
	return err
}
