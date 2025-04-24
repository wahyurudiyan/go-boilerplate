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

func main() {
	var cfg env.ServiceConfig
	config.Load(&cfg)

	runApp := map[string]graceful.ExecCallback{
		"http-server": func(ctx context.Context) (graceful.ShutdownCallback, error) {
			srv := rest.NewServer(&cfg)
			if err := srv.Listen(); err != nil {
				slog.ErrorContext(context.Background(), "unable to run server", "error", err)
				return nil, err
			}

			return func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			}, nil
		},
	}

	graceful.Runner(context.Background(), time.Duration(10*time.Second), runApp)
}
