package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"algo/internal/modules/user/application/service"
	"algo/internal/modules/user/domain/model"
	"algo/tests/modules/user/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewUserService(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)

	// Act
	userService := service.NewUserService(userRepo, roleRepo)

	// Assert
	assert.NotNil(t, userService)
}

func TestRegister(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	name := "Test User"
	username := "testuser"
	email := "test@example.com"
	phoneNumber := "+123456789"
	password := "password123"

	// Mock behavior: No existing users with this email, username, or phone number
	userRepo.On("FindByEmail", ctx, email).Return(nil, nil)
	userRepo.On("FindByUsername", ctx, username).Return(nil, nil)
	userRepo.On("FindByUsername", ctx, phoneNumber).Return(nil, nil)

	// Mock user creation
	userRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	// Mock role assignment
	freeRole := &model.Role{
		ID:          1,
		Name:        "free_user",
		Description: "Free tier user",
	}
	roleRepo.On("FindByName", ctx, "free_user").Return(freeRole, nil)
	roleRepo.On("AssignRoleToUser", ctx, mock.AnythingOfType("uuid.UUID"), freeRole.ID).Return(nil)

	// Act
	user, err := userService.Register(ctx, name, username, email, phoneNumber, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, phoneNumber, user.PhoneNumber)
	assert.NotEmpty(t, user.Password) // Password should be hashed
	assert.NotEqual(t, password, user.Password) // Password should not be stored in plain text
	assert.False(t, user.Verified)
	assert.True(t, user.IsActive)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)
	assert.Len(t, user.Roles, 1)
	assert.Equal(t, "free_user", user.Roles[0].Name)

	// Verify all mocks were called as expected
	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	name := "Test User"
	username := "testuser"
	email := "existing@example.com"
	phoneNumber := "+123456789"
	password := "password123"

	// Mock behavior: Email already exists
	existingUser := &model.User{
		ID:    uuid.New(),
		Email: email,
	}
	userRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	// Act
	user, err := userService.Register(ctx, name, username, email, phoneNumber, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrEmailAlreadyExists, err)
	assert.Nil(t, user)

	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
	// The FindByUsername should not be called if FindByEmail fails
	userRepo.AssertNotCalled(t, "FindByUsername", ctx, username)
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	name := "Test User"
	username := "existinguser"
	email := "test@example.com"
	phoneNumber := "+123456789"
	password := "password123"

	// Mock behavior: Email check passes, but username already exists
	userRepo.On("FindByEmail", ctx, email).Return(nil, nil)
	
	existingUser := &model.User{
		ID:       uuid.New(),
		Username: username,
	}
	userRepo.On("FindByUsername", ctx, username).Return(existingUser, nil)

	// Act
	user, err := userService.Register(ctx, name, username, email, phoneNumber, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrUsernameAlreadyExists, err)
	assert.Nil(t, user)

	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
}

func TestAuthenticate(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	
	// Create a user with bcrypt hashed password
	user := &model.User{
		ID:       uuid.New(),
		Email:    email,
		Password: "$2a$10$Nt8DerURaXZx3P7eGr73nO3BtTJnmva.0BDPt3yfwQDzOQLB4hGGy", // bcrypt hash of "password123"
		IsActive: true,
	}
	
	userRepo.On("FindByEmail", ctx, email).Return(user, nil)
	
	// Mock role loading
	roles := []model.Role{
		{
			ID:   1,
			Name: "user",
		},
	}
	roleRepo.On("GetUserRoles", ctx, user.ID).Return(roles, nil)

	// Act
	authenticatedUser, err := userService.Authenticate(ctx, email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, authenticatedUser)
	assert.Equal(t, user.ID, authenticatedUser.ID)
	assert.Equal(t, user.Email, authenticatedUser.Email)
	assert.Equal(t, roles, authenticatedUser.Roles)

	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

func TestAuthenticate_InvalidCredentials(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	email := "test@example.com"
	wrongPassword := "wrongpassword"
	
	// Create a user with bcrypt hashed password
	user := &model.User{
		ID:       uuid.New(),
		Email:    email,
		Password: "$2a$10$Nt8DerURaXZx3P7eGr73nO3BtTJnmva.0BDPt3yfwQDzOQLB4hGGy", // bcrypt hash of "password123"
		IsActive: true,
	}
	
	userRepo.On("FindByEmail", ctx, email).Return(user, nil)

	// Act
	authenticatedUser, err := userService.Authenticate(ctx, email, wrongPassword)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrInvalidCredentials, err)
	assert.Nil(t, authenticatedUser)

	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
}

