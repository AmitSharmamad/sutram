package service

import (
	"context"
	"errors"
	"fmt"

	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/domain/model"
	"algo/internal/modules/deployment/domain/repository"
	strategyRepository "algo/internal/modules/strategy/domain/repository"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrDeploymentNotFound      = errors.New("deployment not found")
	ErrStrategyNotFound        = errors.New("strategy not found")
	ErrInvalidRequest          = errors.New("invalid request")
	ErrUnauthorized            = errors.New("unauthorized access to deployment")
	ErrDeploymentCreateFailed  = errors.New("failed to create deployment")
	ErrDeploymentUpdateFailed  = errors.New("failed to update deployment")
	ErrDeploymentDeleteFailed  = errors.New("failed to delete deployment")
	ErrDeploymentAlreadyExists = errors.New("deployment already exists for this strategy and instrument")
)

// DeploymentService defines the service for deployment operations
type DeploymentService interface {
	// CreateDeployment creates a new deployment
	CreateDeployment(ctx context.Context, req *dto.CreateDeploymentRequest, userID uuid.UUID) (*dto.DeploymentResponse, error)

	// GetDeploymentByID retrieves a deployment by its ID
	GetDeploymentByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.DeploymentResponse, error)

	// GetUserDeployments retrieves all deployments for a user
	GetUserDeployments(ctx context.Context, userID uuid.UUID) ([]*dto.DeploymentResponse, error)

	// GetDeploymentsByStrategyID retrieves all deployments for a strategy
	GetDeploymentsByStrategyID(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]*dto.DeploymentResponse, error)

	// UpdateDeployment updates an existing deployment
	UpdateDeployment(ctx context.Context, req *dto.UpdateDeploymentRequest, userID uuid.UUID) (*dto.DeploymentResponse, error)

	// DeleteDeployment deletes a deployment
	DeleteDeployment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// ActivateDeployment activates a deployment
	ActivateDeployment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// DeactivateDeployment deactivates a deployment
	DeactivateDeployment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// SearchDeployments searches deployments with filter criteria
	SearchDeployments(ctx context.Context, filter *dto.DeploymentFilter, userID uuid.UUID) ([]*dto.DeploymentResponse, error)

	// StopAllDeployments deactivates all deployments for a user
	StopAllDeployments(ctx context.Context, userID uuid.UUID) error

	// StartAllDeployments activates all paused deployments for a user
	StartAllDeployments(ctx context.Context, userID uuid.UUID) error

	// DeployAllInstruments creates deployments for all instruments of a strategy
	DeployAllInstruments(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]*dto.DeploymentResponse, error)
}

// deploymentService implements DeploymentService
type deploymentService struct {
	deploymentRepo repository.DeploymentRepository
	strategyRepo   strategyRepository.StrategyRepository
}

// NewDeploymentService creates a new deployment service
func NewDeploymentService(
	deploymentRepo repository.DeploymentRepository,
	strategyRepo strategyRepository.StrategyRepository,
) DeploymentService {
	return &deploymentService{
		deploymentRepo: deploymentRepo,
		strategyRepo:   strategyRepo,
	}
}

// CreateDeployment creates a new deployment
func (s *deploymentService) CreateDeployment(ctx context.Context, req *dto.CreateDeploymentRequest, userID uuid.UUID) (*dto.DeploymentResponse, error) {
	// Verify the strategy exists and belongs to the user
	strategy, err := s.strategyRepo.GetByID(ctx, req.StrategyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStrategyNotFound, err)
	}

	if strategy.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Check if deployment already exists
	existingDeployment, err := s.deploymentRepo.GetByUserStrategyInstrument(ctx, userID, req.StrategyID, req.InstrumentID)
	if err == nil && existingDeployment != nil {
		return nil, ErrDeploymentAlreadyExists
	}

	// Create new deployment
	deployment := model.NewDeployment(
		userID,
		req.StrategyID,
		req.InstrumentID,
		req.Name,
		req.Description,
	)

	// Save to repository
	if err := s.deploymentRepo.Create(ctx, deployment); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeploymentCreateFailed, err)
	}

	// Return response
	return &dto.DeploymentResponse{
		ID:           deployment.ID,
		UserID:       deployment.UserID,
		StrategyID:   deployment.StrategyID,
		InstrumentID: deployment.InstrumentID,
		Name:         deployment.Name,
		Description:  deployment.Description,
		IsActive:     deployment.IsActive,
		CreatedAt:    deployment.CreatedAt,
		UpdatedAt:    deployment.UpdatedAt,
	}, nil
}

