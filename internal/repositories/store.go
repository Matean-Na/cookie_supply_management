package repositories

import (
	"cookie_supply_management/internal/models"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type StoreRepositoryInterface interface {
	InsertStore(name, address, contact, phoneNumber string) error
	SelectStore(id uint) (models.Store, error)
	SelectStores() ([]models.Store, error)
	DeleteStore(id uint) error
}

type StoreRepository struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewStoreRepository(db *gorm.DB, rds *redis.Client) *StoreRepository {
	return &StoreRepository{
		db:  db,
		rds: rds,
	}
}

func (r *StoreRepository) InsertStore(name, address, contact, phoneNumber string) error {
	if err := r.db.Create(&models.Store{
		Name:        name,
		Address:     address,
		Contact:     contact,
		PhoneNumber: phoneNumber,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *StoreRepository) SelectStore(id uint) (models.Store, error) {
	var store models.Store
	if err := r.db.Where("id = ?", id).First(&store).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Store{}, errors.New("магазин не найден")
	}

	return store, nil
}

func (r *StoreRepository) SelectStores() ([]models.Store, error) {
	var store []models.Store
	if err := r.db.Find(&store).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return []models.Store{}, errors.New("магазинов не найдено")
	}

	return store, nil
}

func (r *StoreRepository) DeleteStore(id uint) error {
	var store models.Store
	if err := r.db.Where("id = ?", id).First(&store).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("магазин не найден")
	}

	if store.Debt.LessThan(decimal.NewFromInt(0)) {
		return errors.New("магазин не может быть удален, есть задолжность")
	}
	if err := r.db.Delete(&store).Error; err != nil {
		return err
	}

	return nil
}
