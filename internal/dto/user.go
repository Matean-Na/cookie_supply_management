package dto

type UserCreateDTO struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	Role            string `json:"role"`
}

type UserLoginDTO struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
