package di

import (
	"algo/src/bt"
	env "algo/src/core/config"
	"algo/src/core/db"
	"algo/src/stock"
	"algo/src/strategy"
	"algo/src/template"

	"github.com/google/wire"
)

type DIContainer struct{}

func NewDIContainer(
	StrategyRoot *strategy.Root,
) *DIContainer {
	return &DIContainer{}
}

var Set = wire.NewSet(
	NewDIContainer,
	env.Set,
	db.Set,
	bt.Set,
	strategy.Set,
	stock.Set,
	template.Set,
)
