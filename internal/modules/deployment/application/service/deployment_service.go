package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	backtestEngine "algo/internal/modules/backtest/application/service"
	backtestModel "algo/internal/modules/backtest/domain/model"
	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/domain/model"
	"algo/internal/modules/deployment/domain/repository"
	instrumentModel "algo/internal/modules/instrument/domain/model"
	instrumentRepository "algo/internal/modules/instrument/domain/repository"
	parser "algo/internal/modules/parser/domain/model"
	strategyDto "algo/internal/modules/strategy/application/dto"
	strategyModel "algo/internal/modules/strategy/domain/model"
	strategyRepository "algo/internal/modules/strategy/domain/repository"

	"github.com/google/uuid"
)

// DeploymentWithStrategy holds a deployment along with its associated strategy and current position
type DeploymentWithStrategy struct {
	Deployment      *model.Deployment
	Strategy        *strategyModel.Strategy
	CurrentPosition *model.DeploymentPosition // Pre-fetched open position for efficient trading
}

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
	// Deployment core operations
	CreateDeployment(ctx context.Context, req *dto.CreateDeploymentRequest, userID uuid.UUID) (*dto.DeploymentResponse, error)
	GetDeploymentByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.DeploymentResponse, error)
	GetUserDeployments(ctx context.Context, userID uuid.UUID) ([]*dto.DeploymentResponse, error)
	GetDeploymentsByStrategyID(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]*dto.DeploymentResponse, error)
	UpdateDeployment(ctx context.Context, req *dto.UpdateDeploymentRequest, userID uuid.UUID) (*dto.DeploymentResponse, error)
	DeleteDeployment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	ActivateDeployment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	DeactivateDeployment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	SearchDeployments(ctx context.Context, filter *dto.DeploymentFilter, userID uuid.UUID) ([]*dto.DeploymentResponse, error)
	StopAllDeployments(ctx context.Context, userID uuid.UUID) error
	StartAllDeployments(ctx context.Context, userID uuid.UUID) error
	DeployAllInstruments(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]*dto.DeploymentResponse, error)

	// Deployment log operations
	CreateDeploymentLog(ctx context.Context, req *dto.DeploymentLogRequest, userID uuid.UUID) (*dto.DeploymentLogResponse, error)
	GetDeploymentLogs(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) ([]*dto.DeploymentLogResponse, error)
	GetLatestDeploymentLog(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) (*dto.DeploymentLogResponse, error)

	// Deployment position operations
	CreatePosition(ctx context.Context, req *dto.CreatePositionRequest, userID uuid.UUID) (*dto.PositionResponse, error)
	GetPositionByID(ctx context.Context, id int64, userID uuid.UUID) (*dto.PositionResponse, error)
	GetPositionsByDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) ([]*dto.PositionResponse, error)
	GetOpenPositionsByDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) ([]*dto.PositionResponse, error)
	ClosePosition(ctx context.Context, id int64, req *dto.ClosePositionRequest, userID uuid.UUID) (*dto.PositionResponse, error)
	CancelPosition(ctx context.Context, id int64, userID uuid.UUID) error

	// Deployment metrics operations
	CreateOrUpdateMetrics(ctx context.Context, req *dto.CreateMetricsRequest, userID uuid.UUID) (*dto.MetricsResponse, error)
	GetMetricsByDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) ([]*dto.MetricsResponse, error)
	GetMetricsByDeploymentAndDate(ctx context.Context, deploymentID uuid.UUID, date time.Time, userID uuid.UUID) (*dto.MetricsResponse, error)
	GetMetricsTimeRange(ctx context.Context, filter *dto.MetricsFilter, userID uuid.UUID) ([]*dto.MetricsResponse, error)
	GetLatestMetrics(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) (*dto.MetricsResponse, error)

	// Live trading operations
	Initialize(ctx context.Context) error
	Run(ctx context.Context, deploymentWithStrategy *DeploymentWithStrategy, instrumentData *instrumentModel.Instrument) error
	ProcessInstrumentData(ctx context.Context, instrumentData *instrumentModel.Instrument) error
	GetActiveInstrumentIDs() []int
}

// deploymentService implements DeploymentService
type deploymentService struct {
	deploymentRepo        repository.DeploymentRepository
	strategyRepo          strategyRepository.StrategyRepository
	logRepo               repository.DeploymentLogRepository
	positionRepo          repository.DeploymentPositionRepository
	metricsRepo           repository.DeploymentMetricsRepository
	instrumentRepo        instrumentRepository.InstrumentRepository
	engine                *backtestEngine.BacktestEngine
	instrumentDeployments map[int][]*DeploymentWithStrategy
}

