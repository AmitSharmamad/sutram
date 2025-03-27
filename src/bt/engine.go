package bt

import (
	"context"
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

func (engine *Engine) Start(ctx context.Context, strategyID uuid.UUID) *shared.RunGroup {

	strategy := engine.strategyService.Get(ctx, strategyID)

	runGroup := shared.NewRunGroup(strategyID)

	instruments := strategy.Instruments

	var wg sync.WaitGroup

	runs := make([]*shared.Run, 0)

	for _, instrument := range instruments {
		run := shared.NewRun(
			runGroup.ID,
			strategyID,
			instrument,
		)

		runs = append(runs, run)

		wg.Add(1)

		// Execute run asynchronously as goroutine
		go engine.Run(ctx, run, strategy, &wg)
	}

	wg.Wait()

	engine.service.runGroupRepo.Create(runGroup)

	engine.service.runRepo.Create(runs)

	positions := make([]*shared.Position, 0)
	runMetrics := make([]*shared.RunMetrics, 0)
	for _, run := range runs {
		runMetrics = append(runMetrics, run.RunMetrics)

		positions = append(positions, run.Positions...)
	}
	engine.service.positionRepo.Create(positions)

	engine.service.runMetricsRepo.Create(runMetrics)

	runGroup.Runs = runs

	return runGroup
}

func (engine *Engine) Run(
	ctx context.Context, run *shared.Run, strategy *shared.Strategy, wg *sync.WaitGroup,
) {
	defer wg.Done()

	// Fetch the strategy
	// Fetch the datarange
	// Fetch the instrument data based off of the datarange
	// Iterate through the data. Start taking positions

	// Update the Run Start Time
	startTime := util.GetTime()
	run.StartedTime = &startTime

	// TODO: Fix this
	entryLookBack := strategy.GetEntryDataRange()
	exitLookBack := strategy.GetExitDataRange()

	lookBack := util.Max(entryLookBack, exitLookBack)

	println("Fetching the data for the range ", lookBack)

	// instrument := run.Instrument

	// TODO: Fetch instrment id
	data := engine.stockService.Get(
		0, strategy.StartPeriod, strategy.EndPeriod, strategy.ChartTimeframe, lookBack,
	)

	for i := lookBack; i < int32(len(*data)); i++ {
		current := (*data)[i]
		position := engine.GetLatestPosition(run, strategy)

		switch position.State {
		case shared.NotEntered:
			// TODO:
			// check if the entry expression is true
			entry_condition := true
			if entry_condition {

				position = engine.UpdatePositionOnEntrySignal(
					position,
					&current.Close,
					&current.Timestamp,
				)

				engine.UpdatePositionOnEntryOrder(
					position,
					&current.Close,
					&current.Timestamp,
				)
			}

		case shared.Entered:
			var stopLossPrice float64
			var targetProfitPrice float64

			// TODO: Handle Absolute and Percentage for stoploss
			stopLossPrice = *position.EntryPriceSignal * (1 - strategy.StopLoss)
			targetProfitPrice = *position.EntryPriceSignal * (1 + strategy.TargetProfit)
			exit_condition := true

			// TODO: Fix this
			if current.Close < stopLossPrice {
				// Stop Loss (highest priority)
				position = engine.UpdatePositionOnExitSignal(
					position,
					&current.Close,
					&current.Timestamp,
					shared.HIT_STOP_LOSS,
				)
			} else if current.Close >= targetProfitPrice {
				// Target Profit (second priority)
				position = engine.UpdatePositionOnExitSignal(
					position,
					&current.Close,
					&current.Timestamp,
					shared.HIT_TAKE_PROFIT,
				)
			} else if exit_condition {
				// Exit Condition (lowest priority)
				position = engine.UpdatePositionOnExitSignal(
					position,
					&current.Close,
					&current.Timestamp,
					shared.EXIT_CONDITION_MET,
				)
			}

			if position.IsDone() {
				// TODO: if condition matches, Make Order(Sell if LONG, else Buy)

				position = engine.UpdatePositionOnExitOrder(
					position,
					&current.Close,
					&current.Timestamp,
				)
			}
		default:
			panic("Error in position Eval")
		}
	}

	// Update the Run State and End Time
	run.State = shared.Ended
	endTime := util.GetTime()
	run.EndedTime = &endTime

	// Calculate the Run Metrics
	run.RunMetrics = engine.CalculateRunMetrics(run)

	run.State = shared.Ended
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
	exitPriceSignal *float64,
	exitTimeStampSignal *time.Time,
	exitReason shared.ExitReason,
) *shared.Position {

	position.State = shared.Exited
	position.ExitPriceSignal = exitPriceSignal
	position.ExitTimeSignal = exitTimeStampSignal
	position.ExitReason = &exitReason

	return position

}

func (engine *Engine) UpdatePositionOnEntrySignal(
	position *shared.Position,
	entryPriceSignal *float64,
	entryTimeStampSignal *time.Time,
) *shared.Position {
	position.State = shared.Entered
	position.EntryPriceSignal = entryPriceSignal
	position.EntryTimeSignal = entryTimeStampSignal

	return position
}

func (engine *Engine) UpdatePositionOnEntryOrder(
	position *shared.Position,
	entryPriceOrder *float64,
	entryTimeStampOrder *time.Time,
) *shared.Position {
	position.State = shared.Entered
	position.EntryPriceOrder = entryPriceOrder
	position.EntryTimeOrder = entryTimeStampOrder

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

func (engine *Engine) GetLatestPosition(
	run *shared.Run,
	strategy *shared.Strategy,
) *shared.Position {
	position := shared.NewPosition(
		run.ID,
		strategy.PositionType,
	)

	length := len(run.Positions)

	if length != 0 {
		latest := run.Positions[length-1]
		if !latest.IsDone() {
			return latest
		}
	}

	run.Positions = append(run.Positions, position)
	return position
}

func (engine *Engine) CalculateRunMetrics(run *shared.Run) *shared.RunMetrics {

	// TODO: Check this logic

	positions := run.Positions

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
