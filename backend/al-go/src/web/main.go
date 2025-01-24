package main

import "github.com/gin-gonic/gin"

func main() {
	println("")
	println("")

	println("starting web")

	engine := gin.Default()

	container := BuildContainer(engine)

	println(container)

	engine.Run(":9999")

	// container.ConfigRoot.Start()

	println("running web")

	println("")
	println("")
}
