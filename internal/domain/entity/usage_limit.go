package entity

import (
	"time"

	"github.com/google/uuid"
)

// UsageLimitType represents the type of usage limit
type UsageLimitType string

// Usage limit type constants
const (
	BacktestPerMonth UsageLimitType = "backtest_per_month"
	DeploymentCount  UsageLimitType = "deployment_count"
)

// ResetPeriodType represents the reset period for usage limits
type ResetPeriodType string

// Reset period type constants
const (
	Hourly  ResetPeriodType = "hourly"
	Daily   ResetPeriodType = "daily"
	Weekly  ResetPeriodType = "weekly"
	Monthly ResetPeriodType = "monthly"
	Yearly  ResetPeriodType = "yearly"
	Never   ResetPeriodType = "never"
)

// UsageLimit represents a usage limit for a user
type UsageLimit struct {
	ID           int             `json:"id" db:"id"`
	UserID       uuid.UUID       `json:"user_id" db:"user_id"`
	LimitType    UsageLimitType  `json:"limit_type" db:"limit_type"`
	LimitValue   int             `json:"limit_value" db:"limit_value"`
	CurrentUsage int             `json:"current_usage" db:"current_usage"`
	ResetPeriod  ResetPeriodType `json:"reset_period" db:"reset_period"`
	LastResetAt  time.Time       `json:"last_reset_at" db:"last_reset_at"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

// UsageLog represents a log entry for usage tracking
type UsageLog struct {
	ID            int64                  `json:"id" db:"id"`
	UserID        uuid.UUID              `json:"user_id" db:"user_id"`
	LimitType     UsageLimitType         `json:"limit_type" db:"limit_type"`
	ActionType    string                 `json:"action_type" db:"action_type"`
	PreviousValue int                    `json:"previous_value" db:"previous_value"`
	NewValue      int                    `json:"new_value" db:"new_value"`
	Metadata      map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

// CanPerformAction checks if the user can perform an action based on current usage
func (ul *UsageLimit) CanPerformAction() bool {
	return ul.CurrentUsage < ul.LimitValue
}

// GetRemainingUsage returns the remaining usage for this limit
func (ul *UsageLimit) GetRemainingUsage() int {
	remaining := ul.LimitValue - ul.CurrentUsage
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ShouldReset checks if the usage should be reset based on the reset period
func (ul *UsageLimit) ShouldReset(now time.Time) bool {
	switch ul.ResetPeriod {
	case Never:
		return false
	case Hourly:
		return now.Sub(ul.LastResetAt) >= time.Hour
	case Daily:
		return now.Truncate(24 * time.Hour).After(ul.LastResetAt.Truncate(24 * time.Hour))
	case Weekly:
		// Reset on Monday
		weekStart := now.Truncate(24 * time.Hour)
		for weekStart.Weekday() != time.Monday {
			weekStart = weekStart.AddDate(0, 0, -1)
		}
		lastResetWeekStart := ul.LastResetAt.Truncate(24 * time.Hour)
		for lastResetWeekStart.Weekday() != time.Monday {
			lastResetWeekStart = lastResetWeekStart.AddDate(0, 0, -1)
		}
		return weekStart.After(lastResetWeekStart)
	case Monthly:
		return now.Year() > ul.LastResetAt.Year() ||
			(now.Year() == ul.LastResetAt.Year() && now.Month() > ul.LastResetAt.Month())
	case Yearly:
		return now.Year() > ul.LastResetAt.Year()
	default:
		return false
	}
}
