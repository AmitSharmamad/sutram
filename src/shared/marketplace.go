package shared

import (
	"time"

	"github.com/google/uuid"
)

// MarketplaceListing represents a strategy listing in the marketplace
type MarketplaceListing struct {
	ID               uuid.UUID             `json:"id" db:"id"`
	StrategyID       uuid.UUID             `json:"strategy_id" db:"strategy_id"`
	UserID           uuid.UUID             `json:"user_id" db:"user_id"`
	Title            string                `json:"title" db:"title"`
	Description      string                `json:"description" db:"description"`
	Price            float64               `json:"price" db:"price"`
	SubscriptionType string                `json:"subscription_type" db:"subscription_type"` // 'one-time', 'monthly', 'yearly'
	IsFeatured       bool                  `json:"is_featured" db:"is_featured"`
	IsActive         bool                  `json:"is_active" db:"is_active"`
	Views            int                   `json:"views" db:"views"`
	Ratings          float64               `json:"ratings" db:"ratings"`
	Features         []string              `json:"features"`
	Categories       []MarketplaceCategory `json:"categories,omitempty"`
	Performance      *Performance          `json:"performance,omitempty"`
	CreatedAt        time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at" db:"updated_at"`
}

// MarketplaceListings represents a collection of marketplace listings
type MarketplaceListings []MarketplaceListing

// Performance represents performance metrics for a strategy
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

// MarketplaceCategory represents a marketplace category
type MarketplaceCategory struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Icon        *string   `json:"icon" db:"icon"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Categories represents a collection of categories
type Categories []Category

// Subscription represents a user's subscription to a marketplace listing
type Subscription struct {
	ID               uuid.UUID `json:"id" db:"id"`
	ListingID        uuid.UUID `json:"listing_id" db:"listing_id"`
	SubscriberID     uuid.UUID `json:"subscriber_id" db:"subscriber_id"`
	Status           string    `json:"status" db:"status"`                       // 'active', 'expired', 'cancelled'
	SubscriptionType string    `json:"subscription_type" db:"subscription_type"` // 'one-time', 'monthly', 'yearly'
	PricePaid        float64   `json:"price_paid" db:"price_paid"`
	StartDate        time.Time `json:"start_date" db:"start_date"`
	EndDate          time.Time `json:"end_date" db:"end_date"`
	RenewalDate      time.Time `json:"renewal_date" db:"renewal_date"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Subscriptions represents a collection of subscriptions
type Subscriptions []Subscription

// Review represents a user's review of a marketplace listing
type Review struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ListingID uuid.UUID `json:"listing_id" db:"listing_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Rating    int       `json:"rating" db:"rating"`
	Comment   string    `json:"comment" db:"comment"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Reviews represents a collection of reviews
type Reviews []Review

// Transaction represents a marketplace transaction
type Transaction struct {
	ID             uuid.UUID `json:"id" db:"id"`
	SubscriptionID uuid.UUID `json:"subscription_id" db:"subscription_id"`
	Amount         float64   `json:"amount" db:"amount"`
	Status         string    `json:"status" db:"status"` // 'completed', 'pending', 'failed', 'refunded'
	PaymentMethod  string    `json:"payment_method" db:"payment_method"`
	PaymentID      string    `json:"payment_id" db:"payment_id"` // External payment reference ID
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// Transactions represents a collection of transactions
type Transactions []Transaction

// MarketplaceSetting represents a feature toggle or setting for the marketplace
type MarketplaceSetting struct {
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Description string    `json:"description" db:"description"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy   uuid.UUID `json:"updated_by" db:"updated_by"`
}

// SellerAnalytics represents analytics data for a seller
type SellerAnalytics struct {
	TotalEarnings    float64 `json:"total_earnings"`
	ActiveListings   int     `json:"active_listings"`
	TotalSales       int     `json:"total_sales"`
	AverageRating    float64 `json:"average_rating"`
	TotalSubscribers int     `json:"total_subscribers"`
}

// CreateListingRequest represents the request payload for creating a new marketplace listing
type CreateListingRequest struct {
	StrategyID       uuid.UUID `json:"strategy_id" binding:"required"`
	Title            string    `json:"title" binding:"required"`
	Description      string    `json:"description"`
	Price            float64   `json:"price" binding:"required"`
	SubscriptionType string    `json:"subscription_type" binding:"required"`
	CategoryIDs      []string  `json:"category_ids"`
	Features         []string  `json:"features"`
}

// UpdateListingRequest represents the request payload for updating a marketplace listing
type UpdateListingRequest struct {
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Price            float64  `json:"price"`
	SubscriptionType string   `json:"subscription_type"`
	CategoryIDs      []string `json:"category_ids"`
	Features         []string `json:"features"`
	IsActive         bool     `json:"is_active"`
}

// UpdatePerformanceRequest represents the request payload for updating performance metrics
type UpdatePerformanceRequest struct {
	SharpeRatio        float64 `json:"sharpe_ratio"`
	SortinoRatio       float64 `json:"sortino_ratio"`
	TotalReturnPct     float64 `json:"total_return_pct"`
	AnnualizedReturn   float64 `json:"annualized_return"`
	MaxDrawdown        float64 `json:"max_drawdown"`
	WinRate            float64 `json:"win_rate"`
	ProfitFactor       float64 `json:"profit_factor"`
	BacktestPeriodDays int     `json:"backtest_period_days"`
}

// CreateSubscriptionRequest represents the request payload for creating a subscription
type CreateSubscriptionRequest struct {
	ListingID      uuid.UUID              `json:"listing_id" binding:"required"`
	PaymentMethod  string                 `json:"payment_method" binding:"required"`
	PaymentDetails map[string]interface{} `json:"payment_details"`
}

// CancelSubscriptionRequest represents the request payload for cancelling a subscription
type CancelSubscriptionRequest struct {
	Status             string `json:"status" binding:"required"`
	CancellationReason string `json:"cancellation_reason"`
}

// CreateReviewRequest represents the request payload for creating a review
type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

// UpdateReviewRequest represents the request payload for updating a review
type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}
