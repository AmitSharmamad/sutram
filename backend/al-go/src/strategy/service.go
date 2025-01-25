package strategy

import "algo/src/stock"

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
	if err != nil {
		panic("")
	}
	return strategy
}

func (strategy *Service) GetData() stock.StockData {
	return nil
}
