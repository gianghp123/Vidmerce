package videos

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/repositories"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, dbClient *dynamodb.Client) {
	repo := repositories.NewVideoRepository(dbClient)
	svc := services.NewVideoService(repo)
	ctrl := NewVideoController(svc)

	group := r.Group("/videos")
	{
		group.POST("", ctrl.CreateVideo)
		group.GET("", ctrl.ListVideos)
		group.GET("/:id", ctrl.GetVideo)
	}
}
