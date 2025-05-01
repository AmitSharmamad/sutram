package http

import (
	"context"
	"errors"

	"algo/internal/modules/user/application/service"
)

// AuthServiceAdapter adapts the user module's auth service to the middleware's AuthService interface
type AuthServiceAdapter struct {
	authService *service.AuthService
}

// NewAuthServiceAdapter creates a new auth service adapter
func NewAuthServiceAdapter(authService *service.AuthService) *AuthServiceAdapter {
	return &AuthServiceAdapter{
		authService: authService,
	}
}

// ValidateToken implements middleware.AuthService
func (a *AuthServiceAdapter) ValidateToken(ctx interface{}, tokenString string) (interface{}, error) {
	// Try to convert the interface to context.Context
	contextObj, ok := ctx.(context.Context)
	if !ok {
		return nil, errors.New("invalid context type")
	}

	// Call the auth service's ValidateToken method
	user, err := a.authService.ValidateToken(contextObj, tokenString)
	if err != nil {
		return nil, err
	}

	// Return the user as an interface{} as required by the middleware interface
	return user, nil
}

// GetUserIDFromToken implements middleware.AuthService
func (a *AuthServiceAdapter) GetUserIDFromToken(tokenString string) (interface{}, error) {
	return a.authService.GetUserIDFromToken(tokenString)
}

// HasRole implements middleware.AuthService
func (a *AuthServiceAdapter) HasRole(tokenString string, role string) (bool, error) {
	return a.authService.HasRole(tokenString, role)
}
