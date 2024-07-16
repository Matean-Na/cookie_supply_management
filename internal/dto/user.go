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

type UserUpdateDTO struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	NewPasswordConfirm string `json:"new_password_confirm"`

	NewUserName string `json:"username"`
	NewRole     string `json:"role"`
}
