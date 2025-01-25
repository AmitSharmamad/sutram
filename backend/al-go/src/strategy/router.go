package strategy

import "github.com/gin-gonic/gin"

// https://medium.com/@wahyubagus1910/build-scalable-restful-api-with-golang-gin-gonic-framework-43793c730d10

type Router struct{}

func NewRoutesProv(
	engine *gin.Engine,
	controller *Controller,
) *Router {
	api := engine.Group("/strategy")
	{
		api.GET("", controller.GetStrategies)
		api.GET("/:id", controller.GetStrategy)
		api.GET("/:id/run", controller.GetStrategyRuns)
		api.GET("/:id/run/:runId", controller.GetStrategyRun)

		api.POST("/:strategyId/run", controller.Run)
	}

	return &Router{}
}
