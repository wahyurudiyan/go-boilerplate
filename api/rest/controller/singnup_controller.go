package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userDTO "github.com/wahyurudiyan/go-boilerplate/core/dto/user"
	"github.com/wahyurudiyan/go-boilerplate/pkg/common"
)

// SignUp is an controller endpoint that handle request from RESTFul by end-user
// @Summary SignUp user endpoint.
// @Description endpoint that handle user register.
// @Tags User Endpoint
// @Accept */*
// @Param Authorization header string true "Request body type"
// @Param request body userDTO.SignUpDTO true "Request Body"
// @Produce json
// @Success 200 {object} common.RESTBody[any] "Success"
// @Failure 400 {object} common.RESTBody[any] "Bad request"
// @Router /users/signup [POST]
func (b *ControllerBootstrap) SignUp(c *gin.Context) {
	var body userDTO.SignUpDTO
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.RESTErrorResponse[any](1022, "incoming request body invalid"))
		return
	}

	if err := b.UserService.SignUp(c.Request.Context(), body); err != nil {
		c.JSON(http.StatusInternalServerError, common.RESTErrorResponse[any](1034, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, common.RESTSuccessResponse[any]("sign-up success", nil))
}
