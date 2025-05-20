package main

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/wahyurudiyan/go-boilerplate/app"
	env "github.com/wahyurudiyan/go-boilerplate/config"
	"github.com/wahyurudiyan/go-boilerplate/pkg/configz"
	"github.com/wahyurudiyan/go-boilerplate/pkg/graceful"
	"github.com/wahyurudiyan/go-boilerplate/pkg/telemetry"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

// @title Golang Boilerplate Example API
// @version 0.1
// @description This is a sample boilerplate project for golang backend service
// @termsOfService http://swagger.io/terms/
// @contact.name Go Boilerplate API Support
// @contact.email wahyurudiyan@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func main() {
	parentCtx := context.Background()

	var cfg *env.ServiceConfig
	if err := configz.LoadFromDotenv(".env.example", &cfg); err != nil {
		panic(err)
	}

	// Setup Opentelemetry SDK
	telemetryShutdown, err := telemetry.SetupOpentelemetry(parentCtx, telemetry.TelemetrySetup{
		Interval:           cfg.TelemetryMeterInterval,
		ServiceName:        cfg.ApplicationName,
		ServiceVersion:     cfg.ApplicationVersion,
		EnableRuntimeMeter: cfg.TelemetryEnableRuntimeMeter,
	})
	if err != nil {
		panic(err)
	}

	// Run application gracefully
	runApp := map[string]graceful.ExecCallback{
		"http-server": app.RestBootstrap(cfg),
		"opentelemetry": func(ctx context.Context) (graceful.ShutdownCallback, error) {
			return func(ctx context.Context) error {
				err = errors.Join(err, telemetryShutdown(ctx))
				slog.Error("Service shutting down with error", "error", err)
				return err
			}, nil
		},
	}
	if err := graceful.Run(parentCtx, time.Duration(10*time.Second), runApp); err != nil {
		panic(err)
	}
}
