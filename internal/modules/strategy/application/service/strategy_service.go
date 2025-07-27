package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	deploymentRepository "algo/internal/modules/deployment/domain/repository"
	"algo/internal/modules/strategy/application/dto"
	"algo/internal/modules/strategy/domain/model"
	"algo/internal/modules/strategy/domain/repository"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrStrategyNotFound     = errors.New("strategy not found")
	ErrInvalidRequest       = errors.New("invalid request")
	ErrUnauthorized         = errors.New("unauthorized access to strategy")
	ErrInvalidExpression    = errors.New("invalid expression")
	ErrStrategyCreateFailed = errors.New("failed to create strategy")
	ErrStrategyUpdateFailed = errors.New("failed to update strategy")
	ErrStrategyDeleteFailed = errors.New("failed to delete strategy")
	ErrStrategyReadonly     = errors.New("strategy cannot be modified as deployments exist")
)

// StrategyService defines the service for strategy operations
type StrategyService interface {
	// CreateStrategy creates a new strategy
	CreateStrategy(ctx context.Context, req *dto.CreateStrategyRequest, userID uuid.UUID) (*dto.StrategyResponse, error)

	// GetStrategyByID retrieves a strategy by its ID
	GetStrategyByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.StrategyResponse, error)

	// GetUserStrategies retrieves all strategies for a user
	GetUserStrategies(ctx context.Context, userID uuid.UUID) ([]*dto.StrategyResponse, error)

	// GetDefaultStrategies retrieves all strategies for the system user
	GetDefaultStrategies(ctx context.Context) ([]*dto.StrategyResponse, error)

	// UpdateStrategy updates an existing strategy
	UpdateStrategy(ctx context.Context, req *dto.UpdateStrategyRequest, userID uuid.UUID) (*dto.StrategyResponse, error)

	// DeleteStrategy deletes a strategy
	DeleteStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// ActivateStrategy activates a strategy
	ActivateStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// DeactivateStrategy deactivates a strategy
	DeactivateStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// SearchStrategies searches strategies with filter criteria
	SearchStrategies(ctx context.Context, filter *dto.StrategyFilter, userID uuid.UUID) ([]*dto.StrategyResponse, error)

	// GetRepository returns the strategy repository
	GetRepository() repository.StrategyRepository

	// DuplicateStrategy creates a copy of an existing strategy
	DuplicateStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.StrategyResponse, error)

	// CloneDefaultStrategy clones a default strategy for a user
	CloneDefaultStrategy(ctx context.Context, defaultStrategyID uuid.UUID, userID uuid.UUID) (*dto.StrategyResponse, error)
}

// strategyService implements the StrategyService interface
type strategyService struct {
	strategyRepo     repository.StrategyRepository
	deploymentRepo   deploymentRepository.DeploymentRepository
	optionLegService OptionLegService
}

// NewStrategyService creates a new strategy service
func NewStrategyService(strategyRepo repository.StrategyRepository, deploymentRepo deploymentRepository.DeploymentRepository, optionLegService OptionLegService) StrategyService {
	return &strategyService{
		strategyRepo:     strategyRepo,
		deploymentRepo:   deploymentRepo,
		optionLegService: optionLegService,
	}
}

