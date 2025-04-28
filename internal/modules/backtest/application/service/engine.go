package service

import (
	"context"
	stdmath "math"
	"strconv"
	"sync"
	"time"

	"algo/internal/modules/backtest/domain/model"
	backtestRepo "algo/internal/modules/backtest/domain/repository"
	instrumentService "algo/internal/modules/instrument/application/service"
	instrumentModel "algo/internal/modules/instrument/domain/model"
	strategyDto "algo/internal/modules/strategy/application/dto"
	strategyService "algo/internal/modules/strategy/application/service"

	"github.com/google/uuid"
)

// EngineService defines the service for backtest engine operations
type EngineService interface {
	// Start a backtest run for the provided strategy
	StartBacktest(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) (*model.RunGroup, error)
}

// BacktestEngine implements the engine for running backtests
type BacktestEngine struct {
	backtestService   *BacktestService
	strategyService   strategyService.StrategyService
	instrumentService *instrumentService.InstrumentService
	runGroupRepo      backtestRepo.RunGroupRepository
	runRepo           backtestRepo.RunRepository
	positionRepo      backtestRepo.PositionRepository
	runMetricsRepo    backtestRepo.RunMetricsRepository
}

// Instrument represents OHLCV data for backtest processing
type Instrument struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// NewBacktestEngine creates a new backtest engine
func NewBacktestEngine(
	backtestService *BacktestService,
	strategyService strategyService.StrategyService,
	instrumentService *instrumentService.InstrumentService,
	runGroupRepo backtestRepo.RunGroupRepository,
	runRepo backtestRepo.RunRepository,
	positionRepo backtestRepo.PositionRepository,
	runMetricsRepo backtestRepo.RunMetricsRepository,
) EngineService {
	return &BacktestEngine{
		backtestService:   backtestService,
		strategyService:   strategyService,
		instrumentService: instrumentService,
		runGroupRepo:      runGroupRepo,
		runRepo:           runRepo,
		positionRepo:      positionRepo,
		runMetricsRepo:    runMetricsRepo,
	}
}

// StartBacktest starts a backtest run for the provided strategy
func (e *BacktestEngine) StartBacktest(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) (*model.RunGroup, error) {
	// Get the strategy
	strategyResp, err := e.strategyService.GetStrategyByID(ctx, strategyID, userID)
	if err != nil {
		return nil, err
	}

	// Create a run group
	runGroup := model.NewRunGroup(strategyID)
	if err := e.runGroupRepo.Create(ctx, runGroup); err != nil {
		return nil, err
	}

	// Process each instrument in the strategy
	var wg sync.WaitGroup
	for _, instrumentID := range strategyResp.Instruments {
		wg.Add(1)

		go func(instrumentID int32) {
			defer wg.Done()

			// Convert int32 to string for the symbol lookup
			symbol := strconv.FormatInt(int64(instrumentID), 10)

			// Fetch instrument data
			instrumentMaster, err := e.instrumentService.GetInstrumentMasterBySymbol(ctx, symbol)
			if err != nil {
				// Log error and continue with other instruments
				return
			}

			// Get historical data for the instrument
			// Handle nil time values for date range
			var startDate, endDate time.Time
			if strategyResp.StartDate != nil {
				startDate = *strategyResp.StartDate
			} else {
				startDate = time.Now().AddDate(-1, 0, 0) // Default to 1 year back
			}

			if strategyResp.EndDate != nil {
				endDate = *strategyResp.EndDate
			} else {
				endDate = time.Now()
			}

			instruments, err := e.instrumentService.GetInstrumentsByTimeRange(ctx, symbol, startDate, endDate)
			if err != nil {
				// Log error and continue with other instruments
				return
			}

			// Create a run for this instrument
			// Assuming instrumentMaster.ID is uint64, convert to int for Run
			instrumentIDInt := int(instrumentMaster.ID)
			qty := 1 // Default quantity
			if strategyResp.Qty > 0 {
				qty = strategyResp.Qty
			}

			run := model.NewRun(
				runGroup.ID,
				strategyID,
				instrumentMaster.Exchange,
				instrumentIDInt,
				model.Backtest,
				qty,
				startDate,
				endDate,
			)

			if err := e.runRepo.Create(ctx, run); err != nil {
				// Log error and continue with other instruments
				return
			}

			// Convert instruments to engine instrument format
			instrumentData := e.getInstrumentById(instruments)

			// Execute backtest for this instrument
			e.executeBacktest(ctx, strategyResp, instrumentData, run)
		}(instrumentID)
	}

	wg.Wait()
	return runGroup, nil
}

