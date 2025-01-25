package user

type Service struct {
	repo *Repo
}

func NewServiceProv(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (this *Service) Get() (Users, error) {
	return this.repo.Get()
}

func (this *Service) Create(user *User) (bool, error) {
	return this.repo.Create(user)
}
