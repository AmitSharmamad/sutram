package model

import (
	"encoding/json"
	"fmt"
)

// Segment represents the allowed trading segments
//go:generate stringer -type=Segment

type Segment string

const (
	SegmentEquities    Segment = "Equities"
	SegmentDerivatives Segment = "Derivatives"
	SegmentCommodities Segment = "Commodities"
	SegmentCurrencies  Segment = "Currencies"
	SegmentIndices     Segment = "Indices"
)

// IsValid checks if the segment value is valid
func (s Segment) IsValid() bool {
	switch s {
	case SegmentEquities, SegmentDerivatives, SegmentCommodities, SegmentCurrencies, SegmentIndices:
		return true
	default:
		return false
	}
}

// String returns the string representation of the segment
func (s Segment) String() string {
	return string(s)
}

// ParseSegment parses a string into a Segment enum
func ParseSegment(s string) (Segment, bool) {
	segment := Segment(s)
	return segment, segment.IsValid()
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (s *Segment) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, valid := ParseSegment(str)
	if !valid {
		return fmt.Errorf("invalid segment value: %s", str)
	}

	*s = parsed
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (s Segment) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
