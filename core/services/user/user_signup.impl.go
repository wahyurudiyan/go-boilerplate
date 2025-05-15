package user

import (
	"context"
	"log/slog"

	userDto "github.com/wahyurudiyan/go-boilerplate/core/dto/user"
)

// func (u *UserServicesImpl) signToken() {}

func (u *UserServicesImpl) SignUp(ctx context.Context, registerUser userDto.SignUpDTO) error {
	user, err := registerUser.ToUserEntity()
	if err != nil {
		fields := []any{"name", registerUser.Fullname, "email", registerUser.Email} // because of error, let's get user data from parameter
		slog.Error("Error convert sign-up user to entity", fields...)
		return err
	}

	if err := u.UserRepo.SaveUser(ctx, user); err != nil {
		fields := []any{"name", registerUser.Fullname, "email", registerUser.Email} // more secure if getting data from parameter
		slog.Error("Error convert sign-up user to entity", fields...)
		return err
	}

	return nil
}
