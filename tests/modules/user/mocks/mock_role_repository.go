package mocks

import (
	"context"

	"algo/internal/modules/user/domain/model"
	"algo/internal/modules/user/domain/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockRoleRepository is a mock implementation of the RoleRepository interface
type MockRoleRepository struct {
	mock.Mock
}

// Ensure MockRoleRepository implements the RoleRepository interface
var _ repository.RoleRepository = (*MockRoleRepository)(nil)

// FindByID mocks the FindByID method
func (m *MockRoleRepository) FindByID(ctx context.Context, id int) (*model.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

// FindByName mocks the FindByName method
func (m *MockRoleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

// FindAll mocks the FindAll method
func (m *MockRoleRepository) FindAll(ctx context.Context) ([]model.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Role), args.Error(1)
}

// GetUserRoles mocks the GetUserRoles method
func (m *MockRoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Role), args.Error(1)
}

// AssignRoleToUser mocks the AssignRoleToUser method
func (m *MockRoleRepository) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

// RemoveRoleFromUser mocks the RemoveRoleFromUser method
func (m *MockRoleRepository) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID int) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}