package controllers

import (
	"cookie_supply_management/internal/constants"
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
	g.POST("login", AppHandler(uc.Login).Handle)
	g.POST("registration", uc.middleware.Authentication(), AppHandler(uc.Registration).Handle)
	g.POST("logout", uc.middleware.Authentication(constants.Accountant), AppHandler(uc.Logout).Handle)
	g.PATCH("update", uc.middleware.Authentication(constants.Accountant), AppHandler(uc.Update).Handle)

	g.GET("", uc.middleware.Authentication(constants.Accountant), AppHandler(uc.GetOne).Handle)
	return g
}

func (uc *UserController) GetOne(ctx *gin.Context) *AppError {
	username, exists := ctx.Get("username")
	if !exists {
		return BaseError("имя пользователя не найдено в контексте", http.StatusInternalServerError)
	}

	user, err := uc.service.GetByUsername(username.(string))
	if err != nil {
		return BaseError(err.Error(), http.StatusBadRequest)
	}

	return Ok(ctx, user)
}

func (uc *UserController) Update(ctx *gin.Context) *AppError {
	var input dto.UserUpdateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return BaseError(err.Error(), http.StatusBadRequest)
	}

	username, exists := ctx.Get("username")
	if !exists {
		return BaseError("имя пользователя не найдено в контексте", http.StatusInternalServerError)
	}

	httpCode, err := uc.service.Update(username.(string), input)
	if err != nil {
		return BaseError(err.Error(), httpCode)
	}

	return Ok(ctx, gin.H{"message": "success"})
}

func (uc *UserController) Logout(ctx *gin.Context) *AppError {
	e, _ := ctx.Get("username")
	username := e.(string)

	if err := uc.service.Logout(username); err != nil {
		return BaseError(err.Error(), http.StatusInternalServerError)
	}

	return Ok(ctx, gin.H{"message": "success"})
}

func (uc *UserController) Login(ctx *gin.Context) *AppError {
	var input dto.UserLoginDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return BaseError(err.Error(), http.StatusBadRequest)
	}

	token, code, err := uc.service.Login(input)
	if err != nil {
		return BaseError(err.Error(), code)
	}

	return Ok(ctx, gin.H{"token": token})
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
