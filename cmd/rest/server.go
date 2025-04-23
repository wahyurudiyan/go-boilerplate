package rest

import (
	"context"
	"fmt"
	"time"

	"github.com/wahyurudiyan/go-boilerplate/api/rest/routes"
	"github.com/wahyurudiyan/go-boilerplate/config"

	"github.com/gofiber/fiber/v2"
)

type httpServer struct {
	srv *fiber.App
	cfg *config.ServiceConfig
}

type IHttpServer interface {
	Listen() error
	Shutdown(ctx context.Context) error
}

func NewServer(cfg *config.ServiceConfig) IHttpServer {
	httpServer := &httpServer{
		cfg: cfg,
	}
	httpServer.srv = httpServer.newServer()

	return httpServer
}

func (s *httpServer) newServer() *fiber.App {
	/**
	Setup the server configuration
	*/

	opt := fiber.Config{
		AppName:       s.cfg.ApplicationName,
		BodyLimit:     s.cfg.RestBodyLimit * 1024 * 1024,
		CaseSensitive: s.cfg.RestRouteCaseSensitive,
		StrictRouting: s.cfg.RestStrictRoute,
	}

	readTimeout := s.cfg.RestReadTimeout
	if readTimeout != 0 {
		opt.ReadTimeout = time.Duration(readTimeout)
	}

	writeTimeout := s.cfg.RestWriteTimeout
	if writeTimeout != 0 {
		opt.WriteTimeout = time.Duration(writeTimeout)
	}

	/**
	Initialize http server
	*/
	app := fiber.New(opt)
	routes.Routes(app)

	return app
}

func (s *httpServer) Listen() error {
	addr := fmt.Sprintf(":%s", s.cfg.RestPort)
	return s.srv.Listen(addr)
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	return s.srv.ShutdownWithContext(ctx)
}
