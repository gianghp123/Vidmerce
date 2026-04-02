package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Printf("Gin cold start")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	dbClient := dynamodb.NewFromConfig(cfg)

	r := gin.Default()

	api := r.Group("/api")

	videos.RegisterRoutes(api, dbClient)

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
