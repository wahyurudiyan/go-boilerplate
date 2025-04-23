package main

import (
	"context"
	"log/slog"

	"github.com/wahyurudiyan/go-boilerplate/cmd/rest"
	env "github.com/wahyurudiyan/go-boilerplate/config"
	"github.com/wahyurudiyan/go-boilerplate/pkg/config"
)

func main() {
	var cfg env.ServiceConfig
	config.Load(&cfg)

	if err := rest.NewServer(&cfg).Listen(); err != nil {
		slog.ErrorContext(context.Background(), "unable to run server", "error", err)
	}
}
