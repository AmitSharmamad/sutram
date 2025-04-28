package service

import (
	"context"
	stdmath "math"
	"sync"
	"time"

	backtest "algo/internal/modules/backtest/domain/model"
	"algo/internal/modules/backtest/domain/repository"
	strategymodel "algo/internal/modules/strategy/domain/model"

	"github.com/google/uuid"
)

// EngineService defines the service for backtest engine operations
type EngineService interface {
	// Start a backtest run for the provided strategy
	StartBacktest(ctx context.Context, strategyID uuid.UUID) (*backtest.RunGroup, error)
}

// BacktestEngine implements the engine for running backtests
type BacktestEngine struct {
	backtestService   *backtestService
	strategyService   StrategyService
	stockService      StockService
	expressionService ExpressionService
	indicatorService  IndicatorService
	runGroupRepo      repository.RunGroupRepository
	runRepo           repository.RunRepository
	positionRepo      repository.PositionRepository
	runMetricsRepo    repository.RunMetricsRepository
}

// Define Stock type for the engine
type Stock struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
	Direction int
}

// Interface definitions for external services
type StrategyService interface {
	Get(ctx context.Context, id uuid.UUID) *strategymodel.Strategy
}

type StockService interface {
	Get(instrumentID uint32, startTime, endTime time.Time, timeframe string, lookback int32) (*[]Stock, error)
}

type ExpressionService interface {
	GetEntryDataRange(strategy *strategymodel.Strategy) int32
	GetExitDataRange(strategy *strategymodel.Strategy) int32
}

type IndicatorService interface {
	// Add indicator-related methods as needed
}

// NewBacktestEngine creates a new backtest engine
func NewBacktestEngine(
	backtestService *backtestService,
	strategyService StrategyService,
	stockService StockService,
	expressionService ExpressionService,
	indicatorService IndicatorService,
	runGroupRepo repository.RunGroupRepository,
	runRepo repository.RunRepository,
	positionRepo repository.PositionRepository,
	runMetricsRepo repository.RunMetricsRepository,
) EngineService {
	return &BacktestEngine{
		backtestService:   backtestService,
		strategyService:   strategyService,
		stockService:      stockService,
		expressionService: expressionService,
		indicatorService:  indicatorService,
		runGroupRepo:      runGroupRepo,
		runRepo:           runRepo,
		positionRepo:      positionRepo,
		runMetricsRepo:    runMetricsRepo,
	}
}

// StartBacktest runs a backtest for the provided strategy
func (engine *BacktestEngine) StartBacktest(ctx context.Context, strategyID uuid.UUID) (*backtest.RunGroup, error) {
	strategy := engine.strategyService.Get(ctx, strategyID)

	// Create a new run group
	runGroup := backtest.NewRunGroup(strategyID)
	if err := engine.runGroupRepo.Create(ctx, runGroup); err != nil {
		return nil, err
	}

	instruments := strategy.Instruments
	var wg sync.WaitGroup
	runs := make([]*backtest.Run, 0, len(instruments))

	// Get period times from strategy
	// Default to 30 days ago to now if not available from strategy
	periodStart := time.Now().Add(-30 * 24 * time.Hour)
	periodEnd := time.Now()

	// Note: This is an adaptation since the actual Strategy model may have different field names
	// Adjust these field accesses according to your actual Strategy model structure
	if strategy.StartDate != nil {
		periodStart = *strategy.StartDate
	}

	if strategy.EndDate != nil {
		periodEnd = *strategy.EndDate
	}

	for _, instrumentID := range instruments {
		// Create a new run for each instrument
		run := backtest.NewRun(
			runGroup.ID,
			strategyID,
			"",                // Exchange placeholder - would need to be populated from strategy
			int(instrumentID), // Convert to int
			backtest.Backtest,
			1, // Default quantity
			periodStart,
			periodEnd,
		)

		runs = append(runs, run)
		if err := engine.runRepo.Create(ctx, run); err != nil {
			return nil, err
		}

		wg.Add(1)
		// Execute run asynchronously
		go engine.executeRun(ctx, run, strategy, &wg)
	}

	wg.Wait()

	// Persist positions and metrics
	positions := make([]*backtest.Position, 0)
	runMetrics := make([]*backtest.RunMetrics, 0)

	for _, run := range runs {
		// We need to retrieve run from database to get latest state
		updatedRun, err := engine.runRepo.GetByID(ctx, run.ID)
		if err != nil {
			continue
		}

		// Get positions and metrics
		runPositions, err := engine.positionRepo.GetByRunID(ctx, updatedRun.ID)
		if err == nil && len(runPositions) > 0 {
			positions = append(positions, runPositions...)
		}

		metrics, err := engine.runMetricsRepo.GetByRunID(ctx, updatedRun.ID)
		if err == nil {
			runMetrics = append(runMetrics, metrics)
		}
	}

	return runGroup, nil
}

