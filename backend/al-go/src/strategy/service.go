package strategy

import "algo/src/data"

type Service struct {
	Repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		Repo: repo,
	}
}

func (strategy *Service) GetData() data.StockData {
	return nil
}
