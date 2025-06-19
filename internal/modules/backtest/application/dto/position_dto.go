package dto

import (
	"time"

	"github.com/google/uuid"
)

// PositionState defines the state of a position
type PositionState string

const (
	NotEntered PositionState = "NOT_ENTERED"
	Entered    PositionState = "ENTERED"
	Exited     PositionState = "EXITED"
)

// Direction defines the direction of a position
type Direction string

const (
	Long  Direction = "LONG"
	Short Direction = "SHORT"
)

// PositionDTO represents a data transfer object for positions
type PositionDTO struct {
	ID                   uuid.UUID     `json:"id"`
	RunID                uuid.UUID     `json:"run_id"`
	State                PositionState `json:"state"`
	Direction            Direction     `json:"direction"`
	EntryPriceSignal     float64       `json:"entry_price_signal"`
	EntryTimeSignal      time.Time     `json:"entry_time_signal"`
	EntryPriceOrder      float64       `json:"entry_price_order,omitempty"`
	EntryTimeOrder       time.Time     `json:"entry_time_order,omitempty"`
	EntrySignal          string        `json:"entry_signal"`
	ExitPriceSignal      float64       `json:"exit_price_signal,omitempty"`
	ExitTimeSignal       time.Time     `json:"exit_time_signal,omitempty"`
	ExitPriceOrder       float64       `json:"exit_price_order,omitempty"`
	ExitTimeOrder        time.Time     `json:"exit_time_order,omitempty"`
	ExitSignal           string        `json:"exit_signal,omitempty"`
	ExitReason           string        `json:"exit_reason,omitempty"`
	ProfitLossAbsolute   float64       `json:"profit_loss_absolute,omitempty"`
	ProfitLossPercentage float64       `json:"profit_loss_percentage,omitempty"`
	TransactionCost      float64       `json:"transaction_cost,omitempty"`
	EntryFees            float64       `json:"entry_fees,omitempty"`
	ExitFees             float64       `json:"exit_fees,omitempty"`
	NetProfitLoss        float64       `json:"net_profit_loss,omitempty"`
	PositionSize         float64       `json:"position_size"`
	EntryOrderID         *uuid.UUID    `json:"entry_order_id,omitempty"`
	ExitOrderID          *uuid.UUID    `json:"exit_order_id,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
}

// CreatePositionRequest represents a request to create a new position
type CreatePositionRequest struct {
	RunID            uuid.UUID `json:"run_id" binding:"required"`
	Direction        Direction `json:"direction" binding:"required"`
	EntryPriceSignal float64   `json:"entry_price_signal" binding:"required"`
	EntryTimeSignal  time.Time `json:"entry_time_signal" binding:"required"`
	EntrySignal      string    `json:"entry_signal" binding:"required"`
	PositionSize     float64   `json:"position_size" binding:"required"`
}

// EnterPositionRequest represents a request to mark a position as entered
type EnterPositionRequest struct {
	EntryPriceOrder float64    `json:"entry_price_order" binding:"required"`
	EntryTimeOrder  time.Time  `json:"entry_time_order" binding:"required"`
	EntryFees       float64    `json:"entry_fees"`
	EntryOrderID    *uuid.UUID `json:"entry_order_id"`
}

// ExitPositionRequest represents a request to mark a position as exited
type ExitPositionRequest struct {
	ExitPriceSignal float64    `json:"exit_price_signal" binding:"required"`
	ExitTimeSignal  time.Time  `json:"exit_time_signal" binding:"required"`
	ExitPriceOrder  float64    `json:"exit_price_order" binding:"required"`
	ExitTimeOrder   time.Time  `json:"exit_time_order" binding:"required"`
	ExitSignal      string     `json:"exit_signal" binding:"required"`
	ExitReason      string     `json:"exit_reason"`
	ExitFees        float64    `json:"exit_fees"`
	ExitOrderID     *uuid.UUID `json:"exit_order_id"`
}
