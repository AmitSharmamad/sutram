//go:build wireinject
// +build wireinject

package main

import (
	"algo/src/auth"
	"algo/src/bt"
	"algo/src/core/http"
	"algo/src/di"
	"algo/src/stock"
	"algo/src/strategy"
	"algo/src/template"
	"algo/src/user"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type WebContainer struct{}

func NewWebContainer(
	DiContainer *di.DIContainer,
	Engine *gin.Engine,
	StrategyRoot *strategy.Root,
	StockRoot *stock.Root,
	TemplateRoot *template.Root,
	UserRoot *user.Root,
	AuthRoot *auth.Root,
	HTTPRoot *http.Root,
	BTRoot *bt.Root,
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
		),
	)
}
