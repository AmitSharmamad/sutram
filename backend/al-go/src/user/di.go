package user

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
	wire.Bind(new(shared.UserService), new(*Service)),
	NewServiceProv,
	NewControllerProv,
	NewRoutesProv,
)
