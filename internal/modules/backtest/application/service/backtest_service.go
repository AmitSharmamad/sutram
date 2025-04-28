package service

import (
	"context"
	"errors"
	"time"

	"algo/internal/modules/backtest/application/dto"
	"algo/internal/modules/backtest/domain/model"
	"algo/internal/modules/backtest/domain/repository"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrRunGroupNotFound     = errors.New("run group not found")
	ErrRunNotFound          = errors.New("run not found")
	ErrPositionNotFound     = errors.New("position not found")
	ErrRunMetricsNotFound   = errors.New("run metrics not found")
	ErrInvalidRequest       = errors.New("invalid request")
	ErrRunGroupCreateFailed = errors.New("failed to create run group")
	ErrRunCreateFailed      = errors.New("failed to create run")
	ErrPositionCreateFailed = errors.New("failed to create position")
	ErrMetricsCreateFailed  = errors.New("failed to create run metrics")
)

// BacktestService defines the service for backtest operations
type BacktestService interface {
	// Run Group operations
	CreateRunGroup(ctx context.Context, req *dto.CreateRunGroupRequest) (*dto.RunGroupDTO, error)
	GetRunGroup(ctx context.Context, id uuid.UUID) (*dto.RunGroupDTO, error)
	GetRunGroupsByStrategy(ctx context.Context, strategyID uuid.UUID) ([]*dto.RunGroupDTO, error)
	DeleteRunGroup(ctx context.Context, id uuid.UUID) error

	// Run operations
	CreateRun(ctx context.Context, req *dto.CreateRunRequest) (*dto.RunDTO, error)
	GetRun(ctx context.Context, id uuid.UUID) (*dto.RunDTO, error)
	GetRunsByRunGroup(ctx context.Context, runGroupID uuid.UUID) ([]*dto.RunDTO, error)
	GetRunsByStrategy(ctx context.Context, strategyID uuid.UUID) ([]*dto.RunDTO, error)
	UpdateRunState(ctx context.Context, id uuid.UUID, req *dto.UpdateRunStateRequest) (*dto.RunDTO, error)
	DeleteRun(ctx context.Context, id uuid.UUID) error
	SearchRuns(ctx context.Context, filter *dto.RunFilter) ([]*dto.RunDTO, error)

	// Position operations
	CreatePosition(ctx context.Context, req *dto.CreatePositionRequest) (*dto.PositionDTO, error)
	GetPosition(ctx context.Context, id uuid.UUID) (*dto.PositionDTO, error)
	GetPositionsByRun(ctx context.Context, runID uuid.UUID) ([]*dto.PositionDTO, error)
	EnterPosition(ctx context.Context, id uuid.UUID, req *dto.EnterPositionRequest) (*dto.PositionDTO, error)
	ExitPosition(ctx context.Context, id uuid.UUID, req *dto.ExitPositionRequest) (*dto.PositionDTO, error)
	DeletePosition(ctx context.Context, id uuid.UUID) error

	// Run Metrics operations
	CreateRunMetrics(ctx context.Context, req *dto.CreateRunMetricsRequest) (*dto.RunMetricsDTO, error)
	GetRunMetrics(ctx context.Context, runID uuid.UUID) (*dto.RunMetricsDTO, error)
	UpdateRunMetrics(ctx context.Context, id uuid.UUID, req *dto.CreateRunMetricsRequest) (*dto.RunMetricsDTO, error)
	DeleteRunMetrics(ctx context.Context, id uuid.UUID) error
}

// backtestService implements the BacktestService interface
type backtestService struct {
	runGroupRepo   repository.RunGroupRepository
	runRepo        repository.RunRepository
	positionRepo   repository.PositionRepository
	runMetricsRepo repository.RunMetricsRepository
}

// NewBacktestService creates a new backtest service
func NewBacktestService(
	runGroupRepo repository.RunGroupRepository,
	runRepo repository.RunRepository,
	positionRepo repository.PositionRepository,
	runMetricsRepo repository.RunMetricsRepository,
) BacktestService {
	return &backtestService{
		runGroupRepo:   runGroupRepo,
		runRepo:        runRepo,
		positionRepo:   positionRepo,
		runMetricsRepo: runMetricsRepo,
	}
}

