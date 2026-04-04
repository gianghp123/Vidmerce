package videos

import (
	"net/http"
	"strconv"

	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/videos/dtos/req"
	_ "github.com/gianghp123/Vidmerce/backend/internal/modules/videos/dtos/res"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/videos/services"
	"github.com/gin-gonic/gin"
)

type VideoController struct {
	svc services.VideoService
}

func NewVideoController(svc services.VideoService) *VideoController {
	return &VideoController{svc: svc}
}

// CreateVideo godoc
// @Summary      Create a new video
// @Description  Create a video by combining multiple assets with transitions
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        body  body      req.CreateVideoReq  true  "Create video request"
// @Success      201   {object}  response.BaseResponse[res.CreateVideoRes]
// @Failure      400   {object}  response.BaseResponse[any]
// @Failure      500   {object}  response.BaseResponse[any]
// @Router       /videos [post]
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

// ListVideos godoc
// @Summary      List videos
// @Description  Get a paginated list of videos
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        limit   query     int     false  "Number of items per page"                Format(int32)
// @Param        last_key query    string  false  "Pagination cursor (last video ID) for next page"
// @Success      200     {object}  response.BaseResponse[[]res.VideoResponse]
// @Failure      400     {object}  response.BaseResponse[any]
// @Failure      500     {object}  response.BaseResponse[any]
// @Router       /videos [get]
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

// GetVideo godoc
// @Summary      Get video by ID
// @Description  Retrieve detailed information about a specific video including hotspots
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Video ID"
// @Success      200  {object}  response.BaseResponse[res.VideoDetailResponse]
// @Failure      400  {object}  response.BaseResponse[any]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /videos/{id} [get]
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
