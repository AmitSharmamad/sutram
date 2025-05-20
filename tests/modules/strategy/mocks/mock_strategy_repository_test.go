package mocks_test

import (
	"context"
	"errors"
	"testing"

	"algo/internal/modules/strategy/domain/model"
	"algo/tests/modules/strategy/mocks"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestMockStrategyRepository_Create(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockStrategyRepository{}
	ctx := context.Background()
	strategy := &model.Strategy{ID: uuid.New()}
	testErr := errors.New("test error")

	t.Run("With function set", func(t *testing.T) {
		// Arrange
		mockRepo.CreateFunc = func(ctx context.Context, s *model.Strategy) error {
			assert.Equal(t, strategy, s)
			return testErr
		}

		// Act
		err := mockRepo.Create(ctx, strategy)

		// Assert
		assert.Equal(t, testErr, err)
	})

	t.Run("Without function set", func(t *testing.T) {
		// Arrange
		mockRepo.CreateFunc = nil

		// Act
		err := mockRepo.Create(ctx, strategy)

		// Assert
		assert.Nil(t, err)
	})
}

func TestMockStrategyRepository_GetByID(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockStrategyRepository{}
	ctx := context.Background()
	id := uuid.New()
	testStrategy := &model.Strategy{ID: id}
	testErr := errors.New("test error")

	t.Run("With function set", func(t *testing.T) {
		// Arrange
		mockRepo.GetByIDFunc = func(ctx context.Context, i uuid.UUID) (*model.Strategy, error) {
			assert.Equal(t, id, i)
			return testStrategy, testErr
		}

		// Act
		strategy, err := mockRepo.GetByID(ctx, id)

		// Assert
		assert.Equal(t, testStrategy, strategy)
		assert.Equal(t, testErr, err)
	})

	t.Run("Without function set", func(t *testing.T) {
		// Arrange
		mockRepo.GetByIDFunc = nil

		// Act
		strategy, err := mockRepo.GetByID(ctx, id)

		// Assert
		assert.Nil(t, strategy)
		assert.Nil(t, err)
	})
}

func TestMockStrategyRepository_GetByUserID(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockStrategyRepository{}
	ctx := context.Background()
	userID := uuid.New()
	testStrategies := []*model.Strategy{{ID: uuid.New()}}
	testErr := errors.New("test error")

	t.Run("With function set", func(t *testing.T) {
		// Arrange
		mockRepo.GetByUserIDFunc = func(ctx context.Context, id uuid.UUID) ([]*model.Strategy, error) {
			assert.Equal(t, userID, id)
			return testStrategies, testErr
		}

		// Act
		strategies, err := mockRepo.GetByUserID(ctx, userID)

		// Assert
		assert.Equal(t, testStrategies, strategies)
		assert.Equal(t, testErr, err)
	})

	t.Run("Without function set", func(t *testing.T) {
		// Arrange
		mockRepo.GetByUserIDFunc = nil

		// Act
		strategies, err := mockRepo.GetByUserID(ctx, userID)

		// Assert
		assert.Nil(t, strategies)
		assert.Nil(t, err)
	})
}

func TestMockStrategyRepository_Update(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockStrategyRepository{}
	ctx := context.Background()
	strategy := &model.Strategy{ID: uuid.New()}
	testErr := errors.New("test error")

	t.Run("With function set", func(t *testing.T) {
		// Arrange
		mockRepo.UpdateFunc = func(ctx context.Context, s *model.Strategy) error {
			assert.Equal(t, strategy, s)
			return testErr
		}

		// Act
		err := mockRepo.Update(ctx, strategy)

		// Assert
		assert.Equal(t, testErr, err)
	})

	t.Run("Without function set", func(t *testing.T) {
		// Arrange
		mockRepo.UpdateFunc = nil

		// Act
		err := mockRepo.Update(ctx, strategy)

		// Assert
		assert.Nil(t, err)
	})
}

func TestMockStrategyRepository_Delete(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockStrategyRepository{}
	ctx := context.Background()
	id := uuid.New()
	testErr := errors.New("test error")

	t.Run("With function set", func(t *testing.T) {
		// Arrange
		mockRepo.DeleteFunc = func(ctx context.Context, i uuid.UUID) error {
			assert.Equal(t, id, i)
			return testErr
		}

		// Act
		err := mockRepo.Delete(ctx, id)

		// Assert
		assert.Equal(t, testErr, err)
	})

	t.Run("Without function set", func(t *testing.T) {
		// Arrange
		mockRepo.DeleteFunc = nil

		// Act
		err := mockRepo.Delete(ctx, id)

		// Assert
		assert.Nil(t, err)
	})
}

func TestMockStrategyRepository_ActivateStrategy(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockStrategyRepository{}
	ctx := context.Background()
	id := uuid.New()
	testErr := errors.New("test error")

	t.Run("With function set", func(t *testing.T) {
		// Arrange
		mockRepo.ActivateStrategyFunc = func(ctx context.Context, i uuid.UUID) error {
			assert.Equal(t, id, i)
			return testErr
		}

		// Act
		err := mockRepo.ActivateStrategy(ctx, id)

		// Assert
		assert.Equal(t, testErr, err)
	})

	t.Run("Without function set", func(t *testing.T) {
		// Arrange
		mockRepo.ActivateStrategyFunc = nil

		// Act
		err := mockRepo.ActivateStrategy(ctx, id)

		// Assert
		assert.Nil(t, err)
	})
}

func TestMockStrategyRepository_DeactivateStrategy(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockStrategyRepository{}
	ctx := context.Background()
	id := uuid.New()
	testErr := errors.New("test error")

	t.Run("With function set", func(t *testing.T) {
		// Arrange
		mockRepo.DeactivateStrategyFunc = func(ctx context.Context, i uuid.UUID) error {
			assert.Equal(t, id, i)
			return testErr
		}

		// Act
		err := mockRepo.DeactivateStrategy(ctx, id)

		// Assert
		assert.Equal(t, testErr, err)
	})

	t.Run("Without function set", func(t *testing.T) {
		// Arrange
		mockRepo.DeactivateStrategyFunc = nil

		// Act
		err := mockRepo.DeactivateStrategy(ctx, id)

		// Assert
		assert.Nil(t, err)
	})
}

func TestMockStrategyRepository_GetActiveStrategies(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockStrategyRepository{}
	ctx := context.Background()
	testStrategies := []*model.Strategy{{ID: uuid.New()}}
	testErr := errors.New("test error")

	t.Run("With function set", func(t *testing.T) {
		// Arrange
		mockRepo.GetActiveStrategiesFunc = func(ctx context.Context) ([]*model.Strategy, error) {
			return testStrategies, testErr
		}

		// Act
		strategies, err := mockRepo.GetActiveStrategies(ctx)

		// Assert
		assert.Equal(t, testStrategies, strategies)
		assert.Equal(t, testErr, err)
	})

	t.Run("Without function set", func(t *testing.T) {
		// Arrange
		mockRepo.GetActiveStrategiesFunc = nil

		// Act
		strategies, err := mockRepo.GetActiveStrategies(ctx)

		// Assert
		assert.Nil(t, strategies)
		assert.Nil(t, err)
	})
}

func TestMockStrategyRepository_GetActiveStrategiesByUserID(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockStrategyRepository{}
	ctx := context.Background()
	userID := uuid.New()
	testStrategies := []*model.Strategy{{ID: uuid.New()}}
	testErr := errors.New("test error")

	t.Run("With function set", func(t *testing.T) {
		// Arrange
		mockRepo.GetActiveStrategiesByUserIDFunc = func(ctx context.Context, id uuid.UUID) ([]*model.Strategy, error) {
			assert.Equal(t, userID, id)
			return testStrategies, testErr
		}

		// Act
		strategies, err := mockRepo.GetActiveStrategiesByUserID(ctx, userID)

		// Assert
		assert.Equal(t, testStrategies, strategies)
		assert.Equal(t, testErr, err)
	})

	t.Run("Without function set", func(t *testing.T) {
		// Arrange
		mockRepo.GetActiveStrategiesByUserIDFunc = nil

		// Act
		strategies, err := mockRepo.GetActiveStrategiesByUserID(ctx, userID)

		// Assert
		assert.Nil(t, strategies)
		assert.Nil(t, err)
	})
}

func TestMockStrategyRepository_SearchStrategies(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockStrategyRepository{}
	ctx := context.Background()
	filters := map[string]interface{}{"name": "test"}
	testStrategies := []*model.Strategy{{ID: uuid.New()}}
	testErr := errors.New("test error")

	t.Run("With function set", func(t *testing.T) {
		// Arrange
		mockRepo.SearchStrategiesFunc = func(ctx context.Context, f map[string]interface{}) ([]*model.Strategy, error) {
			assert.Equal(t, filters, f)
			return testStrategies, testErr
		}

		// Act
		strategies, err := mockRepo.SearchStrategies(ctx, filters)

		// Assert
		assert.Equal(t, testStrategies, strategies)
		assert.Equal(t, testErr, err)
	})

	t.Run("Without function set", func(t *testing.T) {
		// Arrange
		mockRepo.SearchStrategiesFunc = nil

		// Act
		strategies, err := mockRepo.SearchStrategies(ctx, filters)

		// Assert
		assert.Nil(t, strategies)
		assert.Nil(t, err)
	})
}

// Helper function to create a test strategy
func createTestStrategy() *model.Strategy {
	return model.NewStrategy(
		uuid.New(),
		"Test Strategy",
		"Test Description",
		model.Long,
		pq.Int32Array{1, 2, 3},
		model.Candlestick,
		model.OneHour,
		model.TradeType("MIS"),
		10,
		"close > open",
		"close < open",
		5.0,
		model.Percentage,
		2.0,
		model.Percentage,
		10000.0,
	)
}