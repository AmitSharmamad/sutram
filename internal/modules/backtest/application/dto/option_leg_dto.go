package dto

import (
	"time"

	"algo/internal/modules/backtest/domain/model"

	"github.com/google/uuid"
)

// CreateOptionLegRequest represents the request to create an option leg
type CreateOptionLegRequest struct {
	LegSequence  int                   `json:"leg_sequence" validate:"required,gt=0"`
	SelectIndex  string                `json:"select_index" validate:"required"`
	Segment      model.Segment         `json:"segment" validate:"required,oneof=F&O EQ"`
	OptionType   *model.OptionType     `json:"option_type,omitempty" validate:"omitempty,oneof=CALL PUT"`
	ActionType   model.TransactionType `json:"action_type" validate:"required,oneof=BUY SELL"`
	StrikeType   model.StrikePriceType `json:"strike_type" validate:"required,oneof=OTM20 OTM19 OTM18 OTM17 OTM16 OTM15 OTM14 OTM13 OTM12 OTM11 OTM10 OTM9 OTM8 OTM7 OTM6 OTM5 OTM4 OTM3 OTM2 OTM1 ATM ITM1 ITM2 ITM3 ITM4 ITM5 ITM6 ITM7 ITM8 ITM9 ITM10 ITM11 ITM12 ITM13 ITM14 ITM15 ITM16 ITM17 ITM18 ITM19 ITM20"`
	StrikePrice  *float64              `json:"strike_price,omitempty"`
	StrikeOffset *int                  `json:"strike_offset,omitempty"`
	TotalLot     int                   `json:"total_lot" validate:"required,gt=0"`
	LotSize      *int                  `json:"lot_size,omitempty"`
	OrderType    model.OrderType       `json:"order_type" validate:"required,oneof=MARKET LIMIT"`
	TickOffset   *float64              `json:"tick_offset,omitempty"`
	ExpiryType   model.ExpiryType      `json:"expiry_type" validate:"required,oneof=Weekly Monthly Quarterly"`
	ExpiryCycle  model.ExpiryCycle     `json:"expiry_cycle" validate:"required,oneof=CURRENT NEXT NEXT2"`
	ExpiryDate   *time.Time            `json:"expiry_date,omitempty"`
	DaysToExpiry *int                  `json:"days_to_expiry,omitempty"`
}

// UpdateOptionLegRequest represents the request to update an option leg
type UpdateOptionLegRequest struct {
	LegSequence  *int                   `json:"leg_sequence,omitempty" validate:"omitempty,gt=0"`
	SelectIndex  *string                `json:"select_index,omitempty"`
	Segment      *model.Segment         `json:"segment,omitempty" validate:"omitempty,oneof=F&O EQ"`
	OptionType   *model.OptionType      `json:"option_type,omitempty" validate:"omitempty,oneof=CALL PUT"`
	ActionType   *model.TransactionType `json:"action_type,omitempty" validate:"omitempty,oneof=BUY SELL"`
	StrikeType   *model.StrikePriceType `json:"strike_type,omitempty" validate:"omitempty,oneof=OTM20 OTM19 OTM18 OTM17 OTM16 OTM15 OTM14 OTM13 OTM12 OTM11 OTM10 OTM9 OTM8 OTM7 OTM6 OTM5 OTM4 OTM3 OTM2 OTM1 ATM ITM1 ITM2 ITM3 ITM4 ITM5 ITM6 ITM7 ITM8 ITM9 ITM10 ITM11 ITM12 ITM13 ITM14 ITM15 ITM16 ITM17 ITM18 ITM19 ITM20"`
	StrikePrice  *float64               `json:"strike_price,omitempty"`
	StrikeOffset *int                   `json:"strike_offset,omitempty"`
	TotalLot     *int                   `json:"total_lot,omitempty" validate:"omitempty,gt=0"`
	LotSize      *int                   `json:"lot_size,omitempty"`
	OrderType    *model.OrderType       `json:"order_type,omitempty" validate:"omitempty,oneof=MARKET LIMIT"`
	TickOffset   *float64               `json:"tick_offset,omitempty"`
	ExpiryType   *model.ExpiryType      `json:"expiry_type,omitempty" validate:"omitempty,oneof=Weekly Monthly Quarterly"`
	ExpiryCycle  *model.ExpiryCycle     `json:"expiry_cycle,omitempty" validate:"omitempty,oneof=CURRENT NEXT NEXT2"`
	ExpiryDate   *time.Time             `json:"expiry_date,omitempty"`
	DaysToExpiry *int                   `json:"days_to_expiry,omitempty"`
}

