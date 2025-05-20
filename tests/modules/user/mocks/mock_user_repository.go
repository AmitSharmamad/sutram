package mocks

import (
	"context"

	"algo/internal/modules/user/domain/model"
	"algo/internal/modules/user/domain/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

// Ensure MockUserRepository implements the UserRepository interface
var _ repository.UserRepository = (*MockUserRepository)(nil)

// Create mocks the Create method
func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// FindByEmail mocks the FindByEmail method
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// FindByUsername mocks the FindByUsername method
func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// Update mocks the Update method
func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}