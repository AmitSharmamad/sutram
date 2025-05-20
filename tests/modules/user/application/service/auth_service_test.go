package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"algo/internal/modules/user/application/service"
	"algo/internal/modules/user/application/service/auth"
	"algo/internal/modules/user/domain/model"
	"algo/tests/modules/user/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewAuthService(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)

	// Act
	authService := service.NewAuthService(userService, jwtService)

	// Assert
	assert.NotNil(t, authService)
}

func TestLogin(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	userID := uuid.New()
	
	// User with roles
	user := &model.User{
		ID:    userID,
		Email: email,
		Roles: []model.Role{
			{
				ID:   1,
				Name: "user",
			},
			{
				ID:   2,
				Name: "admin",
			},
		},
	}
	
	// Mock user authentication
	userService.On("Authenticate", ctx, email, password).Return(user, nil)
	
	// Mock token generation
	expectedToken := "jwt-token-string"
	jwtService.On("GenerateToken", userID, email, []string{"user", "admin"}).Return(expectedToken, nil)

	// Act
	authResponse, err := authService.Login(ctx, email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, authResponse)
	assert.Equal(t, expectedToken, authResponse.Token)
	assert.Equal(t, user, authResponse.User)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), authResponse.ExpiresAt, 2*time.Second)

	// Verify mocks were called as expected
	userService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestLogin_AuthenticationFailed(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	ctx := context.Background()
	email := "test@example.com"
	password := "wrong-password"
	
	// Mock authentication failure
	authError := errors.New("invalid credentials")
	userService.On("Authenticate", ctx, email, password).Return(nil, authError)

	// Act
	authResponse, err := authService.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, authError, err)
	assert.Nil(t, authResponse)

	// Verify mocks were called as expected
	userService.AssertExpectations(t)
	jwtService.AssertNotCalled(t, "GenerateToken")
}

func TestRegister(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	ctx := context.Background()
	name := "Test User"
	username := "testuser"
	email := "test@example.com"
	phoneNumber := "+123456789"
	password := "password123"
	userID := uuid.New()
	
	// User with roles
	user := &model.User{
		ID:          userID,
		Name:        name,
		Username:    username,
		Email:       email,
		PhoneNumber: phoneNumber,
		Roles: []model.Role{
			{
				ID:   1,
				Name: "free_user",
			},
		},
	}
	
	// Mock user registration
	userService.On("Register", ctx, name, username, email, phoneNumber, password).Return(user, nil)
	
	// Mock token generation
	expectedToken := "jwt-token-string"
	jwtService.On("GenerateToken", userID, email, []string{"free_user"}).Return(expectedToken, nil)

	// Act
	authResponse, err := authService.Register(ctx, name, username, email, phoneNumber, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, authResponse)
	assert.Equal(t, expectedToken, authResponse.Token)
	assert.Equal(t, user, authResponse.User)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), authResponse.ExpiresAt, 2*time.Second)

	// Verify mocks were called as expected
	userService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestRegister_RegistrationFailed(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	ctx := context.Background()
	name := "Test User"
	username := "testuser"
	email := "test@example.com"
	phoneNumber := "+123456789"
	password := "password123"
	
	// Mock registration failure
	registerError := errors.New("registration failed")
	userService.On("Register", ctx, name, username, email, phoneNumber, password).Return(nil, registerError)

	// Act
	authResponse, err := authService.Register(ctx, name, username, email, phoneNumber, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, registerError, err)
	assert.Nil(t, authResponse)

	// Verify mocks were called as expected
	userService.AssertExpectations(t)
	jwtService.AssertNotCalled(t, "GenerateToken")
}

func TestValidateToken(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	ctx := context.Background()
	token := "valid-jwt-token"
	userID := uuid.New()
	
	// Mock token validation
	claims := &auth.TokenClaims{
		UserID: userID,
		Email:  "test@example.com",
		Roles:  []string{"user"},
	}
	jwtService.On("ValidateToken", token).Return(claims, nil)
	
	// Mock user retrieval
	user := &model.User{
		ID:    userID,
		Email: "test@example.com",
	}
	userService.On("GetUserByID", ctx, userID).Return(user, nil)

	// Act
	validatedUser, err := authService.ValidateToken(ctx, token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, validatedUser)
	assert.Equal(t, user, validatedUser)

	// Verify mocks were called as expected
	jwtService.AssertExpectations(t)
	userService.AssertExpectations(t)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	ctx := context.Background()
	token := "invalid-jwt-token"
	
	// Mock token validation failure
	tokenError := errors.New("invalid token")
	jwtService.On("ValidateToken", token).Return(nil, tokenError)

	// Act
	validatedUser, err := authService.ValidateToken(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, tokenError, err)
	assert.Nil(t, validatedUser)

	// Verify mocks were called as expected
	jwtService.AssertExpectations(t)
	userService.AssertNotCalled(t, "GetUserByID")
}

