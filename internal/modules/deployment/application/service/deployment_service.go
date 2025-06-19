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
	instrumentDto "algo/internal/modules/instrument/application/dto"
	instrumentModel "algo/internal/modules/instrument/domain/model"
	instrumentRepository "algo/internal/modules/instrument/domain/repository"
	parser "algo/internal/modules/parser/domain/model"
	strategyDto "algo/internal/modules/strategy/application/dto"
	strategyModel "algo/internal/modules/strategy/domain/model"
	strategyRepository "algo/internal/modules/strategy/domain/repository"
	marketHours "algo/internal/util/time"

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
	GetUserDeploymentsWithFilter(ctx context.Context, userID uuid.UUID, filter *dto.DeploymentFilter) ([]*dto.DeploymentResponse, error)
	GetDeploymentsByStrategyID(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]*dto.DeploymentResponse, error)
	GetDeploymentStatus(ctx context.Context, strategyID uuid.UUID, instrumentIDs []int, userID uuid.UUID) (*dto.DeploymentStatusResponse, error)
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
	instrumentMasterRepo  instrumentRepository.InstrumentMasterRepository
	engine                *backtestEngine.BacktestEngine
	instrumentDeployments map[int][]*DeploymentWithStrategy
	eventPublisher        DeploymentEventPublisher
}

// NewDeploymentService creates a new deployment service
func NewDeploymentService(
	deploymentRepo repository.DeploymentRepository,
	strategyRepo strategyRepository.StrategyRepository,
	logRepo repository.DeploymentLogRepository,
	positionRepo repository.DeploymentPositionRepository,
	metricsRepo repository.DeploymentMetricsRepository,
	instrumentRepo instrumentRepository.InstrumentRepository,
	instrumentMasterRepo instrumentRepository.InstrumentMasterRepository,
	engine *backtestEngine.BacktestEngine,
	eventPublisher DeploymentEventPublisher,
) DeploymentService {
	return &deploymentService{
		deploymentRepo:        deploymentRepo,
		strategyRepo:          strategyRepo,
		logRepo:               logRepo,
		positionRepo:          positionRepo,
		metricsRepo:           metricsRepo,
		instrumentRepo:        instrumentRepo,
		instrumentMasterRepo:  instrumentMasterRepo,
		engine:                engine,
		instrumentDeployments: make(map[int][]*DeploymentWithStrategy),
		eventPublisher:        eventPublisher,
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
	response := &dto.DeploymentResponse{
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

	// Get latest metrics for this deployment (should include the initial metrics we just created)
	latestMetrics, err := s.metricsRepo.GetLatestByDeploymentID(ctx, deployment.ID)
	if err == nil && latestMetrics != nil {
		response.LatestMetrics = dto.ToMetricsResponse(latestMetrics)
	}

	// Get instrument master data for this deployment
	instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, int64(deployment.InstrumentID))
	if err == nil && instrumentMaster != nil {
		response.InstrumentMaster = instrumentMasterToDTO(instrumentMaster)
	}

	return response, nil
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

	response := &dto.DeploymentResponse{
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

	// Get latest metrics for this deployment
	latestMetrics, err := s.metricsRepo.GetLatestByDeploymentID(ctx, deployment.ID)
	if err == nil && latestMetrics != nil {
		response.LatestMetrics = dto.ToMetricsResponse(latestMetrics)
	}

	// Get instrument master data for this deployment
	instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, int64(deployment.InstrumentID))
	if err == nil && instrumentMaster != nil {
		response.InstrumentMaster = instrumentMasterToDTO(instrumentMaster)
	}

	return response, nil
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

		// Get latest metrics for this deployment
		latestMetrics, err := s.metricsRepo.GetLatestByDeploymentID(ctx, deployment.ID)
		if err == nil && latestMetrics != nil {
			response[i].LatestMetrics = dto.ToMetricsResponse(latestMetrics)
		}

		// Get instrument master data for this deployment
		instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, int64(deployment.InstrumentID))
		if err == nil && instrumentMaster != nil {
			response[i].InstrumentMaster = instrumentMasterToDTO(instrumentMaster)
		}
	}

	return response, nil
}

