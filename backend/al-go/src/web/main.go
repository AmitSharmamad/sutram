package main

import (
	"algo/src/di"

	"github.com/gin-gonic/gin"
)

func init() {
	di.InitialiseDI()
}

func main() {
	server := gin.Default()

	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello",
		})
	})

	server.Run(":8080")
}
