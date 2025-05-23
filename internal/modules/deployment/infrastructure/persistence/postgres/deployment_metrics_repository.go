package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/deployment/domain/model"
	"algo/internal/modules/deployment/domain/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PgDeploymentMetricsRepository implements the deployment metrics repository interface
type PgDeploymentMetricsRepository struct {
	db *sqlx.DB
}

// NewPgDeploymentMetricsRepository creates a new Postgres deployment metrics repository
func NewPgDeploymentMetricsRepository(db *sqlx.DB) repository.DeploymentMetricsRepository {
	return &PgDeploymentMetricsRepository{
		db: db,
	}
}

// Create saves new metrics in the database
func (r *PgDeploymentMetricsRepository) Create(ctx context.Context, metrics *model.DeploymentMetrics) error {
	query := `
		INSERT INTO deployment_metrics (
			deployment_id, metric_date, total_trades, win_trades, loss_trades,
			total_pnl, max_drawdown, sharpe_ratio, sortino_ratio, win_rate,
			avg_win, avg_loss, metadata, created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
		ON CONFLICT (deployment_id, metric_date) DO UPDATE SET
			total_trades = EXCLUDED.total_trades,
			win_trades = EXCLUDED.win_trades,
			loss_trades = EXCLUDED.loss_trades,
			total_pnl = EXCLUDED.total_pnl,
			max_drawdown = EXCLUDED.max_drawdown,
			sharpe_ratio = EXCLUDED.sharpe_ratio,
			sortino_ratio = EXCLUDED.sortino_ratio,
			win_rate = EXCLUDED.win_rate,
			avg_win = EXCLUDED.avg_win,
			avg_loss = EXCLUDED.avg_loss,
			metadata = EXCLUDED.metadata,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(
		ctx,
		query,
		metrics.DeploymentID,
		metrics.MetricDate,
		metrics.TotalTrades,
		metrics.WinTrades,
		metrics.LossTrades,
		metrics.TotalPNL,
		metrics.MaxDrawdown,
		metrics.SharpeRatio,
		metrics.SortinoRatio,
		metrics.WinRate,
		metrics.AvgWin,
		metrics.AvgLoss,
		metrics.Metadata,
		metrics.CreatedAt,
		metrics.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to create deployment metrics: %w", err)
	}

	metrics.ID = id
	return nil
}

// GetByID retrieves metrics by ID
func (r *PgDeploymentMetricsRepository) GetByID(ctx context.Context, id int64) (*model.DeploymentMetrics, error) {
	query := `
		SELECT 
			id, deployment_id, metric_date, total_trades, win_trades, loss_trades,
			total_pnl, max_drawdown, sharpe_ratio, sortino_ratio, win_rate,
			avg_win, avg_loss, metadata, created_at, updated_at
		FROM deployment_metrics
		WHERE id = $1
	`

	var metrics model.DeploymentMetrics
	err := r.db.GetContext(ctx, &metrics, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get deployment metrics by ID: %w", err)
	}

	return &metrics, nil
}

// GetByDeploymentID retrieves all metrics for a specific deployment
func (r *PgDeploymentMetricsRepository) GetByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*model.DeploymentMetrics, error) {
	query := `
		SELECT 
			id, deployment_id, metric_date, total_trades, win_trades, loss_trades,
			total_pnl, max_drawdown, sharpe_ratio, sortino_ratio, win_rate,
			avg_win, avg_loss, metadata, created_at, updated_at
		FROM deployment_metrics
		WHERE deployment_id = $1
		ORDER BY metric_date DESC
	`

	var metrics []*model.DeploymentMetrics
	err := r.db.SelectContext(ctx, &metrics, query, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment metrics by deployment ID: %w", err)
	}

	return metrics, nil
}

// GetByDeploymentIDAndDate retrieves metrics for a specific deployment and date
func (r *PgDeploymentMetricsRepository) GetByDeploymentIDAndDate(ctx context.Context, deploymentID uuid.UUID, date time.Time) (*model.DeploymentMetrics, error) {
	query := `
		SELECT 
			id, deployment_id, metric_date, total_trades, win_trades, loss_trades,
			total_pnl, max_drawdown, sharpe_ratio, sortino_ratio, win_rate,
			avg_win, avg_loss, metadata, created_at, updated_at
		FROM deployment_metrics
		WHERE deployment_id = $1 AND metric_date::date = $2::date
	`

	var metrics model.DeploymentMetrics
	err := r.db.GetContext(ctx, &metrics, query, deploymentID, date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get deployment metrics by deployment ID and date: %w", err)
	}

	return &metrics, nil
}

// GetLatestByDeploymentID retrieves the latest metrics for a specific deployment
func (r *PgDeploymentMetricsRepository) GetLatestByDeploymentID(ctx context.Context, deploymentID uuid.UUID) (*model.DeploymentMetrics, error) {
	query := `
		SELECT 
			id, deployment_id, metric_date, total_trades, win_trades, loss_trades,
			total_pnl, max_drawdown, sharpe_ratio, sortino_ratio, win_rate,
			avg_win, avg_loss, metadata, created_at, updated_at
		FROM deployment_metrics
		WHERE deployment_id = $1
		ORDER BY metric_date DESC
		LIMIT 1
	`

	var metrics model.DeploymentMetrics
	err := r.db.GetContext(ctx, &metrics, query, deploymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest deployment metrics by deployment ID: %w", err)
	}

	return &metrics, nil
}

// Update updates existing metrics in the database
func (r *PgDeploymentMetricsRepository) Update(ctx context.Context, metrics *model.DeploymentMetrics) error {
	query := `
		UPDATE deployment_metrics
		SET
			total_trades = $1,
			win_trades = $2,
			loss_trades = $3,
			total_pnl = $4,
			max_drawdown = $5,
			sharpe_ratio = $6,
			sortino_ratio = $7,
			win_rate = $8,
			avg_win = $9,
			avg_loss = $10,
			metadata = $11,
			updated_at = $12
		WHERE id = $13
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		metrics.TotalTrades,
		metrics.WinTrades,
		metrics.LossTrades,
		metrics.TotalPNL,
		metrics.MaxDrawdown,
		metrics.SharpeRatio,
		metrics.SortinoRatio,
		metrics.WinRate,
		metrics.AvgWin,
		metrics.AvgLoss,
		metrics.Metadata,
		metrics.UpdatedAt,
		metrics.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update deployment metrics: %w", err)
	}

	return nil
}

