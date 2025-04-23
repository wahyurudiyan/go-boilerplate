package routes

import "github.com/gofiber/fiber/v2"

func Routes(app fiber.Router) {
	v1 := app.Group("/api/v1")

	productRoutes := v1.Group("/products")
	productRoutes.Get("/all", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
}
