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
			*strategy.TrailingStopLossType,
		)
	}

	ctx := parser.NewEvaluationContext(data)

	exitConditionResult := false

	if strategy.ExitExpression != "" {
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
	var totalCommission float64 = 0
	var totalSlippage float64 = 0

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

		// Calculate position P&L based on direction
		var pnl float64
		if pos.Direction == shared.LONG {
			pnl = *pos.ExitPriceSignal - *pos.EntryPriceSignal
		} else if pos.Direction == shared.SHORT {
			pnl = *pos.EntryPriceSignal - *pos.ExitPriceSignal
		}
		totalPnL += pnl

		// Calculate daily return for this position (direction-aware)
		var dailyReturn float64
		if pos.Direction == shared.LONG {
			dailyReturn = pnl / positionSize
		} else {
			// For short positions, positive return means price went down
			dailyReturn = -pnl / positionSize
		}
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

		// Win/Loss statistics (direction-aware)
		if (pos.Direction == shared.LONG && pnl > 0) || (pos.Direction == shared.SHORT && pnl < 0) {
			totalWins++
			winAmount += util.Abs(pnl)
		} else if (pos.Direction == shared.LONG && pnl < 0) || (pos.Direction == shared.SHORT && pnl > 0) {
			totalLosses++
			lossAmount += util.Abs(pnl)
		}

		// Trade duration
		duration := pos.ExitTimeSignal.Sub(*pos.EntryTimeSignal)
		totalTradeDuration += duration.Hours()
	}

	// ... rest of the existing metrics calculation ...
	totalTrades := int64(len(positions))
	if totalTrades > 0 {
		metrics.TotalTrades = totalTrades
		metrics.WinRate = float64(totalWins) / float64(totalTrades)
		metrics.AvgTradeDuration = totalTradeDuration / float64(totalTrades)
		metrics.MaxPositionSize = maxPositionSize
		metrics.AvgPositionSize = totalPositionSize / float64(totalTrades)

		// P&L and return metrics
		initialCapital := positions[0].EntryPriceSignal
		if initialCapital != nil && *initialCapital > 0 {
			metrics.TotalReturnPCT = (totalPnL / *initialCapital) * 100

			totalDuration := run.EndedTime.Sub(*run.StartedTime).Hours() / 24
			if totalDuration > 0 {
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
			metrics.MaxDrawdownDuration = maxDrawdown / peakValue * 100
		}

		// Standard deviation and ratios
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

			if metrics.StdDeviation > 0 {
				metrics.SharpeRatio = mean / metrics.StdDeviation * util.Sqrt(252)
			}

			downsideVariance := 0.0
			for _, ret := range dailyReturns {
				if ret < 0 {
					downsideVariance += ret * ret
				}
			}
			downsideVariance /= float64(len(dailyReturns))
			metrics.DownsideDeviation = util.Sqrt(downsideVariance)

			if metrics.DownsideDeviation > 0 {
				metrics.SortinoRatio = mean / metrics.DownsideDeviation * util.Sqrt(252)
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