// CreateRunGroup creates a new run group
func (s *backtestService) CreateRunGroup(ctx context.Context, req *dto.CreateRunGroupRequest) (*dto.RunGroupDTO, error) {
	// Create domain model
	runGroup := model.NewRunGroup(req.StrategyID)

	// Save to database
	if err := s.runGroupRepo.Create(ctx, runGroup); err != nil {
		return nil, errors.Join(ErrRunGroupCreateFailed, err)
	}

	// Return DTO
	return &dto.RunGroupDTO{
		ID:         runGroup.ID,
		StrategyID: runGroup.StrategyID,
		CreatedAt:  runGroup.CreatedAt,
		UpdatedAt:  runGroup.UpdatedAt,
	}, nil
}

// GetRunGroup retrieves a run group by ID
func (s *backtestService) GetRunGroup(ctx context.Context, id uuid.UUID) (*dto.RunGroupDTO, error) {
	// Get from database
	runGroup, err := s.runGroupRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Join(ErrRunGroupNotFound, err)
	}

	// Return DTO
	return &dto.RunGroupDTO{
		ID:         runGroup.ID,
		StrategyID: runGroup.StrategyID,
		CreatedAt:  runGroup.CreatedAt,
		UpdatedAt:  runGroup.UpdatedAt,
	}, nil
}

// GetRunGroupsByStrategy retrieves all run groups for a specific strategy
func (s *backtestService) GetRunGroupsByStrategy(ctx context.Context, strategyID uuid.UUID) ([]*dto.RunGroupDTO, error) {
	// Get from database
	runGroups, err := s.runGroupRepo.GetByStrategyID(ctx, strategyID)
	if err != nil {
		return nil, errors.Join(ErrRunGroupNotFound, err)
	}

	// Convert to DTOs
	result := make([]*dto.RunGroupDTO, 0, len(runGroups))
	for _, runGroup := range runGroups {
		result = append(result, &dto.RunGroupDTO{
			ID:         runGroup.ID,
			StrategyID: runGroup.StrategyID,
			CreatedAt:  runGroup.CreatedAt,
			UpdatedAt:  runGroup.UpdatedAt,
		})
	}

	return result, nil
}

// DeleteRunGroup deletes a run group
func (s *backtestService) DeleteRunGroup(ctx context.Context, id uuid.UUID) error {
	// Delete from database
	if err := s.runGroupRepo.Delete(ctx, id); err != nil {
		return errors.Join(ErrRunGroupNotFound, err)
	}

	return nil
}

// CreateRun creates a new run
func (s *backtestService) CreateRun(ctx context.Context, req *dto.CreateRunRequest) (*dto.RunDTO, error) {
	// Create domain model
	run := model.NewRun(
		req.RunGroupID,
		req.StrategyID,
		req.Exchange,
		req.InstrumentID,
		model.RunType(req.Type),
		req.Qty,
		req.PeriodStartTime,
		req.PeriodEndTime,
	)

	// Save to database
	if err := s.runRepo.Create(ctx, run); err != nil {
		return nil, errors.Join(ErrRunCreateFailed, err)
	}

	// Return DTO
	return &dto.RunDTO{
		ID:              run.ID,
		RunGroupID:      run.RunGroupID,
		StrategyID:      run.StrategyID,
		Exchange:        run.Exchange,
		InstrumentID:    run.InstrumentID,
		Type:            dto.RunType(run.Type),
		State:           dto.RunState(run.State),
		Qty:             run.Qty,
		PeriodStartTime: run.PeriodStartTime,
		PeriodEndTime:   run.PeriodEndTime,
		CreatedAt:       run.CreatedAt,
		UpdatedAt:       run.UpdatedAt,
	}, nil
}

