package bt

import (
	"context"

	"algo/src/shared"
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
		strategyService: strategyService,
		stockService:    stockService,
	}
}

func (engine *Engine) Run(run *shared.Run) {
	ctx := context.Background()

	/// get the data range from the strategy
	s := engine.strategyService.Get(ctx, run.StrategyID)

	dataRange := s.GetEntryDataRange()

	/// fetch the data
	data := engine.stockService.Get("", dataRange)

	id := run.ID

	engine.service.UpdateRunStatus(id, InProgress)

	/// start taking positions
	/// begin from period start date and iterate till the end
	for i := 0; i < len(*data); i++ {
		position := engine.service.GetRunPosition(id)
		position.Eval(data)
		position.IsDone()
	}

	engine.service.UpdateRunStatus(id, Ended)
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
