package repositories

import (
	"context"
	"cookie_supply_management/internal/models"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	SelectExistByUserName(username string) error
	SelectUserByUsername(username string) (models.User, error)
	InsertUser(username, password, role string) error

	SetToken(username, token string) error
	GetToken(username string) (string, error)
	DeleteToken(username string) error
}

type UserRepository struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) SelectExistByUserName(username string) error {
	if err := r.db.First(&models.User{}, "username = ?", username).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("пользователь не найден")
	}

	return nil
}

func (r *UserRepository) SelectUserByUsername(username string) (models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) InsertUser(username, password, role string) error {
	if err := r.db.Create(models.User{
		Username: username,
		Password: password,
		Role:     role,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) SetToken(username, token string) error {
	return r.rds.Set(context.Background(), username, token, 0).Err()
}

func (r *UserRepository) GetToken(username string) (string, error) {
	return r.rds.Get(context.Background(), username).Result()
}

func (r *UserRepository) DeleteToken(username string) error {
	return r.rds.Del(context.Background(), username).Err()
}
