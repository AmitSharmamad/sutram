package repository

import (
	"context"

	"algo/internal/modules/backtest/domain/model"

	"github.com/google/uuid"
)

// RunMetricsRepository defines the interface for run metrics data access
type RunMetricsRepository interface {
	// Create saves new run metrics in the database
	Create(ctx context.Context, metrics *model.RunMetrics) error

	// GetByRunID retrieves metrics for a specific run
	GetByRunID(ctx context.Context, runID uuid.UUID) (*model.RunMetrics, error)

	// Update updates existing run metrics in the database
	Update(ctx context.Context, metrics *model.RunMetrics) error

	// Delete removes run metrics from the database
	Delete(ctx context.Context, id uuid.UUID) error
}
