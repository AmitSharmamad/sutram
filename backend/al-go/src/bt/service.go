package bt

import "algo/src/shared"

type Service struct {
	runRepo         *RunRepo
	positionRepo    *PositionRepo
	positionLogRepo *PositionLogRepo
}

func NewServiceProv(
	repo *RunRepo,
	positionRepo *PositionRepo,
	positionLogRepo *PositionLogRepo,
) *Service {
	return &Service{
		runRepo:         repo,
		positionRepo:    positionRepo,
		positionLogRepo: positionLogRepo,
	}
}

func (service *Service) UpdateRunStatus(id string, status RunState) {
	service.runRepo.UpdateStatus(id, status)
}

func (service *Service) CreateRun(strategyID string) *shared.Run {
	// create a run
	err := service.runRepo.Create(strategyID)
	if err != nil {
		panic("")
	}
	run, err := service.runRepo.Get("")
	if err != nil {
		panic("")
	}
	// TODO: emit an event to asynchronously run the strategy back testing
	return run
}

func (service *Service) GetRunPosition(id string) *shared.Position {
	position, err := service.positionRepo.GetLatest(id)
	if err != nil {
		panic("")
	}
	return position
}

func (service *Service) GetRuns(strategyID string) shared.Runs {
	runs, err := service.runRepo.GetByStrategyID(strategyID)
	if err != nil {
		panic("Error fetching the data")
	}
	return runs
}

func (service *Service) GetRun(id string) (*shared.Run, *shared.Positions) {
	run, err := service.runRepo.Get(id)
	if err != nil {
		panic("Error fetching the data")
	}
	positions, err := service.positionRepo.GetByRunID(id)
	if err != nil {
		panic("Error fetching the data")
	}
	return run, positions
}

func (service *Service) GetRunLogs(id string) *shared.PositionLogs {
	logs, err := service.positionLogRepo.GetByPositionID(id)
	if err != nil {
		panic("Error fetching the data")
	}
	return logs
}
