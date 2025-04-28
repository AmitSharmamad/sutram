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
	ID                   uuid.UUID     `json:"id" db:"id"`
	RunID                uuid.UUID     `json:"run_id" db:"run_id"`
	State                PositionState `json:"state" db:"state"`
	Direction            Direction     `json:"direction" db:"direction"`
	EntryPriceSignal     float64       `json:"entry_price_signal" db:"entry_price_signal"`
	EntryTimeSignal      time.Time     `json:"entry_time_signal" db:"entry_time_signal"`
	EntryPriceOrder      *float64      `json:"entry_price_order,omitempty" db:"entry_price_order"`
	EntryTimeOrder       *time.Time    `json:"entry_time_order,omitempty" db:"entry_time_order"`
	EntrySignal          string        `json:"entry_signal" db:"entry_signal"`
	ExitPriceSignal      *float64      `json:"exit_price_signal,omitempty" db:"exit_price_signal"`
	ExitTimeSignal       *time.Time    `json:"exit_time_signal,omitempty" db:"exit_time_signal"`
	ExitPriceOrder       *float64      `json:"exit_price_order,omitempty" db:"exit_price_order"`
	ExitTimeOrder        *time.Time    `json:"exit_time_order,omitempty" db:"exit_time_order"`
	ExitSignal           *string       `json:"exit_signal,omitempty" db:"exit_signal"`
	ExitReason           *string       `json:"exit_reason,omitempty" db:"exit_reason"`
	ProfitLossAbsolute   *float64      `json:"profit_loss_absolute,omitempty" db:"profit_loss_absolute"`
	ProfitLossPercentage *float64      `json:"profit_loss_percentage,omitempty" db:"profit_loss_percentage"`
	TransactionCost      *float64      `json:"transaction_cost,omitempty" db:"transaction_cost"`
	EntryFees            *float64      `json:"entry_fees,omitempty" db:"entry_fees"`
	ExitFees             *float64      `json:"exit_fees,omitempty" db:"exit_fees"`
	NetProfitLoss        *float64      `json:"net_profit_loss,omitempty" db:"net_profit_loss"`
	PositionSize         float64       `json:"position_size" db:"position_size"`
	EntryOrderID         *uuid.UUID    `json:"entry_order_id,omitempty" db:"entry_order_id"`
	ExitOrderID          *uuid.UUID    `json:"exit_order_id,omitempty" db:"exit_order_id"`
	CreatedAt            time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at" db:"updated_at"`
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
	p.EntryPriceOrder = &entryPriceOrder
	p.EntryTimeOrder = &entryTimeOrder

	if entryFees > 0 {
		p.EntryFees = &entryFees
	}

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
	p.ExitPriceSignal = &exitPriceSignal
	p.ExitTimeSignal = &exitTimeSignal
	p.ExitPriceOrder = &exitPriceOrder
	p.ExitTimeOrder = &exitTimeOrder
	p.ExitSignal = &exitSignal
	p.ExitReason = &exitReason

	if exitFees > 0 {
		p.ExitFees = &exitFees
	}

	p.ExitOrderID = exitOrderID

	// Calculate profit/loss
	p.CalculateProfitLoss()
	p.UpdatedAt = time.Now()
}

// CalculateProfitLoss calculates profit/loss for the position
func (p *Position) CalculateProfitLoss() {
	if p.State != Exited || p.EntryPriceOrder == nil || p.ExitPriceOrder == nil {
		return
	}

	entryPrice := *p.EntryPriceOrder
	exitPrice := *p.ExitPriceOrder

	// Calculate absolute P/L
	var profitLossAbsolute float64
	switch p.Direction {
	case Long:
		profitLossAbsolute = (exitPrice - entryPrice) * p.PositionSize
	case Short:
		profitLossAbsolute = (entryPrice - exitPrice) * p.PositionSize
	}
	p.ProfitLossAbsolute = &profitLossAbsolute

	// Calculate percentage P/L
	profitLossPercentage := (profitLossAbsolute / (entryPrice * p.PositionSize)) * 100
	p.ProfitLossPercentage = &profitLossPercentage

	// Calculate transaction costs
	var transactionCost float64
	if p.EntryFees != nil {
		transactionCost += *p.EntryFees
	}
	if p.ExitFees != nil {
		transactionCost += *p.ExitFees
	}
	p.TransactionCost = &transactionCost

	// Calculate net P/L
	netProfitLoss := profitLossAbsolute - transactionCost
	p.NetProfitLoss = &netProfitLoss
}
