package repository

import (
	"context"

	"algo/internal/modules/strategy/domain/model"

	"github.com/google/uuid"
)

// StrategyRepository defines the interface for strategy data access
type StrategyRepository interface {
	// Create saves a new strategy in the database
	Create(ctx context.Context, strategy *model.Strategy) error

	// GetByID retrieves a strategy by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Strategy, error)

	// GetByUserID retrieves all strategies for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Strategy, error)

	// Update updates an existing strategy in the database
	Update(ctx context.Context, strategy *model.Strategy) error

	// Delete removes a strategy from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// ActivateStrategy sets the strategy as active
	ActivateStrategy(ctx context.Context, id uuid.UUID) error

	// DeactivateStrategy sets the strategy as inactive
	DeactivateStrategy(ctx context.Context, id uuid.UUID) error

	// GetActiveStrategies retrieves all active strategies
	GetActiveStrategies(ctx context.Context) ([]*model.Strategy, error)

	// GetActiveStrategiesByUserID retrieves all active strategies for a specific user
	GetActiveStrategiesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Strategy, error)

	// SearchStrategies searches strategies with filter criteria
	SearchStrategies(ctx context.Context, filters map[string]interface{}) ([]*model.Strategy, error)
}
