package di

import (
	"algo/src/bt"
	env "algo/src/core/config"
	"algo/src/core/db"
	"algo/src/strategy"

	"github.com/google/wire"
)

type DIContainer struct{}

func NewDIContainer() *DIContainer {
	return &DIContainer{}
}

var Set = wire.NewSet(
	NewDIContainer,
	env.Set,
	db.Set,
	bt.Set,
	strategy.Set,
)
