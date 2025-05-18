// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/dto/performance_dto.go
package dto

// UpdatePerformanceRequest represents a request to update performance metrics
type UpdatePerformanceRequest struct {
	SharpeRatio        float64 `json:"sharpe_ratio" binding:"required"`
	SortinoRatio       float64 `json:"sortino_ratio" binding:"required"`
	TotalReturnPct     float64 `json:"total_return_pct" binding:"required"`
	AnnualizedReturn   float64 `json:"annualized_return" binding:"required"`
	MaxDrawdown        float64 `json:"max_drawdown" binding:"required"`
	WinRate            float64 `json:"win_rate" binding:"required"`
	ProfitFactor       float64 `json:"profit_factor" binding:"required"`
	BacktestPeriodDays int     `json:"backtest_period_days" binding:"required,min=1"`
}

// PerformanceResponse represents performance metrics
type PerformanceResponse struct {
	SharpeRatio        float64 `json:"sharpe_ratio"`
	SortinoRatio       float64 `json:"sortino_ratio"`
	TotalReturnPct     float64 `json:"total_return_pct"`
	AnnualizedReturn   float64 `json:"annualized_return"`
	MaxDrawdown        float64 `json:"max_drawdown"`
	WinRate            float64 `json:"win_rate"`
	ProfitFactor       float64 `json:"profit_factor"`
	BacktestPeriodDays int     `json:"backtest_period_days"`
}
