package dto_test

import (
	"testing"
	"time"

	"algo/internal/modules/strategy/application/dto"
	"algo/internal/modules/strategy/domain/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateStrategyRequestValidate_ValidRequest(t *testing.T) {
	// Arrange
	req := createValidCreateStrategyRequest()

	// Act
	err := req.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestCreateStrategyRequestValidate_InvalidRequest(t *testing.T) {
	testCases := []struct {
		name        string
		modifyReq   func(*dto.CreateStrategyRequest)
		expectError bool
	}{
		{
			name: "Empty name",
			modifyReq: func(req *dto.CreateStrategyRequest) {
				req.Name = ""
			},
			expectError: true,
		},
		{
			name: "Name too short",
			modifyReq: func(req *dto.CreateStrategyRequest) {
				req.Name = "AB"
			},
			expectError: true,
		},
		{
			name: "Invalid position type",
			modifyReq: func(req *dto.CreateStrategyRequest) {
				req.PositionType = "INVALID"
			},
			expectError: true,
		},
		{
			name: "Empty instruments",
			modifyReq: func(req *dto.CreateStrategyRequest) {
				req.Instruments = pq.Int32Array{}
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := createValidCreateStrategyRequest()
			tc.modifyReq(req)

			// Act
			err := req.Validate()

			// Assert
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateStrategyRequestValidate_ValidRequest(t *testing.T) {
	// Arrange
	req := createValidUpdateStrategyRequest()

	// Act
	err := req.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestUpdateStrategyRequestValidate_InvalidRequest(t *testing.T) {
	testCases := []struct {
		name        string
		modifyReq   func(*dto.UpdateStrategyRequest)
		expectError bool
	}{
		{
			name: "Empty ID",
			modifyReq: func(req *dto.UpdateStrategyRequest) {
				req.ID = uuid.Nil
			},
			expectError: true,
		},
		{
			name: "Empty name",
			modifyReq: func(req *dto.UpdateStrategyRequest) {
				req.Name = ""
			},
			expectError: true,
		},
		{
			name: "Name too short",
			modifyReq: func(req *dto.UpdateStrategyRequest) {
				req.Name = "AB"
			},
			expectError: true,
		},
		{
			name: "Invalid position type",
			modifyReq: func(req *dto.UpdateStrategyRequest) {
				req.PositionType = "INVALID"
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := createValidUpdateStrategyRequest()
			tc.modifyReq(req)

			// Act
			err := req.Validate()

			// Assert
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateStrategyRequestToDomain(t *testing.T) {
	// Arrange
	req := createValidCreateStrategyRequest()
	userID := uuid.New()

	// Act
	strategy := req.ToDomain(userID)

	// Assert
	assert.NotNil(t, strategy)
	assert.Equal(t, userID, strategy.UserID)
	assert.Equal(t, req.Name, strategy.Name)
	assert.Equal(t, req.Description, strategy.Description)
	assert.Equal(t, req.PositionType, strategy.PositionType)
	assert.Equal(t, req.Instruments, strategy.Instruments)
	assert.Equal(t, req.ChartType, strategy.ChartType)
	assert.Equal(t, req.ChartTimeframe, strategy.ChartTimeframe)
	assert.Equal(t, req.TradeType, strategy.TradeType)
	assert.Equal(t, req.Qty, strategy.Qty)
	assert.Equal(t, req.EntryExpression, strategy.EntryExpression)
	assert.Equal(t, req.ExitExpression, strategy.ExitExpression)
	// ToDomain derferences the pointers
	assert.Equal(t, *req.TargetProfit, strategy.TargetProfit)
	assert.Equal(t, *req.TargetProfitType, strategy.TargetProfitType)
	assert.Equal(t, *req.StopLoss, strategy.StopLoss)
	assert.Equal(t, *req.StopLossType, strategy.StopLossType)
	assert.Equal(t, req.TrailingStopLoss, strategy.TrailingStopLoss)
	assert.Equal(t, req.TrailingStopLossType, strategy.TrailingStopLossType)
	assert.Equal(t, req.ExitType, strategy.ExitType)
	assert.Equal(t, *req.Capital, strategy.Capital)
	assert.Equal(t, req.MaxTrades, strategy.MaxTrades)
	assert.Equal(t, req.EntryStartTime, strategy.EntryStartTime)
	assert.Equal(t, req.EntryEndTime, strategy.EntryEndTime)
	assert.Equal(t, req.StartDate, strategy.StartDate)
	assert.Equal(t, req.EndDate, strategy.EndDate)
	assert.NotEmpty(t, strategy.ID)
	assert.True(t, strategy.IsActive)
}

func TestUpdateStrategyRequestToDomain(t *testing.T) {
	// Arrange
	req := createValidUpdateStrategyRequest()

	// Act
	strategy := req.ToDomain()

	// Assert
	assert.NotNil(t, strategy)
	assert.Equal(t, req.ID, strategy.ID)
	assert.Equal(t, req.Name, strategy.Name)
	assert.Equal(t, req.Description, strategy.Description)
	assert.Equal(t, req.PositionType, strategy.PositionType)
	assert.Equal(t, req.Instruments, strategy.Instruments)
	assert.Equal(t, req.ChartType, strategy.ChartType)
	assert.Equal(t, req.ChartTimeframe, strategy.ChartTimeframe)
	assert.Equal(t, req.TradeType, strategy.TradeType)
	assert.Equal(t, req.Qty, strategy.Qty)
	assert.Equal(t, req.EntryExpression, strategy.EntryExpression)
	assert.Equal(t, req.ExitExpression, strategy.ExitExpression)
	// ToDomain derferences the pointers
	assert.Equal(t, *req.TargetProfit, strategy.TargetProfit)
	assert.Equal(t, *req.TargetProfitType, strategy.TargetProfitType)
	assert.Equal(t, *req.StopLoss, strategy.StopLoss)
	assert.Equal(t, *req.StopLossType, strategy.StopLossType)
	assert.Equal(t, req.TrailingStopLoss, strategy.TrailingStopLoss)
	assert.Equal(t, req.TrailingStopLossType, strategy.TrailingStopLossType)
	assert.Equal(t, req.ExitType, strategy.ExitType)
	assert.Equal(t, *req.Capital, strategy.Capital)
	assert.Equal(t, req.MaxTrades, strategy.MaxTrades)
	assert.Equal(t, req.EntryStartTime, strategy.EntryStartTime)
	assert.Equal(t, req.EntryEndTime, strategy.EntryEndTime)
	assert.Equal(t, req.StartDate, strategy.StartDate)
	assert.Equal(t, req.EndDate, strategy.EndDate)
	// UserID, CreatedAt, UpdatedAt, and IsActive are set by the service
}

func TestFromDomain(t *testing.T) {
	// Arrange
	strategy := createTestStrategy()

	// Act
	response := dto.FromDomain(strategy)

	// Assert
	assert.NotNil(t, response)
	assert.Equal(t, strategy.ID, response.ID)
	assert.Equal(t, strategy.UserID, response.UserID)
	assert.Equal(t, strategy.Name, response.Name)
	assert.Equal(t, strategy.Description, response.Description)
	assert.Equal(t, strategy.PositionType, response.PositionType)
	assert.Equal(t, strategy.Instruments, response.Instruments)
	assert.Equal(t, strategy.ChartType, response.ChartType)
	assert.Equal(t, strategy.ChartTimeframe, response.ChartTimeframe)
	assert.Equal(t, strategy.TradeType, response.TradeType)
	assert.Equal(t, strategy.Qty, response.Qty)
	assert.Equal(t, strategy.EntryExpression, response.EntryExpression)
	assert.Equal(t, strategy.ExitExpression, response.ExitExpression)
	assert.Equal(t, strategy.TargetProfit, response.TargetProfit)
	assert.Equal(t, strategy.TargetProfitType, response.TargetProfitType)
	assert.Equal(t, strategy.StopLoss, response.StopLoss)
	assert.Equal(t, strategy.StopLossType, response.StopLossType)
	// FromDomain copies pointer value to direct value
	if strategy.TrailingStopLoss != nil {
		assert.Equal(t, *strategy.TrailingStopLoss, response.TrailingStopLoss)
	}
	if strategy.TrailingStopLossType != nil {
		assert.Equal(t, *strategy.TrailingStopLossType, response.TrailingStopLossType)
	}
	if strategy.ExitType != nil {
		assert.Equal(t, *strategy.ExitType, response.ExitType)
	}
	assert.Equal(t, strategy.Capital, response.Capital)
	assert.Equal(t, strategy.MaxAllocation, response.MaxAllocation)
	if strategy.MaxTrades != nil {
		assert.Equal(t, *strategy.MaxTrades, response.MaxTrades)
	}
	// Date pointer fields are preserved without dereferencing
	assert.Equal(t, strategy.EntryStartTime, response.EntryStartTime)
	assert.Equal(t, strategy.EntryEndTime, response.EntryEndTime)
	assert.Equal(t, strategy.StartDate, response.StartDate)
	assert.Equal(t, strategy.EndDate, response.EndDate)
	assert.Equal(t, strategy.CreatedAt, response.CreatedAt)
	assert.Equal(t, strategy.UpdatedAt, response.UpdatedAt)
	assert.Equal(t, strategy.IsActive, response.IsActive)
}

// Helper functions to create test objects

func createValidCreateStrategyRequest() *dto.CreateStrategyRequest {
	trailingStopLoss := 3.0
	trailingStopLossType := model.Percentage
	exitType := "EOD"
	maxTrades := 5
	maxAllocation := 8000.0
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
		MaxAllocation:        &maxAllocation,
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
	maxAllocation := 15000.0
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
		MaxAllocation:        &maxAllocation,
		MaxTrades:            &maxTrades,
		EntryStartTime:       &entryStartTime,
		EntryEndTime:         &entryEndTime,
		StartDate:            &startDate,
		EndDate:              &endDate,
	}
}

// Helper function to create a test strategy
func createTestStrategy() *model.Strategy {
	strategy := model.NewStrategy(
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
		0.0, // Default maxAllocation
	)

	trailingStopLoss := 3.0
	trailingStopLossType := model.Percentage
	exitType := "EOD"
	maxAllocation := 8000.0
	maxTrades := 5
	entryStartTime := time.Now().Add(-1 * time.Hour)
	entryEndTime := time.Now()
	startDate := time.Now().Add(-7 * 24 * time.Hour)
	endDate := time.Now().Add(7 * 24 * time.Hour)

	strategy.TrailingStopLoss = &trailingStopLoss
	strategy.TrailingStopLossType = &trailingStopLossType
	strategy.ExitType = &exitType
	strategy.MaxAllocation = maxAllocation
	strategy.MaxTrades = &maxTrades
	strategy.EntryStartTime = &entryStartTime
	strategy.EntryEndTime = &entryEndTime
	strategy.StartDate = &startDate
	strategy.EndDate = &endDate

	return strategy
}
