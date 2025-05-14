package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SampleNotImplement godoc
// @Summary Show not implement code error.
// @Description not implement code error.
// @Tags Test Endpoint
// @Accept */*
// @Produce json
// @Failure 501 {string} string "Not Implemented"
// @Router /test [GET]
func NotImplementController(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "handler not implemented")
}
