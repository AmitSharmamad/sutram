package repository

import (
	"context"

	"algo/internal/modules/backtest/domain/model"

	"github.com/google/uuid"
)

// RunRepository defines the interface for run data access
type RunRepository interface {
	// Create saves a new run in the database
	Create(ctx context.Context, run *model.Run) error

	// BulkCreate saves multiple runs in the database in a single operation
	BulkCreate(ctx context.Context, runs []*model.Run) error

	// GetByID retrieves a run by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Run, error)

	// GetByRunGroupID retrieves all runs for a specific run group
	GetByRunGroupID(ctx context.Context, runGroupID uuid.UUID) ([]*model.Run, error)

	// GetByStrategyID retrieves all runs for a specific strategy
	GetByStrategyID(ctx context.Context, strategyID uuid.UUID) ([]*model.Run, error)

	// Update updates an existing run in the database
	Update(ctx context.Context, run *model.Run) error

	// Delete removes a run from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateState updates the state of a run
	UpdateState(ctx context.Context, id uuid.UUID, state model.RunState) error

	// SearchRuns searches runs with filter criteria
	SearchRuns(ctx context.Context, filters map[string]interface{}) ([]*model.Run, error)
}
