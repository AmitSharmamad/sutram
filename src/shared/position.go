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
	HighestPrice         *float64
	LowestPrice          *float64
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

func (position *Position) CalculateStopLossPrice(stopLoss float64, stopLossType *ValueType) float64 {
	entryPrice := position.EntryPriceOrder

	var stopLossPrice float64

	if *stopLossType == ABSOLUTE {
		if position.Direction == LONG {
			stopLossPrice = *entryPrice - stopLoss
		} else if position.Direction == SHORT {
			stopLossPrice = *entryPrice + stopLoss
		} else {
			panic("Invalid PositionType")
		}
	} else if *stopLossType == PERCENTAGE {
		if position.Direction == LONG {
			stopLossPrice = *entryPrice * (1 - stopLoss)
		} else if position.Direction == SHORT {
			stopLossPrice = *entryPrice * (1 + stopLoss)
		} else {
			panic("Invalid PositionType")
		}
	} else {
		panic("Invalid StopLossType")
	}

	return stopLossPrice
}

func (position *Position) CalculateTargetProfitPrice(targetProfit float64, targetProfitType *ValueType) float64 {
	entryPrice := position.EntryPriceOrder

	var targetProfitPrice float64

	if *targetProfitType == ABSOLUTE {
		if position.Direction == LONG {
			targetProfitPrice = *entryPrice + targetProfit
		} else if position.Direction == SHORT {
			targetProfitPrice = *entryPrice - targetProfit
		} else {
			panic("Invalid PositionType")
		}
	} else if *targetProfitType == PERCENTAGE {
		if position.Direction == LONG {
			targetProfitPrice = *entryPrice * (1 + targetProfit)
		} else if position.Direction == SHORT {
			targetProfitPrice = *entryPrice * (1 - targetProfit)
		} else {
			panic("Invalid PositionType")
		}
	} else {
		panic("Invalid StopLossType")
	}

	return targetProfitPrice
}

func (position *Position) CalculateTrailingStopLossPrice(
	currentClose float64,
	trailingStopLoss *float64,
	trailingStopLossType ValueType,
) float64 {
	var trailingStopLossPrice float64

	if position.Direction == LONG {
		f := util.Max(*position.HighestPrice, currentClose)
		position.HighestPrice = &f

		if trailingStopLossType == ABSOLUTE {
			trailingStopLossPrice = *position.HighestPrice - *trailingStopLoss
		} else if trailingStopLossType == PERCENTAGE {
			trailingStopLossPrice = *position.HighestPrice * (1 - *trailingStopLoss)
		} else {
			panic("Invalid TrailingStopLossType")
		}

	} else if position.Direction == SHORT {
		f := util.Min(*position.LowestPrice, currentClose)
		position.LowestPrice = &f

		if trailingStopLossType == ABSOLUTE {
			trailingStopLossPrice = *position.LowestPrice + *trailingStopLoss
		} else if trailingStopLossType == PERCENTAGE {
			trailingStopLossPrice = *position.LowestPrice * (1 + *trailingStopLoss)
		} else {
			panic("Invalid TrailingStopLossType")
		}
	} else {
		panic("Invalid PositionType")
	}

	return trailingStopLossPrice
}

func (position *Position) CheckForExit(
	currentClose float64,
	stopLossPrice float64,
	targetProfitPrice float64,
	trailingStopLossPrice float64,
	exitConditionResult bool,
) bool {
	stopLossTriggered := ((position.IsLong() && (currentClose < stopLossPrice)) ||
		(!position.IsLong() && (currentClose > stopLossPrice)))
	trailingStopTriggered := ((position.IsLong() && currentClose < trailingStopLossPrice) ||
		(!position.IsLong() && (currentClose > trailingStopLossPrice)))
	takeProfitTriggered := ((position.IsLong() && (currentClose >= targetProfitPrice)) ||
		(!position.IsLong() && (currentClose <= targetProfitPrice)))

	var exitReason ExitReason
	if stopLossTriggered {
		position.State = Exited
		exitReason = HIT_STOP_LOSS
	} else if trailingStopTriggered {
		position.State = Exited
		exitReason = HIT_TRAILING_STOP_LOSS
	} else if takeProfitTriggered {
		position.State = Exited
		exitReason = HIT_TAKE_PROFIT
	} else if exitConditionResult {
		position.State = Exited
		exitReason = EXIT_CONDITION_MET
	}

	if exitReason != "" {
		position.ExitReason = &exitReason
	}

	return position.IsDone()
}

func (position *Position) UpdateOnEntrySignal(
	entryPriceSignal *float64,
	entryTimeStampSignal *time.Time,
) {
	// NOTE: Check this out for Live Trading
	// position.State = Entered
	position.EntryPriceSignal = entryPriceSignal
	position.EntryTimeSignal = entryTimeStampSignal
}

func (position *Position) UpdateOnEntryOrder(
	entryPriceOrder *float64,
	entryTimeStampOrder *time.Time,
) {
	position.State = Entered
	position.EntryPriceOrder = entryPriceOrder
	position.EntryTimeOrder = entryTimeStampOrder

}

func (position *Position) UpdateOnExitSignal(
	exitPriceSignal *float64,
	exitTimeStampSignal *time.Time,
) {
	position.ExitPriceSignal = exitPriceSignal
	position.ExitTimeSignal = exitTimeStampSignal
}

func (position *Position) UpdateOnExitOrder(
	exitPriceOrder *float64,
	exitTimeStampOrder *time.Time,
) {

	// TODO: Add Slippage to exitPriceOrder
	position.ExitPriceOrder = exitPriceOrder
	position.ExitTimeOrder = exitTimeStampOrder

	// Calculate P&L - Absolute and Percentage
	position.UpdateProfitLoss(
		position.Direction,
		*position.EntryPriceOrder,
		*position.ExitPriceOrder,
	)

	// Calculate Entry Fees and Exit Fees
	var fees float64 = 20.0 // Assuming fixed entry fees for now (can be updated later)
	position.EntryFees = &fees
	position.ExitFees = &fees

	// Calculate Transaction Costs
	transactionCost := *position.EntryFees + *position.ExitFees
	position.TransactionCost = &transactionCost

	netProfitLoss := *position.ProfitLossAbsolute - *position.TransactionCost
	position.NetProfitLoss = &netProfitLoss
}

func (position *Position) UpdateProfitLoss(
	direction Direction,
	entryPrice float64,
	exitPrice float64,
) {

	var profitLossAbsolute float64
	var profitLossPercentage float64

	if direction == LONG {
		profitLossAbsolute = exitPrice - entryPrice
	} else {
		profitLossAbsolute = entryPrice - exitPrice
	}

	// Calculate Profit/Loss Percentage
	profitLossPercentage = (profitLossAbsolute / entryPrice) * 100

	position.ProfitLossAbsolute = &profitLossAbsolute
	position.ProfitLossPercentage = &profitLossPercentage
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
