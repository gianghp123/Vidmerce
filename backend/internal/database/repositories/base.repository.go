package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
)

type BaseRepository interface {
	TransactWriteItems(ctx context.Context, items ...interface{}) error
	GetItem(ctx context.Context, pk, sk string) (map[string]types.AttributeValue, error)
	PutItem(ctx context.Context, item interface{}) error
	DeleteItem(ctx context.Context, pk, sk string) error
}

type baseRepository struct {
	dbClient *dynamodb.Client
}

func NewBaseRepository(dbClient *dynamodb.Client) BaseRepository {
	return &baseRepository{dbClient: dbClient}
}

func (r *baseRepository) TransactWriteItems(ctx context.Context, items ...interface{}) error {
	if len(items) == 0 {
		return nil
	}

	transactItems := make([]types.TransactWriteItem, len(items))
	for i, item := range items {
		itemMap, err := attributevalue.MarshalMap(item)
		if err != nil {
			return err
		}
		transactItems[i] = types.TransactWriteItem{
			Put: &types.Put{
				TableName: aws.String(core.TableName),
				Item:      itemMap,
			},
		}
	}

	_, err := r.dbClient.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	return err
}

func (r *baseRepository) GetItem(ctx context.Context, pk, sk string) (map[string]types.AttributeValue, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": pk,
		"SK": sk,
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

	return resp.Item, nil
}

func (r *baseRepository) PutItem(ctx context.Context, item interface{}) error {
	itemMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = r.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(core.TableName),
		Item:      itemMap,
	})
	return err
}

func (r *baseRepository) DeleteItem(ctx context.Context, pk, sk string) error {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": pk,
		"SK": sk,
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
