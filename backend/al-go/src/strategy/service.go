package strategy

import (
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

func (this *Service) Get(id string) *Strategy {
	strategy, err := this.repo.Get(id)
	fmt.Printf("%+v\n", err)
	return strategy
}

func (strategy *Service) GetData() stock.StockData {
	return nil
}

func (this *Service) GetAll() (Strategies, error) {
	return this.repo.GetAll()
}

func (this *Service) Validate(strategy *Strategy) bool {
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

func (this *Service) Create(strategy *Strategy) (bool, error) {
	return this.repo.Create(strategy)
}
