package assets

import (
	"net/http"
	"strconv"

	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets/dtos/req"
	_ "github.com/gianghp123/Vidmerce/backend/internal/modules/assets/dtos/res"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets/services"
	"github.com/gianghp123/Vidmerce/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type AssetController struct {
	svc services.AssetService
}

func NewAssetController(svc services.AssetService) *AssetController {
	return &AssetController{svc: svc}
}

// CreateAsset godoc
// @Summary      Create a new asset
// @Description  Create a new asset with name, price, product URL and image count
// @Security Bearer
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        body  body      req.CreateAssetReq  true  "Create asset request"
// @Success      201   {object}  response.BaseResponse[res.CreateAssetRes]
// @Failure      400   {object}  response.BaseResponse[any]
// @Failure      500   {object}  response.BaseResponse[any]
// @Router       /assets [post]
func (ctrl *AssetController) CreateAsset(c *gin.Context) {
	utils.LogRequestHeaders(c, configs.GetLogger())
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

// ConfirmUpload godoc
// @Summary      Confirm asset upload
// @Description  Confirm upload for an asset and generate image URLs
// @Security Bearer
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Asset ID"
// @Success      200  {object}  response.BaseResponse[res.ConfirmAssetRes]
// @Failure      400  {object}  response.BaseResponse[any]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /assets/{id}/confirm [post]
func (ctrl *AssetController) ConfirmUpload(c *gin.Context) {
	utils.LogRequestHeaders(c, configs.GetLogger())
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

// ListAssets godoc
// @Summary      List assets
// @Description  Get a paginated list of assets
// @Security Bearer
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        limit  query     int     false  "Number of items per page"                Format(int32)
// @Param        cursor query     string  false  "Pagination cursor for next page"
// @Success      200    {object}  response.BaseResponse[[]res.AssetRes]
// @Failure      400    {object}  response.BaseResponse[any]
// @Failure      500    {object}  response.BaseResponse[any]
// @Router       /assets [get]
func (ctrl *AssetController) ListAssets(c *gin.Context) {
	utils.LogRequestHeaders(c, configs.GetLogger())
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

// GetAsset godoc
// @Summary      Get asset by ID
// @Description  Retrieve detailed information about a specific asset
// @Security Bearer
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Asset ID"
// @Success      200  {object}  response.BaseResponse[res.AssetRes]
// @Failure      400  {object}  response.BaseResponse[any]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /assets/{id} [get]
func (ctrl *AssetController) GetAsset(c *gin.Context) {
	utils.LogRequestHeaders(c, configs.GetLogger())
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

// GetImageUploadUrl godoc
// @Summary      Get image upload URL
// @Description  Generate a presigned URL for uploading an image to an asset
// @Security Bearer
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        id         path      string  true  "Asset ID"
// @Param        fileName   query     string  true  "File name"
// @Success      200  {object}  response.BaseResponse[res.AssetUpload]
// @Failure      400  {object}  response.BaseResponse[any]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /assets/{id}/images/upload-url [get]
func (ctrl *AssetController) GetImageUploadUrl(c *gin.Context) {
	utils.LogRequestHeaders(c, configs.GetLogger())
	assetID := c.Param("id")
	if assetID == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("asset ID is required")))
		return
	}

	var query req.GetImageUploadUrlReq
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest(err.Error())))
		return
	}

	result, appErr := ctrl.svc.GetImageUploadUrl(c.Request.Context(), assetID, query.FileName)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// DeleteAssetImage godoc
// @Summary      Delete asset image
// @Description  Delete an image from an asset
// @Security Bearer
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Asset ID"
// @Param        imageId  path      string  true  "Image ID"
// @Success      200  {object}  response.BaseResponse[any]
// @Failure      400  {object}  response.BaseResponse[any]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /assets/{id}/images/{imageId} [delete]
func (ctrl *AssetController) DeleteAssetImage(c *gin.Context) {
	utils.LogRequestHeaders(c, configs.GetLogger())
	assetID := c.Param("id")
	imageID := c.Param("imageId")

	if assetID == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("asset ID is required")))
		return
	}
	if imageID == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("image ID is required")))
		return
	}

	appErr := ctrl.svc.DeleteAssetImage(c.Request.Context(), assetID, imageID)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil))
}

// ImportAsset godoc
// @Summary      Import asset from product URL
// @Description  Create an asset with IMPORTING status and a SCRAPE_PRODUCT job
// @Security Bearer
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        body  body      req.ImportAssetReq  true  "Import asset request"
// @Success      201   {object}  response.BaseResponse[res.ImportAssetRes]
// @Failure      400   {object}  response.BaseResponse[any]
// @Failure      500   {object}  response.BaseResponse[any]
// @Router       /assets/import [post]
func (ctrl *AssetController) ImportAsset(c *gin.Context) {
	utils.LogRequestHeaders(c, configs.GetLogger())
	var body req.ImportAssetReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest(err.Error())))
		return
	}

	result, appErr := ctrl.svc.ImportAsset(c.Request.Context(), body)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusCreated, response.Success(result))
}
