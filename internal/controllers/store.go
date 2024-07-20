package controllers

import (
	"cookie_supply_management/internal/constants"
	"cookie_supply_management/internal/middlewares"
	"cookie_supply_management/internal/models"
	"cookie_supply_management/internal/services"
	"cookie_supply_management/pkg/base/base_controller"
	"cookie_supply_management/pkg/base/base_model"
	"cookie_supply_management/pkg/base/base_service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StoreController struct {
	services.StoreServiceInterface
	middlewares.AuthMiddlewareInterface

	*base_controller.CrudController
}

func NewStoreController(
	service services.StoreServiceInterface,
	middleware middlewares.AuthMiddlewareInterface,
	crudService base_service.CrudServiceInterface) *StoreController {
	c := &StoreController{
		CrudController:          base_controller.NewCrudController(crudService),
		AuthMiddlewareInterface: middleware,
	}
	c.CrudInterface = c
	c.ModelInterface = c
	return c
}

func (c *StoreController) GetAll() interface{} {
	return &[]models.Store{}
}

func (c *StoreController) GetOne() base_model.HasId {
	return &models.Store{}
}

func (c *StoreController) ScopeAll(db *gorm.DB) *gorm.DB {
	return db
}

func (c *StoreController) ScopeOne(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations)
}

func (c *StoreController) Register(r *gin.RouterGroup, s string) *gin.RouterGroup {
	g := r.Group(s)
	g.GET("", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.FindAll).Handle)
	g.GET(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.FindOne).Handle)
	g.POST("", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Create).Handle)
	g.PATCH(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Update).Handle)
	g.DELETE(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Delete).Handle)
	return g
}
