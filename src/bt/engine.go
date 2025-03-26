package bt

import (
	"context"
	"fmt"
	"sync"
	"time"

	"algo/src/shared"
	"algo/src/util"

	"github.com/google/uuid"
)

type Engine struct {
	service           *Service
	strategyService   shared.StrategyService
	stockService      shared.StockService
	expressionService shared.ExpressionService
	indicatorService  shared.IndicatorService
}

func NewEngineProv(
	service *Service,
	strategyService shared.StrategyService,
	stockService shared.StockService,
	indicatorService shared.IndicatorService,
) *Engine {
	return &Engine{
		service:          service,
		strategyService:  strategyService,
		stockService:     stockService,
		indicatorService: indicatorService,
	}
}

func (engine *Engine) Start(ctx context.Context, strategyID uuid.UUID) []string {

	strategy := engine.strategyService.Get(ctx, strategyID)

	// get the instruments for the strategy
	instruments := strategy.Instruments

	// for each of the instruments in the strategy:
	// 1. Create a run
	// 2. Asynchronously run the run

	var wg sync.WaitGroup

	runs := []shared.Run{}
	for _, instrument := range instruments {
		runObj := shared.Run{}

		runObj.Instrument = instrument
		runObj.StrategyID = strategyID
		runObj.RunGroupID = nil
		runObj.Exchange = "MCX"
		runObj.RunType = shared.BackTest
		runObj.State = shared.NotStarted
		runObj.Qty = 10
		runObj.PeriodStartTime = nil
		runObj.PeriodEndTime = nil
		runObj.StartedTime = nil
		runObj.EndedTime = nil

		runs = append(runs, runObj)
		run := engine.service.CreateRun(&runObj)

		wg.Add(1)

		// Execute run asynchronously as goroutine
		go engine.Run(ctx, run, &wg)
	}

	wg.Wait()

	println(runs)

	return instruments

}