func (engine *BacktestEngine) executeRun(
	ctx context.Context,
	run *backtest.Run,
	strategy *strategymodel.Strategy,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	// Update run state to Started
	run.State = backtest.Started
	startTime := time.Now()
	run.StartedTime = &startTime
	run.UpdatedAt = startTime

	if err := engine.runRepo.Update(ctx, run); err != nil {
		return
	}

	// Set default lookback
	var lookBack int32 = 13

	// Get timeframe data for entry conditions
	// Convert chart timeframe to string
	timeframeStr := string(strategy.ChartTimeframe)

	timeframeData, err := engine.stockService.Get(
		uint32(run.InstrumentID),
		run.PeriodStartTime,
		run.PeriodEndTime,
		timeframeStr,
		lookBack,
	)

	if err != nil || timeframeData == nil || len(*timeframeData) == 0 {
		// Handle case where data is not available
		engine.finishRun(ctx, run)
		return
	}

	// Get 1-second data for more granular exit condition checking
	secondData, _ := engine.stockService.Get(
		uint32(run.InstrumentID),
		run.PeriodStartTime,
		run.PeriodEndTime,
		"1s", // 1-second timeframe
		lookBack,
	)

	// Process entry and exit conditions
	positions := engine.processBacktest(run, strategy, timeframeData, secondData, lookBack)

	// Save positions
	for _, position := range positions {
		if err := engine.positionRepo.Create(ctx, position); err != nil {
			continue
		}
	}

	// Calculate metrics
	metrics := engine.calculateRunMetrics(run, positions)
	if err := engine.runMetricsRepo.Create(ctx, metrics); err != nil {
		// Log error but continue
	}

	// Mark run as complete
	engine.finishRun(ctx, run)
}

func (engine *BacktestEngine) finishRun(ctx context.Context, run *backtest.Run) {
	run.State = backtest.Ended
	endTime := time.Now()
	run.EndedTime = &endTime
	run.UpdatedAt = endTime

	// Save completed run
	_ = engine.runRepo.Update(ctx, run)
}

