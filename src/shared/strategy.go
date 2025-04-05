package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"algo/src/util"

	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ValueType string

var (
	ABSOLUTE   ValueType = "ABSOLUTE"
	PERCENTAGE ValueType = "PERCENTAGE"
)

type Strategy struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`

	Name           string        `json:"name" db:"name"`
	Description    string        `json:"description" db:"description"`
	PositionType   Direction     `json:"position_type" db:"position_type"`
	Instruments    pq.Int32Array `json:"instruments" db:"instruments"`
	ChartType      string        `json:"chart_type" db:"chart_type"`
	ChartTimeframe string        `json:"chart_timeframe" db:"chart_timeframe"`
	Qty            int32         `json:"qty" db:"qty"`

	EntryExpression string `json:"entry_expression" db:"entry_expression"`
	ExitExpression  string `json:"exit_expression" db:"exit_expression"`

	TargetProfit         float64    `json:"target_profit" db:"target_profit"`
	TargetProfitType     ValueType  `json:"target_profit_type" db:"target_profit_type"`
	StopLoss             float64    `json:"stop_loss" db:"stop_loss"`
	StopLossType         ValueType  `json:"stop_loss_type" db:"stop_loss_type"`
	TrailingStopLoss     *float64   `json:"trailing_stop_loss" db:"trailing_stop_loss"`
	TrailingStopLossType *ValueType `json:"trailing_stop_loss_type" db:"trailing_stop_loss_type"`
	ExitType             string     `json:"exit_type" db:"exit_type"`

	Capital        float64    `json:"capital" db:"capital"`
	MaxTrades      int32      `json:"max_trades" db:"max_trades"`
	EntryStartTime CustomTime `json:"entry_start_time" db:"entry_start_time"`
	EntryEndTime   CustomTime `json:"entry_end_time" db:"entry_end_time"`
	StartPeriod    CustomTime `json:"start_date" db:"start_date"`
	EndPeriod      CustomTime `json:"end_date" db:"end_date"`
	IsActive       bool       `json:"is_active" db:"is_active"`

	// Audit
	Audit
}

type Strategies []Strategy

func (strategy *Strategy) Validate() {
	path := "../data/llm.json"

	exp := util.GetFileContent(path)

	util.Validate(exp)
}

type StrategyService interface {
	Get(ctx context.Context, id uuid.UUID) *Strategy
	Create(ctx context.Context, strategy *Strategy) (*Strategy, error)
	GetData() StockData
	GetAll(ctx context.Context) (*Strategies, error)
	Validate(strategy *Strategy) bool
}

// CustomTime is a wrapper around time.Time that implements json.Unmarshaler and database/sql interfaces
type CustomTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Handle null/nil case
	if string(b) == "null" || len(b) == 0 {
		ct.Time = time.Time{}
		return nil
	}

	// Try parsing as string first
	var timeStr string
	if err := json.Unmarshal(b, &timeStr); err == nil {
		// Handle empty string case
		if timeStr == "" {
			ct.Time = time.Time{}
			return nil
		}

		// Try RFC3339 format first
		if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
			ct.Time = t
			return nil
		}

		// Try other formats
		formats := []string{
			time.RFC3339Nano,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, timeStr); err == nil {
				ct.Time = t
				return nil
			}
		}

		return fmt.Errorf("time string '%s' doesn't match any supported format", timeStr)
	}

	// Try parsing as Unix timestamp
	var timestamp int64
	if err := json.Unmarshal(b, &timestamp); err == nil {
		// Handle zero timestamp
		if timestamp == 0 {
			ct.Time = time.Time{}
			return nil
		}
		ct.Time = time.Unix(timestamp, 0)
		return nil
	}

	return fmt.Errorf("time must be a string in RFC3339 format (e.g. '2025-01-26T15:45:23Z') or a Unix timestamp. Got: %s", string(b))
}

// MarshalJSON implements json.Marshaler
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Time.Format(time.RFC3339))
}

// Value implements the driver.Valuer interface for database/sql
func (ct CustomTime) Value() (driver.Value, error) {
	if ct.Time.IsZero() {
		return nil, nil
	}
	return ct.Time.Format("2006-01-02 15:04:05"), nil
}

// Scan implements the sql.Scanner interface for database/sql
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		ct.Time = v
	case []byte:
		t, err := time.Parse(time.RFC3339, string(v))
		if err != nil {
			// Try parsing PostgreSQL timestamp format if RFC3339 fails
			t, err = time.Parse("2006-01-02 15:04:05.999999-07", string(v))
			if err != nil {
				t, err = time.Parse("2006-01-02 15:04:05", string(v))
				if err != nil {
					return err
				}
			}
		}
		ct.Time = t
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			// Try parsing PostgreSQL timestamp format if RFC3339 fails
			t, err = time.Parse("2006-01-02 15:04:05.999999-07", v)
			if err != nil {
				t, err = time.Parse("2006-01-02 15:04:05", v)
				if err != nil {
					return err
				}
			}
		}
		ct.Time = t
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type CustomTime", value)
	}
	return nil
}