// GetUserDeploymentsWithFilter retrieves deployments for a user with optional status filtering
func (s *deploymentService) GetUserDeploymentsWithFilter(ctx context.Context, userID uuid.UUID, filter *dto.DeploymentFilter) ([]*dto.DeploymentResponse, error) {
	deployments, err := s.deploymentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user deployments: %w", err)
	}

	// If no status filter is specified, return all deployments
	if filter == nil || filter.Status == nil {
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

			// Get instrument master data for this deployment
			instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, int64(deployment.InstrumentID))
			if err == nil && instrumentMaster != nil {
				response[i].InstrumentMaster = instrumentMasterToDTO(instrumentMaster)
			}
		}
		return response, nil
	}

	// Filter deployments based on position status
	var filteredDeployments []*model.Deployment

	for _, deployment := range deployments {
		// Get open positions for this deployment
		openPositions, err := s.positionRepo.GetOpenPositionsByDeploymentID(ctx, deployment.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get positions for deployment %s: %w", deployment.ID, err)
		}

		hasOpenPosition := len(openPositions) > 0

		// Apply status filter
		switch *filter.Status {
		case "waiting":
			// Include deployments with no open positions (waiting to enter)
			if !hasOpenPosition {
				filteredDeployments = append(filteredDeployments, deployment)
			}
		case "entered":
			// Include deployments with open positions (already entered)
			if hasOpenPosition {
				filteredDeployments = append(filteredDeployments, deployment)
			}
		}
	}

	// Convert to response DTOs
	response := make([]*dto.DeploymentResponse, len(filteredDeployments))
	for i, deployment := range filteredDeployments {
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

		// Get latest metrics for this deployment
		latestMetrics, err := s.metricsRepo.GetLatestByDeploymentID(ctx, deployment.ID)
		if err == nil && latestMetrics != nil {
			response[i].LatestMetrics = dto.ToMetricsResponse(latestMetrics)
		}

		// Get instrument master data for this deployment
		instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, int64(deployment.InstrumentID))
		if err == nil && instrumentMaster != nil {
			response[i].InstrumentMaster = instrumentMasterToDTO(instrumentMaster)
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

		// Get latest metrics for this deployment
		latestMetrics, err := s.metricsRepo.GetLatestByDeploymentID(ctx, deployment.ID)
		if err == nil && latestMetrics != nil {
			response[i].LatestMetrics = dto.ToMetricsResponse(latestMetrics)
		}

		// Get instrument master data for this deployment
		instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, int64(deployment.InstrumentID))
		if err == nil && instrumentMaster != nil {
			response[i].InstrumentMaster = instrumentMasterToDTO(instrumentMaster)
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

	response := &dto.DeploymentResponse{
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

	// Get latest metrics for this deployment
	latestMetrics, err := s.metricsRepo.GetLatestByDeploymentID(ctx, deployment.ID)
	if err == nil && latestMetrics != nil {
		response.LatestMetrics = dto.ToMetricsResponse(latestMetrics)
	}

	// Get instrument master data for this deployment
	instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, int64(deployment.InstrumentID))
	if err == nil && instrumentMaster != nil {
		response.InstrumentMaster = instrumentMasterToDTO(instrumentMaster)
	}

	return response, nil
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
	filters := make(map[string]any)

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

	// If status filter is specified, filter by position status
	var filteredDeployments []*model.Deployment
	if filter.Status != nil {
		for _, deployment := range deployments {
			// Get open positions for this deployment
			openPositions, err := s.positionRepo.GetOpenPositionsByDeploymentID(ctx, deployment.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get positions for deployment %s: %w", deployment.ID, err)
			}

			hasOpenPosition := len(openPositions) > 0

			// Apply status filter
			switch *filter.Status {
			case "waiting":
				// Include deployments with no open positions (waiting to enter)
				if !hasOpenPosition {
					filteredDeployments = append(filteredDeployments, deployment)
				}
			case "entered":
				// Include deployments with open positions (already entered)
				if hasOpenPosition {
					filteredDeployments = append(filteredDeployments, deployment)
				}
			}
		}
	} else {
		filteredDeployments = deployments
	}

	// Map to response
	response := make([]*dto.DeploymentResponse, len(filteredDeployments))
	for i, deployment := range filteredDeployments {
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

		// Fetch and attach latest metrics
		latestMetrics, err := s.metricsRepo.GetLatestByDeploymentID(ctx, deployment.ID)
		if err == nil && latestMetrics != nil {
			response[i].LatestMetrics = dto.ToMetricsResponse(latestMetrics)
		}

		// Get instrument master data for this deployment
		instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, int64(deployment.InstrumentID))
		if err == nil && instrumentMaster != nil {
			response[i].InstrumentMaster = instrumentMasterToDTO(instrumentMaster)
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

	// Check if instrument timestamp matches strategy timeframe
	if !s.isTimestampAlignedWithTimeframe(instrumentData.Timestamp, string(strategy.ChartTimeframe)) {
		fmt.Printf("Skipping execution - timestamp %v does not align with strategy timeframe %s\n",
			instrumentData.Timestamp.Format("15:04:05"), strategy.ChartTimeframe)
		return nil
	}

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
		MaxTrades:        0,   // Set default value
		Brokerage:        0,   // Set default value
		Slippage:         0.0, // Set default value
		StopLoss:         strategy.StopLoss,
		StopLossType:     strategy.StopLossType,
		TargetProfit:     strategy.TargetProfit,
		TargetProfitType: strategy.TargetProfitType,
		EntryStartTime:   strategy.EntryStartTime,
		EntryEndTime:     strategy.EntryEndTime,
	}

	// Set optional fields if they exist
	if strategy.MaxTrades != nil {
		strategyResponse.MaxTrades = *strategy.MaxTrades
	}
	if strategy.Brokerage != nil {
		strategyResponse.Brokerage = *strategy.Brokerage
	}
	if strategy.Slippage != nil {
		strategyResponse.Slippage = *strategy.Slippage
	}

	// Calculate lookback from strategy expressions (like backtesting engine)
	lookbackFromExpressions := s.getLookBackFromExpressions(strategyResponse)

	// Fetch historical data with proper lookback period using the same method as backtesting
	// Get lookback data (historical data before current timestamp)
	lookbackData, err := s.instrumentRepo.FindByExchangeSecurityIDWithLookback(
		ctx,
		instrumentData.ExchangeSecurityID,
		instrumentData.Timestamp,
		string(strategy.ChartTimeframe),
		int(lookbackFromExpressions),
	)

	if err != nil {
		fmt.Printf("Error getting lookback data for instrument %d: %v\n",
			instrumentData.ExchangeSecurityID, err)
		// Use current data point only as fallback
		lookbackData = []*instrumentModel.Instrument{}
	}

	// Reverse lookback data since it comes in descending order (newest first)
	for i, j := 0, len(lookbackData)-1; i < j; i, j = i+1, j-1 {
		lookbackData[i], lookbackData[j] = lookbackData[j], lookbackData[i]
	}

	// Combine lookback data with current data (lookback data first, then current data)
	historicalData := make([]*instrumentModel.Instrument, 0, len(lookbackData)+1)
	historicalData = append(historicalData, lookbackData...)
	historicalData = append(historicalData, instrumentData)

	// Pre-calculate indicators like in backtest engine
	evalContext := parser.NewEvaluationContext(nil)
	evalContext.Strategy = strategyResponse
	evalContext.Data = s.toDataContext(historicalData)

	// Extract indicators from both entry and exit expressions
	entryIndicators, err := parser.ExtractIndicators(strategy.EntryExpression)
	if err != nil {
		fmt.Printf("Error extracting indicators from entry expression: %v\n", err)
		entryIndicators = []parser.IndicatorInfo{}
	}

	var exitIndicators []parser.IndicatorInfo
	if strategy.ExitExpression != "" {
		exitIndicators, err = parser.ExtractIndicators(strategy.ExitExpression)
		if err != nil {
			fmt.Printf("Error extracting indicators from exit expression: %v\n", err)
			exitIndicators = []parser.IndicatorInfo{}
		}
	}

	// Add basic OHLCV indicators explicitly
	basicIndicators := []parser.IndicatorInfo{
		{Name: "OPEN", Arguments: []any{0.0}},
		{Name: "HIGH", Arguments: []any{0.0}},
		{Name: "LOW", Arguments: []any{0.0}},
		{Name: "CLOSE", Arguments: []any{0.0}},
		{Name: "VOLUME", Arguments: []any{0.0}},
	}

	// Combine all indicators (entry + exit + basic)
	allIndicators := append(basicIndicators, entryIndicators...)
	allIndicators = append(allIndicators, exitIndicators...)

	// Pre-calculate all indicators into indicator array map
	indicatorArrayMap := make(map[string][]parser.IndicatorValues)
	calculatedKeys := make(map[string]bool) // Track calculated indicators to avoid duplicates

	for _, indicator := range allIndicators {
		calculator := parser.GetIndFactory(indicator.Name)
		if calculator == nil {
			continue
		}

		uniqueKey, err := calculator.GenerateIndicatorKey(indicator.Arguments, evalContext)
		if err != nil {
			fmt.Printf("Error generating unique key for indicator %s: %v\n", indicator.Name, err)
			continue
		}

		// Skip if already calculated
		if calculatedKeys[uniqueKey] {
			continue
		}
		calculatedKeys[uniqueKey] = true

		// Calculate indicator values using CalculateAll for live trading (single point calculation)
		// Since we only have current data, we calculate for the current timestamp
		indicatorValues := calculator.CalculateAll(indicator.Arguments, evalContext, instrumentData.Timestamp.Add(-24*time.Hour))
		indicatorArrayMap[uniqueKey] = indicatorValues
	}

	currentTimestamp := instrumentData.Timestamp

	fmt.Printf("Current timestamp: %v\n", currentTimestamp)

	// Publish evaluation started event
	// s.publishEventWithPrice(ctx, deployment, strategy, instrumentData, dto.EventTypeEvaluationStarted, dto.EvaluationStartedData{
	// 	MarketPrice: instrumentData.Close,
	// 	Timestamp:   instrumentData.Timestamp,
	// 	Timeframe:   string(strategy.ChartTimeframe),
	// })

	// Track if any action was taken during this evaluation
	actionTaken := false

	// Handle existing position first - check exit conditions
	if currentPosition != nil {
		// fmt.Printf("Found existing open position %s for deployment %s\n",
		// 	currentPosition.ID, deployment.ID)

		// Check for MIS end-of-day square-off first
		if strategy.TradeType == strategyModel.TradeType("MIS") && s.isMISEndOfDayHit(instrumentData) {
			// Force close the MIS position at end of day
			exitPrice := instrumentData.Close
			exitTime := instrumentData.Timestamp

			// Update position with exit details
			currentPosition.ClosePosition(exitPrice, exitTime)

			// Persist the position update
			if err := s.positionRepo.Update(ctx, currentPosition); err != nil {
				fmt.Printf("Error updating closed MIS position: %v\n", err)
			} else {
				actionTaken = true
				fmt.Printf("MIS EOD Square-off: Closed position %s at price %.2f\n",
					currentPosition.ID, exitPrice)

				// Publish EOD square-off event
				s.publishEventWithPrice(ctx, deployment, strategy, instrumentData, dto.EventTypeEODSquareOff, dto.PositionClosedData{
					PositionID: currentPosition.ID.String(),
					ExitPrice:  exitPrice,
					ExitTime:   exitTime,
					ExitReason: "EOD MIS Square-off",
					PNL: func() float64 {
						if currentPosition.PNL != nil {
							return *currentPosition.PNL
						} else {
							return 0.0
						}
					}(),
					Duration:   exitTime.Sub(currentPosition.EntryTime).String(),
					EntryPrice: currentPosition.EntryPrice,
					Quantity:   currentPosition.Quantity,
				})

				// Log the position closure
				logEntry, err := model.NewDeploymentLog(
					deployment.ID,
					model.StatusCompleted,
					map[string]any{
						"position_id": currentPosition.ID,
						"exit_price":  exitPrice,
						"exit_time":   exitTime,
						"exit_reason": "EOD MIS Square-off",
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
		} else {
			// Update trailing stop loss BEFORE other checks
			s.updateStopLossBasedOnTrailing(currentPosition, strategyResponse, instrumentData)

			// Check stop loss and target profit triggers
			stopLossTriggered := s.isStopLossTriggered(currentPosition, instrumentData)
			targetProfitTriggered := s.isTargetProfitTriggered(currentPosition, strategyResponse, instrumentData)

			// Evaluate normal exit condition using indicator array map
			exitSignal, err := s.evaluateExitCondition(strategyResponse.ExitExpression, indicatorArrayMap, currentTimestamp)
			if err != nil {
				fmt.Printf("Error evaluating exit condition: %v\n", err)
			}

			// Close position if any exit condition is met
			if exitSignal || stopLossTriggered || targetProfitTriggered {
				exitPrice := instrumentData.Close
				exitTime := instrumentData.Timestamp
				var exitReason string

				// Determine exit reason
				if stopLossTriggered {
					exitReason = "Stop Loss"
					if currentPosition.InitialStopLossPrice != nil &&
						currentPosition.StopLossPrice != nil &&
						*currentPosition.StopLossPrice != *currentPosition.InitialStopLossPrice {
						exitReason = "Trailing Stop Loss"
					}
				} else if targetProfitTriggered {
					exitReason = "Target Profit"
				} else {
					exitReason = "Exit Signal"
				}

				// Update position with exit details
				currentPosition.ClosePosition(exitPrice, exitTime)

				// Persist the position update
				if err := s.positionRepo.Update(ctx, currentPosition); err != nil {
					fmt.Printf("Error updating closed position: %v\n", err)
				} else {
					actionTaken = true

					println("----------------------------------------")

					fmt.Printf("Closed position %s at price %.2f (Reason: %s)\n",
						currentPosition.ID, exitPrice, exitReason)

					println("----------------------------------------")

					// Publish position closed event
					eventType := dto.EventTypePositionClosed
					if stopLossTriggered {
						eventType = dto.EventTypeStopLossTriggered
					} else if targetProfitTriggered {
						eventType = dto.EventTypeTargetProfitTriggered
					} else {
						eventType = dto.EventTypeExitSignalDetected
					}

					s.publishEventWithPrice(ctx, deployment, strategy, instrumentData, eventType, dto.PositionClosedData{
						PositionID: currentPosition.ID.String(),
						ExitPrice:  exitPrice,
						ExitTime:   exitTime,
						ExitReason: exitReason,
						PNL: func() float64 {
							if currentPosition.PNL != nil {
								return *currentPosition.PNL
							} else {
								return 0.0
							}
						}(),
						Duration:   exitTime.Sub(currentPosition.EntryTime).String(),
						EntryPrice: currentPosition.EntryPrice,
						Quantity:   currentPosition.Quantity,
					})

					// Log the position closure
					logEntry, err := model.NewDeploymentLog(
						deployment.ID,
						model.StatusCompleted,
						map[string]any{
							"position_id": currentPosition.ID,
							"exit_price":  exitPrice,
							"exit_time":   exitTime,
							"exit_reason": exitReason,
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
	}

	// If no current position exists, check entry conditions
	if currentPosition == nil {
		// fmt.Printf("No open position for deployment %s, checking entry conditions\n", deployment.ID)

		// Check if the current time is within the entry time window
		// If EntryStartTime or EntryEndTime is specified, we only look for entries within that window
		canCheckEntryCondition := true

		if strategy.EntryStartTime != nil || strategy.EntryEndTime != nil {
			// Extract time of day from the current instrument timestamp
			currentTime := instrumentData.Timestamp
			currentTimeOfDay := time.Date(0, 1, 1,
				currentTime.Hour(), currentTime.Minute(), currentTime.Second(), 0, time.UTC)

			// Check if we're within the entry time window
			if strategy.EntryStartTime != nil {
				// Extract time component from EntryStartTime
				entryStartTimeOfDay := time.Date(0, 1, 1,
					strategy.EntryStartTime.Hour, strategy.EntryStartTime.Minute,
					strategy.EntryStartTime.Second, 0, time.UTC)

				if currentTimeOfDay.Before(entryStartTimeOfDay) {
					canCheckEntryCondition = false
				}
			}

			if canCheckEntryCondition && strategy.EntryEndTime != nil {
				// Extract time component from EntryEndTime
				entryEndTimeOfDay := time.Date(0, 1, 1,
					strategy.EntryEndTime.Hour, strategy.EntryEndTime.Minute,
					strategy.EntryEndTime.Second, 0, time.UTC)

				if currentTimeOfDay.After(entryEndTimeOfDay) {
					canCheckEntryCondition = false
				}
			}
		}

		// Only check entry conditions if we're within the allowed entry time window
		if canCheckEntryCondition {
			// Evaluate entry condition using indicator array map
			entrySignal, err := s.evaluateEntryCondition(strategyResponse.EntryExpression, indicatorArrayMap, currentTimestamp)
			if err != nil {
				fmt.Printf("Error evaluating entry condition: %v\n", err)
			} else if entrySignal {
				// Publish entry signal detected event
				s.publishEventWithPrice(ctx, deployment, strategy, instrumentData, dto.EventTypeEntrySignalDetected, dto.SignalDetectedData{
					MarketPrice: instrumentData.Close,
					SignalTime:  instrumentData.Timestamp,
					Expression:  strategyResponse.EntryExpression,
					SignalValue: true,
				})

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
					strategy.PositionType, // Add the direction from strategy
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
						actionTaken = true
						println("----------------------------------------")

						fmt.Printf("Created new position %s at price %.2f with quantity %.2f\n",
							newPosition.ID, entryPrice, quantity)

						println("----------------------------------------")

						// Set stop loss and target profit prices
						if strategyResponse.StopLoss > 0 {
							stopLossPrice := s.calculateStopLoss(entryPrice, strategyResponse.StopLoss, strategyResponse.StopLossType, strategy.PositionType)
							if stopLossPrice != nil {
								newPosition.SetStopLossPrice(*stopLossPrice)
							}
						}

						if strategyResponse.TargetProfit > 0 {
							targetProfitPrice := s.calculateTargetProfit(entryPrice, strategyResponse.TargetProfit, strategyResponse.TargetProfitType, strategy.PositionType)
							if targetProfitPrice != nil {
								newPosition.SetTargetProfitPrice(*targetProfitPrice)
							}
						}

						// Update the position in the database with stop loss and target profit
						if err := s.positionRepo.Update(ctx, newPosition); err != nil {
							fmt.Printf("Error updating position with stop loss/target profit: %v\n", err)
						}

						// Publish position opened event
						s.publishEventWithPrice(ctx, deployment, strategy, instrumentData, dto.EventTypePositionOpened, dto.PositionOpenedData{
							PositionID:        newPosition.ID.String(),
							EntryPrice:        entryPrice,
							Quantity:          quantity,
							Direction:         string(strategy.PositionType),
							EntryTime:         entryTime,
							StopLossPrice:     newPosition.StopLossPrice,
							TargetProfitPrice: newPosition.TargetProfitPrice,
						})

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
		} else {
			fmt.Printf("Entry time window check failed for deployment %s - outside allowed entry time\n", deployment.ID)

			// Publish entry window closed event
			var startTime, endTime *time.Time
			if strategy.EntryStartTime != nil {
				t := time.Date(0, 1, 1, strategy.EntryStartTime.Hour, strategy.EntryStartTime.Minute, strategy.EntryStartTime.Second, 0, time.UTC)
				startTime = &t
			}
			if strategy.EntryEndTime != nil {
				t := time.Date(0, 1, 1, strategy.EntryEndTime.Hour, strategy.EntryEndTime.Minute, strategy.EntryEndTime.Second, 0, time.UTC)
				endTime = &t
			}

			s.publishEventWithPrice(ctx, deployment, strategy, instrumentData, dto.EventTypeEntryWindowClosed, dto.TimeWindowEventData{
				WindowType:   "entry",
				WindowStatus: "closed",
				StartTime:    startTime,
				EndTime:      endTime,
				CurrentTime:  instrumentData.Timestamp,
			})
		}
	}

	// Calculate and update deployment metrics
	metrics := s.calculateDeploymentMetrics(ctx, deployment.ID, currentPosition)
	if metrics != nil {
		// Get detailed metrics for proper storage
		totalTrades, winTrades, lossTrades, totalPnL, maxDrawdown, avgWin, avgLoss := s.calculateDeploymentMetricsDetailed(ctx, deployment.ID, currentPosition)

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
				"avg_win":          avgWin,
				"avg_loss":         avgLoss,
			},
		)

		if err != nil {
			fmt.Printf("Error creating deployment metrics: %v\n", err)
		} else {
			// Update the metrics with calculated values
			deploymentMetrics.UpdateMetrics(
				totalTrades,
				winTrades,
				lossTrades,
				totalPnL,
				maxDrawdown,
			)

			// Set average win/loss directly
			if avgWin > 0 {
				deploymentMetrics.AvgWin = &avgWin
			}
			if avgLoss < 0 {
				deploymentMetrics.AvgLoss = &avgLoss
			}

			// Set ratios if available
			if metrics.SharpeRatio != 0 {
				deploymentMetrics.SetRatios(metrics.SharpeRatio, 0.0) // Sortino not calculated yet
			}

			if err := s.metricsRepo.Create(ctx, deploymentMetrics); err != nil {
				fmt.Printf("Error saving deployment metrics: %v\n", err)
			}
		}
	}

	// If no action was taken, publish a no-signal event
	if !actionTaken {
		s.publishEventWithPrice(ctx, deployment, strategy, instrumentData, dto.EventTypeNoSignal, dto.MarketDataReceivedData{
			Price:     instrumentData.Close,
			Volume:    int64(instrumentData.Volume),
			Timestamp: instrumentData.Timestamp,
			OHLC: dto.OHLC{
				Open:  instrumentData.Open,
				High:  instrumentData.High,
				Low:   instrumentData.Low,
				Close: instrumentData.Close,
			},
		})
	}

	// fmt.Printf("Strategy execution completed for deployment %s on instrument %d\n",
	// 	deployment.ID, instrumentData.ExchangeSecurityID)

	return nil
}

// ProcessInstrumentData processes incoming instrument data and runs the appropriate deployments
func (s *deploymentService) ProcessInstrumentData(ctx context.Context, instrumentData *instrumentModel.Instrument) error {
	// Store the data in the database
	// if err := s.instrumentRepo.Create(ctx, instrumentData); err != nil {
	// 	return fmt.Errorf("failed to store instrument data: %w", err)
	// }

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

// evaluateEntryCondition evaluates entry expression using parser with indicator array map
func (s *deploymentService) evaluateEntryCondition(
	expression string,
	indicatorArrayMap map[string][]parser.IndicatorValues,
	currentTimestamp time.Time,
) (bool, error) {
	return parser.EvaluateIndicatorExpression(
		expression, indicatorArrayMap, currentTimestamp, nil, len(indicatorArrayMap["CLOSE"])-1,
	)
}

// calculatePositionSize calculates position size based on strategy settings
func (s *deploymentService) calculatePositionSize(strategy *strategyDto.StrategyResponse, price float64) float64 {
	if strategy.Qty > 0 {
		return float64(strategy.Qty)
	}

	// Default to 1 unit if no quantity specified
	return 1.0
}

// evaluateExitCondition evaluates exit expression using parser with indicator array map
func (s *deploymentService) evaluateExitCondition(
	expression string,
	indicatorArrayMap map[string][]parser.IndicatorValues,
	currentTimestamp time.Time,
) (bool, error) {
	if expression == "" {
		return false, nil
	}

	return parser.EvaluateIndicatorExpression(
		expression, indicatorArrayMap, currentTimestamp, nil, len(indicatorArrayMap["CLOSE"])-1,
	)
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
		// Calculate win rate as decimal (0.0 to 1.0) to match API documentation
		winRate := float64(winningTrades) / float64(totalTrades)
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

// calculateDeploymentMetricsDetailed calculates detailed metrics for deployment metrics storage
func (s *deploymentService) calculateDeploymentMetricsDetailed(ctx context.Context, deploymentID uuid.UUID, currentPosition *model.DeploymentPosition) (totalTrades, winTrades, lossTrades int, totalPnL, maxDrawdown, avgWin, avgLoss float64) {
	// Get all positions for this deployment to calculate metrics
	positions, err := s.positionRepo.GetByDeploymentID(ctx, deploymentID)
	if err != nil {
		fmt.Printf("Error fetching positions for metrics calculation: %v\n", err)
		return 0, 0, 0, 0.0, 0.0, 0.0, 0.0
	}

	totalPnL = 0.0
	winTrades = 0
	lossTrades = 0
	totalTrades = 0
	maxDrawdown = 0.0
	totalWinAmount := 0.0
	totalLossAmount := 0.0

	for _, pos := range positions {
		if pos.Status == model.PositionClosed && pos.PNL != nil {
			totalTrades++
			pnl := *pos.PNL
			totalPnL += pnl

			if pnl > 0 {
				winTrades++
				totalWinAmount += pnl
			} else if pnl < 0 {
				lossTrades++
				totalLossAmount += pnl
			}
		}
	}

	// Calculate average win and loss
	if winTrades > 0 {
		avgWin = totalWinAmount / float64(winTrades)
	}
	if lossTrades > 0 {
		avgLoss = totalLossAmount / float64(lossTrades)
	}

	// TODO: Calculate actual max drawdown based on position history
	// For now, we'll use a simplified calculation
	if totalPnL < 0 {
		maxDrawdown = math.Abs(totalPnL)
	}

	return totalTrades, winTrades, lossTrades, totalPnL, maxDrawdown, avgWin, avgLoss
}

// GetDeploymentStatus gets deployment status for specified instruments against a strategy
func (s *deploymentService) GetDeploymentStatus(ctx context.Context, strategyID uuid.UUID, instrumentIDs []int, userID uuid.UUID) (*dto.DeploymentStatusResponse, error) {
	// Verify the strategy exists and belongs to the user
	strategy, err := s.strategyRepo.GetByID(ctx, strategyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStrategyNotFound, err)
	}

	if strategy.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Get all deployments for this strategy
	deployments, err := s.deploymentRepo.GetByStrategyID(ctx, strategyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get strategy deployments: %w", err)
	}

	// Create a map for quick lookup of deployments by instrument ID
	deploymentMap := make(map[int]*model.Deployment)
	for _, deployment := range deployments {
		deploymentMap[deployment.InstrumentID] = deployment
	}

	// Build response for each requested instrument
	instrumentStatuses := make([]dto.DeploymentInstrumentStatus, len(instrumentIDs))
	for i, instrumentID := range instrumentIDs {
		instrumentStatus := dto.DeploymentInstrumentStatus{
			InstrumentID:  instrumentID,
			HasDeployment: false,
		}

		if deployment, exists := deploymentMap[instrumentID]; exists {
			instrumentStatus.HasDeployment = true
			instrumentStatus.Deployment = &dto.DeploymentResponse{
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

			// Get instrument master data for this deployment
			instrumentMaster, err := s.instrumentMasterRepo.FindByExchangeSecurityID(ctx, int64(deployment.InstrumentID))
			if err == nil && instrumentMaster != nil {
				instrumentStatus.Deployment.InstrumentMaster = instrumentMasterToDTO(instrumentMaster)
			}

			// Get position summary for this deployment
			positions, err := s.positionRepo.GetByDeploymentID(ctx, deployment.ID)
			if err == nil && len(positions) > 0 {
				positionSummary := &dto.PositionSummary{}
				totalPNL := 0.0
				openCount := 0
				closedCount := 0

				for _, pos := range positions {
					if pos.Status == model.PositionOpen {
						openCount++
					} else if pos.Status == model.PositionClosed {
						closedCount++
						if pos.PNL != nil {
							totalPNL += *pos.PNL
						}
					}
				}

				positionSummary.TotalPositions = len(positions)
				positionSummary.OpenPositions = openCount
				positionSummary.ClosedPositions = closedCount
				positionSummary.TotalPNL = totalPNL
				instrumentStatus.PositionSummary = positionSummary
			}

			// Get latest metrics for this deployment
			latestMetrics, err := s.metricsRepo.GetLatestByDeploymentID(ctx, deployment.ID)
			if err == nil && latestMetrics != nil {
				instrumentStatus.LatestMetrics = &dto.MetricsResponse{
					ID:           latestMetrics.ID,
					DeploymentID: latestMetrics.DeploymentID,
					MetricDate:   latestMetrics.MetricDate,
					TotalTrades:  latestMetrics.TotalTrades,
					WinTrades:    latestMetrics.WinTrades,
					LossTrades:   latestMetrics.LossTrades,
					TotalPNL:     latestMetrics.TotalPNL,
					MaxDrawdown:  latestMetrics.MaxDrawdown,
					SharpeRatio:  latestMetrics.SharpeRatio,
					SortinoRatio: latestMetrics.SortinoRatio,
					WinRate:      latestMetrics.WinRate,
					AvgWin:       latestMetrics.AvgWin,
					AvgLoss:      latestMetrics.AvgLoss,
					Metadata:     latestMetrics.Metadata,
					CreatedAt:    latestMetrics.CreatedAt,
					UpdatedAt:    latestMetrics.UpdatedAt,
				}
			}
		}

		instrumentStatuses[i] = instrumentStatus
	}

	response := &dto.DeploymentStatusResponse{
		StrategyID:  strategyID,
		Instruments: instrumentStatuses,
	}

	return response, nil
}

// getLookBackFromExpressions calculates the maximum lookback period required for strategy expressions
func (s *deploymentService) getLookBackFromExpressions(strategy *strategyDto.StrategyResponse) int64 {
	entryLookBack := s.getMaxLookback(strategy.EntryExpression)

	exitLookBack := int64(0)
	if strategy.ExitExpression != "" {
		exitLookBack = s.getMaxLookback(strategy.ExitExpression)
	}

	lookBack := max(entryLookBack, exitLookBack)

	return lookBack
}

// getMaxLookback extracts indicators from the expression and returns the maximum lookback period required
func (s *deploymentService) getMaxLookback(expression string) int64 {
	// Extract all indicators from the expression
	indicators, err := parser.ExtractIndicators(expression)
	if err != nil {
		println("Error extracting indicators from expression:", err.Error())
		return 20 // Return a default lookback if there's an error
	}

	maxLookback := int64(0)

	// Iterate through all indicators and find the maximum lookback
	for _, indicator := range indicators {
		calculator := parser.GetIndFactory(indicator.Name)
		if calculator != nil {
			lookback := calculator.GetLookBack(indicator.Arguments)
			if lookback > maxLookback {
				maxLookback = lookback
			}
		}
	}

	// Return at least 20 as minimum lookback if no indicators found or all have 0 lookback
	if maxLookback < 20 {
		maxLookback = 20
	}

	return maxLookback
}

// calculateStopLoss calculates the stop loss price based on entry price and strategy settings
func (s *deploymentService) calculateStopLoss(
	entryPrice float64,
	stopLoss float64,
	stopLossType strategyModel.ValueType,
	direction strategyModel.Direction,
) *float64 {
	var stopPrice float64
	if direction == strategyModel.Long {
		if stopLossType == "PERCENTAGE" {
			stopPrice = entryPrice * (1 - stopLoss/100)
		} else {
			stopPrice = entryPrice - stopLoss
		}
	} else { // Short position
		if stopLossType == "PERCENTAGE" {
			stopPrice = entryPrice * (1 + stopLoss/100)
		} else {
			stopPrice = entryPrice + stopLoss
		}
	}

	stopPrice = math.Round(stopPrice*100) / 100 // Round to 2 decimal places
	return &stopPrice
}

// calculateTargetProfit calculates the target profit price based on entry price and strategy settings
func (s *deploymentService) calculateTargetProfit(
	entryPrice float64,
	targetProfit float64,
	targetProfitType strategyModel.ValueType,
	direction strategyModel.Direction,
) *float64 {
	var targetPrice float64
	if direction == strategyModel.Long {
		if targetProfitType == "PERCENTAGE" {
			targetPrice = entryPrice * (1 + targetProfit/100)
		} else {
			targetPrice = entryPrice + targetProfit
		}
	} else { // Short position
		if targetProfitType == "PERCENTAGE" {
			targetPrice = entryPrice * (1 - targetProfit/100)
		} else {
			targetPrice = entryPrice - targetProfit
		}
	}

	targetPrice = math.Round(targetPrice*100) / 100 // Round to 2 decimal places
	return &targetPrice
}

// updateStopLossBasedOnTrailing updates stop loss for trailing stop functionality
func (s *deploymentService) updateStopLossBasedOnTrailing(
	position *model.DeploymentPosition,
	strategy *strategyDto.StrategyResponse,
	instrument *instrumentModel.Instrument,
) bool {
	// Skip if trailing stop loss isn't enabled or position doesn't have a stop loss
	if strategy.TrailingStopLoss <= 0 || position.StopLossPrice == nil {
		return false
	}

	currentPrice := instrument.Close
	entryPrice := position.EntryPrice
	currentStopLoss := *position.StopLossPrice

	initialStopDistance := *s.calculateStopLoss(entryPrice, strategy.StopLoss, strategy.StopLossType, position.Direction)

	// Calculate the trailing amount based on type (percentage or absolute)
	trailingAmount := 0.0
	if strategy.TrailingStopLossType == "PERCENTAGE" {
		// For percentage, calculate the equivalent price distance
		if position.Direction == strategyModel.Long {
			trailingAmount = entryPrice * (strategy.TrailingStopLoss / 100.0)
		} else {
			trailingAmount = entryPrice * (strategy.TrailingStopLoss / 100.0)
		}
	} else {
		// For absolute, use the value directly
		trailingAmount = strategy.TrailingStopLoss
	}

	// Handle long position trailing stop
	if position.Direction == strategyModel.Long {
		// Calculate how much price has moved above entry
		priceMovement := currentPrice - entryPrice

		// Only trail if price has moved up enough (entry + trailingAmount)
		if priceMovement >= trailingAmount {
			// Calculate new stop loss: current price - initial stop distance
			unitsJump := math.Floor(priceMovement / trailingAmount)
			newStopLoss := initialStopDistance + (unitsJump * trailingAmount)

			// Only update stop loss if it's higher (better) than current stop loss
			// This ensures the stop loss only moves up, never down
			if newStopLoss > currentStopLoss {
				position.UpdateStopLossPrice(newStopLoss)
				return true
			}
		}
	} else { // Handle short position trailing stop
		// Calculate how much price has moved below entry
		priceMovement := entryPrice - currentPrice

		// Only trail if price has moved down enough (entry - trailingAmount)
		if priceMovement >= trailingAmount {
			// Calculate new stop loss: current price + initial stop distance
			unitsJump := math.Floor(priceMovement / trailingAmount)
			newStopLoss := initialStopDistance - (unitsJump * trailingAmount)

			// Only update stop loss if it's lower (better) than current stop loss
			// This ensures the stop loss only moves down, never up
			if newStopLoss < currentStopLoss {
				position.UpdateStopLossPrice(newStopLoss)
				return true
			}
		}
	}
	return false
}

// isStopLossTriggered checks if stop loss is triggered
func (s *deploymentService) isStopLossTriggered(
	position *model.DeploymentPosition,
	instrument *instrumentModel.Instrument,
) bool {
	if position.StopLossPrice == nil {
		return false
	}

	if position.Direction == strategyModel.Long {
		return instrument.Close <= *position.StopLossPrice
	} else { // Short position
		return instrument.Close >= *position.StopLossPrice
	}
}

// isTargetProfitTriggered checks if target profit is triggered
func (s *deploymentService) isTargetProfitTriggered(
	position *model.DeploymentPosition,
	strategy *strategyDto.StrategyResponse,
	instrument *instrumentModel.Instrument,
) bool {
	if strategy.TargetProfit <= 0 || position.TargetProfitPrice == nil {
		return false
	}

	if position.Direction == strategyModel.Long {
		return instrument.Close >= *position.TargetProfitPrice
	} else { // Short position
		return instrument.Close <= *position.TargetProfitPrice
	}
}

// isMISEndOfDayHit checks if the current time is at or after the MIS end-of-day cutoff time
func (s *deploymentService) isMISEndOfDayHit(data *instrumentModel.Instrument) bool {
	// Use the market hours utility for consistent timing logic
	marketHoursUtil := marketHours.NewMarketHoursUtil()

	// Extract exchange from the context or use a default
	// For now, we'll use "NSE" as the default exchange
	exchange := "NSE"

	return marketHoursUtil.IsMISEndOfDayHit(data.Timestamp, exchange)
}

// isTimestampAlignedWithTimeframe checks if the given timestamp aligns with the specified timeframe
func (s *deploymentService) isTimestampAlignedWithTimeframe(timestamp time.Time, timeframe string) bool {
	minute := timestamp.Minute()

	switch timeframe {
	case "1m":
		// 1-minute timeframe: every minute is valid
		return true
	case "3m":
		// 3-minute timeframe: minutes 0, 3, 6, 9, 12, 15, 18, etc.
		return minute%3 == 0
	case "5m":
		// 5-minute timeframe: minutes 0, 5, 10, 15, 20, etc.
		return minute%5 == 0
	case "10m":
		// 10-minute timeframe: minutes 0, 10, 20, 30, 40, 50
		return minute%10 == 0
	case "15m":
		// 15-minute timeframe: minutes 0, 15, 30, 45
		return minute%15 == 0
	case "30m":
		// 30-minute timeframe: minutes 0, 30
		return minute%30 == 0
	case "1h":
		// 1-hour timeframe: minute must be 0
		return minute == 0
	case "3h":
		// 3-hour timeframe: minute must be 0 and hour must be divisible by 3
		return minute == 0 && timestamp.Hour()%3 == 0
	case "1d":
		// 1-day timeframe: must be start of day (usually 9:15 AM for Indian markets)
		hour := timestamp.Hour()
		return hour == 9 && minute == 15
	default:
		// Unknown timeframe, default to allowing execution
		return true
	}
}

// Helper methods for publishing events

// publishEvent publishes a deployment event to WebSocket clients
func (s *deploymentService) publishEvent(ctx context.Context, deployment *model.Deployment, strategy *strategyModel.Strategy, eventType dto.DeploymentEventType, data interface{}) {
	if s.eventPublisher == nil {
		return // Event publisher not configured
	}

	event := &dto.DeploymentEvent{
		EventType:    eventType,
		DeploymentID: deployment.ID,
		StrategyID:   deployment.StrategyID,
		InstrumentID: deployment.InstrumentID,
		Timestamp:    time.Now(),
		Data:         data,
		Metadata: dto.EventMetadata{
			UserID:       deployment.UserID,
			StrategyName: strategy.Name,
		},
	}

	// Publish to the specific user (non-blocking)
	go func() {
		if err := s.eventPublisher.PublishEventToUser(ctx, deployment.UserID.String(), event); err != nil {
			fmt.Printf("Error publishing event to user %s: %v\n", deployment.UserID.String(), err)
		}
	}()
}

// publishEventWithPrice publishes a deployment event with market price information
func (s *deploymentService) publishEventWithPrice(ctx context.Context, deployment *model.Deployment, strategy *strategyModel.Strategy, instrumentData *instrumentModel.Instrument, eventType dto.DeploymentEventType, data interface{}) {
	if s.eventPublisher == nil {
		return // Event publisher not configured
	}

	marketPrice := instrumentData.Close

	event := &dto.DeploymentEvent{
		EventType:    eventType,
		DeploymentID: deployment.ID,
		StrategyID:   deployment.StrategyID,
		InstrumentID: deployment.InstrumentID,
		Timestamp:    time.Now(),
		Data:         data,
		Metadata: dto.EventMetadata{
			UserID:       deployment.UserID,
			StrategyName: strategy.Name,
			MarketPrice:  &marketPrice,
		},
	}

	// Publish to the specific user (non-blocking)
	go func() {
		if err := s.eventPublisher.PublishEventToUser(ctx, deployment.UserID.String(), event); err != nil {
			fmt.Printf("Error publishing event to user %s: %v\n", deployment.UserID.String(), err)
		}
	}()
}

// Helper function to convert instrument master model to DTO
func instrumentMasterToDTO(master *instrumentModel.InstrumentMaster) *instrumentDto.InstrumentMasterDTO {
	if master == nil {
		return nil
	}

	return &instrumentDto.InstrumentMasterDTO{
		ExchangeSecurityID: master.ExchangeSecurityID,
		Exchange:           master.Exchange,
		Segment:            master.Segment,
		Symbol:             master.Symbol,
		Name:               master.Name,
		StartDate:          master.StartDate,
		Expiry:             master.Expiry,
		ExmExchID:          master.ExmExchID,
		SmstSecurityID:     master.SmstSecurityID,
		InstrumentName:     master.InstrumentName,
		TradingSymbol:      master.TradingSymbol,
		LotUnits:           master.LotUnits,
		StrikePrice:        master.StrikePrice,
		OptionType:         master.OptionType,
		TickSize:           master.TickSize,
		ExpiryFlag:         master.ExpiryFlag,
		ExchInstrumentType: master.ExchInstrumentType,
		SymbolName:         master.SymbolName,
	}
}
