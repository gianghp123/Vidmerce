package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/auth"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS SDK config: %v", err)
	}

	dbClient := dynamodb.NewFromConfig(sdkConfig, awsCfg.DynamoDBOptions)

	clerkConfig := configs.GetClerkConfig()
	clerk.SetKey(clerkConfig.ClerkSecret)

	r := gin.Default()
	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Status(200)
	})
	api := r.Group("/api")
	auth.RegisterRoutes(api, dbClient)

	return r, awsCfg
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Full request", req)
	return ginLambda.ProxyWithContext(ctx, req)
}

// @title           Swagger Auth API
// @version         1.0
// @description     This is the Vidmerce Auth API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3002
// @BasePath  /api
// @schemes   http
func main() {
	router, awsCfg := setup()

	if awsCfg.IsLocal {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
		log.Printf("Running in LOCAL SERVER mode on http://localhost:3002")
		if err := router.Run(":3002"); err != nil {
			log.Fatalf("Failed to run local server: %v", err)
		}
	} else {
		log.Printf("Running in LAMBDA mode")
		ginLambda = ginadapter.NewV2(router)
		lambda.Start(Handler)
	}
}
