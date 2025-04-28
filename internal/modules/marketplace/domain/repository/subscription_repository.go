// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/domain/repository/subscription_repository.go
package repository

import (
	"context"

	"algo/internal/modules/marketplace/domain/model"

	"github.com/google/uuid"
)

// SubscriptionRepository defines the interface for marketplace subscription repository
type SubscriptionRepository interface {
	// Create creates a new subscription
	Create(ctx context.Context, subscription *model.Subscription) error

	// FindByID returns a subscription by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)

	// FindByUserID returns all subscriptions for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error)

	// FindByListingID returns all subscriptions for a listing
	FindByListingID(ctx context.Context, listingID uuid.UUID) ([]*model.Subscription, error)

	// Update updates a subscription
	Update(ctx context.Context, subscription *model.Subscription) error

	// IsUserSubscribed checks if a user is subscribed to a listing
	IsUserSubscribed(ctx context.Context, userID uuid.UUID, listingID uuid.UUID) (bool, error)

	// GetActiveSubscriptions returns all active subscriptions
	GetActiveSubscriptions(ctx context.Context) ([]*model.Subscription, error)

	// GetExpiredSubscriptions returns all expired subscriptions
	GetExpiredSubscriptions(ctx context.Context) ([]*model.Subscription, error)
}
