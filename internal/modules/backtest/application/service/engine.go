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
	timeUtil "algo/internal/util/time"

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
	optionLegService  OptionLegService
	runGroupRepo      backtestRepo.RunGroupRepository
	runRepo           backtestRepo.RunRepository
	positionRepo      backtestRepo.PositionRepository
	runMetricsRepo    backtestRepo.RunMetricsRepository
	metricsService    MetricsCalculationService
	marketHoursUtil   *timeUtil.MarketHoursUtil
}

// We're using instrumentModel.Instrument instead of defining our own struct

// NewBacktestEngine creates a new backtest engine
func NewBacktestEngine(
	backtestService *BacktestService,
	strategyService strategyService.StrategyService,
	instrumentService *instrumentService.InstrumentService,
	optionLegService OptionLegService,
	runGroupRepo backtestRepo.RunGroupRepository,
	runRepo backtestRepo.RunRepository,
	positionRepo backtestRepo.PositionRepository,
	runMetricsRepo backtestRepo.RunMetricsRepository,
) EngineService {
	return &BacktestEngine{
		backtestService:   backtestService,
		strategyService:   strategyService,
		instrumentService: instrumentService,
		optionLegService:  optionLegService,
		runGroupRepo:      runGroupRepo,
		runRepo:           runRepo,
		positionRepo:      positionRepo,
		runMetricsRepo:    runMetricsRepo,
		metricsService:    NewMetricsCalculationService(),
		marketHoursUtil:   timeUtil.NewMarketHoursUtil(),
	}
}

