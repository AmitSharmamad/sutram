package stock

type Service struct {
	repo *Repo
}

func NewServiceProv(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}
