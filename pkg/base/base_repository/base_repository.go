package base_repository

import (
	"context"
	"cookie_supply_management/pkg/base/base_model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"time"
)

type CrudRepositoryInterface interface {
	FindAll(pager Pager, order OrderFilter, scope Scope, total *int64, model interface{}, searcher Searcher) error
	FindOne(id uint, scope Scope, model interface{}) error
	Create(model base_model.HasId) error
	Update(model base_model.HasId) error
	Delete(model base_model.HasId) error

	GetCache(key string) (*string, error)
	SetCache(key string, data []byte) error
	DeleteCacheWithKey(key string) error
}

type CrudRepository struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewCrudRepository(db *gorm.DB, rds *redis.Client) CrudRepositoryInterface {
	return &CrudRepository{
		db:  db,
		rds: rds,
	}
}

func (cr *CrudRepository) GetCache(key string) (*string, error) {
	cachedData, err := cr.rds.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	return &cachedData, nil
}

func (cr *CrudRepository) SetCache(key string, data []byte) error {
	if err := cr.rds.Set(context.Background(), key, data, time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func (cr *CrudRepository) DeleteCacheWithKey(key string) error {
	keys, err := cr.rds.Keys(context.Background(), key+"*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		if err = cr.rds.Del(context.Background(), keys...).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (cr *CrudRepository) FindAll(pager Pager, order OrderFilter, scope Scope, total *int64, model interface{}, searcher Searcher) error {
	if err := cr.db.Model(model).Scopes(scope).Count(total).Error; err != nil {
		return err
	}

	if searcher != nil {
		if err := cr.db.Where(searcher.getQueryJoin()).Joins(searcher.getJoinModels()).Scopes(pager.Paginate(), order.sort(), scope).Find(model).Error; err != nil {
			return err
		}
	} else {
		if err := cr.db.Scopes(pager.Paginate(), order.sort(), scope).Find(model).Error; err != nil {
			return err
		}
	}
	return nil
}

func (cr *CrudRepository) FindOne(id uint, scope Scope, model interface{}) error {
	if err := cr.db.Scopes(scope).Where("id = ?", id).First(model).Error; err != nil {
		return err
	}
	return nil
}

func (cr *CrudRepository) Create(model base_model.HasId) error {
	if err := cr.db.Create(model).Error; err != nil {
		return err
	}
	return nil
}

func (cr *CrudRepository) Update(model base_model.HasId) error {
	if err := cr.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(model).Error; err != nil {
		return err
	}
	return nil
}

func (cr *CrudRepository) Delete(model base_model.HasId) error {
	if err := cr.db.Delete(model, model.GetId()).Error; err != nil {
		return err
	}
	return nil
}
