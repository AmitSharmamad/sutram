// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/domain/model/listing.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionType defines the type of subscription
type SubscriptionType string

const (
	OneTime SubscriptionType = "one-time"
	Monthly SubscriptionType = "monthly"
	Yearly  SubscriptionType = "yearly"
)

// Listing represents a marketplace listing
type Listing struct {
	ID               uuid.UUID        `json:"id" db:"id"`
	StrategyID       uuid.UUID        `json:"strategy_id" db:"strategy_id"`
	UserID           uuid.UUID        `json:"user_id" db:"user_id"`
	Title            string           `json:"title" db:"title"`
	Description      string           `json:"description" db:"description"`
	Price            float64          `json:"price" db:"price"`
	SubscriptionType SubscriptionType `json:"subscription_type" db:"subscription_type"`
	IsFeatured       bool             `json:"is_featured" db:"is_featured"`
	IsActive         bool             `json:"is_active" db:"is_active"`
	Views            int              `json:"views" db:"views"`
	Ratings          float32          `json:"ratings" db:"ratings"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`

	// Relationships
	Categories  []Category   `json:"categories,omitempty" db:"-"`
	Performance *Performance `json:"performance,omitempty" db:"-"`
	Features    []string     `json:"features,omitempty" db:"-"` // Additional features (stored as JSON)
}

// Performance represents performance metrics for a listing
type Performance struct {
	ListingID          uuid.UUID `json:"listing_id" db:"listing_id"`
	SharpeRatio        float64   `json:"sharpe_ratio" db:"sharpe_ratio"`
	SortinoRatio       float64   `json:"sortino_ratio" db:"sortino_ratio"`
	TotalReturnPct     float64   `json:"total_return_pct" db:"total_return_pct"`
	AnnualizedReturn   float64   `json:"annualized_return" db:"annualized_return"`
	MaxDrawdown        float64   `json:"max_drawdown" db:"max_drawdown"`
	WinRate            float64   `json:"win_rate" db:"win_rate"`
	ProfitFactor       float64   `json:"profit_factor" db:"profit_factor"`
	BacktestPeriodDays int       `json:"backtest_period_days" db:"backtest_period_days"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// Category represents a marketplace category
type Category struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
