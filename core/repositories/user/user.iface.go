package user

import (
	"context"

	userEnt "github.com/wahyurudiyan/go-boilerplate/core/entities/user"
)

type IUserRepository interface {
	SaveUser(ctx context.Context, user userEnt.User) error
	SaveUsers(ctx context.Context, users []userEnt.User) error
	UpdateUser(ctx context.Context, user userEnt.User) error
	DeleteUserById(ctx context.Context, id int64) error
	DeleteUserByEmail(ctx context.Context, email string) error
	DeleteUserByUniqueId(ctx context.Context, uniqueId string) error

	RetrieveAllUser(ctx context.Context, offset, limit int) ([]userEnt.User, error)
	RetrieveUserById(ctx context.Context, id int64) (userEnt.User, error)
	RetrieveUserByIds(ctx context.Context, id []int64) ([]userEnt.User, error)
	RetrieveUserByEmail(ctx context.Context, email string) (userEnt.User, error)
	RetrieveUserByEmails(ctx context.Context, emails []string) ([]userEnt.User, error)
	RetrieveUserByUniqueId(ctx context.Context, uniqueId string) (userEnt.User, error)
	RetrieveUserByUniqueIds(ctx context.Context, uniqueId []string) ([]userEnt.User, error)
}