// NewDeploymentService creates a new deployment service
func NewDeploymentService(
	deploymentRepo repository.DeploymentRepository,
	strategyRepo strategyRepository.StrategyRepository,
	logRepo repository.DeploymentLogRepository,
	positionRepo repository.DeploymentPositionRepository,
	metricsRepo repository.DeploymentMetricsRepository,
	instrumentRepo instrumentRepository.InstrumentRepository,
	engine *backtestEngine.BacktestEngine,
) DeploymentService {
	return &deploymentService{
		deploymentRepo:        deploymentRepo,
		strategyRepo:          strategyRepo,
		logRepo:               logRepo,
		positionRepo:          positionRepo,
		metricsRepo:           metricsRepo,
		instrumentRepo:        instrumentRepo,
		engine:                engine,
		instrumentDeployments: make(map[int][]*DeploymentWithStrategy),
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

	// Create initial deployment metrics
	initialMetrics, err := model.NewDeploymentMetrics(
		deployment.ID,
		time.Now(),
		map[string]interface{}{
			"initial_setup": true,
			"status":        "created",
		},
	)
	if err != nil {
		// Log the error but don't fail the deployment creation
		fmt.Printf("Error creating initial deployment metrics: %v\n", err)
	} else {
		// Save the initial metrics
		if err := s.metricsRepo.Create(ctx, initialMetrics); err != nil {
			// Log the error but don't fail the deployment creation
			fmt.Printf("Error saving initial deployment metrics: %v\n", err)
		}
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

// Initialize initializes the deployment service for live trading
func (s *deploymentService) Initialize(ctx context.Context) error {
	// Get all active deployments for all users
	deployments, err := s.deploymentRepo.GetActiveDeployments(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active deployments: %w", err)
	}

	if len(deployments) == 0 {
		fmt.Printf("No active deployments found\n")
		s.instrumentDeployments = make(map[int][]*DeploymentWithStrategy)
		return nil
	}

	// Collect unique strategy IDs from all deployments
	uniqueStrategyIDs := make(map[uuid.UUID]bool)
	deploymentIDs := make([]uuid.UUID, len(deployments))

	for i, deployment := range deployments {
		uniqueStrategyIDs[deployment.StrategyID] = true
		deploymentIDs[i] = deployment.ID
	}

	// Convert map keys to slice for the batch call
	strategyIDs := make([]uuid.UUID, 0, len(uniqueStrategyIDs))
	for strategyID := range uniqueStrategyIDs {
		strategyIDs = append(strategyIDs, strategyID)
	}

	// Fetch all strategies at once using batch call
	strategies, err := s.strategyRepo.GetByIDs(ctx, strategyIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch strategies: %w", err)
	}

	// Create a map of strategy ID to strategy for quick lookup
	strategyMap := make(map[uuid.UUID]*strategyModel.Strategy)
	for _, strategy := range strategies {
		strategyMap[strategy.ID] = strategy
	}

	// Fetch all open positions for all deployments using bulk query
	fmt.Printf("Fetching open positions for %d deployments\n", len(deploymentIDs))
	positionMap, err := s.positionRepo.GetOpenPositionsByDeploymentIDs(ctx, deploymentIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch open positions: %w", err)
	}
	fmt.Printf("Found %d open positions across all deployments\n", len(positionMap))

	// Group deployments by instrument ID and attach their strategies and positions
	instrumentDeployments := make(map[int][]*DeploymentWithStrategy)
	for _, deployment := range deployments {
		// Get the strategy from the pre-fetched map
		strategy, exists := strategyMap[deployment.StrategyID]
		if !exists {
			// Log error but continue with other deployments
			fmt.Printf("Strategy %s not found for deployment %s\n",
				deployment.StrategyID, deployment.ID)
			continue
		}

		// Get the current open position for this deployment (if any)
		currentPosition := positionMap[deployment.ID]
		if currentPosition != nil {
			fmt.Printf("Found open position %s for deployment %s\n",
				currentPosition.ID, deployment.ID)
		}

		// Create DeploymentWithStrategy object with pre-fetched position
		deploymentWithStrategy := &DeploymentWithStrategy{
			Deployment:      deployment,
			Strategy:        strategy,
			CurrentPosition: currentPosition,
		}

		instrumentDeployments[deployment.InstrumentID] = append(
			instrumentDeployments[deployment.InstrumentID],
			deploymentWithStrategy,
		)
	}

	// Store the deployments map for later use when handling events
	// This approach allows us to efficiently find deployments when we receive market data
	s.instrumentDeployments = instrumentDeployments

	fmt.Printf("Initialized %d instruments with deployments\n", len(instrumentDeployments))
	return nil
}

// Run executes a deployment's strategy against live market data
func (s *deploymentService) Run(ctx context.Context, deploymentWithStrategy *DeploymentWithStrategy, instrumentData *instrumentModel.Instrument) error {
	deployment := deploymentWithStrategy.Deployment
	strategy := deploymentWithStrategy.Strategy

	// Log the execution
	logEntry, err := model.NewDeploymentLog(
		deployment.ID,
		model.StatusStarted,
		map[string]any{
			"instrument_data_time": instrumentData.Timestamp,
			"instrument_id":        instrumentData.ExchangeSecurityID,
		},
	)
	if err != nil {
		// Log but continue processing
		fmt.Printf("Error creating deployment log: %v\n", err)
	} else {
		if err := s.logRepo.Create(ctx, logEntry); err != nil {
			// Log but continue processing
			fmt.Printf("Error saving deployment log: %v\n", err)
		}
	}

	// Use pre-fetched position from DeploymentWithStrategy
	currentPosition := deploymentWithStrategy.CurrentPosition

	// Prepare strategy data
	strategyResponse := &strategyDto.StrategyResponse{
		ID:               strategy.ID,
		UserID:           strategy.UserID,
		Name:             strategy.Name,
		Description:      strategy.Description,
		EntryExpression:  strategy.EntryExpression,
		ExitExpression:   strategy.ExitExpression,
		ChartTimeframe:   strategyModel.ChartTimeframe(strategy.ChartTimeframe),
		PositionType:     strategy.PositionType,
		TradeType:        strategyModel.TradeType(strategy.TradeType),
		Capital:          strategy.Capital,
		Instruments:      strategy.Instruments,
		Qty:              strategy.Qty,
		MaxAllocation:    strategy.MaxAllocation,
		StopLoss:         strategy.StopLoss,
		StopLossType:     strategy.StopLossType,
		TargetProfit:     strategy.TargetProfit,
		TargetProfitType: strategy.TargetProfitType,
		EntryStartTime:   strategy.EntryStartTime,
		EntryEndTime:     strategy.EntryEndTime,
	}

	// Fetch historical data needed for strategy execution
	historicalData, err := s.instrumentRepo.FindByExchangeSecurityIDAndTimeRange(
		ctx,
		instrumentData.ExchangeSecurityID,
		instrumentData.Timestamp.Add(-24*time.Hour), // Look back 24 hours for indicators
		instrumentData.Timestamp,
		string(strategy.ChartTimeframe),
	)
	if err != nil {
		fmt.Printf("Error getting historical data for instrument %d: %v\n",
			instrumentData.ExchangeSecurityID, err)
		// Use current data point only as fallback
		historicalData = []*instrumentModel.Instrument{instrumentData}
	}

	// Add current instrument data to the end of historical data
	historicalData = append(historicalData, instrumentData)

	// Handle existing position first - check exit conditions
	if currentPosition != nil {
		fmt.Printf("Found existing open position %s for deployment %s\n",
			currentPosition.ID, deployment.ID)

		// Evaluate exit condition using current market data
		exitSignal, err := s.evaluateExitCondition(strategyResponse, historicalData, currentPosition)
		if err != nil {
			fmt.Printf("Error evaluating exit condition: %v\n", err)
		} else if exitSignal {
			// Close the position
			exitPrice := instrumentData.Close
			exitTime := instrumentData.Timestamp

			// Update position with exit details
			currentPosition.ClosePosition(exitPrice, exitTime)

			// Persist the position update
			if err := s.positionRepo.Update(ctx, currentPosition); err != nil {
				fmt.Printf("Error updating closed position: %v\n", err)
			} else {
				fmt.Printf("Closed position %s at price %.2f\n",
					currentPosition.ID, exitPrice)

				// Log the position closure
				logEntry, err := model.NewDeploymentLog(
					deployment.ID,
					model.StatusCompleted,
					map[string]any{
						"position_id": currentPosition.ID,
						"exit_price":  exitPrice,
						"exit_time":   exitTime,
						"pnl":         currentPosition.PNL,
					},
				)
				if err == nil {
					s.logRepo.Create(ctx, logEntry)
				}
			}

			// Clear current position since it's now closed
			currentPosition = nil
			// Update the DeploymentWithStrategy's position state
			deploymentWithStrategy.CurrentPosition = nil
		}
	}

	// If no current position exists, check entry conditions
	if currentPosition == nil {
		fmt.Printf("No open position for deployment %s, checking entry conditions\n", deployment.ID)

		// Evaluate entry condition
		dataContext := s.toDataContext(historicalData)
		entrySignal := s.evaluateEntryCondition(strategyResponse.EntryExpression, dataContext)
		if entrySignal {
			// Create new position
			entryPrice := instrumentData.Close
			entryTime := instrumentData.Timestamp
			quantity := s.calculatePositionSize(strategyResponse, entryPrice)

			newPosition, err := model.NewDeploymentPosition(
				deployment.ID,
				deployment.InstrumentID,
				quantity,
				entryPrice,
				entryTime,
				map[string]any{
					"strategy_id": strategy.ID.String(),
					"created_by":  "live_trading",
					"timestamp":   time.Now().Format(time.RFC3339),
				},
			)

			if err != nil {
				fmt.Printf("Error creating new deployment position: %v\n", err)
			} else {
				// Persist the new position
				if err := s.positionRepo.Create(ctx, newPosition); err != nil {
					fmt.Printf("Error saving new position: %v\n", err)
				} else {
					fmt.Printf("Created new position %s at price %.2f with quantity %.2f\n",
						newPosition.ID, entryPrice, quantity)

					// Log the position creation
					logEntry, err := model.NewDeploymentLog(
						deployment.ID,
						model.StatusCreated,
						map[string]any{
							"position_id":   newPosition.ID,
							"entry_price":   entryPrice,
							"entry_time":    entryTime,
							"quantity":      quantity,
							"instrument_id": deployment.InstrumentID,
						},
					)
					if err == nil {
						s.logRepo.Create(ctx, logEntry)
					}

					// Set as current position for metrics calculation
					currentPosition = newPosition
					// Update the DeploymentWithStrategy's position state
					deploymentWithStrategy.CurrentPosition = newPosition
				}
			}
		}
	}

	// Calculate and update deployment metrics
	metrics := s.calculateDeploymentMetrics(ctx, deployment.ID, currentPosition)
	if metrics != nil {
		deploymentMetrics, err := model.NewDeploymentMetrics(
			deployment.ID,
			time.Now(),
			map[string]any{
				"total_return_pct": metrics.TotalReturnPct,
				"total_trades":     metrics.TotalTrades,
				"win_rate":         metrics.WinRate,
				"sharpe_ratio":     metrics.SharpeRatio,
				"max_drawdown":     metrics.MaxDrawdown,
				"total_commission": metrics.TotalCommission,
				"profit_factor":    metrics.ProfitFactor,
				"avg_win_size":     metrics.AvgWinSize,
				"avg_loss_size":    metrics.AvgLossSize,
			},
		)

		if err != nil {
			fmt.Printf("Error creating deployment metrics: %v\n", err)
		} else {
			if err := s.metricsRepo.Create(ctx, deploymentMetrics); err != nil {
				fmt.Printf("Error saving deployment metrics: %v\n", err)
			}
		}
	}

	fmt.Printf("Strategy execution completed for deployment %s on instrument %d\n",
		deployment.ID, instrumentData.ExchangeSecurityID)

	return nil
}

// ProcessInstrumentData processes incoming instrument data and runs the appropriate deployments
func (s *deploymentService) ProcessInstrumentData(ctx context.Context, instrumentData *instrumentModel.Instrument) error {
	// Store the data in the database
	if err := s.instrumentRepo.Create(ctx, instrumentData); err != nil {
		return fmt.Errorf("failed to store instrument data: %w", err)
	}

	// Find deployments for this instrument
	instrumentID := int(instrumentData.ExchangeSecurityID)
	deployments, exists := s.instrumentDeployments[instrumentID]
	if !exists || len(deployments) == 0 {
		// No deployments for this instrument, nothing to do
		return nil
	}

	// Run each deployment for this instrument
	for _, deploymentWithStrategy := range deployments {
		if err := s.Run(ctx, deploymentWithStrategy, instrumentData); err != nil {
			// Log the error but continue with other deployments
			fmt.Printf("Error running deployment %s for instrument %d: %v\n",
				deploymentWithStrategy.Deployment.ID, instrumentID, err)
			continue
		}
	}

	return nil
}

// GetActiveInstrumentIDs returns the unique instrument IDs from all active deployments
func (s *deploymentService) GetActiveInstrumentIDs() []int {
	instrumentIDs := make([]int, 0, len(s.instrumentDeployments))
	for instrumentID := range s.instrumentDeployments {
		instrumentIDs = append(instrumentIDs, instrumentID)
	}
	return instrumentIDs
}

// toDataContext converts instruments to data context for expression evaluation
func (s *deploymentService) toDataContext(instruments []*instrumentModel.Instrument) *instrumentModel.StockData {
	stockData := make(instrumentModel.StockData, len(instruments))
	for i, instrument := range instruments {
		stockData[i] = *instrument
	}
	return &stockData
}

// evaluateEntryCondition evaluates entry expression using parser
func (s *deploymentService) evaluateEntryCondition(expression string, data *instrumentModel.StockData) bool {
	return parser.Execute(expression, parser.NewEvaluationContext(data))
}

// calculatePositionSize calculates position size based on strategy settings
func (s *deploymentService) calculatePositionSize(strategy *strategyDto.StrategyResponse, price float64) float64 {
	if strategy.Qty > 0 {
		return float64(strategy.Qty)
	}

	// Default to 1 unit if no quantity specified
	return 1.0
}

// evaluateExitCondition evaluates exit expression using parser for current position
func (s *deploymentService) evaluateExitCondition(strategy *strategyDto.StrategyResponse, historicalData []*instrumentModel.Instrument, currentPosition *model.DeploymentPosition) (bool, error) {
	if strategy.ExitExpression == "" {
		return false, nil
	}

	// Convert instruments to data context for expression evaluation
	dataContext := s.toDataContext(historicalData)

	// Execute the exit expression
	return parser.Execute(strategy.ExitExpression, parser.NewEvaluationContext(dataContext)), nil
}

// calculateDeploymentMetrics calculates metrics for a deployment based on its positions
func (s *deploymentService) calculateDeploymentMetrics(ctx context.Context, deploymentID uuid.UUID, currentPosition *model.DeploymentPosition) *backtestModel.RunMetrics {
	// Get all positions for this deployment to calculate metrics
	positions, err := s.positionRepo.GetByDeploymentID(ctx, deploymentID)
	if err != nil {
		fmt.Printf("Error fetching positions for metrics calculation: %v\n", err)
		return nil
	}

	totalPnL := 0.0
	winningTrades := 0
	losingTrades := 0
	totalTrades := 0
	totalWinAmount := 0.0
	totalLossAmount := 0.0

	for _, pos := range positions {
		if pos.Status == model.PositionClosed && pos.PNL != nil {
			totalTrades++
			pnl := *pos.PNL
			totalPnL += pnl

			if pnl > 0 {
				winningTrades++
				totalWinAmount += pnl
			} else if pnl < 0 {
				losingTrades++
				totalLossAmount += pnl
			}
		}
	}

	// Create basic metrics
	metrics := &backtestModel.RunMetrics{
		TotalTrades: totalTrades,
	}

	if totalTrades > 0 {
		winRate := float64(winningTrades) / float64(totalTrades) * 100
		metrics.WinRate = winRate

		// Calculate total return percentage (simplified)
		baseCapital := 100000.0 // Default base capital for percentage calculation
		if totalPnL != 0 {
			totalReturnPct := (totalPnL / baseCapital) * 100
			metrics.TotalReturnPct = totalReturnPct
		}

		// Calculate profit factor
		profitFactor := 1.0
		if totalLossAmount != 0 {
			profitFactor = totalWinAmount / math.Abs(totalLossAmount)
		}
		metrics.ProfitFactor = &profitFactor

		// Calculate average win and loss sizes
		avgWinSize := 0.0
		if winningTrades > 0 {
			avgWinSize = totalWinAmount / float64(winningTrades)
		}
		metrics.AvgWinSize = &avgWinSize

		avgLossSize := 0.0
		if losingTrades > 0 {
			avgLossSize = totalLossAmount / float64(losingTrades)
		}
		metrics.AvgLossSize = &avgLossSize

		// Basic metrics (simplified)
		metrics.MaxDrawdown = 0.0
		metrics.SharpeRatio = 0.0
		metrics.TotalCommission = 0.0
	}

	return metrics
}
