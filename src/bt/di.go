package bt

import "github.com/google/wire"

type Root struct{}

func NewRoot(*Router) *Root {
	return &Root{}
}

var Set = wire.NewSet(
	NewRoot,

	NewRunRepoProv,
	NewRunMetricsRepoProv,
	NewPositionRepoProv,
	NewPositionLogRepoProv,

	NewServiceProv,
	NewControllerProv,

	NewRoutesProv,

	NewEngineProv,
)
