package campaigns

import (
	"net/http"
	"strconv"

	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/dtos/req"
	_ "github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/dtos/res"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/services"
	"github.com/gin-gonic/gin"
)

type CampaignController struct {
	svc services.CampaignService
}

func NewCampaignController(svc services.CampaignService) *CampaignController {
	return &CampaignController{svc: svc}
}

// CreateCampaign godoc
// @Summary      Create a new campaign
// @Description  Create a new campaign with an asset and generate a job
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        body  body      req.CreateCampaignReq  true  "Create campaign request"
// @Success      201   {object}  response.BaseResponse[res.CampaignRes]
// @Failure      400   {object}  response.BaseResponse[any]
// @Failure      500   {object}  response.BaseResponse[any]
// @Router       /campaigns [post]
func (ctrl *CampaignController) CreateCampaign(c *gin.Context) {
	var body req.CreateCampaignReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest(err.Error())))
		return
	}

	result, appErr := ctrl.svc.CreateCampaign(c.Request.Context(), body)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusCreated, response.Success(result))
}

// GetCampaign godoc
// @Summary      Get campaign by ID
// @Description  Retrieve detailed information about a specific campaign
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Campaign ID"
// @Success      200  {object}  response.BaseResponse[res.CampaignRes]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /campaigns/{id} [get]
func (ctrl *CampaignController) GetCampaign(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("campaign ID is required")))
		return
	}

	result, appErr := ctrl.svc.GetCampaign(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// ListCampaigns godoc
// @Summary      List campaigns
// @Description  Get a paginated list of campaigns
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        limit  query     int     false  "Number of items per page"                Format(int32)
// @Param        cursor query     string  false  "Pagination cursor for next page"
// @Success      200    {object}  response.BaseResponse[[]res.CampaignRes]
// @Failure      500    {object}  response.BaseResponse[any]
// @Router       /campaigns [get]
func (ctrl *CampaignController) ListCampaigns(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	cursor := c.Query("cursor")

	result, appErr := ctrl.svc.ListCampaigns(c.Request.Context(), limit, cursor)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMeta(result.Data, result.Meta))
}

// UpdateCampaign godoc
// @Summary      Update campaign
// @Description  Update campaign storyboard or slide order
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      string                 true  "Campaign ID"
// @Param        body body      req.UpdateCampaignReq  true  "Update campaign request"
// @Success      200  {object}  response.BaseResponse[res.CampaignRes]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /campaigns/{id} [patch]
func (ctrl *CampaignController) UpdateCampaign(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("campaign ID is required")))
		return
	}

	var body req.UpdateCampaignReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest(err.Error())))
		return
	}

	result, appErr := ctrl.svc.UpdateCampaign(c.Request.Context(), id, body)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// DeleteCampaign godoc
// @Summary      Delete campaign
// @Description  Delete a campaign and its active jobs
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Campaign ID"
// @Success      200  {object}  response.BaseResponse[any]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /campaigns/{id} [delete]
func (ctrl *CampaignController) DeleteCampaign(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("campaign ID is required")))
		return
	}

	appErr := ctrl.svc.DeleteCampaign(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil))
}
