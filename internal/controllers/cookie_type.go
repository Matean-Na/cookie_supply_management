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

type CookieTypeController struct {
	*base_controller.CrudController
	middlewares.AuthMiddlewareInterface
}

func NewCookieTypeController(service base_service.CrudServiceInterface, middleware middlewares.AuthMiddlewareInterface) *CookieTypeController {
	c := &CookieTypeController{
		CrudController:          base_controller.NewCrudController(service),
		AuthMiddlewareInterface: middleware,
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

func (c *CookieTypeController) Register(r *gin.RouterGroup, s string) *gin.RouterGroup {
	g := r.Group(s)
	g.GET("", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.FindAll).Handle)
	g.GET(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.FindOne).Handle)
	g.POST("", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Create).Handle)
	g.PATCH(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Update).Handle)
	g.DELETE(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Delete).Handle)
	return g
}
