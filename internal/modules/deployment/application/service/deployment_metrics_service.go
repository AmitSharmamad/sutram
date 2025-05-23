package service

import (
	"context"
	"fmt"
	"time"

	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// CreateOrUpdateMetrics creates or updates metrics for a deployment
func (s *deploymentService) CreateOrUpdateMetrics(ctx context.Context, req *dto.CreateMetricsRequest, userID uuid.UUID) (*dto.MetricsResponse, error) {
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

	// Check if metrics already exist for this date
	existingMetrics, err := s.metricsRepo.GetByDeploymentIDAndDate(ctx, req.DeploymentID, req.MetricDate)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing metrics: %w", err)
	}

	// Create or update the metrics
	var metrics *model.DeploymentMetrics
	if existingMetrics != nil {
		// Update existing metrics
		existingMetrics.TotalTrades = req.TotalTrades
		existingMetrics.WinTrades = req.WinTrades
		existingMetrics.LossTrades = req.LossTrades
		existingMetrics.TotalPNL = req.TotalPNL
		existingMetrics.MaxDrawdown = req.MaxDrawdown
		existingMetrics.SharpeRatio = req.SharpeRatio
		existingMetrics.SortinoRatio = req.SortinoRatio

		// Update derived metrics
		if req.TotalTrades > 0 {
			winRate := float64(req.WinTrades) / float64(req.TotalTrades)
			existingMetrics.WinRate = &winRate
		}

		if req.WinTrades > 0 {
			avgWin := req.TotalPNL / float64(req.WinTrades)
			existingMetrics.AvgWin = &avgWin
		}

		if req.LossTrades > 0 {
			avgLoss := req.TotalPNL / float64(req.LossTrades)
			existingMetrics.AvgLoss = &avgLoss
		}

		existingMetrics.UpdatedAt = time.Now()

		// Update in database
		if err := s.metricsRepo.Update(ctx, existingMetrics); err != nil {
			return nil, fmt.Errorf("failed to update metrics: %w", err)
		}

		metrics = existingMetrics
	} else {
		// Create new metrics
		newMetrics, err := model.NewDeploymentMetrics(
			req.DeploymentID,
			req.MetricDate,
			req.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create metrics: %w", err)
		}

		// Set the metrics values
		newMetrics.TotalTrades = req.TotalTrades
		newMetrics.WinTrades = req.WinTrades
		newMetrics.LossTrades = req.LossTrades
		newMetrics.TotalPNL = req.TotalPNL
		newMetrics.MaxDrawdown = req.MaxDrawdown
		newMetrics.SharpeRatio = req.SharpeRatio
		newMetrics.SortinoRatio = req.SortinoRatio

		// Calculate derived metrics
		if req.TotalTrades > 0 {
			winRate := float64(req.WinTrades) / float64(req.TotalTrades)
			newMetrics.WinRate = &winRate
		}

		if req.WinTrades > 0 {
			avgWin := req.TotalPNL / float64(req.WinTrades)
			newMetrics.AvgWin = &avgWin
		}

		if req.LossTrades > 0 {
			avgLoss := req.TotalPNL / float64(req.LossTrades)
			newMetrics.AvgLoss = &avgLoss
		}

		// Save in database
		if err := s.metricsRepo.Create(ctx, newMetrics); err != nil {
			return nil, fmt.Errorf("failed to save metrics: %w", err)
		}

		metrics = newMetrics
	}

	// Create a log entry for the metrics update
	metadata := map[string]interface{}{
		"metric_date":  req.MetricDate,
		"total_trades": req.TotalTrades,
		"total_pnl":    req.TotalPNL,
		"action":       "METRICS_UPDATED",
	}

	log, err := model.NewDeploymentLog(
		req.DeploymentID,
		model.StatusCompleted,
		metadata,
	)
	if err != nil {
		// Just log the error, don't fail the entire operation
		fmt.Printf("Error creating metrics log: %v\n", err)
	} else {
		if err := s.logRepo.Create(ctx, log); err != nil {
			fmt.Printf("Error saving metrics log: %v\n", err)
		}
	}

	return dto.ToMetricsResponse(metrics), nil
}

// GetMetricsByDeployment retrieves all metrics for a deployment
func (s *deploymentService) GetMetricsByDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) ([]*dto.MetricsResponse, error) {
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

	// Get the metrics
	metrics, err := s.metricsRepo.GetByDeploymentID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return dto.ToMetricsResponses(metrics), nil
}

// GetMetricsByDeploymentAndDate retrieves metrics for a deployment and date
func (s *deploymentService) GetMetricsByDeploymentAndDate(ctx context.Context, deploymentID uuid.UUID, date time.Time, userID uuid.UUID) (*dto.MetricsResponse, error) {
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

	// Get the metrics
	metrics, err := s.metricsRepo.GetByDeploymentIDAndDate(ctx, deploymentID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	if metrics == nil {
		return nil, nil
	}

	return dto.ToMetricsResponse(metrics), nil
}

// GetMetricsTimeRange retrieves metrics for a deployment within a time range
func (s *deploymentService) GetMetricsTimeRange(ctx context.Context, filter *dto.MetricsFilter, userID uuid.UUID) ([]*dto.MetricsResponse, error) {
	// Verify the deployment exists and belongs to the user
	deployment, err := s.deploymentRepo.GetByID(ctx, filter.DeploymentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	if deployment == nil {
		return nil, ErrDeploymentNotFound
	}

	if deployment.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Set default dates if not provided
	startDate := time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	if filter.StartDate != nil {
		startDate = *filter.StartDate
	}

	endDate := time.Now() // Default to now
	if filter.EndDate != nil {
		endDate = *filter.EndDate
	}

	// Get the metrics
	metrics, err := s.metricsRepo.GetMetricsTimeRange(ctx, filter.DeploymentID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return dto.ToMetricsResponses(metrics), nil
}

// GetLatestMetrics retrieves the latest metrics for a deployment
func (s *deploymentService) GetLatestMetrics(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) (*dto.MetricsResponse, error) {
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

	// Get the latest metrics
	metrics, err := s.metricsRepo.GetLatestByDeploymentID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest metrics: %w", err)
	}

	if metrics == nil {
		return nil, nil
	}

	return dto.ToMetricsResponse(metrics), nil
}