// GetRun retrieves a run by ID
func (s *backtestService) GetRun(ctx context.Context, id uuid.UUID) (*dto.RunDTO, error) {
	// Get from database
	run, err := s.runRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Join(ErrRunNotFound, err)
	}

	// Return DTO
	return &dto.RunDTO{
		ID:              run.ID,
		RunGroupID:      run.RunGroupID,
		StrategyID:      run.StrategyID,
		Exchange:        run.Exchange,
		InstrumentID:    run.InstrumentID,
		Type:            dto.RunType(run.Type),
		State:           dto.RunState(run.State),
		Qty:             run.Qty,
		PeriodStartTime: run.PeriodStartTime,
		PeriodEndTime:   run.PeriodEndTime,
		CreatedAt:       run.CreatedAt,
		UpdatedAt:       run.UpdatedAt,
	}, nil
}

// GetRunsByRunGroup retrieves all runs for a specific run group
func (s *backtestService) GetRunsByRunGroup(ctx context.Context, runGroupID uuid.UUID) ([]*dto.RunDTO, error) {
	// Get from database
	runs, err := s.runRepo.GetByRunGroupID(ctx, runGroupID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	result := make([]*dto.RunDTO, 0, len(runs))
	for _, run := range runs {
		result = append(result, &dto.RunDTO{
			ID:              run.ID,
			RunGroupID:      run.RunGroupID,
			StrategyID:      run.StrategyID,
			Exchange:        run.Exchange,
			InstrumentID:    run.InstrumentID,
			Type:            dto.RunType(run.Type),
			State:           dto.RunState(run.State),
			Qty:             run.Qty,
			PeriodStartTime: run.PeriodStartTime,
			PeriodEndTime:   run.PeriodEndTime,
			CreatedAt:       run.CreatedAt,
			UpdatedAt:       run.UpdatedAt,
		})
	}

	return result, nil
}

// GetRunsByStrategy retrieves all runs for a specific strategy
func (s *backtestService) GetRunsByStrategy(ctx context.Context, strategyID uuid.UUID) ([]*dto.RunDTO, error) {
	// Get from database
	runs, err := s.runRepo.GetByStrategyID(ctx, strategyID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	result := make([]*dto.RunDTO, 0, len(runs))
	for _, run := range runs {
		result = append(result, &dto.RunDTO{
			ID:              run.ID,
			RunGroupID:      run.RunGroupID,
			StrategyID:      run.StrategyID,
			Exchange:        run.Exchange,
			InstrumentID:    run.InstrumentID,
			Type:            dto.RunType(run.Type),
			State:           dto.RunState(run.State),
			Qty:             run.Qty,
			PeriodStartTime: run.PeriodStartTime,
			PeriodEndTime:   run.PeriodEndTime,
			CreatedAt:       run.CreatedAt,
			UpdatedAt:       run.UpdatedAt,
		})
	}

	return result, nil
}

// UpdateRunState updates the state of a run
func (s *backtestService) UpdateRunState(ctx context.Context, id uuid.UUID, req *dto.UpdateRunStateRequest) (*dto.RunDTO, error) {
	// Get run
	run, err := s.runRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Join(ErrRunNotFound, err)
	}

	// Update state
	run.State = model.RunState(req.State)
	run.UpdatedAt = time.Now()

	// Update started/ended times based on state
	now := time.Now()
	switch run.State {
	case model.Started:
		run.StartedTime = &now
	case model.Ended:
		run.EndedTime = &now
	}

	// Save to database
	if err := s.runRepo.Update(ctx, run); err != nil {
		return nil, err
	}

	// Convert nullable time fields
	var startedTime, endedTime time.Time
	if run.StartedTime != nil {
		startedTime = *run.StartedTime
	}
	if run.EndedTime != nil {
		endedTime = *run.EndedTime
	}

	// Return updated run
	return &dto.RunDTO{
		ID:              run.ID,
		RunGroupID:      run.RunGroupID,
		StrategyID:      run.StrategyID,
		Exchange:        run.Exchange,
		InstrumentID:    run.InstrumentID,
		Type:            dto.RunType(run.Type),
		State:           dto.RunState(run.State),
		Qty:             run.Qty,
		PeriodStartTime: run.PeriodStartTime,
		PeriodEndTime:   run.PeriodEndTime,
		StartedTime:     startedTime,
		EndedTime:       endedTime,
		CreatedAt:       run.CreatedAt,
		UpdatedAt:       run.UpdatedAt,
	}, nil
}

// DeleteRun deletes a run
func (s *backtestService) DeleteRun(ctx context.Context, id uuid.UUID) error {
	// Delete from database
	if err := s.runRepo.Delete(ctx, id); err != nil {
		return errors.Join(ErrRunNotFound, err)
	}

	return nil
}

// SearchRuns searches runs with filter criteria
func (s *backtestService) SearchRuns(ctx context.Context, filter *dto.RunFilter) ([]*dto.RunDTO, error) {
	// Convert DTO filter to map for repository
	filters := make(map[string]interface{})

	if filter.StrategyID != nil {
		filters["strategy_id"] = *filter.StrategyID
	}

	if filter.Type != nil {
		filters["type"] = *filter.Type
	}

	if filter.State != nil {
		filters["state"] = *filter.State
	}

	if filter.Exchange != nil {
		filters["exchange"] = *filter.Exchange
	}

	if filter.InstrumentID != nil {
		filters["instrument_id"] = *filter.InstrumentID
	}

	// Search in database
	runs, err := s.runRepo.SearchRuns(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	result := make([]*dto.RunDTO, 0, len(runs))
	for _, run := range runs {
		// Create base DTO
		runDTO := &dto.RunDTO{
			ID:              run.ID,
			RunGroupID:      run.RunGroupID,
			StrategyID:      run.StrategyID,
			Exchange:        run.Exchange,
			InstrumentID:    run.InstrumentID,
			Type:            dto.RunType(run.Type),
			State:           dto.RunState(run.State),
			Qty:             run.Qty,
			PeriodStartTime: run.PeriodStartTime,
			PeriodEndTime:   run.PeriodEndTime,
			CreatedAt:       run.CreatedAt,
			UpdatedAt:       run.UpdatedAt,
		}

		// Handle nullable time fields
		if run.StartedTime != nil {
			runDTO.StartedTime = *run.StartedTime
		}

		if run.EndedTime != nil {
			runDTO.EndedTime = *run.EndedTime
		}

		result = append(result, runDTO)
	}

	return result, nil
}

// CreatePosition creates a new position
func (s *backtestService) CreatePosition(ctx context.Context, req *dto.CreatePositionRequest) (*dto.PositionDTO, error) {
	// Create domain model
	position := model.NewPosition(
		req.RunID,
		model.Direction(req.Direction),
		req.EntryPriceSignal,
		req.EntryTimeSignal,
		req.EntrySignal,
		req.PositionSize,
	)

	// Save to database
	if err := s.positionRepo.Create(ctx, position); err != nil {
		return nil, errors.Join(ErrPositionCreateFailed, err)
	}

	// Return DTO
	return s.positionToDTO(position), nil
}

// GetPosition retrieves a position by ID
func (s *backtestService) GetPosition(ctx context.Context, id uuid.UUID) (*dto.PositionDTO, error) {
	// Get from database
	position, err := s.positionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Join(ErrPositionNotFound, err)
	}

	// Return DTO
	return s.positionToDTO(position), nil
}

// GetPositionsByRun retrieves all positions for a specific run
func (s *backtestService) GetPositionsByRun(ctx context.Context, runID uuid.UUID) ([]*dto.PositionDTO, error) {
	// Get from database
	positions, err := s.positionRepo.GetByRunID(ctx, runID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	result := make([]*dto.PositionDTO, 0, len(positions))
	for _, position := range positions {
		result = append(result, s.positionToDTO(position))
	}

	return result, nil
}

// EnterPosition marks a position as entered
func (s *backtestService) EnterPosition(ctx context.Context, id uuid.UUID, req *dto.EnterPositionRequest) (*dto.PositionDTO, error) {
	// Get position
	position, err := s.positionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Join(ErrPositionNotFound, err)
	}

	// Update position
	position.Enter(
		req.EntryPriceOrder,
		req.EntryTimeOrder,
		req.EntryFees,
		req.EntryOrderID,
	)

	// Save to database
	if err := s.positionRepo.Update(ctx, position); err != nil {
		return nil, err
	}

	// Return updated position
	return s.positionToDTO(position), nil
}

// ExitPosition marks a position as exited
func (s *backtestService) ExitPosition(ctx context.Context, id uuid.UUID, req *dto.ExitPositionRequest) (*dto.PositionDTO, error) {
	// Get position
	position, err := s.positionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Join(ErrPositionNotFound, err)
	}

	// Update position
	position.Exit(
		req.ExitPriceSignal,
		req.ExitTimeSignal,
		req.ExitPriceOrder,
		req.ExitTimeOrder,
		req.ExitSignal,
		req.ExitReason,
		req.ExitFees,
		req.ExitOrderID,
	)

	// Save to database
	if err := s.positionRepo.Update(ctx, position); err != nil {
		return nil, err
	}

	// Return updated position
	return s.positionToDTO(position), nil
}

