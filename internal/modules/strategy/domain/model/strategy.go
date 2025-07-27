package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"algo/internal/util/types"

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
	ThirtyMinutes  ChartTimeframe = "30m"
	OneHour        ChartTimeframe = "1h"
	ThreeHours     ChartTimeframe = "3h"
	OneDay         ChartTimeframe = "1d"
)

// Trade type enum
type TradeType string

const (
	MIS      = "MIS"
	CNC_NRML = "CNC/NRML"
)

// Exchange enum
type Exchange string

const (
	NSE Exchange = "NSE"
	BSE Exchange = "BSE"
	MCX Exchange = "MCX"
)

// Segment enum
type Segment string

const (
	FO Segment = "F&O"
	EQ Segment = "EQ"
)

// Asset enum
type Asset string

const (
	Index Asset = "INDEX"
	Stock Asset = "STOCK"
)

// StrategyType enum
type StrategyType string

const (
	Conditional StrategyType = "CONDITIONAL"
	TimeBased   StrategyType = "TIME_BASED"
)

// Value type enum
type ValueType string

const (
	Absolute   ValueType = "ABSOLUTE"
	Percentage ValueType = "PERCENTAGE"
)

// Trading day enum
type TradingDay string

const (
	Monday    TradingDay = "MON"
	Tuesday   TradingDay = "TUE"
	Wednesday TradingDay = "WED"
	Thursday  TradingDay = "THU"
	Friday    TradingDay = "FRI"
)

// TradingDayArray is a custom type for PostgreSQL array of trading days
type TradingDayArray []TradingDay

