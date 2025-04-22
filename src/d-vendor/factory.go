package dvendor

import (
	"fmt"
)

// DataVendorFactory creates appropriate DataVendor implementations based on vendor type
type DataVendorFactory struct{}

// NewDataVendorFactory creates a new data vendor factory
func NewDataVendorFactory() *DataVendorFactory {
	return &DataVendorFactory{}
}

// CreateVendor creates a DataVendor implementation based on the provided configuration
func (f *DataVendorFactory) CreateVendor(config VendorConfig) (DataVendor, error) {
	switch config.Type {
	case DhanVendor:
		// Create Dhan-specific configuration
		dhanConfig := DhanVendorConfig{
			VendorConfig: config,
		}

		// Check if Dhan client ID is provided in parameters
		if dhanClientID, ok := config.Parameters["dhanClientID"]; ok {
			dhanConfig.DhanClientID = dhanClientID
		}

		return NewDhanDataVendor(dhanConfig), nil
	case ZerodhaVendor:
		return NewZerodhaDataVendor(config), nil
	default:
		return nil, fmt.Errorf("unsupported vendor type: %s", config.Type)
	}
}
