package model

import (
	"time"
)

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
