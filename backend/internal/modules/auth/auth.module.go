package auth

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/auth/services"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/users/repositories"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, dbClient *dynamodb.Client) {
	userRepo := repositories.NewUserRepository(dbClient)
	svc := services.NewAuthService(userRepo)
	ctrl := NewAuthController(svc)

	r.POST("/signup", ctrl.SignUp)
}
