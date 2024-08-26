package base_service

import (
	"cookie_supply_management/pkg/base/base_model"
	"cookie_supply_management/pkg/base/base_repository"
	"encoding/json"
	"fmt"
)

type FindAll func(model *interface{}, scope base_repository.Scope, pager base_repository.Pager, order base_repository.OrderFilter, total *int64, searcher base_repository.Searcher, keyAll string) error
type FindOne func(model base_model.HasId, scope base_repository.Scope) error
type Create func(model base_model.HasId, keyAll string) error
type Update func(model base_model.HasId, keyAll string) error
type Delete func(model base_model.HasId, keyAll string) error

type FindAllInterface interface {
	FindAll(model *interface{}, scope base_repository.Scope, pager base_repository.Pager, order base_repository.OrderFilter, total *int64, searcher base_repository.Searcher, keyAll string) error
}

type FindOneInterface interface {
	FindOne(model base_model.HasId, scope base_repository.Scope) error
}

type CreateInterface interface {
	Create(model base_model.HasId, keyAll string) error
}

type UpdateInterface interface {
	Update(model base_model.HasId, keyAll string) error
}

type DeleteInterface interface {
	Delete(model base_model.HasId, keyAll string) error
}

type CrudServiceInterface interface {
	FindAllInterface
	FindOneInterface
	CreateInterface
	UpdateInterface
	DeleteInterface
}

type CrudService struct {
	repo base_repository.CrudRepositoryInterface
}

type CacheResponse struct {
	Model interface{}
	Total int64
}

func NewCrudService(repo base_repository.CrudRepositoryInterface) *CrudService {
	return &CrudService{
		repo: repo,
	}
}

func (c *CrudService) FindAll(model *interface{}, scope base_repository.Scope, pager base_repository.Pager, order base_repository.OrderFilter, total *int64, searcher base_repository.Searcher, keyAll string) error {
	cacheKey := fmt.Sprintf("%s_%d_%d", keyAll, pager.GetPage(), pager.GetPageSize())
	cachedData, err := c.repo.GetCache(cacheKey)
	if err == nil && searcher == nil {
		var cache CacheResponse
		if err = json.Unmarshal([]byte(*cachedData), &cache); err != nil {
			return err
		}
		*model = cache.Model
		*total = cache.Total
		return nil

	}

	if err = c.repo.FindAll(pager, order, scope, total, model, searcher); err != nil {
		return err
	}

	cache := CacheResponse{
		Total: *total,
		Model: model,
	}

	jsonData, err := json.Marshal(&cache)
	if err != nil {
		return err
	}
	if err = c.repo.SetCache(cacheKey, jsonData); err != nil {
		return err
	}

	return nil
}

func (c *CrudService) FindOne(model base_model.HasId, scope base_repository.Scope) error {
	return c.repo.FindOne(model.GetId(), scope, model)
}

func (c *CrudService) Create(model base_model.HasId, keyAll string) error {
	if err := c.repo.Create(model); err != nil {
		return err
	}

	if err := c.repo.DeleteCacheWithKey(keyAll); err != nil {
		return err
	}

	return nil
}

func (c *CrudService) Update(model base_model.HasId, keyAll string) error {
	if err := c.repo.Update(model); err != nil {
		return err
	}

	if err := c.repo.DeleteCacheWithKey(keyAll); err != nil {
		return err
	}

	return nil
}

func (c *CrudService) Delete(one base_model.HasId, keyAll string) error {
	if err := c.repo.Delete(one); err != nil {
		return err
	}

	if err := c.repo.DeleteCacheWithKey(keyAll); err != nil {
		return err
	}

	return nil
}