func TestRefreshToken(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	ctx := context.Background()
	oldToken := "old-jwt-token"
	userID := uuid.New()
	email := "test@example.com"
	
	// Mock token validation
	claims := &auth.TokenClaims{
		UserID: userID,
		Email:  email,
		Roles:  []string{"user", "admin"},
	}
	jwtService.On("ValidateToken", oldToken).Return(claims, nil)
	
	// Mock user retrieval
	user := &model.User{
		ID:    userID,
		Email: email,
		IsActive: true,
		Roles: []model.Role{
			{
				ID:   1,
				Name: "user",
			},
			{
				ID:   2,
				Name: "admin",
			},
		},
	}
	userService.On("GetUserByID", ctx, userID).Return(user, nil)
	
	// Mock token generation
	newToken := "new-jwt-token"
	jwtService.On("GenerateToken", userID, email, []string{"user", "admin"}).Return(newToken, nil)

	// Act
	authResponse, err := authService.RefreshToken(ctx, oldToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, authResponse)
	assert.Equal(t, newToken, authResponse.Token)
	assert.Equal(t, user, authResponse.User)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), authResponse.ExpiresAt, 2*time.Second)

	// Verify mocks were called as expected
	jwtService.AssertExpectations(t)
	userService.AssertExpectations(t)
}

func TestRefreshToken_InactiveUser(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	ctx := context.Background()
	oldToken := "old-jwt-token"
	userID := uuid.New()
	
	// Mock token validation
	claims := &auth.TokenClaims{
		UserID: userID,
		Email:  "test@example.com",
	}
	jwtService.On("ValidateToken", oldToken).Return(claims, nil)
	
	// Mock user retrieval (inactive user)
	user := &model.User{
		ID:       userID,
		Email:    "test@example.com",
		IsActive: false,
	}
	userService.On("GetUserByID", ctx, userID).Return(user, nil)

	// Act
	authResponse, err := authService.RefreshToken(ctx, oldToken)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not active")
	assert.Nil(t, authResponse)

	// Verify mocks were called as expected
	jwtService.AssertExpectations(t)
	userService.AssertExpectations(t)
}

func TestGetUserFromToken(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	ctx := context.Background()
	token := "jwt-token"
	userID := uuid.New()
	
	// Mock token extraction
	jwtService.On("GetUserIDFromToken", token).Return(userID, nil)
	
	// Mock user retrieval
	user := &model.User{
		ID:    userID,
		Email: "test@example.com",
	}
	userService.On("GetUserByID", ctx, userID).Return(user, nil)

	// Act
	extractedUser, err := authService.GetUserFromToken(ctx, token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, extractedUser)
	assert.Equal(t, user, extractedUser)

	// Verify mocks were called as expected
	jwtService.AssertExpectations(t)
	userService.AssertExpectations(t)
}

func TestChangePassword(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	ctx := context.Background()
	userID := uuid.New()
	currentPassword := "old-password"
	newPassword := "new-password"
	
	// Mock password change
	userService.On("ChangePassword", ctx, userID, currentPassword, newPassword).Return(nil)
	
	// Mock user retrieval after password change
	user := &model.User{
		ID:    userID,
		Email: "test@example.com",
		Roles: []model.Role{
			{
				ID:   1,
				Name: "user",
			},
		},
	}
	userService.On("GetUserByID", ctx, userID).Return(user, nil)
	
	// Mock token generation
	newToken := "new-jwt-token"
	jwtService.On("GenerateToken", userID, user.Email, []string{"user"}).Return(newToken, nil)

	// Act
	authResponse, err := authService.ChangePassword(ctx, userID, currentPassword, newPassword)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, authResponse)
	assert.Equal(t, newToken, authResponse.Token)
	assert.Equal(t, user, authResponse.User)

	// Verify mocks were called as expected
	userService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestHasRole(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	token := "jwt-token"
	role := "admin"
	
	// Mock role check
	jwtService.On("HasRole", token, role).Return(true, nil)

	// Act
	hasRole, err := authService.HasRole(token, role)

	// Assert
	assert.NoError(t, err)
	assert.True(t, hasRole)

	// Verify mocks were called as expected
	jwtService.AssertExpectations(t)
}

func TestGetUserIDFromToken(t *testing.T) {
	// Arrange
	userService := new(mocks.MockUserService)
	jwtService := new(mocks.MockJWTService)
	authService := service.NewAuthService(userService, jwtService)

	token := "jwt-token"
	userID := uuid.New()
	
	// Mock user ID extraction
	jwtService.On("GetUserIDFromToken", token).Return(userID, nil)

	// Act
	extractedUserID, err := authService.GetUserIDFromToken(token)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)

	// Verify mocks were called as expected
	jwtService.AssertExpectations(t)
}