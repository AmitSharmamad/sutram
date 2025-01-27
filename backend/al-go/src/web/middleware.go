package main

import "github.com/gin-gonic/gin"

func RegisterMiddleware(engine *gin.Engine) {
	engine.Use(GlobalErrorHandler)
}
