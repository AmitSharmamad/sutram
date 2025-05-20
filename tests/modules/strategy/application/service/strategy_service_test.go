package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"algo/internal/modules/strategy/application/dto"
	"algo/internal/modules/strategy/application/service"
	"algo/internal/modules/strategy/domain/model"
	"algo/tests/modules/strategy/mocks"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateStrategy_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := createValidCreateStrategyRequest()
	userID := uuid.New()

	mockRepo := &mocks.MockStrategyRepository{
		CreateFunc: func(ctx context.Context, strategy *model.Strategy) error {
			return nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.CreateStrategy(ctx, req, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.Name, response.Name)
	assert.Equal(t, req.Description, response.Description)
	assert.Equal(t, userID, response.UserID)
}

func TestCreateStrategy_ValidationError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := createValidCreateStrategyRequest()
	req.Name = "" // Invalid name
	userID := uuid.New()

	mockRepo := &mocks.MockStrategyRepository{}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.CreateStrategy(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, service.ErrInvalidRequest, err)
}

func TestCreateStrategy_EmptyExpressionError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := createValidCreateStrategyRequest()
	req.EntryExpression = "" // Empty expression - this will fail validation
	userID := uuid.New()

	mockRepo := &mocks.MockStrategyRepository{}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.CreateStrategy(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, service.ErrInvalidRequest, err) // Expression validation is part of the request validation
}

func TestCreateStrategy_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	req := createValidCreateStrategyRequest()
	userID := uuid.New()

	mockRepo := &mocks.MockStrategyRepository{
		CreateFunc: func(ctx context.Context, strategy *model.Strategy) error {
			return errors.New("database error")
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.CreateStrategy(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, service.ErrStrategyCreateFailed, err)
}

func TestGetStrategyByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()
	strategy := createTestStrategy(userID)
	strategy.ID = strategyID

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			assert.Equal(t, strategyID, id)
			return strategy, nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.GetStrategyByID(ctx, strategyID, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, strategyID, response.ID)
	assert.Equal(t, userID, response.UserID)
}

func TestGetStrategyByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			return nil, errors.New("strategy not found")
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.GetStrategyByID(ctx, strategyID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, service.ErrStrategyNotFound, err)
}

func TestGetStrategyByID_Unauthorized(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	strategy := createTestStrategy(otherUserID)
	strategy.ID = strategyID

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			return strategy, nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.GetStrategyByID(ctx, strategyID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, service.ErrUnauthorized, err)
}

func TestGetUserStrategies_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	strategies := []*model.Strategy{
		createTestStrategy(userID),
		createTestStrategy(userID),
	}

	mockRepo := &mocks.MockStrategyRepository{
		GetByUserIDFunc: func(ctx context.Context, id uuid.UUID) ([]*model.Strategy, error) {
			assert.Equal(t, userID, id)
			return strategies, nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.GetUserStrategies(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response, 2)
}

func TestGetUserStrategies_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := &mocks.MockStrategyRepository{
		GetByUserIDFunc: func(ctx context.Context, id uuid.UUID) ([]*model.Strategy, error) {
			return nil, errors.New("database error")
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.GetUserStrategies(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestUpdateStrategy_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()
	
	existingStrategy := createTestStrategy(userID)
	existingStrategy.ID = strategyID
	
	req := createValidUpdateStrategyRequest()
	req.ID = strategyID

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			assert.Equal(t, strategyID, id)
			return existingStrategy, nil
		},
		UpdateFunc: func(ctx context.Context, strategy *model.Strategy) error {
			assert.Equal(t, strategyID, strategy.ID)
			assert.Equal(t, userID, strategy.UserID)
			return nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.UpdateStrategy(ctx, req, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, strategyID, response.ID)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, req.Name, response.Name)
}

func TestUpdateStrategy_ValidationError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	
	req := createValidUpdateStrategyRequest()
	req.Name = "" // Invalid name

	mockRepo := &mocks.MockStrategyRepository{}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.UpdateStrategy(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, service.ErrInvalidRequest, err)
}

func TestUpdateStrategy_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()
	
	req := createValidUpdateStrategyRequest()
	req.ID = strategyID

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			return nil, errors.New("strategy not found")
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.UpdateStrategy(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, service.ErrStrategyNotFound, err)
}

