package shared

import (
	"time"

	"github.com/google/uuid"
)

const (
	NotEntered PositionState = "NOT_ENTERED"
	Entered    PositionState = "ENTERED"
	Exited     PositionState = "EXITED"

	NotStarted RunState = "NOT_STARTED"
	Started    RunState = "STARTED"
	Ended      RunState = "ENDED"
)

type RunState string

type PositionState string

const (
	BackTest   RunType = "BACKTEST"
	PaperTrade RunType = "PAPER_TRADE"
	LiveTrade  RunType = "LIVE_TRADE"
)

type RunType string

type Runs []Run

type Run struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	RunGroupID      *uuid.UUID `json:"run_group_id" db:"run_group_id"`
	StrategyID      uuid.UUID  `json:"strategy_id" db:"strategy_id"`
	Exchange        string     `json:"exchange" db:"exchange"`
	Instrument      string     `json:"instrument" db:"instrument"`
	RunType         RunType    `json:"type" db:"type"`
	State           RunState   `json:"state" db:"state"`
	Qty             int        `json:"qty" db:"qty"`
	PeriodStartTime *time.Time `json:"period_start_time" db:"period_start_time"`
	PeriodEndTime   *time.Time `json:"period_end_time" db:"period_end_time"`
	StartedTime     *time.Time `json:"started_time" db:"started_time"`
	EndedTime       *time.Time `json:"ended_time" db:"ended_time"`
	Audit
}

func NewRun() *Run {
	return new(Run)
}

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

type RunMetrics struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	RunID               uuid.UUID `json:"run_id" db:"run_id"`
	TotalReturnPCT      float64   `json:"total_return_pct" db:"total_return_pct"`
	AnnualizedReturn    float64   `json:"annualized_return" db:"annualized_return"`
	RiskAdjustedReturn  float64   `json:"risk_adjusted_return" db:"risk_adjusted_return"`
	DailyReturns        float64   `json:"daily_returns" db:"daily_returns"`
	MonthlyReturns      float64   `json:"monthly_returns" db:"monthly_returns"`
	SharpeRatio         float64   `json:"sharpe_ratio" db:"sharpe_ratio"`
	SortinoRatio        float64   `json:"sortino_ratio" db:"sortino_ratio"`
	MaxDrawdown         float64   `json:"max_drawdown" db:"max_drawdown"`
	MaxDrawdownDuration float64   `json:"max_drawdown_duration" db:"max_drawdown_duration"`
	ValueAtRisk         float64   `json:"value_at_risk" db:"value_at_risk"`
	StdDeviation        float64   `json:"std_deviation" db:"std_deviation"`
	DownsideDeviation   float64   `json:"downside_deviation" db:"downside_deviation"`
	TotalTrades         int64     `json:"total_trades" db:"total_trades"`
	WinRate             float64   `json:"win_rate" db:"win_rate"`
	AvgWinSize          float64   `json:"avg_win_size" db:"avg_win_size"`
	AvgLossSize         float64   `json:"avg_loss_size" db:"avg_loss_size"`
	AvgTradeDuration    float64   `json:"avg_trade_duration" db:"avg_trade_duration"`
	ProfitFactor        float64   `json:"profit_factor" db:"profit_factor"`
	AvgPositionSize     float64   `json:"avg_position_size" db:"avg_position_size"`
	MaxPositionSize     float64   `json:"max_position_size" db:"max_position_size"`
	PositionTurnover    float64   `json:"position_turnover" db:"position_turnover"`
	AvgExposure         float64   `json:"avg_exposure" db:"avg_exposure"`
	TotalCommission     float64   `json:"total_commission" db:"total_commission"`
	TotalSlippage       float64   `json:"total_slippage" db:"total_slippage"`
	TransactionCostsPCT float64   `json:"transaction_costs_pct" db:"transaction_costs_pct"`
	MetricJson          string    `json:"metric_json" db:"metric_json"`
}
