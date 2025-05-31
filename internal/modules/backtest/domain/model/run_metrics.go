package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// RunMetrics represents a set of performance metrics for a run
type RunMetrics struct {
	ID                  uuid.UUID   `json:"id" db:"id"`
	RunID               uuid.UUID   `json:"run_id" db:"run_id"`
	TotalReturnPct      float64     `json:"total_return_pct" db:"total_return_pct"`
	AnnualizedReturn    float64     `json:"annualized_return" db:"annualized_return"`
	RiskAdjustedReturn  *float64    `json:"risk_adjusted_return,omitempty" db:"risk_adjusted_return"`
	DailyReturns        []float64   `json:"daily_returns" db:"daily_returns"`
	MonthlyReturns      []float64   `json:"monthly_returns" db:"monthly_returns"`
	SharpeRatio         float64     `json:"sharpe_ratio" db:"sharpe_ratio"`
	SortinoRatio        float64     `json:"sortino_ratio" db:"sortino_ratio"`
	MaxDrawdown         float64     `json:"max_drawdown" db:"max_drawdown"`
	MaxDrawdownDuration float64     `json:"max_drawdown_duration" db:"max_drawdown_duration"`
	ValueAtRisk         *float64    `json:"value_at_risk,omitempty" db:"value_at_risk"`
	StdDeviation        float64     `json:"std_deviation" db:"std_deviation"`
	DownsideDeviation   *float64    `json:"downside_deviation,omitempty" db:"downside_deviation"`
	TotalTrades         int         `json:"total_trades" db:"total_trades"`
	WinRate             float64     `json:"win_rate" db:"win_rate"`
	AvgWinSize          *float64    `json:"avg_win_size,omitempty" db:"avg_win_size"`
	AvgLossSize         *float64    `json:"avg_loss_size,omitempty" db:"avg_loss_size"`
	AvgTradeDuration    *float64    `json:"avg_trade_duration,omitempty" db:"avg_trade_duration"`
	ProfitFactor        *float64    `json:"profit_factor,omitempty" db:"profit_factor"`
	AvgPositionSize     *float64    `json:"avg_position_size,omitempty" db:"avg_position_size"`
	MaxPositionSize     *float64    `json:"max_position_size,omitempty" db:"max_position_size"`
	PositionTurnover    *float64    `json:"position_turnover,omitempty" db:"position_turnover"`
	AvgExposure         *float64    `json:"avg_exposure,omitempty" db:"avg_exposure"`
	TotalCommission     float64     `json:"total_commission" db:"total_commission"`
	TotalSlippage       float64     `json:"total_slippage" db:"total_slippage"`
	TransactionCostsPct float64     `json:"transaction_costs_pct" db:"transaction_costs_pct"`
	MetricJSON          string      `json:"metric_json" db:"metric_json"`
	metricData          interface{} // Cached data from MetricJSON
	CreatedAt           time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at" db:"updated_at"`
}

// NewRunMetrics creates a new run metrics for the given run
func NewRunMetrics(runID uuid.UUID) *RunMetrics {
	now := time.Now()
	return &RunMetrics{
		ID:             uuid.New(),
		RunID:          runID,
		DailyReturns:   []float64{},
		MonthlyReturns: []float64{},
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// GetMetricData parses and returns the JSON metrics data
func (rm *RunMetrics) GetMetricData(data interface{}) error {
	if rm.MetricJSON == "" {
		return nil
	}

	return json.Unmarshal([]byte(rm.MetricJSON), data)
}

// SetMetricData sets the metric JSON data
func (rm *RunMetrics) SetMetricData(data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	rm.MetricJSON = string(bytes)
	rm.metricData = data
	return nil
}

// ReturnsData represents the different types of returns calculated from positions
type ReturnsData struct {
	DailyReturns     []float64 // Array of daily returns
	MonthlyReturns   []float64 // Array of monthly returns
	AnnualizedReturn float64   // Annualized return value
}
