package repository

import (
	"context"

	"algo/internal/modules/backtest/domain/model"

	"github.com/google/uuid"
)

// RunGroupRepository defines the interface for run group data access
type RunGroupRepository interface {
	// Create saves a new run group in the database
	Create(ctx context.Context, runGroup *model.RunGroup) error

	// GetByID retrieves a run group by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.RunGroup, error)

	// GetByStrategyID retrieves all run groups for a specific strategy
	GetByStrategyID(ctx context.Context, strategyID uuid.UUID) ([]*model.RunGroup, error)

	// Delete removes a run group from the database
	Delete(ctx context.Context, id uuid.UUID) error
}
