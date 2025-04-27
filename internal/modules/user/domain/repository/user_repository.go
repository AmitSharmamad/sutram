package repository

import (
	"context"

	"algo/internal/modules/user/domain/model"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user operations within the user module
type UserRepository interface {
	// Create a new user
	Create(ctx context.Context, user *model.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// FindByUsername finds a user by username
	FindByUsername(ctx context.Context, username string) (*model.User, error)

	// Update updates a user
	Update(ctx context.Context, user *model.User) error

	// Delete deletes a user
	Delete(ctx context.Context, id uuid.UUID) error
}

// RoleRepository defines the interface for role operations
type RoleRepository interface {
	// FindByID finds a role by ID
	FindByID(ctx context.Context, id int) (*model.Role, error)

	// FindByName finds a role by name
	FindByName(ctx context.Context, name string) (*model.Role, error)

	// FindAll retrieves all roles
	FindAll(ctx context.Context) ([]model.Role, error)

	// GetUserRoles gets all roles for a user
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error)

	// AssignRoleToUser assigns a role to a user
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int) error

	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID int) error
}

// PermissionRepository defines the interface for permission operations
type PermissionRepository interface {
	// GetRolePermissions gets all permissions for a role
	GetRolePermissions(ctx context.Context, roleID int) ([]model.Permission, error)

	// GetUserPermissions gets all permissions for a user
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error)
}