// OptionLegResponse represents the response for option leg operations
type OptionLegResponse struct {
	ID            uuid.UUID             `json:"id"`
	StrategyID    uuid.UUID             `json:"strategy_id"`
	LegSequence   int                   `json:"leg_sequence"`
	SelectIndex   string                `json:"select_index"`
	Segment       model.Segment         `json:"segment"`
	OptionType    *model.OptionType     `json:"option_type,omitempty"`
	ActionType    model.TransactionType `json:"action_type"`
	StrikeType    model.StrikePriceType `json:"strike_type"`
	StrikePrice   *float64              `json:"strike_price,omitempty"`
	StrikeOffset  *int                  `json:"strike_offset,omitempty"`
	TotalLot      int                   `json:"total_lot"`
	LotSize       *int                  `json:"lot_size,omitempty"`
	TotalQuantity *int                  `json:"total_quantity,omitempty"`
	OrderType     model.OrderType       `json:"order_type"`
	TickOffset    *float64              `json:"tick_offset,omitempty"`
	ExpiryType    model.ExpiryType      `json:"expiry_type"`
	ExpiryCycle   model.ExpiryCycle     `json:"expiry_cycle"`
	ExpiryDate    *time.Time            `json:"expiry_date,omitempty"`
	DaysToExpiry  *int                  `json:"days_to_expiry,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

// OptionLegFilter represents filter criteria for searching option legs
type OptionLegFilter struct {
	StrategyID  *uuid.UUID             `json:"strategy_id,omitempty"`
	SelectIndex *string                `json:"select_index,omitempty"`
	Segment     *model.Segment         `json:"segment,omitempty"`
	OptionType  *model.OptionType      `json:"option_type,omitempty"`
	ActionType  *model.TransactionType `json:"action_type,omitempty"`
	ExpiryType  *model.ExpiryType      `json:"expiry_type,omitempty"`
	Page        int                    `json:"page" validate:"min=1"`
	PerPage     int                    `json:"per_page" validate:"min=1,max=100"`
}

// ToModel converts CreateOptionLegRequest to domain model
func (req *CreateOptionLegRequest) ToModel(strategyID uuid.UUID) *model.OptionLeg {
	optionLeg := model.NewOptionLeg(
		strategyID,
		req.LegSequence,
		req.SelectIndex,
		req.Segment,
		req.ActionType,
		req.StrikeType,
		req.ExpiryType,
		req.ExpiryCycle,
	)

	optionLeg.OptionType = req.OptionType
	optionLeg.StrikePrice = req.StrikePrice
	optionLeg.StrikeOffset = req.StrikeOffset
	optionLeg.TotalLot = req.TotalLot
	optionLeg.LotSize = req.LotSize
	optionLeg.OrderType = req.OrderType
	optionLeg.TickOffset = req.TickOffset
	optionLeg.ExpiryDate = req.ExpiryDate
	optionLeg.DaysToExpiry = req.DaysToExpiry

	return optionLeg
}

// ToResponse converts domain model to OptionLegResponse
func ToResponse(optionLeg *model.OptionLeg) *OptionLegResponse {
	return &OptionLegResponse{
		ID:            optionLeg.ID,
		StrategyID:    optionLeg.StrategyID,
		LegSequence:   optionLeg.LegSequence,
		SelectIndex:   optionLeg.SelectIndex,
		Segment:       optionLeg.Segment,
		OptionType:    optionLeg.OptionType,
		ActionType:    optionLeg.ActionType,
		StrikeType:    optionLeg.StrikeType,
		StrikePrice:   optionLeg.StrikePrice,
		StrikeOffset:  optionLeg.StrikeOffset,
		TotalLot:      optionLeg.TotalLot,
		LotSize:       optionLeg.LotSize,
		TotalQuantity: optionLeg.TotalQuantity,
		OrderType:     optionLeg.OrderType,
		TickOffset:    optionLeg.TickOffset,
		ExpiryType:    optionLeg.ExpiryType,
		ExpiryCycle:   optionLeg.ExpiryCycle,
		ExpiryDate:    optionLeg.ExpiryDate,
		DaysToExpiry:  optionLeg.DaysToExpiry,
		CreatedAt:     optionLeg.CreatedAt,
		UpdatedAt:     optionLeg.UpdatedAt,
	}
}
