package service

import (
	"context"
	"fmt"

	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// CreateDeploymentLog creates a new deployment log entry
func (s *deploymentService) CreateDeploymentLog(ctx context.Context, req *dto.DeploymentLogRequest, userID uuid.UUID) (*dto.DeploymentLogResponse, error) {
	// Verify the deployment exists and belongs to the user
	deployment, err := s.deploymentRepo.GetByID(ctx, req.DeploymentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	if deployment == nil {
		return nil, ErrDeploymentNotFound
	}

	if deployment.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Create the log entry
	log, err := model.NewDeploymentLog(req.DeploymentID, req.Status, req.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment log: %w", err)
	}

	// Save the log
	if err := s.logRepo.Create(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to save deployment log: %w", err)
	}

	return dto.ToDeploymentLogResponse(log), nil
}

// GetDeploymentLogs retrieves all log entries for a deployment
func (s *deploymentService) GetDeploymentLogs(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) ([]*dto.DeploymentLogResponse, error) {
	// Verify the deployment exists and belongs to the user
	deployment, err := s.deploymentRepo.GetByID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	if deployment == nil {
		return nil, ErrDeploymentNotFound
	}

	if deployment.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Get the logs
	logs, err := s.logRepo.GetByDeploymentID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment logs: %w", err)
	}

	return dto.ToDeploymentLogResponses(logs), nil
}

// GetLatestDeploymentLog retrieves the latest log entry for a deployment
func (s *deploymentService) GetLatestDeploymentLog(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) (*dto.DeploymentLogResponse, error) {
	// Verify the deployment exists and belongs to the user
	deployment, err := s.deploymentRepo.GetByID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	if deployment == nil {
		return nil, ErrDeploymentNotFound
	}

	if deployment.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Get the latest log
	log, err := s.logRepo.GetLatestByDeploymentID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest deployment log: %w", err)
	}

	if log == nil {
		return nil, nil
	}

	return dto.ToDeploymentLogResponse(log), nil
}
