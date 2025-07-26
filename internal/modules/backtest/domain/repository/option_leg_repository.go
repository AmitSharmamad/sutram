package repository

import (
	"context"

	"algo/internal/modules/backtest/domain/model"

	"github.com/google/uuid"
)

// OptionLegRepository defines the interface for option leg data access
type OptionLegRepository interface {
	// Create saves a new option leg in the database
	Create(ctx context.Context, optionLeg *model.OptionLeg) error

	// GetByID retrieves an option leg by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.OptionLeg, error)

	// GetByStrategyID retrieves all option legs for a specific strategy
	GetByStrategyID(ctx context.Context, strategyID uuid.UUID) ([]*model.OptionLeg, error)

	// GetByStrategyIDAndSequence retrieves an option leg by strategy ID and leg sequence
	GetByStrategyIDAndSequence(ctx context.Context, strategyID uuid.UUID, legSequence int) (*model.OptionLeg, error)

	// Update updates an existing option leg in the database
	Update(ctx context.Context, optionLeg *model.OptionLeg) error

	// Delete removes an option leg from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByStrategyID removes all option legs for a specific strategy
	DeleteByStrategyID(ctx context.Context, strategyID uuid.UUID) error

	// GetByIndex retrieves all option legs for a specific index/underlying
	GetByIndex(ctx context.Context, selectIndex string) ([]*model.OptionLeg, error)

	// GetByExpiryDate retrieves all option legs with a specific expiry date
	GetByExpiryDate(ctx context.Context, expiryDate string) ([]*model.OptionLeg, error)

	// GetByStatus retrieves all option legs with a specific status
	GetByStatus(ctx context.Context, status model.LegStatus) ([]*model.OptionLeg, error)

	// List retrieves all option legs with pagination
	List(ctx context.Context, offset, limit int) ([]*model.OptionLeg, int64, error)

	// UpdateStatus updates the status of an option leg
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.LegStatus) error
}
