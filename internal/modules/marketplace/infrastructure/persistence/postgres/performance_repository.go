// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/infrastructure/persistence/postgres/performance_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// PerformanceRepository implements repository.PerformanceRepository
type PerformanceRepository struct {
	db *sql.DB
}

// NewPerformanceRepository creates a new performance repository
func NewPerformanceRepository(db *sql.DB) repository.PerformanceRepository {
	return &PerformanceRepository{
		db: db,
	}
}

// Create creates performance metrics for a listing
func (r *PerformanceRepository) Create(ctx context.Context, performance *model.Performance) error {
	query := `
		INSERT INTO marketplace_performance (
			listing_id, sharpe_ratio, sortino_ratio, total_return_pct,
			annualized_return, max_drawdown, win_rate, profit_factor,
			backtest_period_days, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		performance.ListingID,
		performance.SharpeRatio,
		performance.SortinoRatio,
		performance.TotalReturnPct,
		performance.AnnualizedReturn,
		performance.MaxDrawdown,
		performance.WinRate,
		performance.ProfitFactor,
		performance.BacktestPeriodDays,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create performance metrics: %w", err)
	}

	return nil
}

// FindByListingID returns performance metrics for a listing
func (r *PerformanceRepository) FindByListingID(ctx context.Context, listingID uuid.UUID) (*model.Performance, error) {
	query := `
		SELECT 
			listing_id, sharpe_ratio, sortino_ratio, total_return_pct,
			annualized_return, max_drawdown, win_rate, profit_factor,
			backtest_period_days, updated_at
		FROM marketplace_performance
		WHERE listing_id = $1
	`

	var performance model.Performance
	err := r.db.QueryRowContext(ctx, query, listingID).Scan(
		&performance.ListingID,
		&performance.SharpeRatio,
		&performance.SortinoRatio,
		&performance.TotalReturnPct,
		&performance.AnnualizedReturn,
		&performance.MaxDrawdown,
		&performance.WinRate,
		&performance.ProfitFactor,
		&performance.BacktestPeriodDays,
		&performance.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find performance: %w", err)
	}

	return &performance, nil
}

// Update updates performance metrics for a listing
func (r *PerformanceRepository) Update(ctx context.Context, performance *model.Performance) error {
	query := `
		UPDATE marketplace_performance SET
			sharpe_ratio = $1,
			sortino_ratio = $2,
			total_return_pct = $3,
			annualized_return = $4,
			max_drawdown = $5,
			win_rate = $6,
			profit_factor = $7,
			backtest_period_days = $8,
			updated_at = $9
		WHERE listing_id = $10
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		performance.SharpeRatio,
		performance.SortinoRatio,
		performance.TotalReturnPct,
		performance.AnnualizedReturn,
		performance.MaxDrawdown,
		performance.WinRate,
		performance.ProfitFactor,
		performance.BacktestPeriodDays,
		time.Now(),
		performance.ListingID,
	)

	if err != nil {
		return fmt.Errorf("failed to update performance metrics: %w", err)
	}

	return nil
}

// Delete deletes performance metrics for a listing
func (r *PerformanceRepository) Delete(ctx context.Context, listingID uuid.UUID) error {
	query := `DELETE FROM marketplace_performance WHERE listing_id = $1`
	_, err := r.db.ExecContext(ctx, query, listingID)
	if err != nil {
		return fmt.Errorf("failed to delete performance metrics: %w", err)
	}

	return nil
}
