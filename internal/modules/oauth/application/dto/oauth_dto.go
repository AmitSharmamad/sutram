package dto

import (
	"time"

	"github.com/google/uuid"
)

// OAuthResponse represents a response from the OAuth provider
type OAuthResponse struct {
	Provider     string    `json:"provider"`
	ProviderID   string    `json:"provider_id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
}

// OAuthLoginURL represents the URL to redirect to for OAuth login
type OAuthLoginURL struct {
	URL        string `json:"url"`
	Provider   string `json:"provider"`
	StateToken string `json:"state_token"`
}

// OAuthCallback represents the callback data from an OAuth provider
type OAuthCallback struct {
	Code  string `json:"code"`
	State string `json:"state"`
	Error string `json:"error,omitempty"`
}

// ApiResponse represents a standard API response
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// OAuthProviderInfo represents information about an OAuth provider
type OAuthProviderInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	LoginURL    string `json:"login_url,omitempty"`
}

// UserOAuthInfo represents information about a user's OAuth connection
type UserOAuthInfo struct {
	UserID         uuid.UUID `json:"user_id"`
	Provider       string    `json:"provider"`
	ProviderUserID string    `json:"provider_user_id"`
	Connected      bool      `json:"connected"`
	ConnectedAt    time.Time `json:"connected_at,omitempty"`
}
