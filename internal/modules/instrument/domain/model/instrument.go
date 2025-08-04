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
	ExchangeSecurityID int32          `json:"exchange_security_id" db:"exchange_security_id"`
	Exchange           string         `json:"exchange" db:"exchange"`
	Segment            Segment        `json:"segment" db:"segment"`
	Symbol             string         `json:"symbol" db:"symbol"`
	Name               string         `json:"name" db:"instrument"` // Named instrument in DB
	StartDate          time.Time      `json:"start_date" db:"start_date"`
	Expiry             *time.Time     `json:"expiry,omitempty" db:"expiry"`
	ExmExchID          int32          `json:"exm_exch_id" db:"exm_exch_id"`
	SmstSecurityID     int32          `json:"smst_security_id" db:"smst_security_id"`
	InstrumentName     string         `json:"instrument_name" db:"instrument_name"`
	TradingSymbol      string         `json:"trading_symbol" db:"trading_symbol"`
	LotUnits           int32          `json:"lot_units" db:"lot_units"`
	StrikePrice        float64        `json:"strike_price" db:"strike_price"`
	OptionType         string         `json:"option_type" db:"option_type"`
	TickSize           float64        `json:"tick_size" db:"tick_size"`
	ExpiryFlag         string         `json:"expiry_flag" db:"expiry_flag"`
	ExchInstrumentType InstrumentType `json:"exch_instrument_type" db:"exch_instrument_type"`
	SymbolName         string         `json:"symbol_name" db:"symbol_name"`
}

// InstrumentWithMetadata represents an instrument with its metadata
type InstrumentWithMetadata struct {
	Instrument       *Instrument       `json:"instrument"`
	InstrumentMaster *InstrumentMaster `json:"instrument_master"`
}
