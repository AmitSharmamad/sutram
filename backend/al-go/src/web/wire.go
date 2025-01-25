//go:build wireinject
// +build wireinject

package main

import (
	"algo/src/di"
	"algo/src/stock"
	"algo/src/strategy"
	"algo/src/template"
	"algo/src/user"
	"algo/src/web/config"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type WebContainer struct{}

func NewWebContainer(
	DiContainer *di.DIContainer,
	ConfigRoot *config.Server,
	Engine *gin.Engine,
	StrategyRoot *strategy.Root,
	StockRoot *stock.Root,
	TemplateRoot *template.Root,
	UserRoot *user.Root,
) *WebContainer {
	return &WebContainer{}
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
