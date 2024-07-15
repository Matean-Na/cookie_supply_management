package dto

type UserCreateDTO struct {
	UserName        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
	Role            string `json:"role" validate:"required,role"`
}