// CreateStrategy creates a new strategy
func (s *strategyService) CreateStrategy(ctx context.Context, req *dto.CreateStrategyRequest, userID uuid.UUID) (*dto.StrategyResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidRequest
	}

	// Validate entry expression
	if err := validateExpression(req.EntryExpression); err != nil {
		return nil, ErrInvalidExpression
	}

	// Validate exit expression if provided
	if req.ExitExpression != "" {
		if err := validateExpression(req.ExitExpression); err != nil {
			return nil, ErrInvalidExpression
		}
	}

	// Convert DTO to domain model
	strategy := req.ToDomain(userID)

	// Save strategy
	err := s.strategyRepo.Create(ctx, strategy)
	if err != nil {
		return nil, ErrStrategyCreateFailed
	}

	// Create option legs if segment is F&O and option legs are provided
	response := dto.FromDomain(strategy)
	// Always include option legs field (empty array for non-F&O segments)
	response.OptionLegs = []*dto.OptionLegResponse{}
	if req.Segment == "F&O" && len(req.OptionLegs) > 0 {
		var optionLegResponses []*dto.OptionLegResponse
		for i, legReq := range req.OptionLegs {
			// Set strategy ID and leg sequence if not provided
			legReq.StrategyID = strategy.ID
			if legReq.LegSequence == 0 {
				legReq.LegSequence = i + 1
			}

			optionLeg, err := s.optionLegService.CreateOptionLeg(ctx, &legReq)
			if err != nil {
				// If option leg creation fails, we should consider whether to delete the strategy
				// For now, we'll return an error but keep the strategy
				return nil, fmt.Errorf("failed to create option leg: %w", err)
			}
			optionLegResponses = append(optionLegResponses, optionLeg)
		}
		response.OptionLegs = optionLegResponses
	}

	// Convert domain model to response DTO
	return response, nil
}

// GetStrategyByID retrieves a strategy by its ID
func (s *strategyService) GetStrategyByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.StrategyResponse, error) {
	strategy, err := s.strategyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrStrategyNotFound
	}

	// Check if user has access to the strategy
	if strategy.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Check if strategy has deployments to determine readonly status
	hasDeployments, err := s.hasDeployments(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check deployments: %w", err)
	}

	response := dto.FromDomainWithReadonly(strategy, hasDeployments)

	// Always include option legs field (empty array for non-F&O segments)
	response.OptionLegs = []*dto.OptionLegResponse{}
	if strategy.Segment == "F&O" {
		optionLegs, err := s.optionLegService.GetOptionLegsByStrategyID(ctx, id)
		if err == nil {
			response.OptionLegs = optionLegs
		}
	}

	return response, nil
}

// GetUserStrategies retrieves all strategies for a user
func (s *strategyService) GetUserStrategies(ctx context.Context, userID uuid.UUID) ([]*dto.StrategyResponse, error) {
	strategies, err := s.strategyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert domain models to response DTOs with readonly status
	var strategyResponses []*dto.StrategyResponse
	for _, strategy := range strategies {
		hasDeployments, err := s.hasDeployments(ctx, strategy.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check deployments for strategy %s: %w", strategy.ID, err)
		}
		response := dto.FromDomainWithReadonly(strategy, hasDeployments)

		// Always include option legs field (empty array for non-F&O segments)
		response.OptionLegs = []*dto.OptionLegResponse{}
		if strategy.Segment == "F&O" {
			optionLegs, err := s.optionLegService.GetOptionLegsByStrategyID(ctx, strategy.ID)
			if err == nil {
				response.OptionLegs = optionLegs
			}
		}

		strategyResponses = append(strategyResponses, response)
	}

	return strategyResponses, nil
}

// GetDefaultStrategies retrieves all strategies for the system user
func (s *strategyService) GetDefaultStrategies(ctx context.Context) ([]*dto.StrategyResponse, error) {
	// System user ID from default.sql
	systemUserID, err := uuid.Parse("e65c67f1-2b43-4581-b33b-ec26295e741b")
	if err != nil {
		return nil, fmt.Errorf("invalid system user ID: %w", err)
	}

	// Get strategies for system user
	strategies, err := s.strategyRepo.GetByUserID(ctx, systemUserID)
	if err != nil {
		return nil, err
	}

	// Convert domain models to response DTOs with readonly status
	var strategyResponses []*dto.StrategyResponse
	for _, strategy := range strategies {
		hasDeployments, err := s.hasDeployments(ctx, strategy.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check deployments for strategy %s: %w", strategy.ID, err)
		}
		response := dto.FromDomainWithReadonly(strategy, hasDeployments)

		// Always include option legs field (empty array for non-F&O segments)
		response.OptionLegs = []*dto.OptionLegResponse{}
		if strategy.Segment == "F&O" {
			optionLegs, err := s.optionLegService.GetOptionLegsByStrategyID(ctx, strategy.ID)
			if err == nil {
				response.OptionLegs = optionLegs
			}
		}

		strategyResponses = append(strategyResponses, response)
	}

	return strategyResponses, nil
}

