package middlewares

import "cookie_supply_management/internal/services"

type Middlewares struct {
	AuthMiddlewareInterface
}

func NewMiddlewares(service *services.Service) *Middlewares {
	return &Middlewares{
		AuthMiddlewareInterface: NewAuthMiddleware(service.UserServiceInterface),
	}
}
