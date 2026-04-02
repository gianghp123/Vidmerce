package assets

import (
	"net/http"
	"strconv"

	"github.com/gianghp123/Vidmerce/backend/services/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/dtos/req"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/services"
	"github.com/gin-gonic/gin"
)

type AssetController struct {
	svc services.AssetService
}

func NewAssetController(svc services.AssetService) *AssetController {
	return &AssetController{svc: svc}
}

func (ctrl *AssetController) GenerateUploadUrls(c *gin.Context) {
	var body req.GenerateUploadUrlsReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest(err.Error())))
		return
	}

	result, appErr := ctrl.svc.GenerateUploadUrls(c, body.Count)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

func (ctrl *AssetController) CreateAsset(c *gin.Context) {
	var body req.CreateAssetReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest(err.Error())))
		return
	}

	result, appErr := ctrl.svc.CreateAsset(c, body)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusCreated, response.Success(result))
}

func (ctrl *AssetController) ListAssets(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	cursor := c.Query("cursor")

	result, appErr := ctrl.svc.ListAssets(c, limit, cursor)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMeta(result.Data, result.Meta))
}
