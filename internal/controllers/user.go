package controllers

import (
	"cookie_supply_management/internal/dto"
	"cookie_supply_management/internal/middlewares"
	"cookie_supply_management/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	service    services.UserServiceInterface
	middleware middlewares.AuthMiddlewareInterface
}

func NewUserController(service services.UserServiceInterface, middleware middlewares.AuthMiddlewareInterface) *UserController {
	return &UserController{
		service:    service,
		middleware: middleware,
	}
}

func (uc *UserController) Register(r *gin.RouterGroup, s string) *gin.RouterGroup {
	g := r.Group(s)
	//g.POST("login", AppHandler(uc.Login).Handle)
	g.POST("registration", uc.middleware.Authentication(), AppHandler(uc.Registration).Handle)
	//g.POST("logout", AppHandler(uc.Logout).Handle)
	return g
}

func (uc *UserController) Registration(ctx *gin.Context) *AppError {
	var input dto.UserCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return BaseError(err.Error(), http.StatusBadRequest)
	}

	token, httpCode, err := uc.service.Registration(input)
	if err != nil {
		return BaseError(err.Error(), httpCode)
	}

	return Ok(ctx, gin.H{"token": token})
}