// UpdateStrategy updates an existing strategy
func (s *strategyService) UpdateStrategy(ctx context.Context, req *dto.UpdateStrategyRequest, userID uuid.UUID) (*dto.StrategyResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidRequest
	}

	// Retrieve the strategy
	existingStrategy, err := s.strategyRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, ErrStrategyNotFound
	}

	// Check if user has access to the strategy
	if existingStrategy.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Check if deployments exist for this strategy (only if deployment repo is available)
	if s.deploymentRepo != nil {
		deployments, err := s.deploymentRepo.GetByStrategyID(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check strategy deployments: %w", err)
		}

		// If any deployments exist, don't allow strategy modification
		if len(deployments) > 0 {
			return nil, ErrStrategyReadonly
		}
	} // Validate entry expression
	if err := validateExpression(req.EntryExpression); err != nil {
		return nil, ErrInvalidExpression
	}

	// Validate exit expression if provided
	if req.ExitExpression != "" {
		if err := validateExpression(req.ExitExpression); err != nil {
			return nil, ErrInvalidExpression
		}
	}

	// Convert request DTO to domain model
	updatedStrategy := req.ToDomain()

	// Preserve the original user ID and creation timestamp
	updatedStrategy.UserID = existingStrategy.UserID
	updatedStrategy.CreatedAt = existingStrategy.CreatedAt
	updatedStrategy.UpdatedAt = time.Now()
	updatedStrategy.IsActive = existingStrategy.IsActive

	// Save updated strategy
	err = s.strategyRepo.Update(ctx, updatedStrategy)
	if err != nil {
		return nil, ErrStrategyUpdateFailed
	}

	// Handle option legs update if segment is F&O
	response := dto.FromDomainWithReadonly(updatedStrategy, false)
	// Always include option legs field (empty array for non-F&O segments)
	response.OptionLegs = []*dto.OptionLegResponse{}
	if req.Segment == "F&O" && len(req.OptionLegs) > 0 {
		// First, delete existing option legs for this strategy
		err := s.optionLegService.DeleteOptionLegsByStrategyID(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing option legs: %w", err)
		}

		// Create new option legs
		var optionLegResponses []*dto.OptionLegResponse
		for i, legReq := range req.OptionLegs {
			// Set strategy ID and leg sequence if not provided
			legReq.StrategyID = req.ID
			if legReq.LegSequence == 0 {
				legReq.LegSequence = i + 1
			}

			optionLeg, err := s.optionLegService.CreateOptionLeg(ctx, &legReq)
			if err != nil {
				return nil, fmt.Errorf("failed to create option leg: %w", err)
			}
			optionLegResponses = append(optionLegResponses, optionLeg)
		}
		response.OptionLegs = optionLegResponses
	} else if req.Segment == "F&O" {
		// If F&O segment but no option legs provided, get existing ones
		existingLegs, err := s.optionLegService.GetOptionLegsByStrategyID(ctx, req.ID)
		if err == nil {
			response.OptionLegs = existingLegs
		}
	}

	// Note: Since we already checked that no deployments exist, readonly should be false
	return response, nil
}

// DeleteStrategy deletes a strategy
func (s *strategyService) DeleteStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Retrieve the strategy
	strategy, err := s.strategyRepo.GetByID(ctx, id)
	if err != nil {
		return ErrStrategyNotFound
	}

	// Check if user has access to the strategy
	if strategy.UserID != userID {
		return ErrUnauthorized
	}

	// Delete the strategy
	err = s.strategyRepo.Delete(ctx, id)
	if err != nil {
		return ErrStrategyDeleteFailed
	}

	return nil
}

