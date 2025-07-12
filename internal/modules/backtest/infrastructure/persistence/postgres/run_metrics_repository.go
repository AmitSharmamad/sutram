package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"algo/internal/modules/backtest/domain/model"
	"algo/internal/modules/backtest/domain/repository"
	mathUtil "algo/internal/util/math"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

// clampDecimal10_4 clamps a float64 to fit in DECIMAL(10, 4) database column
func clampDecimal10_4(value float64) float64 {
	return mathUtil.ClampFloat(value, -999999.9999, 999999.9999)
}

// clampDecimal18_8 clamps a float64 to fit in DECIMAL(18, 8) database column
func clampDecimal18_8(value float64) float64 {
	return mathUtil.ClampFloat(value, -9999999999.99999999, 9999999999.99999999)
}

// clampDecimal10_4Array clamps all values in a float64 slice to fit in DECIMAL(10, 4) database column
func clampDecimal10_4Array(values []float64) []float64 {
	if values == nil {
		return nil
	}
	result := make([]float64, len(values))
	for i, val := range values {
		result[i] = clampDecimal10_4(val)
	}
	return result
}

// clampMetrics applies clamping to all metrics fields to prevent database precision errors
func (r *RunMetricsRepository) clampMetrics(metrics *model.RunMetrics) {
	// Clamp DECIMAL(10, 4) fields
	metrics.TotalReturnPct = clampDecimal10_4(metrics.TotalReturnPct)
	metrics.AnnualizedReturn = clampDecimal10_4(metrics.AnnualizedReturn)
	if metrics.RiskAdjustedReturn != nil {
		clamped := clampDecimal10_4(*metrics.RiskAdjustedReturn)
		metrics.RiskAdjustedReturn = &clamped
	}

	// Clamp arrays
	metrics.DailyReturns = clampDecimal10_4Array(metrics.DailyReturns)
	metrics.MonthlyReturns = clampDecimal10_4Array(metrics.MonthlyReturns)

	// Clamp more DECIMAL(10, 4) fields
	metrics.SharpeRatio = clampDecimal10_4(metrics.SharpeRatio)
	metrics.SortinoRatio = clampDecimal10_4(metrics.SortinoRatio)
	metrics.MaxDrawdown = clampDecimal10_4(metrics.MaxDrawdown)
	metrics.MaxDrawdownDuration = clampDecimal10_4(metrics.MaxDrawdownDuration)
	if metrics.ValueAtRisk != nil {
		clamped := clampDecimal10_4(*metrics.ValueAtRisk)
		metrics.ValueAtRisk = &clamped
	}
	metrics.StdDeviation = clampDecimal10_4(metrics.StdDeviation)
	if metrics.DownsideDeviation != nil {
		clamped := clampDecimal10_4(*metrics.DownsideDeviation)
		metrics.DownsideDeviation = &clamped
	}

	// Clamp trade metrics
	metrics.WinRate = clampDecimal10_4(metrics.WinRate)
	if metrics.AvgTradeDuration != nil {
		clamped := clampDecimal10_4(*metrics.AvgTradeDuration)
		metrics.AvgTradeDuration = &clamped
	}
	if metrics.ProfitFactor != nil {
		clamped := clampDecimal10_4(*metrics.ProfitFactor)
		metrics.ProfitFactor = &clamped
	}
	if metrics.PositionTurnover != nil {
		clamped := clampDecimal10_4(*metrics.PositionTurnover)
		metrics.PositionTurnover = &clamped
	}
	metrics.TransactionCostsPct = clampDecimal10_4(metrics.TransactionCostsPct)

	// Clamp DECIMAL(18, 8) fields
	if metrics.AvgWinSize != nil {
		clamped := clampDecimal18_8(*metrics.AvgWinSize)
		metrics.AvgWinSize = &clamped
	}
	if metrics.AvgLossSize != nil {
		clamped := clampDecimal18_8(*metrics.AvgLossSize)
		metrics.AvgLossSize = &clamped
	}
	if metrics.AvgPositionSize != nil {
		clamped := clampDecimal18_8(*metrics.AvgPositionSize)
		metrics.AvgPositionSize = &clamped
	}
	if metrics.MaxPositionSize != nil {
		clamped := clampDecimal18_8(*metrics.MaxPositionSize)
		metrics.MaxPositionSize = &clamped
	}
	if metrics.AvgExposure != nil {
		clamped := clampDecimal18_8(*metrics.AvgExposure)
		metrics.AvgExposure = &clamped
	}
	metrics.TotalCommission = clampDecimal18_8(metrics.TotalCommission)
	metrics.TotalSlippage = clampDecimal18_8(metrics.TotalSlippage)
}

// Create saves new run metrics in the database
func (r *RunMetricsRepository) Create(ctx context.Context, metrics *model.RunMetrics) error {
	// Apply clamping to prevent database precision errors
	r.clampMetrics(metrics)

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
	// Apply clamping to prevent database precision errors
	r.clampMetrics(metrics)

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

// BulkCreate saves multiple run metrics in the database in a single operation
func (r *RunMetricsRepository) BulkCreate(ctx context.Context, metrics []*model.RunMetrics) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback()

	// Prepare the statement
	stmt, err := tx.PrepareContext(ctx, `
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
	`)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer stmt.Close()

	// Execute for each metric
	for _, metric := range metrics {
		// Apply clamping to prevent database precision errors
		r.clampMetrics(metric)

		_, err := stmt.ExecContext(
			ctx,
			metric.ID,
			metric.RunID,
			metric.TotalReturnPct,
			metric.AnnualizedReturn,
			metric.RiskAdjustedReturn,
			pq.Array(metric.DailyReturns),
			pq.Array(metric.MonthlyReturns),
			metric.SharpeRatio,
			metric.SortinoRatio,
			metric.MaxDrawdown,
			metric.MaxDrawdownDuration,
			metric.ValueAtRisk,
			metric.StdDeviation,
			metric.DownsideDeviation,
			metric.TotalTrades,
			metric.WinRate,
			metric.AvgWinSize,
			metric.AvgLossSize,
			metric.AvgTradeDuration,
			metric.ProfitFactor,
			metric.AvgPositionSize,
			metric.MaxPositionSize,
			metric.PositionTurnover,
			metric.AvgExposure,
			metric.TotalCommission,
			metric.TotalSlippage,
			metric.TransactionCostsPct,
			metric.MetricJSON,
			metric.CreatedAt,
			metric.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}
