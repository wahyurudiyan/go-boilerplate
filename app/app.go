package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/wahyurudiyan/go-boilerplate/config"
	userRepo "github.com/wahyurudiyan/go-boilerplate/core/repositories/user"
	userSvc "github.com/wahyurudiyan/go-boilerplate/core/services/user"
	"github.com/wahyurudiyan/go-boilerplate/pkg/configz"
	"github.com/wahyurudiyan/go-boilerplate/pkg/sql"
)

type appBoostraper struct {
	db          *sqlx.DB
	cfg         *config.ServiceConfig
	userService userSvc.IUserServices
}

func NewApp() *appBoostraper {
	var cfg *config.ServiceConfig
	if err := configz.LoadFromAWSParameterStore("/go-boilerplate/", "user", &cfg); err != nil {
		panic(err)
	}

	db, err := sql.NewClient(&cfg.Database)
	if err != nil {
		panic(err)
	}

	// User repositories contruction
	userRepo := userRepo.NewUserSQLRepository(db)

	// User services construction
	repoDependency := userSvc.UserServicesImpl{
		UserRepo: userRepo,
	}
	userService := userSvc.NewUserService(repoDependency)

	return &appBoostraper{
		db:          db,
		cfg:         cfg,
		userService: userService,
	}
}

func (a *appBoostraper) GetServiceConfig() *config.ServiceConfig {
	return a.cfg
}
