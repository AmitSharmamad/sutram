package auth_test

import (
	"testing"
	"time"

	"algo/internal/modules/user/application/service/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTService(t *testing.T) {
	// Arrange
	secretKey := "test-secret"
	tokenDuration := 24 * time.Hour

	// Act
	service := auth.NewJWTService(secretKey, tokenDuration)

	// Assert
	assert.NotNil(t, service)
}

func TestGenerateToken(t *testing.T) {
	// Arrange
	secretKey := "test-secret"
	tokenDuration := 24 * time.Hour
	service := auth.NewJWTService(secretKey, tokenDuration)
	userID := uuid.New()
	email := "test@example.com"
	roles := []string{"user", "admin"}

	// Act
	token, err := service.GenerateToken(userID, email, roles)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken_ValidToken(t *testing.T) {
	// Arrange
	secretKey := "test-secret"
	tokenDuration := 24 * time.Hour
	service := auth.NewJWTService(secretKey, tokenDuration)
	userID := uuid.New()
	email := "test@example.com"
	roles := []string{"user", "admin"}

	// Generate a token
	token, err := service.GenerateToken(userID, email, roles)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Act
	claims, err := service.ValidateToken(token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, roles, claims.Roles)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	secretKey := "test-secret"
	tokenDuration := 24 * time.Hour
	service := auth.NewJWTService(secretKey, tokenDuration)
	invalidToken := "invalid.token.string"

	// Act
	claims, err := service.ValidateToken(invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetUserIDFromToken(t *testing.T) {
	// Arrange
	secretKey := "test-secret"
	tokenDuration := 24 * time.Hour
	service := auth.NewJWTService(secretKey, tokenDuration)
	userID := uuid.New()
	email := "test@example.com"
	roles := []string{"user"}

	// Generate a token
	token, err := service.GenerateToken(userID, email, roles)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Act
	extractedUserID, err := service.GetUserIDFromToken(token)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)
}

func TestGetRolesFromToken(t *testing.T) {
	// Arrange
	secretKey := "test-secret"
	tokenDuration := 24 * time.Hour
	service := auth.NewJWTService(secretKey, tokenDuration)
	userID := uuid.New()
	email := "test@example.com"
	roles := []string{"user", "admin", "moderator"}

	// Generate a token
	token, err := service.GenerateToken(userID, email, roles)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Act
	extractedRoles, err := service.GetRolesFromToken(token)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, roles, extractedRoles)
}

func TestHasRole(t *testing.T) {
	// Arrange
	secretKey := "test-secret"
	tokenDuration := 24 * time.Hour
	service := auth.NewJWTService(secretKey, tokenDuration)
	userID := uuid.New()
	email := "test@example.com"
	roles := []string{"user", "admin", "moderator"}

	// Generate a token
	token, err := service.GenerateToken(userID, email, roles)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test cases
	testCases := []struct {
		name     string
		role     string
		expected bool
	}{
		{
			name:     "User has role",
			role:     "admin",
			expected: true,
		},
		{
			name:     "User doesn't have role",
			role:     "superadmin",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			hasRole, err := service.HasRole(token, tc.role)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, hasRole)
		})
	}
}