// GetDeploymentByID retrieves a deployment by its ID
func (s *deploymentService) GetDeploymentByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.DeploymentResponse, error) {
	deployment, err := s.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	// Check authorization
	if deployment.UserID != userID {
		return nil, ErrUnauthorized
	}

	return &dto.DeploymentResponse{
		ID:           deployment.ID,
		UserID:       deployment.UserID,
		StrategyID:   deployment.StrategyID,
		InstrumentID: deployment.InstrumentID,
		Name:         deployment.Name,
		Description:  deployment.Description,
		IsActive:     deployment.IsActive,
		CreatedAt:    deployment.CreatedAt,
		UpdatedAt:    deployment.UpdatedAt,
	}, nil
}

// GetUserDeployments retrieves all deployments for a user
func (s *deploymentService) GetUserDeployments(ctx context.Context, userID uuid.UUID) ([]*dto.DeploymentResponse, error) {
	deployments, err := s.deploymentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user deployments: %w", err)
	}

	response := make([]*dto.DeploymentResponse, len(deployments))
	for i, deployment := range deployments {
		response[i] = &dto.DeploymentResponse{
			ID:           deployment.ID,
			UserID:       deployment.UserID,
			StrategyID:   deployment.StrategyID,
			InstrumentID: deployment.InstrumentID,
			Name:         deployment.Name,
			Description:  deployment.Description,
			IsActive:     deployment.IsActive,
			CreatedAt:    deployment.CreatedAt,
			UpdatedAt:    deployment.UpdatedAt,
		}
	}

	return response, nil
}

// GetDeploymentsByStrategyID retrieves all deployments for a strategy
func (s *deploymentService) GetDeploymentsByStrategyID(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]*dto.DeploymentResponse, error) {
	// Verify the strategy exists and belongs to the user
	strategy, err := s.strategyRepo.GetByID(ctx, strategyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStrategyNotFound, err)
	}

	if strategy.UserID != userID {
		return nil, ErrUnauthorized
	}

	deployments, err := s.deploymentRepo.GetByStrategyID(ctx, strategyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get strategy deployments: %w", err)
	}

	response := make([]*dto.DeploymentResponse, len(deployments))
	for i, deployment := range deployments {
		response[i] = &dto.DeploymentResponse{
			ID:           deployment.ID,
			UserID:       deployment.UserID,
			StrategyID:   deployment.StrategyID,
			InstrumentID: deployment.InstrumentID,
			Name:         deployment.Name,
			Description:  deployment.Description,
			IsActive:     deployment.IsActive,
			CreatedAt:    deployment.CreatedAt,
			UpdatedAt:    deployment.UpdatedAt,
		}
	}

	return response, nil
}

// UpdateDeployment updates an existing deployment
func (s *deploymentService) UpdateDeployment(ctx context.Context, req *dto.UpdateDeploymentRequest, userID uuid.UUID) (*dto.DeploymentResponse, error) {
	// Retrieve existing deployment
	deployment, err := s.deploymentRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	// Check authorization
	if deployment.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Update fields
	deployment.Update(req.Name, req.Description)

	// Save to repository
	if err := s.deploymentRepo.Update(ctx, deployment); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeploymentUpdateFailed, err)
	}

	return &dto.DeploymentResponse{
		ID:           deployment.ID,
		UserID:       deployment.UserID,
		StrategyID:   deployment.StrategyID,
		InstrumentID: deployment.InstrumentID,
		Name:         deployment.Name,
		Description:  deployment.Description,
		IsActive:     deployment.IsActive,
		CreatedAt:    deployment.CreatedAt,
		UpdatedAt:    deployment.UpdatedAt,
	}, nil
}

// DeleteDeployment deletes a deployment
func (s *deploymentService) DeleteDeployment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Retrieve existing deployment
	deployment, err := s.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	// Check authorization
	if deployment.UserID != userID {
		return ErrUnauthorized
	}

	// Delete from repository
	if err := s.deploymentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrDeploymentDeleteFailed, err)
	}

	return nil
}

// ActivateDeployment activates a deployment
func (s *deploymentService) ActivateDeployment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Retrieve existing deployment
	deployment, err := s.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	// Check authorization
	if deployment.UserID != userID {
		return ErrUnauthorized
	}

	// Activate deployment
	if err := s.deploymentRepo.ActivateDeployment(ctx, id); err != nil {
		return fmt.Errorf("failed to activate deployment: %w", err)
	}

	return nil
}

// DeactivateDeployment deactivates a deployment
func (s *deploymentService) DeactivateDeployment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Retrieve existing deployment
	deployment, err := s.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeploymentNotFound, err)
	}

	// Check authorization
	if deployment.UserID != userID {
		return ErrUnauthorized
	}

	// Deactivate deployment
	if err := s.deploymentRepo.DeactivateDeployment(ctx, id); err != nil {
		return fmt.Errorf("failed to deactivate deployment: %w", err)
	}

	return nil
}

