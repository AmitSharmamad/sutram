package auth

import (
	"errors"
	"fmt"
	"time"

	env "algo/src/core/config"
	"algo/src/shared"
	"algo/src/util"

	"github.com/golang-jwt/jwt"
)

type Service struct {
	secret string
	repo   *Repo
}

func NewServiceProv(
	envProv *env.EnvProv,
	repo *Repo,
) *Service {
	return &Service{
		secret: envProv.GetProp(env.JWT_SIGNING_SECRET),
		repo:   repo,
	}
}

func (service *Service) GetToken(user *shared.User) string {
	secret := service.secret

	claims := &jwt.MapClaims{
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iss":      "sutram",
		"user_id":  user.ID,
		"username": user.Username,
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwt.SignedString(secret)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JWT Token: %v\n", token)

	return token
}

func (service *Service) ValidateToken(token string) (*shared.User, error) {
	secret := service.secret

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, errors.New("INVALID TOKEN")
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		fmt.Println("Token is valid")
		return &shared.User{
			ID:       util.Parse(claims["user_id"].(string)),
			Username: claims["username"].(string),
		}, nil
	} else {
		return nil, errors.New("INVALID TOKEN")
	}
}
