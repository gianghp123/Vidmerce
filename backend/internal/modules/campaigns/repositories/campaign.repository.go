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
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/internal/utils"
)

type CampaignRepository interface {
	FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.CampaignEntity], error)
	FindByID(ctx context.Context, id string) (*models.CampaignEntity, error)
	Create(ctx context.Context, campaign models.CampaignEntity) error
	Update(ctx context.Context, campaign models.CampaignEntity) error
	Delete(ctx context.Context, id string) error
	FindActiveJobsByTargetID(ctx context.Context, targetID string) ([]models.JobEntity, error)
	DBClient() *dynamodb.Client
}

type campaignRepository struct {
	dbClient *dynamodb.Client
}

func NewCampaignRepository(dbClient *dynamodb.Client) CampaignRepository {
	return &campaignRepository{dbClient: dbClient}
}

func (r *campaignRepository) DBClient() *dynamodb.Client {
	return r.dbClient
}

func (r *campaignRepository) FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.CampaignEntity], error) {
	exclusiveStartKey, err := core.DecodeCursor(lastKey)
	if err != nil {
		return nil, err
	}

	keyCond := expression.Key("GSI1PK").Equal(expression.Value(string(enums.GSI1PKEntityCampaign)))
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

	var items []models.CampaignEntity
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &items); err != nil {
		return nil, err
	}

	for i := range items {
		items[i].PK = strings.TrimPrefix(items[i].PK, string(enums.EntityTypeCampaign)+enums.KeySeparator)
	}

	nextCursor, err := core.EncodeCursor(resp.LastEvaluatedKey)
	if err != nil {
		return nil, err
	}

	hasMore := resp.LastEvaluatedKey != nil
	return &response.PaginatedResult[models.CampaignEntity]{
		Data: items,
		Meta: response.NewCursorMeta(limit, nextCursor, hasMore),
	}, nil
}

func (r *campaignRepository) FindByID(ctx context.Context, id string) (*models.CampaignEntity, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": utils.BuildPK(enums.EntityTypeCampaign, id),
		"SK": string(enums.SortKeyMetadata),
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

	var item models.CampaignEntity
	if err := attributevalue.UnmarshalMap(resp.Item, &item); err != nil {
		return nil, err
	}
	item.PK = strings.TrimPrefix(item.PK, string(enums.EntityTypeCampaign)+enums.KeySeparator)
	return &item, nil
}

func (r *campaignRepository) Create(ctx context.Context, campaign models.CampaignEntity) error {
	item, err := attributevalue.MarshalMap(campaign)
	if err != nil {
		return err
	}

	_, err = r.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(core.TableName),
		Item:      item,
	})
	return err
}

func (r *campaignRepository) Update(ctx context.Context, campaign models.CampaignEntity) error {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": utils.BuildPK(enums.EntityTypeCampaign, campaign.PK),
		"SK": string(enums.SortKeyMetadata),
	})
	if err != nil {
		return err
	}

	updateExpr := "SET "
	exprNames := map[string]string{}
	exprValues := map[string]types.AttributeValue{}

	if campaign.Storyboard.Headline != "" {
		updateExpr += "#storyboard.#headline = :headline, "
		exprNames["#storyboard"] = "storyboard"
		exprNames["#headline"] = "headline"
		exprValues[":headline"] = &types.AttributeValueMemberS{Value: campaign.Storyboard.Headline}
	}
	if campaign.Storyboard.BodyCopy != "" {
		updateExpr += "#storyboard.#body = :body, "
		exprNames["#body"] = "bodyCopy"
		exprValues[":body"] = &types.AttributeValueMemberS{Value: campaign.Storyboard.BodyCopy}
	}
	if campaign.Storyboard.CTA != "" {
		updateExpr += "#storyboard.#cta = :cta, "
		exprNames["#cta"] = "cta"
		exprValues[":cta"] = &types.AttributeValueMemberS{Value: campaign.Storyboard.CTA}
	}
	if campaign.Storyboard.Slides != nil {
		updateExpr += "#storyboard.#slides = :slides, "
		exprNames["#slides"] = "slides"
		slides, _ := attributevalue.MarshalList(campaign.Storyboard.Slides)
		exprValues[":slides"] = &types.AttributeValueMemberL{Value: slides}
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

func (r *campaignRepository) Delete(ctx context.Context, id string) error {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": utils.BuildPK(enums.EntityTypeCampaign, id),
		"SK": string(enums.SortKeyMetadata),
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

func (r *campaignRepository) FindActiveJobsByTargetID(ctx context.Context, targetID string) ([]models.JobEntity, error) {
	keyCond := expression.Key("TargetID").Equal(expression.Value(targetID))
	filterCond := expression.Name("Status").Equal(expression.Value(string(enums.StatusJobPending))).
		Or(expression.Name("Status").Equal(expression.Value(string(enums.StatusJobProcessing))))
	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
		WithFilter(filterCond).
		Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(core.TableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	resp, err := r.dbClient.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	var items []models.JobEntity
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &items); err != nil {
		return nil, err
	}

	return items, nil
}
