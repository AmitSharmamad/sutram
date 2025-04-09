package auth

import "github.com/gin-gonic/gin"

type Router struct{}

func NewRoutesProv(
	engine *gin.Engine,
	controller *Controller,
) *Router {
	auth := engine.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/logout", controller.Logout)
		auth.GET("/google/login", controller.GoogleLogin)
		auth.GET("/google/callback", controller.GoogleCallback)
	}

	return &Router{}
}
