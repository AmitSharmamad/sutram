// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/service/seller_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"algo/internal/modules/marketplace/application/dto"
	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// SellerService handles marketplace seller business logic
type SellerService struct {
	listingRepo      repository.ListingRepository
	subscriptionRepo repository.SubscriptionRepository
	transactionRepo  repository.TransactionRepository
}

// NewSellerService creates a new seller service
func NewSellerService(
	listingRepo repository.ListingRepository,
	subscriptionRepo repository.SubscriptionRepository,
	transactionRepo repository.TransactionRepository,
) *SellerService {
	return &SellerService{
		listingRepo:      listingRepo,
		subscriptionRepo: subscriptionRepo,
		transactionRepo:  transactionRepo,
	}
}

// GetSellerListings gets all listings for a seller
func (s *SellerService) GetSellerListings(ctx context.Context, sellerID uuid.UUID) ([]*model.Listing, error) {
	listings, err := s.listingRepo.FindByUserID(ctx, sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller listings: %w", err)
	}

	return listings, nil
}

// GetSellerAnalytics gets analytics data for a seller
func (s *SellerService) GetSellerAnalytics(ctx context.Context, sellerID uuid.UUID) (*dto.SellerAnalyticsResponse, error) {
	// Get all listings for the seller
	listings, err := s.listingRepo.FindByUserID(ctx, sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller listings: %w", err)
	}

	// Calculate metrics
	var totalListings int = len(listings)
	var activeListings int
	var totalViews int
	var sumRatings float32
	var ratedListings int

	listingIDs := make([]uuid.UUID, len(listings))
	for i, listing := range listings {
		listingIDs[i] = listing.ID
		if listing.IsActive {
			activeListings++
		}
		totalViews += listing.Views
		if listing.Ratings > 0 {
			sumRatings += listing.Ratings
			ratedListings++
		}
	}

	// Get the revenue data
	// For a simplified version, get the last 12 months of revenue
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0) // 1 year ago
	endDate := now

	totalRevenue, err := s.transactionRepo.GetRevenueBySeller(ctx, sellerID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue data: %w", err)
	}

	// Get revenue breakdown by month (last 12 months)
	revenueByMonth := make(map[string]float64)

	// In a real implementation, you would query the database for this data
	// For now, we'll return a simplified structure
	currentMonth := now.Format("2006-01")
	revenueByMonth[currentMonth] = totalRevenue * 0.4

	lastMonth := now.AddDate(0, -1, 0).Format("2006-01")
	revenueByMonth[lastMonth] = totalRevenue * 0.3

	twoMonthsAgo := now.AddDate(0, -2, 0).Format("2006-01")
	revenueByMonth[twoMonthsAgo] = totalRevenue * 0.2

	threeMonthsAgo := now.AddDate(0, -3, 0).Format("2006-01")
	revenueByMonth[threeMonthsAgo] = totalRevenue * 0.1

	// Get subscription counts
	var totalSubscriptions int
	var activeSubscriptions int

	// TODO: Implement proper counting of subscriptions
	// This would require additional repository methods

	// For now, just provide some placeholder values
	totalSubscriptions = totalViews / 10
	activeSubscriptions = totalSubscriptions / 2

	// Calculate average rating
	averageRating := float32(0)
	if ratedListings > 0 {
		averageRating = sumRatings / float32(ratedListings)
	}

	return &dto.SellerAnalyticsResponse{
		TotalRevenue:        totalRevenue,
		TotalSubscriptions:  totalSubscriptions,
		ActiveSubscriptions: activeSubscriptions,
		TotalListings:       totalListings,
		ActiveListings:      activeListings,
		AverageRating:       averageRating,
		TotalViews:          totalViews,
		RevenueByMonth:      revenueByMonth,
	}, nil
}
