package marketplace

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"algo/src/shared"
)

// Service defines the interface for marketplace business operations
type Service interface {
	// Listing operations
	GetListings(ctx context.Context, params ListingQueryParams) (shared.MarketplaceListings, error)
	GetListingByID(ctx context.Context, id uuid.UUID) (*shared.MarketplaceListing, error)
	CreateListing(ctx context.Context, userID uuid.UUID, req *shared.CreateListingRequest) (*shared.MarketplaceListing, error)
	UpdateListing(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *shared.UpdateListingRequest) (*shared.MarketplaceListing, error)
	DeleteListing(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// Category operations
	GetCategories(ctx context.Context) (shared.MarketplaceCategories, error)

	// Performance metrics operations
	UpdatePerformance(ctx context.Context, listingID uuid.UUID, userID uuid.UUID, req *shared.UpdatePerformanceRequest) error

	// Subscription operations
	CreateSubscription(ctx context.Context, userID uuid.UUID, req *shared.CreateSubscriptionRequest) (*shared.Subscription, error)
	GetUserSubscriptions(ctx context.Context, userID uuid.UUID) (shared.Subscriptions, error)
	CancelSubscription(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *shared.CancelSubscriptionRequest) error

	// Review operations
	GetReviews(ctx context.Context, listingID uuid.UUID) (shared.Reviews, error)
	CreateReview(ctx context.Context, listingID uuid.UUID, userID uuid.UUID, req *shared.CreateReviewRequest) (*shared.Review, error)
	UpdateReview(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *shared.UpdateReviewRequest) (*shared.Review, error)
	DeleteReview(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// Transaction operations
	GetTransactions(ctx context.Context, userID uuid.UUID) (shared.Transactions, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*shared.Transaction, error)

	// Seller operations
	GetSellerAnalytics(ctx context.Context, userID uuid.UUID) (*shared.SellerAnalytics, error)
	GetSellerListings(ctx context.Context, userID uuid.UUID) (shared.MarketplaceListings, error)
}

type service struct {
	repo Repository
}

// NewService creates a new instance of marketplace service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetListings retrieves marketplace listings with filters
func (s *service) GetListings(ctx context.Context, params ListingQueryParams) (shared.MarketplaceListings, error) {
	return s.repo.GetListings(ctx, params)
}

// GetListingByID retrieves a specific marketplace listing
func (s *service) GetListingByID(ctx context.Context, id uuid.UUID) (*shared.MarketplaceListing, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid listing ID")
	}
	return s.repo.GetListingByID(ctx, id)
}

// CreateListing creates a new marketplace listing
func (s *service) CreateListing(ctx context.Context, userID uuid.UUID, req *shared.CreateListingRequest) (*shared.MarketplaceListing, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	if req.StrategyID == uuid.Nil {
		return nil, errors.New("strategy ID is required")
	}

	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	if req.Price < 0 {
		return nil, errors.New("price must be non-negative")
	}

	if req.SubscriptionType != "one-time" && req.SubscriptionType != "monthly" && req.SubscriptionType != "yearly" {
		return nil, errors.New("subscription type must be 'one-time', 'monthly', or 'yearly'")
	}

	// Create a new listing object
	listing := &shared.MarketplaceListing{
		ID:               uuid.New(),
		StrategyID:       req.StrategyID,
		UserID:           userID,
		Title:            req.Title,
		Description:      req.Description,
		Price:            req.Price,
		SubscriptionType: req.SubscriptionType,
		IsFeatured:       false,
		IsActive:         true,
		Views:            0,
		Ratings:          0,
	}

	return s.repo.CreateListing(ctx, listing, req.CategoryIDs, req.Features)
}

// UpdateListing updates an existing marketplace listing
func (s *service) UpdateListing(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *shared.UpdateListingRequest) (*shared.MarketplaceListing, error) {
	if id == uuid.Nil || userID == uuid.Nil {
		return nil, errors.New("invalid ID")
	}

	// Verify ownership
	listing, err := s.repo.GetListingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if listing.UserID != userID {
		return nil, errors.New("access denied: you are not the owner of this listing")
	}

	// Validate subscription type if provided
	if req.SubscriptionType != "" && req.SubscriptionType != "one-time" && req.SubscriptionType != "monthly" && req.SubscriptionType != "yearly" {
		return nil, errors.New("subscription type must be 'one-time', 'monthly', or 'yearly'")
	}

	// Validate price if provided
	if req.Price < 0 {
		return nil, errors.New("price must be non-negative")
	}

	return s.repo.UpdateListing(ctx, id, req)
}

// DeleteListing removes a marketplace listing
func (s *service) DeleteListing(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if id == uuid.Nil || userID == uuid.Nil {
		return errors.New("invalid ID")
	}

	// Verify ownership
	listing, err := s.repo.GetListingByID(ctx, id)
	if err != nil {
		return err
	}

	if listing.UserID != userID {
		return errors.New("access denied: you are not the owner of this listing")
	}

	return s.repo.DeleteListing(ctx, id)
}

// GetCategories retrieves all marketplace categories
func (s *service) GetCategories(ctx context.Context) (shared.MarketplaceCategories, error) {
	return s.repo.GetCategories(ctx)
}

// UpdatePerformance updates performance metrics for a listing
func (s *service) UpdatePerformance(ctx context.Context, listingID uuid.UUID, userID uuid.UUID, req *shared.UpdatePerformanceRequest) error {
	if listingID == uuid.Nil || userID == uuid.Nil {
		return errors.New("invalid ID")
	}

	// Verify ownership
	listing, err := s.repo.GetListingByID(ctx, listingID)
	if err != nil {
		return err
	}

	if listing.UserID != userID {
		return errors.New("access denied: you are not the owner of this listing")
	}

	// Validate performance metrics
	if req.SharpeRatio < -10 || req.SharpeRatio > 10 {
		return errors.New("sharpe ratio must be between -10 and 10")
	}

	if req.SortinoRatio < -10 || req.SortinoRatio > 10 {
		return errors.New("sortino ratio must be between -10 and 10")
	}

	if req.MaxDrawdown < 0 || req.MaxDrawdown > 100 {
		return errors.New("max drawdown must be between 0 and 100")
	}

	if req.WinRate < 0 || req.WinRate > 100 {
		return errors.New("win rate must be between 0 and 100")
	}

	if req.ProfitFactor < 0 {
		return errors.New("profit factor must be non-negative")
	}

	if req.BacktestPeriodDays <= 0 {
		return errors.New("backtest period must be positive")
	}

	return s.repo.UpdatePerformance(ctx, listingID, req)
}

// processPayment is a simplified payment processing function
// In a real application, this would integrate with a payment gateway
func (s *service) processPayment(paymentMethod string, paymentDetails map[string]interface{}, amount float64) (string, error) {
	// Mock payment processing
	// In a real app, this would integrate with a payment provider

	// Simple validation
	if paymentMethod == "credit_card" {
		if _, exists := paymentDetails["card_token"]; !exists {
			return "", errors.New("card_token is required for credit card payments")
		}
	}

	// Generate a mock payment ID
	paymentID := uuid.New().String()

	return paymentID, nil
}

// CreateSubscription creates a new subscription to a marketplace listing
func (s *service) CreateSubscription(ctx context.Context, userID uuid.UUID, req *shared.CreateSubscriptionRequest) (*shared.Subscription, error) {
	if userID == uuid.Nil || req.ListingID == uuid.Nil {
		return nil, errors.New("invalid ID")
	}

	// Get the listing to check availability and price
	listing, err := s.repo.GetListingByID(ctx, req.ListingID)
	if err != nil {
		return nil, err
	}

	if !listing.IsActive {
		return nil, errors.New("this listing is no longer available")
	}

	// Check if the user already has an active subscription to this listing
	subs, err := s.repo.GetUserSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, sub := range subs {
		if sub.ListingID == req.ListingID && sub.Status == "active" {
			return nil, errors.New("you already have an active subscription to this listing")
		}
	}

	// Process payment
	paymentID, err := s.processPayment(req.PaymentMethod, req.PaymentDetails, listing.Price)
	if err != nil {
		return nil, fmt.Errorf("payment processing failed: %w", err)
	}

	// Create subscription
	subscription := &shared.Subscription{
		ID:               uuid.New(),
		ListingID:        req.ListingID,
		SubscriberID:     userID,
		Status:           "active",
		SubscriptionType: listing.SubscriptionType,
		PricePaid:        listing.Price,
		StartDate:        time.Now(),
	}

	// Create transaction record
	transaction := &shared.Transaction{
		ID:            uuid.New(),
		Amount:        listing.Price,
		Status:        "completed",
		PaymentMethod: req.PaymentMethod,
		PaymentID:     paymentID,
	}

	return s.repo.CreateSubscription(ctx, subscription, transaction)
}

// GetUserSubscriptions retrieves all subscriptions for a user
func (s *service) GetUserSubscriptions(ctx context.Context, userID uuid.UUID) (shared.Subscriptions, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetUserSubscriptions(ctx, userID)
}

// CancelSubscription cancels an active subscription
func (s *service) CancelSubscription(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *shared.CancelSubscriptionRequest) error {
	if id == uuid.Nil || userID == uuid.Nil {
		return errors.New("invalid ID")
	}

	if req.Status != "cancelled" {
		return errors.New("invalid status for cancellation")
	}

	// Get user's subscriptions to verify ownership
	subs, err := s.repo.GetUserSubscriptions(ctx, userID)
	if err != nil {
		return err
	}

	var isOwner bool
	for _, sub := range subs {
		if sub.ID == id {
			isOwner = true
			if sub.Status != "active" {
				return errors.New("subscription is not active")
			}
			break
		}
	}

	if !isOwner {
		return errors.New("subscription not found or access denied")
	}

	return s.repo.UpdateSubscription(ctx, id, req.Status, req.CancellationReason)
}

// GetReviews retrieves all reviews for a listing
func (s *service) GetReviews(ctx context.Context, listingID uuid.UUID) (shared.Reviews, error) {
	if listingID == uuid.Nil {
		return nil, errors.New("invalid listing ID")
	}

	return s.repo.GetReviews(ctx, listingID)
}

// CreateReview creates a new review for a listing
func (s *service) CreateReview(ctx context.Context, listingID uuid.UUID, userID uuid.UUID, req *shared.CreateReviewRequest) (*shared.Review, error) {
	if listingID == uuid.Nil || userID == uuid.Nil {
		return nil, errors.New("invalid ID")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	// Verify the listing exists
	_, err := s.repo.GetListingByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	// Check if the user has an active subscription to this listing
	// This ensures that only actual subscribers can leave reviews
	subs, err := s.repo.GetUserSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var hasSubscription bool
	for _, sub := range subs {
		if sub.ListingID == listingID && (sub.Status == "active" || sub.Status == "expired") {
			hasSubscription = true
			break
		}
	}

	if !hasSubscription {
		return nil, errors.New("you must be a subscriber to leave a review")
	}

	review := &shared.Review{
		ID:        uuid.New(),
		ListingID: listingID,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	return s.repo.CreateReview(ctx, review)
}

// UpdateReview updates an existing review
func (s *service) UpdateReview(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *shared.UpdateReviewRequest) (*shared.Review, error) {
	if id == uuid.Nil || userID == uuid.Nil {
		return nil, errors.New("invalid ID")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	return s.repo.UpdateReview(ctx, id, userID, req)
}

// DeleteReview removes a review
func (s *service) DeleteReview(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if id == uuid.Nil || userID == uuid.Nil {
		return errors.New("invalid ID")
	}

	return s.repo.DeleteReview(ctx, id, userID)
}

// GetTransactions retrieves transaction history for a user
func (s *service) GetTransactions(ctx context.Context, userID uuid.UUID) (shared.Transactions, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetTransactions(ctx, userID)
}

// GetTransactionByID retrieves a specific transaction with ownership verification
func (s *service) GetTransactionByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*shared.Transaction, error) {
	if id == uuid.Nil || userID == uuid.Nil {
		return nil, errors.New("invalid ID")
	}

	// Get the transaction
	transaction, err := s.repo.GetTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get the subscription to verify ownership
	subs, err := s.repo.GetUserSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	for _, sub := range subs {
		if sub.ID == transaction.SubscriptionID {
			isOwner = true
			break
		}
	}

	if !isOwner {
		return nil, errors.New("transaction not found or access denied")
	}

	return transaction, nil
}

// GetSellerAnalytics retrieves analytics for a seller
func (s *service) GetSellerAnalytics(ctx context.Context, userID uuid.UUID) (*shared.SellerAnalytics, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetSellerAnalytics(ctx, userID)
}

// GetSellerListings retrieves all listings created by a seller
func (s *service) GetSellerListings(ctx context.Context, userID uuid.UUID) (shared.MarketplaceListings, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetSellerListings(ctx, userID)
}
