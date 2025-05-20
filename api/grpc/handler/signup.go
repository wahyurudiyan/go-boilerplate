package handler

import (
	"context"

	userPb "github.com/wahyurudiyan/go-boilerplate/api/grpc/service-user"
	userDto "github.com/wahyurudiyan/go-boilerplate/core/dto/user"
)

func (h *grpcHandler) SignUp(ctx context.Context, m *userPb.SignUpRequest) (*userPb.SignUpResponse, error) {
	userDto := userDto.SignUpDTO{
		Role:     m.GetRole(),
		Email:    m.GetEmail(),
		Fullname: m.GetFullname(),
		Username: m.GetUsername(),
		Password: m.GetPassword(),
	}

	if err := h.userService.SignUp(ctx, userDto); err != nil {
		return nil, err
	}

	return &userPb.SignUpResponse{}, nil
}
