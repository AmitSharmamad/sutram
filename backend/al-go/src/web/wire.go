//go:build wireinject
// +build wireinject

package main

import (
	"algo/src/di"

	"github.com/google/wire"
)

type WebContainer struct{}

func NewWebContainer(
	*di.DIContainer,
) *WebContainer {
	return &WebContainer{}
}

func BuildContainer() *WebContainer {
	panic(
		wire.Build(
			NewWebContainer,
			di.Set,
		),
	)
}
