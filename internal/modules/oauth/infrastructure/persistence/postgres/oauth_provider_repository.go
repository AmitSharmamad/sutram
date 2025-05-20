package postgres

import (
	"context"
	"database/sql"
	"errors"

	"algo/internal/modules/oauth/domain/model"
	"algo/internal/modules/oauth/domain/repository"
)

// OAuthProviderRepository implements the OAuth provider repository interface using Postgres
type OAuthProviderRepository struct {
	db *sql.DB
}

// NewOAuthProviderRepository creates a new OAuth provider repository
func NewOAuthProviderRepository(db *sql.DB) repository.OAuthProviderRepository {
	return &OAuthProviderRepository{db: db}
}

// GetByID retrieves an OAuth provider by ID
func (r *OAuthProviderRepository) GetByID(ctx context.Context, id int) (*model.OAuthProvider, error) {
	query := `SELECT id, name, display_name, is_active, created_at, updated_at 
			  FROM oauth_provider 
			  WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var provider model.OAuthProvider
	err := row.Scan(
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

	return &provider, nil
}

// GetByName retrieves an OAuth provider by name
func (r *OAuthProviderRepository) GetByName(ctx context.Context, name string) (*model.OAuthProvider, error) {
	query := `SELECT id, name, display_name, is_active, created_at, updated_at 
			  FROM oauth_provider 
			  WHERE name = $1`

	row := r.db.QueryRowContext(ctx, query, name)

	var provider model.OAuthProvider
	err := row.Scan(
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

	return &provider, nil
}

// GetAll retrieves all OAuth providers
func (r *OAuthProviderRepository) GetAll(ctx context.Context) ([]*model.OAuthProvider, error) {
	query := `SELECT id, name, display_name, is_active, created_at, updated_at 
			  FROM oauth_provider 
			  WHERE is_active = true
			  ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*model.OAuthProvider
	for rows.Next() {
		var provider model.OAuthProvider
		err := rows.Scan(
			&provider.ID,
			&provider.Name,
			&provider.DisplayName,
			&provider.IsActive,
			&provider.CreatedAt,
			&provider.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, &provider)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return providers, nil
}
