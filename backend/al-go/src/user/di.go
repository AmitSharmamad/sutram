package user

import "github.com/google/wire"

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
	NewControllerProv,
	NewRoutesProv,
)