func (engine *BacktestEngine) processBacktest(
	run *backtest.Run,
	strategy *strategymodel.Strategy,
	timeframeData *[]Stock,
	secondData *[]Stock,
	lookBack int32,
) []*backtest.Position {
	positions := make([]*backtest.Position, 0)
	var currentPosition *backtest.Position

	// Return empty if data is missing
	if timeframeData == nil || len(*timeframeData) == 0 {
		return positions
	}

	// Main loop through timeframe data for entry signals
	for i := lookBack; i < int32(len(*timeframeData)); i++ {
		currentTimeframe := (*timeframeData)[i]
		slicedTimeframeData := (*timeframeData)[:i+1]

		// Check if we need a new position or use existing position
		if currentPosition == nil || currentPosition.State == backtest.Exited {
			// Create a new position
			direction := backtest.Long
			if strategy.PositionType == "SHORT" {
				direction = backtest.Short
			}

			currentPosition = backtest.NewPosition(
				run.ID,
				direction,
				currentTimeframe.Close,
				currentTimeframe.Timestamp,
				"entry_condition",
				currentTimeframe.Close, // Using close price as position size
			)
			positions = append(positions, currentPosition)
		}

		if currentPosition.State == backtest.NotEntered {
			// Check entry signal
			entrySignal := engine.checkEntryCondition(currentPosition, strategy, &currentTimeframe, &slicedTimeframeData)

			if entrySignal && secondData != nil && len(*secondData) > 0 {
				// Position was entered, determine the boundary time for exit checking
				var nextTimeframeBoundary time.Time
				if i < int32(len(*timeframeData))-1 {
					nextTimeframeBoundary = (*timeframeData)[i+1].Timestamp
				} else {
					nextTimeframeBoundary = run.PeriodEndTime
				}

				// Check for exit using second data
				engine.processExit(currentPosition, strategy, secondData, currentTimeframe.Timestamp, nextTimeframeBoundary)

				// Optimize by skipping to next relevant timeframe after exit
				if currentPosition.State == backtest.Exited && currentPosition.ExitTimeOrder != nil {
					exitTime := *currentPosition.ExitTimeOrder
					for j := i + 1; j < int32(len(*timeframeData)); j++ {
						if (*timeframeData)[j].Timestamp.After(exitTime) {
							i = j - 1 // -1 because the loop will increment i
							break
						}
					}
				}
			}
		}
	}

	// Handle the last open position
	if len(positions) > 0 {
		lastPosition := positions[len(positions)-1]
		if lastPosition.State == backtest.Entered {
			// Force exit the last position
			lastData := (*timeframeData)[len(*timeframeData)-1]
			exitPrice := lastData.Close
			exitTime := lastData.Timestamp
			exitReason := "end_of_backtest"

			// Create a nil UUID pointer for the order ID
			var nilOrderID *uuid.UUID

			lastPosition.Exit(
				exitPrice,
				exitTime,
				exitPrice,
				exitTime,
				exitReason,
				exitReason,
				0, // No fees in backtest
				nilOrderID,
			)
		} else if lastPosition.State == backtest.NotEntered {
			// Remove the position that was never entered
			positions = positions[:len(positions)-1]
		}
	}

	return positions
}

func (engine *BacktestEngine) checkEntryCondition(
	position *backtest.Position,
	strategy *strategymodel.Strategy,
	current *Stock,
	data *[]Stock,
) bool {
	// In the actual implementation, use the parser module to evaluate conditions
	entryCondition := engine.evaluateExpression(strategy.EntryExpression, data)

	if entryCondition {
		// Create a nil UUID pointer for the order ID
		var nilOrderID *uuid.UUID

		// Mark position as entered
		position.Enter(
			current.Close,
			current.Timestamp,
			0, // No fees in backtest
			nilOrderID,
		)
		return true
	}

	return false
}

// Helper method to evaluate expressions using the parser
func (engine *BacktestEngine) evaluateExpression(expression string, data *[]Stock) bool {
	// This is a simplified placeholder implementation
	// In the actual implementation, use the parser module to evaluate the expression
	return true // For demonstration purposes
}

func (engine *BacktestEngine) processExit(
	position *backtest.Position,
	strategy *strategymodel.Strategy,
	secondData *[]Stock,
	startTime time.Time,
	endTime time.Time,
) {
	// Return if no second data
	if secondData == nil || len(*secondData) == 0 {
		return
	}

	// Find the starting point in the second data
	var startIndex int32 = -1
	for i := int32(0); i < int32(len(*secondData)); i++ {
		if !(*secondData)[i].Timestamp.Before(startTime) {
			startIndex = i
			break
		}
	}

	// If no valid start point, use fallback
	if startIndex == -1 {
		return
	}

	// Find the end index
	var endIndex int32 = int32(len(*secondData)) - 1
	for i := startIndex; i < int32(len(*secondData)); i++ {
		if !(*secondData)[i].Timestamp.Before(endTime) {
			endIndex = i - 1
			break
		}
	}

	// If invalid range, return
	if endIndex < startIndex {
		return
	}

	// Process each second data point
	for i := startIndex; i <= endIndex; i++ {
		current := (*secondData)[i]

		// Skip if position already exited
		if position.State == backtest.Exited {
			break
		}

		// Slice the data for context
		slicedData := (*secondData)[:i+1]

		// Try to exit the position
		engine.checkExitCondition(position, strategy, &current, &slicedData)
	}
}