func TestUpdateStrategy_Unauthorized(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	
	existingStrategy := createTestStrategy(otherUserID)
	existingStrategy.ID = strategyID
	
	req := createValidUpdateStrategyRequest()
	req.ID = strategyID

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			return existingStrategy, nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.UpdateStrategy(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, service.ErrUnauthorized, err)
}

func TestDeleteStrategy_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()
	
	existingStrategy := createTestStrategy(userID)
	existingStrategy.ID = strategyID

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			assert.Equal(t, strategyID, id)
			return existingStrategy, nil
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			assert.Equal(t, strategyID, id)
			return nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	err := strategyService.DeleteStrategy(ctx, strategyID, userID)

	// Assert
	assert.NoError(t, err)
}

func TestDeleteStrategy_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			return nil, errors.New("strategy not found")
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	err := strategyService.DeleteStrategy(ctx, strategyID, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrStrategyNotFound, err)
}

func TestDeleteStrategy_Unauthorized(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	
	existingStrategy := createTestStrategy(otherUserID)
	existingStrategy.ID = strategyID

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			return existingStrategy, nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	err := strategyService.DeleteStrategy(ctx, strategyID, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrUnauthorized, err)
}

func TestActivateStrategy_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()
	
	existingStrategy := createTestStrategy(userID)
	existingStrategy.ID = strategyID
	existingStrategy.IsActive = false

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			assert.Equal(t, strategyID, id)
			return existingStrategy, nil
		},
		ActivateStrategyFunc: func(ctx context.Context, id uuid.UUID) error {
			assert.Equal(t, strategyID, id)
			return nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	err := strategyService.ActivateStrategy(ctx, strategyID, userID)

	// Assert
	assert.NoError(t, err)
}

func TestActivateStrategy_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			return nil, errors.New("strategy not found")
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	err := strategyService.ActivateStrategy(ctx, strategyID, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrStrategyNotFound, err)
}

func TestActivateStrategy_Unauthorized(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	
	existingStrategy := createTestStrategy(otherUserID)
	existingStrategy.ID = strategyID

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			return existingStrategy, nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	err := strategyService.ActivateStrategy(ctx, strategyID, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrUnauthorized, err)
}

func TestDeactivateStrategy_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	strategyID := uuid.New()
	userID := uuid.New()
	
	existingStrategy := createTestStrategy(userID)
	existingStrategy.ID = strategyID

	mockRepo := &mocks.MockStrategyRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Strategy, error) {
			assert.Equal(t, strategyID, id)
			return existingStrategy, nil
		},
		DeactivateStrategyFunc: func(ctx context.Context, id uuid.UUID) error {
			assert.Equal(t, strategyID, id)
			return nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	err := strategyService.DeactivateStrategy(ctx, strategyID, userID)

	// Assert
	assert.NoError(t, err)
}

func TestSearchStrategies_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	isActive := true
	filter := &dto.StrategyFilter{
		IsActive: &isActive,
	}
	
	strategies := []*model.Strategy{
		createTestStrategy(userID),
		createTestStrategy(userID),
	}

	mockRepo := &mocks.MockStrategyRepository{
		SearchStrategiesFunc: func(ctx context.Context, filters map[string]interface{}) ([]*model.Strategy, error) {
			assert.Contains(t, filters, "user_id")
			assert.Contains(t, filters, "is_active")
			return strategies, nil
		},
	}

	strategyService := service.NewStrategyService(mockRepo)

	// Act
	response, err := strategyService.SearchStrategies(ctx, filter, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response, 2)
}

// Helper functions to create test objects

