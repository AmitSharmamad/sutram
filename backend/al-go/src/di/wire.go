//go:build wireinject
// +build wireinject

package di

import (
	"algo/src/bt"
	env "algo/src/core/config"
	"algo/src/core/db"
	"algo/src/strategy"

	"github.com/google/wire"
)

type DIContainer struct {
	EnvRoot      env.Root
	DbRoot       db.Root
	StrategyRoot strategy.Root
}

func NewDIContainer(
	env *env.Root,
	db *db.Root,
	bt *bt.Root,
	strategy *strategy.Root,
) *DIContainer {
	return &DIContainer{}
}

func InitialiseDI() *DIContainer {

	wire.Build(
		NewDIContainer,

		env.Set,

		db.Set,

		bt.Set,

		strategy.Set,
	)

	return &DIContainer{}
}