// Value implements the driver.Valuer interface for database storage
func (tda TradingDayArray) Value() (driver.Value, error) {
	if tda == nil {
		return nil, nil
	}

	values := make([]string, len(tda))
	for i, td := range tda {
		values[i] = string(td)
	}

	return "{" + strings.Join(values, ",") + "}", nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (tda *TradingDayArray) Scan(value interface{}) error {
	if value == nil {
		*tda = nil
		return nil
	}

	switch s := value.(type) {
	case string:
		// Remove braces and split by comma
		s = strings.Trim(s, "{}")
		if s == "" {
			*tda = TradingDayArray{}
			return nil
		}

		parts := strings.Split(s, ",")
		result := make(TradingDayArray, len(parts))
		for i, part := range parts {
			result[i] = TradingDay(strings.TrimSpace(part))
		}
		*tda = result
		return nil
	case []byte:
		return tda.Scan(string(s))
	default:
		return fmt.Errorf("cannot scan %T into TradingDayArray", value)
	}
}

// Strategy represents a trading strategy with entry/exit conditions and risk parameters
type Strategy struct {
	ID                   uuid.UUID       `json:"id" db:"id"`
	UserID               uuid.UUID       `json:"user_id" db:"user_id"`
	Name                 string          `json:"name" db:"name"`
	Description          string          `json:"description" db:"description"`
	PositionType         *Direction      `json:"position_type,omitempty" db:"position_type"`
	Instruments          pq.Int32Array   `json:"instruments" db:"instruments"`
	ChartType            ChartType       `json:"chart_type" db:"chart_type"`
	ChartTimeframe       ChartTimeframe  `json:"chart_timeframe" db:"chart_timeframe"`
	TradingDays          TradingDayArray `json:"trading_days" db:"trading_days"`
	TradeType            *TradeType      `json:"trade_type,omitempty" db:"trade_type"`
	Qty                  *int            `json:"qty,omitempty" db:"qty"`
	EntryExpression      string          `json:"entry_expression" db:"entry_expression"`
	ExitExpression       string          `json:"exit_expression,omitempty" db:"exit_expression"`
	TargetProfit         float64         `json:"target_profit" db:"target_profit"`
	TargetProfitType     ValueType       `json:"target_profit_type" db:"target_profit_type"`
	StopLoss             float64         `json:"stop_loss" db:"stop_loss"`
	StopLossType         ValueType       `json:"stop_loss_type" db:"stop_loss_type"`
	TrailingStopLoss     *float64        `json:"trailing_stop_loss,omitempty" db:"trailing_stop_loss"`
	TrailingStopLossType *ValueType      `json:"trailing_stop_loss_type,omitempty" db:"trailing_stop_loss_type"`
	ExitType             *string         `json:"exit_type,omitempty" db:"exit_type"`
	Capital              *float64        `json:"capital,omitempty" db:"capital"`
	MaxAllocation        *float64        `json:"max_allocation,omitempty" db:"max_allocation"`
	MaxTrades            *int            `json:"max_trades,omitempty" db:"max_trades"`
	Brokerage            *int            `json:"brokerage,omitempty" db:"brokerage"`
	Slippage             *float64        `json:"slippage,omitempty" db:"slippage"`
	EntryStartTime       *types.TimeOnly `json:"entry_start_time,omitempty" db:"entry_start_time"`
	EntryEndTime         *types.TimeOnly `json:"entry_end_time,omitempty" db:"entry_end_time"`
	Segment              Segment         `json:"segment" db:"segment"`
	StrategyType         StrategyType    `json:"strategy_type" db:"strategy_type"`
	EntryTime            *types.TimeOnly `json:"entry_time,omitempty" db:"entry_time"`
	ExitTime             *types.TimeOnly `json:"exit_time,omitempty" db:"exit_time"`
	StartDate            *time.Time      `json:"start_date,omitempty" db:"start_date"`
	EndDate              *time.Time      `json:"end_date,omitempty" db:"end_date"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
	IsActive             bool            `json:"is_active" db:"is_active"`
}

// NewStrategy creates a new strategy with the given parameters
func NewStrategy(
	userID uuid.UUID,
	name string,
	description string,
	positionType *Direction,
	instruments pq.Int32Array,
	chartType ChartType,
	chartTimeframe ChartTimeframe,
	tradingDays TradingDayArray,
	tradeType *TradeType,
	qty *int,
	entryExpression string,
	exitExpression string,
	targetProfit float64,
	targetProfitType ValueType,
	stopLoss float64,
	stopLossType ValueType,
	capital *float64,
	maxAllocation *float64,
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
		TradingDays:      tradingDays,
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
		Segment:          FO,          // Default segment
		StrategyType:     Conditional, // Default strategy type
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsActive:         true,
	}
}

// Update updates the strategy fields
func (s *Strategy) Update(
	name string,
	description string,
	positionType *Direction,
	instruments pq.Int32Array,
	chartType ChartType,
	chartTimeframe ChartTimeframe,
	tradingDays TradingDayArray,
	tradeType *TradeType,
	qty *int,
	entryExpression string,
	exitExpression string,
	targetProfit float64,
	targetProfitType ValueType,
	stopLoss float64,
	stopLossType ValueType,
	trailingStopLoss *float64,
	trailingStopLossType *ValueType,
	exitType *string,
	capital *float64,
	maxAllocation *float64,
	maxTrades *int,
	brokerage *int,
	slippage *float64,
	entryStartTime *types.TimeOnly,
	entryEndTime *types.TimeOnly,
	startDate *time.Time,
	endDate *time.Time,
	segment Segment,
	strategyType StrategyType,
	entryTime *types.TimeOnly,
	exitTime *types.TimeOnly,
) {
	s.Name = name
	s.Description = description
	s.PositionType = positionType
	s.Instruments = instruments
	s.ChartType = chartType
	s.ChartTimeframe = chartTimeframe
	s.TradingDays = tradingDays
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
	s.Brokerage = brokerage
	s.Slippage = slippage
	s.EntryStartTime = entryStartTime
	s.EntryEndTime = entryEndTime
	s.StartDate = startDate
	s.EndDate = endDate
	s.Segment = segment
	s.StrategyType = strategyType
	s.EntryTime = entryTime
	s.ExitTime = exitTime
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
