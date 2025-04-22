package broker

import (
	"algo/src/shared"
	"context"
)

type OrderType string

const (
	MARKET OrderType = "MARKET"
	LIMIT  OrderType = "LIMIT"
)

type OrderStatus string

const (
	PENDING   OrderStatus = "PENDING"
	EXECUTED  OrderStatus = "EXECUTED"
	REJECTED  OrderStatus = "REJECTED"
	CANCELLED OrderStatus = "CANCELLED"
)

// Order represents a trading order
type Order struct {
	ID            string
	InstrumentID  int32
	Type          OrderType
	Direction     shared.Direction
	Quantity      int
	Price         *float64 // Optional for market orders
	Status        OrderStatus
	ErrorMessage  string
	EntryTime     string
	ExecutionTime string
}

// OrderBroker interface defines methods for order execution that any broker implementation must support
type OrderBroker interface {
	// Connect establishes connection with the broker
	Connect(ctx context.Context) error

	// PlaceOrder places a new order and returns the order ID
	PlaceOrder(ctx context.Context, order Order) (string, error)

	// CancelOrder cancels an existing order
	CancelOrder(ctx context.Context, orderID string) error

	// GetOrderStatus retrieves the current status of an order
	GetOrderStatus(ctx context.Context, orderID string) (OrderStatus, error)

	// GetPositions retrieves all current positions
	GetPositions(ctx context.Context) ([]shared.Position, error)
}

// DataVendor interface defines methods for market data operations
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

// LiveTradingBroker interface combines OrderBroker and DataVendor interfaces
// This is kept for backward compatibility with existing code
type LiveTradingBroker interface {
	OrderBroker
	DataVendor
}