// DeletePosition deletes a position
func (s *backtestService) DeletePosition(ctx context.Context, id uuid.UUID) error {
	// Delete from database
	if err := s.positionRepo.Delete(ctx, id); err != nil {
		return errors.Join(ErrPositionNotFound, err)
	}

	return nil
}

// CreateRunMetrics creates new run metrics
func (s *backtestService) CreateRunMetrics(ctx context.Context, req *dto.CreateRunMetricsRequest) (*dto.RunMetricsDTO, error) {
	// Create domain model
	metrics := model.NewRunMetrics(req.RunID)

	// Set non-pointer values
	metrics.TotalReturnPct = req.TotalReturnPct
	metrics.AnnualizedReturn = req.AnnualizedReturn
	metrics.DailyReturns = req.DailyReturns
	metrics.MonthlyReturns = req.MonthlyReturns
	metrics.SharpeRatio = req.SharpeRatio
	metrics.SortinoRatio = req.SortinoRatio
	metrics.MaxDrawdown = req.MaxDrawdown
	metrics.MaxDrawdownDuration = req.MaxDrawdownDuration
	metrics.StdDeviation = req.StdDeviation
	metrics.TotalTrades = req.TotalTrades
	metrics.WinRate = req.WinRate
	metrics.TotalCommission = req.TotalCommission
	metrics.TotalSlippage = req.TotalSlippage
	metrics.TransactionCostsPct = req.TransactionCostsPct
	metrics.MetricJSON = req.MetricJSON

	// Set pointer values
	valueAtRisk := req.ValueAtRisk
	metrics.ValueAtRisk = &valueAtRisk

	avgWinSize := req.AvgWinSize
	metrics.AvgWinSize = &avgWinSize

	avgLossSize := req.AvgLossSize
	metrics.AvgLossSize = &avgLossSize

	avgTradeDuration := req.AvgTradeDuration
	metrics.AvgTradeDuration = &avgTradeDuration

	profitFactor := req.ProfitFactor
	metrics.ProfitFactor = &profitFactor

	avgPositionSize := req.AvgPositionSize
	metrics.AvgPositionSize = &avgPositionSize

	maxPositionSize := req.MaxPositionSize
	metrics.MaxPositionSize = &maxPositionSize

	positionTurnover := req.PositionTurnover
	metrics.PositionTurnover = &positionTurnover

	avgExposure := req.AvgExposure
	metrics.AvgExposure = &avgExposure

	// Save to database
	if err := s.runMetricsRepo.Create(ctx, metrics); err != nil {
		return nil, errors.Join(ErrMetricsCreateFailed, err)
	}

	// Return DTO
	return s.runMetricsToDTO(metrics), nil
}

