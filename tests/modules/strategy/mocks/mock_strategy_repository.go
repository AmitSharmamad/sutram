package mocks

import (
	"context"

	"algo/internal/modules/strategy/domain/model"

	"github.com/google/uuid"
)

// MockStrategyRepository is a mock implementation of the StrategyRepository interface
type MockStrategyRepository struct {
	CreateFunc               func(ctx context.Context, strategy *model.Strategy) error
	GetByIDFunc              func(ctx context.Context, id uuid.UUID) (*model.Strategy, error)
	GetByUserIDFunc          func(ctx context.Context, userID uuid.UUID) ([]*model.Strategy, error)
	UpdateFunc               func(ctx context.Context, strategy *model.Strategy) error
	DeleteFunc               func(ctx context.Context, id uuid.UUID) error
	ActivateStrategyFunc     func(ctx context.Context, id uuid.UUID) error
	DeactivateStrategyFunc   func(ctx context.Context, id uuid.UUID) error
	GetActiveStrategiesFunc  func(ctx context.Context) ([]*model.Strategy, error)
	GetActiveStrategiesByUserIDFunc func(ctx context.Context, userID uuid.UUID) ([]*model.Strategy, error)
	SearchStrategiesFunc     func(ctx context.Context, filters map[string]interface{}) ([]*model.Strategy, error)
}

// Create implements the repository interface
func (m *MockStrategyRepository) Create(ctx context.Context, strategy *model.Strategy) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, strategy)
	}
	return nil
}

// GetByID implements the repository interface
func (m *MockStrategyRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

// GetByUserID implements the repository interface
func (m *MockStrategyRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Strategy, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

// Update implements the repository interface
func (m *MockStrategyRepository) Update(ctx context.Context, strategy *model.Strategy) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, strategy)
	}
	return nil
}

// Delete implements the repository interface
func (m *MockStrategyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// ActivateStrategy implements the repository interface
func (m *MockStrategyRepository) ActivateStrategy(ctx context.Context, id uuid.UUID) error {
	if m.ActivateStrategyFunc != nil {
		return m.ActivateStrategyFunc(ctx, id)
	}
	return nil
}

// DeactivateStrategy implements the repository interface
func (m *MockStrategyRepository) DeactivateStrategy(ctx context.Context, id uuid.UUID) error {
	if m.DeactivateStrategyFunc != nil {
		return m.DeactivateStrategyFunc(ctx, id)
	}
	return nil
}

// GetActiveStrategies implements the repository interface
func (m *MockStrategyRepository) GetActiveStrategies(ctx context.Context) ([]*model.Strategy, error) {
	if m.GetActiveStrategiesFunc != nil {
		return m.GetActiveStrategiesFunc(ctx)
	}
	return nil, nil
}

// GetActiveStrategiesByUserID implements the repository interface
func (m *MockStrategyRepository) GetActiveStrategiesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Strategy, error) {
	if m.GetActiveStrategiesByUserIDFunc != nil {
		return m.GetActiveStrategiesByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

// SearchStrategies implements the repository interface
func (m *MockStrategyRepository) SearchStrategies(ctx context.Context, filters map[string]interface{}) ([]*model.Strategy, error) {
	if m.SearchStrategiesFunc != nil {
		return m.SearchStrategiesFunc(ctx, filters)
	}
	return nil, nil
}