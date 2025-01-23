//go:build wireinject
// +build wireinject

package main

import (
	env "algo/src/core/config"
	"algo/src/core/db"

	"github.com/google/wire"
)

type DIContainer struct {
	EnvRoot env.Root
	DbRoot  db.Root
}

func NewDIContainer(
	env *env.Root,
	db *db.Root,
) *DIContainer {
	return &DIContainer{}
}

func InitialiseDI() *DIContainer {

	wire.Build(
		NewDIContainer,

		env.Set,

		db.Set,
	)

	return &DIContainer{}
}