// GetRunMetrics retrieves metrics for a specific run
func (s *backtestService) GetRunMetrics(ctx context.Context, runID uuid.UUID) (*dto.RunMetricsDTO, error) {
	// Get from database
	metrics, err := s.runMetricsRepo.GetByRunID(ctx, runID)
	if err != nil {
		return nil, errors.Join(ErrRunMetricsNotFound, err)
	}

	// Return DTO
	return s.runMetricsToDTO(metrics), nil
}

// UpdateRunMetrics updates existing run metrics
func (s *backtestService) UpdateRunMetrics(ctx context.Context, id uuid.UUID, req *dto.CreateRunMetricsRequest) (*dto.RunMetricsDTO, error) {
	// Get metrics
	metrics, err := s.runMetricsRepo.GetByRunID(ctx, req.RunID)
	if err != nil {
		return nil, errors.Join(ErrRunMetricsNotFound, err)
	}

	// Ensure we're updating the correct metrics
	if metrics.ID != id {
		return nil, ErrRunMetricsNotFound
	}

	// Update non-pointer values
	metrics.TotalReturnPct = req.TotalReturnPct
	metrics.AnnualizedReturn = req.AnnualizedReturn
	metrics.DailyReturns = req.DailyReturns
	metrics.MonthlyReturns = req.MonthlyReturns
	metrics.SharpeRatio = req.SharpeRatio
	metrics.SortinoRatio = req.SortinoRatio
	metrics.MaxDrawdown = req.MaxDrawdown
	metrics.MaxDrawdownDuration = req.MaxDrawdownDuration
	metrics.StdDeviation = req.StdDeviation
	metrics.TotalTrades = req.TotalTrades
	metrics.WinRate = req.WinRate
	metrics.TotalCommission = req.TotalCommission
	metrics.TotalSlippage = req.TotalSlippage
	metrics.TransactionCostsPct = req.TransactionCostsPct
	metrics.MetricJSON = req.MetricJSON
	metrics.UpdatedAt = time.Now()

	// Update pointer values
	valueAtRisk := req.ValueAtRisk
	metrics.ValueAtRisk = &valueAtRisk

	avgWinSize := req.AvgWinSize
	metrics.AvgWinSize = &avgWinSize

	avgLossSize := req.AvgLossSize
	metrics.AvgLossSize = &avgLossSize

	avgTradeDuration := req.AvgTradeDuration
	metrics.AvgTradeDuration = &avgTradeDuration

	profitFactor := req.ProfitFactor
	metrics.ProfitFactor = &profitFactor

	avgPositionSize := req.AvgPositionSize
	metrics.AvgPositionSize = &avgPositionSize

	maxPositionSize := req.MaxPositionSize
	metrics.MaxPositionSize = &maxPositionSize

	positionTurnover := req.PositionTurnover
	metrics.PositionTurnover = &positionTurnover

	avgExposure := req.AvgExposure
	metrics.AvgExposure = &avgExposure

	// Save to database
	if err := s.runMetricsRepo.Update(ctx, metrics); err != nil {
		return nil, err
	}

	// Return updated metrics
	return s.runMetricsToDTO(metrics), nil
}

