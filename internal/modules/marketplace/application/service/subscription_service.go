// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/service/subscription_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// Additional errors
var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrPaymentFailed        = errors.New("payment processing failed")
	ErrAlreadySubscribed    = errors.New("user is already subscribed to this listing")
)

// TransactionStatus defines the transaction status
type TransactionStatus string

const (
	TransactionCompleted TransactionStatus = "completed"
	TransactionPending   TransactionStatus = "pending"
	TransactionFailed    TransactionStatus = "failed"
	TransactionRefunded  TransactionStatus = "refunded"
)

// SubscriptionService handles marketplace subscription business logic
type SubscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository
	listingRepo      repository.ListingRepository
	transactionRepo  repository.TransactionRepository
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(
	subscriptionRepo repository.SubscriptionRepository,
	listingRepo repository.ListingRepository,
	transactionRepo repository.TransactionRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		listingRepo:      listingRepo,
		transactionRepo:  transactionRepo,
	}
}

// CreateSubscription creates a new subscription
func (s *SubscriptionService) CreateSubscription(
	ctx context.Context,
	userID uuid.UUID,
	listingID uuid.UUID,
	paymentMethod string,
	paymentDetails map[string]interface{},
) (*model.Subscription, error) {
	// Check if listing exists and is active
	listing, err := s.listingRepo.FindByID(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find listing: %w", err)
	}
	if listing == nil {
		return nil, ErrListingNotFound
	}
	if !listing.IsActive {
		return nil, errors.New("listing is not active")
	}

	// Check if user is already subscribed
	isSubscribed, err := s.subscriptionRepo.IsUserSubscribed(ctx, userID, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to check subscription status: %w", err)
	}
	if isSubscribed {
		return nil, ErrAlreadySubscribed
	}

	// Create subscription
	now := time.Now()
	subscription := &model.Subscription{
		ID:               uuid.New(),
		ListingID:        listingID,
		SubscriberID:     userID,
		Status:           model.StatusActive,
		SubscriptionType: listing.SubscriptionType,
		PricePaid:        listing.Price,
		StartDate:        now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// Set end date and renewal date based on subscription type
	if listing.SubscriptionType == model.Monthly {
		endDate := now.AddDate(0, 1, 0)
		subscription.EndDate = &endDate
		subscription.RenewalDate = &endDate
	} else if listing.SubscriptionType == model.Yearly {
		endDate := now.AddDate(1, 0, 0)
		subscription.EndDate = &endDate
		subscription.RenewalDate = &endDate
	} else if listing.SubscriptionType == model.OneTime {
		// For one-time subscriptions, no renewal date
		// But we'll set an end date far in the future to indicate "forever" access
		endDate := now.AddDate(100, 0, 0) // 100 years
		subscription.EndDate = &endDate
	}

	// Process payment (simplified, in a real app this would interact with a payment gateway)
	paymentID := ""
	if paymentDetails != nil && paymentMethod != "" {
		// In a real implementation, process payment through payment gateway
		// For now, just simulate with a fake payment ID
		paymentID = uuid.New().String()
	}

	// Save subscription
	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Create transaction record
	transaction := &model.Transaction{
		ID:             uuid.New(),
		SubscriptionID: subscription.ID,
		Amount:         listing.Price,
		Status:         model.TransactionStatusCompleted,
		PaymentMethod:  paymentMethod,
		PaymentID:      paymentID,
		CreatedAt:      now,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Load listing data into the subscription
	subscription.Listing = listing

	return subscription, nil
}

// GetUserSubscriptions gets all subscriptions for a user
func (s *SubscriptionService) GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error) {
	subscriptions, err := s.subscriptionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user subscriptions: %w", err)
	}

	// Load listing data for each subscription
	for _, subscription := range subscriptions {
		listing, err := s.listingRepo.FindByID(ctx, subscription.ListingID)
		if err == nil && listing != nil {
			subscription.Listing = listing
		}
	}

	return subscriptions, nil
}

// CancelSubscription cancels an active subscription
func (s *SubscriptionService) CancelSubscription(
	ctx context.Context,
	subscriptionID uuid.UUID,
	userID uuid.UUID,
) (*model.Subscription, error) {
	// Find subscription
	subscription, err := s.subscriptionRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find subscription: %w", err)
	}
	if subscription == nil {
		return nil, ErrSubscriptionNotFound
	}

	// Check if subscription belongs to user
	if subscription.SubscriberID != userID {
		return nil, ErrUnauthorized
	}

	// Check if subscription is already cancelled
	if subscription.Status == model.StatusCancelled {
		return subscription, nil
	}

	// Update subscription
	subscription.Status = model.StatusCancelled
	subscription.UpdatedAt = time.Now()

	// Save changes
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return subscription, nil
}

// GetSubscriptionByID gets a subscription by ID
func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, subscriptionID uuid.UUID) (*model.Subscription, error) {
	subscription, err := s.subscriptionRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find subscription: %w", err)
	}
	if subscription == nil {
		return nil, ErrSubscriptionNotFound
	}

	// Load listing data
	listing, err := s.listingRepo.FindByID(ctx, subscription.ListingID)
	if err == nil && listing != nil {
		subscription.Listing = listing
	}

	return subscription, nil
}
