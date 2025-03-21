package stock

import (
	"algo/src/shared"

	"github.com/google/wire"
)

type Root struct{}

func NewRoot(router *Router) *Root {
	return &Root{}
}

var Set = wire.NewSet(
	NewRoot,
	NewMasterRepoProv,
	NewRepoProv,
	NewServiceProv,
	wire.Bind(new(shared.StockService), new(*Service)),
	NewControllerProv,
	NewRoutesProv,
)
