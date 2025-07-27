package repository

import (
	"context"

	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// DeploymentRepository defines the interface for deployment data access
type DeploymentRepository interface {
	// Create saves a new deployment in the database
	Create(ctx context.Context, deployment *model.Deployment) error

	// GetByID retrieves a deployment by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Deployment, error)

	// GetByUserID retrieves all deployments for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Deployment, error)

	// GetByStrategyID retrieves all deployments for a specific strategy
	GetByStrategyID(ctx context.Context, strategyID uuid.UUID) ([]*model.Deployment, error)

	// GetByInstrumentID retrieves all deployments for a specific instrument
	GetByInstrumentID(ctx context.Context, instrumentID int) ([]*model.Deployment, error)

	// GetByUserAndStrategy retrieves all deployments for a specific user and strategy
	GetByUserAndStrategy(ctx context.Context, userID, strategyID uuid.UUID) ([]*model.Deployment, error)

	// GetByUserStrategyInstrument retrieves a deployment by user, strategy, and instrument
	GetByUserStrategyInstrument(ctx context.Context, userID, strategyID uuid.UUID, instrumentID int) (*model.Deployment, error)

	// GetActiveByUserStrategyInstrument retrieves an active deployment by user, strategy, and instrument
	GetActiveByUserStrategyInstrument(ctx context.Context, userID, strategyID uuid.UUID, instrumentID int) (*model.Deployment, error)

	// Update updates an existing deployment in the database
	Update(ctx context.Context, deployment *model.Deployment) error

	// Delete removes a deployment from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// ActivateDeployment sets the deployment as active
	ActivateDeployment(ctx context.Context, id uuid.UUID) error

	// DeactivateDeployment sets the deployment as inactive
	DeactivateDeployment(ctx context.Context, id uuid.UUID) error

	// GetActiveDeployments retrieves all active deployments
	GetActiveDeployments(ctx context.Context) ([]*model.Deployment, error)

	// GetActiveDeploymentsByUserID retrieves all active deployments for a specific user
	GetActiveDeploymentsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Deployment, error)

	// GetInactiveDeploymentsByUserID retrieves all inactive deployments for a specific user
	GetInactiveDeploymentsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Deployment, error)

	// SearchDeployments searches deployments with filter criteria
	SearchDeployments(ctx context.Context, filters map[string]interface{}) ([]*model.Deployment, error)
}