func TestFindUserByEmail(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	email := "test@example.com"
	
	user := &model.User{
		ID:    uuid.New(),
		Email: email,
	}
	
	userRepo.On("FindByEmail", ctx, email).Return(user, nil)
	
	// Mock role loading
	roles := []model.Role{
		{
			ID:   1,
			Name: "user",
		},
	}
	roleRepo.On("GetUserRoles", ctx, user.ID).Return(roles, nil)

	// Act
	foundUser, err := userService.FindUserByEmail(ctx, email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, roles, foundUser.Roles)

	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

func TestFindUserByEmail_UserNotFound(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	email := "nonexistent@example.com"
	
	userRepo.On("FindByEmail", ctx, email).Return(nil, nil)

	// Act
	foundUser, err := userService.FindUserByEmail(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrUserNotFound, err)
	assert.Nil(t, foundUser)

	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	userID := uuid.New()
	
	user := &model.User{
		ID:    userID,
		Email: "test@example.com",
	}
	
	userRepo.On("FindByID", ctx, userID).Return(user, nil)

	// Act
	foundUser, err := userService.GetUserByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)

	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
}

func TestGetUserByID_UserNotFound(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	userID := uuid.New()
	
	userRepo.On("FindByID", ctx, userID).Return(nil, nil)

	// Act
	foundUser, err := userService.GetUserByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrUserNotFound, err)
	assert.Nil(t, foundUser)

	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	userID := uuid.New()
	
	existingUser := &model.User{
		ID:          userID,
		Name:        "Old Name",
		Username:    "oldusername",
		Email:       "old@example.com",
		PhoneNumber: "+111111111",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
	}
	
	newName := "New Name"
	newUsername := "newusername"
	newEmail := "new@example.com"
	newPhoneNumber := "+222222222"
	
	userRepo.On("FindByID", ctx, userID).Return(existingUser, nil)
	
	// Check for email uniqueness
	userRepo.On("FindByEmail", ctx, newEmail).Return(nil, nil)
	
	// Check for username uniqueness
	userRepo.On("FindByUsername", ctx, newUsername).Return(nil, nil)
	
	// Check for phone number uniqueness
	userRepo.On("FindByUsername", ctx, newPhoneNumber).Return(nil, nil)
	
	// Mock update
	userRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	// Act
	updatedUser, err := userService.UpdateUser(ctx, userID, newName, newUsername, newEmail, newPhoneNumber)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, userID, updatedUser.ID)
	assert.Equal(t, newName, updatedUser.Name)
	assert.Equal(t, newUsername, updatedUser.Username)
	assert.Equal(t, newEmail, updatedUser.Email)
	assert.Equal(t, newPhoneNumber, updatedUser.PhoneNumber)
	assert.Equal(t, existingUser.CreatedAt, updatedUser.CreatedAt) // Should be unchanged
	assert.NotEqual(t, existingUser.UpdatedAt, updatedUser.UpdatedAt) // Should be updated

	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
}

func TestChangePassword(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	userID := uuid.New()
	
	currentPassword := "password123"
	newPassword := "newpassword456"
	
	// Create a user with bcrypt hashed password for current password
	user := &model.User{
		ID:       userID,
		Email:    "test@example.com",
		Password: "$2a$10$Nt8DerURaXZx3P7eGr73nO3BtTJnmva.0BDPt3yfwQDzOQLB4hGGy", // bcrypt hash of "password123"
	}
	
	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	// Act
	err := userService.ChangePassword(ctx, userID, currentPassword, newPassword)

	// Assert
	assert.NoError(t, err)
	
	// Verify the password was updated (indirectly via the Update call)
	userRepo.AssertCalled(t, "Update", ctx, mock.AnythingOfType("*model.User"))
	
	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
}

