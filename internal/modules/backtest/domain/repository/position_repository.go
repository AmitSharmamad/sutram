package repository

import (
	"context"

	"algo/internal/modules/backtest/domain/model"

	"github.com/google/uuid"
)

// PositionRepository defines the interface for position data access
type PositionRepository interface {
	// Create saves a new position in the database
	Create(ctx context.Context, position *model.Position) error

	// BulkCreate saves multiple positions in the database in a single operation
	BulkCreate(ctx context.Context, positions []*model.Position) error

	// GetByID retrieves a position by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Position, error)

	// GetByRunID retrieves all positions for a specific run
	GetByRunID(ctx context.Context, runID uuid.UUID) ([]*model.Position, error)

	// Update updates an existing position in the database
	Update(ctx context.Context, position *model.Position) error

	// Delete removes a position from the database
	Delete(ctx context.Context, id uuid.UUID) error
}
