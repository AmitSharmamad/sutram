package algo

import (
	"algo/src/algo/expression"
	"algo/src/algo/indicator"
	"algo/src/shared"

	"github.com/google/wire"
)

type Root struct{}

func NewRootProv(
	service *expression.Service,
) *Root {
	return &Root{}
}

var Set = wire.NewSet(
	expression.NewService,
	indicator.NewServiceProv,
	wire.Bind(new(shared.ExpressionService), new(*expression.Service)),
	wire.Bind(new(shared.IndicatorService), new(*indicator.Service)),
)
