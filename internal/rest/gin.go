package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahyurudiyan/go-boilerplate/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type ginServer struct {
	srv    *http.Server
	cfg    *config.ServiceConfig
	router *gin.Engine
}

func newGinServer(cfg *config.ServiceConfig) *ginServer {
	var s ginServer
	s.cfg = cfg
	ginEngine := gin.Default()
	ginEngine.Use(otelgin.Middleware(cfg.ApplicationName))

	s.router = ginEngine
	return &s
}

func (s *ginServer) RegisterRoutes(routesFn func(*gin.Engine)) {
	if s.router != nil {
		routesFn(s.router)
	}
}

func (s *ginServer) Run() error {
	handler := otelhttp.NewHandler(s.router.Handler(), "")
	listener := http.Server{
		Addr:         fmt.Sprintf(":%v", s.cfg.RestPort),
		Handler:      handler,
		ReadTimeout:  time.Duration(s.cfg.RestReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.cfg.RestWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.cfg.RestIdleTimeout) * time.Second,
	}
	s.srv = &listener

	slog.Info("[SERVER] running...", "port", s.cfg.RestPort)
	if err := listener.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *ginServer) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return http.ErrServerClosed
	}

	if err := s.srv.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}
