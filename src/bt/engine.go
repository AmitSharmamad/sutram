package bt

import (
	"context"
	"sync"
	"time"

	"algo/src/broker"
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
	broker            broker.LiveTradingBroker
}

func NewEngineProv(
	service *Service,
	strategyService shared.StrategyService,
	stockService shared.StockService,
	expressionService shared.ExpressionService,
	indicatorService shared.IndicatorService,
	broker broker.LiveTradingBroker,
) *Engine {
	return &Engine{
		service:           service,
		strategyService:   strategyService,
		stockService:      stockService,
		expressionService: expressionService,
		indicatorService:  indicatorService,
		broker:            broker,
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

	// Calculate lookback period for both entry and exit conditions
	// TODO: fix this
	// entryLookBack := engine.expressionService.GetEntryDataRange(strategy)
	// exitLookBack := engine.expressionService.GetExitDataRange(strategy)
	// lookBack := util.Max(entryLookBack, exitLookBack)

	// TODO: fix this for GA, leave it like this for now
	var lookBack int32 = 13

	// Get timeframe data for entry conditions
	timeframeData := engine.stockService.Get(
		run.InstrumentID,
		strategy.StartPeriod,
		strategy.EndPeriod,
		strategy.ChartTimeframe,
		lookBack,
	)

	// Get 1-second data for more granular exit condition checking
	secondData := engine.stockService.Get(
		run.InstrumentID,
		strategy.StartPeriod,
		strategy.EndPeriod,
		string(shared.StockGranularity1s), // Use the constant for 1-second data
		lookBack,
	)

	if len(*timeframeData) == 0 {
		// Handle case where data is not available
		run.State = shared.Ended
		endTime := util.GetTime()
		run.EndedTime = &endTime
		return
	}

	// Main loop through timeframe data for entry signals
	for i := lookBack; i < int32(len(*timeframeData)); i++ {
		currentTimeframe := (*timeframeData)[i]
		position := engine.GetLatestPosition(run, strategy)
		slicedTimeframeData := (*timeframeData)[:i+1]

		if position.State == shared.NotEntered {
			// Check entry using timeframe data
			entrySignal := engine.CheckEntryCondition(position, strategy, &currentTimeframe, &slicedTimeframeData)

			if entrySignal && len(*secondData) > 0 {
				// Determine the end time for exit condition checking
				var nextTimeframeBoundary time.Time
				if i < int32(len(*timeframeData))-1 {
					nextTimeframeBoundary = (*timeframeData)[i+1].Timestamp
				} else {
					nextTimeframeBoundary = strategy.EndPeriod.Time
				}

				// Position was entered, now need to check exit using second data
				engine.ProcessExitWithSecondData(position, strategy, secondData, currentTimeframe.Timestamp, nextTimeframeBoundary)

				// If position exited, we can optionally optimize by advancing the counter
				// to the next timeframe after the exit time
				if position.State == shared.Exited && position.ExitTimeOrder != nil {
					// Find the next possible entry point in timeframe data after this exit
					exitTime := *position.ExitTimeOrder
					for j := i + 1; j < int32(len(*timeframeData)); j++ {
						// If this timeframe data point is after the exit time, set counter to it
						if (*timeframeData)[j].Timestamp.After(exitTime) {
							i = j - 1 // -1 because the loop will increment i
							break
						}
					}
				}
			}
		}
	}

	// handle the last data point
	if len(run.Positions) > 0 {
		lastPosition := run.Positions[len(run.Positions)-1]
		if lastPosition.State == shared.Entered {
			lastPosition.State = shared.Exited
			// Force exit the last position using the last data point
			lastData := (*timeframeData)[len(*timeframeData)-1]
			lastPosition.UpdateOnExitSignal(
				&lastData.Close,
				&lastData.Timestamp,
			)
			lastPosition.UpdateOnExitOrder(
				&lastData.Close,
				&lastData.Timestamp,
			)
		} else if lastPosition.State == shared.NotEntered {
			// Remove the last position if it was never entered
			run.Positions = run.Positions[:len(run.Positions)-1]
		}
	}

	run.State = shared.Ended
	endTime := util.GetTime()
	run.EndedTime = &endTime

	run.RunMetrics = engine.CalculateRunMetrics(run)
}

func (engine *Engine) CheckEntryCondition(
	position *shared.Position,
	strategy *shared.Strategy,
	current *shared.Stock,
	data *shared.StockData,
) bool {
	ctx := parser.NewEvaluationContext(data)
	entryCondition := parser.Execute(strategy.EntryExpression, ctx)

	if entryCondition {
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

		return true
	}

	return false
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

	// Removed trailing stop loss calculation
	var trailingStopLossPrice float64 = 0

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
				// TODO: Fix this
				// annualFactor := 365.0 / totalDuration
				// metrics.AnnualizedReturn = (util.Pow(1+totalPnL/(*initialCapital), annualFactor) - 1) * 100
				metrics.AnnualizedReturn = 0
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

<<<<<<< Updated upstream
func (engine *Engine) ProcessExitWithMinuteData(
	position *shared.Position,
	strategy *shared.Strategy,
	minuteData *shared.StockData,
	startTime time.Time,
	endTime time.Time,
) {
	// Find the starting point in the 1-minute data (where timestamp >= position entry time)
	var startIndex int32 = -1
	for i := int32(0); i < int32(len(*minuteData)); i++ {
		if !(*minuteData)[i].Timestamp.Before(startTime) {
			startIndex = i
			break
		}
	}

	// If we couldn't find a valid starting point, use timeframe data for exit
	if startIndex == -1 {
		// Create a single data point for exit evaluation using the entry data
		fallbackStock := shared.Stock{
			Timestamp: *position.EntryTimeOrder,
			Open:      *position.EntryPriceOrder,
			High:      *position.EntryPriceOrder,
			Low:       *position.EntryPriceOrder,
			Close:     *position.EntryPriceOrder,
			Volume:    0,
			Direction: 0,
		}

		// Fallback to using the entry data for exit decision
		slicedData := shared.StockData{fallbackStock}
		engine.CheckExitCondition(position, strategy, &fallbackStock, &slicedData)
		return
	}

	// Find the end index in the 1-minute data (last point before next candle or endTime)
	var endIndex int32 = int32(len(*minuteData)) - 1
	for i := startIndex; i < int32(len(*minuteData)); i++ {
		if !(*minuteData)[i].Timestamp.Before(endTime) {
			endIndex = i - 1
			break
		}
	}

	// If we couldn't find a valid range, return
	if endIndex < startIndex {
		return
	}

	// Process each 1-minute data point until we exit the position or reach the end time
	for i := startIndex; i <= endIndex; i++ {
		current := (*minuteData)[i]

		// Skip processing if position is already exited
		if position.State == shared.Exited {
			break
		}

		// Slice the data from the beginning to the current point for context
		slicedData := (*minuteData)[:i+1]

		// Try to exit the position using current 1-minute data
		engine.CheckExitCondition(position, strategy, &current, &slicedData)
	}
}

func (engine *Engine) ProcessExitWithSecondData(
	position *shared.Position,
	strategy *shared.Strategy,
	secondData *shared.StockData,
	startTime time.Time,
	endTime time.Time,
) {
	// Find the starting point in the 1-second data (where timestamp >= position entry time)
	var startIndex int32 = -1
	for i := int32(0); i < int32(len(*secondData)); i++ {
		if !(*secondData)[i].Timestamp.Before(startTime) {
			startIndex = i
			break
		}
	}

	// TODO: Check this
	// If we couldn't find a valid starting point, use timeframe data for exit
	if startIndex == -1 {
		// Create a single data point for exit evaluation using the entry data
		fallbackStock := shared.Stock{
			Timestamp: *position.EntryTimeOrder,
			Open:      *position.EntryPriceOrder,
			High:      *position.EntryPriceOrder,
			Low:       *position.EntryPriceOrder,
			Close:     *position.EntryPriceOrder,
			Volume:    0,
			Direction: 0,
		}

		// Fallback to using the entry data for exit decision
		slicedData := shared.StockData{fallbackStock}
		engine.CheckExitCondition(position, strategy, &fallbackStock, &slicedData)
		return
	}

	// Find the end index in the 1-second data (last point before next candle or endTime)
	var endIndex int32 = int32(len(*secondData)) - 1
	for i := startIndex; i < int32(len(*secondData)); i++ {
		if !(*secondData)[i].Timestamp.Before(endTime) {
			endIndex = i - 1
			break
		}
	}

	// If we couldn't find a valid range, return
	if endIndex < startIndex {
		return
	}

	// Process each 1-second data point until we exit the position or reach the end time
	for i := startIndex; i <= endIndex; i++ {
		current := (*secondData)[i]

		// Skip processing if position is already exited
		if position.State == shared.Exited {
			break
		}

		// Slice the data from the beginning to the current point for context
		slicedData := (*secondData)[:i+1]

		// Try to exit the position using current 1-second data
		engine.CheckExitCondition(position, strategy, &current, &slicedData)
	}
}

func (engine *Engine) CheckExitCondition(
=======
func (engine *Engine) RunLiveTrading(
	ctx context.Context,
	run *shared.Run,
	strategy *shared.Strategy,
) error {
	// Connect to broker
	if err := engine.broker.Connect(ctx); err != nil {
		return err
	}

	// Subscribe to real-time updates for the instrument
	if err := engine.broker.Subscribe(ctx, []int32{run.InstrumentID}); err != nil {
		return err
	}
	defer engine.broker.Unsubscribe(ctx, []int32{run.InstrumentID})

	// Initialize the run
	startTime := util.GetTime()
	run.StartedTime = &startTime
	run.State = shared.Started
	run.RunType = shared.LiveTrade

	// Get initial position
	position := engine.GetLatestPosition(run, strategy)

	// Start with some historical data for indicators
	var lookBack int32 = 100 // Reduced lookback for live trading
	historicalData := engine.stockService.Get(
		run.InstrumentID,
		strategy.StartPeriod,
		strategy.EndPeriod,
		strategy.ChartTimeframe,
		lookBack,
	)

	data := *historicalData
	
	for {
		select {
		case <-ctx.Done():
			// Handle graceful shutdown
			endTime := util.GetTime()
			run.EndedTime = &endTime
			run.State = shared.Ended
			return ctx.Err()
			
		default:
			// Get latest quote
			quote, err := engine.broker.GetQuote(ctx, run.InstrumentID)
			if err != nil {
				continue
			}

			// Update our data series with the new quote
			current := shared.Stock{
				Close:     quote.Close,
				Open:      quote.Open,
				High:      quote.High,
				Low:       quote.Low,
				Volume:    quote.Volume,
				Timestamp: quote.Timestamp,
			}
			data = append(data, current)

			// Maintain rolling window of data
			if len(data) > 1000 { // Keep last 1000 data points
				data = data[1:]
			}

			switch position.State {
			case shared.NotEntered:
				if engine.tryEnterLive(ctx, position, strategy, &current, &data) {
					// Place actual order through broker
					order := broker.Order{
						InstrumentID: run.InstrumentID,
						Type:        broker.MARKET,
						Direction:   strategy.PositionType,
						Quantity:    strategy.Qty,
					}
					
					orderID, err := engine.broker.PlaceOrder(ctx, order)
					if err != nil {
						continue
					}

					// Monitor order status
					for {
						status, err := engine.broker.GetOrderStatus(ctx, orderID)
						if err != nil {
							break
						}
						if status == broker.EXECUTED {
							// Order was filled, update position
							position.State = shared.Entered
							position.EntryOrderID = util.ParsePtr(orderID)
							break
						} else if status == broker.REJECTED || status == broker.CANCELLED {
							// Order failed, reset position
							position.State = shared.NotEntered
							break
						}
					}
				}

			case shared.Entered:
				if engine.tryExitLive(ctx, position, strategy, &current, &data) {
					// Place exit order
					order := broker.Order{
						InstrumentID: run.InstrumentID,
						Type:        broker.MARKET,
						Direction:   engine.getOppositeDirection(strategy.PositionType),
						Quantity:    strategy.Qty,
					}
					
					orderID, err := engine.broker.PlaceOrder(ctx, order)
					if err != nil {
						continue
					}

					// Monitor exit order status
					for {
						status, err := engine.broker.GetOrderStatus(ctx, orderID)
						if err != nil {
							break
						}
						if status == broker.EXECUTED {
							// Order was filled, update position
							position.State = shared.Exited
							position.ExitOrderID = util.ParsePtr(orderID)
							// Create new position for next trade
							position = engine.GetLatestPosition(run, strategy)
							break
						} else if status == broker.REJECTED || status == broker.CANCELLED {
							// Order failed, keep position open
							break
						}
					}
				}
			}
		}
	}
}

func (engine *Engine) tryEnterLive(
	ctx context.Context,
>>>>>>> Stashed changes
	position *shared.Position,
	strategy *shared.Strategy,
	current *shared.Stock,
	data *shared.StockData,
<<<<<<< Updated upstream
) {
=======
) bool {
	evalCtx := parser.NewEvaluationContext(data)
	return parser.Execute(strategy.EntryExpression, evalCtx)
}

func (engine *Engine) tryExitLive(
	ctx context.Context,
	position *shared.Position,
	strategy *shared.Strategy,
	current *shared.Stock,
	data *shared.StockData,
) bool {
>>>>>>> Stashed changes
	stopLossPrice := position.CalculateStopLossPrice(
		strategy.StopLoss,
		strategy.StopLossType,
	)

	targetProfitPrice := position.CalculateTargetProfitPrice(
		strategy.TargetProfit,
		strategy.TargetProfitType,
	)

<<<<<<< Updated upstream
	// Removed trailing stop loss calculation
	var trailingStopLossPrice float64 = 0

	ctx := parser.NewEvaluationContext(data)

	exitConditionResult := false

	if strategy.ExitExpression != "" {
		exitConditionResult = parser.Execute(strategy.ExitExpression, ctx)
	}

	isExit := position.CheckForExit(
=======
	var trailingStopLossPrice float64
	if strategy.TrailingStopLoss != nil {
		trailingStopLossPrice = position.CalculateTrailingStopLossPrice(
			current.Close,
			strategy.TrailingStopLoss,
			*strategy.TrailingStopLossType,
		)
	}

	evalCtx := parser.NewEvaluationContext(data)
	exitConditionResult := false
	if strategy.ExitExpression != "" {
		exitConditionResult = parser.Execute(strategy.ExitExpression, evalCtx)
	}

	return position.CheckForExit(
>>>>>>> Stashed changes
		current.Close,
		stopLossPrice,
		targetProfitPrice,
		trailingStopLossPrice,
		exitConditionResult,
	)
<<<<<<< Updated upstream

	if isExit {
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
=======
}

func (engine *Engine) getOppositeDirection(direction shared.Direction) shared.Direction {
	if direction == shared.LONG {
		return shared.SHORT
	}
	return shared.LONG
>>>>>>> Stashed changes
}
