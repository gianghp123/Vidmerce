package auth

import (
	"net/http"

	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/auth/dtos/req"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/auth/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	svc *services.AuthService
}

func NewAuthController(svc *services.AuthService) *AuthController {
	return &AuthController{svc: svc}
}

// SignUp godoc
// @Summary      Sign up a new user
// @Description  Register a new user with email and password using Clerk
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      req.SignUpReq  true  "Sign up request"
// @Success      201  {object}  response.BaseResponse[any]
// @Failure      400  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router       /auth/signup [post]
func (ctrl *AuthController) SignUp(c *gin.Context) {
	var body req.SignUpReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest(err.Error())))
		return
	}

	result, appErr := ctrl.svc.SignUp(c.Request.Context(), body)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"userId": result.ID,
		"email":  result.EmailAddresses,
	}))
}
