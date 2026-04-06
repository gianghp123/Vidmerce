package assets

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets/services"
	"github.com/gianghp123/Vidmerce/backend/internal/storage"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, dbClient *dynamodb.Client, store storage.Storage) {
	assetRepo := repositories.NewAssetRepository(dbClient)
	imageRepo := repositories.NewImageRepository(dbClient)
	svc := services.NewAssetService(assetRepo, imageRepo, store)
	ctrl := NewAssetController(svc)

	group := r.Group("/assets")
	{
		group.POST("", ctrl.CreateAsset)
		group.POST("/:id/confirm", ctrl.ConfirmUpload)
		group.GET("/:id", ctrl.GetAsset)
		group.GET("", ctrl.ListAssets)

		group.GET("/:id/images/upload-url", ctrl.GetImageUploadUrl)
		group.DELETE("/:id/images/:imageId", ctrl.DeleteAssetImage)
	}
}
