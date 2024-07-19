package base_controller

import (
	"cookie_supply_management/core/connect"
	"cookie_supply_management/pkg/base/base_model"
	"cookie_supply_management/pkg/localizer"
	"cookie_supply_management/utils/parsers"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

type AppHandler func(ctx *gin.Context) *AppError

func (a AppHandler) Handle(ctx *gin.Context) {
	if err := a(ctx); err != nil {
		ctx.JSON(err.Code, gin.H{"error": err})
	}
}

func Ok(ctx *gin.Context, i interface{}) *AppError {
	ctx.JSON(http.StatusOK, i)
	return nil
}

func ErrNotExist(ctx *gin.Context, i interface{}) *AppError {
	ctx.JSON(http.StatusNotFound, i)
	return nil
}

func OkT(ctx *gin.Context, t int64, i interface{}) *AppError {
	ctx.JSON(http.StatusOK, gin.H{"total": t, "data": i})
	return nil
}

func OkJson(ctx *gin.Context, s string) *AppError {
	ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(s))
	return nil
}

func okMessage(ctx *gin.Context, SuccessMessage string, Message string) *AppError {
	ctx.JSON(http.StatusOK, gin.H{"Success": SuccessMessage, "Message": Message})
	return nil
}

func BaseError(message string, code int) *AppError {
	appError := &AppError{
		Code:    code,
		Message: message,
	}

	return appError
}

func ErrNotFound(ctx *gin.Context, instance interface{}, id uint, params string) *AppError {
	var message string
	if params != "" {
		message = localizer.Localize(
			ctx,
			"exception:failed-to-fetch-records-with-params",
			map[string]interface{}{
				"Table":      base_model.GetTableName(instance, connect.DB),
				"Parameters": params,
			},
		)
	} else {
		if id != 0 {
			message = localizer.Localize(ctx, "exception:failed-to-fetch-one-record", map[string]interface{}{
				"Table": base_model.GetTableName(instance, connect.DB),
				"ID":    id,
			})
		} else {
			message = localizer.Localize(ctx, "exception:could-not-fetch-records", map[string]string{
				"Table": base_model.GetTableName(instance, connect.DB),
			})
		}
	}

	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func ErrNotCreated(ctx *gin.Context, err error, instance interface{}) *AppError {
	message := localizer.Localize(ctx, "exception:failed-to-create-record", map[string]interface{}{
		"Table": base_model.GetTableName(instance, connect.DB),
	})
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func ErrNotUpdated(err error) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Record not updated",
	}
}

func ErrNotDeleted(err error) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Record not deleted",
	}
}

func ErrDenied(err error) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: "Access denied",
	}
}

func ErrBadRequest(ctx *gin.Context, err error, data map[string]interface{}) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: localizer.Localize(ctx, err.Error(), data),
	}
}

func LocalizeError(ctx *gin.Context, err error) *AppError {
	var pgxError *pgconn.PgError
	var fieldName string
	var detailedInfo string

	switch e := err.(type) {
	case localizer.LocalizeError:
		errors.As(e.Source, &pgxError)
		if pgxError != nil {
			fieldName = parsers.ParseFieldName(pgxError.Detail)
			detailedInfo = pgxError.Detail
		}
		if pgxError != nil {
			fieldName = parsers.ParseFieldName(pgxError.Detail)
			detailedInfo = pgxError.Detail
		}
		return &AppError{
			Message: e.Localize(ctx),
			Code:    http.StatusBadRequest,
		}
	case *localizer.LocalizeError:
		errors.As(e.Source, &pgxError)
		if pgxError != nil {
			fieldName = parsers.ParseFieldName(pgxError.Detail)
			detailedInfo = pgxError.Detail
		}
		return &AppError{
			Message: e.Localize(ctx),
			Code:    http.StatusBadRequest,
		}
	}
	errors.As(err, &pgxError)
	if pgxError != nil {
		fieldName = parsers.ParseFieldName(pgxError.Detail)
		detailedInfo = pgxError.Detail
	}
	return DefaultError(ctx, err, fieldName, detailedInfo)
}

func DefaultError(ctx *gin.Context, err error, fieldName string, detailedInfo string) *AppError {
	return &AppError{
		Message: localizer.Localize(ctx, err.Error(), nil),
		Code:    http.StatusBadRequest,
	}
}

func I18nError(c *gin.Context, model interface{}, errCode string) *AppError {
	codes := map[string]int{
		"exception:could-not-count-records": http.StatusBadRequest,
		"exception:could-not-fetch-records": http.StatusNotFound,
		"exception:failed-to-create-record": http.StatusBadRequest,
		"exception:failed-to-update-record": http.StatusBadRequest,
		"exception:failed-to-delete-record": http.StatusBadRequest,
	}

	var code = http.StatusBadRequest
	if c, found := codes[errCode]; found {
		code = c
	}

	table := base_model.GetTableName(model, connect.DB)
	return &AppError{
		Message: localizer.Localize(c, errCode, map[string]interface{}{
			"Table": table,
		}),
		Code: code,
	}
}
