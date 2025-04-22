package broker

import (
	env "algo/src/core/config"

	"github.com/google/wire"
)

// ProvideOrderBroker provides an implementation of OrderBroker based on configuration
func ProvideOrderBroker(envProv *env.EnvProv) OrderBroker {
	// Default to Dhan broker, but can be extended to support multiple broker types
	// based on configuration
	return NewDhan(envProv)
}

// ProvideDataVendor provides an implementation of DataVendor based on configuration
func ProvideDataVendor(envProv *env.EnvProv) DataVendor {
	// Default to Dhan as data vendor, but can be extended to support other data vendors
	// For now, reusing the broker implementation as it implements both interfaces
	return NewDhan(envProv)
}

// ProvideLiveTradingBroker provides an implementation of LiveTradingBroker for backward compatibility
func ProvideLiveTradingBroker(envProv *env.EnvProv) LiveTradingBroker {
	// For backward compatibility, providing a full LiveTradingBroker implementation
	return NewDhan(envProv)
}

// Set provides the Wire provider set for broker components
var Set = wire.NewSet(
	ProvideOrderBroker,
	ProvideDataVendor,
	ProvideLiveTradingBroker,
)
