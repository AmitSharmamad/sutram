package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Position type enum (direction)
type Direction string

const (
	Long  Direction = "LONG"
	Short Direction = "SHORT"
	Both  Direction = "BOTH"
)

// Chart type enum
type ChartType string

const (
	Candlestick ChartType = "CANDLESTICK"
	Heikinashi  ChartType = "HEIKINASHI"
	Renko       ChartType = "RENKO"
)

// Chart timeframe enum
type ChartTimeframe string

const (
	OneMinute      ChartTimeframe = "1m"
	ThreeMinutes   ChartTimeframe = "3m"
	FiveMinutes    ChartTimeframe = "5m"
	TenMinutes     ChartTimeframe = "10m"
	FifteenMinutes ChartTimeframe = "15m"
	OneHour        ChartTimeframe = "1h"
	OneDay         ChartTimeframe = "1d"
)

// Trade type enum
type TradeType string

const (
	MIS      = "MIS"
	CNC_NRML = "CNC/NRML"
)

// Value type enum
type ValueType string

const (
	Absolute   ValueType = "ABSOLUTE"
	Percentage ValueType = "PERCENTAGE"
)

// Strategy represents a trading strategy with entry/exit conditions and risk parameters
type Strategy struct {
	ID                   uuid.UUID      `json:"id" db:"id"`
	UserID               uuid.UUID      `json:"user_id" db:"user_id"`
	Name                 string         `json:"name" db:"name"`
	Description          string         `json:"description" db:"description"`
	PositionType         Direction      `json:"position_type" db:"position_type"`
	Instruments          pq.Int32Array  `json:"instruments" db:"instruments"`
	ChartType            ChartType      `json:"chart_type" db:"chart_type"`
	ChartTimeframe       ChartTimeframe `json:"chart_timeframe" db:"chart_timeframe"`
	TradeType            TradeType      `json:"trade_type" db:"trade_type"`
	Qty                  int            `json:"qty" db:"qty"`
	EntryExpression      string         `json:"entry_expression" db:"entry_expression"`
	ExitExpression       string         `json:"exit_expression,omitempty" db:"exit_expression"`
	TargetProfit         float64        `json:"target_profit" db:"target_profit"`
	TargetProfitType     ValueType      `json:"target_profit_type" db:"target_profit_type"`
	StopLoss             float64        `json:"stop_loss" db:"stop_loss"`
	StopLossType         ValueType      `json:"stop_loss_type" db:"stop_loss_type"`
	TrailingStopLoss     *float64       `json:"trailing_stop_loss,omitempty" db:"trailing_stop_loss"`
	TrailingStopLossType *ValueType     `json:"trailing_stop_loss_type,omitempty" db:"trailing_stop_loss_type"`
	ExitType             *string        `json:"exit_type,omitempty" db:"exit_type"`
	Capital              float64        `json:"capital" db:"capital"`
	MaxAllocation        float64        `json:"max_allocation" db:"max_allocation"`
	MaxTrades            *int           `json:"max_trades,omitempty" db:"max_trades"`
	EntryStartTime       *time.Time     `json:"entry_start_time,omitempty" db:"entry_start_time"`
	EntryEndTime         *time.Time     `json:"entry_end_time,omitempty" db:"entry_end_time"`
	StartDate            *time.Time     `json:"start_date,omitempty" db:"start_date"`
	EndDate              *time.Time     `json:"end_date,omitempty" db:"end_date"`
	CreatedAt            time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at" db:"updated_at"`
	IsActive             bool           `json:"is_active" db:"is_active"`
}

// NewStrategy creates a new strategy with the given parameters
func NewStrategy(
	userID uuid.UUID,
	name string,
	description string,
	positionType Direction,
	instruments pq.Int32Array,
	chartType ChartType,
	chartTimeframe ChartTimeframe,
	tradeType TradeType,
	qty int,
	entryExpression string,
	exitExpression string,
	targetProfit float64,
	targetProfitType ValueType,
	stopLoss float64,
	stopLossType ValueType,
	capital float64,
	maxAllocation float64,
) *Strategy {
	return &Strategy{
		ID:               uuid.New(),
		UserID:           userID,
		Name:             name,
		Description:      description,
		PositionType:     positionType,
		Instruments:      instruments,
		ChartType:        chartType,
		ChartTimeframe:   chartTimeframe,
		TradeType:        tradeType,
		Qty:              qty,
		EntryExpression:  entryExpression,
		ExitExpression:   exitExpression,
		TargetProfit:     targetProfit,
		TargetProfitType: targetProfitType,
		StopLoss:         stopLoss,
		StopLossType:     stopLossType,
		Capital:          capital,
		MaxAllocation:    maxAllocation,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsActive:         true,
	}
}

// Update updates the strategy fields
func (s *Strategy) Update(
	name string,
	description string,
	positionType Direction,
	instruments pq.Int32Array,
	chartType ChartType,
	chartTimeframe ChartTimeframe,
	tradeType TradeType,
	qty int,
	entryExpression string,
	exitExpression string,
	targetProfit float64,
	targetProfitType ValueType,
	stopLoss float64,
	stopLossType ValueType,
	trailingStopLoss *float64,
	trailingStopLossType *ValueType,
	exitType *string,
	capital float64,
	maxAllocation float64,
	maxTrades *int,
	entryStartTime *time.Time,
	entryEndTime *time.Time,
	startDate *time.Time,
	endDate *time.Time,
) {
	s.Name = name
	s.Description = description
	s.PositionType = positionType
	s.Instruments = instruments
	s.ChartType = chartType
	s.ChartTimeframe = chartTimeframe
	s.TradeType = tradeType
	s.Qty = qty
	s.EntryExpression = entryExpression
	s.ExitExpression = exitExpression
	s.TargetProfit = targetProfit
	s.TargetProfitType = targetProfitType
	s.StopLoss = stopLoss
	s.StopLossType = stopLossType
	s.TrailingStopLoss = trailingStopLoss
	s.TrailingStopLossType = trailingStopLossType
	s.ExitType = exitType
	s.Capital = capital
	s.MaxAllocation = maxAllocation
	s.MaxTrades = maxTrades
	s.EntryStartTime = entryStartTime
	s.EntryEndTime = entryEndTime
	s.StartDate = startDate
	s.EndDate = endDate
	s.UpdatedAt = time.Now()
}

// Activate activates the strategy
func (s *Strategy) Activate() {
	s.IsActive = true
	s.UpdatedAt = time.Now()
}

// Deactivate deactivates the strategy
func (s *Strategy) Deactivate() {
	s.IsActive = false
	s.UpdatedAt = time.Now()
}
