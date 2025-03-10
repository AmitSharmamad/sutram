package stock

import "github.com/gin-gonic/gin"

type Router struct{}

func NewRoutesProv(
	engine *gin.Engine,
	controller *Controller,
) *Router {
	api := engine.Group("/stock")
	{
		api.GET("", controller.Get)
	}

	println(api)

	return &Router{}
}
