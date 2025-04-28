package routes

import (
	"github.com/gofiber/swagger"
	"github.com/wahyurudiyan/go-boilerplate/api/rest/controller"

	_ "github.com/wahyurudiyan/go-boilerplate/docs" // change with your own project docs path

	"github.com/gofiber/fiber/v2"
)

func openAPI(app fiber.Router) {
	app.Get("/swagger/*", swagger.HandlerDefault) // default handler
}

func Routes(app fiber.Router) {
	// Init swagger endpoint
	openAPI(app)

	v1 := app.Group("/api/v1")
	v1.Get("/health-check", controller.HealthCheck)

	productRoutes := v1.Group("/products")
	productRoutes.Get("/all", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
}
