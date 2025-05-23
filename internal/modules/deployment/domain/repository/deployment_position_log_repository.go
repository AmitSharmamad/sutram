package repository

import (
	"context"
	"time"

	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// DeploymentPositionLogRepository defines the interface for deployment position log data access
type DeploymentPositionLogRepository interface {
	// Create saves a new position log in the database
	Create(ctx context.Context, positionLog *model.DeploymentPositionLog) error

	// GetByID retrieves a position log by its ID
	GetByID(ctx context.Context, id string) (*model.DeploymentPositionLog, error)

	// GetByPositionID retrieves all logs for a specific position
	GetByPositionID(ctx context.Context, positionID int64) ([]*model.DeploymentPositionLog, error)

	// GetByPositionIDAndTimeRange retrieves logs for a position within a specific time range
	GetByPositionIDAndTimeRange(ctx context.Context, positionID int64, startTime, endTime time.Time) ([]*model.DeploymentPositionLog, error)

	// GetByDeploymentID retrieves all logs for all positions in a deployment
	GetByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*model.DeploymentPositionLog, error)
}
