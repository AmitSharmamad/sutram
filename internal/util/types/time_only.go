package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// TimeOnly represents a time-of-day without date information
// This corresponds to PostgreSQL's TIME type
type TimeOnly struct {
	Hour   int
	Minute int
	Second int
}

// NewTimeOnly creates a new TimeOnly from hour, minute, second
func NewTimeOnly(hour, minute, second int) *TimeOnly {
	return &TimeOnly{
		Hour:   hour,
		Minute: minute,
		Second: second,
	}
}

// NewTimeOnlyFromTime creates a TimeOnly from a time.Time, extracting only the time portion
func NewTimeOnlyFromTime(t time.Time) *TimeOnly {
	return &TimeOnly{
		Hour:   t.Hour(),
		Minute: t.Minute(),
		Second: t.Second(),
	}
}

// ParseTimeOnly parses a time string in format "HH:MM:SS" or "HH:MM" into TimeOnly
func ParseTimeOnly(s string) (*TimeOnly, error) {
	// Try parsing with seconds first (HH:MM:SS)
	if t, err := time.Parse("15:04:05", s); err == nil {
		return NewTimeOnlyFromTime(t), nil
	}

	// Try parsing without seconds (HH:MM)
	if t, err := time.Parse("15:04", s); err == nil {
		return NewTimeOnlyFromTime(t), nil
	}

	return nil, fmt.Errorf("invalid time format: expected HH:MM or HH:MM:SS, got %s", s)
}

// String returns the string representation in "HH:MM:SS" format
func (t TimeOnly) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
}

// ToTime converts TimeOnly to time.Time using a base date of 0000-01-01
func (t TimeOnly) ToTime() time.Time {
	return time.Date(0, 1, 1, t.Hour, t.Minute, t.Second, 0, time.UTC)
}

// ToTimeOnDate converts TimeOnly to time.Time using the provided date
func (t TimeOnly) ToTimeOnDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour, t.Minute, t.Second, 0, date.Location())
}

// Value implements the driver.Valuer interface for database storage
func (t TimeOnly) Value() (driver.Value, error) {
	return t.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		parsed, err := ParseTimeOnly(v)
		if err != nil {
			return err
		}
		*t = *parsed
		return nil
	case []byte:
		parsed, err := ParseTimeOnly(string(v))
		if err != nil {
			return err
		}
		*t = *parsed
		return nil
	case time.Time:
		*t = *NewTimeOnlyFromTime(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into TimeOnly", value)
	}
}

// MarshalJSON implements json.Marshaler interface
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" {
		return nil
	}

	parsed, err := ParseTimeOnly(s)
	if err != nil {
		return err
	}
	*t = *parsed
	return nil
}

// Before checks if this time is before another time
func (t TimeOnly) Before(other TimeOnly) bool {
	return t.ToTime().Before(other.ToTime())
}

// After checks if this time is after another time
func (t TimeOnly) After(other TimeOnly) bool {
	return t.ToTime().After(other.ToTime())
}

// Equal checks if this time equals another time
func (t TimeOnly) Equal(other TimeOnly) bool {
	return t.Hour == other.Hour && t.Minute == other.Minute && t.Second == other.Second
}
