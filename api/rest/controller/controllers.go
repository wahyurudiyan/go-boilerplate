package controller

import userSvc "github.com/wahyurudiyan/go-boilerplate/core/services/user"

type ControllerBootstrap struct {
	// Service interfaces
	UserService userSvc.IUserServices
}

func Bootstrap(cb ControllerBootstrap) *ControllerBootstrap {
	return &cb
}
