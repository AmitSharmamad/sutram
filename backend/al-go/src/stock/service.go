package stock

type Service struct {
	repo *Repo
}

func NewServiceProv(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (this *Service) Get(id string) (any, error) {
	return this.repo.Select()
}
