package stock

import "algo/src/shared"

type Service struct {
	repo *Repo
}

func NewServiceProv(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (service *Service) Get(name string, lastNRecords int32) *shared.StockData {
	return service.repo.Get(lastNRecords)
}
