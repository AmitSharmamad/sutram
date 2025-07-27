package model

import (
	"encoding/json"
	"fmt"
)

// InstrumentType represents the allowed instrument types
//go:generate stringer -type=InstrumentType

type InstrumentType string

const (
	InstrumentTypeStockCash        InstrumentType = "Stocks (Cash Market)"
	InstrumentTypeStockFutures     InstrumentType = "Stock Futures"
	InstrumentTypeStockOptions     InstrumentType = "Stock Options"
	InstrumentTypeIndexSpot        InstrumentType = "Index Spot"
	InstrumentTypeIndexFutures     InstrumentType = "Index Futures"
	InstrumentTypeIndexOptions     InstrumentType = "Index Options"
	InstrumentTypeCommodityFutures InstrumentType = "Commodity Futures"
	InstrumentTypeCommodityOptions InstrumentType = "Commodity Options"
)

// IsValid checks if the instrument type value is valid
func (it InstrumentType) IsValid() bool {
	switch it {
	case InstrumentTypeStockCash, InstrumentTypeStockFutures, InstrumentTypeStockOptions,
		InstrumentTypeIndexSpot, InstrumentTypeIndexFutures, InstrumentTypeIndexOptions,
		InstrumentTypeCommodityFutures, InstrumentTypeCommodityOptions:
		return true
	default:
		return false
	}
}

// String returns the string representation of the instrument type
func (it InstrumentType) String() string {
	return string(it)
}

// ParseInstrumentType parses a string into an InstrumentType enum
func ParseInstrumentType(s string) (InstrumentType, bool) {
	instrumentType := InstrumentType(s)
	return instrumentType, instrumentType.IsValid()
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (it *InstrumentType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, valid := ParseInstrumentType(str)
	if !valid {
		return fmt.Errorf("invalid instrument type value: %s", str)
	}

	*it = parsed
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (it InstrumentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(it.String())
}
