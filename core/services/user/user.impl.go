package user

import (
	"aidanwoods.dev/go-paseto"
	userRepository "github.com/wahyurudiyan/go-boilerplate/core/repositories/user"
)

var _ IUserServices = (*UserServicesImpl)(nil)

type UserServicesImpl struct {
	// Tokenizer will generate PASETO Token
	tokenizer paseto.Token

	// Add service dependency below
	UserRepo userRepository.IUserRepository
}

func NewUserService(userSvc UserServicesImpl) IUserServices {
	userSvc.tokenizer = paseto.NewToken()
	return &userSvc
}