// SearchDeployments searches deployments with filter criteria
func (s *deploymentService) SearchDeployments(ctx context.Context, filter *dto.DeploymentFilter, userID uuid.UUID) ([]*dto.DeploymentResponse, error) {
	// Build filters map
	filters := make(map[string]interface{})

	// Always filter by user ID for security
	filters["user_id"] = userID

	if filter.StrategyID != nil {
		filters["strategy_id"] = *filter.StrategyID
	}

	if filter.InstrumentID != nil {
		filters["instrument_id"] = *filter.InstrumentID
	}

	if filter.IsActive != nil {
		filters["is_active"] = *filter.IsActive
	}

	if filter.Name != nil {
		filters["name"] = *filter.Name
	}

	// Perform search
	deployments, err := s.deploymentRepo.SearchDeployments(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search deployments: %w", err)
	}

	// Map to response
	response := make([]*dto.DeploymentResponse, len(deployments))
	for i, deployment := range deployments {
		response[i] = &dto.DeploymentResponse{
			ID:           deployment.ID,
			UserID:       deployment.UserID,
			StrategyID:   deployment.StrategyID,
			InstrumentID: deployment.InstrumentID,
			Name:         deployment.Name,
			Description:  deployment.Description,
			IsActive:     deployment.IsActive,
			CreatedAt:    deployment.CreatedAt,
			UpdatedAt:    deployment.UpdatedAt,
		}
	}

	return response, nil
}

// StopAllDeployments deactivates all deployments for a user
func (s *deploymentService) StopAllDeployments(ctx context.Context, userID uuid.UUID) error {
	// Get all active deployments for the user
	deployments, err := s.deploymentRepo.GetActiveDeploymentsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get active deployments: %w", err)
	}

	// If there are no active deployments, return early
	if len(deployments) == 0 {
		return nil
	}

	// Deactivate each deployment
	for _, deployment := range deployments {
		if err := s.deploymentRepo.DeactivateDeployment(ctx, deployment.ID); err != nil {
			return fmt.Errorf("failed to deactivate deployment %s: %w", deployment.ID, err)
		}
	}

	return nil
}

// StartAllDeployments activates all paused deployments for a user
func (s *deploymentService) StartAllDeployments(ctx context.Context, userID uuid.UUID) error {
	// Get all inactive deployments for the user
	deployments, err := s.deploymentRepo.GetInactiveDeploymentsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get inactive deployments: %w", err)
	}

	// If there are no inactive deployments, return early
	if len(deployments) == 0 {
		return nil
	}

	// Activate each deployment
	for _, deployment := range deployments {
		if err := s.deploymentRepo.ActivateDeployment(ctx, deployment.ID); err != nil {
			return fmt.Errorf("failed to activate deployment %s: %w", deployment.ID, err)
		}
	}

	return nil
}

// DeployAllInstruments creates deployments for all instruments of a strategy
func (s *deploymentService) DeployAllInstruments(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]*dto.DeploymentResponse, error) {
	// Verify the strategy exists and belongs to the user
	strategy, err := s.strategyRepo.GetByID(ctx, strategyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStrategyNotFound, err)
	}

	if strategy.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Get existing deployments for this strategy
	existingDeployments, err := s.deploymentRepo.GetByStrategyID(ctx, strategyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing deployments: %w", err)
	}

	// Create a map of existing deployments by instrument ID for quick lookup
	existingDeploymentsByInstrument := make(map[int]bool)
	for _, deployment := range existingDeployments {
		existingDeploymentsByInstrument[deployment.InstrumentID] = true
	}

	// Create deployments for each instrument that doesn't already have one
	var newDeployments []*dto.DeploymentResponse
	for _, instrumentID := range strategy.Instruments {
		// Skip if deployment already exists for this instrument
		if existingDeploymentsByInstrument[int(instrumentID)] {
			continue
		}

		// Create deployment request
		deploymentName := fmt.Sprintf("%s - Instrument %d", strategy.Name, instrumentID)
		deploymentDesc := fmt.Sprintf("Auto-deployed for strategy: %s", strategy.Name)
		req := &dto.CreateDeploymentRequest{
			StrategyID:   strategyID,
			InstrumentID: int(instrumentID),
			Name:         deploymentName,
			Description:  deploymentDesc,
		}

		// Create the deployment
		deployment, err := s.CreateDeployment(ctx, req, userID)
		if err != nil {
			// Log the error but continue with other instruments
			fmt.Printf("Error creating deployment for instrument %d: %v\n", instrumentID, err)
			continue
		}

		newDeployments = append(newDeployments, deployment)
	}

	return newDeployments, nil
}
