package domain

import (
	"context"

	"algo/internal/modules/broker/application/dto"
)

// Broker defines the common interface that all broker implementations must satisfy
type Broker interface {
	// PlaceOrder places a new order with the broker
	PlaceOrder(ctx context.Context, request *dto.PlaceOrderRequest) (*dto.OrderResponse, error)
	
	// UpdateOrder modifies an existing order (also known as ModifyOrder)
	UpdateOrder(ctx context.Context, request *dto.ModifyOrderRequest) (*dto.OrderResponse, error)
	
	// CancelOrder cancels an existing order
	CancelOrder(ctx context.Context, request *dto.CancelOrderRequest) (*dto.OrderResponse, error)
}