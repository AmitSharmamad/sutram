package service

import (
	"context"
	"math"
	"sync"
	"time"

	"algo/internal/modules/backtest/domain/model"
	backtestRepo "algo/internal/modules/backtest/domain/repository"
	instrumentService "algo/internal/modules/instrument/application/service"
	instrumentModel "algo/internal/modules/instrument/domain/model"
	parser "algo/internal/modules/parser/domain/model"
	strategyDto "algo/internal/modules/strategy/application/dto"
	strategyService "algo/internal/modules/strategy/application/service"

	"github.com/google/uuid"
)

// EngineService defines the service for backtest engine operations
type EngineService interface {
	// Start a backtest run for the provided strategy
	Start(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) (*model.RunGroup, error)
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

// We're using instrumentModel.Instrument instead of defining our own struct

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

// Start starts a backtest run for the provided strategy
func (e *BacktestEngine) Start(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) (*model.RunGroup, error) {
	// Get the strategy
	strategyResp, err := e.strategyService.GetStrategyByID(ctx, strategyID, userID)
	if err != nil {
		return nil, err
	}

	// Create a run group in memory
	runGroup := model.NewRunGroup(strategyID)

	// Create runs for each instrument
	runs := make([]*model.Run, 0)
	var wg sync.WaitGroup

	for _, instrumentID := range strategyResp.Instruments {
		// Fetch instrument data using ID directly
		// TODO: Move this outside of the loop
		instrumentMaster, err := e.instrumentService.GetInstrumentMasterByID(
			ctx,
			uint64(instrumentID),
		)
		if err != nil {
			// Log error and continue with other instruments
			continue
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

		instrumentIDInt := int(instrumentMaster.ID)
		qty := 1 // Default quantity
		if strategyResp.Qty > 0 {
			qty = strategyResp.Qty
		}

		// Create a run for this instrument (in memory)
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

		runs = append(runs, run)

		wg.Add(1)

		// Execute run asynchronously as goroutine
		go e.executeRunBacktest(ctx, run, int64(instrumentID), strategyResp, &wg)
	}

	// Wait for all backtest runs to complete
	wg.Wait()

	// Save run group to database
	if err := e.runGroupRepo.Create(ctx, runGroup); err != nil {
		return nil, err
	}

	// Bulk create all runs
	if len(runs) > 0 {
		if err := e.runRepo.BulkCreate(ctx, runs); err != nil {
			return nil, err
		}
	}

	// Collect all positions and metrics from runs
	positions := make([]*model.Position, 0)
	metrics := make([]*model.RunMetrics, 0)

	for _, run := range runs {
		if run.RunMetrics != nil {
			metrics = append(metrics, run.RunMetrics)
		}
		positions = append(positions, run.Positions...)
	}

	// Bulk create all positions
	if len(positions) > 0 {
		if err := e.positionRepo.BulkCreate(ctx, positions); err != nil {
			return nil, err
		}
	}

	// Bulk create all metrics
	if len(metrics) > 0 {
		if err := e.runMetricsRepo.BulkCreate(ctx, metrics); err != nil {
			return nil, err
		}
	}

	return runGroup, nil
}

// executeRunBacktest runs a backtest for a single instrument run
func (e *BacktestEngine) executeRunBacktest(ctx context.Context, run *model.Run, instrumentID int64, strategy *strategyDto.StrategyResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	instrumentMaster, err := e.instrumentService.GetInstrumentMasterByID(
		ctx,
		uint64(instrumentID),
	)
	if err != nil {
		return
	}

	// Get data for the main chart timeframe
	timeFrameData, err := e.instrumentService.GetInstrumentsByTimeRange(
		ctx,
		instrumentMaster.ID,
		run.PeriodStartTime,
		run.PeriodEndTime,
		string(strategy.ChartTimeframe),
	)
	if err != nil {
		return
	}

	// Get second-level data for more accurate exit condition checking
	secondLevelData, err := e.instrumentService.GetInstrumentsByTimeRange(
		ctx,
		instrumentMaster.ID,
		run.PeriodStartTime,
		run.PeriodEndTime,
		"1s", // 1-second data for granular exit condition checking
	)
	if err != nil {
		// If we can't get second-level data, we'll just use the timeframe data
		// Continue with the backtest using only timeframe data
		secondLevelData = nil
	}

	// Execute backtest for this instrument and gather positions and metrics
	positions, metrics := e.executeBacktestInMemory(strategy, timeFrameData, secondLevelData, run)

	// Store the positions and metrics in the run
	run.Positions = positions
	run.RunMetrics = metrics
}

func (e *BacktestEngine) toDataContext(
	instruments []*instrumentModel.Instrument,
) *instrumentModel.StockData {
	// Convert the slice of instruments to a StockData type
	stockData := make(instrumentModel.StockData, len(instruments))
	for i, instrument := range instruments {
		stockData[i] = *instrument
	}
	return &stockData
}

// executeBacktestInMemory runs the backtest for a single instrument and returns the results without saving to database
func (e *BacktestEngine) executeBacktestInMemory(strategy *strategyDto.StrategyResponse, stockData []*instrumentModel.Instrument, secondLevelData []*instrumentModel.Instrument, run *model.Run) ([]*model.Position, *model.RunMetrics) {
	// Initialize variables to track positions and capital
	var positions []*model.Position
	var openPosition *model.Position
	capital := strategy.Capital
	initialCapital := capital

	// Mark run as started
	startTime := time.Now()
	run.StartedTime = &startTime
	run.State = model.Started

	// Minimum lookback period for indicators (for example, SMA(20) would need at least 20 candles)
	minLookback := 20

	// Process each candle from the minimum lookback period to the end of the data
	for i := minLookback; i < len(stockData); i++ {
		// Create data context using the data slice from the beginning up to the current position
		data := e.toDataContext(stockData[:i+1])

		// Check entry condition if no position is open
		if openPosition == nil {
			entrySignalTriggered := e.evaluateEntryCondition(strategy.EntryExpression, data)

			if entrySignalTriggered {
				// Calculate position size based on strategy
				positionSize := e.calculatePositionSize(strategy, capital, stockData[i].Close)

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
					stockData[i].Close,
					stockData[i].Timestamp,
					strategy.EntryExpression,
					positionSize,
				)

				// Simulate order execution
				entryPrice := stockData[i].Close
				entryFees := calculateFees(entryPrice, positionSize)
				openPosition.Enter(entryPrice, stockData[i].Timestamp, entryFees, nil)

				// Update capital for position entry
				capital = e.updateCapitalOnPositionChange(capital, openPosition, entryPrice, entryFees, true)

			}
		}

		// If position is still open, check exit condition using timeframe data
		if openPosition != nil {
				e.TryExit(openPosition, strategy, secondLevelData)

				// Collect position
				positions = append(positions, openPosition)

				// Update capital for position exit
				capital = e.updateCapitalOnPositionChange(
					capital, openPosition, *openPosition.ExitPriceOrder, *openPosition.ExitFees, false,
				)

				// Add net profit/loss if available
				if openPosition.NetProfitLoss != nil {
					capital += *openPosition.NetProfitLoss
				}

				// Reset open position
				openPosition = nil
			}
	}

	// Mark run as ended
	endTime := time.Now()
	run.EndedTime = &endTime
	run.State = model.Ended

	// Calculate metrics
	metrics := e.calculateMetrics(positions, initialCapital, capital, run)

	return positions, metrics
}

// executeBacktest runs the backtest for a single instrument and saves results to the database
// This legacy method is kept for backward compatibility
func (e *BacktestEngine) executeBacktest(ctx context.Context, strategy *strategyDto.StrategyResponse, instruments []*instrumentModel.Instrument, run *model.Run) {
	// Mark run as started
	startTime := time.Now()
	run.StartedTime = &startTime
	run.State = model.Started
	e.runRepo.Update(ctx, run)

	// Run the backtest in memory
	positions, metrics := e.executeBacktestInMemory(strategy, instruments, nil, run)

	// Save positions one by one
	for _, position := range positions {
		e.positionRepo.Create(ctx, position)
	}

	// Update run state
	e.runRepo.Update(ctx, run)

	// Save metrics
	if metrics != nil {
		e.runMetricsRepo.Create(ctx, metrics)
	}
}

// createDataContext creates a data context for expression evaluation
func (e *BacktestEngine) createDataContext(instruments []*instrumentModel.Instrument, currentIndex int) map[string]interface{} {
	context := make(map[string]any)
	context["data"] = instruments
	return context
}

func (e *BacktestEngine) evaluateEntryCondition(expression string, data *instrumentModel.StockData) bool {
	return parser.Execute(expression, parser.NewEvaluationContext(data))
}

func (e *BacktestEngine) evaluateExitCondition(expression string, data *instrumentModel.StockData) bool {
	return parser.Execute(expression, parser.NewEvaluationContext(data))
}

// checkStopLossOrTarget checks if a position should be closed due to stop loss or target price
func (e *BacktestEngine) checkStopLossOrTarget(position *model.Position, strategy *strategyDto.StrategyResponse, instrument *instrumentModel.Instrument) bool {
	return e.isStopLossTriggered(position, strategy, instrument) ||
		e.isTargetProfitTriggered(position, strategy, instrument)
}

// isStopLossTriggered checks if stop loss is triggered
func (e *BacktestEngine) isStopLossTriggered(position *model.Position, strategy *strategyDto.StrategyResponse, instrument *instrumentModel.Instrument) bool {
	if strategy.StopLoss <= 0 || position.EntryPriceOrder == nil {
		return false
	}

	entryPrice := *position.EntryPriceOrder
	stopPrice := 0.0

	if position.Direction == model.Long {
		if strategy.StopLossType == "PERCENTAGE" {
			stopPrice = entryPrice * (1 - strategy.StopLoss/100)
		} else {
			stopPrice = entryPrice - strategy.StopLoss
		}
		return instrument.Low <= stopPrice
	} else { // Short position
		if strategy.StopLossType == "PERCENTAGE" {
			stopPrice = entryPrice * (1 + strategy.StopLoss/100)
		} else {
			stopPrice = entryPrice + strategy.StopLoss
		}
		return instrument.High >= stopPrice
	}
}

// isTargetProfitTriggered checks if target profit is triggered
func (e *BacktestEngine) isTargetProfitTriggered(position *model.Position, strategy *strategyDto.StrategyResponse, instrument *instrumentModel.Instrument) bool {
	if strategy.TargetProfit <= 0 || position.EntryPriceOrder == nil {
		return false
	}

	entryPrice := *position.EntryPriceOrder
	targetPrice := 0.0

	if position.Direction == model.Long {
		if strategy.TargetProfitType == "PERCENTAGE" {
			targetPrice = entryPrice * (1 + strategy.TargetProfit/100)
		} else {
			targetPrice = entryPrice + strategy.TargetProfit
		}
		return instrument.High >= targetPrice
	} else { // Short position
		if strategy.TargetProfitType == "PERCENTAGE" {
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

// updateCapitalOnPositionChange updates the capital based on the position direction and event type (entry or exit)
func (e *BacktestEngine) updateCapitalOnPositionChange(capital float64, position *model.Position, price float64, fees float64, isEntry bool) float64 {
	if isEntry {
		// For position entry
		if position.Direction == model.Long {
			// For long positions, capital decreases (buying)
			capital -= (price*position.PositionSize + fees)
		} else {
			// For short positions, capital increases (selling)
			capital += (price*position.PositionSize - fees)
		}
	} else {
		// For position exit
		if position.Direction == model.Long {
			// For long positions, capital increases when exiting (selling)
			capital += (price*position.PositionSize - fees)
		} else {
			// For short positions, capital decreases when exiting (buying to cover)
			capital -= (price*position.PositionSize + fees)
		}
	}

	return capital
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

// processSecondLevelData processes second-level data for granular exit condition checking
func (e *BacktestEngine) processSecondLevelData(position *model.Position, strategy *strategyDto.StrategyResponse, secondLevelData []*instrumentModel.Instrument, currentTime, nextTime time.Time) bool {
	// Find the range of second-level data within the current timeframe window
	var relevantData []*instrumentModel.Instrument
	for _, data := range secondLevelData {
		if data.Timestamp.After(currentTime) && data.Timestamp.Before(nextTime) {
			relevantData = append(relevantData, data)
		}
	}

	if len(relevantData) == 0 {
		return false
	}

	// Process each second-level data point
	for i, data := range relevantData {
		// Create data context for evaluation using available second-level data
		// We use the toDataContext method to create a compatible StockData object
		dataContext := e.toDataContext(relevantData[:i+1])

		// Check exit conditions
		exitSignalTriggered := false
		if strategy.ExitExpression != "" {
			exitSignalTriggered = e.evaluateExitCondition(strategy.ExitExpression, dataContext)
		}

		// Check stop loss or target profit hit
		stopLossHit := e.isStopLossTriggered(position, strategy, data)
		targetProfitHit := e.isTargetProfitTriggered(position, strategy, data)

		if exitSignalTriggered || stopLossHit || targetProfitHit {
			// Determine exit reason
			exitReason := "Exit signal"
			if stopLossHit {
				exitReason = "Stop loss"
			} else if targetProfitHit {
				exitReason = "Target profit"
			}

			// Calculate fees
			exitFees := calculateFees(data.Close, position.PositionSize)

			// Close the position
			position.Exit(
				data.Close,
				data.Timestamp,
				data.Close,
				data.Timestamp,
				strategy.ExitExpression,
				exitReason,
				exitFees,
				nil,
			)

			return true
		}
	}

	return false
}

// createSecondLevelContext creates a data context for second-level data evaluation
func (e *BacktestEngine) createSecondLevelContext(instruments []*instrumentModel.Instrument) map[string]interface{} {
	if len(instruments) == 0 {
		return make(map[string]interface{})
	}

	context := make(map[string]interface{})
	currentIndex := len(instruments) - 1

	// Add basic OHLCV data from current tick
	context["open"] = instruments[currentIndex].Open
	context["high"] = instruments[currentIndex].High
	context["low"] = instruments[currentIndex].Low
	context["close"] = instruments[currentIndex].Close
	context["volume"] = instruments[currentIndex].Volume

	// Add recent price trends for analysis
	// Average price of the last n ticks
	lastNPrices := min(10, len(instruments))
	if lastNPrices > 0 {
		avgPrice := 0.0
		for i := 0; i < lastNPrices; i++ {
			avgPrice += instruments[currentIndex-i].Close
		}
		avgPrice /= float64(lastNPrices)
		context["avgPrice"] = avgPrice
	}

	// Calculate simple volatility metric (standard deviation of recent prices)
	if len(instruments) > 5 {
		mean := 0.0
		for i := 0; i < 5; i++ {
			mean += instruments[currentIndex-i].Close
		}
		mean /= 5.0

		variance := 0.0
		for i := 0; i < 5; i++ {
			diff := instruments[currentIndex-i].Close - mean
			variance += diff * diff
		}
		variance /= 5.0
		context["recentVolatility"] = math.Sqrt(variance)
	}

	// Add price change rate for micro-trend detection
	if len(instruments) > 1 {
		priceChangePct := (instruments[currentIndex].Close - instruments[currentIndex-1].Close) /
			instruments[currentIndex-1].Close * 100
		context["priceChangePct"] = priceChangePct
	}

	return context
}

func (e *BacktestEngine) TryExit(
	position *model.Position,
	strategy *strategyDto.StrategyResponse,
	secondLevelData []*instrumentModel.Instrument,
) bool {
	// If the position is already closed, return
	if position.State == model.Exited {
		return false
	}

	// Cut the data right after the entry time till the end
	var relevantData []*instrumentModel.Instrument
	if position.EntryTimeOrder == nil {
		return false
	}

	for _, data := range secondLevelData {
		if data.Timestamp.After(*position.EntryTimeOrder) {
			relevantData = append(relevantData, data)
		}
	}

	if len(relevantData) == 0 {
		return false
	}

	// Iterate through the second-level data to check for exit conditions
	for i, data := range relevantData {
		// Create data context for evaluation using available second-level data
		dataContext := e.toDataContext(relevantData[:i+1])

		// Check exit conditions
		exitSignalTriggered := false
		if strategy.ExitExpression != "" {
			exitSignalTriggered = e.evaluateExitCondition(strategy.ExitExpression, dataContext)
		}

		// Check stop loss or target profit hit
		stopLossHit := e.isStopLossTriggered(position, strategy, data)
		targetProfitHit := e.isTargetProfitTriggered(position, strategy, data)

		if exitSignalTriggered || stopLossHit || targetProfitHit {
			// Determine exit reason
			exitReason := "Exit signal"
			if stopLossHit {
				exitReason = "Stop loss"
			} else if targetProfitHit {
				exitReason = "Target profit"
			}

			// Calculate fees
			exitFees := calculateFees(data.Close, position.PositionSize)

			// Calculate exit order price with slippage
			exitOrderPrice := calculateExitOrderSlippage(
				data.Close,
				position.Direction,
				data,
			)

			// Close the position
			position.Exit(
				data.Close,
				data.Timestamp,
				exitOrderPrice,
				data.Timestamp,
				strategy.ExitExpression,
				exitReason,
				exitFees,
				nil,
			)

			return true
		}
	}

	// If no exit condition was met, return false
	return false
}

// calculateExitOrderSlippage simulates slippage between exit signal and exit order execution
func calculateExitOrderSlippage(
	signalPrice float64,
	direction model.Direction,
	instrument *instrumentModel.Instrument,
) float64 {
	// Default slippage percentage (0.1% of price)
	slippagePct := 0.1

	// Calculate base slippage amount
	baseSlippage := signalPrice * (slippagePct / 100.0)

	// Adjust slippage based on volatility (using high-low range as a simple volatility proxy)
	volatilityFactor := (instrument.High - instrument.Low) / instrument.Low
	// Scale volatility factor to not be too extreme (cap at 2x default slippage)
	adjustedVolatilityFactor := math.Min(volatilityFactor*10, 2.0)

	// Calculate final slippage amount accounting for volatility
	slippageAmount := baseSlippage * adjustedVolatilityFactor

	// For long positions, exit price is typically worse (lower) due to slippage
	// For short positions, exit price is typically worse (higher) due to slippage
	if direction == model.Long {
		return signalPrice - slippageAmount
	} else {
		return signalPrice + slippageAmount
	}
}
