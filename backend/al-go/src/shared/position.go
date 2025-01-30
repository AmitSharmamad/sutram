package shared

import (
	"time"

	"algo/src/shared/algo"

	"github.com/google/uuid"
)

const (
	NotEntered RunState = "NOT_ENTERED"
	Entered    RunState = "ENTERED"
	Exited     RunState = "EXITED"
)

type RunState string

type Runs []Run

type Positions []Position

type PositionLogs []PositionLog

type RunPositions []RunPosition

type Run struct {
	ID              uuid.UUID `json:"id" db:"id"`
	StrategyID      uuid.UUID `json:"strategy_id" db:"strategy_id"`
	State           RunState  `json:"state" db:"state"`
	Qty             float64   `json:"qty" db:"qty"`
	PeriodStartTime time.Time `json:"period_start_time" db:"period_start_time"`
	PeriodEndTime   time.Time `json:"period_end_time" db:"period_end_time"`
	StartedTime     time.Time `json:"started_time" db:"started_time"`
	EndedTime       time.Time `json:"ended_time" db:"ended_time"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func NewRun() *Run {
	return new(Run)
}

type PositionType string

const (
	LATEST PositionType = "LATEST"
	NEW    PositionType = "NEW"
)

func (run *Run) GetPositionType() PositionType {
	switch run.State {
	case "START":
		run.State = "IN_PROGRESS"
		return NEW
	case "IN_PROGRESS":
		// latest := position.GetLatest()
		// if latest.State == "EXITED" {
		// 	// create a new position
		// 	return position.NewPosition()
		// }
		// return latest
		return LATEST
	case "ENDED":
		panic("Run ended, but loop going on")
	}
	panic("Error in run get position")
}

type Position struct {
	ID         string    `json:"id" db:"id"`
	RunID      int64     `json:"run_id" db:"run_id"`
	State      RunState  `json:"state" db:"state"`
	EntryPrice float64   `json:"entry_price" db:"entry_price"`
	EntryTime  time.Time `json:"entry_time" db:"entry_time"`
	ExitPrice  float64   `json:"exit_price" db:"exit_price"`
	ExitTime   time.Time `json:"exit_time" db:"exit_time"`
	ExitReason string    `json:"exit_reason" db:"exit_reason"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func NewPosition() *Position {
	return new(Position)
}

func (position *Position) TryEnter(
	strategy *Strategy,
	data *StockData,
) bool {
	return algo.NewExpressionStatement(strategy.EntryExpression).Eval(data)
}

func (position *Position) TryExit(
	strategy *Strategy,
	data *StockData,
) bool {
	// TODO: also check for stop loss, update trailing stop loss
	return algo.NewExpressionStatement(strategy.ExitExpression).Eval(data)
}

func (position *Position) IsDone() bool {
	return position.State == Exited
}

func (position *Position) Eval(
	strategy *Strategy,
	data *StockData,
) {
	switch position.State {
	case NotEntered:
		if position.TryEnter(strategy, data) {
			position.State = Entered
		}
	case Entered:
		if position.TryExit(strategy, data) {
			position.State = Exited
		}
	default:
		panic("Error in position Eval")
	}
}

type RunPosition struct {
	Run       *Run       `json:"run"`
	Positions *Positions `json:"positions"`
}

func NewRunPosition() *RunPosition {
	return new(RunPosition)
}

type PositionLog struct {
	ID         string    `json:"id" db:"id"`
	PositionID int64     `json:"position_id" db:"position_id"`
	State      string    `json:"state" db:"state"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

func NewPositionLog() *PositionLog {
	return new(PositionLog)
}
