//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

type DIContainer struct{}

func NewDIContainer() *DIContainer {
	return &DIContainer{}
}

func InitialiseDI() *DIContainer {

	wire.Build(
		NewDIContainer,
	)

	return &DIContainer{}
}
