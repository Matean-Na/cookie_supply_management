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

type CookieController struct {
	*base_controller.CrudController
	middlewares.AuthMiddlewareInterface
}

func NewCookieController(service base_service.CrudServiceInterface, middleware middlewares.AuthMiddlewareInterface) *CookieController {
	c := &CookieController{
		CrudController:          base_controller.NewCrudController(service),
		AuthMiddlewareInterface: middleware,
	}
	c.CrudInterface = c
	c.ModelInterface = c
	return c
}

func (c *CookieController) GetAll() interface{} {
	return &[]models.Cookie{}
}

func (c *CookieController) GetOne() base_model.HasId {
	return &models.Cookie{}
}

func (c *CookieController) ScopeAll(db *gorm.DB) *gorm.DB {
	return db
}

func (c *CookieController) ScopeOne(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations)
}

func (c *CookieController) Register(r *gin.RouterGroup, s string) *gin.RouterGroup {
	g := r.Group(s)
	g.GET("", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.FindAll).Handle)
	g.GET(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.FindOne).Handle)
	g.POST("", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Create).Handle)
	g.PATCH(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Update).Handle)
	g.DELETE(":id", c.Authentication(constants.Accountant), base_controller.AppHandler(c.CrudInterface.Delete).Handle)
	return g
}