// getInstrumentById converts instrument domain model to engine instrument model
func (e *BacktestEngine) getInstrumentById(instruments []*instrumentModel.Instrument) []Instrument {
	instrumentData := make([]Instrument, len(instruments))
	for i, inst := range instruments {
		instrumentData[i] = Instrument{
			Timestamp: inst.Timestamp,
			Open:      inst.Open,
			High:      inst.High,
			Low:       inst.Low,
			Close:     inst.Close,
			Volume:    float64(inst.Volume),
		}
	}
	return instrumentData
}

// executeBacktest runs the backtest for a single instrument
func (e *BacktestEngine) executeBacktest(ctx context.Context, strategy *strategyDto.StrategyResponse, instruments []Instrument, run *model.Run) {
	// Initialize a parser for entry/exit expressions (assuming a factory exists in your parser domain)
	// This is a placeholder - you'll need to implement actual expression parsing and evaluation
	// based on your parser module structure

	// Track positions
	var positions []*model.Position
	var openPosition *model.Position
	capital := strategy.Capital
	initialCapital := capital

	// Mark run as started
	startTime := time.Now()
	run.StartedTime = &startTime
	run.State = model.Started
	e.runRepo.Update(ctx, run)

	for i, instrument := range instruments {
		// If we have at least 20 candles for indicators, evaluate expressions
		if i >= 20 {
			// Create data context for expression evaluation
			dataContext := e.createDataContext(instruments, i)

			// Check entry condition if no position is open
			// This is simplified - you'd use your actual parser module here
			entrySignalTriggered := e.evaluateEntryCondition(strategy.EntryExpression, dataContext)

			if openPosition == nil && entrySignalTriggered {
				// Calculate position size based on strategy
				positionSize := e.calculatePositionSize(strategy, capital, instrument.Close)

				// Determine position direction from strategy
				var direction model.Direction
				if strategy.PositionType == "LONG" || strategy.PositionType == "long" {
					direction = model.Long
				} else {
					direction = model.Short
				}

				// Create new position
				openPosition = model.NewPosition(
					run.ID,
					direction,
					instrument.Close,
					instrument.Timestamp,
					strategy.EntryExpression,
					positionSize,
				)

				// Simulate order execution
				entryPrice := instrument.Close
				entryFees := calculateFees(entryPrice, positionSize)
				openPosition.Enter(entryPrice, instrument.Timestamp, entryFees, nil)

				// Update capital
				capital -= (entryPrice*positionSize + entryFees)
			}

			// Check exit condition if position is open
			exitSignalTriggered := false
			if strategy.ExitExpression != "" {
				exitSignalTriggered = e.evaluateExitCondition(strategy.ExitExpression, dataContext)
			}

			if openPosition != nil && (exitSignalTriggered ||
				e.checkStopLossOrTarget(openPosition, strategy, instrument) ||
				i == len(instruments)-1) { // Force exit on last candle

				// Calculate exit price and reason
				exitPrice := instrument.Close
				exitReason := "Signal"

				if exitSignalTriggered {
					exitReason = "Exit signal triggered"
				} else if i == len(instruments)-1 {
					exitReason = "End of backtest period"
				} else {
					// Check if it was stop loss or target profit
					if e.isStopLossTriggered(openPosition, strategy, instrument) {
						exitReason = "Stop loss"
					} else if e.isTargetProfitTriggered(openPosition, strategy, instrument) {
						exitReason = "Target profit"
					}
				}

				// Calculate fees
				exitFees := calculateFees(exitPrice, openPosition.PositionSize)

				// Close the position
				openPosition.Exit(
					instrument.Close,
					instrument.Timestamp,
					exitPrice,
					instrument.Timestamp,
					strategy.ExitExpression,
					exitReason,
					exitFees,
					nil,
				)

				// Save position
				positions = append(positions, openPosition)
				e.positionRepo.Create(ctx, openPosition)

				// Update capital if this is a real exit (not end of backtest)
				if openPosition.NetProfitLoss != nil {
					capital += (exitPrice*openPosition.PositionSize - exitFees)
					if *openPosition.NetProfitLoss > 0 {
						capital += *openPosition.NetProfitLoss
					} else {
						capital += *openPosition.NetProfitLoss // This is negative already
					}
				}

				// Reset open position
				openPosition = nil
			}
		}
	}

	// Mark run as ended
	endTime := time.Now()
	run.EndedTime = &endTime
	run.State = model.Ended
	e.runRepo.Update(ctx, run)

	// Calculate and save metrics
	metrics := e.calculateMetrics(positions, initialCapital, capital, run)
	e.runMetricsRepo.Create(ctx, metrics)
}

