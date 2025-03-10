package shared

import (
	"context"
	"time"

	"algo/src/util"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Strategy struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	// Position
	Name           string         `json:"name" db:"name"`
	Description    string         `json:"description" db:"description"`
	PositionType   string         `json:"position_type" db:"position_type"`
	Instruments    pq.StringArray `json:"instruments" db:"instruments"`
	ChartType      string         `json:"chart_type" db:"chart_type"`
	ChartTimeframe string         `json:"chart_timeframe" db:"chart_timeframe"`
	ChartUnknown   string         `json:"chart_unknown" db:"chart_unknown"`
	Qty            int32          `json:"qty" db:"qty"`

	// Entry Condition
	EntryExpression string `json:"entry_expression" db:"entry_expression"`

	// Exit Condition
	ExitExpression   string  `json:"exit_expression" db:"exit_expression"`
	TargetProfit     float64 `json:"target_profit" db:"target_profit"`
	StopLoss         float64 `json:"stop_loss" db:"stop_loss"`
	TrailingStopLoss float64 `json:"trailing_stop_loss" db:"trailing_stop_loss"`
	ExitType         string  `json:"exit_type" db:"exit_type"`

	// Backtest params
	Capital        float64   `json:"capital" db:"capital"`
	MaxTrades      int32     `json:"max_trades" db:"max_trades"`
	EntryStartTime time.Time `json:"entry_start_time" db:"entry_start_time"`
	EntryEndTime   time.Time `json:"entry_end_time" db:"entry_end_time"`
	StartDate      time.Time `json:"start_date" db:"start_date"`
	EndDate        time.Time `json:"end_date" db:"end_date"`
	IsActive       bool      `json:"is_active" db:"is_active"`

	// Audit
	Audit
}

type Strategies []Strategy

func (strategy *Strategy) GetEntryDataRange() int32 {
	return 0
	// return NewExpressionStatement(
	// 	strategy.EntryExpression,
	// ).GetDataRange()
}

func (strategy *Strategy) GetExitDataRange() int32 {
	return 0
	// return NewExpressionStatement(
	// 	strategy.ExitExpression,
	// ).GetDataRange()
}

func (strategy *Strategy) TryEnter(data StockData) bool {
	// TODO: handle this
	return true
}

func (strategy *Strategy) TryExit(data StockData) bool {
	// TODO: handle this
	return true
}

func (strategy *Strategy) Validate() {
	path := "../data/llm.json"

	exp := util.GetFileContent(path)

	util.Validate(exp)
}

type StrategyService interface {
	Get(ctx context.Context, id uuid.UUID) *Strategy
	Create(ctx context.Context, strategy *Strategy,) (*Strategy, error)
	GetData() StockData
	GetAll(ctx context.Context) (*Strategies, error)
	Validate(strategy *Strategy) bool
}
