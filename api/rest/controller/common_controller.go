package controller

import (
	"github.com/gofiber/fiber/v2"
)

// SampleNotImplement godoc
// @Summary Show not implement code error.
// @Description not implement code error.
// @Tags Test Endpoint
// @Accept */*
// @Produce json
// @Failure 501 {string} string "Not Implemented"
// @Router /test [GET]
func NotImplementController(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
