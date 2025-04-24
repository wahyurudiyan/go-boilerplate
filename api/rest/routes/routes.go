package routes

import (
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/gofiber/fiber/v2"
	_ "github.com/swaggo/fiber-swagger/example/docs"
)

func swagger(app fiber.Router) {
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
}

func Routes(app fiber.Router) {
	// Init swagger endpoint
	swagger(app)

	v1 := app.Group("/api/v1")
	productRoutes := v1.Group("/products")
	productRoutes.Get("/all", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
}
