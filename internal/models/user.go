package models

type User struct {
	Entity
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
