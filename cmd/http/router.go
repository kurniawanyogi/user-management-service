package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-management-service/cmd/middleware"
	"user-management-service/common"
	"user-management-service/delivery"
)

type Router interface {
	Register() *gin.Engine
}

type router struct {
	engine     *gin.Engine
	common     common.IRegistry
	delivery   delivery.IRegistry
	middleware middleware.IMiddleware
}

func NewRouter(
	common common.IRegistry,
	delivery delivery.IRegistry,
	middleware middleware.IMiddleware,
) Router {
	return &router{
		engine:     gin.Default(),
		common:     common,
		delivery:   delivery,
		middleware: middleware,
	}
}

func (r *router) Register() *gin.Engine {

	// Landing
	r.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})

	// v1
	v1 := r.engine.Group("/v1")
	v1.POST("/users/registration", r.delivery.GetUserDelivery().Register)
	v1.POST("/users/login", r.delivery.GetUserDelivery().Login)

	// pengecekan token untuk endpoint get detail, update, dan delete melalui middleware AuthenticateToken()
	v1.GET("/users/:id", r.middleware.AuthenticateToken(), r.delivery.GetUserDelivery().Detail)
	v1.PUT("/users/:id", r.middleware.AuthenticateToken(), r.delivery.GetUserDelivery().Update)
	v1.DELETE("/users/:id", r.middleware.AuthenticateToken(), r.delivery.GetUserDelivery().Delete)

	return r.engine
}
