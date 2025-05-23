package service

import (
	"context"
	"time"

	"algo/internal/modules/deployment/domain/model"
	"algo/internal/modules/deployment/domain/repository"
)

// DeploymentPositionLogService handles business logic for deployment position logs
type DeploymentPositionLogService struct {
	repository repository.DeploymentPositionLogRepository
}

// NewDeploymentPositionLogService creates a new deployment position log service
func NewDeploymentPositionLogService(repo repository.DeploymentPositionLogRepository) *DeploymentPositionLogService {
	return &DeploymentPositionLogService{
		repository: repo,
	}
}

// CreateLog creates a new position log
func (s *DeploymentPositionLogService) CreateLog(ctx context.Context, positionLog *model.DeploymentPositionLog) error {
	return s.repository.Create(ctx, positionLog)
}

// GetLogByID retrieves a position log by its ID
func (s *DeploymentPositionLogService) GetLogByID(ctx context.Context, id string) (*model.DeploymentPositionLog, error) {
	return s.repository.GetByID(ctx, id)
}

// GetLogsByPositionID retrieves all logs for a specific position
func (s *DeploymentPositionLogService) GetLogsByPositionID(ctx context.Context, positionID int64) ([]*model.DeploymentPositionLog, error) {
	return s.repository.GetByPositionID(ctx, positionID)
}

// GetLogsByPositionIDAndTimeRange retrieves logs for a position within a specific time range
func (s *DeploymentPositionLogService) GetLogsByPositionIDAndTimeRange(ctx context.Context, positionID int64, startTime, endTime time.Time) ([]*model.DeploymentPositionLog, error) {
	return s.repository.GetByPositionIDAndTimeRange(ctx, positionID, startTime, endTime)
}