func (engine *Engine) Run(ctx context.Context, run *shared.Run, wg *sync.WaitGroup) *shared.StockIndicator {
	defer wg.Done()

	strategy := engine.strategyService.Get(ctx, run.StrategyID)

	// Update the Run Start Time
	startTime := util.GetTime()
	run.StartedTime = &startTime

	// entryDataRange := strategy.GetEntryDataRange()
	// exitDataRange := strategy.GetExitDataRange()

	// dataRange := util.Max(entryDataRange, exitDataRange)

	// println("Fetching the data for the range ", dataRange)

	instrument := run.Instrument

	data := engine.stockService.Get(instrument, "", "", 1000) // Remove hardcoded values
	stockDataColumn := engine.stockService.Converter(*data)
	println(stockDataColumn.Close)
	closeData := stockDataColumn.Close
	println(closeData[0])
	ema14 := engine.indicatorService.GetIndicatorData(ctx, "EMA", 14, data)
	ema50 := engine.indicatorService.GetIndicatorData(ctx, "EMA", 50, data)

	stockIndicator := &shared.StockIndicator{
		StockColumn: *stockDataColumn,
		Indicators: map[string][]float64{
			"ema14": ema14,
			"ema50": ema50,
		},
	}

	for i := 50; i < len(stockIndicator.Close); i++ {
		position := engine.GetLatestPosition(run.ID, strategy)
		currPos := position.State
		switch currPos {
		case "NOT_ENTERED":
			// LONG POSITION

			// check if the entry expression is true
			entry_condition := (ema14[i] > ema50[i]) && (ema14[i-1] < ema50[i-1])
			if entry_condition {
				fmt.Println("test", stockIndicator.Timestamp[i].Format("2006-01-02 15:04:05"))
				fmt.Printf("close -> %.6f", stockIndicator.Close[i])
				fmt.Printf("ema14[i] -> %.6f ema50[i] -> %.6f\n", ema14[i], ema50[i])
				position.State = shared.Entered
				position.EntryPriceSignal = &stockIndicator.Close[i]
				position.EntryTimeSignal = &stockIndicator.Timestamp[i]

				// NOTE: Simulate Order - Remove once the order is implemented
				position.EntryPriceOrder = &stockIndicator.Close[i]
				position.EntryTimeOrder = &stockIndicator.Timestamp[i]

				// TODO: Make Order(Buy if LONG, else Sell)
			}

		case "ENTERED":
			// Exit logic using if-else conditions to ensure mutual exclusivity

			// Check all exit conditions in priority order
			currentPrice := stockIndicator.Close[i]
			var stopLossPrice float64
			var targetProfitPrice float64

			// NOTE: Remove hardcoded values and replace 'EntryPriceSignal' with 'EntryPriceOrder'
			stopLossPrice = *position.EntryPriceSignal * (1 - 0.01)
			targetProfitPrice = *position.EntryPriceSignal * (1 + 0.02)
			exit_condition := (ema14[i] < ema50[i]) && (ema14[i-1] < ema50[i-1])

			// Apply exit logic in priority order with if-else
			if currentPrice < stopLossPrice {
				// Stop Loss (highest priority)
				position = engine.UpdatePositionOnExitSignal(
					position,
					&stockIndicator.Close[i],
					&stockIndicator.Timestamp[i],
					shared.HIT_STOP_LOSS,
				)
			} else if currentPrice >= targetProfitPrice {
				// Target Profit (second priority)
				position = engine.UpdatePositionOnExitSignal(
					position,
					&stockIndicator.Close[i],
					&stockIndicator.Timestamp[i],
					shared.HIT_TAKE_PROFIT,
				)
			} else if exit_condition {
				// Exit Condition (lowest priority)
				position = engine.UpdatePositionOnExitSignal(
					position,
					&stockIndicator.Close[i],
					&stockIndicator.Timestamp[i],
					shared.EXIT_CONDITION_MET,
				)
			}

			// TODO: if condition matches, Make Order(Sell if LONG, else Buy)
			if position.State == shared.Exited {
				// Simulate Order

				// NOTE: Simulate Order - Remove once the order is implemented
				position = engine.UpdatePositionOnExitOrder(
					position,
					&stockIndicator.Close[i],
					&stockIndicator.Timestamp[i],
				)
			}
		default:
			panic("Error in position Eval")
		}

		// Update the Position
		updatedPos := position.State
		if currPos != updatedPos {
			engine.service.UpdatePosition(position)
		}
	}

	// Update the Run State and End Time
	run.State = shared.Ended
	endTime := util.GetTime()
	run.EndedTime = &endTime

	// Finally Update the changes in db
	engine.service.UpdateRunState(run)

	// Calculate the Run Metrics
	runMetrics := engine.CalculateRunMetrics(run)

	// Store the metrics
	engine.service.SaveRunMetrics(runMetrics)

	// Change the Run Status to Ended
	run.State = shared.Ended

	return stockIndicator
}

func (engine *Engine) UpdatePositionOnExitOrder(
	position *shared.Position,
	exitPriceOrder *float64,
	exitTimeStampOrder *time.Time,
) *shared.Position {

	position.ExitPriceOrder = exitPriceOrder
	position.ExitTimeOrder = exitTimeStampOrder

	// Calculate P&L - Absolute and Percentage
	var profitlossAbsolute, profitLossPercentage float64
	profitlossAbsolute, profitLossPercentage = engine.CalculateProfitLoss(
		position.Direction,
		*position.EntryPriceOrder,
		*position.ExitPriceOrder,
	)
	position.ProfitLossAbsolute = &profitlossAbsolute
	position.ProfitLossPercentage = &profitLossPercentage

	// Calculate Entry Fees and Exit Fees
	var fees float64 = 20.0 // Assuming fixed entry fees for now (can be updated later)
	position.EntryFees = &fees
	position.ExitFees = &fees

	// Calculate Transaction Costs
	transactionCost := *position.EntryFees + *position.ExitFees
	position.TransactionCost = &transactionCost

	netProfitLoss := *position.ProfitLossAbsolute - *position.TransactionCost
	position.NetProfitLoss = &netProfitLoss

	return position
}

