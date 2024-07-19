package models

import "cookie_supply_management/pkg/base/base_model"

type User struct {
	base_model.Entity
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
