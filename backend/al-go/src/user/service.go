package user

import (
	"errors"
	"fmt"
	"time"

	env "algo/src/core/config"
	"algo/src/shared"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	env  *env.EnvProv
	repo *Repo
}

func NewServiceProv(
	env *env.EnvProv,
	repo *Repo,
) *Service {
	return &Service{
		env:  env,
		repo: repo,
	}
}

func (this *Service) Get() (shared.Users, error) {
	return this.repo.Get()
}

func (this *Service) GetByUsername(username string) (shared.User, error) {
	res, _ := this.repo.GetUsername(username)

	if len(res) != 1 {
		return shared.User{}, errors.New("User not found with the username")
	}

	return res[0], nil
}

func (this *Service) Create(user *User) (bool, error) {
	return this.repo.Create(user)
}

func (this *Service) GetToken(user *User) string {
	mySigningKey := []byte(this.env.GetProp(env.JWT_SIGNING_SECRET))

	// Create the Claims
	claims := &jwt.MapClaims{
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iss":      "sutram",
		"user_id":  user.ID,
		"username": user.Username,
	}

	// Generate the JWT
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwt.SignedString(mySigningKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JWT Token: %v\n", token)

	return token
}
