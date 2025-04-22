package bt

import (
	"algo/src/broker"

	"github.com/google/wire"
)

type Root struct{}

func NewRoot(router *Router) *Root {
	return &Root{}
}

// NewBrokerProvider provides a broker implementation
func NewBrokerProvider() broker.LiveTradingBroker {
	// In a real implementation, you'd select a broker based on configuration
	// For now, return nil to avoid errors but this needs to be implemented with a real broker
	return nil
}

var Set = wire.NewSet(
	NewRoot,
	NewServiceProv,
	NewRunGroupRepoProv,
	NewRunRepoProv,
	NewPositionRepoProv,
	NewPositionLogRepoProv,
	NewRunMetricsRepoProv,
	NewRoutesProv,
	NewControllerProv,
	NewBrokerProvider,
	NewEngineProv,
)
