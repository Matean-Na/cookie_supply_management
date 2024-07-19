package controllers

import (
	"cookie_supply_management/internal/constants"
	"cookie_supply_management/internal/dto"
	"cookie_supply_management/internal/middlewares"
	"cookie_supply_management/internal/services"
	"cookie_supply_management/utils/parsers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StoreController struct {
	service    services.StoreServiceInterface
	middleware middlewares.AuthMiddlewareInterface
}

func NewStoreController(service services.StoreServiceInterface, middleware middlewares.AuthMiddlewareInterface) *StoreController {
	return &StoreController{
		service:    service,
		middleware: middleware,
	}
}

func (sc *StoreController) Register(r *gin.RouterGroup, s string) *gin.RouterGroup {
	g := r.Group(s)
	g.POST("create", sc.middleware.Authentication(constants.Accountant), AppHandler(sc.CreateStore).Handle)
	g.GET(":id", sc.middleware.Authentication(constants.Accountant), AppHandler(sc.GetOne).Handle)
	g.GET("", sc.middleware.Authentication(constants.Accountant), AppHandler(sc.GetAll).Handle)
	g.DELETE(":id", sc.middleware.Authentication(constants.Accountant), AppHandler(sc.DeleteStore).Handle)
	return g
}

func (sc *StoreController) CreateStore(ctx *gin.Context) *AppError {
	var input dto.StoreCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return BaseError(err.Error(), http.StatusBadRequest)
	}

	httpCode, err := sc.service.CreateStore(input)
	if err != nil {
		return BaseError(err.Error(), httpCode)
	}

	return Ok(ctx, gin.H{"message": "success"})
}

func (sc *StoreController) GetOne(ctx *gin.Context) *AppError {
	storeId := ctx.Param("id")
	id := parsers.ParamUint(storeId)

	store, httpCode, err := sc.service.GetOneStore(id)
	if err != nil {
		return BaseError(err.Error(), httpCode)
	}

	return Ok(ctx, store)
}

func (sc *StoreController) GetAll(ctx *gin.Context) *AppError {
	//page, page_size := GetPager(ctx)

	store, httpCode, err := sc.service.GetAllStore()
	if err != nil {
		return BaseError(err.Error(), httpCode)
	}

	return Ok(ctx, store)
}

func (sc *StoreController) DeleteStore(ctx *gin.Context) *AppError {
	//page, page_size := GetPager(ctx)
	storeId := ctx.Param("id")
	id := parsers.ParamUint(storeId)
	httpCode, err := sc.service.DeleteOneStore(id)
	if err != nil {
		return BaseError(err.Error(), httpCode)
	}

	return Ok(ctx, gin.H{"message": "success"})
}
