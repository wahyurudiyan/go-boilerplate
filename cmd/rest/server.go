package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/wahyurudiyan/go-boilerplate/config"
)

type IHttpServer interface {
	Run() error
	Shutdown(ctx context.Context) error
	RegisterRoutes(routesFn func(*gin.Engine))
}

func NewGinServer(cfg *config.ServiceConfig) IHttpServer {
	return newGinServer(cfg)
}
