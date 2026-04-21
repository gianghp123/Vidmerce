package users

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/users/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/users/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, dbClient *dynamodb.Client) {
	repo := repositories.NewUserRepository(dbClient)
	svc := services.NewUserService(repo)
	ctrl := NewUserController(svc)

	r.POST("", ctrl.CreateUser)
	r.GET("/:id", ctrl.GetUser)
	r.GET("", ctrl.ListUsers)
	r.PATCH("/:id", ctrl.UpdateUser)
	r.DELETE("/:id", ctrl.DeleteUser)
}
