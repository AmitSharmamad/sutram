package repository

import (
	"context"
	"time"

	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// DeploymentPositionRepository defines the interface for deployment position data access
type DeploymentPositionRepository interface {
	// Create saves a new position in the database
	Create(ctx context.Context, position *model.DeploymentPosition) error

	// GetByID retrieves a position by its ID
	GetByID(ctx context.Context, id int64) (*model.DeploymentPosition, error)

	// GetByDeploymentID retrieves all positions for a specific deployment
	GetByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*model.DeploymentPosition, error)

	// GetOpenPositionsByDeploymentID retrieves all open positions for a specific deployment
	GetOpenPositionsByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*model.DeploymentPosition, error)

	// GetByDeploymentAndInstrument retrieves all positions for a specific deployment and instrument
	GetByDeploymentAndInstrument(ctx context.Context, deploymentID uuid.UUID, instrumentID int) ([]*model.DeploymentPosition, error)

	// Update updates an existing position in the database
	Update(ctx context.Context, position *model.DeploymentPosition) error

	// ClosePosition updates a position with exit details
	ClosePosition(ctx context.Context, id int64, exitPrice float64, exitTime time.Time) error

	// CancelPosition marks a position as canceled
	CancelPosition(ctx context.Context, id int64) error

	// CountPositions counts positions by deployment ID and optional status
	CountPositions(ctx context.Context, deploymentID uuid.UUID, status *model.PositionStatus) (int, error)
}
