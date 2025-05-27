package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	backtestEngine "algo/internal/modules/backtest/application/service"
	backtestModel "algo/internal/modules/backtest/domain/model"
	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/domain/model"
	"algo/internal/modules/deployment/domain/repository"
	instrumentModel "algo/internal/modules/instrument/domain/model"
	instrumentRepository "algo/internal/modules/instrument/domain/repository"
	strategyDto "algo/internal/modules/strategy/application/dto"
	strategyModel "algo/internal/modules/strategy/domain/model"
	strategyRepository "algo/internal/modules/strategy/domain/repository"

	"github.com/google/uuid"
)

// DeploymentWithStrategy holds a deployment along with its associated strategy
type DeploymentWithStrategy struct {
	Deployment *model.Deployment
	Strategy   *strategyModel.Strategy
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

	// Collect unique strategy IDs from all deployments
	uniqueStrategyIDs := make(map[uuid.UUID]bool)
	for _, deployment := range deployments {
		uniqueStrategyIDs[deployment.StrategyID] = true
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

	// Group deployments by instrument ID and attach their strategies
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

		// Create DeploymentWithStrategy object
		deploymentWithStrategy := &DeploymentWithStrategy{
			Deployment: deployment,
			Strategy:   strategy,
		}

		instrumentDeployments[deployment.InstrumentID] = append(
			instrumentDeployments[deployment.InstrumentID],
			deploymentWithStrategy,
		)
	}

	// Store the deployments map for later use when handling events
	// This approach allows us to efficiently find deployments when we receive market data
	s.instrumentDeployments = instrumentDeployments

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

	// Prepare strategy data
	// This will be used in the future when we implement proper strategy execution
	_ = &strategyDto.StrategyResponse{
		ID:              strategy.ID,
		UserID:          strategy.UserID,
		Name:            strategy.Name,
		Description:     strategy.Description,
		EntryExpression: strategy.EntryExpression,
		ExitExpression:  strategy.ExitExpression,
		ChartTimeframe:  strategyModel.ChartTimeframe(strategy.ChartTimeframe),
		PositionType:    strategy.PositionType,
		TradeType:       strategyModel.TradeType(strategy.TradeType),
		Capital:         strategy.Capital,
		Instruments:     strategy.Instruments,
	}

	// Create a dummy run for future use
	// Currently not utilized but will be needed for strategy execution
	_ = &backtestModel.Run{
		ID:           uuid.New(),
		StrategyID:   strategy.ID,
		InstrumentID: deployment.InstrumentID,
		State:        backtestModel.Started,
	}

	// Create a single-element slice with the current instrument data
	// This will be used in the future for strategy execution
	_ = []*instrumentModel.Instrument{instrumentData}

	// Execute the strategy with a simulated backtest run
	// TODO: Implement proper strategy execution method for live trading
	// For now, use a simple placeholder
	positions := []*backtestModel.Position{} // Initialize as empty slice instead of nil
	// metrics will be nil initially, only create if we have actual metrics to report

	// Placeholder for actual strategy execution logic
	fmt.Printf("Running strategy for instrument %d at %s\n",
		instrumentData.ExchangeSecurityID, instrumentData.Timestamp)

	// Process any positions that were generated
	if len(positions) > 0 {
		// Create a deployment position for each position
		for _, position := range positions {
			// Only create a position if it's in "Entered" state
			if position.State == backtestModel.Entered {
				deploymentPosition, err := model.NewDeploymentPosition(
					deployment.ID,
					deployment.InstrumentID,
					position.PositionSize,
					position.EntryPriceSignal,
					position.EntryTimeSignal,
					map[string]interface{}{
						"strategy_id": strategy.ID.String(),
						"created_by":  "live_trading",
						"timestamp":   time.Now().Format(time.RFC3339),
					},
				)

				if err != nil {
					fmt.Printf("Error creating deployment position: %v\n", err)
					continue
				}

				if err := s.positionRepo.Create(ctx, deploymentPosition); err != nil {
					fmt.Printf("Error creating deployment position: %v\n", err)
					continue
				}

				// Log the position creation
				logEntry, err := model.NewDeploymentLog(
					deployment.ID,
					model.StatusCreated,
					map[string]interface{}{
						"position_id":   deploymentPosition.ID,
						"entry_price":   deploymentPosition.EntryPrice,
						"entry_time":    deploymentPosition.EntryTime,
						"quantity":      deploymentPosition.Quantity,
						"instrument_id": deployment.InstrumentID,
					},
				)
				if err != nil {
					fmt.Printf("Error creating deployment log for position: %v\n", err)
				} else if err := s.logRepo.Create(ctx, logEntry); err != nil {
					fmt.Printf("Error saving deployment log for position: %v\n", err)
				}
			}
		}
	}

	// TODO: Update deployment metrics when actual strategy execution is implemented
	// For now, we skip metrics creation since this is placeholder code

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