func createValidCreateStrategyRequest() *dto.CreateStrategyRequest {
	trailingStopLoss := 3.0
	trailingStopLossType := model.Percentage
	exitType := "EOD"
	maxTrades := 5
	entryStartTime := time.Now().Add(-1 * time.Hour)
	entryEndTime := time.Now()
	startDate := time.Now().Add(-7 * 24 * time.Hour)
	endDate := time.Now().Add(7 * 24 * time.Hour)
	targetProfit := 5.0
	targetProfitType := model.Percentage
	stopLoss := 2.0
	stopLossType := model.Percentage
	capital := 10000.0

	return &dto.CreateStrategyRequest{
		Name:                 "Test Strategy",
		Description:          "Strategy for testing",
		PositionType:         model.Long,
		Instruments:          pq.Int32Array{1, 2, 3},
		ChartType:            model.Candlestick,
		ChartTimeframe:       model.OneHour,
		TradeType:            "MIS",
		Qty:                  10,
		EntryExpression:      "close > open",
		ExitExpression:       "close < open",
		TargetProfit:         &targetProfit,
		TargetProfitType:     &targetProfitType,
		StopLoss:             &stopLoss,
		StopLossType:         &stopLossType,
		TrailingStopLoss:     &trailingStopLoss,
		TrailingStopLossType: &trailingStopLossType,
		ExitType:             &exitType,
		Capital:              &capital,
		MaxTrades:            &maxTrades,
		EntryStartTime:       &entryStartTime,
		EntryEndTime:         &entryEndTime,
		StartDate:            &startDate,
		EndDate:              &endDate,
	}
}

func createValidUpdateStrategyRequest() *dto.UpdateStrategyRequest {
	trailingStopLoss := 3.0
	trailingStopLossType := model.Percentage
	exitType := "EOD"
	maxTrades := 5
	entryStartTime := time.Now().Add(-1 * time.Hour)
	entryEndTime := time.Now()
	startDate := time.Now().Add(-7 * 24 * time.Hour)
	endDate := time.Now().Add(7 * 24 * time.Hour)
	targetProfit := 10.0
	targetProfitType := model.Absolute
	stopLoss := 5.0
	stopLossType := model.Absolute
	capital := 20000.0

	return &dto.UpdateStrategyRequest{
		ID:                   uuid.New(),
		Name:                 "Updated Strategy",
		Description:          "Updated Strategy for testing",
		PositionType:         model.Short,
		Instruments:          pq.Int32Array{4, 5, 6},
		ChartType:            model.Heikinashi,
		ChartTimeframe:       model.FiveMinutes,
		TradeType:            "CNC/NRML",
		Qty:                  20,
		EntryExpression:      "macd > signal",
		ExitExpression:       "macd < signal",
		TargetProfit:         &targetProfit,
		TargetProfitType:     &targetProfitType,
		StopLoss:             &stopLoss,
		StopLossType:         &stopLossType,
		TrailingStopLoss:     &trailingStopLoss,
		TrailingStopLossType: &trailingStopLossType,
		ExitType:             &exitType,
		Capital:              &capital,
		MaxTrades:            &maxTrades,
		EntryStartTime:       &entryStartTime,
		EntryEndTime:         &entryEndTime,
		StartDate:            &startDate,
		EndDate:              &endDate,
	}
}

// Helper function to create a test strategy
func createTestStrategy(userID uuid.UUID) *model.Strategy {
	strategy := model.NewStrategy(
		userID,
		"Test Strategy",
		"Test Description",
		model.Long,
		pq.Int32Array{1, 2, 3},
		model.Candlestick,
		model.OneHour,
		"MIS",
		10,
		"close > open",
		"close < open",
		5.0,
		model.Percentage,
		2.0,
		model.Percentage,
		10000.0,
	)

	trailingStopLoss := 3.0
	trailingStopLossType := model.Percentage
	exitType := "EOD"
	maxTrades := 5
	entryStartTime := time.Now().Add(-1 * time.Hour)
	entryEndTime := time.Now()
	startDate := time.Now().Add(-7 * 24 * time.Hour)
	endDate := time.Now().Add(7 * 24 * time.Hour)

	strategy.TrailingStopLoss = &trailingStopLoss
	strategy.TrailingStopLossType = &trailingStopLossType
	strategy.ExitType = &exitType
	strategy.MaxTrades = &maxTrades
	strategy.EntryStartTime = &entryStartTime
	strategy.EntryEndTime = &entryEndTime
	strategy.StartDate = &startDate
	strategy.EndDate = &endDate

	return strategy
}