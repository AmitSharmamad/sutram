package repository

import (
	"context"
	"time"

	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// DeploymentMetricsRepository defines the interface for deployment metrics data access
type DeploymentMetricsRepository interface {
	// Create saves new metrics in the database
	Create(ctx context.Context, metrics *model.DeploymentMetrics) error

	// GetByID retrieves metrics by ID
	GetByID(ctx context.Context, id int64) (*model.DeploymentMetrics, error)

	// GetByDeploymentID retrieves all metrics for a specific deployment
	GetByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*model.DeploymentMetrics, error)

	// GetByDeploymentIDAndDate retrieves metrics for a specific deployment and date
	GetByDeploymentIDAndDate(ctx context.Context, deploymentID uuid.UUID, date time.Time) (*model.DeploymentMetrics, error)

	// GetLatestByDeploymentID retrieves the latest metrics for a specific deployment
	GetLatestByDeploymentID(ctx context.Context, deploymentID uuid.UUID) (*model.DeploymentMetrics, error)

	// Update updates existing metrics in the database
	Update(ctx context.Context, metrics *model.DeploymentMetrics) error

	// GetMetricsTimeRange retrieves metrics for a deployment within a date range
	GetMetricsTimeRange(ctx context.Context, deploymentID uuid.UUID, startDate, endDate time.Time) ([]*model.DeploymentMetrics, error)

	// AggregateMetrics calculates aggregate metrics for a deployment over a time period
	AggregateMetrics(ctx context.Context, deploymentID uuid.UUID, startDate, endDate time.Time) (*model.DeploymentMetrics, error)
}
