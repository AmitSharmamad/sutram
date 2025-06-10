package dto

import (
	"time"

	"github.com/google/uuid"
)

// Common DTOs for broker operations - these are used by the Broker interface

// PlaceOrderRequest represents a standardized request for placing orders across all brokers
type PlaceOrderRequest struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id" binding:"required"`
	PositionID      uuid.UUID `json:"position_id" binding:"required"`
	Exchange        string    `json:"exchange" binding:"required"`
	Segment         string    `json:"segment" binding:"required"`
	Instrument      string    `json:"instrument" binding:"required"`
	TransactionType string    `json:"transaction_type" binding:"required"` // BUY, SELL
	Quantity        int       `json:"quantity" binding:"required"`
	Price           float64   `json:"price" binding:"required"`
	ProductType     string    `json:"product_type" binding:"required"`     // CNC, MIS, NRML
	OrderType       string    `json:"order_type" binding:"required"`       // LIMIT, MARKET, SL, SL-M
	AlgoID          string    `json:"algo_id,omitempty"`
}

// ModifyOrderRequest represents a standardized request for modifying orders across all brokers
type ModifyOrderRequest struct {
	OrderID         string  `json:"order_id" binding:"required"`
	Quantity        *int    `json:"quantity,omitempty"`
	Price           *float64 `json:"price,omitempty"`
	OrderType       *string `json:"order_type,omitempty"`       // LIMIT, MARKET, SL, SL-M
	TriggerPrice    *float64 `json:"trigger_price,omitempty"`
}

// CancelOrderRequest represents a standardized request for canceling orders across all brokers
type CancelOrderRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

// OrderResponse represents a standardized response for order operations across all brokers
type OrderResponse struct {
	ID                   uuid.UUID `json:"id"`
	UserID               uuid.UUID `json:"user_id"`
	PositionID           uuid.UUID `json:"position_id"`
	Exchange             string    `json:"exchange"`
	Segment              string    `json:"segment"`
	Instrument           string    `json:"instrument"`
	TransactionType      string    `json:"transaction_type"`
	Quantity             int       `json:"quantity"`
	Price                float64   `json:"price"`
	State                string    `json:"state"`
	ProductType          string    `json:"product_type"`
	OrderType            string    `json:"order_type"`
	BrokerOrderID        string    `json:"broker_order_id,omitempty"`
	BrokerOrderResponse  any       `json:"broker_order_response,omitempty"`
	AlgoID               string    `json:"algo_id,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	ExecutedAt           *time.Time `json:"executed_at,omitempty"`
}