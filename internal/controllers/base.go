package controllers

import (
	"cookie_supply_management/utils/parsers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AppError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type AppHandler func(ctx *gin.Context) *AppError

func (a AppHandler) Handle(ctx *gin.Context) {
	if err := a(ctx); err != nil {
		ctx.JSON(err.Code, err)
	}
}

func Ok(ctx *gin.Context, i interface{}) *AppError {
	ctx.JSON(http.StatusOK, i)
	return nil
}

func OkT(ctx *gin.Context, t int64, i interface{}) *AppError {
	ctx.JSON(http.StatusOK, gin.H{"total": t, "data": i})
	return nil
}

func OkMessage(ctx *gin.Context, message string) *AppError {
	ctx.JSON(http.StatusOK, gin.H{"message": message})
	return nil
}

func BaseError(err string, code int) *AppError {
	return &AppError{
		Error: err,
		Code:  code,
	}
}

func GetPager(ctx *gin.Context) (int, int) {
	page := parsers.ParamInt(ctx.Query("page"))
	pageSize := parsers.ParamInt(ctx.Query("page_size"))

	if page == 0 {
		page = 1
	}
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	return page, pageSize
}