func (engine *BacktestEngine) checkExitCondition(
	position *backtest.Position,
	strategy *strategymodel.Strategy,
	current *Stock,
	data *[]Stock,
) {
	// Calculate stop loss and take profit levels
	var stopLossPrice float64
	if strategy.StopLoss > 0 {
		if position.Direction == backtest.Long {
			// Handle pointer dereference
			if position.EntryPriceOrder != nil {
				entryPrice := *position.EntryPriceOrder // Correctly dereference pointer
				stopLossPrice = entryPrice * (1 - strategy.StopLoss/100)
			}
		} else {
			// Handle pointer dereference
			if position.EntryPriceOrder != nil {
				entryPrice := *position.EntryPriceOrder // Correctly dereference pointer
				stopLossPrice = entryPrice * (1 + strategy.StopLoss/100)
			}
		}
	}

	var targetProfitPrice float64
	if strategy.TargetProfit > 0 {
		if position.Direction == backtest.Long {
			// Handle pointer dereference
			if position.EntryPriceOrder != nil {
				entryPrice := *position.EntryPriceOrder // Correctly dereference pointer
				targetProfitPrice = entryPrice * (1 + strategy.TargetProfit/100)
			}
		} else {
			// Handle pointer dereference
			if position.EntryPriceOrder != nil {
				entryPrice := *position.EntryPriceOrder // Correctly dereference pointer
				targetProfitPrice = entryPrice * (1 - strategy.TargetProfit/100)
			}
		}
	}

	// Check exit conditions using the parser
	exitConditionResult := engine.evaluateExpression(strategy.ExitExpression, data)

	// Determine if we should exit
	var shouldExit bool
	var exitReason string

	// Check stop loss
	if strategy.StopLoss > 0 {
		if (position.Direction == backtest.Long && current.Close <= stopLossPrice) ||
			(position.Direction == backtest.Short && current.Close >= stopLossPrice) {
			shouldExit = true
			exitReason = "stop_loss"
		}
	}

	// Check take profit
	if !shouldExit && strategy.TargetProfit > 0 {
		if (position.Direction == backtest.Long && current.Close >= targetProfitPrice) ||
			(position.Direction == backtest.Short && current.Close <= targetProfitPrice) {
			shouldExit = true
			exitReason = "take_profit"
		}
	}

	// Check custom exit condition
	if !shouldExit && exitConditionResult {
		shouldExit = true
		exitReason = "exit_condition"
	}

	// Exit the position if needed
	if shouldExit {
		// Create a nil UUID pointer for the order ID
		var nilOrderID *uuid.UUID

		position.Exit(
			current.Close,
			current.Timestamp,
			current.Close,
			current.Timestamp,
			exitReason,
			exitReason,
			0, // No fees in backtest
			nilOrderID,
		)
	}
}