func (engine *Engine) UpdatePositionOnExitSignal(
	position *shared.Position,
	exitPrice *float64,
	exitTimeStamp *time.Time,
	exitReason shared.ExitReason,
) *shared.Position {

	position.State = shared.Exited
	position.ExitPriceSignal = exitPrice
	position.ExitTimeSignal = exitTimeStamp
	position.ExitReason = &exitReason

	return position

}

func (engine *Engine) CalculateProfitLoss(
	direction shared.Direction,
	entryPrice float64,
	exitPrice float64,
) (float64, float64) {

	var profitLossAbsolute float64
	var profitLossPercentage float64

	if direction == shared.LONG {
		profitLossAbsolute = exitPrice - entryPrice
	} else {
		profitLossAbsolute = entryPrice - exitPrice
	}

	// Calculate Profit/Loss Percentage
	profitLossPercentage = (profitLossAbsolute / entryPrice) * 100

	return profitLossAbsolute, profitLossPercentage
}

func (engine *Engine) HandlePositionIndicator(
	runID uuid.UUID,
	strategy *shared.Strategy,
	data shared.StockIndicator,
) {

}

func (engine *Engine) HandlePosition(
	runID uuid.UUID,
	strategy *shared.Strategy,
	data shared.StockData,
) {
	position := engine.GetLatestPosition(runID, strategy)
	currPos := position.State
	switch position.State {
	case shared.NotEntered:
		if engine.expressionService.Eval(strategy.EntryExpression, data) {
			position.State = shared.Entered
		}
	case shared.Entered:
		if engine.expressionService.Eval(strategy.ExitExpression, data) {
			position.State = shared.Exited
		}
	default:
		panic("Error in position Eval")
	}

	updatedPos := position.State
	engine.service.CreateLog(shared.NewPositionLog(
		position.ID,
		position.State,
	))
	if currPos != updatedPos {
		engine.service.UpdatePosition(position)
	}
}

func (engine *Engine) GetLatestPosition(
	id uuid.UUID,
	strategy *shared.Strategy,
) *shared.Position {
	position := engine.service.GetLatestPosition(id)
	if position == nil || position.IsDone() {
		// NOTE: direction is position is equivalent to strategy.PositionType
		direction := strategy.PositionType
		// create a new position
		position = engine.service.CreatePosition(id, direction)
	}
	return position
}

