package middlewares

import (
	"cookie_supply_management/core/config"
	"cookie_supply_management/internal/constants"
	"cookie_supply_management/internal/services"
	"cookie_supply_management/pkg/security"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthMiddlewareInterface interface {
	Authentication(allowedRoles ...string) gin.HandlerFunc
}

type AuthMiddleware struct {
	service services.UserServiceInterface
}

func NewAuthMiddleware(service services.UserServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{service}
}

// Authentication Middleware
func (m *AuthMiddleware) Authentication(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := security.ExtractToken(ctx.Request)

		claims, err := security.ValidateToken(tokenString, config.Get().Token.SecretKey)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "неверный токен",
				"Code":  http.StatusUnauthorized,
			})
			ctx.Abort()
			return
		}

		_, err = m.service.GetToken(claims.Username)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "токен не найден",
				"Code":  http.StatusUnauthorized,
			})
			ctx.Abort()
			return
		}

		ctx.Set("username", claims.Username)
		username := claims.Username

		user, err := m.service.GetByUsername(username)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Error": "пользователь не найден",
				"Code":  http.StatusInternalServerError,
			})
			return
		}

		//check permission
		if !isRoleAllowed(user.Role, allowedRoles...) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func isRoleAllowed(role string, allowedRoles ...string) bool {
	allowedRolesMap := make(map[string]bool)
	for _, allowedRole := range allowedRoles {
		allowedRolesMap[allowedRole] = true
	}

	if role == constants.Admin {
		return true
	}

	if allowedRolesMap[role] {
		return true
	}

	return false
}
