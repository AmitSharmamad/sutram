package bt

import (
	"context"

	"algo/src/shared"

	"github.com/google/uuid"
)

type Engine struct {
	service         *Service
	strategyService shared.StrategyService
	stockService    shared.StockService
}

func NewEngineProv(
	service *Service,
	strategyService shared.StrategyService,
	stockService shared.StockService,
) *Engine {
	return &Engine{
		service:         service,
		strategyService: strategyService,
		stockService:    stockService,
	}
}

func (engine *Engine) Run(ctx context.Context, run *shared.Run) {
	strategy := engine.strategyService.Get(ctx, run.StrategyID)

	// dataRange := strategy.GetEntryDataRange()
	dataRange := 1000

	data := engine.stockService.Get("TESLA", dataRange)

	id := run.ID

	run.State = "IN_PROGRESS"

	engine.service.UpdateState(run)

	start := 100
	end := 200

	counter := 0

	println(start)
	println(end)

	for i := start; i < end; i++ {
		println(start, end)
		engine.HandlePosition(
			id,
			strategy,
			(*data)[counter:start],
		)
		counter++
	}

	run.State = "ENDED"

	engine.service.UpdateState(run)
}

func (engine *Engine) HandlePosition(
	runID uuid.UUID,
	strategy *shared.Strategy,
	data shared.StockData,
) {
	position := engine.GetLatestPosition(runID)
	currPos := position.State
	position.Eval(strategy, &data)
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
