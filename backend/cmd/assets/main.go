package main

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/gianghp123/Vidmerce/backend/cmd/assets/docs" // swagger docs initialization
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns"
	"github.com/gianghp123/Vidmerce/backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var ginLambda *ginadapter.GinLambdaV2

func setup() (*gin.Engine, *configs.AWSConfig) {
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Found .env file, loading local configurations...")
		_ = godotenv.Load()
	}

	logger := configs.InitLogger()
	logger.Info("Initializing application", zap.String("mode", "startup"))

	awsCfg := configs.LoadAWSConfig()
	s3Cfg := configs.LoadS3Config()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS SDK config: %v", err)
	}

	dbClient := dynamodb.NewFromConfig(sdkConfig)
	s3Client := storage.NewS3Storage(
		s3.NewFromConfig(sdkConfig),
		s3Cfg.BucketName,
	)

	r := gin.Default()
	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Status(200)
	})
	api := r.Group("/api")
	assets.RegisterRoutes(api, dbClient, s3Client)
	campaigns.RegisterRoutes(api, dbClient)

	return r, awsCfg
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Full request", req)
	return ginLambda.ProxyWithContext(ctx, req)
}

// @title           Swagger Assets API
// @version         1.0
// @description     This is the Vidmerce Assets API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api
// @schemes   http
func main() {
	router, awsCfg := setup()

	if awsCfg.IsLocal {
		// Swagger is configured via swag annotations in main.go and generated docs
		// The docs.SwaggerInfo is already properly set by swag init
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, // hoặc domain cụ thể
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
		// Run as a standard HTTP server locally
		log.Printf("Running in LOCAL SERVER mode on http://localhost:3000")
		if err := router.Run(":3000"); err != nil {
			log.Fatalf("Failed to run local server: %v", err)
		}
	} else {
		// Run as an AWS Lambda function
		log.Printf("Running in LAMBDA mode")
		ginLambda = ginadapter.NewV2(router)
		lambda.Start(Handler)
	}
}
