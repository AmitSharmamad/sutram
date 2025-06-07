package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"algo/internal/modules/backtest/domain/model"
	backtestRepo "algo/internal/modules/backtest/domain/repository"
	instrumentService "algo/internal/modules/instrument/application/service"
	instrumentModel "algo/internal/modules/instrument/domain/model"
	parser "algo/internal/modules/parser/domain/model"
	strategyDto "algo/internal/modules/strategy/application/dto"
	strategyService "algo/internal/modules/strategy/application/service"
	strategyModel "algo/internal/modules/strategy/domain/model"

	"github.com/google/uuid"
)

// Response struct represents the response from a backtest run
type Response struct {
	InstrumentMaster *instrumentModel.InstrumentMaster `json:"instrument_master"`
	InstrumentID     string                            `json:"instrument_id"`
	Positions        []*model.Position                 `json:"positions"`
	Metrics          *model.RunMetrics                 `json:"metrics"`
}

// IndicatorValue represents a single indicator value with timestamp
type IndicatorValue struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// EngineService defines the service for backtest engine operations
type EngineService interface {
	// Start a backtest run for the provided strategy
	Start(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]Response, error)
	// GetIndicatorData retrieves calculated indicator data for a strategy
	GetIndicatorData(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) (map[string]map[string][]IndicatorValue, error)
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
	metricsService    MetricsCalculationService
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
		metricsService:    NewMetricsCalculationService(),
	}
}

