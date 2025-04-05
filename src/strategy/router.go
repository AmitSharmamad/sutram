package strategy

import (
	"algo/src/core/http"

	"github.com/gin-gonic/gin"
)

// https://medium.com/@wahyubagus1910/build-scalable-restful-api-with-golang-gin-gonic-framework-43793c730d10

type Router struct{}

func NewRoutesProv(
	engine *gin.Engine,
	controller *Controller,
	authMiddleware *http.AuthMiddleware,
) *Router {
	api := engine.Group("/strategy")
	{
		api.POST("", controller.Create)
		api.GET("", controller.GetAll)
		api.GET("/:id", controller.Get)
		api.PUT("/:id", controller.Update)
	}

	router := &Router{}
	return router
}
