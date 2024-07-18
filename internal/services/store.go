package services

import (
	"cookie_supply_management/core/config"
	"cookie_supply_management/internal/dto"
	"cookie_supply_management/internal/models"
	"cookie_supply_management/internal/repositories"
	"net/http"
)

type StoreServiceInterface interface {
	CreateStore(input dto.StoreCreateDTO) (httpCode int, err error)
	GetOneStore(id uint) (store models.Store, httpCode int, err error)
	GetAllStore() (store []models.Store, httpCode int, err error)
	DeleteOneStore(id uint) (httpCode int, err error)
}

type StoreService struct {
	repo repositories.StoreRepositoryInterface
	conf *config.Config
}

func NewStoreService(repo repositories.StoreRepositoryInterface, conf *config.Config) *StoreService {
	return &StoreService{
		repo: repo,
		conf: conf,
	}
}

func (s *StoreService) CreateStore(input dto.StoreCreateDTO) (httpCode int, err error) {

	if err = s.repo.InsertStore(input.Name, input.Address, input.Contact, input.PhoneNumber); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *StoreService) GetOneStore(id uint) (store models.Store, httpCode int, err error) {
	var output models.Store
	if output, err = s.repo.SelectStore(id); err != nil {
		return models.Store{}, http.StatusInternalServerError, err
	}
	return output, http.StatusOK, nil
}
func (s *StoreService) GetAllStore() (store []models.Store, httpCode int, err error) {
	var output []models.Store
	if output, err = s.repo.SelectStores(); err != nil {
		return []models.Store{}, http.StatusInternalServerError, err
	}
	return output, http.StatusOK, nil
}

func (s *StoreService) DeleteOneStore(id uint) (httpCode int, err error) {
	if err = s.repo.DeleteStore(id); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