// createDataContext creates a data context for expression evaluation
func (e *BacktestEngine) createDataContext(instruments []Instrument, currentIndex int) map[string]interface{} {
	context := make(map[string]interface{})

	// Add basic OHLCV data
	context["open"] = instruments[currentIndex].Open
	context["high"] = instruments[currentIndex].High
	context["low"] = instruments[currentIndex].Low
	context["close"] = instruments[currentIndex].Close
	context["volume"] = instruments[currentIndex].Volume

	// Add indicator calculations (simplified)
	// In a real implementation, these would use your parser module

	// Calculate SMA(20)
	sma20 := 0.0
	for i := 0; i < 20; i++ {
		sma20 += instruments[currentIndex-i].Close
	}
	sma20 /= 20.0
	context["sma20"] = sma20

	// Calculate SMA(50)
	sma50 := 0.0
	if currentIndex >= 50 {
		for i := 0; i < 50; i++ {
			sma50 += instruments[currentIndex-i].Close
		}
		sma50 /= 50.0
	}
	context["sma50"] = sma50

	return context
}

// evaluateEntryCondition evaluates the entry condition expression
// This is a placeholder - you'd use your actual parser module
func (e *BacktestEngine) evaluateEntryCondition(expression string, context map[string]interface{}) bool {
	// Simple evaluation for testing
	// In a real implementation, this would use your parser domain
	if context["close"].(float64) > context["sma20"].(float64) &&
		context["sma20"].(float64) > context["sma50"].(float64) {
		return true
	}
	return false
}

// evaluateExitCondition evaluates the exit condition expression
// This is a placeholder - you'd use your actual parser module
func (e *BacktestEngine) evaluateExitCondition(expression string, context map[string]interface{}) bool {
	// Simple evaluation for testing
	// In a real implementation, this would use your parser domain
	return context["close"].(float64) < context["sma20"].(float64)
}

// checkStopLossOrTarget checks if a position should be closed due to stop loss or target price
func (e *BacktestEngine) checkStopLossOrTarget(position *model.Position, strategy *strategyDto.StrategyResponse, instrument Instrument) bool {
	return e.isStopLossTriggered(position, strategy, instrument) ||
		e.isTargetProfitTriggered(position, strategy, instrument)
}

// isStopLossTriggered checks if stop loss is triggered
func (e *BacktestEngine) isStopLossTriggered(position *model.Position, strategy *strategyDto.StrategyResponse, instrument Instrument) bool {
	if strategy.StopLoss <= 0 || position.EntryPriceOrder == nil {
		return false
	}

	entryPrice := *position.EntryPriceOrder
	stopPrice := 0.0

	if position.Direction == model.Long {
		if strategy.StopLossType == "percent" {
			stopPrice = entryPrice * (1 - strategy.StopLoss/100)
		} else {
			stopPrice = entryPrice - strategy.StopLoss
		}
		return instrument.Low <= stopPrice
	} else { // Short position
		if strategy.StopLossType == "percent" {
			stopPrice = entryPrice * (1 + strategy.StopLoss/100)
		} else {
			stopPrice = entryPrice + strategy.StopLoss
		}
		return instrument.High >= stopPrice
	}
}

// isTargetProfitTriggered checks if target profit is triggered
func (e *BacktestEngine) isTargetProfitTriggered(position *model.Position, strategy *strategyDto.StrategyResponse, instrument Instrument) bool {
	if strategy.TargetProfit <= 0 || position.EntryPriceOrder == nil {
		return false
	}

	entryPrice := *position.EntryPriceOrder
	targetPrice := 0.0

	if position.Direction == model.Long {
		if strategy.TargetProfitType == "percent" {
			targetPrice = entryPrice * (1 + strategy.TargetProfit/100)
		} else {
			targetPrice = entryPrice + strategy.TargetProfit
		}
		return instrument.High >= targetPrice
	} else { // Short position
		if strategy.TargetProfitType == "percent" {
			targetPrice = entryPrice * (1 - strategy.TargetProfit/100)
		} else {
			targetPrice = entryPrice - strategy.TargetProfit
		}
		return instrument.Low <= targetPrice
	}
}

