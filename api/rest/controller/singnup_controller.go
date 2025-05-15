package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userDto "github.com/wahyurudiyan/go-boilerplate/core/dto/user"
	"github.com/wahyurudiyan/go-boilerplate/pkg/common"
)

// SampleNotImplement godoc
// @Summary Show not implement code error.
// @Description not implement code error.
// @Tags Test Mot Implemented Endpoint
// @Accept */*
// @Produce json
// @Failure 501 {string} string "Not Implemented"
// @Router /test [GET]
func (b *ControllerBootstrap) SignUp(c *gin.Context) {
	requestBody, err := common.BindJSON[userDto.SignUpDTO](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RESTErrorResponse[any](1022, "incoming request body invalid"))
		return
	}

	if err := b.UserService.SignUp(c.Request.Context(), requestBody.Data); err != nil {
		c.JSON(http.StatusInternalServerError, common.RESTErrorResponse[any](1034, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, common.RESTSuccessResponse[any]("sign-up success", nil))
}