func (engine *Engine) CalculateRunMetrics(run *shared.Run) *shared.RunMetrics {
	// Get all positions for this run
	positions := engine.service.GetPositions(run.ID)

	metrics := &shared.RunMetrics{
		ID:    util.NewUUID(),
		RunID: run.ID,
	}

	if len(positions) == 0 {
		return metrics
	}

	// Initialize metric variables
	var totalPnL float64
	var totalWins int64
	var totalLosses int64
	var totalTradeDuration float64
	var winAmount float64
	var lossAmount float64
	var maxDrawdown float64
	var currentDrawdown float64
	var peakValue float64
	var maxPositionSize float64
	var totalPositionSize float64
	var totalCommission float64 = 0 // Assuming zero commission for now (can be updated later)
	var totalSlippage float64 = 0   // Assuming zero slippage for now (can be updated later)

	// Daily returns for standard deviation calculation
	dailyReturns := []float64{}

	// Calculate metrics from positions
	for i, pos := range positions {
		// Skip incomplete positions
		if pos.EntryPriceSignal == nil || pos.ExitPriceSignal == nil || pos.EntryTimeSignal == nil || pos.ExitTimeSignal == nil {
			continue
		}

		// Position size and P&L calculations
		positionSize := *pos.EntryPriceSignal
		if positionSize > maxPositionSize {
			maxPositionSize = positionSize
		}
		totalPositionSize += positionSize

		// Calculate position P&L
		pnl := *pos.ExitPriceSignal - *pos.EntryPriceSignal
		totalPnL += pnl

		// Calculate daily return for this position (simplified)
		dailyReturn := pnl / positionSize
		dailyReturns = append(dailyReturns, dailyReturn)

		// Update peak value and drawdown
		currentValue := totalPnL
		if i == 0 || currentValue > peakValue {
			peakValue = currentValue
		}
		currentDrawdown = peakValue - currentValue
		if currentDrawdown > maxDrawdown {
			maxDrawdown = currentDrawdown
		}

		// Win/Loss statistics
		if pnl > 0 {
			totalWins++
			winAmount += pnl
		} else if pnl < 0 {
			totalLosses++
			lossAmount += -pnl // Make loss amount positive for calculations
		}

		// Trade duration
		duration := pos.ExitTimeSignal.Sub(*pos.EntryTimeSignal)
		totalTradeDuration += duration.Hours()
	}

	totalTrades := int64(len(positions))
	if totalTrades > 0 {
		// Basic metrics
		metrics.TotalTrades = totalTrades
		metrics.WinRate = float64(totalWins) / float64(totalTrades)
		metrics.AvgTradeDuration = totalTradeDuration / float64(totalTrades)
		metrics.MaxPositionSize = maxPositionSize
		metrics.AvgPositionSize = totalPositionSize / float64(totalTrades)

		// P&L and return metrics
		// If we have initial capital information from the strategy or run, use it
		// For now, we'll use the first position's entry price as a proxy for initial capital
		initialCapital := positions[0].EntryPriceSignal
		if initialCapital != nil && *initialCapital > 0 {
			metrics.TotalReturnPCT = (totalPnL / *initialCapital) * 100

			// Calculate annualized return (simplified)
			totalDuration := run.EndedTime.Sub(*run.StartedTime).Hours() / 24 // in days
			if totalDuration > 0 {
				// Annualized return = (1 + total_return)^(365/days) - 1
				annualFactor := 365.0 / totalDuration
				metrics.AnnualizedReturn = (util.Pow(1+totalPnL/(*initialCapital), annualFactor) - 1) * 100
			}
		}

		// Win/Loss metrics
		if totalWins > 0 {
			metrics.AvgWinSize = winAmount / float64(totalWins)
		}
		if totalLosses > 0 {
			metrics.AvgLossSize = lossAmount / float64(totalLosses)
		}
		if totalLosses > 0 && lossAmount > 0 {
			metrics.ProfitFactor = winAmount / lossAmount
		}

		// Risk metrics
		metrics.MaxDrawdown = maxDrawdown
		if peakValue > 0 {
			metrics.MaxDrawdownDuration = maxDrawdown / peakValue * 100 // As percentage
		}

		// Standard deviation calculation
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
			metrics.StdDeviation = util.Sqrt(variance)

			// Calculate Sharpe Ratio (simplified, assuming risk-free rate of 0)
			if metrics.StdDeviation > 0 {
				metrics.SharpeRatio = mean / metrics.StdDeviation * util.Sqrt(252) // Annualized
			}

			// Calculate downside deviation (for Sortino ratio)
			downsideVariance := 0.0
			for _, ret := range dailyReturns {
				if ret < 0 {
					downsideVariance += ret * ret
				}
			}
			downsideVariance /= float64(len(dailyReturns))
			metrics.DownsideDeviation = util.Sqrt(downsideVariance)

			// Calculate Sortino Ratio
			if metrics.DownsideDeviation > 0 {
				metrics.SortinoRatio = mean / metrics.DownsideDeviation * util.Sqrt(252) // Annualized
			}
		}

		// Transaction costs
		metrics.TotalCommission = totalCommission
		metrics.TotalSlippage = totalSlippage
		if totalPositionSize > 0 {
			metrics.TransactionCostsPCT = (totalCommission + totalSlippage) / totalPositionSize * 100
		}
	}

	return metrics
}
