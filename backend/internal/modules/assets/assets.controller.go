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

func (ctrl *AssetController) CreateAsset(c *gin.Context) {
	var body req.CreateAssetReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest(err.Error())))
		return
	}

	result, appErr := ctrl.svc.CreateAsset(c.Request.Context(), body)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusCreated, response.Success(result))
}

func (ctrl *AssetController) ConfirmUpload(c *gin.Context) {
	assetID := c.Param("id")
	if assetID == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("asset ID is required")))
		return
	}

	result, appErr := ctrl.svc.ConfirmUpload(c.Request.Context(), assetID)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

func (ctrl *AssetController) ListAssets(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	cursor := c.Query("cursor")

	result, appErr := ctrl.svc.ListAssets(c.Request.Context(), limit, cursor)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMeta(result.Data, result.Meta))
}

func (ctrl *AssetController) GetAsset(c *gin.Context) {
	assetID := c.Param("id")
	if assetID == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("asset ID is required")))
		return
	}

	result, appErr := ctrl.svc.GetAsset(c.Request.Context(), assetID)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}
