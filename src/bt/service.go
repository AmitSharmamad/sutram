package bt

import (
	"algo/src/shared"

	"github.com/google/uuid"
)

type Service struct {
	runGroupRepo    *RunGroupRepo
	runRepo         *RunRepo
	positionRepo    *PositionRepo
	positionLogRepo *PositionLogRepo
	runMetricsRepo  *RunMetricsRepo
	strategyService shared.StrategyService
}

func NewServiceProv(
	runGroupRepo *RunGroupRepo,
	runRepo *RunRepo,
	positionRepo *PositionRepo,
	positionLogRepo *PositionLogRepo,
	runMetricsRepo *RunMetricsRepo,
	strategyService shared.StrategyService,
) *Service {
	return &Service{
		runGroupRepo:    runGroupRepo,
		runRepo:         runRepo,
		positionRepo:    positionRepo,
		positionLogRepo: positionLogRepo,
		runMetricsRepo:  runMetricsRepo,
		strategyService: strategyService,
	}
}

func (service *Service) UpdateState(run *shared.Run) {
	service.runRepo.UpdateState(run)
}

func (service *Service) GetLatestPosition(runID uuid.UUID) *shared.Position {
	position, err := service.positionRepo.GetLatest(runID)
	if err != nil {
		panic(err)
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
	run, _ := service.runRepo.Get(id)
	// if err != nil {
	// 	panic(err)
	// }
	// positions, err := service.positionRepo.GetByRunID(id)
	// if err != nil {
	// 	panic("Error fetching the data")
	// }
	return run, nil
}

func (service *Service) UpdateRunState(run *shared.Run) {
	err := service.runRepo.UpdateState(run)
	if err != nil {
		panic(err)
	}
}

func (service *Service) UpdatePosition(position *shared.Position) {
	// update the position value
	err := service.positionRepo.Update(position)
	if err != nil {
		panic(err)
	}

}

func (service *Service) GetRunLogs(id string) *shared.PositionLogs {
	logs, err := service.positionLogRepo.GetByPositionID(id)
	if err != nil {
		panic("Error fetching the data")
	}
	return logs
}

func (service *Service) CreateLog(log *shared.PositionLog) {
	service.positionLogRepo.Create(log)
}

func (service *Service) GetPositions(runID uuid.UUID) shared.Positions {
	positions, err := service.positionRepo.GetByRunID(runID.String())
	if err != nil {
		panic("Error fetching positions")
	}
	return *positions
}
