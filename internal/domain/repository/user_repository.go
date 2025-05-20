package repository

import (
	"context"

	"algo/internal/domain/entity"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user in the repository
	Create(ctx context.Context, user *entity.User) error

	// FindByID finds a user by their ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// FindByEmail finds a user by their email address
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindByUsername finds a user by their username
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// FindByPhoneNumber finds a user by their phone number
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error)

	// Update updates an existing user in the repository
	Update(ctx context.Context, user *entity.User) error

	// Delete removes a user from the repository
	Delete(ctx context.Context, id uuid.UUID) error

	// GetUserRoles gets all roles assigned to a user
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]entity.Role, error)

	// AssignRoleToUser assigns a role to a user
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int) error

	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID int) error
}
