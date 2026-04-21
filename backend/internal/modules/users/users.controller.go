package users

import (
	"net/http"
	"strconv"

	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/users/dtos/req"
	_ "github.com/gianghp123/Vidmerce/backend/internal/modules/users/dtos/res"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/users/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	svc services.UserService
}

func NewUserController(svc services.UserService) *UserController {
	return &UserController{svc: svc}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with email and role
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      req.CreateUserReq  true  "Create user request"
// @Success      201  {object}  response.BaseResponse[res.UserRes]
// @Failure      400  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var body req.CreateUserReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest(err.Error())))
		return
	}

	result, appErr := ctrl.svc.CreateUser(c.Request.Context(), body)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusCreated, response.Success(result))
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Retrieve detailed information about a specific user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.BaseResponse[res.UserRes]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /users/{id} [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("user ID is required")))
		return
	}

	result, appErr := ctrl.svc.GetUser(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// ListUsers godoc
// @Summary      List users
// @Description  Get a paginated list of users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit  query     int     false  "Number of items per page"                Format(int32)
// @Param        cursor query     string  false  "Pagination cursor for next page"
// @Success      200   {object}  response.BaseResponse[[]res.UserRes]
// @Failure      500   {object}  response.BaseResponse[any]
// @Router       /users [get]
func (ctrl *UserController) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	cursor := c.Query("cursor")

	result, appErr := ctrl.svc.ListUsers(c.Request.Context(), limit, cursor)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMeta(result.Data, result.Meta))
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user information
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string               true  "User ID"
// @Param        body body      req.UpdateUserReq   true  "Update user request"
// @Success      200  {object}  response.BaseResponse[res.UserRes]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /users/{id} [patch]
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("user ID is required")))
		return
	}

	var body req.UpdateUserReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest(err.Error())))
		return
	}

	result, appErr := ctrl.svc.UpdateUser(c.Request.Context(), id, body)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.BaseResponse[any]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest("user ID is required")))
		return
	}

	appErr := ctrl.svc.DeleteUser(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil))
}
