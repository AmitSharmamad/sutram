package stock

import (
	"algo/src/shared"

	"github.com/google/wire"
)

type Root struct{}

func NewRootProv(
	router *Router,
) *Root {
	return &Root{}
}

var Set = wire.NewSet(
	NewRootProv,
	NewRepoProv,
	NewServiceProv,
	wire.Bind(new(shared.StockService), new(*Service)),
	NewControllerProv,
	NewRoutesProv,
)
