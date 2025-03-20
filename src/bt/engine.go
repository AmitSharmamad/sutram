package bt

import (
	"context"
	"fmt"

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
	// 1. create a run
	// 2. emit an event to asynchronously run the strategy back testing
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

		run := engine.service.CreateRun(&runObj)
		println(util.ToString(run.ID))

		// TODO: emit an event to asynchronously run the strategy back testing
	}

	return instruments

}

func (engine *Engine) Run(ctx context.Context, run *shared.Run) *shared.StockIndicator {
	strategy := engine.strategyService.Get(ctx, run.StrategyID)

	entryDataRange := strategy.GetEntryDataRange()
	exitDataRange := strategy.GetExitDataRange()

	dataRange := util.Max(entryDataRange, exitDataRange)

	println("Fetching the data for the range ", dataRange)

	instrument := strategy.Instruments[0]

	data := engine.stockService.Get(instrument, "", "", dataRange)
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
				// TODO: Make Order(Buy if LONG, else Sell)
			}

		case "ENTERED":
			// check if the exit expression is true
			exit_condition := (ema14[i] < ema50[i]) && (ema14[i-1] < ema50[i-1])

			// TODO: Check the Risk Metric Conditions(Stop Loss, Trailing Stop, etc.)
			// TODO: Check the Profit Conditions for Exit

			if exit_condition {
				position.State = shared.Exited
				position.ExitPriceSignal = &stockIndicator.Close[i]
				position.ExitTimeSignal = &stockIndicator.Timestamp[i]
				// TODO: Make the Exit Reason Dynamic
				exitReason := "HIT_PROFIT|HIT_STOP_LOSS|MANUAL"
				position.ExitReason = &exitReason

				// TODO: Make Order(Sell if LONG, else Buy)
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

	// Calculate the Run Metrics
	runMetrics := engine.CalculateRunMetrics(run)

	// Store the metrics
	engine.service.SaveRunMetrics(runMetrics)

	// Change the Run Status to Ended
	run.State = shared.Ended
	engine.service.UpdateRunState(run)

	return stockIndicator

	// response := map[string][]float64{}
	// response["close"] = closeData
	// response["ema10"] = indicators

	// return response

	// id := run.ID

	// run.State = "IN_PROGRESS"

	// engine.service.UpdateState(run)

	// start := 100
	// end := 200

	// counter := 0

	// println(start)
	// println(end)

	// for i := start; i < end; i++ {
	// 	println(start, end)
	// 	engine.HandlePosition(
	// 		id,
	// 		strategy,
	// 		(*data)[counter:start],
	// 	)
	// 	counter++
	// }

	// run.State = "ENDED"

	// engine.service.UpdateState(run)
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
