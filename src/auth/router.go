package auth

import "github.com/gin-gonic/gin"

type Router struct{}

func NewRoutesProv(
	engine *gin.Engine,
	controller *Controller,
) *Router {
	auth := engine.Group("/")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/logout", controller.Logout)
	}

	return &Router{}
}
