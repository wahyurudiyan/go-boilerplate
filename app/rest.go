package app

import (
	"context"
	"log/slog"

	"github.com/wahyurudiyan/go-boilerplate/api/rest/controller"
	"github.com/wahyurudiyan/go-boilerplate/api/rest/routes"
	"github.com/wahyurudiyan/go-boilerplate/internal/rest"
	"github.com/wahyurudiyan/go-boilerplate/pkg/graceful"
)

func (a *appBoostraper) RestBootstrap() graceful.ExecCallback {
	// Controller bootstraping
	controllerDependency := controller.ControllerBootstrap{
		UserService: a.userService,
	}
	controller := controller.Bootstrap(controllerDependency)
	// Setup router
	router := routes.NewRouter(controller)
	return func(ctx context.Context) (graceful.ShutdownCallback, error) {
		srv := rest.NewGinServer(a.cfg)
		srv.RegisterRoutes(router.Routes)

		// Run the server
		if err := srv.Run(); err != nil {
			slog.ErrorContext(context.Background(), "unable to run server", "error", err)
			return nil, err
		}

		return func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		}, nil
	}
}