func (engine *BacktestEngine) calculateRunMetrics(run *backtest.Run, positions []*backtest.Position) *backtest.RunMetrics {
	metrics := backtest.NewRunMetrics(run.ID)

	if len(positions) == 0 {
		return metrics
	}

	// Calculate basic metrics
	var totalPnL float64
	var totalWins int
	var totalLosses int
	var totalTradeDuration float64
	var winAmount float64
	var lossAmount float64
	var maxDrawdown float64
	var peakValue float64
	var maxPositionSize float64
	var totalPositionSize float64

	// Daily returns for standard deviation
	dailyReturns := make([]float64, 0, len(positions))

	// Process each position
	for i, pos := range positions {
		if pos.State != backtest.Exited || pos.EntryTimeOrder == nil || pos.ExitTimeOrder == nil {
			continue
		}

		// Position size
		posSize := pos.PositionSize
		if posSize > maxPositionSize {
			maxPositionSize = posSize
		}
		totalPositionSize += posSize

		// Calculate P&L - handle pointer dereferencing
		var pnl float64
		if pos.Direction == backtest.Long && pos.ExitPriceOrder != nil && pos.EntryPriceOrder != nil {
			// Dereference the pointers for calculation
			exitPrice := *pos.ExitPriceOrder
			entryPrice := *pos.EntryPriceOrder
			pnl = exitPrice - entryPrice
		} else if pos.ExitPriceOrder != nil && pos.EntryPriceOrder != nil {
			// Dereference the pointers for calculation
			exitPrice := *pos.ExitPriceOrder
			entryPrice := *pos.EntryPriceOrder
			pnl = entryPrice - exitPrice
		}
		totalPnL += pnl

		// Calculate return
		dailyReturn := pnl / posSize
		dailyReturns = append(dailyReturns, dailyReturn)

		// Update drawdown tracking
		currentValue := totalPnL
		if i == 0 || currentValue > peakValue {
			peakValue = currentValue
		}
		currentDrawdown := peakValue - currentValue
		if currentDrawdown > maxDrawdown {
			maxDrawdown = currentDrawdown
		}

		// Win/loss stats
		if (pos.Direction == backtest.Long && pnl > 0) || (pos.Direction == backtest.Short && pnl < 0) {
			totalWins++
			winAmount += stdmath.Abs(pnl)
		} else {
			totalLosses++
			lossAmount += stdmath.Abs(pnl)
		}

		// Trade duration
		duration := pos.ExitTimeOrder.Sub(*pos.EntryTimeOrder)
		totalTradeDuration += duration.Hours()
	}

	// Set calculated metrics
	totalTrades := totalWins + totalLosses
	if totalTrades > 0 {
		metrics.TotalTrades = int(totalTrades) // Convert to int as required by model
		metrics.WinRate = float64(totalWins) / float64(totalTrades)

		avgTradeDuration := totalTradeDuration / float64(totalTrades)
		metrics.AvgTradeDuration = &avgTradeDuration

		maxPositionSizeVal := maxPositionSize
		metrics.MaxPositionSize = &maxPositionSizeVal

		avgPositionSize := totalPositionSize / float64(totalTrades)
		metrics.AvgPositionSize = &avgPositionSize

		// P&L metrics
		if totalPositionSize > 0 {
			metrics.TotalReturnPct = (totalPnL / totalPositionSize) * 100

			// Calculate annualized return
			if run.StartedTime != nil && run.EndedTime != nil {
				durationDays := run.EndedTime.Sub(*run.StartedTime).Hours() / 24
				if durationDays > 0 {
					annualFactor := 365.0 / durationDays
					metrics.AnnualizedReturn = (stdmath.Pow(1+totalPnL/totalPositionSize, annualFactor) - 1) * 100
				}
			}
		}

		// Win/loss metrics
		if totalWins > 0 {
			avgWinSize := winAmount / float64(totalWins)
			metrics.AvgWinSize = &avgWinSize
		}

		if totalLosses > 0 {
			avgLossSize := lossAmount / float64(totalLosses)
			metrics.AvgLossSize = &avgLossSize
		}

		if totalLosses > 0 && lossAmount > 0 {
			profitFactor := winAmount / lossAmount
			metrics.ProfitFactor = &profitFactor
		}

		// Risk metrics
		metrics.MaxDrawdown = maxDrawdown
		if peakValue > 0 {
			metrics.MaxDrawdownDuration = maxDrawdown / peakValue * 100
		}

		// Calculate standard deviation
		if len(dailyReturns) > 0 {
			mean := 0.0
			for _, ret := range dailyReturns {
				mean += ret
			}
			mean /= float64(len(dailyReturns))

			variance := 0.0
			for _, ret := range dailyReturns {
				variance += (ret - mean) * (ret - mean)
			}
			variance /= float64(len(dailyReturns))
			stdDev := stdmath.Sqrt(variance)
			metrics.StdDeviation = stdDev

			// Calculate Sharpe ratio (assuming risk-free rate of 0)
			if stdDev > 0 {
				metrics.SharpeRatio = mean / stdDev * stdmath.Sqrt(252)
			}

			// Calculate Sortino ratio
			downsideVariance := 0.0
			for _, ret := range dailyReturns {
				if ret < 0 {
					downsideVariance += ret * ret
				}
			}
			downsideVariance /= float64(len(dailyReturns))
			downsideDeviation := stdmath.Sqrt(downsideVariance)

			if downsideDeviation > 0 {
				metrics.SortinoRatio = mean / downsideDeviation * stdmath.Sqrt(252)
			}
		}
	}

	return metrics
}
