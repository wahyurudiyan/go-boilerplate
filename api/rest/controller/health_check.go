package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Show the status of server.
// @Description the health check endpoint provide the status of server.
// @Tags Health Check Endpoint
// @Accept */*
// @Produce json
// @Success 200 {object} object
// @Router /health-check [GET]
func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}
	c.JSON(http.StatusOK, res)

	return
}
