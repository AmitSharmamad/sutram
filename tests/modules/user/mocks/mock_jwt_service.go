package mocks

import (
	"algo/internal/modules/user/application/service/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockJWTService is a mock implementation of the JWTService
type MockJWTService struct {
	mock.Mock
}

// GenerateToken mocks the GenerateToken method
func (m *MockJWTService) GenerateToken(userID uuid.UUID, email string, roles []string) (string, error) {
	args := m.Called(userID, email, roles)
	return args.String(0), args.Error(1)
}

// ValidateToken mocks the ValidateToken method
func (m *MockJWTService) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenClaims), args.Error(1)
}

// GetUserIDFromToken mocks the GetUserIDFromToken method
func (m *MockJWTService) GetUserIDFromToken(tokenString string) (uuid.UUID, error) {
	args := m.Called(tokenString)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

// GetRolesFromToken mocks the GetRolesFromToken method
func (m *MockJWTService) GetRolesFromToken(tokenString string) ([]string, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// HasRole mocks the HasRole method
func (m *MockJWTService) HasRole(tokenString string, role string) (bool, error) {
	args := m.Called(tokenString, role)
	return args.Bool(0), args.Error(1)
}