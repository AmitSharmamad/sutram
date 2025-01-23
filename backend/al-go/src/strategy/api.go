package strategy

import "github.com/gin-gonic/gin"

// https://medium.com/@wahyubagus1910/build-scalable-restful-api-with-golang-gin-gonic-framework-43793c730d10

func Init(
	controller *Controller,
) *gin.Engine {

	router := gin.Default()

	api := router.Group("/strategy")

	// get all strategies
	api.GET("", controller.GetStrategies)
	// get a specific strategy
	api.GET("/:id", controller.GetStrategy)
	// get runs for a strategy
	api.GET("/:id/run", controller.GetStrategyRuns)
	// get results for a run
	api.GET("/:id/run/:runId", controller.GetStrategyRun)

	// run a strategy
	api.POST("/:strategyId/run", controller.Run)

	return router
}
