package user

import (
	"time"

	"github.com/rs/xid"
	userEnt "github.com/wahyurudiyan/go-boilerplate/core/entities/user"
	"golang.org/x/crypto/bcrypt"
)

type SignUpDTO struct {
	Role     string `json:"role,omitempty"`
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (s SignUpDTO) ToUserEntity() (userEnt.User, error) {
	// Generate unique_id for this user
	uniqueId := xid.New().String()

	// Hash inputted password from request
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(s.Password), 10)
	if err != nil {
		return userEnt.User{}, err
	}

	return userEnt.User{
		Role:      s.Role,
		Email:     s.Email,
		UniqueId:  uniqueId,
		Fullname:  s.Fullname,
		Username:  s.Username,
		Password:  string(hashedPass),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
