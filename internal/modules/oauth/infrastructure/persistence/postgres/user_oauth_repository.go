package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"algo/internal/modules/oauth/domain/model"
	"algo/internal/modules/oauth/domain/repository"

	"github.com/google/uuid"
)

// UserOAuthRepository implements the user OAuth repository interface using Postgres
type UserOAuthRepository struct {
	db *sql.DB
}

// NewUserOAuthRepository creates a new user OAuth repository
func NewUserOAuthRepository(db *sql.DB) repository.UserOAuthRepository {
	return &UserOAuthRepository{db: db}
}

// Create creates a new user OAuth connection
func (r *UserOAuthRepository) Create(ctx context.Context, userOAuth *model.UserOAuth) error {
	query := `
		INSERT INTO user_oauth (
			id, user_id, provider_id, provider_user_id, 
			access_token, refresh_token, token_expires_at, user_data,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	userOAuth.ID = uuid.New()
	userOAuth.CreatedAt = time.Now()
	userOAuth.UpdatedAt = time.Now()

	userData, err := json.Marshal(userOAuth.UserData)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(
		ctx,
		query,
		userOAuth.ID,
		userOAuth.UserID,
		userOAuth.ProviderID,
		userOAuth.ProviderUserID,
		userOAuth.AccessToken,
		userOAuth.RefreshToken,
		userOAuth.TokenExpiresAt,
		userData,
		userOAuth.CreatedAt,
		userOAuth.UpdatedAt,
	)

	return err
}

// GetByUserIDAndProvider retrieves a user OAuth connection by user ID and provider ID
func (r *UserOAuthRepository) GetByUserIDAndProvider(ctx context.Context, userID uuid.UUID, providerID int) (*model.UserOAuth, error) {
	query := `
		SELECT 
			uo.id, uo.user_id, uo.provider_id, uo.provider_user_id, 
			uo.access_token, uo.refresh_token, uo.token_expires_at, uo.user_data,
			uo.created_at, uo.updated_at,
			op.id, op.name, op.display_name, op.is_active, op.created_at, op.updated_at
		FROM user_oauth uo
		JOIN oauth_provider op ON uo.provider_id = op.id
		WHERE uo.user_id = $1 AND uo.provider_id = $2
	`

	row := r.db.QueryRowContext(ctx, query, userID, providerID)

	var userOAuth model.UserOAuth
	var provider model.OAuthProvider
	var userData []byte

	err := row.Scan(
		&userOAuth.ID,
		&userOAuth.UserID,
		&userOAuth.ProviderID,
		&userOAuth.ProviderUserID,
		&userOAuth.AccessToken,
		&userOAuth.RefreshToken,
		&userOAuth.TokenExpiresAt,
		&userData,
		&userOAuth.CreatedAt,
		&userOAuth.UpdatedAt,
		&provider.ID,
		&provider.Name,
		&provider.DisplayName,
		&provider.IsActive,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	userOAuth.UserData = userData
	userOAuth.Provider = &provider

	return &userOAuth, nil
}

// GetByProviderUserID retrieves a user OAuth connection by provider ID and provider user ID
func (r *UserOAuthRepository) GetByProviderUserID(ctx context.Context, providerID int, providerUserID string) (*model.UserOAuth, error) {
	query := `
		SELECT 
			uo.id, uo.user_id, uo.provider_id, uo.provider_user_id, 
			uo.access_token, uo.refresh_token, uo.token_expires_at, uo.user_data,
			uo.created_at, uo.updated_at,
			op.id, op.name, op.display_name, op.is_active, op.created_at, op.updated_at
		FROM user_oauth uo
		JOIN oauth_provider op ON uo.provider_id = op.id
		WHERE uo.provider_id = $1 AND uo.provider_user_id = $2
	`

	row := r.db.QueryRowContext(ctx, query, providerID, providerUserID)

	var userOAuth model.UserOAuth
	var provider model.OAuthProvider
	var userData []byte

	err := row.Scan(
		&userOAuth.ID,
		&userOAuth.UserID,
		&userOAuth.ProviderID,
		&userOAuth.ProviderUserID,
		&userOAuth.AccessToken,
		&userOAuth.RefreshToken,
		&userOAuth.TokenExpiresAt,
		&userData,
		&userOAuth.CreatedAt,
		&userOAuth.UpdatedAt,
		&provider.ID,
		&provider.Name,
		&provider.DisplayName,
		&provider.IsActive,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	userOAuth.UserData = userData
	userOAuth.Provider = &provider

	return &userOAuth, nil
}

// UpdateTokens updates the access and refresh tokens for a user OAuth connection
func (r *UserOAuthRepository) UpdateTokens(ctx context.Context, id uuid.UUID, accessToken, refreshToken string, expiresAt int64) error {
	query := `
		UPDATE user_oauth
		SET access_token = $1, refresh_token = $2, token_expires_at = $3, updated_at = $4
		WHERE id = $5
	`

	now := time.Now()
	expiresTime := time.Unix(expiresAt, 0)

	_, err := r.db.ExecContext(ctx, query, accessToken, refreshToken, expiresTime, now, id)
	return err
}

// DeleteByID deletes a user OAuth connection by ID
func (r *UserOAuthRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_oauth WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
