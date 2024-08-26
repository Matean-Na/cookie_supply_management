package controllers

import (
	"cookie_supply_management/internal/constants"
	"cookie_supply_management/internal/middlewares"
	"cookie_supply_management/internal/models"
	"cookie_supply_management/pkg/base/base_controller"
	"cookie_supply_management/pkg/base/base_model"
	"cookie_supply_management/pkg/base/base_service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentController struct {
	middlewares.AuthMiddlewareInterface

	*base_controller.CrudController
}

func NewPaymentController(
	middleware middlewares.AuthMiddlewareInterface,
	crudService base_service.CrudServiceInterface) *PaymentController {
	c := &PaymentController{
		CrudController:          base_controller.NewCrudController(crudService),
		AuthMiddlewareInterface: middleware,
	}
	c.CrudInterface = c
	c.ModelInterface = c
	return c
}

func (c *PaymentController) GetAll() interface{} {
	return &[]models.Payment{}
}

func (c *PaymentController) GetOne() base_model.HasId {
	return &models.Payment{}
}

func (c *PaymentController) ScopeAll(db *gorm.DB) *gorm.DB {
	return db
}

func (c *PaymentController) ScopeOne(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations)
}

func (c *PaymentController) Register(r *gin.RouterGroup, s string) *gin.RouterGroup {
	g := r.Group(s)
	g.GET("", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.FindAll).Handle)
	g.GET(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.FindOne).Handle)
	g.POST("", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Create).Handle)
	g.PATCH(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Update).Handle)
	g.DELETE(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Delete).Handle)
	return g
}
