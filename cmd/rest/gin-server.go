package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahyurudiyan/go-boilerplate/api/rest/routes"
	"github.com/wahyurudiyan/go-boilerplate/config"
)

type httpGinServer struct {
	srv    *http.Server
	cfg    *config.ServiceConfig
	router *gin.Engine
}

func (s *httpGinServer) newGinServer() *gin.Engine {
	engine := gin.Default()
	return engine
}

func (s *httpGinServer) Listen() error {
	listener := http.Server{
		Addr:         fmt.Sprintf(":%v", s.cfg.RestPort),
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.cfg.RestReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.cfg.RestWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.cfg.RestIdleTimeout) * time.Second,
	}
	s.srv = &listener

	routes.Routes(s.router)

	slog.Info("[SERVER] running...", "port", s.cfg.RestPort)
	if err := listener.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *httpGinServer) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return http.ErrServerClosed
	}

	if err := s.srv.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}
