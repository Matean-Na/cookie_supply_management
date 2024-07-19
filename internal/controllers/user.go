package controllers

import (
	"cookie_supply_management/internal/constants"
	"cookie_supply_management/internal/dto"
	"cookie_supply_management/internal/middlewares"
	"cookie_supply_management/internal/services"
	"cookie_supply_management/pkg/base/base_controller"
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
	g.POST("login", base_controller.AppHandler(uc.Login).Handle)
	g.POST("registration", uc.middleware.Authentication(), base_controller.AppHandler(uc.Registration).Handle)
	g.POST("logout", uc.middleware.Authentication(constants.Accountant), base_controller.AppHandler(uc.Logout).Handle)
	g.PATCH("update", uc.middleware.Authentication(constants.Accountant), base_controller.AppHandler(uc.Update).Handle)

	g.GET("", uc.middleware.Authentication(constants.Accountant), base_controller.AppHandler(uc.GetOne).Handle)
	return g
}

func (uc *UserController) GetOne(ctx *gin.Context) *base_controller.AppError {
	username, exists := ctx.Get("username")
	if !exists {
		return base_controller.BaseError("имя пользователя не найдено в контексте", http.StatusInternalServerError)
	}

	user, err := uc.service.GetByUsername(username.(string))
	if err != nil {
		return base_controller.BaseError(err.Error(), http.StatusBadRequest)
	}

	return base_controller.Ok(ctx, user)
}

func (uc *UserController) Update(ctx *gin.Context) *base_controller.AppError {
	var input dto.UserUpdateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return base_controller.BaseError(err.Error(), http.StatusBadRequest)
	}

	username, exists := ctx.Get("username")
	if !exists {
		return base_controller.BaseError("имя пользователя не найдено в контексте", http.StatusInternalServerError)
	}

	httpCode, err := uc.service.Update(username.(string), input)
	if err != nil {
		return base_controller.BaseError(err.Error(), httpCode)
	}

	return base_controller.Ok(ctx, gin.H{"message": "success"})
}

func (uc *UserController) Logout(ctx *gin.Context) *base_controller.AppError {
	e, _ := ctx.Get("username")
	username := e.(string)

	if err := uc.service.Logout(username); err != nil {
		return base_controller.BaseError(err.Error(), http.StatusInternalServerError)
	}

	return base_controller.Ok(ctx, gin.H{"message": "success"})
}

func (uc *UserController) Login(ctx *gin.Context) *base_controller.AppError {
	var input dto.UserLoginDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return base_controller.BaseError(err.Error(), http.StatusBadRequest)
	}

	token, code, err := uc.service.Login(input)
	if err != nil {
		return base_controller.BaseError(err.Error(), code)
	}

	return base_controller.Ok(ctx, gin.H{"token": token})
}

func (uc *UserController) Registration(ctx *gin.Context) *base_controller.AppError {
	var input dto.UserCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return base_controller.BaseError(err.Error(), http.StatusBadRequest)
	}

	token, httpCode, err := uc.service.Registration(input)
	if err != nil {
		return base_controller.BaseError(err.Error(), httpCode)
	}

	return base_controller.Ok(ctx, gin.H{"token": token})
}
