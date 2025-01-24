package config

import (
	"algo/src/web/router"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer() *Server {
	return &Server{}
}

var Set = wire.NewSet(
	NewServer,

	router.NewEngine,

	router.HealthCheck,
)
