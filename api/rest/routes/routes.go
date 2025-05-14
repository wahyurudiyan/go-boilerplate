package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wahyurudiyan/go-boilerplate/api/rest/controller"
	"github.com/wahyurudiyan/go-boilerplate/docs" // change with your own project docs path
	// _ "github.com/wahyurudiyan/go-boilerplate/docs" // change with your own project docs path
)

func swaggerAPIDoc(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // default handler
}

func Routes(router *gin.Engine) {
	// Init base path
	rootPathV1 := router.Group("/api/v1")
	rootPathV1.GET("/health-check", controller.HealthCheck)
	rootPathV1.GET("/test", controller.NotImplementController)

	// Init swagger endpoint
	docs.SwaggerInfo.BasePath = rootPathV1.BasePath()
	swaggerAPIDoc(router)

	productRoutes := rootPathV1.Group("/products")
	productRoutes.GET("/all", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
}
