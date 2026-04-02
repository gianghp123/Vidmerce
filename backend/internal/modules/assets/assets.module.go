package assets

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/repositories"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/services"
	"github.com/gianghp123/Vidmerce/backend/services/internal/storage"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, dbClient *dynamodb.Client, store storage.Storage) {
	repo := repositories.NewAssetRepository(dbClient)
	svc := services.NewAssetService(repo, store)
	ctrl := NewAssetController(svc)

	group := r.Group("/assets")
	{
		group.POST("/upload-urls", ctrl.GenerateUploadUrls)
		group.POST("", ctrl.CreateAsset)
		group.GET("", ctrl.ListAssets)
	}
}
