package handler

import (
	userPb "github.com/wahyurudiyan/go-boilerplate/api/grpc/service-user"
	userSvc "github.com/wahyurudiyan/go-boilerplate/core/services/user"
)

type grpcHandler struct {
	userService userSvc.IUserServices
}

func NewGRPCHandler(userService userSvc.IUserServices) userPb.ServiceUserServer {
	return &grpcHandler{
		userService: userService,
	}
}
