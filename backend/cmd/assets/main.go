package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gianghp123/Vidmerce/backend/services/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets"
	"github.com/gianghp123/Vidmerce/backend/services/internal/storage"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Printf("Gin cold start")

	awsCfg := configs.LoadAWSConfig()

	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("Couldn't load default configuration. Have you set up your AWS account?: %v", err)
	}

	dbClient := dynamodb.NewFromConfig(sdkConfig, awsCfg.DynamoDBOptions)

	s3Config := configs.LoadS3Config()
	s3Client := storage.NewS3Storage(s3.NewFromConfig(sdkConfig, awsCfg.S3Options), s3Config.BucketName)

	r := gin.Default()

	api := r.Group("/api")

	assets.RegisterRoutes(api, dbClient, s3Client)

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
