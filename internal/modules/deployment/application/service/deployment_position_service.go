package service

import (
	"context"
	"fmt"
	"time"

	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/domain/model"
	strategyModel "algo/internal/modules/strategy/domain/model"

	"github.com/google/uuid"
)

// CreatePosition creates a new position for a deployment
func (s *deploymentService) CreatePosition(ctx context.Context, req *dto.CreatePositionRequest, userID uuid.UUID) (*dto.PositionResponse, error) {
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

	// Create the position
	// TODO: Add direction to CreatePositionRequest DTO to allow specifying Long/Short
	position, err := model.NewDeploymentPosition(
		req.DeploymentID,
		req.InstrumentID,
		req.Quantity,
		req.EntryPrice,
		req.EntryTime,
		strategyModel.Long, // Default to Long for manual position creation
		req.Metadata,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create position: %w", err)
	}

	// Save the position
	if err := s.positionRepo.Create(ctx, position); err != nil {
		return nil, fmt.Errorf("failed to save position: %w", err)
	}

	// Create a log entry for the position
	metadata := map[string]interface{}{
		"position_id":   position.ID,
		"instrument_id": position.InstrumentID,
		"quantity":      position.Quantity,
		"entry_price":   position.EntryPrice,
		"action":        "POSITION_OPENED",
	}

	log, err := model.NewDeploymentLog(
		req.DeploymentID,
		model.StatusStarted,
		metadata,
	)
	if err != nil {
		// Just log the error, don't fail the entire operation
		fmt.Printf("Error creating position log: %v\n", err)
	} else {
		if err := s.logRepo.Create(ctx, log); err != nil {
			fmt.Printf("Error saving position log: %v\n", err)
		}
	}

	return dto.ToPositionResponse(position), nil
}

// GetPositionByID retrieves a position by its ID
func (s *deploymentService) GetPositionByID(ctx context.Context, id int64, userID uuid.UUID) (*dto.PositionResponse, error) {
	// Get the position
	position, err := s.positionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	if position == nil {
		return nil, fmt.Errorf("position not found")
	}

	// Verify the deployment belongs to the user
	deployment, err := s.deploymentRepo.GetByID(ctx, position.DeploymentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	if deployment == nil {
		return nil, ErrDeploymentNotFound
	}

	if deployment.UserID != userID {
		return nil, ErrUnauthorized
	}

	return dto.ToPositionResponse(position), nil
}

// GetPositionsByDeployment retrieves all positions for a deployment
func (s *deploymentService) GetPositionsByDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) ([]*dto.PositionResponse, error) {
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

	// Get the positions
	positions, err := s.positionRepo.GetByDeploymentID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	return dto.ToPositionResponses(positions), nil
}

// GetOpenPositionsByDeployment retrieves all open positions for a deployment
func (s *deploymentService) GetOpenPositionsByDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) ([]*dto.PositionResponse, error) {
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

	// Get the open positions
	positions, err := s.positionRepo.GetOpenPositionsByDeploymentID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get open positions: %w", err)
	}

	return dto.ToPositionResponses(positions), nil
}

// ClosePosition closes a position with exit details
func (s *deploymentService) ClosePosition(ctx context.Context, id int64, req *dto.ClosePositionRequest, userID uuid.UUID) (*dto.PositionResponse, error) {
	// Get the position
	position, err := s.positionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	if position == nil {
		return nil, fmt.Errorf("position not found")
	}

	// Verify the deployment belongs to the user
	deployment, err := s.deploymentRepo.GetByID(ctx, position.DeploymentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	if deployment == nil {
		return nil, ErrDeploymentNotFound
	}

	if deployment.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Handle nil request
	var exitPrice float64
	var exitTime time.Time
	var exitReason string

	if req != nil {
		exitPrice = req.ExitPrice
		exitTime = req.ExitTime
		exitReason = req.ExitReason
	} else {
		// Default values if request is nil
		exitTime = time.Now()
		// Use entry price as exit price (fallback)
		exitPrice = position.EntryPrice
		exitReason = "MANUAL"
	}

	// Set default exit reason if not provided
	if exitReason == "" {
		exitReason = "MANUAL"
	}

	// Close the position
	if err := s.positionRepo.ClosePosition(ctx, id, exitPrice, exitTime, exitReason); err != nil {
		return nil, fmt.Errorf("failed to close position: %w", err)
	}

	// Get the updated position
	updatedPosition, err := s.positionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated position: %w", err)
	}

	// Create a log entry for the position
	metadata := map[string]interface{}{
		"position_id":   updatedPosition.ID,
		"instrument_id": updatedPosition.InstrumentID,
		"quantity":      updatedPosition.Quantity,
		"entry_price":   updatedPosition.EntryPrice,
		"exit_price":    exitPrice,
		"pnl":           updatedPosition.PNL,
		"action":        "POSITION_CLOSED",
	}

	log, err := model.NewDeploymentLog(
		position.DeploymentID,
		model.StatusCompleted,
		metadata,
	)
	if err != nil {
		// Just log the error, don't fail the entire operation
		fmt.Printf("Error creating position close log: %v\n", err)
	} else {
		if err := s.logRepo.Create(ctx, log); err != nil {
			fmt.Printf("Error saving position close log: %v\n", err)
		}
	}

	return dto.ToPositionResponse(updatedPosition), nil
}

// CancelPosition cancels a position
func (s *deploymentService) CancelPosition(ctx context.Context, id int64, userID uuid.UUID) error {
	// Get the position
	position, err := s.positionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get position: %w", err)
	}

	if position == nil {
		return fmt.Errorf("position not found")
	}

	// Verify the deployment belongs to the user
	deployment, err := s.deploymentRepo.GetByID(ctx, position.DeploymentID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	if deployment == nil {
		return ErrDeploymentNotFound
	}

	if deployment.UserID != userID {
		return ErrUnauthorized
	}

	// Cancel the position
	if err := s.positionRepo.CancelPosition(ctx, id); err != nil {
		return fmt.Errorf("failed to cancel position: %w", err)
	}

	// Create a log entry for the position
	metadata := map[string]interface{}{
		"position_id":   position.ID,
		"instrument_id": position.InstrumentID,
		"quantity":      position.Quantity,
		"entry_price":   position.EntryPrice,
		"action":        "POSITION_CANCELED",
	}

	log, err := model.NewDeploymentLog(
		position.DeploymentID,
		model.StatusStopped,
		metadata,
	)
	if err != nil {
		// Just log the error, don't fail the entire operation
		fmt.Printf("Error creating position cancel log: %v\n", err)
	} else {
		if err := s.logRepo.Create(ctx, log); err != nil {
			fmt.Printf("Error saving position cancel log: %v\n", err)
		}
	}

	return nil
}

// ClosePositionWithReason closes a position with a specific exit reason
func (s *deploymentService) ClosePositionWithReason(ctx context.Context, id int64, exitPrice float64, exitTime time.Time, exitReason string, userID uuid.UUID) (*dto.PositionResponse, error) {
	req := &dto.ClosePositionRequest{
		ExitPrice:  exitPrice,
		ExitTime:   exitTime,
		ExitReason: exitReason,
	}
	return s.ClosePosition(ctx, id, req, userID)
}
