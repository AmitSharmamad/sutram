package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"algo/internal/modules/user/domain/model"
	"algo/internal/modules/user/domain/repository"

	"github.com/google/uuid"
)

// RoleRepository is the PostgreSQL implementation of the RoleRepository interface
type RoleRepository struct {
	db *sql.DB
}

// NewRoleRepository creates a new PostgreSQL role repository
func NewRoleRepository(db *sql.DB) repository.RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

// FindByID finds a role by ID
func (r *RoleRepository) FindByID(ctx context.Context, id int) (*model.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM role WHERE id = $1`

	var role model.Role
	var roleName string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&roleName,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find role by ID: %w", err)
	}

	role.Name = roleName

	// Get role permissions
	permissions, err := r.getRolePermissions(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions

	return &role, nil
}

// FindByName finds a role by name
func (r *RoleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM role WHERE name = $1`

	var role model.Role
	var roleName string

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&roleName,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find role by name: %w", err)
	}

	role.Name = roleName

	// Get role permissions
	permissions, err := r.getRolePermissions(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions

	return &role, nil
}

// FindAll retrieves all roles
func (r *RoleRepository) FindAll(ctx context.Context) ([]model.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM role ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all roles: %w", err)
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

		// Get role permissions
		permissions, err := r.getRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		role.Permissions = permissions

		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating role rows: %w", err)
	}

	return roles, nil
}

// GetUserRoles gets all roles for a user
func (r *RoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
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

		// Get role permissions
		permissions, err := r.getRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		role.Permissions = permissions

		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating role rows: %w", err)
	}

	return roles, nil
}

// AssignRoleToUser assigns a role to a user
func (r *RoleRepository) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int) error {
	query := `INSERT INTO user_role (user_id, role_id) VALUES ($1, $2) ON CONFLICT (user_id, role_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (r *RoleRepository) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID int) error {
	query := `DELETE FROM user_role WHERE user_id = $1 AND role_id = $2`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	return nil
}

// getRolePermissions is a helper method to get all permissions for a role
func (r *RoleRepository) getRolePermissions(ctx context.Context, roleID int) ([]model.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.resource_type, p.action, p.created_at
		FROM permission p
		JOIN role_permission rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var permission model.Permission
		var permName string

		err := rows.Scan(
			&permission.ID,
			&permName,
			&permission.Description,
			&permission.ResourceType,
			&permission.Action,
			&permission.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}

		permission.Name = permName
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permission rows: %w", err)
	}

	return permissions, nil
}
