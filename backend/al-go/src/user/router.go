package user

import (
	"github.com/gin-gonic/gin"
)

type Router struct{}

func NewRoutesProv(
	engine *gin.Engine,
	controller *Controller,
) *Router {
	user := engine.Group("/user")
	{
		user.GET("", controller.Get)
		user.POST("", controller.Create)
	}

	return &Router{}
}
