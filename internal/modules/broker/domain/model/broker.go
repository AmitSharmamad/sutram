package model

import (
	"time"

	"github.com/google/uuid"
)

// Broker represents a broker entity
type Broker struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	ClientID     string    `json:"client_id" db:"client_id"`
	ClientSecret string    `json:"client_secret" db:"client_secret"`
	AuthType     string    `json:"auth_type" db:"auth_type"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// BrokerUser represents a user's broker connection
type BrokerUser struct {
	ID       uuid.UUID  `json:"id" db:"id"`
	UserID   uuid.UUID  `json:"user_id" db:"user_id"`
	BrokerID uuid.UUID  `json:"broker_id" db:"broker_id"`
	Token    string     `json:"token" db:"token"`
	Expiry   *time.Time `json:"expiry" db:"expiry"`
	IsActive bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Order represents an order entity
type Order struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	UserID              uuid.UUID  `json:"user_id" db:"user_id"`
	PositionID          uuid.UUID  `json:"position_id" db:"position_id"`
	Exchange            string     `json:"exchange" db:"exchange"`
	Segment             string     `json:"segment" db:"segment"`
	Instrument          string     `json:"instrument" db:"instrument"`
	TransactionType     string     `json:"transaction_type" db:"transaction_type"`
	Quantity            int        `json:"quantity" db:"quantity"`
	Price               float64    `json:"price" db:"price"`
	ExecutedAt          *time.Time `json:"executed_at" db:"executed_at"`
	State               string     `json:"state" db:"state"`
	ProductType         string     `json:"product_type" db:"product_type"`
	OrderType           string     `json:"order_type" db:"order_type"`
	BrokerOrderID       *string    `json:"broker_order_id" db:"broker_order_id"`
	BrokerOrderResponse any        `json:"broker_order_response" db:"broker_order_response"`
	AlgoID              *string    `json:"algo_id" db:"algo_id"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}