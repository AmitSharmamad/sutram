package auth

import (
	"fmt"
	"time"

	env "algo/src/core/config"
	"algo/src/shared"

	"github.com/golang-jwt/jwt"
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

func (this *Service) GetToken(user *shared.User) string {
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