// GetMetricsTimeRange retrieves metrics for a deployment within a date range
func (r *PgDeploymentMetricsRepository) GetMetricsTimeRange(ctx context.Context, deploymentID uuid.UUID, startDate, endDate time.Time) ([]*model.DeploymentMetrics, error) {
	query := `
		SELECT 
			id, deployment_id, metric_date, total_trades, win_trades, loss_trades,
			total_pnl, max_drawdown, sharpe_ratio, sortino_ratio, win_rate,
			avg_win, avg_loss, metadata, created_at, updated_at
		FROM deployment_metrics
		WHERE deployment_id = $1 AND metric_date >= $2 AND metric_date <= $3
		ORDER BY metric_date
	`

	var metrics []*model.DeploymentMetrics
	err := r.db.SelectContext(ctx, &metrics, query, deploymentID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics within time range: %w", err)
	}

	return metrics, nil
}

// AggregateMetrics calculates aggregate metrics for a deployment over a time period
func (r *PgDeploymentMetricsRepository) AggregateMetrics(ctx context.Context, deploymentID uuid.UUID, startDate, endDate time.Time) (*model.DeploymentMetrics, error) {
	query := `
		SELECT 
			$1 as deployment_id,
			$4 as metric_date,
			SUM(total_trades) as total_trades,
			SUM(win_trades) as win_trades,
			SUM(loss_trades) as loss_trades,
			SUM(total_pnl) as total_pnl,
			MAX(max_drawdown) as max_drawdown,
			AVG(sharpe_ratio) as sharpe_ratio,
			AVG(sortino_ratio) as sortino_ratio,
			CASE WHEN SUM(total_trades) > 0 THEN SUM(win_trades)::float / SUM(total_trades) ELSE NULL END as win_rate,
			CASE WHEN SUM(win_trades) > 0 THEN SUM(CASE WHEN win_trades > 0 THEN total_pnl ELSE 0 END) / SUM(win_trades) ELSE NULL END as avg_win,
			CASE WHEN SUM(loss_trades) > 0 THEN SUM(CASE WHEN loss_trades > 0 THEN total_pnl ELSE 0 END) / SUM(loss_trades) ELSE NULL END as avg_loss
		FROM deployment_metrics
		WHERE deployment_id = $1 AND metric_date >= $2 AND metric_date <= $3
		GROUP BY deployment_id
	`

	aggregateDate := time.Now() // This is just a placeholder for the aggregate result
	var metrics model.DeploymentMetrics
	err := r.db.GetContext(ctx, &metrics, query, deploymentID, startDate, endDate, aggregateDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to aggregate metrics: %w", err)
	}

	return &metrics, nil
}
