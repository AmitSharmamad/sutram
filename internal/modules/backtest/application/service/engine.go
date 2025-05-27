package service

import (
	"context"
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

// EngineService defines the service for backtest engine operations
type EngineService interface {
	// Start a backtest run for the provided strategy
	Start(ctx context.Context, strategyID uuid.UUID, userID uuid.UUID) ([]Response, error)
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

// calculateLookbackDuration calculates the duration to look back based on timeframe and number of candles
func (e *BacktestEngine) calculateLookbackDuration(timeframe string, candles int) time.Duration {
	var baseDuration time.Duration

	switch timeframe {
	case "1m":
		baseDuration = time.Minute
	case "3m":
		baseDuration = 3 * time.Minute
	case "5m":
		baseDuration = 5 * time.Minute
	case "15m":
		baseDuration = 15 * time.Minute
	case "30m":
		baseDuration = 30 * time.Minute
	case "1h":
		baseDuration = time.Hour
	case "3h":
		baseDuration = 3 * time.Hour
	case "1d":
		baseDuration = 24 * time.Hour
	default:
		// Default to 1 minute if unknown timeframe
		baseDuration = time.Minute
	}

	return time.Duration(candles) * baseDuration
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
	fetchStartTime := run.PeriodStartTime

	// Calculate lookback duration for 100 candles
	lookbackCandles := 100
	lookbackDuration := e.calculateLookbackDuration(string(strategy.ChartTimeframe), lookbackCandles)

	// Adjust fetch start time to include lookback period
	lookbackStartTime := fetchStartTime.Add(-lookbackDuration)

	// Get data for the main chart timeframe with lookback
	timeFrameDataStartTime := time.Now()
	println("\tFetching timeframe data with lookback for instrument", instrumentID, "- Start:", timeFrameDataStartTime.Format(time.RFC3339))
	println("\t\tLookback period:", lookbackCandles, "candles,", lookbackDuration.String())
	println("\t\tFetch range:", lookbackStartTime.Format(time.RFC3339), "to", run.PeriodEndTime.Format(time.RFC3339))

	timeFrameDataWithLookback, err := e.instrumentService.GetInstrumentsByTimeRange(
		ctx,
		instrumentMaster.ExchangeSecurityID,
		lookbackStartTime,
		run.PeriodEndTime,
		string(strategy.ChartTimeframe),
	)

	timeFrameDataEndTime := time.Now()
	timeFrameDataDuration := timeFrameDataEndTime.Sub(timeFrameDataStartTime)
	println("\tTimeframe data with lookback fetched for instrument", instrumentID, "- Duration:", timeFrameDataDuration.String(),
		"- Records:", len(timeFrameDataWithLookback))

	if err != nil {
		println("Error getting timeframe data for instrument ID:", instrumentID, "Error:", err.Error())
		return
	}

	// Use the data with lookback as our timeFrameData
	timeFrameData := timeFrameDataWithLookback

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
	initialCapital := capital

	// Mark run as started
	startTime := time.Now()
	if run != nil {
		run.StartedTime = &startTime
		run.State = model.Started
	}

	// Find the index where the actual backtest period starts
	// (after the lookback period)
	backtestStartIndex := 0
	if run != nil {
		for i, candle := range stockData {
			if candle.Timestamp.Equal(run.PeriodStartTime) || candle.Timestamp.After(run.PeriodStartTime) {
				backtestStartIndex = i
				break
			}
		}
	}

	// Minimum lookback period for indicators (use max of 20 or backtestStartIndex to ensure enough data)
	minLookback := 20
	if backtestStartIndex > minLookback {
		minLookback = backtestStartIndex
	}

	// Keep track of the current trading day for MIS trades
	var currentTradingDay time.Time
	if len(stockData) > 0 {
		currentTradingDay = stockData[0].Timestamp.Truncate(24 * time.Hour)
	}

	// Process each candle from the minimum lookback period to the end of the data
	for i := minLookback; i < len(stockData); i++ {
		// Safety check to avoid infinite loops
		if i >= len(stockData) {
			break
		}

		// Create data context using the data slice from the beginning up to the current position
		data := e.toDataContext(stockData[:i+1])

		// Check if we've moved to a new trading day
		dateTruncated := stockData[i].Timestamp.Truncate(24 * time.Hour)
		if !dateTruncated.Equal(currentTradingDay) {
			// If we have an open position and this is an MIS (intraday) trade, close it at the end of the previous day
			if openPosition != nil && openPosition.State == model.Entered && strategy.TradeType == strategyModel.TradeType("MIS") {
				// Find the last data point from the previous day - limit the search to avoid too much iteration
				var lastDataFromPreviousDay *instrumentModel.Instrument
				maxSearchDepth := min(100, i) // Limit search depth to avoid excessive looping
				for j := i - 1; j >= i-maxSearchDepth && j >= 0; j-- {
					if stockData[j].Timestamp.Truncate(24 * time.Hour).Equal(currentTradingDay) {
						lastDataFromPreviousDay = stockData[j]
						break
					}
				}

				if lastDataFromPreviousDay != nil {
					// Calculate fees
					exitFees := 0.0

					exitOrderPrice := lastDataFromPreviousDay.Close // For backtest, we can use the close price directly

					// Close the position with "EOD MIS Square-off" as the reason
					exitReason := "EOD MIS Square-off" // Special reason for MIS end-of-day square-off

					openPosition.Exit(
						lastDataFromPreviousDay.Close,
						lastDataFromPreviousDay.Timestamp,
						exitOrderPrice,
						lastDataFromPreviousDay.Timestamp,
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

					// Reset open position
					openPosition = nil
				}
			}
			// Update the current trading day
			currentTradingDay = dateTruncated
		}

		// Check entry condition if no position is open
		if openPosition == nil {
			// Check if the current time is within the entry time window
			// If EntryStartTime or EntryEndTime is specified, we only look for entries within that window
			canCheckEntryCondition := true

			if strategy.EntryStartTime != nil || strategy.EntryEndTime != nil {
				// Extract time of day from the current candle timestamp
				currentTime := stockData[i].Timestamp
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

				if canCheckEntryCondition && strategy.EntryEndTime != nil {
					// Extract time component from EntryEndTime
					entryEndTimeOfDay := time.Date(0, 1, 1,
						strategy.EntryEndTime.Hour(), strategy.EntryEndTime.Minute(),
						strategy.EntryEndTime.Second(), 0, time.UTC)

					if currentTimeOfDay.After(entryEndTimeOfDay) {
						canCheckEntryCondition = false
					}
				}
			}

			// Only check entry conditions if we're within the allowed entry time window
			if canCheckEntryCondition {
				entrySignalTriggered := e.evaluateEntryCondition(strategy.EntryExpression, data)

				if entrySignalTriggered {
					// Calculate position size based on strategy
					positionSize := e.calculatePositionSize(strategy, stockData[i].Close)

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
					entryFees := 0.0
					openPosition.Enter(entryPrice, stockData[i].Timestamp, entryFees, nil)

					openPosition.StopLossPrice = e.calculateStopLoss(
						*openPosition.EntryPriceOrder,
						strategy.StopLoss,
						strategy.StopLossType,
						strategy.PositionType,
					)

					// Update capital for position entry
					capital = e.updateCapitalOnPositionChange(capital, openPosition, entryPrice, entryFees, true)

				}
			}
		}

		// If position is still open, check exit condition using timeframe data
		if openPosition != nil {
			isExit := e.TryExit(openPosition, strategy, secondLevelData)

			// If the position was exited
			if isExit && openPosition.ExitTimeOrder != nil {
				// Collect position
				positions = append(positions, openPosition)

				// Add net profit/loss if available
				if openPosition.NetProfitLoss != nil {
					capital += *openPosition.NetProfitLoss
				}

				// Find the next index to resume entry checks based on the timeframe
				// This ensures we don't check for entry until the next candle in the chart timeframe
				nextIndex := e.findNextCandleIndex(stockData, *openPosition.ExitTimeOrder, string(strategy.ChartTimeframe))

				// Safety check to ensure we don't go out of bounds or get stuck
				if nextIndex <= i {
					nextIndex = i + 1 // Simply advance by one if we couldn't find a good next index
				}

				// Ensure we never go beyond the array bounds
				if nextIndex >= len(stockData) {
					nextIndex = len(stockData) - 1
				}

				i = nextIndex - 1 // -1 because the loop will increment i

				// Reset open position
				openPosition = nil

				// Skip to the next candle - this will be incremented by the for loop
				continue
			}
		}
	}

	// Handle any remaining open position at the end of backtest period
	if openPosition != nil && openPosition.State == model.Entered {
		// If this is a MIS position, close it at the end of the backtest
		if strategy.TradeType == strategyModel.TradeType("MIS") && len(stockData) > 0 {
			lastData := stockData[len(stockData)-1]

			// Calculate fees
			exitFees := 0.0

			exitOrderPrice := lastData.Close // For backtest, we can use the close price directly

			// Close the position with "EOD MIS Square-off" as the reason
			exitReason := "EOD MIS Square-off" // Special reason for MIS end-of-day square-off

			openPosition.Exit(
				lastData.Close,
				lastData.Timestamp,
				exitOrderPrice,
				lastData.Timestamp,
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

	// Calculate metrics
	metrics := e.calculateMetrics(positions, initialCapital, capital, run)

	// Calculate and log function execution time
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
	data *instrumentModel.Instrument,
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

	// Extract the time from the data timestamp
	currentTime := data.Timestamp

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

func (e *BacktestEngine) evaluateEntryCondition(expression string, data *instrumentModel.StockData) bool {
	return parser.Execute(expression, parser.NewEvaluationContext(data))
}

func (e *BacktestEngine) evaluateExitCondition(expression string, data *instrumentModel.StockData) bool {
	return parser.Execute(expression, parser.NewEvaluationContext(data))
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

func (e *BacktestEngine) updateStopLossBasedOnTrailing(
	position *model.Position,
	strategy *strategyDto.StrategyResponse,
	instrument *instrumentModel.Instrument,
) {
	// Skip if trailing stop loss isn't enabled or position doesn't have a stop loss
	if strategy.TrailingStopLoss <= 0 || position.StopLossPrice == nil || position.EntryPriceOrder == nil {
		return
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
			}
		}
	}
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
		return instrument.Low <= *position.StopLossPrice
	} else { // Short position
		return instrument.High >= *position.StopLossPrice
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
		// Clamp win rate to avoid PostgreSQL numeric field overflow
		winRate := float64(winningTrades) / float64(totalTrades) * 100
		metrics.WinRate = e.clampFloat(winRate, -9999.9999, 9999.9999)
	}

	// Calculate profit factor
	profitFactor := 0.0
	if totalLossAmount != 0 {
		profitFactor = totalWinAmount / math.Abs(totalLossAmount)
		// Clamp profit factor to fit in database DECIMAL(10, 4)
		profitFactor = e.clampFloat(profitFactor, -9999.9999, 9999.9999)
	}
	metrics.ProfitFactor = &profitFactor

	// Calculate average win and loss sizes
	avgWinSize := 0.0
	if winningTrades > 0 {
		avgWinSize = totalWinAmount / float64(winningTrades)
		// Clamp to fit in database DECIMAL(18, 8)
		avgWinSize = e.clampFloat(avgWinSize, -9999999999.99999999, 9999999999.99999999)
	}
	metrics.AvgWinSize = &avgWinSize

	avgLossSize := 0.0
	if losingTrades > 0 {
		avgLossSize = totalLossAmount / float64(losingTrades)
		// Clamp to fit in database DECIMAL(18, 8)
		avgLossSize = e.clampFloat(avgLossSize, -9999999999.99999999, 9999999999.99999999)
	}
	metrics.AvgLossSize = &avgLossSize

	// Set transaction costs
	metrics.TotalCommission = e.clampFloat(totalCommission, -9999999999.99999999, 9999999999.99999999)
	metrics.TotalSlippage = e.clampFloat(totalSlippage, -9999999999.99999999, 9999999999.99999999)

	if initialCapital > 0 {
		transactionCostsPct := (totalCommission + totalSlippage) / initialCapital * 100
		metrics.TransactionCostsPct = e.clampFloat(transactionCostsPct, -9999.9999, 9999.9999)

		totalReturnPct := (finalCapital - initialCapital) / initialCapital * 100
		metrics.TotalReturnPct = e.clampFloat(totalReturnPct, -9999.9999, 9999.9999)
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
		// Clamp to fit in database DECIMAL(10, 4)
		avgTradeDuration = e.clampFloat(avgTradeDuration, -9999.9999, 9999.9999)
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

// clampFloat restricts a float value to be within the specified range
func (e *BacktestEngine) clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
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
		stopLossHit := e.isStopLossTriggered(position, data)
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
			exitFees := 0.0

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

		// Check exit conditions
		exitSignalTriggered := false
		if strategy.ExitExpression != "" {
			exitSignalTriggered = e.evaluateExitCondition(strategy.ExitExpression, dataContext)
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
			exitReason := "Exit signal"
			if stopLossHit {
				exitReason = "Stop loss"
			} else if targetProfitHit {
				exitReason = "Target profit"
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
