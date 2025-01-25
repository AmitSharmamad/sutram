package user

import "github.com/gin-gonic/gin"

type Router struct{}

func NewRoutesProv(
	engine *gin.Engine,
	controller *Controller,
) *Router {
	api := engine.Group("/user")
	{
		api.GET("", controller.Get)
		api.POST("", controller.Create)
	}

	return &Router{}
}
