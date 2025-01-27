package strategy

import (
	"context"
	"fmt"

	"algo/src/stock"
	"algo/src/util"
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
	strategy *Strategy,
) (bool, error) {
	return service.repo.Create(ctx, strategy)
}

func (service *Service) Get(
	ctx context.Context,
	id string,
) *Strategy {
	strategy, err := service.repo.Get(ctx, id)
	fmt.Printf("%+v\n", err)
	return strategy
}

func (strategy *Service) GetData() stock.StockData {
	return nil
}

func (service *Service) GetAll(
	ctx context.Context,
) (Strategies, error) {
	return service.repo.GetAll(ctx)
}

func (service *Service) Validate(strategy *Strategy) bool {
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
