package model

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

// Position represents a trading position within a run
type Position struct {
	ID                   uuid.UUID     `json:"id"`
	RunID                uuid.UUID     `json:"run_id"`
	State                PositionState `json:"state"`
	Direction            Direction     `json:"direction"`
	EntryPriceSignal     float64       `json:"entry_price_signal"`
	EntryTimeSignal      time.Time     `json:"entry_time_signal"`
	EntryPriceOrder      float64       `json:"entry_price_order"`
	EntryTimeOrder       time.Time     `json:"entry_time_order"`
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

// NewPosition creates a new trading position
func NewPosition(
	runID uuid.UUID,
	direction Direction,
	entryPriceSignal float64,
	entryTimeSignal time.Time,
	entrySignal string,
	positionSize float64,
) *Position {
	now := time.Now()
	return &Position{
		ID:               uuid.New(),
		RunID:            runID,
		State:            NotEntered,
		Direction:        direction,
		EntryPriceSignal: entryPriceSignal,
		EntryTimeSignal:  entryTimeSignal,
		EntrySignal:      entrySignal,
		PositionSize:     positionSize,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Enter sets the position as entered
func (p *Position) Enter(entryPriceOrder float64, entryTimeOrder time.Time, entryFees float64, entryOrderID *uuid.UUID) {
	p.State = Entered
	p.EntryPriceOrder = entryPriceOrder
	p.EntryTimeOrder = entryTimeOrder
	p.EntryFees = entryFees
	p.EntryOrderID = entryOrderID
	p.UpdatedAt = time.Now()
}

// Exit sets the position as exited
func (p *Position) Exit(
	exitPriceSignal float64,
	exitTimeSignal time.Time,
	exitPriceOrder float64,
	exitTimeOrder time.Time,
	exitSignal string,
	exitReason string,
	exitFees float64,
	exitOrderID *uuid.UUID,
) {
	p.State = Exited
	p.ExitPriceSignal = exitPriceSignal
	p.ExitTimeSignal = exitTimeSignal
	p.ExitPriceOrder = exitPriceOrder
	p.ExitTimeOrder = exitTimeOrder
	p.ExitSignal = exitSignal
	p.ExitReason = exitReason
	p.ExitFees = exitFees
	p.ExitOrderID = exitOrderID

	// Calculate profit/loss
	p.CalculateProfitLoss()
	p.UpdatedAt = time.Now()
}

// CalculateProfitLoss calculates profit/loss for the position
func (p *Position) CalculateProfitLoss() {
	if p.State != Exited {
		return
	}

	entryPrice := p.EntryPriceOrder
	exitPrice := p.ExitPriceOrder

	// Calculate absolute P/L
	switch p.Direction {
	case Long:
		p.ProfitLossAbsolute = (exitPrice - entryPrice) * p.PositionSize
	case Short:
		p.ProfitLossAbsolute = (entryPrice - exitPrice) * p.PositionSize
	}

	// Calculate percentage P/L
	p.ProfitLossPercentage = (p.ProfitLossAbsolute / (entryPrice * p.PositionSize)) * 100

	// Calculate transaction costs
	p.TransactionCost = p.EntryFees + p.ExitFees

	// Calculate net P/L
	p.NetProfitLoss = p.ProfitLossAbsolute - p.TransactionCost
}