// ActivateStrategy activates a strategy
func (s *strategyService) ActivateStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Retrieve the strategy
	strategy, err := s.strategyRepo.GetByID(ctx, id)
	if err != nil {
		return ErrStrategyNotFound
	}

	// Check if user has access to the strategy
	if strategy.UserID != userID {
		return ErrUnauthorized
	}

	// Activate the strategy
	strategy.Activate()

	// Save updated strategy
	err = s.strategyRepo.Update(ctx, strategy)
	if err != nil {
		return ErrStrategyUpdateFailed
	}

	return nil
}

// DeactivateStrategy deactivates a strategy
func (s *strategyService) DeactivateStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Retrieve the strategy
	strategy, err := s.strategyRepo.GetByID(ctx, id)
	if err != nil {
		return ErrStrategyNotFound
	}

	// Check if user has access to the strategy
	if strategy.UserID != userID {
		return ErrUnauthorized
	}

	// Deactivate the strategy
	strategy.Deactivate()

	// Save updated strategy
	err = s.strategyRepo.Update(ctx, strategy)
	if err != nil {
		return ErrStrategyUpdateFailed
	}

	return nil
}

// SearchStrategies searches strategies with filter criteria
func (s *strategyService) SearchStrategies(ctx context.Context, filter *dto.StrategyFilter, userID uuid.UUID) ([]*dto.StrategyResponse, error) {
	// Convert filter DTO to map for repository
	filters := make(map[string]interface{})

	// Always filter by userID for security
	filters["user_id"] = userID

	if filter.Name != nil {
		filters["name"] = *filter.Name
	}

	if filter.PositionType != nil {
		filters["position_type"] = *filter.PositionType
	}

	if filter.Instrument != nil {
		filters["instrument"] = *filter.Instrument
	}

	if filter.ChartType != nil {
		filters["chart_type"] = *filter.ChartType
	}

	if filter.ChartTimeframe != nil {
		filters["chart_timeframe"] = *filter.ChartTimeframe
	}

	if filter.IsActive != nil {
		filters["is_active"] = *filter.IsActive
	}

	if filter.CreatedAfter != nil {
		filters["created_after"] = *filter.CreatedAfter
	}

	if filter.CreatedBefore != nil {
		filters["created_before"] = *filter.CreatedBefore
	}

	// Search strategies
	strategies, err := s.strategyRepo.SearchStrategies(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Convert domain models to response DTOs
	var strategyResponses []*dto.StrategyResponse
	for _, strategy := range strategies {
		strategyResponses = append(strategyResponses, dto.FromDomain(strategy))
	}

	return strategyResponses, nil
}

// DuplicateStrategy creates a copy of an existing strategy
func (s *strategyService) DuplicateStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.StrategyResponse, error) {
	// Retrieve the original strategy
	original, err := s.strategyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrStrategyNotFound
	}

	// Check if user has access to the strategy
	if original.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Create a new strategy with the same properties
	duplicate := model.NewStrategy(
		userID,
		original.Name+" (Copy)",
		original.Description,
		original.PositionType,
		original.Instruments,
		original.ChartType,
		original.ChartTimeframe,
		original.TradingDays,
		original.TradeType,
		original.Qty,
		original.EntryExpression,
		original.ExitExpression,
		original.TargetProfit,
		original.TargetProfitType,
		original.StopLoss,
		original.StopLossType,
		original.Capital,
		original.MaxAllocation,
	)

	// Copy additional optional fields
	duplicate.TrailingStopLoss = original.TrailingStopLoss
	duplicate.TrailingStopLossType = original.TrailingStopLossType
	duplicate.ExitType = original.ExitType
	duplicate.MaxTrades = original.MaxTrades
	duplicate.Brokerage = original.Brokerage
	duplicate.Slippage = original.Slippage
	duplicate.EntryStartTime = original.EntryStartTime
	duplicate.EntryEndTime = original.EntryEndTime
	duplicate.StartDate = original.StartDate
	duplicate.EndDate = original.EndDate

	// Save the duplicated strategy
	err = s.strategyRepo.Create(ctx, duplicate)
	if err != nil {
		return nil, ErrStrategyCreateFailed
	}

	return dto.FromDomain(duplicate), nil
}

