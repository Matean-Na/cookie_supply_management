package services

import (
	"cookie_supply_management/core/config"
	"cookie_supply_management/internal/repositories"
)

type Service struct {
	UserServiceInterface
	StoreServiceInterface
}

func NewService(
	repos *repositories.Repository,
	conf *config.Config,
) *Service {
	return &Service{
		UserServiceInterface:  NewUserService(repos, conf),
		StoreServiceInterface: NewStoreService(repos, conf),
	}
}
