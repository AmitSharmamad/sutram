package di

import (
	"algo/src/algo"
	"algo/src/auth"
	"algo/src/bt"
	env "algo/src/core/config"
	"algo/src/core/db"
	"algo/src/core/http"
	dvendor "algo/src/d-vendor"
	"algo/src/stock"
	"algo/src/strategy"
	"algo/src/template"
	"algo/src/user"

	"github.com/google/wire"
)

type DIContainer struct{}

func NewDIContainer(
	StrategyRoot *strategy.Root,
	AuthRoot *auth.Root,
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
	user.Set,
	auth.Set,
	http.Set,
	algo.Set,
	dvendor.Set,
)
