package main

import (
	web "algo/src/web/config"
	"algo/src/web/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	println("starting web")

	engine := gin.Default()

	web.RegisterMiddleware(engine)

	container := BuildContainer(engine)

	println(container)

	// Start the WebSocket server in a separate goroutine
	go websocket.StartWebSocketServer()

	engine.Run(":9090")

	println("running web")
}
