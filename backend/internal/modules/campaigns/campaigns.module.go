package campaigns

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, dbClient *dynamodb.Client) {
	repo := repositories.NewCampaignRepository(dbClient)
	svc := services.NewCampaignService(repo)
	ctrl := NewCampaignController(svc)

	r.POST("", ctrl.CreateCampaign)
	r.GET("/:id", ctrl.GetCampaign)
	r.GET("", ctrl.ListCampaigns)
	r.PATCH("/:id", ctrl.UpdateCampaign)
	r.DELETE("/:id", ctrl.DeleteCampaign)
}
