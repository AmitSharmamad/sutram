// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/dto/subscription_dto.go
package dto

import (
	"time"

	"algo/internal/modules/marketplace/domain/model"

	"github.com/google/uuid"
)

// CreateSubscriptionRequest represents a request to create a subscription
type CreateSubscriptionRequest struct {
	ListingID      uuid.UUID              `json:"listing_id" binding:"required,uuid4"`
	PaymentMethod  string                 `json:"payment_method" binding:"required"`
	PaymentDetails map[string]interface{} `json:"payment_details" binding:"required"`
}

// UpdateSubscriptionRequest represents a request to update a subscription
type UpdateSubscriptionRequest struct {
	Status             string `json:"status" binding:"required,oneof=cancelled"`
	CancellationReason string `json:"cancellation_reason" binding:"omitempty"`
}

// SubscriptionResponse represents a subscription response
type SubscriptionResponse struct {
	ID               uuid.UUID       `json:"id"`
	ListingID        uuid.UUID       `json:"listing_id"`
	Listing          ListingResponse `json:"listing,omitempty"`
	Status           string          `json:"status"`
	SubscriptionType string          `json:"subscription_type"`
	PricePaid        float64         `json:"price_paid"`
	StartDate        time.Time       `json:"start_date"`
	EndDate          *time.Time      `json:"end_date,omitempty"`
	RenewalDate      *time.Time      `json:"renewal_date,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// ToSubscriptionResponse converts a model.Subscription to a SubscriptionResponse
func ToSubscriptionResponse(subscription *model.Subscription) SubscriptionResponse {
	var listingResponse ListingResponse
	if subscription.Listing != nil {
		listingResponse = ToListingResponse(subscription.Listing)
	}

	return SubscriptionResponse{
		ID:               subscription.ID,
		ListingID:        subscription.ListingID,
		Listing:          listingResponse,
		Status:           string(subscription.Status),
		SubscriptionType: string(subscription.SubscriptionType),
		PricePaid:        subscription.PricePaid,
		StartDate:        subscription.StartDate,
		EndDate:          subscription.EndDate,
		RenewalDate:      subscription.RenewalDate,
		CreatedAt:        subscription.CreatedAt,
		UpdatedAt:        subscription.UpdatedAt,
	}
}
