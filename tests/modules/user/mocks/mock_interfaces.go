package mocks

import (
	"context"
	
	"algo/internal/modules/user/application/service/auth"
	"algo/internal/modules/user/domain/model"
	
	"github.com/google/uuid"
)

// IUserService defines the interface for the user service
type IUserService interface {
	Register(ctx context.Context, name, username, email, phoneNumber, password string) (*model.User, error)
	Authenticate(ctx context.Context, email, password string) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, name, username, email, phoneNumber string) (*model.User, error)
	ChangePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error
	VerifyUser(ctx context.Context, id uuid.UUID) error
	DeactivateUser(ctx context.Context, id uuid.UUID) error
	ActivateUser(ctx context.Context, id uuid.UUID) error
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) error
	RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleName string) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

// IJWTService defines the interface for the JWT service
type IJWTService interface {
	GenerateToken(userID uuid.UUID, email string, roles []string) (string, error)
	ValidateToken(tokenString string) (*auth.TokenClaims, error)
	GetUserIDFromToken(tokenString string) (uuid.UUID, error)
	GetRolesFromToken(tokenString string) ([]string, error)
	HasRole(tokenString string, role string) (bool, error)
}