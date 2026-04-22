package webhooks

import (
	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/webhooks/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	clerkCfg := configs.GetClerkConfig()
	svc := services.NewWebhookService(clerkCfg.ClerkWebhookSecret)
	ctrl := NewWebhookController(svc)

	r.POST("/webhooks/clerk", ctrl.HandleClerkWebhook)
}