// Start starts a backtest run for the provided strategy
func (e *BacktestEngine) Start(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]Response, error) {
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
		// Fetch instrument data using exchange security ID directly
		// TODO: Move this outside of the loop
		instrumentMaster, err := e.instrumentService.GetInstrumentMasterByExchangeSecurityID(
			ctx,
			int64(instrumentID),
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

		instrumentIDInt := int(instrumentMaster.ExchangeSecurityID)
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
			instrumentMaster,
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

	// Create the response
	var response []Response

	// Create a response for each run
	for _, run := range runs {
		// Convert the instrument ID to a string for the response
		instrumentIDStr := strconv.Itoa(run.InstrumentID)

		resp := Response{
			InstrumentMaster: run.InstrumentMaster,
			InstrumentID:     instrumentIDStr,
			Positions:        run.Positions,
			Metrics:          run.RunMetrics,
		}

		response = append(response, resp)
	}

	return response, nil
}

func (e *BacktestEngine) getLookBackFromExpressions(strategy *strategyDto.StrategyResponse) int64 {

	entryLookBack := e.getMaxLookback(strategy.EntryExpression)

	exitLookBack := int64(0)
	if strategy.ExitExpression != "" {
		exitLookBack = e.getMaxLookback(strategy.ExitExpression)
	}

	lookBack := max(entryLookBack, exitLookBack)

	return lookBack
}

// getMaxLookback extracts indicators from the expression and returns the maximum lookback period required
func (e *BacktestEngine) getMaxLookback(expression string) int64 {
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

// executeRunBacktest runs a backtest for a single instrument run
func (e *BacktestEngine) executeRunBacktest(ctx context.Context, run *model.Run, instrumentID int64, strategy *strategyDto.StrategyResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()

	instrumentMaster, err := e.instrumentService.GetInstrumentMasterByExchangeSecurityID(
		ctx,
		int64(instrumentID),
	)
	if err != nil {
		println("Error getting instrument master for exchange security ID:", instrumentID, "Error:", err.Error())
		return
	}

	if instrumentMaster == nil {
		println("Instrument master is nil for exchange security ID:", instrumentID)
		return
	}

	// Determine the fetch start time based on entry_start_time if available

	// Calculate lookback from expressions
	lookbackFromExpressions := e.getLookBackFromExpressions(strategy)

	// Get historical data with exact lookback count
	timeFrameDataStartTime := time.Now()
	println("\tFetching timeframe data with lookback for instrument", instrumentID, "- Start:", timeFrameDataStartTime.Format(time.RFC3339))
	println("\t\tLookback period:", lookbackFromExpressions, "candles")
	println("\t\tFetch range: Before", strategy.StartDate.Format(time.RFC3339), "to", run.PeriodEndTime.Format(time.RFC3339))

	// Get lookback data (historical data before strategy start date)
	lookbackData, err := e.instrumentService.GetInstrumentsByExchangeSecurityIDWithLookback(
		ctx,
		instrumentMaster.ExchangeSecurityID,
		*strategy.StartDate,
		string(strategy.ChartTimeframe),
		int(lookbackFromExpressions),
	)
	if err != nil {
		println("Error getting lookback data for instrument ID:", instrumentID, "Error:", err.Error())
		lookbackData = []*instrumentModel.Instrument{} // Continue with empty lookback if error
	}

	// Get main data for the strategy period
	timeFrameData, err := e.instrumentService.GetInstrumentsByTimeRange(
		ctx,
		instrumentMaster.ExchangeSecurityID,
		*strategy.StartDate,
		run.PeriodEndTime,
		string(strategy.ChartTimeframe),
	)

	if err != nil {
		println("Error getting timeframe data for instrument ID:", instrumentID, "Error:", err.Error())
		return
	}

	// Combine lookback data with main data (lookback data first, then main data)
	// Reverse lookback data since it comes in descending order (newest first)
	for i, j := 0, len(lookbackData)-1; i < j; i, j = i+1, j-1 {
		lookbackData[i], lookbackData[j] = lookbackData[j], lookbackData[i]
	}

	// Combine lookback and main data
	timeFrameDataWithLookback := make([]*instrumentModel.Instrument, 0, len(lookbackData)+len(timeFrameData))
	timeFrameDataWithLookback = append(timeFrameDataWithLookback, lookbackData...)
	timeFrameDataWithLookback = append(timeFrameDataWithLookback, timeFrameData...)

	timeFrameDataEndTime := time.Now()
	timeFrameDataDuration := timeFrameDataEndTime.Sub(timeFrameDataStartTime)
	println("\tTimeframe data with lookback fetched for instrument", instrumentID, "- Duration:", timeFrameDataDuration.String(),
		"- Records: Lookback:", len(lookbackData), "Main:", len(timeFrameData), "Total:", len(timeFrameDataWithLookback))

	// Use the data with lookback as our timeFrameData
	timeFrameData = timeFrameDataWithLookback

	// Get second-level data for more accurate exit condition checking
	secondLevelDataStartTime := time.Now()
	println("\tFetching second-level data for instrument", instrumentID, "- Start:", secondLevelDataStartTime.Format(time.RFC3339))

	// secondLevelData, err := e.instrumentService.GetInstrumentsByTimeRange(
	// 	ctx,
	// 	instrumentMaster.ExchangeSecurityID,
	// 	fetchStartTime,
	// 	run.PeriodEndTime,
	// 	string(strategy.ChartTimeframe), // 1-second data for granular exit condition checking
	// )

	// secondLevelDataEndTime := time.Now()
	// secondLevelDataDuration := secondLevelDataEndTime.Sub(secondLevelDataStartTime)
	// secondLevelDataCount := 0
	// if secondLevelData != nil {
	// 	secondLevelDataCount = len(secondLevelData)
	// }
	// println("\tSecond-level data fetched for instrument", instrumentID, "- Duration:", secondLevelDataDuration.String(),
	// 	"- Records:", secondLevelDataCount)

	// if err != nil {
	// 	// If we can't get second-level data, we'll just use the timeframe data
	// 	// Continue with the backtest using only timeframe data
	// 	println("\tWarning: Could not get second-level data for instrument ID:", instrumentID, "Error:", err.Error())
	// 	secondLevelData = nil
	// }

	// Execute backtest for this instrument and gather positions and metrics
	positions, metrics := e.executeBacktestInMemory(strategy, timeFrameData, timeFrameData, run)

	// Store the positions and metrics in the run
	run.Positions = positions
	run.RunMetrics = metrics

	endTime := time.Now()
	executionDuration := endTime.Sub(startTime)

	// Format: InstrumentID (Symbol): Start -> End (Duration)
	instrumentIDStr := strconv.FormatInt(instrumentID, 10)
	symbolStr := instrumentMaster.Symbol
	if symbolStr == "" {
		symbolStr = "Unknown"
	}

	println("Backtest execution completed for instrument", instrumentIDStr, "("+symbolStr+"):",
		"\n\tStart:", startTime.Format(time.RFC3339),
		"\n\tEnd:", endTime.Format(time.RFC3339),
		"\n\tDuration:", executionDuration.String())
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
	// Start tracking execution time for this function
	functionStartTime := time.Now()

	// Initialize variables to track positions and capital
	var positions []*model.Position
	var openPosition *model.Position
	capital := strategy.Capital
	initialCapital := 0.0

	// Mark run as started
	startTime := time.Now()
	if run != nil {
		run.StartedTime = &startTime
		run.State = model.Started
	}

	// Keep track of the current trading day for MIS trades
	// Initialize current trading day using first timestamp from indicators
	var currentTradingDay time.Time

	evalContext := parser.NewEvaluationContext(nil)
	evalContext.Strategy = strategy

	// ********* Remove this piece of code *********
	evalContext.Data = e.toDataContext(stockData)

	// Extract indicators from both entry and exit expressions
	entryIndicators, err := parser.ExtractIndicators(strategy.EntryExpression)
	if err != nil {
		println("Error extracting indicators from entry expression:", err.Error())
	}

	var exitIndicators []parser.IndicatorInfo
	if strategy.ExitExpression != "" {
		exitIndicators, err = parser.ExtractIndicators(strategy.ExitExpression)
		if err != nil {
			println("Error extracting indicators from exit expression:", err.Error())
		}
	}

	// Combine all indicators (entry + exit)
	allIndicators := append(entryIndicators, exitIndicators...)

	indicatorArrayMap := make(map[string][]parser.IndicatorValues)
	length := 0

	// Iterate through all indicators and calculate their values
	for _, indicator := range allIndicators {
		calculator := parser.GetIndFactory(indicator.Name)
		uniqueString, err := calculator.GenerateIndicatorKey(indicator.Arguments, evalContext)
		if err != nil {
			println("Error generating unique key for indicator:", err.Error())
		}

		// Skip if already calculated (avoid duplicates between entry and exit)
		if _, exists := indicatorArrayMap[uniqueString]; exists {
			continue
		}

		indicatorArrayMap[uniqueString] = calculator.CalculateAll(indicator.Arguments, evalContext, *strategy.StartDate)
		if len(indicatorArrayMap[uniqueString]) > length {
			length = len(indicatorArrayMap[uniqueString])
		}
	}

	// Process each candle using pre-calculated indicators
	// Get timestamps from indicator arrays to iterate over
	var timestamps []time.Time
	timestampSet := make(map[time.Time]bool)

	// Collect all unique timestamps from indicator arrays
	for _, values := range indicatorArrayMap {
		for _, val := range values {
			if !timestampSet[val.Timestamp] {
				timestamps = append(timestamps, val.Timestamp)
				timestampSet[val.Timestamp] = true
			}
		}
	}

	// Sort timestamps chronologically
	for i := range len(timestamps) - 1 {
		for j := i + 1; j < len(timestamps); j++ {
			if timestamps[i].After(timestamps[j]) {
				timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
			}
		}
	}

	// Initialize current trading day using first timestamp
	if len(timestamps) > 0 {
		currentTradingDay = timestamps[0].Truncate(24 * time.Hour)
	}

	// Use timestamp array length for the main loop
	for i := range len(timestamps) {
		// Get current timestamp for this iteration
		currentTimestamp := timestamps[i]

		// Check if we've moved to a new trading day
		dateTruncated := currentTimestamp.Truncate(24 * time.Hour)
		if !dateTruncated.Equal(currentTradingDay) {
			// If we have an open position and this is an MIS (intraday) trade, close it at the end of the previous day
			if openPosition != nil && openPosition.State == model.Entered && strategy.TradeType == strategyModel.TradeType("MIS") {
				// Find the last timestamp from the previous day
				var lastTimestampFromPreviousDay *time.Time

				// Look for the most recent timestamp from the previous trading day
				for j := i - 1; j >= 0 && j >= max(0, i-100); j-- { // Limit search to avoid excessive iteration
					checkTimestamp := timestamps[j]
					if checkTimestamp.Truncate(24 * time.Hour).Equal(currentTradingDay) {
						lastTimestampFromPreviousDay = &checkTimestamp
						break
					}
				}

				if lastTimestampFromPreviousDay != nil {
					// Use the consolidated MIS square-off function
					if e.handleMISSquareOff(openPosition, strategy, indicatorArrayMap, *lastTimestampFromPreviousDay, &positions, &capital) {
						openPosition = nil
					}
				}
			}
			// Update the current trading day
			currentTradingDay = dateTruncated
		}

		// Check entry condition if no position is open
		if openPosition == nil {
			newPosition, updatedCapital := e.checkAndCreateNewPosition(
				strategy,
				indicatorArrayMap,
				currentTimestamp,
				i,
				run,
				capital,
			)

			if newPosition != nil {
				openPosition = newPosition
				capital = updatedCapital
			}
		}

		// If position is still open, check exit condition
		if openPosition != nil {
			// Process exit conditions and check if position was closed
			if e.processExitConditions(openPosition, strategy, indicatorArrayMap, currentTimestamp, i, &positions, &capital) {
				// Position was closed, reset open position and continue to next iteration
				openPosition = nil
				continue
			}

			// Close the position if MIS end-of-day is hit
			if strategy.TradeType == strategyModel.TradeType("MIS") && e.IsMISEndOfDayHit(currentTimestamp) {
				// Use the consolidated MIS square-off function
				if e.handleMISSquareOff(openPosition, strategy, indicatorArrayMap, currentTimestamp, &positions, &capital) {
					openPosition = nil
				}
			}
		}
	}

	// Handle any remaining open position at the end of backtest period
	if openPosition != nil && openPosition.State == model.Entered {
		// If this is a MIS position, close it at the end of the backtest
		if strategy.TradeType == strategyModel.TradeType("MIS") && len(timestamps) > 0 {
			// Get the last timestamp and close price from indicators
			lastTimestamp := timestamps[len(timestamps)-1]
			var lastClosePrice float64

			// Get close price from CLOSE indicator for the last timestamp
			if closeIndicator, exists := indicatorArrayMap["CLOSE"]; exists {
				for _, val := range closeIndicator {
					if val.Timestamp.Equal(lastTimestamp) {
						lastClosePrice = val.Value
						break
					}
				}
			}

			// Calculate fees
			exitFees := 0.0
			exitOrderPrice := lastClosePrice // Use close price from indicator

			// Close the position with "EOD MIS Square-off" as the reason
			exitReason := string(model.ExitReasonMISSquareOff) // Special reason for MIS end-of-day square-off

			openPosition.Exit(
				lastClosePrice,
				lastTimestamp,
				exitOrderPrice,
				lastTimestamp,
				"",         // No exit expression for auto square-offs
				exitReason, // Use our defined exit reason
				exitFees,
				nil,
			)

			// Collect position
			positions = append(positions, openPosition)

			// Add net profit/loss if available
			if openPosition.NetProfitLoss != nil {
				capital += *openPosition.NetProfitLoss
			}
		} else {
			// For CNC/NRML positions, just add to the list without closing
			positions = append(positions, openPosition)
		}

		// Reset open position in all cases
		openPosition = nil
	}

	// Mark run as ended
	endTime := time.Now()
	if run != nil {
		run.EndedTime = &endTime
		run.State = model.Ended
	}

	// Remove the positions that are not closed
	closedPositions := make([]*model.Position, 0)
	for _, position := range positions {
		if position.State == model.Exited {
			closedPositions = append(closedPositions, position)
		}
	}
	positions = closedPositions

	// Calculate metrics
	metrics := e.calculateMetrics(positions, initialCapital, capital, run)

	functionEndTime := time.Now()
	executionDuration := functionEndTime.Sub(functionStartTime)

	// Format a descriptive message about the backtest execution
	totalPositions := len(positions)
	instrumentInfo := ""
	if run != nil && run.InstrumentID > 0 {
		instrumentInfo = strconv.Itoa(run.InstrumentID)
	}

	println("\tBacktest in-memory calculation completed for instrument", instrumentInfo, ":",
		totalPositions, "positions generated,",
		"\n\t\tStart:", functionStartTime.Format(time.RFC3339),
		"\n\t\tEnd:", functionEndTime.Format(time.RFC3339),
		"\n\t\tDuration:", executionDuration.String())

	return positions, metrics
}

func (e *BacktestEngine) IsMISEndOfDayHit(
	timestamp time.Time,
) bool {
	// NOTE: MIS close is 15 minutes prior to market close
	exchange := "NSE"
	endHour := 15
	endMinute := 15

	if exchange == "MCX" {
		// MCX has different timings, so we need to adjust accordingly
		// MCX trading ends at 11:30 PM IST
		endHour = 23
		endMinute = 30
	}

	// Extract the time from the provided timestamp
	currentTime := timestamp

	// Create a time representing 3:15 PM on the same date
	eodTime := time.Date(
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		endHour,   // 3 PM (24-hour format)
		endMinute, // 15 minutes
		0,         // 0 seconds
		0,         // 0 nanoseconds
		currentTime.Location(),
	)

	// Return true if current time is at or after 3:15 PM
	return currentTime.Equal(eodTime) || currentTime.After(eodTime)
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

func (e *BacktestEngine) evaluateEntryCondition(expression string, evalContext *parser.EvaluationContext) bool {
	return parser.Execute(expression, evalContext)
}

func (e *BacktestEngine) evaluateExitCondition(expression string, evalContext *parser.EvaluationContext) bool {
	return parser.Execute(expression, evalContext)
}

func (e *BacktestEngine) calculateStopLoss(
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

func (e *BacktestEngine) calculateTargetProfit(
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

func (e *BacktestEngine) updateStopLossBasedOnTrailing(
	position *model.Position,
	strategy *strategyDto.StrategyResponse,
	instrument *instrumentModel.Instrument,
) bool {
	// Skip if trailing stop loss isn't enabled or position doesn't have a stop loss
	if strategy.TrailingStopLoss <= 0 || position.StopLossPrice == nil || position.EntryPriceOrder == nil {
		return false
	}

	currentPrice := instrument.Close
	entryPrice := *position.EntryPriceOrder
	currentStopLoss := *position.StopLossPrice

	initialStopDistance := *e.calculateStopLoss(entryPrice, strategy.StopLoss, strategy.StopLossType, strategy.PositionType)

	// Calculate the trailing amount based on type (percentage or absolute)
	trailingAmount := 0.0
	if strategy.TrailingStopLossType == "PERCENTAGE" {
		// For percentage, calculate the equivalent price distance
		if position.Direction == model.Long {
			trailingAmount = entryPrice * (strategy.TrailingStopLoss / 100.0)
		} else {
			trailingAmount = entryPrice * (strategy.TrailingStopLoss / 100.0)
		}
	} else {
		// For absolute, use the value directly
		trailingAmount = strategy.TrailingStopLoss
	}

	// Handle long position trailing stop
	if position.Direction == model.Long {
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
				*position.StopLossPrice = newStopLoss
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
				*position.StopLossPrice = newStopLoss
				return true
			}
		}
	}
	return false
}

// isStopLossTriggered checks if stop loss is triggered
func (e *BacktestEngine) isStopLossTriggered(
	position *model.Position,
	instrument *instrumentModel.Instrument,
) bool {
	if position.EntryPriceOrder == nil {
		return false
	}

	if position.Direction == model.Long {
		return instrument.Close <= *position.StopLossPrice
	} else { // Short position
		return instrument.Close >= *position.StopLossPrice
	}
}

// isTargetProfitTriggered checks if target profit is triggered
func (e *BacktestEngine) isTargetProfitTriggered(position *model.Position, strategy *strategyDto.StrategyResponse, instrument *instrumentModel.Instrument) bool {
	if strategy.TargetProfit <= 0 || position.EntryPriceOrder == nil {
		return false
	}

	if position.Direction == model.Long {
		return instrument.Close >= *position.TargetProfitPrice
	} else { // Short position
		return instrument.Close <= *position.TargetProfitPrice
	}
}

// calculatePositionSize calculates the position size based on strategy settings
func (e *BacktestEngine) calculatePositionSize(strategy *strategyDto.StrategyResponse, price float64) float64 {
	// 1. If strategy.qty <= 0 and strategy.max_allocation <= 0 throw an error
	if strategy.Qty <= 0 && strategy.MaxAllocation <= 0 {
		panic("Invalid position size parameters")
	}

	// 2. If strategy.qty > 0 and strategy.max_allocation <= 0 then return strategy.qty
	if strategy.Qty > 0 && strategy.MaxAllocation <= 0 {
		return float64(strategy.Qty)
	}

	// 3. If strategy.qty <= 0 and strategy.max_allocation > 0 then return strategy.max_allocation / price as integer
	// Ensure at least 1 share is purchased if the calculation results in less than 0
	if strategy.Qty <= 0 && strategy.MaxAllocation > 0 {
		quantity := math.Floor(strategy.MaxAllocation / price)
		if quantity <= 0 {
			return 1
		}
		return quantity
	}

	// 4. If strategy.qty > 0 and strategy.max_allocation > 0,
	// then return the lowest between (strategy.qty, strategy.max_allocation / price) as integer
	qtyValue := float64(strategy.Qty) * price
	if qtyValue < strategy.MaxAllocation {
		return float64(strategy.Qty)
	} else {
		quantity := math.Floor(strategy.MaxAllocation / price)
		if quantity <= 0 {
			return 1
		}
		return quantity
	}
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

// calculateMetrics calculates backtest performance metrics using the metrics calculation service
func (e *BacktestEngine) calculateMetrics(positions []*model.Position, initialCapital, finalCapital float64, run *model.Run) *model.RunMetrics {
	// Delegate all metrics calculation to the dedicated service
	return e.metricsService.CalculateRunMetrics(positions, initialCapital, finalCapital, run)
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
	evalContext *parser.EvaluationContext,
) bool {
	// If the position is already closed, return
	if position.State == model.Exited {
		return true // Return true since it is already exited
	}

	// No second-level data, can't evaluate exit
	if len(secondLevelData) == 0 {
		return false
	}

	// Position must have entry time
	if position.EntryTimeOrder == nil {
		return false
	}

	// Cut the data right after the entry time till the end
	var relevantData []*instrumentModel.Instrument

	// Limit the number of data points to check to avoid excessive processing
	const maxDataPoints = 1000
	pointsAdded := 0

	for _, data := range secondLevelData {
		if data.Timestamp.After(*position.EntryTimeOrder) {
			relevantData = append(relevantData, data)
			pointsAdded++

			// Limit the amount of data to process to avoid excessive computation
			if pointsAdded >= maxDataPoints {
				break
			}
		}
	}

	if len(relevantData) == 0 {
		return false
	}

	// Iterate through the second-level data to check for exit conditions
	for i, data := range relevantData {
		// Create data context for evaluation using available second-level data
		dataContext := e.toDataContext(relevantData[:i+1])
		evalContext.Data = dataContext

		// Check exit conditions
		exitSignalTriggered := false
		if strategy.ExitExpression != "" {
			exitSignalTriggered = e.evaluateExitCondition(strategy.ExitExpression, evalContext)
		}

		// Update the Stop Loss based on trailing stop loss
		e.updateStopLossBasedOnTrailing(
			position,
			strategy,
			data,
		)

		// Check stop loss or target profit hit
		stopLossHit := e.isStopLossTriggered(position, data)
		targetProfitHit := e.isTargetProfitTriggered(position, strategy, data)

		if exitSignalTriggered || stopLossHit || targetProfitHit {
			// Determine exit reason
			exitReason := string(model.ExitReasonSignal)
			if stopLossHit {
				// Check if this is a trailing stop loss exit
				if position.InitialStopLossPrice != nil &&
					position.StopLossPrice != nil &&
					*position.StopLossPrice > *position.InitialStopLossPrice {
					exitReason = string(model.ExitReasonTrailingStopLoss)
				} else {
					exitReason = string(model.ExitReasonStopLoss)
				}
			} else if targetProfitHit {
				exitReason = string(model.ExitReasonTargetProfit)
			}

			// Calculate fees
			exitFees := 0.0

			exitOrderPrice := data.Close // For backtest, we can use the close price directly

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

// findNextCandleIndex finds the index of the next candle matching the chart timeframe after the given timestamp
func (e *BacktestEngine) findNextCandleIndex(stockData []*instrumentModel.Instrument, exitTime time.Time, chartTimeframe string) int {
	// Safety check - empty data or nil
	if len(stockData) == 0 {
		return 0
	}

	// Find the next candle in the chart timeframe after the exit time
	for i, candle := range stockData {
		if candle.Timestamp.After(exitTime) {
			// Round up to the next available candle based on the chart timeframe
			// For example, if exit happened at 11:08 and timeframe is 5m,
			// the next entry should be checked at 11:10 (next 5m candle)
			return i
		}
	}

	// If no next candle is found, return the last valid index (not beyond the array)
	return len(stockData) - 1
}

// GetIndicatorData retrieves calculated indicator data for a strategy
func (e *BacktestEngine) GetIndicatorData(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) (map[string]map[string][]IndicatorValue, error) {
	// Get strategy details
	strategy, err := e.strategyService.GetStrategyByID(ctx, strategyID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get strategy: %w", err)
	}

	result := make(map[string]map[string][]IndicatorValue)

	// Process each instrument in the strategy
	for _, instrumentID := range strategy.Instruments {
		instrumentIDStr := strconv.FormatInt(int64(instrumentID), 10)

		// Get instrument master for metadata
		instrumentMaster, err := e.instrumentService.GetInstrumentMasterByExchangeSecurityID(
			ctx,
			int64(instrumentID),
		)
		if err != nil {
			continue // Skip this instrument if we can't get master data
		}

		// Get historical data with lookback for proper indicator calculation
		lookbackFromExpressions := e.getLookBackFromExpressions(strategy)

		// Get lookback data (historical data before strategy start date)
		lookbackData, err := e.instrumentService.GetInstrumentsByExchangeSecurityIDWithLookback(
			ctx,
			instrumentMaster.ExchangeSecurityID,
			*strategy.StartDate,
			string(strategy.ChartTimeframe),
			int(lookbackFromExpressions),
		)
		if err != nil {
			// Continue with empty lookback if error
			lookbackData = []*instrumentModel.Instrument{}
		}

		// Get main data for the strategy period
		timeFrameData, err := e.instrumentService.GetInstrumentsByTimeRange(
			ctx,
			instrumentMaster.ExchangeSecurityID,
			*strategy.StartDate,
			*strategy.EndDate,
			string(strategy.ChartTimeframe),
		)
		if err != nil {
			continue // Skip this instrument if data fetch fails
		}

		// Combine lookback data with main data
		// Reverse lookback data since it comes in descending order (newest first)
		for i, j := 0, len(lookbackData)-1; i < j; i, j = i+1, j-1 {
			lookbackData[i], lookbackData[j] = lookbackData[j], lookbackData[i]
		}

		// Combine lookback and main data
		allData := make([]*instrumentModel.Instrument, 0, len(lookbackData)+len(timeFrameData))
		allData = append(allData, lookbackData...)
		allData = append(allData, timeFrameData...)

		// Set up evaluation context
		evalContext := parser.NewEvaluationContext(nil)
		evalContext.Strategy = strategy
		evalContext.Data = e.toDataContext(allData)

		// Extract indicators from both entry and exit expressions
		entryIndicators, err := parser.ExtractIndicators(strategy.EntryExpression)
		if err != nil {
			continue
		}

		var exitIndicators []parser.IndicatorInfo
		if strategy.ExitExpression != "" {
			exitIndicators, err = parser.ExtractIndicators(strategy.ExitExpression)
			if err != nil {
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

		// Combine all indicators
		allIndicators := append(basicIndicators, entryIndicators...)
		allIndicators = append(allIndicators, exitIndicators...)

		// Calculate indicators
		indicatorMap := make(map[string][]IndicatorValue)
		calculatedKeys := make(map[string]bool) // Track calculated indicators to avoid duplicates

		for _, indicator := range allIndicators {
			calculator := parser.GetIndFactory(indicator.Name)
			if calculator == nil {
				continue
			}

			uniqueKey, err := calculator.GenerateIndicatorKey(indicator.Arguments, evalContext)
			if err != nil {
				continue
			}

			// Skip if already calculated
			if calculatedKeys[uniqueKey] {
				continue
			}
			calculatedKeys[uniqueKey] = true

			// Calculate indicator values
			indicatorValues := calculator.CalculateAll(indicator.Arguments, evalContext, *strategy.StartDate)

			// Convert to IndicatorValue format
			var convertedValues []IndicatorValue
			for _, val := range indicatorValues {
				convertedValues = append(convertedValues, IndicatorValue{
					Timestamp: val.Timestamp,
					Value:     val.Value,
				})
			}

			indicatorMap[uniqueKey] = convertedValues
		}

		result[instrumentIDStr] = indicatorMap
	}

	return result, nil
}

// processExitConditions checks and processes exit conditions for an open position
// Returns true if position was closed, false if position remains open
func (e *BacktestEngine) processExitConditions(
	openPosition *model.Position,
	strategy *strategyDto.StrategyResponse,
	indicatorArrayMap map[string][]parser.IndicatorValues,
	currentTimestamp time.Time,
	i int,
	positions *[]*model.Position,
	capital *float64,
) bool {
	// 1. Create instrument object first
	currentInstrument := &instrumentModel.Instrument{
		Open:      indicatorArrayMap["OPEN"][i].Value,
		High:      indicatorArrayMap["HIGH"][i].Value,
		Low:       indicatorArrayMap["LOW"][i].Value,
		Close:     indicatorArrayMap["CLOSE"][i].Value,
		Volume:    indicatorArrayMap["VOLUME"][i].Value,
		Timestamp: currentTimestamp,
	}

	// 2. Update trailing stop loss BEFORE other checks
	e.updateStopLossBasedOnTrailing(openPosition, strategy, currentInstrument)

	// 3. Check stop loss and target profit
	stopLossHit := e.isStopLossTriggered(openPosition, currentInstrument)
	targetProfitHit := e.isTargetProfitTriggered(openPosition, strategy, currentInstrument)

	// 4. Evaluate exit expression last (with updated trailing stop)
	var exitResult bool
	var err error
	if strategy.ExitExpression != "" {
		exitResult, err = parser.EvaluateIndicatorExpression(
			strategy.ExitExpression, indicatorArrayMap, currentTimestamp, nil, i,
		)
		if err != nil {
			println("Error evaluating exit expression:", err.Error())
		}
	}

	// Check if exit condition is met
	if exitResult || stopLossHit || targetProfitHit {
		// Determine exit reason

		exitReason := string(model.ExitReasonSignal)
		exitPrice := currentInstrument.Close

		if stopLossHit {
			// Check if this is a trailing stop loss exit
			if openPosition.InitialStopLossPrice != nil &&
				openPosition.StopLossPrice != nil &&
				*openPosition.StopLossPrice != *openPosition.InitialStopLossPrice {
				exitReason = string(model.ExitReasonTrailingStopLoss)
			} else {
				exitReason = string(model.ExitReasonStopLoss)
			}
			exitPrice = *openPosition.StopLossPrice
		} else if targetProfitHit {
			exitReason = string(model.ExitReasonTargetProfit)
			exitPrice = *openPosition.TargetProfitPrice
		}

		// Calculate fees
		exitFees := 0.0

		// Close the position
		openPosition.Exit(
			exitPrice,
			currentTimestamp,
			exitPrice,
			currentTimestamp,
			strategy.ExitExpression,
			exitReason,
			exitFees,
			nil,
		)

		// Collect position
		*positions = append(*positions, openPosition)

		// Add net profit/loss if available
		if openPosition.NetProfitLoss != nil {
			*capital += *openPosition.NetProfitLoss
		}

		// Position was closed
		return true
	}

	// Position remains open
	return false
}

// checkAndCreateNewPosition evaluates entry conditions and creates a new position if conditions are met
// Returns the created position and updated capital, or nil position if no entry was made
func (e *BacktestEngine) checkAndCreateNewPosition(
	strategy *strategyDto.StrategyResponse,
	indicatorArrayMap map[string][]parser.IndicatorValues,
	currentTimestamp time.Time,
	i int,
	run *model.Run,
	capital float64,
) (*model.Position, float64) {
	// Evaluate entry expression using indicator array map
	entryResult, err := parser.EvaluateIndicatorExpression(
		strategy.EntryExpression, indicatorArrayMap, currentTimestamp, nil, i,
	)
	if err != nil {
		println("Error evaluating entry expression:", err.Error())
	}

	// Check if the current time is within the entry time window
	// If EntryStartTime or EntryEndTime is specified, we only look for entries within that window
	canCheckEntryCondition := true

	if strategy.EntryStartTime != nil || strategy.EntryEndTime != nil {
		// Extract time of day from the current candle timestamp
		currentTime := currentTimestamp
		currentTimeOfDay := time.Date(0, 1, 1,
			currentTime.Hour(), currentTime.Minute(), currentTime.Second(), 0, time.UTC)

		// Check if we're within the entry time window
		if strategy.EntryStartTime != nil {
			// Extract time component from EntryStartTime
			entryStartTimeOfDay := time.Date(0, 1, 1,
				strategy.EntryStartTime.Hour(), strategy.EntryStartTime.Minute(),
				strategy.EntryStartTime.Second(), 0, time.UTC)

			if currentTimeOfDay.Before(entryStartTimeOfDay) {
				canCheckEntryCondition = false
			}
		}

		if strategy.EntryEndTime != nil {
			// Extract time component from EntryEndTime
			entryEndTimeOfDay := time.Date(0, 1, 1,
				strategy.EntryEndTime.Hour(), strategy.EntryEndTime.Minute(),
				strategy.EntryEndTime.Second(), 0, time.UTC)

			if currentTimeOfDay.After(entryEndTimeOfDay) {
				canCheckEntryCondition = false
			}
		}
	}

	// Check if MIS square-off time has been reached (3:15 PM IST)
	// Prevent new MIS entries after this time
	if strategy.TradeType == strategyModel.TradeType("MIS") {
		if e.IsMISEndOfDayHit(currentTimestamp) {
			canCheckEntryCondition = false
		}
	}

	// Only check entry conditions if we're within the allowed entry time window
	if canCheckEntryCondition && entryResult {
		// Get current close price from CLOSE indicator
		var currentClosePrice float64
		if closeIndicator, exists := indicatorArrayMap["CLOSE"]; exists {
			for _, val := range closeIndicator {
				if val.Timestamp.Equal(currentTimestamp) {
					currentClosePrice = val.Value
					break
				}
			}
		}

		if currentClosePrice <= 0 {
			// Return nil if we can't get a valid close price
			return nil, capital
		}

		// Calculate position size based on strategy
		positionSize := e.calculatePositionSize(strategy, currentClosePrice)

		// Determine position direction from strategy
		var direction model.Direction
		if strategy.PositionType == "LONG" || strategy.PositionType == "long" {
			direction = model.Long
		} else {
			direction = model.Short
		}

		// Create new position
		newPosition := model.NewPosition(
			run.ID,
			direction,
			currentClosePrice,
			currentTimestamp,
			strategy.EntryExpression,
			positionSize,
		)

		// Simulate order execution
		entryPrice := currentClosePrice
		entryFees := 0.0
		newPosition.Enter(entryPrice, currentTimestamp, entryFees, nil)

		newPosition.StopLossPrice = e.calculateStopLoss(
			*newPosition.EntryPriceOrder,
			strategy.StopLoss,
			strategy.StopLossType,
			strategy.PositionType,
		)

		// Store the initial stop loss price for trailing stop loss detection
		// Create a copy of the value to avoid both pointers pointing to the same memory
		if newPosition.StopLossPrice != nil {
			initialStopLoss := *newPosition.StopLossPrice
			newPosition.InitialStopLossPrice = &initialStopLoss
		}

		newPosition.TargetProfitPrice = e.calculateTargetProfit(
			*newPosition.EntryPriceOrder,
			strategy.TargetProfit,
			strategy.TargetProfitType,
			strategy.PositionType,
		)

		// Update capital for position entry
		updatedCapital := e.updateCapitalOnPositionChange(capital, newPosition, entryPrice, entryFees, true)

		return newPosition, updatedCapital
	}

	// No entry conditions met, return nil position with unchanged capital
	return nil, capital
}

func (e *BacktestEngine) handleMISSquareOff(
	openPosition *model.Position,
	strategy *strategyDto.StrategyResponse,
	indicatorArrayMap map[string][]parser.IndicatorValues,
	currentTimestamp time.Time,
	positions *[]*model.Position,
	capital *float64,
) bool {
	// Only process MIS positions that are entered
	if openPosition == nil || openPosition.State != model.Entered ||
		strategy.TradeType != strategyModel.TradeType("MIS") {
		return false
	}

	// Get close price for current timestamp
	var closePrice float64
	if closeIndicator, exists := indicatorArrayMap["CLOSE"]; exists {
		for _, val := range closeIndicator {
			if val.Timestamp.Equal(currentTimestamp) {
				closePrice = val.Value
				break
			}
		}
	}

	if closePrice <= 0 {
		return false
	}

	// Calculate fees and exit details
	exitFees := 0.0
	exitOrderPrice := closePrice
	exitReason := string(model.ExitReasonMISSquareOff)

	// Close the position
	openPosition.Exit(
		closePrice,
		currentTimestamp,
		exitOrderPrice,
		currentTimestamp,
		"", // No exit expression for auto square-offs
		exitReason,
		exitFees,
		nil,
	)

	// Collect position
	*positions = append(*positions, openPosition)

	// Add net profit/loss if available
	if openPosition.NetProfitLoss != nil {
		*capital += *openPosition.NetProfitLoss
	}

	return true // Position was closed
}
