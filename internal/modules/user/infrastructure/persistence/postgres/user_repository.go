package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"algo/internal/modules/user/domain/model"
	"algo/internal/modules/user/domain/repository"

	"github.com/google/uuid"
)

// UserRepository is the PostgreSQL implementation of the UserRepository interface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO "user" (id, name, username, email, phone_number, password, verified, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.PhoneNumber,
		user.Password,
		user.Verified,
		user.CreatedAt,
		user.UpdatedAt,
		user.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// FindByID finds a user by their ID
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `SELECT id, name, username, email, phone_number, password, verified, created_at, updated_at, is_active FROM "user" WHERE id = $1`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.Password,
		&user.Verified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	// Get user roles
	roles, err := r.getUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

// FindByEmail finds a user by their email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, username, email, phone_number, password, verified, created_at, updated_at, is_active FROM "user" WHERE email = $1`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.Password,
		&user.Verified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	// Get user roles
	roles, err := r.getUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

// FindByUsername finds a user by their username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, name, username, email, phone_number, password, verified, created_at, updated_at, is_active FROM "user" WHERE username = $1`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.Password,
		&user.Verified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	// Get user roles
	roles, err := r.getUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE "user"
		SET name = $1, username = $2, email = $3, phone_number = $4, password = $5, 
			verified = $6, updated_at = $7, is_active = $8
		WHERE id = $9
	`

	user.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Username,
		user.Email,
		user.PhoneNumber,
		user.Password,
		user.Verified,
		user.UpdatedAt,
		user.IsActive,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user by their ID
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM "user" WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// getUserRoles is a helper method to get all roles for a user
func (r *UserRepository) getUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM role r
		JOIN user_role ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		var roleName string

		err := rows.Scan(
			&role.ID,
			&roleName,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}

		role.Name = roleName
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating role rows: %w", err)
	}

	return roles, nil
}
