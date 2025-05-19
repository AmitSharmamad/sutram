package provider

import (
	"algo/internal/modules/oauth/application/dto"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// OAuthProvider defines the interface for OAuth providers
type OAuthProvider interface {
	GetLoginURL(state string) string
	GetUserInfo(ctx context.Context, code string) (*dto.OAuthResponse, error)
	GetProviderName() string
}

// GoogleOAuthProvider implements the OAuthProvider interface for Google
type GoogleOAuthProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
	scopes       []string
}

// GoogleUserInfo represents the user info returned by Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// NewGoogleOAuthProvider creates a new Google OAuth provider
func NewGoogleOAuthProvider() *GoogleOAuthProvider {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_OAUTH_REDIRECT_URL")

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		panic("Google OAuth configuration is missing. Please set GOOGLE_OAUTH_CLIENT_ID, GOOGLE_OAUTH_CLIENT_SECRET, and GOOGLE_OAUTH_REDIRECT_URL environment variables.")
	}

	return &GoogleOAuthProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		scopes:       []string{"profile", "email"},
	}
}

// GetProviderName returns the provider name
func (p *GoogleOAuthProvider) GetProviderName() string {
	return "google"
}

// GetLoginURL returns the URL to redirect the user to for login
func (p *GoogleOAuthProvider) GetLoginURL(state string) string {
	authURL := "https://accounts.google.com/o/oauth2/v2/auth"

	params := url.Values{}
	params.Add("client_id", p.clientID)
	params.Add("redirect_uri", p.redirectURL)
	params.Add("response_type", "code")
	params.Add("scope", strings.Join(p.scopes, " "))
	params.Add("access_type", "offline")
	params.Add("state", state)
	params.Add("prompt", "consent") // Force to get refresh token

	return authURL + "?" + params.Encode()
}

// GetUserInfo exchanges the authorization code for an access token and gets the user info
func (p *GoogleOAuthProvider) GetUserInfo(ctx context.Context, code string) (*dto.OAuthResponse, error) {
	// Exchange code for token
	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", p.clientID)
	data.Set("client_secret", p.clientSecret)
	data.Set("redirect_uri", p.redirectURL)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get token: %s", body)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		IDToken      string `json:"id_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Get user info
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err = http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", body)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Calculate token expiration time
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return &dto.OAuthResponse{
		Provider:     "google",
		ProviderID:   userInfo.ID,
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
		AvatarURL:    userInfo.Picture,
	}, nil
}
