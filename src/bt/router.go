package bt

import "github.com/gin-gonic/gin"

type Router struct{}

func NewRoutesProv(
	engine *gin.Engine,
	controller *Controller,
) *Router {
	api := engine.Group("/bt")
	{
		api.POST("/start", controller.DoBacktest)
		api.POST("/run/:id/start", controller.Start)
		api.POST("/run/:id/pause", controller.Pause)
		api.POST("/run/:id/resume", controller.Resume)
		api.POST("/run/:id/metrics", controller.GetMetrics)

		api.GET("/run", controller.GetAllRuns)
		api.GET("/run/:id", controller.GetRun)
		api.GET("/run/:id/log", controller.GetRunLogs)

		// Live trading deployment endpoint
		api.POST("/strategy/:id/deploy", controller.DeployLive)
	}

	println(api)

	return &Router{}
}
