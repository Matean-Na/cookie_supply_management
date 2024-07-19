package controllers

import (
	"cookie_supply_management/internal/models"
	"cookie_supply_management/pkg/base/base_controller"
	"cookie_supply_management/pkg/base/base_model"
	"cookie_supply_management/pkg/base/base_service"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CookieTypeController struct {
	*base_controller.CrudController
}

func NewCookieTypeController(service base_service.CrudServiceInterface) *CookieTypeController {
	c := &CookieTypeController{
		CrudController: base_controller.NewCrudController(service),
	}
	c.CrudInterface = c
	c.ModelInterface = c
	return c
}

func (c *CookieTypeController) GetAll() interface{} {
	return &[]models.CookieType{}
}

func (c *CookieTypeController) GetOne() base_model.HasId {
	return &models.CookieType{}
}

func (c *CookieTypeController) ScopeAll(db *gorm.DB) *gorm.DB {
	return db
}

func (c *CookieTypeController) ScopeOne(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations)
}
