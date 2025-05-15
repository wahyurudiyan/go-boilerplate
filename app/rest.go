package app

import (
	"context"
	"log/slog"

	"github.com/wahyurudiyan/go-boilerplate/api/rest/controller"
	"github.com/wahyurudiyan/go-boilerplate/api/rest/routes"
	"github.com/wahyurudiyan/go-boilerplate/cmd/rest"
	"github.com/wahyurudiyan/go-boilerplate/config"
	userRepo "github.com/wahyurudiyan/go-boilerplate/core/repositories/user"
	userSvc "github.com/wahyurudiyan/go-boilerplate/core/services/user"
	"github.com/wahyurudiyan/go-boilerplate/pkg/configz"
	"github.com/wahyurudiyan/go-boilerplate/pkg/graceful"
	"github.com/wahyurudiyan/go-boilerplate/pkg/sql"
)

func RestBootstrap(cfg *config.ServiceConfig) graceful.ExecCallback {

	var sqlConfig *sql.SQLConfig
	configz.MustLoadEnv(&sqlConfig)
	sql, err := sql.NewClient(sqlConfig)
	if err != nil {
		panic(err)
	}

	// User repositories contruction
	userRepo := userRepo.NewUserSQLRepository(sql)

	// User services construction
	repoDependency := userSvc.UserServicesImpl{
		UserRepo: userRepo,
	}
	userService := userSvc.NewUserService(repoDependency)

	// Controller bootstraping
	controllerDependency := controller.ControllerBootstrap{
		UserService: userService,
	}
	controller := controller.Bootstrap(controllerDependency)

	// Setup router
	router := routes.NewRouter(controller)

	return func(ctx context.Context) (graceful.ShutdownCallback, error) {
		srv := rest.NewGinServer(cfg)
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
