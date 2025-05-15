package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wahyurudiyan/go-boilerplate/api/rest/controller"
	"github.com/wahyurudiyan/go-boilerplate/docs" // change with your own project docs path
	// _ "github.com/wahyurudiyan/go-boilerplate/docs" // change with your own project docs path
)

type routerBootstrap struct {
	controller *controller.ControllerBootstrap
}

func NewRouter(c *controller.ControllerBootstrap) *routerBootstrap {
	return &routerBootstrap{
		controller: c,
	}
}

func (r *routerBootstrap) swaggerAPIDoc(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // default handler
}

func (r *routerBootstrap) Routes(router *gin.Engine) {
	// Init base path
	rootPathV1 := router.Group("/api/v1")
	rootPathV1.GET("/health-check", controller.HealthCheck)

	// Init swagger endpoint
	docs.SwaggerInfo.BasePath = rootPathV1.BasePath()
	r.swaggerAPIDoc(router)

	userRoutes := rootPathV1.Group("/users")
	userRoutes.POST("/signup", r.controller.SignUp)
}
