package router

import "github.com/gin-gonic/gin"

func NewEngine() *gin.Engine {
	return gin.Default()
}

type Router struct{}

func HealthCheck(engine *gin.Engine) *Router {
	engine.GET("/", func(c *gin.Context) {
		println("coming inside this")
		c.SecureJSON(200, true)
	})
	return &Router{}
}
