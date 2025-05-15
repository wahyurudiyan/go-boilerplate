package user

type UserDTO struct {
	Role     string `json:"role,omitempty"`
	Email    string `json:"email,omitempty"`
	UniqueId string `json:"unique_id,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	Username string `json:"username,omitempty"`
	Passowrd string `json:"passowrd,omitempty"`
	Status   bool   `json:"status,omitempty"`
}
