package repositories

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/internal/utils"
)

type UserRepository interface {
	FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.UserEntity], error)
	FindByID(ctx context.Context, id string) (*models.UserEntity, error)
	Create(ctx context.Context, user models.UserEntity) error
	Update(ctx context.Context, user models.UserEntity) error
	Delete(ctx context.Context, id string) error
	DBClient() *dynamodb.Client
}

type userRepository struct {
	dbClient *dynamodb.Client
}

func NewUserRepository(dbClient *dynamodb.Client) UserRepository {
	return &userRepository{dbClient: dbClient}
}

func (r *userRepository) DBClient() *dynamodb.Client {
	return r.dbClient
}

func (r *userRepository) FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.UserEntity], error) {
	exclusiveStartKey, err := core.DecodeCursor(lastKey)
	if err != nil {
		return nil, err
	}

	keyCond := expression.Key("Gsi1Pk").Equal(expression.Value(string(core.Gsi1PkEntityUser)))
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

	var items []models.UserEntity
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &items); err != nil {
		return nil, err
	}

	for i := range items {
		items[i].Pk = strings.TrimPrefix(items[i].Pk, string(core.EntityTypeUser)+core.KeySeparator)
	}

	nextCursor, err := core.EncodeCursor(resp.LastEvaluatedKey)
	if err != nil {
		return nil, err
	}

	hasMore := resp.LastEvaluatedKey != nil
	return &response.PaginatedResult[models.UserEntity]{
		Data: items,
		Meta: response.NewCursorMeta(limit, nextCursor, hasMore),
	}, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.UserEntity, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"Pk": utils.BuildPk(core.EntityTypeUser, id),
		"Sk": string(core.SortKeyMetadata),
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

	var item models.UserEntity
	if err := attributevalue.UnmarshalMap(resp.Item, &item); err != nil {
		return nil, err
	}
	item.Pk = strings.TrimPrefix(item.Pk, string(core.EntityTypeUser)+core.KeySeparator)
	return &item, nil
}

func (r *userRepository) Create(ctx context.Context, user models.UserEntity) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(core.TableName),
		Item:      item,
	})
	return err
}

func (r *userRepository) Update(ctx context.Context, user models.UserEntity) error {
	key, err := attributevalue.MarshalMap(map[string]string{
		"Pk": utils.BuildPk(core.EntityTypeUser, user.Pk),
		"Sk": string(core.SortKeyMetadata),
	})
	if err != nil {
		return err
	}

	updateExpr := "SET "
	exprNames := map[string]string{}
	exprValues := map[string]types.AttributeValue{}

	if user.Role != "" {
		updateExpr += "#role = :role, "
		exprNames["#role"] = "role"
		exprValues[":role"] = &types.AttributeValueMemberS{Value: string(user.Role)}
	}

	updateExpr += "#updatedAt = :updatedAt"
	exprNames["#updatedAt"] = "updatedAt"
	exprValues[":updatedAt"] = &types.AttributeValueMemberS{Value: utils.Now()}

	_, err = r.dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(core.TableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprNames,
		ExpressionAttributeValues: exprValues,
	})
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	key, err := attributevalue.MarshalMap(map[string]string{
		"Pk": utils.BuildPk(core.EntityTypeUser, id),
		"Sk": string(core.SortKeyMetadata),
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
