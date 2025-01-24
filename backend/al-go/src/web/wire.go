//go:build wireinject
// +build wireinject

package main

import (
	"algo/src/di"
	"algo/src/web/config"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type WebContainer struct {
	DiContainer *di.DIContainer
	ConfigRoot  *config.Server
	Engine      *gin.Engine
}

func NewWebContainer(
	DiContainer *di.DIContainer,
	ConfigRoot *config.Server,
	Engine *gin.Engine,
) *WebContainer {
	return &WebContainer{
		DiContainer: DiContainer,
		ConfigRoot:  ConfigRoot,
		Engine:      Engine,
	}
}

func NewEngine(
	engine *gin.Engine,
) interface{} {
	return func() *gin.Engine {
		return engine
	}
}

func BuildContainer(
	engine *gin.Engine,
) *WebContainer {
	panic(
		wire.Build(
			NewWebContainer,
			di.Set,
			config.Set,
		),
	)
}
