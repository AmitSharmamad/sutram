package dto

import (
	"time"

	"github.com/google/uuid"
)

// RunMetricsDTO represents a data transfer object for run metrics
type RunMetricsDTO struct {
	ID                  uuid.UUID `json:"id"`
	RunID               uuid.UUID `json:"run_id"`
	TotalReturnPct      float64   `json:"total_return_pct"`
	AnnualizedReturn    float64   `json:"annualized_return"`
	RiskAdjustedReturn  float64   `json:"risk_adjusted_return"`
	DailyReturns        []float64 `json:"daily_returns"`
	MonthlyReturns      []float64 `json:"monthly_returns"`
	SharpeRatio         float64   `json:"sharpe_ratio"`
	SortinoRatio        float64   `json:"sortino_ratio"`
	MaxDrawdown         float64   `json:"max_drawdown"`
	MaxDrawdownDuration float64   `json:"max_drawdown_duration"`
	ValueAtRisk         float64   `json:"value_at_risk"`
	StdDeviation        float64   `json:"std_deviation"`
	DownsideDeviation   float64   `json:"downside_deviation"`
	TotalTrades         int       `json:"total_trades"`
	WinRate             float64   `json:"win_rate"`
	AvgWinSize          float64   `json:"avg_win_size"`
	AvgLossSize         float64   `json:"avg_loss_size"`
	AvgTradeDuration    float64   `json:"avg_trade_duration"`
	ProfitFactor        float64   `json:"profit_factor"`
	AvgPositionSize     float64   `json:"avg_position_size"`
	MaxPositionSize     float64   `json:"max_position_size"`
	PositionTurnover    float64   `json:"position_turnover"`
	AvgExposure         float64   `json:"avg_exposure"`
	TotalCommission     float64   `json:"total_commission"`
	TotalSlippage       float64   `json:"total_slippage"`
	TransactionCostsPct float64   `json:"transaction_costs_pct"`
	MetricJSON          string    `json:"metric_json"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// CreateRunMetricsRequest represents a request to create new run metrics
type CreateRunMetricsRequest struct {
	RunID               uuid.UUID `json:"run_id" binding:"required"`
	TotalReturnPct      float64   `json:"total_return_pct"`
	AnnualizedReturn    float64   `json:"annualized_return"`
	RiskAdjustedReturn  float64   `json:"risk_adjusted_return"`
	DailyReturns        []float64 `json:"daily_returns"`
	MonthlyReturns      []float64 `json:"monthly_returns"`
	SharpeRatio         float64   `json:"sharpe_ratio"`
	SortinoRatio        float64   `json:"sortino_ratio"`
	MaxDrawdown         float64   `json:"max_drawdown"`
	MaxDrawdownDuration float64   `json:"max_drawdown_duration"`
	ValueAtRisk         float64   `json:"value_at_risk"`
	StdDeviation        float64   `json:"std_deviation"`
	DownsideDeviation   float64   `json:"downside_deviation"`
	TotalTrades         int       `json:"total_trades"`
	WinRate             float64   `json:"win_rate"`
	AvgWinSize          float64   `json:"avg_win_size"`
	AvgLossSize         float64   `json:"avg_loss_size"`
	AvgTradeDuration    float64   `json:"avg_trade_duration"`
	ProfitFactor        float64   `json:"profit_factor"`
	AvgPositionSize     float64   `json:"avg_position_size"`
	MaxPositionSize     float64   `json:"max_position_size"`
	PositionTurnover    float64   `json:"position_turnover"`
	AvgExposure         float64   `json:"avg_exposure"`
	TotalCommission     float64   `json:"total_commission"`
	TotalSlippage       float64   `json:"total_slippage"`
	TransactionCostsPct float64   `json:"transaction_costs_pct"`
	MetricJSON          string    `json:"metric_json"`
}
