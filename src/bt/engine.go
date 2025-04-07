package bt

import (
	"context"
	"sync"

	"algo/src/parser"
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
	expressionService shared.ExpressionService,
	indicatorService shared.IndicatorService,
) *Engine {
	return &Engine{
		service:           service,
		strategyService:   strategyService,
		stockService:      stockService,
		expressionService: expressionService,
		indicatorService:  indicatorService,
	}
}

func (engine *Engine) Start(ctx context.Context, strategyID uuid.UUID) *shared.RunGroup {
	strategy := engine.strategyService.Get(ctx, strategyID)

	runGroup := shared.NewRunGroup(strategyID)

	instruments := strategy.Instruments

	var wg sync.WaitGroup

	runs := make([]*shared.Run, 0)

	for _, instrumentID := range instruments {
		run := shared.NewRun(
			runGroup.ID,
			strategyID,
			instrumentID,
		)

		runs = append(runs, run)

		wg.Add(1)

		// Execute run asynchronously as goroutine
		go engine.Run(ctx, run, strategy, &wg)
	}

	wg.Wait()

	runGroup.Runs = runs

	// TODO: Move all the below DB operations to a separate goroutine function

	engine.service.runGroupRepo.Create(runGroup)

	engine.service.runRepo.Create(runs)

	positions := make([]*shared.Position, 0)
	runMetrics := make([]*shared.RunMetrics, 0)
	for _, run := range runs {
		runMetrics = append(runMetrics, run.RunMetrics)
		positions = append(positions, run.Positions...)
	}

	engine.service.positionRepo.Create(positions)

	err := engine.service.runMetricsRepo.Upsert(runMetrics)
	println(err)

	return runGroup
}

func (engine *Engine) Run(
	ctx context.Context,
	run *shared.Run,
	strategy *shared.Strategy,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	startTime := util.GetTime()
	run.StartedTime = &startTime

	// TODO: fix this
	// entryLookBack := engine.expressionService.GetEntryDataRange(strategy)
	// exitLookBack := engine.expressionService.GetExitDataRange(strategy)
	// lookBack := util.Max(entryLookBack, exitLookBack)

	// TODO: fix this for GA, leave it like this for now
	var lookBack int32 = 1000

	data := engine.stockService.Get(
		run.InstrumentID,
		strategy.StartPeriod,
		strategy.EndPeriod,
		strategy.ChartTimeframe,
		lookBack,
	)

	for i := lookBack; i < int32(len(*data)); i++ {
		current := (*data)[i]
		position := engine.GetLatestPosition(run, strategy)
		slicedData := (*data)[:i+1]

		switch position.State {
		case shared.NotEntered:
			engine.TryEnter(position, strategy, &current, &slicedData)
		case shared.Entered:
			engine.TryExit(position, strategy, &current, &slicedData)
		default:
			panic("Error in position Eval")
		}
	}

	run.State = shared.Ended
	endTime := util.GetTime()
	run.EndedTime = &endTime

	run.RunMetrics = engine.CalculateRunMetrics(run)
}

func (engine *Engine) TryEnter(
	position *shared.Position,
	strategy *shared.Strategy,
	current *shared.Stock,
	data *shared.StockData,
) {
	ctx := parser.NewEvaluationContext(data)

	entryCondition := parser.Execute(strategy.EntryExpression, ctx)

	if entryCondition {

		println("entering")

		position.HighestPrice = &current.Close
		position.LowestPrice = &current.Close

		position.UpdateOnEntrySignal(
			&current.Close,
			&current.Timestamp,
		)

		position.UpdateOnEntryOrder(
			&current.Close,
			&current.Timestamp,
		)
	}
}

func (engine *Engine) TryExit(
	position *shared.Position,
	strategy *shared.Strategy,
	current *shared.Stock,
	data *shared.StockData,
) {

	stopLossPrice := position.CalculateStopLossPrice(
		strategy.StopLoss,
		strategy.StopLossType,
	)

	targetProfitPrice := position.CalculateTargetProfitPrice(
		strategy.TargetProfit,
		strategy.TargetProfitType,
	)

	var trailingStopLossPrice float64
	if strategy.TrailingStopLoss != nil {
		trailingStopLossPrice = position.CalculateTrailingStopLossPrice(
			current.Close,
			strategy.TrailingStopLoss,
			strategy.TrailingStopLossType,
		)
	}

	ctx := parser.NewEvaluationContext(data)

	exitConditionResult := false

	if (strategy.ExitExpression != "") {
		exitConditionResult = parser.Execute(strategy.ExitExpression, ctx)
	}

	isExit := position.CheckForExit(
		current.Close,
		stopLossPrice,
		targetProfitPrice,
		trailingStopLossPrice,
		exitConditionResult,
	)

	if isExit {

		println("exiting")

		position.UpdateOnExitSignal(
			&current.Close,
			&current.Timestamp,
		)

		// Simulate Order
		position.UpdateOnExitOrder(
			&current.Close,
			&current.Timestamp,
		)
	}
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