// Start starts a backtest run for the provided strategy
func (e *BacktestEngine) Start(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]Response, error) {
	// Get the strategy
	strategyResp, err := e.strategyService.GetStrategyByID(ctx, strategyID, userID)
	if err != nil {
		return nil, err
	}

	// Temporary Placement
	if strategyResp.Segment == strategyModel.Segment("F&O") {
		e.simpleBackTestFlow(ctx, strategyID, userID)
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
		if strategyResp.Qty != nil && *strategyResp.Qty > 0 {
			qty = *strategyResp.Qty
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

	// Calculate total commission and slippage across all runs
	totalPositionsCount := 0
	for _, run := range runs {
		totalPositionsCount += len(run.Positions)
	}

	// Update TotalCommission and TotalSlippage for each run's metrics
	for _, run := range runs {
		if run.RunMetrics != nil {
			// Calculate commission: number of positions * 2 (entry + exit) * brokerage
			commission := float64(totalPositionsCount * 2 * strategyResp.Brokerage)
			run.RunMetrics.TotalCommission = commission

			// Calculate slippage: slippage percentage calculation
			// Slippage = (slippage_percentage / 100) * total_trade_value
			totalTradeValue := 0.0
			for _, position := range run.Positions {
				if position.EntryPriceOrder != nil {
					totalTradeValue += *position.EntryPriceOrder * position.PositionSize
				}
			}
			slippage := (strategyResp.Slippage / 100.0) * totalTradeValue
			run.RunMetrics.TotalSlippage = slippage

			// Subtract commission and slippage from total P&L
			if run.RunMetrics.TotalReturnPct != 0 {
				// Convert percentage back to absolute value to subtract costs
				initialCapital := 100000.0 // Default capital
				if strategyResp.Capital != nil {
					initialCapital = *strategyResp.Capital
				}
				totalPnL := (run.RunMetrics.TotalReturnPct / 100.0) * initialCapital
				adjustedPnL := totalPnL - commission - slippage
				run.RunMetrics.TotalReturnPct = (adjustedPnL / initialCapital) * 100
			}
		}
	}

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
	for i := range len(timestamps) - 1 {
		// Get current timestamp for this iteration
		currentTimestamp := timestamps[i]

		// Check if we've moved to a new trading day
		dateTruncated := currentTimestamp.Truncate(24 * time.Hour)
		if !dateTruncated.Equal(currentTradingDay) {
			// If we have an open position and this is an MIS (intraday) trade, close it at the end of the previous day
			if openPosition != nil && openPosition.State == model.Entered && strategy.TradeType != nil && *strategy.TradeType == strategyModel.TradeType("MIS") {
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
					if e.handleMISSquareOff(openPosition, strategy, indicatorArrayMap, *lastTimestampFromPreviousDay, &positions, capital) {
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
				*capital,
			)

			if newPosition != nil {
				openPosition = newPosition
				*capital = updatedCapital
			}
		}

		// If position is still open, check exit condition
		if openPosition != nil {
			// Process exit conditions and check if position was closed
			if e.processExitConditions(openPosition, strategy, indicatorArrayMap, currentTimestamp, i, &positions, capital) {
				// Position was closed, reset open position and continue to next iteration
				openPosition = nil
				continue
			}

			// Close the position if MIS end-of-day is hit
			if strategy.TradeType != nil && *strategy.TradeType == strategyModel.TradeType("MIS") && e.marketHoursUtil.IsMISEndOfDayHit(currentTimestamp, run.InstrumentMaster.Exchange) {
				// Use the consolidated MIS square-off function
				if e.handleMISSquareOff(openPosition, strategy, indicatorArrayMap, currentTimestamp, &positions, capital) {
					openPosition = nil
				}
			}
		}
	}

	// Handle any remaining open position at the end of backtest period
	if openPosition != nil && openPosition.State == model.Entered {
		// If this is a MIS position, close it at the end of the backtest
		if strategy.TradeType != nil && *strategy.TradeType == strategyModel.TradeType("MIS") && len(timestamps) > 0 {
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
				*capital += *openPosition.NetProfitLoss
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
	metrics := e.calculateMetrics(positions, initialCapital, *capital, run)

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

	initialStopDistance := *e.calculateStopLoss(entryPrice, strategy.StopLoss, strategy.StopLossType,
		func() strategyModel.Direction {
			if strategy.PositionType != nil {
				return *strategy.PositionType
			}
			return strategyModel.Direction("LONG") // default fallback
		}())

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
	// Get values from pointers, using defaults if nil
	qty := 0
	if strategy.Qty != nil {
		qty = *strategy.Qty
	}

	maxAllocation := 0.0
	if strategy.MaxAllocation != nil {
		maxAllocation = *strategy.MaxAllocation
	}

	// 1. If strategy.qty <= 0 and strategy.max_allocation <= 0 throw an error
	if qty <= 0 && maxAllocation <= 0 {
		panic("Invalid position size parameters")
	}

	// 2. If strategy.qty > 0 and strategy.max_allocation <= 0 then return strategy.qty
	if qty > 0 && maxAllocation <= 0 {
		return float64(qty)
	}

	// 3. If strategy.qty <= 0 and strategy.max_allocation > 0 then return strategy.max_allocation / price as integer
	// Ensure at least 1 share is purchased if the calculation results in less than 0
	if qty <= 0 && maxAllocation > 0 {
		quantity := math.Floor(maxAllocation / price)
		if quantity <= 0 {
			return 1
		}
		return quantity
	}

	// 4. If strategy.qty > 0 and strategy.max_allocation > 0,
	// then return the lowest between (strategy.qty, strategy.max_allocation / price) as integer
	qtyValue := float64(qty) * price
	if qtyValue < maxAllocation {
		return float64(qty)
	} else {
		quantity := math.Floor(maxAllocation / price)
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

	// Check if the current time is within the entry time window using the centralized function
	canCheckEntryCondition := e.marketHoursUtil.IsWithinEntryTimeWindow(currentTimestamp, strategy, run.InstrumentMaster.Exchange)

	// Check if MIS square-off time has been reached (3:15 PM IST)
	// Prevent new MIS entries after this time
	if strategy.TradeType != nil && *strategy.TradeType == strategyModel.TradeType("MIS") {
		if e.marketHoursUtil.IsMISEndOfDayHit(currentTimestamp, run.InstrumentMaster.Exchange) {
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
		if strategy.PositionType != nil && (string(*strategy.PositionType) == "LONG" || string(*strategy.PositionType) == "long") {
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
			func() strategyModel.Direction {
				if strategy.PositionType != nil {
					return *strategy.PositionType
				}
				return strategyModel.Direction("LONG") // default fallback
			}(),
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
			func() strategyModel.Direction {
				if strategy.PositionType != nil {
					return *strategy.PositionType
				}
				return strategyModel.Direction("LONG") // default fallback
			}(),
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
		strategy.TradeType == nil || *strategy.TradeType != strategyModel.TradeType("MIS") {
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

// GetInstrumentService returns the instrumentService
func (e *BacktestEngine) GetInstrumentService() *instrumentService.InstrumentService {
	return e.instrumentService
}

// getLegSpotData retrieves spot data for all option legs in a strategy
func (e *BacktestEngine) getLegSpotData(
	ctx context.Context,
	strategyResp *strategyDto.StrategyResponse,
) (map[uuid.UUID][]*instrumentModel.Instrument, error) {
	legSpotData := make(map[uuid.UUID][]*instrumentModel.Instrument)
	for _, leg := range strategyResp.OptionLegs {
		legInstrumentSpotData, err := e.instrumentService.GetInstrumentsByTimeRange(
			ctx,
			*leg.ExchangeSecurityID,
			*strategyResp.StartDate,
			*strategyResp.EndDate,
			string(strategyResp.ChartTimeframe),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get leg instrument data: %w", err)
		}
		legSpotData[leg.ID] = legInstrumentSpotData
	}
	return legSpotData, nil
}

// EnterFnOPosition resolves the option contract for all legs and enters FnO positions
func (e *BacktestEngine) EnterFnOPosition(
	ctx context.Context,
	legSpotData map[uuid.UUID][]*instrumentModel.Instrument,
	spotTimestamp time.Time,
	strategyResp *strategyDto.StrategyResponse,
	run *model.Run,
) error {
	// for each of the legs, resolve the contract
	for _, leg := range strategyResp.OptionLegs {
		// get the spot for the leg.exchange_security_id
		var legSpotPrice *float64
		if legSpotRecords, ok := legSpotData[leg.ID]; ok {
			for _, legSpotRecord := range legSpotRecords {
				if legSpotRecord.Timestamp.Equal(spotTimestamp) {
					// get the legSpot.Close
					legSpotPrice = &legSpotRecord.Close
					break
				}
			}
		}

		// Use legSpotPrice for processing (placeholder to avoid unused variable error)
		if legSpotPrice != nil {
			// Resolve option contract for this leg
			optionContract, err := e.resolveOptionContract(ctx, leg, *legSpotPrice, spotTimestamp)
			if err != nil {
				// Handle error appropriately - for now, log and continue
				continue
			}

			// Take position in the above resolved optionContract based on leg
			if optionContract != nil {
				// For option contracts, we need to fetch the actual option data to get the premium
				// For now, we'll use a default premium price until we get actual option data
				premiumPrice := 120.0 // This should be fetched from actual option data

				// Calculate the total quantity based on leg configuration
				lotSize := 1
				if leg.LotSize != nil {
					lotSize = *leg.LotSize
				}
				totalQuantity := float64(leg.TotalLot * lotSize)

				// Determine the direction based on leg action type (BUY/SELL)
				var direction model.Direction
				if leg.ActionType == "BUY" {
					direction = model.Long
				} else {
					direction = model.Short
				}

				// Create a new position for this option leg
				strikeDisplay := *legSpotPrice // Default to spot price if no strike available
				if optionContract != nil {
					// Use contract's strike if available, otherwise use leg spot price
					strikeDisplay = *legSpotPrice // This would be replaced with actual strike from contract
				}

				position := model.NewPosition(
					run.ID,
					direction,
					premiumPrice,
					spotTimestamp,
					fmt.Sprintf("Option Leg: %s %s %.2f Strike", string(*leg.OptionType), string(leg.ActionType), strikeDisplay),
					totalQuantity,
				)

				// Simulate order execution with the premium price
				entryFees := 0.0 // Calculate actual fees based on brokerage
				position.Enter(premiumPrice, spotTimestamp, entryFees, nil)

				// Store the option contract ExchangeSecurityID in the entry signal for reference
				position.EntrySignal = fmt.Sprintf("Option Contract ID: %d, %s",
					optionContract.ExchangeSecurityID, position.EntrySignal)

				// Store position in run for later processing
				if run.Positions == nil {
					run.Positions = make([]*model.Position, 0)
				}
				run.Positions = append(run.Positions, position)

				fmt.Printf("Created position for leg %s: %s %s %.0f lots at premium %.2f\n",
					leg.ID, string(*leg.OptionType), string(leg.ActionType),
					totalQuantity/float64(lotSize), premiumPrice)
			}

		}
	}

	return nil
}

// executeTimeBasedStrategy executes a time-based strategy by iterating through dates and checking entry/exit times
func (e *BacktestEngine) executeTimeBasedStrategy(
	ctx context.Context,
	strategyResp *strategyDto.StrategyResponse,
	runGroup *model.RunGroup,
) ([]*model.Run, error) {
	// getLegSpotData
	legSpotData, err := e.getLegSpotData(ctx, strategyResp)
	if err != nil {
		return nil, err
	}

	// Get the ingredients: 'entry_time', 'exit_time'
	// NOTE: fix the below issue, it has to EntryTime and ExitTime
	entry_time := strategyResp.EntryStartTime
	entryTimeAsTime := timeUtil.TimeOnlyToTime(entry_time)

	exitTime := strategyResp.EntryEndTime
	exitTimeAsTime := timeUtil.TimeOnlyToTime(exitTime)

	// Track the current active run and its positions
	var activeRun *model.Run = nil

	// Track all closed runs instead of persisting to database
	var closedRuns []*model.Run

	// Iterate from start date to end date in 1-minute timeframe from 9:15 to 3:29
	currentDate := *strategyResp.StartDate
	endDate := *strategyResp.EndDate

	// Create IST timezone (UTC+5:30)
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		// Fallback to fixed offset if location loading fails
		ist = time.FixedZone("IST", 5*60*60+30*60)
	}

	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		// Skip weekends (Saturday = 6, Sunday = 0)
		if currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday {
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		// Start time: 9:15 AM IST
		startTime := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 9, 15, 0, 0, ist)
		// End time: 3:29 PM IST
		endTime := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 15, 29, 0, 0, ist)

		// Iterate in 1-minute intervals
		for currentTime := startTime; currentTime.Before(endTime) || currentTime.Equal(endTime); currentTime = currentTime.Add(time.Minute) {
			// *************************Remove This**************************
			entryTimeAsHour := entryTimeAsTime.Hour()
			entryTimeAsMinute := entryTimeAsTime.Minute()

			exitTimeAsHour := exitTimeAsTime.Hour()
			exitTimeAsMinute := exitTimeAsTime.Minute()

			currentTimeAsHour := currentTime.Hour()
			currentTimeAsMinute := currentTime.Minute()

			fmt.Print(exitTimeAsHour, ":", exitTimeAsMinute, " - ", entryTimeAsHour, ":", entryTimeAsMinute, " - ", currentTimeAsHour, ":", currentTimeAsMinute, "\n")
			// *************************Remove This**************************

			// Entry check - only if no active run exists
			if activeRun == nil && currentTime.Hour() == entryTimeAsTime.Hour() && currentTime.Minute() == entryTimeAsTime.Minute() {
				// resolve the contract for all legs
				// Create a run object for this entry
				activeRun = model.NewRun(
					runGroup.ID,
					strategyResp.ID,
					"NSE", // Default exchange, should be determined from strategy or instruments
					0,     // InstrumentID will be set based on the option contracts
					model.Backtest,
					1, // Default quantity, should be determined from strategy
					currentTime,
					currentTime.Add(24*time.Hour), // End time, adjust as needed
					nil,                           // InstrumentMaster will be set in EnterFnOPosition
				)

				// Pass the run object to EnterFnOPosition to create and enter all positions
				err := e.EnterFnOPosition(ctx, legSpotData, currentTime, strategyResp, activeRun)
				if err != nil {
					// Handle error appropriately - for now, log and continue
					println("Error entering FnO position:", err.Error())
					activeRun = nil // Reset run on error
					continue
				}

				// Log successful entry
				println("Created new run with", len(activeRun.Positions), "positions at", currentTime.Format(time.RFC3339))
			}

			// Exit check - only if we have an active run with positions
			if activeRun != nil && len(activeRun.Positions) > 0 && currentTime.Hour() == exitTimeAsTime.Hour() && currentTime.Minute() == exitTimeAsTime.Minute() {
				// Exit all positions using the new function
				e.ExitFnOPosition(activeRun, currentTime, &closedRuns)

				// Reset active run after adding to closed runs
				activeRun = nil
			}
		}

		// Move to next day
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Handle any remaining open run at the end of backtest period
	if activeRun != nil && len(activeRun.Positions) > 0 {
		// Force exit all remaining positions at market close
		lastTimestamp := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 15, 30, 0, 0, ist)

		for _, position := range activeRun.Positions {
			if position.State == model.Entered {
				exitPrice := *position.EntryPriceOrder // Simplified
				exitFees := 0.0

				position.Exit(
					exitPrice,
					lastTimestamp,
					exitPrice,
					lastTimestamp,
					"END_OF_BACKTEST_EXIT",
					string(model.ExitReasonSignal),
					exitFees,
					nil,
				)
			}
		}

		// Mark run as ended and persist
		activeRun.EndedTime = &lastTimestamp
		activeRun.State = model.Ended

		// Skip metrics calculation for now

		// Add final completed run to closedRuns list instead of persisting
		closedRuns = append(closedRuns, activeRun)
		println("Added final run with ID:", activeRun.ID.String(), "to closed runs list")
		println("Total closed runs:", len(closedRuns))

		activeRun = nil
	}

	// Log final summary of completed runs
	println("Time-based strategy execution completed.")
	println("Total closed runs collected:", len(closedRuns))
	for i, run := range closedRuns {
		println(fmt.Sprintf("Run %d: ID=%s, Positions=%d, State=%s",
			i+1, run.ID.String(), len(run.Positions), string(run.State)))
	}

	return closedRuns, nil
}

// executeConditionalStrategy executes a conditional strategy by evaluating entry/exit expressions against market data
func (e *BacktestEngine) executeConditionalStrategy(
	ctx context.Context,
	strategyResp *strategyDto.StrategyResponse,
	runGroup *model.RunGroup,
) ([]*model.Run, error) {
	// getLegSpotData
	legSpotData, err := e.getLegSpotData(ctx, strategyResp)
	if err != nil {
		return nil, err
	}

	// Get the main instrument data for conditional evaluation
	if len(strategyResp.Instruments) == 0 {
		return nil, fmt.Errorf("no instruments found in strategy")
	}

	instrumentID := strategyResp.Instruments[0]
	instrumentSpotData, err := e.instrumentService.GetInstrumentsByTimeRange(
		ctx,
		int32(instrumentID),
		*strategyResp.StartDate,
		*strategyResp.EndDate,
		string(strategyResp.ChartTimeframe),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get main instrument data: %w", err)
	}

	// Initialize evaluation context
	evalContext := parser.NewEvaluationContext(nil)
	evalContext.Strategy = strategyResp
	evalContext.Data = e.toDataContext(instrumentSpotData)

	// Extract indicators from both entry and exit expressions
	entryIndicators, err := parser.ExtractIndicators(strategyResp.EntryExpression)
	if err != nil {
		return nil, fmt.Errorf("error extracting indicators from entry expression: %w", err)
	}

	var exitIndicators []parser.IndicatorInfo
	if strategyResp.ExitExpression != "" {
		exitIndicators, err = parser.ExtractIndicators(strategyResp.ExitExpression)
		if err != nil {
			return nil, fmt.Errorf("error extracting indicators from exit expression: %w", err)
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
			return nil, fmt.Errorf("error generating unique key for indicator: %w", err)
		}

		// Skip if already calculated (avoid duplicates between entry and exit)
		if _, exists := indicatorArrayMap[uniqueString]; exists {
			continue
		}

		indicatorArrayMap[uniqueString] = calculator.CalculateAll(indicator.Arguments, evalContext, *strategyResp.StartDate)
		if len(indicatorArrayMap[uniqueString]) > length {
			length = len(indicatorArrayMap[uniqueString])
		}
	}

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

	// Sort timestamps
	// (Assuming timestamps are already in chronological order from data source)

	// Track open positions
	var openPositions []*model.Position

	// Iterate through each timestamp and evaluate conditions
	for i, timestamp := range timestamps {
		// Skip weekends (Saturday = 6, Sunday = 0)
		if timestamp.Weekday() == time.Saturday || timestamp.Weekday() == time.Sunday {
			continue
		}

		// Skip if timestamp is outside strategy date range
		if timestamp.Before(*strategyResp.StartDate) || timestamp.After(*strategyResp.EndDate) {
			continue
		}

		// Check entry conditions if we don't have open positions
		if len(openPositions) == 0 {
			entryResult, err := parser.EvaluateIndicatorExpression(
				strategyResp.EntryExpression,
				indicatorArrayMap,
				timestamp,
				nil,
				i,
			)
			if err != nil {
				// Log error and continue
				continue
			}

			if entryResult {
				// Create a run object for this entry
				run := model.NewRun(
					runGroup.ID,
					strategyResp.ID,
					"NSE", // Default exchange, should be determined from strategy or instruments
					0,     // InstrumentID will be set based on the option contracts
					model.Backtest,
					1, // Default quantity, should be determined from strategy
					timestamp,
					timestamp.Add(24*time.Hour), // End time, adjust as needed
					nil,                         // InstrumentMaster will be set in EnterFnOPosition
				)

				// Enter position - resolve contracts for all legs
				err := e.EnterFnOPosition(ctx, legSpotData, timestamp, strategyResp, run)
				if err != nil {
					// Handle error appropriately - for now, log and continue
					continue
				}

				// Create a position record (simplified)
				position := &model.Position{
					ID:               uuid.New(),
					RunID:            run.ID, // Use the actual run ID
					State:            model.NotEntered,
					Direction:        model.Long, // This should be determined by strategy logic
					EntryPriceSignal: 0.0,        // This should be determined from current price
					EntryTimeSignal:  timestamp,
					EntrySignal:      strategyResp.EntryExpression,
					PositionSize:     1.0, // This should be calculated based on strategy
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				}
				openPositions = append(openPositions, position)
			}
		}

		// Check exit conditions if we have open positions
		if len(openPositions) > 0 && strategyResp.ExitExpression != "" {
			exitResult, err := parser.EvaluateIndicatorExpression(
				strategyResp.ExitExpression,
				indicatorArrayMap,
				timestamp,
				nil,
				i,
			)
			if err != nil {
				// Log error and continue
				continue
			}

			if exitResult {
				// Exit all positions
				openPositions = nil // Clear all positions
			}
		}
	}

	// For conditional strategy, return empty runs list as it's not implemented yet
	return []*model.Run{}, nil
}

func (e *BacktestEngine) simpleBackTestFlow(
	ctx context.Context,
	strategyID uuid.UUID,
	userID uuid.UUID,
) (interface{}, error) {
	// 0.0 Get the necessary ingredients

	// 0.1 Get the strategy data along with legs
	strategyResp, err := e.strategyService.GetStrategyByID(ctx, strategyID, userID)
	if err != nil {
		return nil, err
	}
	if strategyResp.Segment != strategyModel.Segment("F&O") {
		return nil, fmt.Errorf("Not a F&O strategy")
	}

	// 0.2 Create a run group in memory
	runGroup := model.NewRunGroup(strategyID)
	var closedRuns []*model.Run

	if strategyResp.StrategyType == strategyModel.Conditional {
		// Backtest conditional strategy
		runs, err := e.executeConditionalStrategy(ctx, strategyResp, runGroup)
		if err != nil {
			return nil, fmt.Errorf("failed to execute conditional strategy: %w", err)
		}
		closedRuns = runs
	} else if strategyResp.StrategyType == strategyModel.TimeBased {
		// Backtest time-based strategy
		runs, err := e.executeTimeBasedStrategy(ctx, strategyResp, runGroup)
		if err != nil {
			return nil, fmt.Errorf("failed to execute time-based strategy: %w", err)
		}
		closedRuns = runs
	}

	// Persist the run group record
	if err := e.runGroupRepo.Create(ctx, runGroup); err != nil {
		return nil, err
	}

	// Persist the run records
	if len(closedRuns) > 0 {
		if err := e.runRepo.BulkCreate(ctx, closedRuns); err != nil {
			return nil, fmt.Errorf("failed to persist run records: %w", err)
		}
		println("Successfully persisted", len(closedRuns), "run records")
	}

	// Persist the position records
	allPositions := make([]*model.Position, 0)
	for _, run := range closedRuns {
		if len(run.Positions) > 0 {
			allPositions = append(allPositions, run.Positions...)
		}
	}

	if len(allPositions) > 0 {
		if err := e.positionRepo.BulkCreate(ctx, allPositions); err != nil {
			return nil, fmt.Errorf("failed to persist position records: %w", err)
		}
		println("Successfully persisted", len(allPositions), "position records")
	}

	return strategyResp, nil
}

// executeFOBacktestInMemory runs a multi-leg options backtest for a strategy, following the pseudocode in backtest-fno.md
func (e *BacktestEngine) executeFOBacktestInMemory(
	strategy *strategyDto.StrategyResponse,
	indexData []*instrumentModel.Instrument, // Spot/index OHLC data
	optionChainData map[string][]*instrumentModel.Instrument, // Option chain by strike+type
	run *model.Run,
) ([]*model.Position, *model.RunMetrics) {
	// Start tracking execution time
	functionStartTime := time.Now()

	var (
		openPositions  []*model.Position // All legs
		positionStatus string
		entryLegs      []struct {
			Contract *instrumentModel.Instrument
			Leg      *strategyDto.OptionLegResponse
		}
		exitReason string
		totalPnL   float64
		capital    = strategy.Capital
	)

	// 1. ENTRY FLOW: Loop through indexData (spot/index OHLC)
	for _, indexOHLC := range indexData {
		now := indexOHLC.Timestamp

		// Check entry window
		if !e.marketHoursUtil.IsWithinEntryTimeWindow(now, strategy, run.InstrumentMaster.Exchange) {
			continue
		}

		// Evaluate entry condition (on index/spot)
		evalContext := parser.NewEvaluationContext(nil)
		evalContext.Strategy = strategy
		evalContext.Data = e.toDataContext([]*instrumentModel.Instrument{indexOHLC})
		entryOk := parser.Execute(strategy.EntryExpression, evalContext)
		if !entryOk {
			continue
		}

		// Resolve contracts for each leg
		entryLegs = nil
		for _, leg := range strategy.OptionLegs {
			// Determine strike price
			strike := 0.0
			if leg.StrikePrice != nil {
				strike = *leg.StrikePrice
			} else if string(leg.StrikeType) == "ATM" {
				strike = e.resolveATMStrike(indexOHLC.Close, optionChainData)
			}
			// Find contract in optionChainData
			contracts := optionChainData[string(*leg.OptionType)]
			var contract *instrumentModel.Instrument
			for _, c := range contracts {
				// Instrument does not have StrikePrice, so use Close as a proxy for now
				if c.Close == strike {
					contract = c
					break
				}
			}
			if contract == nil {
				continue // skip this leg if contract not found
			}
			entryLegs = append(entryLegs, struct {
				Contract *instrumentModel.Instrument
				Leg      *strategyDto.OptionLegResponse
			}{contract, leg})
		}
		if len(entryLegs) == 0 {
			continue
		}

		// Execute trades for all legs
		for _, leg := range entryLegs {
			pos := model.NewPosition(
				run.ID,
				model.Long, // Only buy for now
				leg.Contract.Close,
				now,
				"ENTRY",
				float64(leg.Leg.TotalLot),
			)
			pos.State = model.Entered
			pos.EntryPriceOrder = &leg.Contract.Close
			pos.EntryTimeOrder = &now
			pos.EntrySignal = "ENTRY"
			// Option meta (store as needed in custom fields or use EntrySignal/ExitSignal for debugging)
			openPositions = append(openPositions, pos)
		}
		positionStatus = "ENTERED"
		break // Only one multi-leg position at a time
	}

	// 2. EXIT FLOW: Monitor exit conditions if position is entered
	if positionStatus == "ENTERED" && len(openPositions) > 0 {
		for _, indexOHLC := range indexData {
			now := indexOHLC.Timestamp

			// Get premium data for all entered positions (legs)
			legPremiums := make([]float64, len(openPositions))
			for i, pos := range openPositions {
				// Find premium for this contract at current time
				var premium float64
				// We don't have OptionType/StrikePrice fields, so match by EntryPriceOrder and time
				for _, cList := range optionChainData {
					for _, c := range cList {
						if pos.EntryPriceOrder != nil && c.Timestamp.Equal(now) && c.Close == *pos.EntryPriceOrder {
							premium = c.Close
							break
						}
					}
				}
				legPremiums[i] = premium
			}

			// Calculate cumulative P&L across all legs
			totalPnL = 0
			for i, pos := range openPositions {
				entry := 0.0
				if pos.EntryPriceOrder != nil {
					entry = *pos.EntryPriceOrder
				}
				pnl := (legPremiums[i] - entry) * pos.PositionSize
				totalPnL += pnl
			}

			// Check exit conditions (priority order)
			exitTriggered := false
			// 1. Strategy-based exit on spot
			evalContext := parser.NewEvaluationContext(nil)
			evalContext.Strategy = strategy
			evalContext.Data = e.toDataContext([]*instrumentModel.Instrument{indexOHLC})
			if parser.Execute(strategy.ExitExpression, evalContext) {
				exitTriggered = true
				exitReason = string(model.ExitReasonSignal)
			} else if strategy.TargetProfit > 0 && totalPnL >= strategy.TargetProfit {
				exitTriggered = true
				exitReason = string(model.ExitReasonTargetProfit)
			} else if strategy.StopLoss > 0 && totalPnL <= -strategy.StopLoss {
				exitTriggered = true
				exitReason = string(model.ExitReasonStopLoss)
			} else if strategy.TrailingStopLoss > 0 && e.checkTrailingStopLoss(totalPnL, strategy) {
				exitTriggered = true
				exitReason = string(model.ExitReasonTrailingStopLoss)
			} else if strategy.TradeType != nil && *strategy.TradeType == strategyModel.TradeType("MIS") && e.marketHoursUtil.IsMISEndOfDayHit(now, run.InstrumentMaster.Exchange) {
				exitTriggered = true
				exitReason = string(model.ExitReasonMISSquareOff)
			}

			if exitTriggered {
				// Close all legs
				for i, pos := range openPositions {
					pos.Exit(
						legPremiums[i],
						now,
						legPremiums[i],
						now,
						"EXIT",
						exitReason,
						0,
						nil,
					)
				}
				positionStatus = "EXITED"
				break
			}
		}
	}

	// 3. Metrics calculation and return
	closedPositions := make([]*model.Position, 0)
	for _, pos := range openPositions {
		if pos.State == model.Exited {
			closedPositions = append(closedPositions, pos)
		}
	}
	metrics := e.calculateMetrics(closedPositions, 0, *capital, run)

	functionEndTime := time.Now()
	println("\tFO Backtest in-memory calculation completed:",
		len(closedPositions), "positions generated,",
		"\n\t\tStart:", functionStartTime.Format(time.RFC3339),
		"\n\t\tEnd:", functionEndTime.Format(time.RFC3339),
		"\n\t\tDuration:", functionEndTime.Sub(functionStartTime).String())

	return closedPositions, metrics
}

// resolveStrikePriceByType resolves the actual strike price based on the strike type and option type
// Supports strike types like: ATM, ITM1, ITM2, ITM3, OTM1, OTM2, OTM3, etc.
func (e *BacktestEngine) resolveStrikePriceByType(
	strikeType string,
	optionType string,
	spotPrice float64,
	strikeOffset float64,
) (float64, error) {
	if strikeOffset <= 0 {
		strikeOffset = 50 // Default strike offset
	}

	// Handle ATM (At-The-Money)
	if strikeType == "ATM" {
		return e.NearestStrikeByOffset(spotPrice, strikeOffset), nil
	}

	// Parse strike type to extract ITM/OTM and the level (1, 2, 3, etc.)
	var isITM bool
	var level int
	var err error

	if len(strikeType) >= 4 {
		prefix := strikeType[:3]
		levelStr := strikeType[3:]

		switch prefix {
		case "ITM":
			isITM = true
		case "OTM":
			isITM = false
		default:
			return 0, fmt.Errorf("unsupported strike type prefix: %s", prefix)
		}

		// Parse the level number
		if levelStr == "" {
			level = 1 // Default to level 1 if no number specified
		} else {
			level, err = strconv.Atoi(levelStr)
			if err != nil {
				return 0, fmt.Errorf("invalid strike level: %s", levelStr)
			}
		}
	} else {
		return 0, fmt.Errorf("invalid strike type format: %s", strikeType)
	}

	// Calculate the strike based on ITM/OTM, level, and option type
	var targetStrike float64

	if optionType == "CE" || optionType == "CALL" {
		// For Call options:
		// ITM: strike < spot (lower strikes are more ITM)
		// OTM: strike > spot (higher strikes are more OTM)
		if isITM {
			targetStrike = spotPrice - (float64(level) * strikeOffset)
		} else {
			targetStrike = spotPrice + (float64(level) * strikeOffset)
		}
	} else if optionType == "PE" || optionType == "PUT" {
		// For Put options:
		// ITM: strike > spot (higher strikes are more ITM)
		// OTM: strike < spot (lower strikes are more OTM)
		if isITM {
			targetStrike = spotPrice + (float64(level) * strikeOffset)
		} else {
			targetStrike = spotPrice - (float64(level) * strikeOffset)
		}
	} else {
		return 0, fmt.Errorf("unsupported option type: %s", optionType)
	}

	// Round to nearest strike using the offset
	finalStrike := e.NearestStrikeByOffset(targetStrike, strikeOffset)

	return finalStrike, nil
}

// NearestStrikeByOffset returns the nearest strike price for a given currentPrice and offset.
// For example, currentPrice=22540, offset=100 => 22500; currentPrice=22540, offset=50 => 22550
func (e *BacktestEngine) NearestStrikeByOffset(currentPrice float64, offset float64) float64 {
	if offset <= 0 {
		return currentPrice
	}
	remainder := math.Mod(currentPrice, offset)
	if remainder < offset/2 {
		return currentPrice - remainder
	}
	return currentPrice - remainder + offset
}

// resolveATMStrike returns the nearest ATM strike from option chain data
func (e *BacktestEngine) resolveATMStrike(spot float64, optionChainData map[string][]*instrumentModel.Instrument) float64 {
	// Find the closest strike in the option chain (assume CE chain exists)
	minDiff := 1e9
	atm := 0.0
	for _, contracts := range optionChainData {
		for _, c := range contracts {
			diff := math.Abs(c.Close - spot)
			if diff < minDiff {
				minDiff = diff
				atm = c.Close
			}
		}
	}
	return atm
}

// checkTrailingStopLoss checks if trailing stop loss is hit (stub, expand as needed)
func (e *BacktestEngine) checkTrailingStopLoss(totalPnL float64, strategy *strategyDto.StrategyResponse) bool {
	// Implement trailing logic as per strategy.TrailingStopLossParams
	// For now, just a stub
	return false
}

// calculateStrikePriceForLeg calculates the strike price for a leg based on its configuration
func (e *BacktestEngine) calculateStrikePriceForLeg(leg *strategyDto.OptionLegResponse, spotPrice float64) (float64, error) {
	// If strike price is explicitly set, use it
	if leg.StrikePrice != nil {
		return *leg.StrikePrice, nil
	}

	// Calculate based on strike type
	switch string(leg.StrikeType) {
	case "ATM":
		// For ATM, use nearest strike based on common offsets (50, 100, etc.)
		return e.NearestStrikeByOffset(spotPrice, 50), nil // Default offset of 50
	case "ITM":
		// In-the-money: for CE, strike < spot; for PE, strike > spot
		offset := 100.0 // Default offset
		if leg.OptionType != nil && string(*leg.OptionType) == "CE" {
			return e.NearestStrikeByOffset(spotPrice-offset, 50), nil
		}
		return e.NearestStrikeByOffset(spotPrice+offset, 50), nil
	case "OTM":
		// Out-of-the-money: for CE, strike > spot; for PE, strike < spot
		offset := 100.0 // Default offset
		if leg.OptionType != nil && string(*leg.OptionType) == "CE" {
			return e.NearestStrikeByOffset(spotPrice+offset, 50), nil
		}
		return e.NearestStrikeByOffset(spotPrice-offset, 50), nil
	case "FIXED":
		// Fixed strike should have been provided in StrikePrice field
		if leg.StrikePrice != nil {
			return *leg.StrikePrice, nil
		}
		return 0, fmt.Errorf("fixed strike type specified but no strike price provided")
	default:
		return 0, fmt.Errorf("unsupported strike type: %s", leg.StrikeType)
	}
}

// calculateExpiryDateForLeg calculates the expiry date for a leg based on its configuration
func (e *BacktestEngine) calculateExpiryDateForLeg(leg *strategyDto.OptionLegResponse, currentTime time.Time) (time.Time, error) {
	// If expiry date is explicitly set, use it
	if leg.ExpiryDate != nil {
		return *leg.ExpiryDate, nil
	}

	// Calculate based on expiry type and cycle
	now := currentTime

	switch leg.ExpiryType {
	case "WEEKLY":
		return e.calculateWeeklyExpiry(now, leg.ExpiryCycle)
	case "MONTHLY":
		return e.calculateMonthlyExpiry(now, leg.ExpiryCycle)
	case "QUARTERLY":
		return e.calculateQuarterlyExpiry(now, leg.ExpiryCycle)
	default:
		return time.Time{}, fmt.Errorf("unsupported expiry type: %s", leg.ExpiryType)
	}
}

// calculateWeeklyExpiry calculates weekly expiry based on cycle
func (e *BacktestEngine) calculateWeeklyExpiry(currentTime time.Time, cycle strategyModel.ExpiryCycle) (time.Time, error) {
	// Find the next Thursday (weekly expiry day)
	daysUntilThursday := (int(time.Thursday) - int(currentTime.Weekday()) + 7) % 7
	if daysUntilThursday == 0 && currentTime.Hour() >= 15 { // After 3:30 PM on Thursday
		daysUntilThursday = 7
	}

	nextThursday := currentTime.AddDate(0, 0, daysUntilThursday)

	switch cycle {
	case "CURRENT":
		return nextThursday, nil
	case "NEXT":
		return nextThursday.AddDate(0, 0, 7), nil
	case "NEXT2":
		return nextThursday.AddDate(0, 0, 14), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported expiry cycle: %s", cycle)
	}
}

// calculateMonthlyExpiry calculates monthly expiry (last Thursday of the month)
func (e *BacktestEngine) calculateMonthlyExpiry(currentTime time.Time, cycle strategyModel.ExpiryCycle) (time.Time, error) {
	var targetMonth time.Time

	switch cycle {
	case "CURRENT":
		targetMonth = currentTime
	case "NEXT":
		targetMonth = currentTime.AddDate(0, 1, 0)
	case "NEXT2":
		targetMonth = currentTime.AddDate(0, 2, 0)
	default:
		return time.Time{}, fmt.Errorf("unsupported expiry cycle: %s", cycle)
	}

	// Find last Thursday of the target month
	lastDay := time.Date(targetMonth.Year(), targetMonth.Month()+1, 0, 15, 30, 0, 0, targetMonth.Location())
	daysBack := (int(lastDay.Weekday()) - int(time.Thursday) + 7) % 7
	lastThursday := lastDay.AddDate(0, 0, -daysBack)

	return lastThursday, nil
}

// calculateQuarterlyExpiry calculates quarterly expiry (last Thursday of March, June, September, December)
func (e *BacktestEngine) calculateQuarterlyExpiry(currentTime time.Time, cycle strategyModel.ExpiryCycle) (time.Time, error) {
	quarterMonths := []time.Month{time.March, time.June, time.September, time.December}
	currentQuarter := -1

	// Find current quarter
	for i, month := range quarterMonths {
		if currentTime.Month() <= month {
			currentQuarter = i
			break
		}
	}

	if currentQuarter == -1 {
		currentQuarter = 0 // December passed, next is March
	}

	var targetMonth time.Month
	var targetYear int = currentTime.Year()

	switch cycle {
	case "CURRENT":
		targetMonth = quarterMonths[currentQuarter]
		if currentTime.Month() > targetMonth {
			targetYear++
			targetMonth = quarterMonths[0] // March of next year
		}
	case "NEXT":
		nextQuarter := (currentQuarter + 1) % 4
		targetMonth = quarterMonths[nextQuarter]
		if nextQuarter <= currentQuarter {
			targetYear++
		}
	case "NEXT2":
		next2Quarter := (currentQuarter + 2) % 4
		targetMonth = quarterMonths[next2Quarter]
		if next2Quarter <= currentQuarter {
			targetYear++
		}
	default:
		return time.Time{}, fmt.Errorf("unsupported expiry cycle: %s", cycle)
	}

	// Find last Thursday of the target month
	lastDay := time.Date(targetYear, targetMonth+1, 0, 15, 30, 0, 0, currentTime.Location())
	daysBack := (int(lastDay.Weekday()) - int(time.Thursday) + 7) % 7
	lastThursday := lastDay.AddDate(0, 0, -daysBack)

	return lastThursday, nil
}

// ExitFnOPosition exits all positions in the active run at the specified time
func (e *BacktestEngine) ExitFnOPosition(
	activeRun *model.Run,
	currentTime time.Time,
	closedRuns *[]*model.Run,
) {
	// Exit all positions in the active run
	for _, position := range activeRun.Positions {
		if position.State == model.Entered {
			// Get current premium price for exit (using entry price for now, should be fetched from market data)
			exitPrice := *position.EntryPriceOrder // Simplified - should get current market premium
			exitFees := 0.0

			// Exit the position
			position.Exit(
				exitPrice,
				currentTime,
				exitPrice,
				currentTime,
				"TIME_BASED_EXIT",              // Exit signal
				string(model.ExitReasonSignal), // Exit reason
				exitFees,
				nil,
			)

			println("Exited position for leg at", currentTime.Format(time.RFC3339), "- P&L:", *position.NetProfitLoss)
		}
	}

	// Mark run as ended
	runEndTime := currentTime
	activeRun.EndedTime = &runEndTime
	activeRun.State = model.Ended

	// Skip metrics calculation for now

	// Add completed run to closedRuns list instead of persisting
	*closedRuns = append(*closedRuns, activeRun)
	println("Added completed run with ID:", activeRun.ID.String(), "to closed runs list")
	println("Total closed runs:", len(*closedRuns))
}

// resolveOptionContract resolves the option contract for a given leg with spot price and timestamp
func (e *BacktestEngine) resolveOptionContract(
	ctx context.Context,
	leg *strategyDto.OptionLegResponse,
	legSpotPrice float64,
	spotTimestamp time.Time,
) (*instrumentModel.InstrumentMaster, error) {
	// Determine strike offset (use TickOffset if available, otherwise default to 50)
	strikeOffset := 100.0 // Default value

	// if leg.InstrumentType == strategyModel.InstrumentType("FUTURES") {
	// 	// Logic to resolve futures contract
	// } else if leg.InstrumentType == strategyModel.InstrumentType("OPTIONS") {
	// 	// Logic to resolve options contract
	// }

	// Resolve the strike price based on strike type and option type
	strikePrice, err := e.resolveStrikePriceByType(
		string(leg.StrikeType),
		string(*leg.OptionType),
		legSpotPrice,
		strikeOffset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve strike price: %w", err)
	}

	// Resolve options contracts from master using the calculated strike price
	var parentSecurityID *int32
	if leg.ExchangeSecurityID != nil {
		id := int32(*leg.ExchangeSecurityID)
		parentSecurityID = &id
	}

	// Calculate the actual expiry date based on ExpiryType and ExpiryCycle
	expiryDate, err := e.calculateExpiryDateForLeg(leg, spotTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate expiry date: %w", err)
	}

	optionContract, err := e.instrumentService.ResolveOptionsContractsFromMaster(
		ctx,
		strikePrice,
		string(*leg.OptionType), // Convert OptionType enum to string
		expiryDate,              // Use calculated expiry date
		parentSecurityID,        // Parent security ID
	)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve option contract: %w", err)
	}

	return optionContract, nil
}
