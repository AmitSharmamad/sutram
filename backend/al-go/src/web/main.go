package main

import (
	web "algo/src/web/config"

	"github.com/gin-gonic/gin"
)

func main() {
	println("")
	println("")

	println("starting web")

	engine := gin.Default()

	container := BuildContainer(engine)

	web.RegisterMiddleware(engine)

	println(container)

	engine.Run(":9999")

	println("running web")

	println("")
	println("")
}
