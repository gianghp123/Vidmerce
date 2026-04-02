package videos

import (
	"net/http"
	"strconv"

	"github.com/gianghp123/Vidmerce/backend/services/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/dtos/req"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/services"
	"github.com/gin-gonic/gin"
)

type VideoController struct {
	svc services.VideoService
}

func NewVideoController(svc services.VideoService) *VideoController {
	return &VideoController{svc: svc}
}

func (ctrl *VideoController) CreateVideo(c *gin.Context) {
	var body req.CreateVideoReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest()))
		return
	}

	result, appErr := ctrl.svc.CreateVideo(c, body)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusCreated, response.Success(result))
}

func (ctrl *VideoController) ListVideos(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	cursor := c.Query("cursor")

	result, appErr := ctrl.svc.ListVideos(c, limit, cursor)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMeta(result.Data, result.Meta))
}

func (ctrl *VideoController) GetVideo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest()))
		return
	}

	result, appErr := ctrl.svc.GetVideo(c, id)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}
