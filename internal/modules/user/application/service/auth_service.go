package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/user/application/service/auth"
	"algo/internal/modules/user/domain/model"

	"github.com/google/uuid"
)

// AuthService handles authentication-related functionality
type AuthService struct {
	userService *UserService
	jwtService  *auth.JWTService
}

// AuthResponse represents the authentication response data
type AuthResponse struct {
	Token     string      `json:"token"`
	User      *model.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// NewAuthService creates a new authentication service
func NewAuthService(userService *UserService, jwtService *auth.JWTService) *AuthService {
	return &AuthService{
		userService: userService,
		jwtService:  jwtService,
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	// Authenticate the user
	// Empty password indicates OAuth authentication
	user, err := s.userService.Authenticate(ctx, email, password)
	if err != nil {
		return nil, err
	}

	// Extract role names for the token
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, roleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(24 * time.Hour) // Assuming 24-hour token validity

	return &AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// Register creates a new user and logs them in
func (s *AuthService) Register(ctx context.Context, name, username, email, phoneNumber, password string) (*AuthResponse, error) {
	// Register the user
	user, err := s.userService.Register(ctx, name, username, email, phoneNumber, password)
	if err != nil {
		return nil, err
	}

	// Extract role names for the token
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, roleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(24 * time.Hour) // Assuming 24-hour token validity

	return &AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateToken validates a JWT token and returns the user information
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*model.User, error) {
	// Validate the token
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Get the user by ID
	user, err := s.userService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// RefreshToken refreshes a user's JWT token
func (s *AuthService) RefreshToken(ctx context.Context, oldToken string) (*AuthResponse, error) {
	// Validate the old token
	claims, err := s.jwtService.ValidateToken(oldToken)
	if err != nil {
		return nil, err
	}

	// Get the user by ID
	user, err := s.userService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Check if the user is active
	if !user.IsActive {
		return nil, errors.New("user account is not active")
	}

	// Extract role names for the token
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	// Generate a new JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, roleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(24 * time.Hour) // Assuming 24-hour token validity

	return &AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// GetUserFromToken extracts user information from a token
func (s *AuthService) GetUserFromToken(ctx context.Context, token string) (*model.User, error) {
	// Get user ID from token
	userID, err := s.jwtService.GetUserIDFromToken(token)
	if err != nil {
		return nil, err
	}

	// Get the user by ID
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password and issues a new token
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) (*AuthResponse, error) {
	// Change the password
	if err := s.userService.ChangePassword(ctx, userID, currentPassword, newPassword); err != nil {
		return nil, err
	}

	// Get the updated user
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Extract role names for the token
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	// Generate a new JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, roleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(24 * time.Hour) // Assuming 24-hour token validity

	return &AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// HasRole checks if a user has a specific role
func (s *AuthService) HasRole(token string, role string) (bool, error) {
	return s.jwtService.HasRole(token, role)
}

// GetUserIDFromToken extracts the user ID from a token
func (s *AuthService) GetUserIDFromToken(token string) (uuid.UUID, error) {
	return s.jwtService.GetUserIDFromToken(token)
}
