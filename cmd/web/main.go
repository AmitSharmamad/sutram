package main

import (
	"algo/src/web"

	"github.com/gin-gonic/gin"
)

func main() {
	println("")
	println("")

	println("starting web")

	engine := gin.Default()

	web.RegisterMiddleware(engine)

	container := BuildContainer(engine)

	println(container)

	engine.Run(":9090")

	println("running web")

	println("")
	println("")
}
