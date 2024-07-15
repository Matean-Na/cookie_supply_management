package repositories

import (
	"gorm.io/gorm"
)

type Repository struct {
	UserRepositoryInterface
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepositoryInterface: NewUserRepository(db),
	}
}
