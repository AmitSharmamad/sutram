package model_test

import (
	"testing"
	"time"

	"algo/internal/modules/strategy/domain/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNewStrategy(t *testing.T) {
	// Arrange
	userID := uuid.New()
	name := "Test Strategy"
	description := "Strategy for testing"
	positionType := model.Long
	instruments := pq.Int32Array{1, 2, 3}
	chartType := model.Candlestick
	chartTimeframe := model.OneHour
	tradeType := model.TradeType("MIS")
	qty := 10
	entryExpression := "close > open"
	exitExpression := "close < open"
	targetProfit := 5.0
	targetProfitType := model.Percentage
	stopLoss := 2.0
	stopLossType := model.Percentage
	capital := 10000.0

	// Act
	strategy := model.NewStrategy(
		userID,
		name,
		description,
		positionType,
		instruments,
		chartType,
		chartTimeframe,
		tradeType,
		qty,
		entryExpression,
		exitExpression,
		targetProfit,
		targetProfitType,
		stopLoss,
		stopLossType,
		capital,
	)

	// Assert
	assert.NotNil(t, strategy)
	assert.Equal(t, userID, strategy.UserID)
	assert.Equal(t, name, strategy.Name)
	assert.Equal(t, description, strategy.Description)
	assert.Equal(t, positionType, strategy.PositionType)
	assert.Equal(t, instruments, strategy.Instruments)
	assert.Equal(t, chartType, strategy.ChartType)
	assert.Equal(t, chartTimeframe, strategy.ChartTimeframe)
	assert.Equal(t, tradeType, strategy.TradeType)
	assert.Equal(t, qty, strategy.Qty)
	assert.Equal(t, entryExpression, strategy.EntryExpression)
	assert.Equal(t, exitExpression, strategy.ExitExpression)
	assert.Equal(t, targetProfit, strategy.TargetProfit)
	assert.Equal(t, targetProfitType, strategy.TargetProfitType)
	assert.Equal(t, stopLoss, strategy.StopLoss)
	assert.Equal(t, stopLossType, strategy.StopLossType)
	assert.Equal(t, capital, strategy.Capital)
	assert.True(t, strategy.IsActive)
	assert.NotEmpty(t, strategy.ID)
	assert.False(t, strategy.CreatedAt.IsZero())
	assert.False(t, strategy.UpdatedAt.IsZero())
}

func TestUpdateStrategy(t *testing.T) {
	// Arrange
	strategy := createTestStrategy()
	originalCreatedAt := strategy.CreatedAt
	
	// New values for update
	newName := "Updated Strategy"
	newDescription := "Updated description"
	newPositionType := model.Short
	newInstruments := pq.Int32Array{4, 5, 6}
	newChartType := model.Heikinashi
	newChartTimeframe := model.FiveMinutes
	newTradeType := model.TradeType("CNC/NRML")
	newQty := 20
	newEntryExpression := "macd > signal"
	newExitExpression := "macd < signal"
	newTargetProfit := 10.0
	newTargetProfitType := model.Absolute
	newStopLoss := 5.0
	newStopLossType := model.Absolute
	trailingStopLoss := 3.0
	trailingStopLossType := model.Percentage
	exitType := "EOD"
	newCapital := 20000.0
	maxTrades := 5
	entryStartTime := time.Now().Add(-1 * time.Hour)
	entryEndTime := time.Now()
	startDate := time.Now().Add(-7 * 24 * time.Hour)
	endDate := time.Now().Add(7 * 24 * time.Hour)

	// Act
	strategy.Update(
		newName,
		newDescription,
		newPositionType,
		newInstruments,
		newChartType,
		newChartTimeframe,
		newTradeType,
		newQty,
		newEntryExpression,
		newExitExpression,
		newTargetProfit,
		newTargetProfitType,
		newStopLoss,
		newStopLossType,
		&trailingStopLoss,
		&trailingStopLossType,
		&exitType,
		newCapital,
		&maxTrades,
		&entryStartTime,
		&entryEndTime,
		&startDate,
		&endDate,
	)

	// Assert
	assert.Equal(t, newName, strategy.Name)
	assert.Equal(t, newDescription, strategy.Description)
	assert.Equal(t, newPositionType, strategy.PositionType)
	assert.Equal(t, newInstruments, strategy.Instruments)
	assert.Equal(t, newChartType, strategy.ChartType)
	assert.Equal(t, newChartTimeframe, strategy.ChartTimeframe)
	assert.Equal(t, newTradeType, strategy.TradeType)
	assert.Equal(t, newQty, strategy.Qty)
	assert.Equal(t, newEntryExpression, strategy.EntryExpression)
	assert.Equal(t, newExitExpression, strategy.ExitExpression)
	assert.Equal(t, newTargetProfit, strategy.TargetProfit)
	assert.Equal(t, newTargetProfitType, strategy.TargetProfitType)
	assert.Equal(t, newStopLoss, strategy.StopLoss)
	assert.Equal(t, newStopLossType, strategy.StopLossType)
	assert.Equal(t, &trailingStopLoss, strategy.TrailingStopLoss)
	assert.Equal(t, &trailingStopLossType, strategy.TrailingStopLossType)
	assert.Equal(t, &exitType, strategy.ExitType)
	assert.Equal(t, newCapital, strategy.Capital)
	assert.Equal(t, &maxTrades, strategy.MaxTrades)
	assert.Equal(t, &entryStartTime, strategy.EntryStartTime)
	assert.Equal(t, &entryEndTime, strategy.EntryEndTime)
	assert.Equal(t, &startDate, strategy.StartDate)
	assert.Equal(t, &endDate, strategy.EndDate)
	assert.Equal(t, originalCreatedAt, strategy.CreatedAt)
	assert.True(t, strategy.UpdatedAt.After(originalCreatedAt) || strategy.UpdatedAt.Equal(originalCreatedAt))
}

func TestActivateStrategy(t *testing.T) {
	// Arrange
	strategy := createTestStrategy()
	strategy.IsActive = false
	originalUpdatedAt := strategy.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure updated time is different

	// Act
	strategy.Activate()

	// Assert
	assert.True(t, strategy.IsActive)
	assert.True(t, strategy.UpdatedAt.After(originalUpdatedAt))
}

func TestDeactivateStrategy(t *testing.T) {
	// Arrange
	strategy := createTestStrategy()
	originalUpdatedAt := strategy.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure updated time is different

	// Act
	strategy.Deactivate()

	// Assert
	assert.False(t, strategy.IsActive)
	assert.True(t, strategy.UpdatedAt.After(originalUpdatedAt))
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