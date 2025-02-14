package algo

import (
	"algo/src/algo/expression"
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
	wire.Bind(new(shared.ExpressionService), new(*expression.Service)),
)
