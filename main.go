package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/wahyurudiyan/go-boilerplate/cmd/rest"
	env "github.com/wahyurudiyan/go-boilerplate/config"
	"github.com/wahyurudiyan/go-boilerplate/pkg/config"
	"github.com/wahyurudiyan/go-boilerplate/pkg/graceful"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

// @title Golang Boilerplate Example API
// @version 0.1
// @description This is a sample boilerplate project for golang backend service
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func main() {
	var cfg env.ServiceConfig
	config.Load(&cfg)

	runApp := map[string]graceful.ExecCallback{
		"http-server": func(ctx context.Context) (graceful.ShutdownCallback, error) {
			srv := rest.NewGinServer(&cfg)
			if err := srv.Listen(); err != nil {
				slog.ErrorContext(context.Background(), "unable to run server", "error", err)
			}

			return func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			}, nil
		},
		"test": func(ctx context.Context) (graceful.ShutdownCallback, error) {
			slog.Info("This only graceful execution test exec")
			return func(ctx context.Context) error {
				slog.Info("This only graceful execution test shutdown")
				return nil
			}, nil
		},
	}

	graceful.Run(context.Background(), time.Duration(10*time.Second), runApp)
}
