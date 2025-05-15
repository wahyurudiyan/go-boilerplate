package user

import (
	"context"

	userDto "github.com/wahyurudiyan/go-boilerplate/core/dto/user"
)

type IUserServices interface {
	SignUp(ctx context.Context, user userDto.SignUpDTO) error
}

type AuthService interface{}
