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

func (engine *Engine) Run(ctx context.Context, run *shared.Run) *shared.StockIndicator {
	strategy := engine.strategyService.Get(ctx, run.StrategyID)

	entryDataRange := strategy.GetEntryDataRange()
	exitDataRange := strategy.GetExitDataRange()

	dataRange := util.Max(entryDataRange, exitDataRange)

	println("Fetching the data for the range ", dataRange)

	instrument := strategy.Instruments[0]

	data := engine.stockService.Get(instrument, "", dataRange)
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
		position := engine.GetLatestPosition(run.ID)
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
				position.EntryPrice = &stockIndicator.Close[i]
				position.EntryTime = &stockIndicator.Timestamp[i]
				// TODO: Make Order(Buy if LONG, else Sell)
			}

		case "ENTERED":
			// check if the exit expression is true
			exit_condition := (ema14[i] < ema50[i]) && (ema14[i-1] < ema50[i-1])

			// TODO: Check the Risk Metric Conditions(Stop Loss, Trailing Stop, etc.)
			// TODO: Check the Profit Conditions for Exit

			if exit_condition {
				position.State = shared.Exited
				position.ExitPrice = &stockIndicator.Close[i]
				position.ExitTime = &stockIndicator.Timestamp[i]
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
	// runMetrics := engine.CalculateRunMetrics(run)

	// Store the metrics
	// engine.service.SaveRunMetrics(runMetrics)

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
	position := engine.GetLatestPosition(runID)
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
) *shared.Position {
	position := engine.service.GetLatestPosition(id)
	if position == nil || position.IsDone() {
		// create a new position
		position = engine.service.CreatePosition(id)
	}
	return position
}

func (engine *Engine) CalculateRunMetrics(run *shared.Run) *shared.RunMetrics {
	// Get all positions for this run
	positions := engine.service.GetPositions(run.ID)

	println(len(positions))

	metrics := &shared.RunMetrics{
		ID:    util.NewUUID(),
		RunID: run.ID,
	}

	return metrics

	if len(positions) == 0 {
		return metrics
	}

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
	var totalCommission float64 = 0 // Assuming zero commission for now
	var totalSlippage float64 = 0   // Assuming zero slippage for now

	// Calculate metrics from positions
	for i, pos := range positions {
		if pos.EntryPrice == nil || pos.ExitPrice == nil || pos.EntryTime == nil || pos.ExitTime == nil {
			continue
		}

		// Position size and P&L calculations
		positionSize := *pos.EntryPrice
		if positionSize > maxPositionSize {
			maxPositionSize = positionSize
		}
		totalPositionSize += positionSize

		pnl := *pos.ExitPrice - *pos.EntryPrice
		totalPnL += pnl

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
			lossAmount += -pnl
		}

		// Trade duration
		duration := pos.ExitTime.Sub(*pos.EntryTime)
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
		initialCapital := positions[0].EntryPrice
		if initialCapital != nil && *initialCapital > 0 {
			metrics.TotalReturnPCT = (totalPnL / *initialCapital) * 100
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

		// Transaction costs
		metrics.TotalCommission = totalCommission
		metrics.TotalSlippage = totalSlippage
		if totalPositionSize > 0 {
			metrics.TransactionCostsPCT = (totalCommission + totalSlippage) / totalPositionSize * 100
		}
	}

	return metrics
}

/*
// TODO: resume execution

func (run *shared.Run) GetPosition() *shared.Position {
	if run.State == "START" {
		run.State = "IN_PROGRESS"
		// create a new position
		return position.NewPosition()
	} else if run.State == "IN_PROGRESS" {
		latest := position.GetLatest()
		if latest.State == "EXITED" {
			// create a new position
			return position.NewPosition()
		}
		return latest
	} else if run.State == "ENDED" {
		panic("Run ended, but loop going on")
	}
	panic("Error in run get position")
}
*/
