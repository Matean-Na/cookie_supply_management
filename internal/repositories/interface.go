package repositories

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepositoryInterface
}

func NewRepository(db *gorm.DB, rds *redis.Client) *Repository {
	return &Repository{
		UserRepositoryInterface: NewUserRepository(db, rds),
	}
}
