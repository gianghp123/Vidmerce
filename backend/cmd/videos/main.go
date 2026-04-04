package main

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/gianghp123/Vidmerce/backend/cmd/videos/docs" // swagger docs initialization
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/videos"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ginLambda *ginadapter.GinLambdaV2

func setup() (*gin.Engine, *configs.AWSConfig) {

	if _, err := os.Stat(".env"); err == nil {
		log.Println("Found .env file, loading local configurations...")
		_ = godotenv.Load()
	}

	awsCfg := configs.LoadAWSConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	dbClient := dynamodb.NewFromConfig(cfg, awsCfg.DynamoDBOptions)

	r := gin.Default()
	api := r.Group("/api")
	videos.RegisterRoutes(api, dbClient)

	return r, awsCfg
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	router, awsCfg := setup()

	if awsCfg.IsLocal {
		// Swagger is configured via swag annotations in main.go and generated docs
		// The docs.SwaggerInfo is already properly set by swag init
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		// Run as a standard HTTP server locally
		log.Printf("Running in LOCAL SERVER mode on http://localhost:3001")
		if err := router.Run(":3001"); err != nil {
			log.Fatalf("Failed to run local server: %v", err)
		}
	} else {
		// Run as an AWS Lambda function
		log.Printf("Running in LAMBDA mode")
		ginLambda = ginadapter.NewV2(router)
		lambda.Start(Handler)
	}
}

// @title           Swagger Videos API
// @version         1.0
// @description     This is the Vidmerce Videos API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3001
// @BasePath  /api
// @schemes   http
