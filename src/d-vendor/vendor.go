package dvendor

import (
	"algo/src/shared"
	"context"
)

// DataVendor interface defines methods for market data operations
// This is separated from broker operations to follow the single responsibility principle
type DataVendor interface {
	// Connect establishes connection with the data vendor
	Connect(ctx context.Context) error

	// GetQuote gets real-time quote for an instrument
	GetQuote(ctx context.Context, instrumentID int32) (*shared.OHLCV, error)

	// Subscribe to real-time updates for instruments
	Subscribe(ctx context.Context, instrumentIDs []int32) error

	// Unsubscribe from real-time updates
	Unsubscribe(ctx context.Context, instrumentIDs []int32) error
}

// VendorType represents the type of data vendor
type VendorType string

const (
	// DhanVendor represents the Dhan data vendor
	DhanVendor VendorType = "DHAN"

	// ZerodhaVendor represents the Zerodha data vendor
	ZerodhaVendor VendorType = "ZERODHA"
)

// VendorConfig holds common configuration for all vendors
type VendorConfig struct {
	Type       VendorType
	ClientID   string
	Token      string
	BaseURL    string
	Parameters map[string]string // Additional vendor-specific parameters
}
