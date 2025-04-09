package auth

import (
	env "algo/src/core/config"
	"algo/src/shared"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleService struct {
	env         *env.EnvProv
	config      *oauth2.Config
	userService shared.UserService
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	Picture       string `json:"picture"`
}

func NewGoogleServiceProv(envProv *env.EnvProv, userService shared.UserService) *GoogleService {
	config := &oauth2.Config{
		ClientID:     envProv.GetProp(env.GOOGLE_CLIENT_ID),
		ClientSecret: envProv.GetProp(env.GOOGLE_CLIENT_SECRET),
		RedirectURL:  envProv.GetProp(env.GOOGLE_REDIRECT_URL),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleService{
		env:         envProv,
		config:      config,
		userService: userService,
	}
}

func (s *GoogleService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *GoogleService) HandleCallback(code string) (*shared.User, error) {
	token, err := s.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	userInfo, err := s.getUserInfo(token.AccessToken)
	if err != nil {
		return nil, err
	}

	// Try to find existing user
	user, err := s.userService.GetByUsername(userInfo.Email)
	if err != nil {
		// Create new user if not found
		user = shared.User{
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			Username: userInfo.Email,
			Verified: userInfo.VerifiedEmail,
			IsActive: true,
		}
		// Create user in database
		_, err = s.userService.Create(&user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *GoogleService) getUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer resp.Body.Close()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		return nil, fmt.Errorf("failed parsing user info: %s", err.Error())
	}

	return &userInfo, nil
}
