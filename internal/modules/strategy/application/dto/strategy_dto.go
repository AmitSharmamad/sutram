package dto

import (
	"errors"
	"time"

	"algo/internal/modules/strategy/domain/model"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// validator is a singleton instance of the validator
var validate = validator.New()

// Common errors
var (
	// ErrInvalidID is returned when the provided ID is invalid
	ErrInvalidID = errors.New("invalid ID")

	// ErrInvalidName is returned when the provided name is invalid
	ErrInvalidName = errors.New("invalid name")

	// ErrInvalidDescription is returned when the provided description is invalid
	ErrInvalidDescription = errors.New("invalid description")

	// ErrInvalidPositionType is returned when the provided position type is invalid
	ErrInvalidPositionType = errors.New("invalid position type")

	// ErrInvalidInstruments is returned when the provided instruments are invalid
	ErrInvalidInstruments = errors.New("invalid instruments")

	// ErrInvalidChartType is returned when the provided chart type is invalid
	ErrInvalidChartType = errors.New("invalid chart type")

	// ErrInvalidChartTimeframe is returned when the provided chart timeframe is invalid
	ErrInvalidChartTimeframe = errors.New("invalid chart timeframe")

	// ErrInvalidTradeType is returned when the provided trade type is invalid
	ErrInvalidTradeType = errors.New("invalid trade type")

	// ErrInvalidQty is returned when the provided qty is invalid
	ErrInvalidQty = errors.New("invalid qty")

	// ErrInvalidEntryExpression is returned when the provided entry expression is invalid
	ErrInvalidEntryExpression = errors.New("invalid entry expression")

	// ErrInvalidExitExpression is returned when the provided exit expression is invalid
	ErrInvalidExitExpression = errors.New("invalid exit expression")

	// ErrInvalidTargetProfit is returned when the provided target profit is invalid
	ErrInvalidTargetProfit = errors.New("invalid target profit")

	// ErrInvalidTargetProfitType is returned when the provided target profit type is invalid
	ErrInvalidTargetProfitType = errors.New("invalid target profit type")

	// ErrInvalidStopLoss is returned when the provided stop loss is invalid
	ErrInvalidStopLoss = errors.New("invalid stop loss")

	// ErrInvalidStopLossType is returned when the provided stop loss type is invalid
	ErrInvalidStopLossType = errors.New("invalid stop loss type")

	// ErrInvalidTrailingStopLoss is returned when the provided trailing stop loss is invalid
	ErrInvalidTrailingStopLoss = errors.New("invalid trailing stop loss")

	// ErrInvalidTrailingStopLossType is returned when the provided trailing stop loss type is invalid
	ErrInvalidTrailingStopLossType = errors.New("invalid trailing stop loss type")

	// ErrInvalidExitType is returned when the provided exit type is invalid
	ErrInvalidExitType = errors.New("invalid exit type")

	// ErrInvalidCapital is returned when the provided capital is invalid
	ErrInvalidCapital = errors.New("invalid capital")

	// ErrInvalidMaxTrades is returned when the provided max trades is invalid
	ErrInvalidMaxTrades = errors.New("invalid max trades")

	// ErrInvalidEntryStartTime is returned when the provided entry start time is invalid
	ErrInvalidEntryStartTime = errors.New("invalid entry start time")

	// ErrInvalidEntryEndTime is returned when the provided entry end time is invalid
	ErrInvalidEntryEndTime = errors.New("invalid entry end time")

	// ErrInvalidStartDate is returned when the provided start date is invalid
	ErrInvalidStartDate = errors.New("invalid start date")

	// ErrInvalidEndDate is returned when the provided end date is invalid
	ErrInvalidEndDate = errors.New("invalid end date")
)

// CreateStrategyRequest represents the request to create a new strategy
type CreateStrategyRequest struct {
	Name                 string               `json:"name" validate:"required,min=3,max=100"`
	Description          string               `json:"description" validate:"max=500"`
	PositionType         model.Direction      `json:"position_type" validate:"required,oneof=LONG SHORT BOTH"`
	Instruments          []int                `json:"instruments" validate:"required,min=1,dive,gt=0"`
	ChartType            model.ChartType      `json:"chart_type" validate:"required,oneof=CANDLESTICK HEIKINASHI RENKO"`
	ChartTimeframe       model.ChartTimeframe `json:"chart_timeframe" validate:"required,oneof=1M 5M 15M 30M 1H 4H 1D 1W"`
	TradeType            model.TradeType      `json:"trade_type" validate:"required,oneof=MARKET LIMIT STOP_LOSS"`
	Qty                  int                  `json:"qty" validate:"required,gt=0"`
	EntryExpression      string               `json:"entry_expression" validate:"required"`
	ExitExpression       string               `json:"exit_expression"`
	TargetProfit         *float64             `json:"target_profit" validate:"required,gte=0"`
	TargetProfitType     *model.ValueType     `json:"target_profit_type" validate:"required,oneof=ABSOLUTE PERCENTAGE"`
	StopLoss             *float64             `json:"stop_loss" validate:"required,gte=0"`
	StopLossType         *model.ValueType     `json:"stop_loss_type" validate:"required,oneof=ABSOLUTE PERCENTAGE"`
	TrailingStopLoss     *float64             `json:"trailing_stop_loss" validate:"omitempty,gte=0"`
	TrailingStopLossType *model.ValueType     `json:"trailing_stop_loss_type" validate:"omitempty,oneof=ABSOLUTE PERCENTAGE"`
	ExitType             *string              `json:"exit_type"`
	Capital              *float64             `json:"capital" validate:"required,gt=0"`
	MaxTrades            *int                 `json:"max_trades" validate:"omitempty,gt=0"`
	EntryStartTime       *time.Time           `json:"entry_start_time"`
	EntryEndTime         *time.Time           `json:"entry_end_time"`
	StartDate            *time.Time           `json:"start_date"`
	EndDate              *time.Time           `json:"end_date"`
}

// UpdateStrategyRequest represents the request to update an existing strategy
type UpdateStrategyRequest struct {
	ID                   uuid.UUID            `json:"id" validate:"required"`
	Name                 string               `json:"name" validate:"required,min=3,max=100"`
	Description          string               `json:"description" validate:"max=500"`
	PositionType         model.Direction      `json:"position_type" validate:"required,oneof=LONG SHORT BOTH"`
	Instruments          []int                `json:"instruments" validate:"required,min=1,dive,gt=0"`
	ChartType            model.ChartType      `json:"chart_type" validate:"required,oneof=CANDLESTICK HEIKINASHI RENKO"`
	ChartTimeframe       model.ChartTimeframe `json:"chart_timeframe" validate:"required,oneof=1M 5M 15M 30M 1H 4H 1D 1W"`
	TradeType            model.TradeType      `json:"trade_type" validate:"required,oneof=MARKET LIMIT STOP_LOSS"`
	Qty                  int                  `json:"qty" validate:"required,gt=0"`
	EntryExpression      string               `json:"entry_expression" validate:"required"`
	ExitExpression       string               `json:"exit_expression"`
	TargetProfit         *float64             `json:"target_profit" validate:"required,gte=0"`
	TargetProfitType     *model.ValueType     `json:"target_profit_type" validate:"required,oneof=ABSOLUTE PERCENTAGE"`
	StopLoss             *float64             `json:"stop_loss" validate:"required,gte=0"`
	StopLossType         *model.ValueType     `json:"stop_loss_type" validate:"required,oneof=ABSOLUTE PERCENTAGE"`
	TrailingStopLoss     *float64             `json:"trailing_stop_loss" validate:"omitempty,gte=0"`
	TrailingStopLossType *model.ValueType     `json:"trailing_stop_loss_type" validate:"omitempty,oneof=ABSOLUTE PERCENTAGE"`
	ExitType             *string              `json:"exit_type"`
	Capital              *float64             `json:"capital" validate:"required,gt=0"`
	MaxTrades            *int                 `json:"max_trades" validate:"omitempty,gt=0"`
	EntryStartTime       *time.Time           `json:"entry_start_time"`
	EntryEndTime         *time.Time           `json:"entry_end_time"`
	StartDate            *time.Time           `json:"start_date"`
	EndDate              *time.Time           `json:"end_date"`
}

// StrategyResponse represents the response for a strategy
type StrategyResponse struct {
	ID                   uuid.UUID            `json:"id"`
	UserID               uuid.UUID            `json:"user_id"`
	Name                 string               `json:"name"`
	Description          string               `json:"description"`
	PositionType         model.Direction      `json:"position_type"`
	Instruments          []int                `json:"instruments"`
	ChartType            model.ChartType      `json:"chart_type"`
	ChartTimeframe       model.ChartTimeframe `json:"chart_timeframe"`
	TradeType            model.TradeType      `json:"trade_type"`
	Qty                  int                  `json:"qty"`
	EntryExpression      string               `json:"entry_expression"`
	ExitExpression       string               `json:"exit_expression"`
	TargetProfit         float64              `json:"target_profit,omitempty"`
	TargetProfitType     model.ValueType      `json:"target_profit_type,omitempty"`
	StopLoss             float64              `json:"stop_loss,omitempty"`
	StopLossType         model.ValueType      `json:"stop_loss_type,omitempty"`
	TrailingStopLoss     float64              `json:"trailing_stop_loss,omitempty"`
	TrailingStopLossType model.ValueType      `json:"trailing_stop_loss_type,omitempty"`
	ExitType             string               `json:"exit_type,omitempty"`
	Capital              float64              `json:"capital,omitempty"`
	MaxTrades            int                  `json:"max_trades,omitempty"`
	EntryStartTime       *time.Time           `json:"entry_start_time,omitempty"`
	EntryEndTime         *time.Time           `json:"entry_end_time,omitempty"`
	StartDate            *time.Time           `json:"start_date,omitempty"`
	EndDate              *time.Time           `json:"end_date,omitempty"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
	IsActive             bool                 `json:"is_active"`
}

// StrategyFilter represents filters for searching strategies
type StrategyFilter struct {
	UserID         *uuid.UUID            `json:"user_id"`
	Name           *string               `json:"name"`
	PositionType   *model.Direction      `json:"position_type"`
	Instrument     *int                  `json:"instrument"`
	ChartType      *model.ChartType      `json:"chart_type"`
	ChartTimeframe *model.ChartTimeframe `json:"chart_timeframe"`
	IsActive       *bool                 `json:"is_active"`
	CreatedAfter   *time.Time            `json:"created_after"`
	CreatedBefore  *time.Time            `json:"created_before"`
}

// FromDomain converts a domain model to a response DTO
func FromDomain(strategy *model.Strategy) *StrategyResponse {
	response := &StrategyResponse{
		ID:               strategy.ID,
		UserID:           strategy.UserID,
		Name:             strategy.Name,
		Description:      strategy.Description,
		PositionType:     strategy.PositionType,
		Instruments:      strategy.Instruments,
		ChartType:        strategy.ChartType,
		ChartTimeframe:   strategy.ChartTimeframe,
		TradeType:        strategy.TradeType,
		Qty:              strategy.Qty,
		EntryExpression:  strategy.EntryExpression,
		ExitExpression:   strategy.ExitExpression,
		TargetProfit:     strategy.TargetProfit,
		TargetProfitType: strategy.TargetProfitType,
		StopLoss:         strategy.StopLoss,
		StopLossType:     strategy.StopLossType,
		Capital:          strategy.Capital,
		CreatedAt:        strategy.CreatedAt,
		UpdatedAt:        strategy.UpdatedAt,
		IsActive:         strategy.IsActive,
	}

	// Handle nullable fields
	if strategy.TrailingStopLoss != nil {
		response.TrailingStopLoss = *strategy.TrailingStopLoss
	}

	if strategy.TrailingStopLossType != nil {
		response.TrailingStopLossType = *strategy.TrailingStopLossType
	}

	if strategy.ExitType != nil {
		response.ExitType = *strategy.ExitType
	}

	if strategy.MaxTrades != nil {
		response.MaxTrades = *strategy.MaxTrades
	}

	response.EntryStartTime = strategy.EntryStartTime
	response.EntryEndTime = strategy.EntryEndTime
	response.StartDate = strategy.StartDate
	response.EndDate = strategy.EndDate

	return response
}

// ToDomain converts a request DTO to a domain model
func (req *CreateStrategyRequest) ToDomain(userID uuid.UUID) *model.Strategy {
	// Set default values for required fields if they're nil
	targetProfit := 0.0
	if req.TargetProfit != nil {
		targetProfit = *req.TargetProfit
	}

	targetProfitType := model.Percentage
	if req.TargetProfitType != nil {
		targetProfitType = *req.TargetProfitType
	}

	stopLoss := 0.0
	if req.StopLoss != nil {
		stopLoss = *req.StopLoss
	}

	stopLossType := model.Percentage
	if req.StopLossType != nil {
		stopLossType = *req.StopLossType
	}

	capital := 0.0
	if req.Capital != nil {
		capital = *req.Capital
	}

	// Create strategy with required parameters
	strategy := model.NewStrategy(
		userID,
		req.Name,
		req.Description,
		req.PositionType,
		req.Instruments,
		req.ChartType,
		req.ChartTimeframe,
		req.TradeType,
		req.Qty,
		req.EntryExpression,
		req.ExitExpression,
		targetProfit,
		targetProfitType,
		stopLoss,
		stopLossType,
		capital,
	)

	// Set optional fields
	if req.TrailingStopLoss != nil {
		strategy.TrailingStopLoss = req.TrailingStopLoss
		if req.TrailingStopLossType != nil {
			strategy.TrailingStopLossType = req.TrailingStopLossType
		}
	}

	if req.ExitType != nil {
		strategy.ExitType = req.ExitType
	}

	if req.MaxTrades != nil {
		strategy.MaxTrades = req.MaxTrades
	}

	if req.EntryStartTime != nil {
		strategy.EntryStartTime = req.EntryStartTime
	}

	if req.EntryEndTime != nil {
		strategy.EntryEndTime = req.EntryEndTime
	}

	if req.StartDate != nil {
		strategy.StartDate = req.StartDate
	}

	if req.EndDate != nil {
		strategy.EndDate = req.EndDate
	}

	return strategy
}

// ToDomain converts a request DTO to a domain model
func (req *UpdateStrategyRequest) ToDomain() *model.Strategy {
	// Set default values for required fields if they're nil
	targetProfit := 0.0
	if req.TargetProfit != nil {
		targetProfit = *req.TargetProfit
	}

	targetProfitType := model.Percentage
	if req.TargetProfitType != nil {
		targetProfitType = *req.TargetProfitType
	}

	stopLoss := 0.0
	if req.StopLoss != nil {
		stopLoss = *req.StopLoss
	}

	stopLossType := model.Percentage
	if req.StopLossType != nil {
		stopLossType = *req.StopLossType
	}

	capital := 0.0
	if req.Capital != nil {
		capital = *req.Capital
	}

	// Create the strategy object with base fields
	strategy := &model.Strategy{
		ID:               req.ID,
		Name:             req.Name,
		Description:      req.Description,
		PositionType:     req.PositionType,
		Instruments:      req.Instruments,
		ChartType:        req.ChartType,
		ChartTimeframe:   req.ChartTimeframe,
		TradeType:        req.TradeType,
		Qty:              req.Qty,
		EntryExpression:  req.EntryExpression,
		ExitExpression:   req.ExitExpression,
		TargetProfit:     targetProfit,
		TargetProfitType: targetProfitType,
		StopLoss:         stopLoss,
		StopLossType:     stopLossType,
		Capital:          capital,
	}

	// Set optional fields
	if req.TrailingStopLoss != nil {
		strategy.TrailingStopLoss = req.TrailingStopLoss
	}

	if req.TrailingStopLossType != nil {
		strategy.TrailingStopLossType = req.TrailingStopLossType
	}

	if req.ExitType != nil {
		strategy.ExitType = req.ExitType
	}

	if req.MaxTrades != nil {
		strategy.MaxTrades = req.MaxTrades
	}

	if req.EntryStartTime != nil {
		strategy.EntryStartTime = req.EntryStartTime
	}

	if req.EntryEndTime != nil {
		strategy.EntryEndTime = req.EntryEndTime
	}

	if req.StartDate != nil {
		strategy.StartDate = req.StartDate
	}

	if req.EndDate != nil {
		strategy.EndDate = req.EndDate
	}

	return strategy
}

// Validate validates the CreateStrategyRequest
func (req *CreateStrategyRequest) Validate() error {
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

// Validate validates the UpdateStrategyRequest
func (req *UpdateStrategyRequest) Validate() error {
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}
