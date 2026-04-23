package repositories

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
)

const MaxImagesPerAsset = 5

var ErrMaxImagesReached = errors.New("maximum images per asset reached")

type AssetRepository interface {
	FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.AssetEntity], error)
	FindByID(ctx context.Context, id string) (*models.AssetEntity, error)
	Create(ctx context.Context, asset models.AssetEntity) error
	IncrementImageCount(ctx context.Context, id string) (int, error)
	CreateWithImages(ctx context.Context, asset models.AssetEntity, images []models.ImageEntity) error
	CreateWithJob(ctx context.Context, asset models.AssetEntity, job models.JobEntity) error
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

	keyCond := expression.Key("Gsi1Pk").Equal(expression.Value(string(core.Gsi1PkEntityAsset)))
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

	for i := range items {
		items[i].Pk = strings.TrimPrefix(items[i].Pk, string(core.SkPrefixAsset)+core.KeySeparator)
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
	keyCond := expression.Key("Gsi1Pk").Equal(expression.Value(string(core.Gsi1PkEntityAsset)))
	keyCond = keyCond.And(expression.Key("Gsi1Sk").Equal(expression.Value(id)))

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
		Limit:                     aws.Int32(1),
	}

	resp, err := r.dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, nil
	}

	var item models.AssetEntity
	if err := attributevalue.UnmarshalMap(resp.Items[0], &item); err != nil {
		return nil, err
	}

	item.Pk = strings.TrimPrefix(item.Pk, string(core.SkPrefixAsset)+core.KeySeparator)
	return &item, nil
}

func (r *assetRepository) Create(ctx context.Context, asset models.AssetEntity) error {
	item, err := attributevalue.MarshalMap(asset)
	if err != nil {
		return err
	}

	_, err = r.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(core.TableName),
		Item:      item,
	})
	return err
}

func (r *assetRepository) IncrementImageCount(ctx context.Context, id string) (int, error) {
	item, err := attributevalue.MarshalMap(models.BaseItem{
		Pk: string(core.PkPrefixAsset) + core.KeySeparator + id,
		Sk: string(core.SkPrefixAsset) + core.KeySeparator + id,
	})
	if err != nil {
		return 0, err
	}

	maxVal := strconv.Itoa(MaxImagesPerAsset)

	resp, err := r.dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:           aws.String(core.TableName),
		Key:                 item,
		UpdateExpression:    aws.String("SET #c = #c + :inc"),
		ConditionExpression: aws.String("#c < :max"),
		ExpressionAttributeNames: map[string]string{
			"#c": "imageCount",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
			":max": &types.AttributeValueMemberN{Value: maxVal},
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		var conditionalErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalErr) {
			return 0, ErrMaxImagesReached
		}
		return 0, err
	}

	var result struct {
		ImageCount int `dynamodbav:"imageCount"`
	}
	if err := attributevalue.UnmarshalMap(resp.Attributes, &result); err != nil {
		return 0, err
	}
	return result.ImageCount, nil
}

func (r *assetRepository) CreateWithImages(ctx context.Context, asset models.AssetEntity, images []models.ImageEntity) error {
	items := make([]types.TransactWriteItem, 0, 1+len(images))

	assetMap, err := attributevalue.MarshalMap(asset)
	if err != nil {
		return err
	}
	items = append(items, types.TransactWriteItem{Put: &types.Put{TableName: aws.String(core.TableName), Item: assetMap}})

	for _, img := range images {
		imgMap, err := attributevalue.MarshalMap(img)
		if err != nil {
			return err
		}
		items = append(items, types.TransactWriteItem{Put: &types.Put{TableName: aws.String(core.TableName), Item: imgMap}})
	}

	_, err = r.dbClient.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{TransactItems: items})
	return err
}

func (r *assetRepository) CreateWithJob(ctx context.Context, asset models.AssetEntity, job models.JobEntity) error {
	assetMap, err := attributevalue.MarshalMap(asset)
	if err != nil {
		return err
	}
	jobMap, err := attributevalue.MarshalMap(job)
	if err != nil {
		return err
	}

	_, err = r.dbClient.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{Put: &types.Put{TableName: aws.String(core.TableName), Item: assetMap}},
			{Put: &types.Put{TableName: aws.String(core.TableName), Item: jobMap}},
		},
	})
	return err
}
