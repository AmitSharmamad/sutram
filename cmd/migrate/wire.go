//go:build wireinject
// +build wireinject

package main

import (
	"algo/cmd/migrate/handler"
	"algo/src/di"
	"algo/src/stock"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type WebContainer struct {
	dbHandler *handler.DBHandler
}

func NewWebContainer(
	DiContainer *di.DIContainer,
	StockRoot *stock.Root,
	dbHandler *handler.DBHandler,
) *WebContainer {
	return &WebContainer{
		dbHandler: dbHandler,
	}
}

func BuildContainer(
	engine *gin.Engine,
) *WebContainer {
	panic(
		wire.Build(
			NewWebContainer,
			di.Set,
			handler.NewDBHandler,
		),
	)
}
