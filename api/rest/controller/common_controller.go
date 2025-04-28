package controller

import "github.com/gofiber/fiber/v2"

// SampleNotImplement godoc
// @Summary Show not implement code error.
// @Description not implement code error.
// @Tags Test Endpoint
// @Accept */*
// @Produce json
// @Success 501 {object} map[string]interface{}
// @Router /test [GET]
func NotImplementController(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
