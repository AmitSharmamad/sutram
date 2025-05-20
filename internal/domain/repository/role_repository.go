package repository

import (
	"context"

	"algo/internal/domain/entity"
)

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	// FindByID finds a role by its ID
	FindByID(ctx context.Context, id int) (*entity.Role, error)

	// FindByName finds a role by its name
	FindByName(ctx context.Context, name entity.RoleType) (*entity.Role, error)

	// FindAll returns all roles
	FindAll(ctx context.Context) ([]entity.Role, error)

	// GetRolePermissions returns all permissions assigned to a role
	GetRolePermissions(ctx context.Context, roleID int) ([]entity.Permission, error)
}
