package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"algo/internal/modules/backtest/domain/model"
	"algo/internal/modules/backtest/domain/repository"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Common errors
var (
	ErrRunMetricsNotFound = errors.New("run metrics not found")
)

// RunMetricsRepository implements the repository.RunMetricsRepository interface
type RunMetricsRepository struct {
	db *sql.DB
}

// NewRunMetricsRepository creates a new PostgreSQL repository for run metrics
func NewRunMetricsRepository(db *sql.DB) repository.RunMetricsRepository {
	return &RunMetricsRepository{
		db: db,
	}
}

// Create saves new run metrics in the database
func (r *RunMetricsRepository) Create(ctx context.Context, metrics *model.RunMetrics) error {
	query := `
		INSERT INTO run_metrics (
			id, run_id, total_return_pct, annualized_return, risk_adjusted_return,
			daily_returns, monthly_returns, sharpe_ratio, sortino_ratio, max_drawdown,
			max_drawdown_duration, value_at_risk, std_deviation, downside_deviation,
			total_trades, win_rate, avg_win_size, avg_loss_size, avg_trade_duration,
			profit_factor, avg_position_size, max_position_size, position_turnover,
			avg_exposure, total_commission, total_slippage, transaction_costs_pct,
			metric_json, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		metrics.ID,
		metrics.RunID,
		metrics.TotalReturnPct,
		metrics.AnnualizedReturn,
		metrics.RiskAdjustedReturn,
		pq.Array(metrics.DailyReturns),
		pq.Array(metrics.MonthlyReturns),
		metrics.SharpeRatio,
		metrics.SortinoRatio,
		metrics.MaxDrawdown,
		metrics.MaxDrawdownDuration,
		metrics.ValueAtRisk,
		metrics.StdDeviation,
		metrics.DownsideDeviation,
		metrics.TotalTrades,
		metrics.WinRate,
		metrics.AvgWinSize,
		metrics.AvgLossSize,
		metrics.AvgTradeDuration,
		metrics.ProfitFactor,
		metrics.AvgPositionSize,
		metrics.MaxPositionSize,
		metrics.PositionTurnover,
		metrics.AvgExposure,
		metrics.TotalCommission,
		metrics.TotalSlippage,
		metrics.TransactionCostsPct,
		metrics.MetricJSON,
		metrics.CreatedAt,
		metrics.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetByRunID retrieves metrics for a specific run
func (r *RunMetricsRepository) GetByRunID(ctx context.Context, runID uuid.UUID) (*model.RunMetrics, error) {
	query := `
		SELECT 
			id, run_id, total_return_pct, annualized_return, risk_adjusted_return,
			daily_returns, monthly_returns, sharpe_ratio, sortino_ratio, max_drawdown,
			max_drawdown_duration, value_at_risk, std_deviation, downside_deviation,
			total_trades, win_rate, avg_win_size, avg_loss_size, avg_trade_duration,
			profit_factor, avg_position_size, max_position_size, position_turnover,
			avg_exposure, total_commission, total_slippage, transaction_costs_pct,
			metric_json, created_at, updated_at
		FROM run_metrics
		WHERE run_id = $1
	`

	var metrics model.RunMetrics
	var dailyReturns, monthlyReturns []float64

	err := r.db.QueryRowContext(ctx, query, runID).Scan(
		&metrics.ID,
		&metrics.RunID,
		&metrics.TotalReturnPct,
		&metrics.AnnualizedReturn,
		&metrics.RiskAdjustedReturn,
		pq.Array(&dailyReturns),
		pq.Array(&monthlyReturns),
		&metrics.SharpeRatio,
		&metrics.SortinoRatio,
		&metrics.MaxDrawdown,
		&metrics.MaxDrawdownDuration,
		&metrics.ValueAtRisk,
		&metrics.StdDeviation,
		&metrics.DownsideDeviation,
		&metrics.TotalTrades,
		&metrics.WinRate,
		&metrics.AvgWinSize,
		&metrics.AvgLossSize,
		&metrics.AvgTradeDuration,
		&metrics.ProfitFactor,
		&metrics.AvgPositionSize,
		&metrics.MaxPositionSize,
		&metrics.PositionTurnover,
		&metrics.AvgExposure,
		&metrics.TotalCommission,
		&metrics.TotalSlippage,
		&metrics.TransactionCostsPct,
		&metrics.MetricJSON,
		&metrics.CreatedAt,
		&metrics.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRunMetricsNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	// Assign arrays
	metrics.DailyReturns = dailyReturns
	metrics.MonthlyReturns = monthlyReturns

	return &metrics, nil
}

// Update updates existing run metrics in the database
func (r *RunMetricsRepository) Update(ctx context.Context, metrics *model.RunMetrics) error {
	query := `
		UPDATE run_metrics
		SET 
			total_return_pct = $3,
			annualized_return = $4,
			risk_adjusted_return = $5,
			daily_returns = $6,
			monthly_returns = $7,
			sharpe_ratio = $8,
			sortino_ratio = $9,
			max_drawdown = $10,
			max_drawdown_duration = $11,
			value_at_risk = $12,
			std_deviation = $13,
			downside_deviation = $14,
			total_trades = $15,
			win_rate = $16,
			avg_win_size = $17,
			avg_loss_size = $18,
			avg_trade_duration = $19,
			profit_factor = $20,
			avg_position_size = $21,
			max_position_size = $22,
			position_turnover = $23,
			avg_exposure = $24,
			total_commission = $25,
			total_slippage = $26,
			transaction_costs_pct = $27,
			metric_json = $28,
			updated_at = $29
		WHERE id = $1 AND run_id = $2
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		metrics.ID,
		metrics.RunID,
		metrics.TotalReturnPct,
		metrics.AnnualizedReturn,
		metrics.RiskAdjustedReturn,
		pq.Array(metrics.DailyReturns),
		pq.Array(metrics.MonthlyReturns),
		metrics.SharpeRatio,
		metrics.SortinoRatio,
		metrics.MaxDrawdown,
		metrics.MaxDrawdownDuration,
		metrics.ValueAtRisk,
		metrics.StdDeviation,
		metrics.DownsideDeviation,
		metrics.TotalTrades,
		metrics.WinRate,
		metrics.AvgWinSize,
		metrics.AvgLossSize,
		metrics.AvgTradeDuration,
		metrics.ProfitFactor,
		metrics.AvgPositionSize,
		metrics.MaxPositionSize,
		metrics.PositionTurnover,
		metrics.AvgExposure,
		metrics.TotalCommission,
		metrics.TotalSlippage,
		metrics.TransactionCostsPct,
		metrics.MetricJSON,
		metrics.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrRunMetricsNotFound
	}

	return nil
}

// Delete removes run metrics from the database
func (r *RunMetricsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM run_metrics
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrRunMetricsNotFound
	}

	return nil
}
