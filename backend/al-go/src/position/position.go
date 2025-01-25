package position

import (
	"algo/src/stock"
	"time"
)

const (
	NOT_ENTERED = "NOT_ENTERED"
	ENTERED     = "ENTERED"
	EXITED      = "EXITED"
)

type Position struct {
	ID int `json:"id" db:"id"`

	RunID int `json:"run_id" db:"run_id"`

	State string `json:"state" db:"state"`

	EntryPrice float64   `json:"entry_price" db:"entry_price"`
	EntryTime  time.Time `json:"entry_time" db:"entry_time"`

	ExitPrice  float64   `json:"exit_price" db:"exit_price"`
	ExitTime   time.Time `json:"exit_time" db:"exit_time"`
	ExitReason string    `json:"exit_reason" db:"exit_reason"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewPosition() *Position {
	return &Position{}
}

// PositionLog represents a log entry for a trading position.
type PositionLog struct {
	ID int `json:"id" db:"id"`

	PositionID int `json:"position_id" db:"position_id"`

	State string `json:"state" db:"state"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (position *Position) TryEnter(data stock.StockData) bool {
	return true
}

func (position *Position) TryExit(data stock.StockData) bool {
	return true
}

func (position *Position) IsDone() bool {
	return position.State == EXITED
}

func (position *Position) Eval(data stock.StockData) {
	switch position.State {
	case NOT_ENTERED:
		if position.TryEnter(data) {
			position.State = ENTERED
		}
	case ENTERED:
		if position.TryExit(data) {
			position.State = EXITED
		}
	default:
		panic("Error in position Eval")
	}
}
