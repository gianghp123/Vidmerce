package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	_ "github.com/gianghp123/Vidmerce/backend/cmd/assets/docs"
	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/middlewares"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

var ginLambda *ginadapter.GinLambdaV2

func setup() (*gin.Engine, *configs.AWSConfig, *core.Clients) {
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Found .env file, loading local configurations...")
		_ = godotenv.Load()
	}

	logger := configs.InitLogger()
	logger.Info("Initializing application", zap.String("mode", "startup"))

	awsCfg := configs.LoadAWSConfig()
	s3Cfg := configs.LoadS3Config()

	awsClients, err := core.NewClients(context.TODO(), awsCfg, s3Cfg.BucketName)
	if err != nil {
		log.Fatalf("Failed to init AWS clients: %v", err)
	}

	r := gin.Default()

	// 1. Apply CORS first
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "apikey"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger (Public)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 2. Apply Auth Middleware
	// In Lambda, auth is usually handled by API Gateway, so we check IsLocal
	if awsCfg.IsLocal {
		clerkCfg := configs.GetClerkConfig()
		if clerkCfg.ClerkSecret != "" {
			r.Use(middlewares.ClerkAuth())
		}
	}

	// 3. Register Routes AFTER middleware
	api := r.Group("/api")
	assets.RegisterRoutes(api, awsClients.DB, awsClients.S3)

	return r, awsCfg, awsClients
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

// @title   	Assets API
// @version  	1.0
// @description This API handles requests for Assets
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api
// @schemes   http

// @query.collection.format multi

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	router, awsCfg, _ := setup()

	if awsCfg.IsLocal {
		log.Printf("Running in LOCAL SERVER mode on http://localhost:3000")
		if err := router.Run(":3000"); err != nil {
			log.Fatalf("Failed to run local server: %v", err)
		}
	} else {
		log.Printf("Running in LAMBDA mode")
		ginLambda = ginadapter.NewV2(router)
		lambda.Start(Handler)
	}
}
