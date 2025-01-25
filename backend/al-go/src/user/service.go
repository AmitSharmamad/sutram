package user

import "errors"

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

func (this *Service) GetByUsername(username string) (User, error) {
	res, _ := this.repo.GetUsername(username)

	if len(res) != 1 {
		return User{}, errors.New("User not found with the username")
	}

	return res[0], nil
}

func (this *Service) Create(user *User) (bool, error) {
	return this.repo.Create(user)
}
