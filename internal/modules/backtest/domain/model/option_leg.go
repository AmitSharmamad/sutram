package model

import (
	"time"

	"github.com/google/uuid"
)

// OptionLeg represents an individual leg in a multi-leg options strategy
type OptionLeg struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	StrategyID      uuid.UUID       `json:"strategy_id" db:"strategy_id"`
	LegSequence     int             `json:"leg_sequence" db:"leg_sequence"`
	SelectIndex     string          `json:"select_index" db:"select_index"`
	Segment         Segment         `json:"segment" db:"segment"`
	OptionType      *OptionType     `json:"option_type,omitempty" db:"option_type"`
	ActionType      TransactionType `json:"action_type" db:"action_type"`
	StrikePriceType StrikePriceType `json:"strike_price_type" db:"strike_price_type"`
	StrikePrice     *float64        `json:"strike_price,omitempty" db:"strike_price"`
	StrikeOffset    *int            `json:"strike_offset,omitempty" db:"strike_offset"`
	TotalLot        int             `json:"total_lot" db:"total_lot"`
	LotSize         *int            `json:"lot_size,omitempty" db:"lot_size"`
	TotalQuantity   *int            `json:"total_quantity,omitempty" db:"total_quantity"`
	OrderType       OrderType       `json:"order_type" db:"order_type"`
	TickOffset      *float64        `json:"tick_offset,omitempty" db:"tick_offset"`
	ExpiryType      ExpiryType      `json:"expiry_type" db:"expiry_type"`
	ExpiryDate      *time.Time      `json:"expiry_date,omitempty" db:"expiry_date"`
	DaysToExpiry    *int            `json:"days_to_expiry,omitempty" db:"days_to_expiry"`
	LegStatus       LegStatus       `json:"leg_status" db:"leg_status"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// Segment enum
type Segment string

const (
	SegmentFO Segment = "F&O"
	SegmentEQ Segment = "EQ"
)

// OptionType enum
type OptionType string

const (
	OptionTypeCall OptionType = "CALL"
	OptionTypePut  OptionType = "PUT"
)

// TransactionType enum
type TransactionType string

const (
	TransactionTypeBuy  TransactionType = "BUY"
	TransactionTypeSell TransactionType = "SELL"
)

// StrikePriceType enum
type StrikePriceType string

const (
	StrikePriceTypeATM    StrikePriceType = "ATM"
	StrikePriceTypeITM    StrikePriceType = "ITM"
	StrikePriceTypeOTM    StrikePriceType = "OTM"
	StrikePriceTypeCustom StrikePriceType = "Custom"
)

// OrderType enum
type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

// ExpiryType enum
type ExpiryType string

const (
	ExpiryTypeWeekly    ExpiryType = "Weekly"
	ExpiryTypeMonthly   ExpiryType = "Monthly"
	ExpiryTypeQuarterly ExpiryType = "Quarterly"
)

// LegStatus enum
type LegStatus string

const (
	LegStatusActive    LegStatus = "ACTIVE"
	LegStatusInactive  LegStatus = "INACTIVE"
	LegStatusExecuted  LegStatus = "EXECUTED"
	LegStatusCancelled LegStatus = "CANCELLED"
)

// NewOptionLeg creates a new option leg
func NewOptionLeg(strategyID uuid.UUID, legSequence int, selectIndex string, segment Segment, actionType TransactionType, strikePriceType StrikePriceType, expiryType ExpiryType) *OptionLeg {
	return &OptionLeg{
		ID:              uuid.New(),
		StrategyID:      strategyID,
		LegSequence:     legSequence,
		SelectIndex:     selectIndex,
		Segment:         segment,
		ActionType:      actionType,
		StrikePriceType: strikePriceType,
		TotalLot:        1,
		OrderType:       OrderTypeMarket,
		ExpiryType:      expiryType,
		LegStatus:       LegStatusActive,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// IsValid validates the option leg data
func (ol *OptionLeg) IsValid() error {
	// Validate option details for F&O segment
	if ol.Segment == SegmentFO && ol.OptionType == nil {
		return ErrInvalidOptionDetails
	}

	// Validate strike price for F&O segment
	if ol.Segment == SegmentFO && ol.StrikePrice == nil {
		return ErrInvalidStrikePrice
	}

	// EQ segment should not have option type or strike price
	if ol.Segment == SegmentEQ && (ol.OptionType != nil || ol.StrikePrice != nil) {
		return ErrInvalidEquityDetails
	}

	return nil
}
