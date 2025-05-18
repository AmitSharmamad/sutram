// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/domain/model/subscription.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionStatus defines the status of a subscription
type SubscriptionStatus string

const (
	StatusActive    SubscriptionStatus = "active"
	StatusExpired   SubscriptionStatus = "expired"
	StatusCancelled SubscriptionStatus = "cancelled"
)

// Subscription represents a user subscription to a marketplace listing
type Subscription struct {
	ID               uuid.UUID          `json:"id" db:"id"`
	ListingID        uuid.UUID          `json:"listing_id" db:"listing_id"`
	SubscriberID     uuid.UUID          `json:"subscriber_id" db:"subscriber_id"`
	Status           SubscriptionStatus `json:"status" db:"status"`
	SubscriptionType SubscriptionType   `json:"subscription_type" db:"subscription_type"`
	PricePaid        float64            `json:"price_paid" db:"price_paid"`
	StartDate        time.Time          `json:"start_date" db:"start_date"`
	EndDate          *time.Time         `json:"end_date,omitempty" db:"end_date"`
	RenewalDate      *time.Time         `json:"renewal_date,omitempty" db:"renewal_date"`
	CreatedAt        time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" db:"updated_at"`

	// Relationships
	Listing      *Listing       `json:"listing,omitempty" db:"-"`
	Transactions []*Transaction `json:"transactions,omitempty" db:"-"`
}
