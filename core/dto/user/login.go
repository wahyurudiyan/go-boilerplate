package user

import "time"

type LoginDTO struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Passowrd string `json:"passowrd,omitempty"`
}

type TokenDTO struct {
	User     *UserDTO   `json:"user,omitempty"`
	Token    string     `json:"token,omitempty"`
	ExpireAt *time.Time `json:"expire_at,omitempty"`
}
