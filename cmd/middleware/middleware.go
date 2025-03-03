package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"user-management-service/common"
	"user-management-service/model"
	"user-management-service/service"
)

type IMiddleware interface {
	AuthenticateToken() gin.HandlerFunc
}
type middleware struct {
	common   common.IRegistry
	services service.IRegistry
}

func NewMiddleware(common common.IRegistry, services service.IRegistry) *middleware {
	return &middleware{
		common:   common,
		services: services,
	}
}

func (m *middleware) AuthenticateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx   = c.Request.Context()
			token string
		)
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.JSON(http.StatusUnauthorized, model.Response{
				Status:  common.StatusError,
				Message: "Authorization token is missing",
			})
			c.Abort()
			return
		}

		parts := strings.Split(bearerToken, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, model.Response{
				Status:  common.StatusError,
				Message: "Authorization header format must be Bearer <token>",
			})
			c.Abort()
			return
		}
		token = parts[1]

		err := m.services.GetUserService().Authenticate(ctx, token)
		if err != nil {
			if err.Error() == common.ErrInvalidToken.Error() {
				c.JSON(http.StatusUnauthorized, model.Response{
					Status:  common.StatusError,
					Message: err.Error(),
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, model.Response{
				Status:  common.StatusError,
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
		return
	}
}
