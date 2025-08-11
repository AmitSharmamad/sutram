package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// Market represents the trading market types
type Market string

const (
	MarketEQ Market = "EQ"  // Equity market
	MarketFO Market = "F&O" // Futures & Options market
)

// IsValid checks if the market value is valid
func (m Market) IsValid() bool {
	switch m {
	case MarketEQ, MarketFO:
		return true
	default:
		return false
	}
}

// String returns the string representation of the market
func (m Market) String() string {
	return string(m)
}

// ParseMarket parses a string into a Market enum
func ParseMarket(s string) (Market, bool) {
	market := Market(s)
	return market, market.IsValid()
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (m *Market) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, valid := ParseMarket(str)
	if !valid {
		return fmt.Errorf("invalid market value: %s", str)
	}

	*m = parsed
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (m Market) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

// InstrumentCategory represents the basic instrument categories for filtering
type InstrumentCategory string

const (
	InstrumentCategoryStock InstrumentCategory = "Stock"
	InstrumentCategoryIndex InstrumentCategory = "Index"
)

// IsValid checks if the instrument category value is valid
func (ic InstrumentCategory) IsValid() bool {
	switch ic {
	case InstrumentCategoryStock, InstrumentCategoryIndex:
		return true
	default:
		return false
	}
}

// String returns the string representation of the instrument category
func (ic InstrumentCategory) String() string {
	return string(ic)
}

// ParseInstrumentCategory parses a string into an InstrumentCategory enum
func ParseInstrumentCategory(s string) (InstrumentCategory, bool) {
	category := InstrumentCategory(s)
	return category, category.IsValid()
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ic *InstrumentCategory) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, valid := ParseInstrumentCategory(str)
	if !valid {
		return fmt.Errorf("invalid instrument category value: %s", str)
	}

	*ic = parsed
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (ic InstrumentCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(ic.String())
}

// Contract represents the contract types for filtering
type Contract string

const (
	ContractSpot    Contract = "Spot"
	ContractFutures Contract = "Futures"
	ContractOptions Contract = "Options"
)

// IsValid checks if the contract value is valid
func (c Contract) IsValid() bool {
	switch c {
	case ContractSpot, ContractFutures, ContractOptions:
		return true
	default:
		return false
	}
}

// String returns the string representation of the contract
func (c Contract) String() string {
	return string(c)
}

// ParseContract parses a string into a Contract enum
func ParseContract(s string) (Contract, bool) {
	contract := Contract(s)
	return contract, contract.IsValid()
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (c *Contract) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, valid := ParseContract(str)
	if !valid {
		return fmt.Errorf("invalid contract value: %s", str)
	}

	*c = parsed
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (c Contract) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// Instrument represents a financial instrument with time series data
type Instrument struct {
	ExchangeSecurityID int32     `json:"exchange_security_id" db:"exchange_security_id"`
	Open               float64   `json:"open" db:"open"`
	High               float64   `json:"high" db:"high"`
	Low                float64   `json:"low" db:"low"`
	Close              float64   `json:"close" db:"close"`
	Volume             float64   `json:"volume" db:"volume"`
	Timestamp          time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type StockData []Instrument

// InstrumentMaster represents the master instrument record with metadata
type InstrumentMaster struct {
	// Core Identification Columns
	ExchangeSecurityID int32  `json:"exchange_security_id" db:"exchange_security_id"`
	Exchange           string `json:"exchange" db:"exchange"`
	Segment            string `json:"segment" db:"segment"`
	Symbol             string `json:"symbol" db:"symbol"`
	Instrument         string `json:"instrument" db:"instrument"`

	// Trading Symbols and Names
	InstrumentName string `json:"instrument_name" db:"instrument_name"`

	// ISIN and Security Identifiers
	ISIN                 string  `json:"isin" db:"isin"`
	UnderlyingSecurityID *int32  `json:"underlying_security_id,omitempty" db:"underlying_security_id"`
	UnderlyingSymbol     *string `json:"underlying_symbol,omitempty" db:"underlying_symbol"`

	// Date and Expiry Information
	Expiry     time.Time `json:"expiry" db:"expiry"`
	ExpiryFlag *string   `json:"expiry_flag,omitempty" db:"expiry_flag"`

	// Trading Parameters
	LotUnits       int32    `json:"lot_units" db:"lot_units"`
	TickSize       float64  `json:"tick_size" db:"tick_size"`
	StrikePrice    *float64 `json:"strike_price,omitempty" db:"strike_price"`
	OptionType     *string  `json:"option_type,omitempty" db:"option_type"`
	InstrumentType string   `json:"instrument_type" db:"instrument_type"`

	// Order Type Availability Flags
	BracketFlag string `json:"bracket_flag" db:"bracket_flag"`
	CoverFlag   string `json:"cover_flag" db:"cover_flag"`

	// Surveillance and Trading Permissions
	AsmGsmFlag       string `json:"asm_gsm_flag" db:"asm_gsm_flag"`
	BuySellIndicator string `json:"buy_sell_indicator" db:"buy_sell_indicator"`

	// Cover Order Margin Parameters
	BuyCoMinMarginPer    *float64 `json:"buy_co_min_margin_per,omitempty" db:"buy_co_min_margin_per"`
	SellCoMinMarginPer   *float64 `json:"sell_co_min_margin_per,omitempty" db:"sell_co_min_margin_per"`
	BuyCoSlRangeMaxPerc  *float64 `json:"buy_co_sl_range_max_perc,omitempty" db:"buy_co_sl_range_max_perc"`
	SellCoSlRangeMinPerc *float64 `json:"sell_co_sl_range_min_perc,omitempty" db:"sell_co_sl_range_min_perc"`

	// Bracket Order Parameters
	BuyBoMinMarginPer        *float64 `json:"buy_bo_min_margin_per,omitempty" db:"buy_bo_min_margin_per"`
	BuyBoSlRangeMaxPerc      *float64 `json:"buy_bo_sl_range_max_perc,omitempty" db:"buy_bo_sl_range_max_perc"`
	SellBoSlMinRange         *float64 `json:"sell_bo_sl_min_range,omitempty" db:"sell_bo_sl_min_range"`
	BuyBoProfitRangeMaxPerc  *float64 `json:"buy_bo_profit_range_max_perc,omitempty" db:"buy_bo_profit_range_max_perc"`
	SellBoProfitRangeMaxPerc *float64 `json:"sell_bo_profit_range_max_perc,omitempty" db:"sell_bo_profit_range_max_perc"`

	// Margin Trading Facility
	MtfLeverage *float64 `json:"mtf_leverage,omitempty" db:"mtf_leverage"`

	// Metadata
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// InstrumentWithMetadata represents an instrument with its metadata
type InstrumentWithMetadata struct {
	Instrument       *Instrument       `json:"instrument"`
	InstrumentMaster *InstrumentMaster `json:"instrument_master"`
}
