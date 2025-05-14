package rest

import (
	"context"

	"github.com/wahyurudiyan/go-boilerplate/config"
)

type IHttpServer interface {
	Listen() error
	Shutdown(ctx context.Context) error
}

func NewGinServer(cfg *config.ServiceConfig) IHttpServer {
	ginServer := &httpGinServer{
		cfg: cfg,
	}
	ginServer.router = ginServer.newGinServer()

	return ginServer
}
