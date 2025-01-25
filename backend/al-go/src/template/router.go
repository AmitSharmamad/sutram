package template

import "github.com/gin-gonic/gin"

type Router struct{}

func RegisterRoutesProv(
	engine *gin.Engine,
	controller *Controller,
) *Router {

	api := engine.Group("/template")
	{
		api.GET("/", controller.Get)
	}

	println(api)

	return &Router{}
}
