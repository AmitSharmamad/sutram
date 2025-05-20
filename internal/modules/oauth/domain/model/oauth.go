package model

import (
	"time"

	"github.com/google/uuid"
)

// OAuthProvider represents an OAuth provider
type OAuthProvider struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserOAuth represents a user's OAuth connection
type UserOAuth struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	UserID         uuid.UUID      `json:"user_id" db:"user_id"`
	ProviderID     int            `json:"provider_id" db:"provider_id"`
	ProviderUserID string         `json:"provider_user_id" db:"provider_user_id"`
	AccessToken    string         `json:"access_token" db:"access_token"`
	RefreshToken   string         `json:"refresh_token" db:"refresh_token"`
	TokenExpiresAt time.Time      `json:"token_expires_at" db:"token_expires_at"`
	UserData       []byte         `json:"user_data" db:"user_data"` // JSONB in database
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
	Provider       *OAuthProvider `json:"provider,omitempty" db:"-"` // Loaded separately
}

// ProviderType is an enum for OAuth provider types
type ProviderType string

// Provider type constants
const (
	ProviderGoogle ProviderType = "google"
	ProviderGitHub ProviderType = "github"
)

// GetProviderTypeByName returns the provider type by name
func GetProviderTypeByName(name string) ProviderType {
	switch name {
	case "google":
		return ProviderGoogle
	case "github":
		return ProviderGitHub
	default:
		return ""
	}
}