// DeleteRunMetrics deletes run metrics
func (s *backtestService) DeleteRunMetrics(ctx context.Context, id uuid.UUID) error {
	// Delete from database
	if err := s.runMetricsRepo.Delete(ctx, id); err != nil {
		return errors.Join(ErrRunMetricsNotFound, err)
	}

	return nil
}

// Helper method to convert position model to DTO
func (s *backtestService) positionToDTO(position *model.Position) *dto.PositionDTO {
	positionDTO := &dto.PositionDTO{
		ID:               position.ID,
		RunID:            position.RunID,
		State:            dto.PositionState(position.State),
		Direction:        dto.Direction(position.Direction),
		EntryPriceSignal: position.EntryPriceSignal,
		EntryTimeSignal:  position.EntryTimeSignal,
		EntrySignal:      position.EntrySignal,
		PositionSize:     position.PositionSize,
		EntryOrderID:     position.EntryOrderID,
		ExitOrderID:      position.ExitOrderID,
		CreatedAt:        position.CreatedAt,
		UpdatedAt:        position.UpdatedAt,
	}

	// Handle nullable fields
	if position.EntryPriceOrder != nil {
		positionDTO.EntryPriceOrder = *position.EntryPriceOrder
	}

	if position.EntryTimeOrder != nil {
		positionDTO.EntryTimeOrder = *position.EntryTimeOrder
	}

	if position.ExitPriceSignal != nil {
		positionDTO.ExitPriceSignal = *position.ExitPriceSignal
	}

	if position.ExitTimeSignal != nil {
		positionDTO.ExitTimeSignal = *position.ExitTimeSignal
	}

	if position.ExitPriceOrder != nil {
		positionDTO.ExitPriceOrder = *position.ExitPriceOrder
	}

	if position.ExitTimeOrder != nil {
		positionDTO.ExitTimeOrder = *position.ExitTimeOrder
	}

	if position.ExitSignal != nil {
		positionDTO.ExitSignal = *position.ExitSignal
	}

	if position.ExitReason != nil {
		positionDTO.ExitReason = *position.ExitReason
	}

	if position.ProfitLossAbsolute != nil {
		positionDTO.ProfitLossAbsolute = *position.ProfitLossAbsolute
	}

	if position.ProfitLossPercentage != nil {
		positionDTO.ProfitLossPercentage = *position.ProfitLossPercentage
	}

	if position.TransactionCost != nil {
		positionDTO.TransactionCost = *position.TransactionCost
	}

	if position.EntryFees != nil {
		positionDTO.EntryFees = *position.EntryFees
	}

	if position.ExitFees != nil {
		positionDTO.ExitFees = *position.ExitFees
	}

	if position.NetProfitLoss != nil {
		positionDTO.NetProfitLoss = *position.NetProfitLoss
	}

	return positionDTO
}