func TestVerifyUser(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	userID := uuid.New()
	
	user := &model.User{
		ID:       userID,
		Verified: false,
	}
	
	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Run(func(args mock.Arguments) {
		updatedUser := args.Get(1).(*model.User)
		assert.True(t, updatedUser.Verified)
	}).Return(nil)

	// Act
	err := userService.VerifyUser(ctx, userID)

	// Assert
	assert.NoError(t, err)
	
	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
}

func TestDeactivateUser(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	userID := uuid.New()
	
	user := &model.User{
		ID:       userID,
		IsActive: true,
	}
	
	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Run(func(args mock.Arguments) {
		updatedUser := args.Get(1).(*model.User)
		assert.False(t, updatedUser.IsActive)
	}).Return(nil)

	// Act
	err := userService.DeactivateUser(ctx, userID)

	// Assert
	assert.NoError(t, err)
	
	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
}

func TestActivateUser(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	userID := uuid.New()
	
	user := &model.User{
		ID:       userID,
		IsActive: false,
	}
	
	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Run(func(args mock.Arguments) {
		updatedUser := args.Get(1).(*model.User)
		assert.True(t, updatedUser.IsActive)
	}).Return(nil)

	// Act
	err := userService.ActivateUser(ctx, userID)

	// Assert
	assert.NoError(t, err)
	
	// Verify mocks were called as expected
	userRepo.AssertExpectations(t)
}

func TestAssignRoleToUser(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	userID := uuid.New()
	roleName := "admin"
	
	role := &model.Role{
		ID:   2,
		Name: roleName,
	}
	
	roleRepo.On("FindByName", ctx, roleName).Return(role, nil)
	roleRepo.On("AssignRoleToUser", ctx, userID, role.ID).Return(nil)

	// Act
	err := userService.AssignRoleToUser(ctx, userID, roleName)

	// Assert
	assert.NoError(t, err)
	
	// Verify mocks were called as expected
	roleRepo.AssertExpectations(t)
}

func TestRemoveRoleFromUser(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	userID := uuid.New()
	roleName := "admin"
	
	role := &model.Role{
		ID:   2,
		Name: roleName,
	}
	
	roleRepo.On("FindByName", ctx, roleName).Return(role, nil)
	roleRepo.On("RemoveRoleFromUser", ctx, userID, role.ID).Return(nil)

	// Act
	err := userService.RemoveRoleFromUser(ctx, userID, roleName)

	// Assert
	assert.NoError(t, err)
	
	// Verify mocks were called as expected
	roleRepo.AssertExpectations(t)
}

func TestGetUserRoles(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	userID := uuid.New()
	
	roles := []model.Role{
		{
			ID:   1,
			Name: "user",
		},
		{
			ID:   2,
			Name: "admin",
		},
	}
	
	roleRepo.On("GetUserRoles", ctx, userID).Return(roles, nil)

	// Act
	userRoles, err := userService.GetUserRoles(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, roles, userRoles)
	
	// Verify mocks were called as expected
	roleRepo.AssertExpectations(t)
}

func TestGetAllRoles(t *testing.T) {
	// Arrange
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	userService := service.NewUserService(userRepo, roleRepo)

	ctx := context.Background()
	
	roles := []model.Role{
		{
			ID:   1,
			Name: "user",
		},
		{
			ID:   2,
			Name: "admin",
		},
		{
			ID:   3,
			Name: "moderator",
		},
	}
	
	roleRepo.On("FindAll", ctx).Return(roles, nil)

	// Act
	allRoles, err := userService.GetAllRoles(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, roles, allRoles)
	
	// Verify mocks were called as expected
	roleRepo.AssertExpectations(t)
}