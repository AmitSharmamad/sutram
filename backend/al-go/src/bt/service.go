package bt

type Service struct {
	repo *Repo
}

func NewService(
	repo *Repo,
) *Service {
	return &Service{
		repo: repo,
	}
}

func (service *Service) Run(strategyId string) int {
	run := service.repo.Create()
	run.Run()
	return run.ID
}
