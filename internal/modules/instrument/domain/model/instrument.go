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
	ID                 uint64    `json:"id" db:"id"`
	ExchangeSecurityID int32     `json:"exchange_security_id" db:"exchange_security_id"`
	Exchange           string    `json:"exchange" db:"exchange"`
	Segment            string    `json:"segment" db:"segment"`
	Symbol             string    `json:"symbol" db:"symbol"`
	Name               string    `json:"name" db:"instrument"` // Named instrument in DB
	StartDate          time.Time `json:"start_date" db:"start_date"`
	Expiry             time.Time `json:"expiry" db:"expiry"`
}

// InstrumentWithMetadata represents an instrument with its metadata
type InstrumentWithMetadata struct {
	Instrument       *Instrument       `json:"instrument"`
	InstrumentMaster *InstrumentMaster `json:"instrument_master"`
}