// Helper method to convert run metrics model to DTO
func (s *backtestService) runMetricsToDTO(metrics *model.RunMetrics) *dto.RunMetricsDTO {
	metricsDTO := &dto.RunMetricsDTO{
		ID:                  metrics.ID,
		RunID:               metrics.RunID,
		TotalReturnPct:      metrics.TotalReturnPct,
		AnnualizedReturn:    metrics.AnnualizedReturn,
		DailyReturns:        metrics.DailyReturns,
		MonthlyReturns:      metrics.MonthlyReturns,
		SharpeRatio:         metrics.SharpeRatio,
		SortinoRatio:        metrics.SortinoRatio,
		MaxDrawdown:         metrics.MaxDrawdown,
		MaxDrawdownDuration: metrics.MaxDrawdownDuration,
		StdDeviation:        metrics.StdDeviation,
		TotalTrades:         metrics.TotalTrades,
		WinRate:             metrics.WinRate,
		TotalCommission:     metrics.TotalCommission,
		TotalSlippage:       metrics.TotalSlippage,
		TransactionCostsPct: metrics.TransactionCostsPct,
		MetricJSON:          metrics.MetricJSON,
		CreatedAt:           metrics.CreatedAt,
		UpdatedAt:           metrics.UpdatedAt,
	}

	// Handle nullable fields
	if metrics.RiskAdjustedReturn != nil {
		metricsDTO.RiskAdjustedReturn = *metrics.RiskAdjustedReturn
	}

	if metrics.ValueAtRisk != nil {
		metricsDTO.ValueAtRisk = *metrics.ValueAtRisk
	}

	if metrics.DownsideDeviation != nil {
		metricsDTO.DownsideDeviation = *metrics.DownsideDeviation
	}

	if metrics.AvgWinSize != nil {
		metricsDTO.AvgWinSize = *metrics.AvgWinSize
	}

	if metrics.AvgLossSize != nil {
		metricsDTO.AvgLossSize = *metrics.AvgLossSize
	}

	if metrics.AvgTradeDuration != nil {
		metricsDTO.AvgTradeDuration = *metrics.AvgTradeDuration
	}

	if metrics.ProfitFactor != nil {
		metricsDTO.ProfitFactor = *metrics.ProfitFactor
	}

	if metrics.AvgPositionSize != nil {
		metricsDTO.AvgPositionSize = *metrics.AvgPositionSize
	}

	if metrics.MaxPositionSize != nil {
		metricsDTO.MaxPositionSize = *metrics.MaxPositionSize
	}

	if metrics.PositionTurnover != nil {
		metricsDTO.PositionTurnover = *metrics.PositionTurnover
	}

	if metrics.AvgExposure != nil {
		metricsDTO.AvgExposure = *metrics.AvgExposure
	}

	return metricsDTO
}
