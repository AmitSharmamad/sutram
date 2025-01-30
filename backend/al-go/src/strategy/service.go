package strategy

import (
	"context"

	"algo/src/shared"
	"algo/src/util"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repo
}

func NewServiceProv(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (service *Service) Create(
	ctx context.Context,
	strategy *shared.Strategy,
) (*shared.Strategy, error) {
	strategy.ID = util.NewUUID()
	strategy.EntryExpression = string(strategy.EntryExpression)
	strategy.ExitExpression = string(strategy.ExitExpression)
	err := service.repo.Create(ctx, strategy)
	return strategy, err
}

func (service *Service) Get(
	ctx context.Context,
	id uuid.UUID,
) *shared.Strategy {
	strategy, err := service.repo.Get(ctx, id)
	if err != nil {
		panic("")
	}
	return strategy
}

func (strategy *Service) GetData() shared.StockData {
	return nil
}

func (service *Service) GetAll(
	ctx context.Context,
) (*shared.Strategies, error) {
	return service.repo.GetAll(ctx)
}

func (service *Service) Validate(strategy *shared.Strategy) bool {
	validEntry := util.Validate(strategy.EntryExpression)
	validExit := util.Validate(strategy.ExitExpression)

	if !validEntry {
		return false
	}

	if !validExit {
		return false
	}

	// TODO: do rest of the validation

	return true
}
