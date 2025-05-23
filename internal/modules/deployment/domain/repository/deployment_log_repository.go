package repository

import (
	"context"

	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// DeploymentLogRepository defines the interface for deployment log data access
type DeploymentLogRepository interface {
	// Create saves a new deployment log entry in the database
	Create(ctx context.Context, log *model.DeploymentLog) error

	// GetByID retrieves a deployment log entry by its ID
	GetByID(ctx context.Context, id int64) (*model.DeploymentLog, error)

	// GetByDeploymentID retrieves all log entries for a specific deployment
	GetByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*model.DeploymentLog, error)

	// GetByDeploymentIDAndStatus retrieves all log entries for a specific deployment and status
	GetByDeploymentIDAndStatus(ctx context.Context, deploymentID uuid.UUID, status model.DeploymentStatus) ([]*model.DeploymentLog, error)

	// GetLatestByDeploymentID retrieves the latest log entry for a specific deployment
	GetLatestByDeploymentID(ctx context.Context, deploymentID uuid.UUID) (*model.DeploymentLog, error)

	// ListLogs retrieves logs with pagination
	ListLogs(ctx context.Context, offset, limit int) ([]*model.DeploymentLog, error)

	// CountLogs counts total logs by deployment ID
	CountLogs(ctx context.Context, deploymentID uuid.UUID) (int, error)
}
