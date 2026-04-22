package webhooks

import (
	"io"
	"net/http"

	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/webhooks/services"
	"github.com/gin-gonic/gin"
)

type WebhookController struct {
	svc services.WebhookService
}

func NewWebhookController(svc services.WebhookService) *WebhookController {
	return &WebhookController{svc: svc}
}

// HandleClerkWebhook godoc
// @Summary      Handle Clerk webhook
// @Description  Process Clerk webhooks and update user metadata
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        body  body      []byte  true  "Clerk webhook payload"
// @Success      200  {object}  response.BaseResponse[any]
// @Failure      400  {object}  response.BaseResponse[any]
// @Failure      401  {object}  response.BaseResponse[any]
// @Router       /webhooks/clerk [post]
func (ctrl *WebhookController) HandleClerkWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("failed to read payload")))
		return
	}

	appErr := ctrl.svc.HandleClerkWebhook(c.Request.Context(), payload, c.Request.Header)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"message": "Webhook processed successfully",
	}))
}