// CloneDefaultStrategy clones a default strategy for a user
func (s *strategyService) CloneDefaultStrategy(ctx context.Context, defaultStrategyID uuid.UUID, userID uuid.UUID) (*dto.StrategyResponse, error) {
	// System user ID from default.sql
	systemUserID, err := uuid.Parse("e65c67f1-2b43-4581-b33b-ec26295e741b")
	if err != nil {
		return nil, fmt.Errorf("invalid system user ID: %w", err)
	}

	// Retrieve the default strategy
	defaultStrategy, err := s.strategyRepo.GetByID(ctx, defaultStrategyID)
	if err != nil {
		return nil, ErrStrategyNotFound
	}

	// Verify this is actually a default strategy (owned by system user)
	if defaultStrategy.UserID != systemUserID {
		return nil, errors.New("strategy is not a default strategy")
	}

	// Create a new strategy with the same properties but owned by the current user
	clonedStrategy := model.NewStrategy(
		userID,
		defaultStrategy.Name,
		defaultStrategy.Description,
		defaultStrategy.PositionType,
		defaultStrategy.Instruments,
		defaultStrategy.ChartType,
		defaultStrategy.ChartTimeframe,
		defaultStrategy.TradingDays,
		defaultStrategy.TradeType,
		defaultStrategy.Qty,
		defaultStrategy.EntryExpression,
		defaultStrategy.ExitExpression,
		defaultStrategy.TargetProfit,
		defaultStrategy.TargetProfitType,
		defaultStrategy.StopLoss,
		defaultStrategy.StopLossType,
		defaultStrategy.Capital,
		defaultStrategy.MaxAllocation,
	)

	// Copy additional optional fields
	clonedStrategy.TrailingStopLoss = defaultStrategy.TrailingStopLoss
	clonedStrategy.TrailingStopLossType = defaultStrategy.TrailingStopLossType
	clonedStrategy.ExitType = defaultStrategy.ExitType
	clonedStrategy.MaxTrades = defaultStrategy.MaxTrades
	clonedStrategy.Brokerage = defaultStrategy.Brokerage
	clonedStrategy.Slippage = defaultStrategy.Slippage
	clonedStrategy.EntryStartTime = defaultStrategy.EntryStartTime
	clonedStrategy.EntryEndTime = defaultStrategy.EntryEndTime
	clonedStrategy.StartDate = defaultStrategy.StartDate
	clonedStrategy.EndDate = defaultStrategy.EndDate

	// Save the cloned strategy
	err = s.strategyRepo.Create(ctx, clonedStrategy)
	if err != nil {
		return nil, ErrStrategyCreateFailed
	}

	return dto.FromDomain(clonedStrategy), nil
}

// validateExpression validates the trading expression syntax
// This is a placeholder for the actual implementation which would
// check the expression syntax, variable references, and operations
func validateExpression(expression string) error {
	// Placeholder for actual expression validation logic
	// In a real implementation, this would parse and validate the expression
	if expression == "" {
		return errors.New("expression cannot be empty")
	}

	return nil
}

// GetRepository returns the strategy repository
func (s *strategyService) GetRepository() repository.StrategyRepository {
	return s.strategyRepo
}

// hasDeployments checks if a strategy has any deployments
func (s *strategyService) hasDeployments(ctx context.Context, strategyID uuid.UUID) (bool, error) {
	// If deployment repository is not available, assume no deployments exist
	if s.deploymentRepo == nil {
		return false, nil
	}

	deployments, err := s.deploymentRepo.GetByStrategyID(ctx, strategyID)
	if err != nil {
		return false, fmt.Errorf("failed to check strategy deployments: %w", err)
	}

	return len(deployments) > 0, nil
}
