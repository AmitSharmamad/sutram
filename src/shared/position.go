package shared

import (
	"time"

	"algo/src/util"

	"github.com/google/uuid"
)

type Positions []Position

type PositionLogs []PositionLog

type RunPositions []RunPosition

type PositionType string

const (
	LATEST PositionType = "LATEST"
	NEW    PositionType = "NEW"
)

type Direction string

const (
	LONG  Direction = "LONG"
	SHORT Direction = "SHORT"
)

type ExitReason string

const (
	HIT_STOP_LOSS          ExitReason = "HIT_STOP_LOSS"
	HIT_TRAILING_STOP_LOSS ExitReason = "HIT_TRAILING_STOP_LOSS"
	HIT_TAKE_PROFIT        ExitReason = "HIT_TAKE_PROFIT"
	EXIT_CONDITION_MET     ExitReason = "EXIT_CONDITION_MET"
)

type Position struct {
	ID                   uuid.UUID     `json:"id" db:"id"`
	RunID                uuid.UUID     `json:"run_id" db:"run_id"`
	State                PositionState `json:"state" db:"state"`
	Direction            Direction     `json:"direction" db:"direction"`
	EntryPriceSignal     *float64      `json:"entry_price_signal" db:"entry_price_signal"`
	EntryTimeSignal      *time.Time    `json:"entry_time_signal" db:"entry_time_signal"`
	EntryPriceOrder      *float64      `json:"entry_price_order" db:"entry_price_order"`
	EntryTimeOrder       *time.Time    `json:"entry_time_order" db:"entry_time_order"`
	EntrySignal          *string       `json:"entry_signal" db:"entry_signal"`
	ExitPriceSignal      *float64      `json:"exit_price_signal" db:"exit_price_signal"`
	ExitTimeSignal       *time.Time    `json:"exit_time_signal" db:"exit_time_signal"`
	ExitPriceOrder       *float64      `json:"exit_price_order" db:"exit_price_order"`
	ExitTimeOrder        *time.Time    `json:"exit_time_order" db:"exit_time_order"`
	ExitSignal           *string       `json:"exit_signal" db:"exit_signal"`
	ExitReason           *ExitReason   `json:"exit_reason" db:"exit_reason"`
	ProfitLossAbsolute   *float64      `json:"profit_loss_absolute" db:"profit_loss_absolute"`
	ProfitLossPercentage *float64      `json:"profit_loss_percentage" db:"profit_loss_percentage"`
	TransactionCost      *float64      `json:"transaction_cost" db:"transaction_cost"`
	EntryFees            *float64      `json:"entry_fees" db:"entry_fees"`
	ExitFees             *float64      `json:"exit_fees" db:"exit_fees"`
	NetProfitLoss        *float64      `json:"net_profit_loss" db:"net_profit_loss"`
	PositionSize         *float64      `json:"position_size" db:"position_size"`
	EntryOrderID         *uuid.UUID    `json:"entry_order_id" db:"entry_order_id"`
	ExitOrderID          *uuid.UUID    `json:"exit_order_id" db:"exit_order_id"`
	Audit
}

func NewPosition(runID uuid.UUID, direction Direction) *Position {
	return &Position{
		ID:        util.NewUUID(),
		RunID:     runID,
		State:     NotEntered,
		Direction: direction, //Long or short
	}
}

func (position *Position) IsLong() bool {
	return position.Direction == LONG
}

func (position *Position) IsDone() bool {
	return position.State == Exited
}

type RunPosition struct {
	Run       *Run       `json:"run"`
	Positions *Positions `json:"positions"`
}

func NewRunPosition() *RunPosition {
	return new(RunPosition)
}
