package controller

import "github.com/gofiber/fiber/v2"

// HealthCheck godoc
// @Summary Show the status of server.
// @Description the health check endpoint provide the status of server.
// @Tags Health Check Endpoint
// @Accept */*
// @Produce json
// @Success 200 {object} object
// @Router /health-check [GET]
func HealthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}
