package server

import (
	"cookie_supply_management/core/config"
	"cookie_supply_management/internal/controllers"
	"cookie_supply_management/internal/middlewares"
	"cookie_supply_management/internal/services"
	"cookie_supply_management/pkg/base/base_service"
	"github.com/elliotchance/orderedmap"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Register(r *gin.RouterGroup, s string) *gin.RouterGroup
}

func SetupRoutesWithDeps(service *services.Service, serviceCrud base_service.CrudServiceInterface) *gin.Engine {
	//init middlewares
	middleware := middlewares.NewMiddlewares(service)

	e := gin.Default()
	e.MaxMultipartMemory = 8 << 20
	e.Use(middlewares.CORSMiddleware())

	api := e.Group("/api")
	{
		r := orderedmap.NewOrderedMap()
		r.Set("user", controllers.NewUserController(service.UserServiceInterface, middleware.AuthMiddlewareInterface))
		r.Set("store", controllers.NewStoreController(service.StoreServiceInterface, middleware.AuthMiddlewareInterface))
		r.Set("cookie_type", controllers.NewCookieTypeController(serviceCrud))

		for g := r.Front(); g != nil; g = g.Next() {
			if c, ok := g.Value.(Controller); ok {
				c.Register(api, g.Key.(string))
			}
		}
	}

	mediaPath := config.Get().Dir.Media
	e.Static("/media", mediaPath)

	return e
}