// calculatePositionSize calculates the position size based on strategy settings
func (e *BacktestEngine) calculatePositionSize(strategy *strategyDto.StrategyResponse, capital float64, price float64) float64 {
	// If strategy specifies fixed quantity, use that
	if strategy.Qty > 0 {
		return float64(strategy.Qty)
	}

	// Otherwise, use a percentage of capital
	// Default to 10% if not specified
	capitalPercentage := 10.0
	capitalToRisk := capital * (capitalPercentage / 100)

	// Calculate position size based on available capital
	return capitalToRisk / price
}

// calculateFees calculates the trading fees
func calculateFees(price float64, quantity float64) float64 {
	// Implement your fee structure
	// For example, 0.1% of transaction value
	return price * quantity * 0.001
}

// calculateMetrics calculates backtest performance metrics
func (e *BacktestEngine) calculateMetrics(positions []*model.Position, initialCapital, finalCapital float64, run *model.Run) *model.RunMetrics {
	metrics := model.NewRunMetrics(run.ID)

	// Calculate basic metrics
	totalTrades := len(positions)
	winningTrades := 0
	losingTrades := 0
	breakEvenTrades := 0
	totalProfitLoss := 0.0
	totalWinAmount := 0.0
	totalLossAmount := 0.0
	totalCommission := 0.0
	totalSlippage := 0.0

	// Process all positions
	for _, position := range positions {
		if position.NetProfitLoss == nil {
			continue
		}

		profitLoss := *position.NetProfitLoss

		if profitLoss > 0 {
			winningTrades++
			totalWinAmount += profitLoss
		} else if profitLoss < 0 {
			losingTrades++
			totalLossAmount += profitLoss // This is negative
		} else {
			breakEvenTrades++
		}

		totalProfitLoss += profitLoss

		// Add transaction costs
		if position.EntryFees != nil {
			totalCommission += *position.EntryFees
		}
		if position.ExitFees != nil {
			totalCommission += *position.ExitFees
		}
	}

	// Set metrics
	metrics.TotalTrades = totalTrades

	if totalTrades > 0 {
		metrics.WinRate = float64(winningTrades) / float64(totalTrades) * 100
	}

	// Calculate profit factor
	profitFactor := 0.0
	if totalLossAmount != 0 {
		profitFactor = totalWinAmount / stdmath.Abs(totalLossAmount)
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

	// Set transaction costs
	metrics.TotalCommission = totalCommission
	metrics.TotalSlippage = totalSlippage

	if initialCapital > 0 {
		metrics.TransactionCostsPct = (totalCommission + totalSlippage) / initialCapital * 100
		metrics.TotalReturnPct = (finalCapital - initialCapital) / initialCapital * 100
	}

	// Calculate duration-based metrics
	// Assuming you have appropriate start/end times in each position
	totalDurationSeconds := 0.0
	count := 0

	for _, position := range positions {
		if position.EntryTimeOrder != nil && position.ExitTimeOrder != nil {
			durationSeconds := position.ExitTimeOrder.Sub(*position.EntryTimeOrder).Seconds()
			totalDurationSeconds += durationSeconds
			count++
		}
	}

	if count > 0 {
		avgTradeDuration := totalDurationSeconds / float64(count) / 3600 // Convert to hours
		metrics.AvgTradeDuration = &avgTradeDuration
	}

	// Calculate more advanced metrics (simplified)
	metrics.SharpeRatio = 0.0  // Would require daily returns calculation
	metrics.SortinoRatio = 0.0 // Would require daily returns calculation
	metrics.MaxDrawdown = 0.0  // Would require equity curve calculation

	// Store additional metrics as JSON
	additionalMetrics := map[string]interface{}{
		"initialCapital":  initialCapital,
		"finalCapital":    finalCapital,
		"netProfit":       finalCapital - initialCapital,
		"winningTrades":   winningTrades,
		"losingTrades":    losingTrades,
		"breakEvenTrades": breakEvenTrades,
	}

	metrics.SetMetricData(additionalMetrics)

	// Update metrics timestamp
	metrics.UpdatedAt = time.Now()

	return metrics
}
