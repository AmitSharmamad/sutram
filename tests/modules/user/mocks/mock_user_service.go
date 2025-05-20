package mocks

import (
	"context"

	"algo/internal/modules/user/domain/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the UserService
type MockUserService struct {
	mock.Mock
}

// Register mocks the Register method
func (m *MockUserService) Register(ctx context.Context, name, username, email, phoneNumber, password string) (*model.User, error) {
	args := m.Called(ctx, name, username, email, phoneNumber, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// Authenticate mocks the Authenticate method
func (m *MockUserService) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// FindUserByEmail mocks the FindUserByEmail method
func (m *MockUserService) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// GetUserByID mocks the GetUserByID method
func (m *MockUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// UpdateUser mocks the UpdateUser method
func (m *MockUserService) UpdateUser(ctx context.Context, id uuid.UUID, name, username, email, phoneNumber string) (*model.User, error) {
	args := m.Called(ctx, id, name, username, email, phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// ChangePassword mocks the ChangePassword method
func (m *MockUserService) ChangePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error {
	args := m.Called(ctx, id, currentPassword, newPassword)
	return args.Error(0)
}

// VerifyUser mocks the VerifyUser method
func (m *MockUserService) VerifyUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// DeactivateUser mocks the DeactivateUser method
func (m *MockUserService) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ActivateUser mocks the ActivateUser method
func (m *MockUserService) ActivateUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// AssignRoleToUser mocks the AssignRoleToUser method
func (m *MockUserService) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	args := m.Called(ctx, userID, roleName)
	return args.Error(0)
}

// RemoveRoleFromUser mocks the RemoveRoleFromUser method
func (m *MockUserService) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	args := m.Called(ctx, userID, roleName)
	return args.Error(0)
}

// GetUserRoles mocks the GetUserRoles method
func (m *MockUserService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Role), args.Error(1)
}

// GetAllRoles mocks the GetAllRoles method
func (m *MockUserService) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Role), args.Error(1)
}