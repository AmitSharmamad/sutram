package repository

import (
	"context"

	"algo/internal/modules/oauth/domain/model"

	"github.com/google/uuid"
)

// OAuthProviderRepository defines the interface for OAuth provider persistence
type OAuthProviderRepository interface {
	GetByID(ctx context.Context, id int) (*model.OAuthProvider, error)
	GetByName(ctx context.Context, name string) (*model.OAuthProvider, error)
	GetAll(ctx context.Context) ([]*model.OAuthProvider, error)
}

// UserOAuthRepository defines the interface for user OAuth persistence
type UserOAuthRepository interface {
	Create(ctx context.Context, userOAuth *model.UserOAuth) error
	GetByUserIDAndProvider(ctx context.Context, userID uuid.UUID, providerID int) (*model.UserOAuth, error)
	GetByProviderUserID(ctx context.Context, providerID int, providerUserID string) (*model.UserOAuth, error)
	UpdateTokens(ctx context.Context, id uuid.UUID, accessToken, refreshToken string, expiresAt int64) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
