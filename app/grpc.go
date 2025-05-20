package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/wahyurudiyan/go-boilerplate/api/grpc/handler"
	userPb "github.com/wahyurudiyan/go-boilerplate/api/grpc/service-user"
	"github.com/wahyurudiyan/go-boilerplate/pkg/graceful"
	"google.golang.org/grpc"
)

func (a *appBoostraper) GRPCBootstrap() graceful.ExecCallback {
	slog.Info("[GRPC] server running", "port", a.cfg.GrpcPort)
	return func(ctx context.Context) (graceful.ShutdownCallback, error) {
		grpcHost := fmt.Sprintf("0.0.0.0:%s", a.cfg.GrpcPort)
		grpcListener, err := net.Listen("tcp", grpcHost)
		if err != nil {
			panic(err)
		}

		grpcservice := handler.NewGRPCHandler(a.userService)

		grpcServer := grpc.NewServer(grpc.ConnectionTimeout(a.cfg.GrpcTimeout))
		userPb.RegisterServiceUserServer(grpcServer, grpcservice)
		grpcServer.Serve(grpcListener)
		return func(ctx context.Context) error {
			slog.Info("[GRPC] server shutting down!")
			grpcServer.GracefulStop()
			return nil
		}, nil
	}
}
