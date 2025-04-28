// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/service/listing_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrListingNotFound      = errors.New("listing not found")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidListingData   = errors.New("invalid listing data")
	ErrListingAlreadyExists = errors.New("listing already exists for this strategy")
)

// ListingService handles marketplace listing business logic
type ListingService struct {
	listingRepo     repository.ListingRepository
	categoryRepo    repository.CategoryRepository
	performanceRepo repository.PerformanceRepository
}

// NewListingService creates a new listing service
func NewListingService(
	listingRepo repository.ListingRepository,
	categoryRepo repository.CategoryRepository,
	performanceRepo repository.PerformanceRepository,
) *ListingService {
	return &ListingService{
		listingRepo:     listingRepo,
		categoryRepo:    categoryRepo,
		performanceRepo: performanceRepo,
	}
}

// CreateListing creates a new marketplace listing
func (s *ListingService) CreateListing(
	ctx context.Context,
	userID uuid.UUID,
	strategyID uuid.UUID,
	title,
	description string,
	price float64,
	subscriptionType string,
	categoryIDs []uuid.UUID,
	features []string,
) (*model.Listing, error) {
	// Prepare the listing model
	listing := &model.Listing{
		ID:               uuid.New(),
		StrategyID:       strategyID,
		UserID:           userID,
		Title:            title,
		Description:      description,
		Price:            price,
		SubscriptionType: model.SubscriptionType(subscriptionType),
		IsFeatured:       false,
		IsActive:         true,
		Views:            0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Create the listing
	if err := s.listingRepo.Create(ctx, listing); err != nil {
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	// Add categories
	if err := s.listingRepo.AddCategories(ctx, listing.ID, categoryIDs); err != nil {
		return nil, fmt.Errorf("failed to add categories: %w", err)
	}

	// Store features if provided
	if len(features) > 0 {
		if err := s.listingRepo.SetFeatures(ctx, listing.ID, features); err != nil {
			return nil, fmt.Errorf("failed to set features: %w", err)
		}
	}

	// Load categories
	categories, err := s.listingRepo.GetCategories(ctx, listing.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load categories: %w", err)
	}
	listing.Categories = categories

	// Load features
	listingFeatures, err := s.listingRepo.GetFeatures(ctx, listing.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load features: %w", err)
	}
	listing.Features = listingFeatures

	return listing, nil
}

// GetListingByID retrieves a listing by ID
func (s *ListingService) GetListingByID(ctx context.Context, id uuid.UUID) (*model.Listing, error) {
	// Get the listing
	listing, err := s.listingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find listing: %w", err)
	}
	if listing == nil {
		return nil, ErrListingNotFound
	}

	// Load categories
	categories, err := s.listingRepo.GetCategories(ctx, listing.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load categories: %w", err)
	}
	listing.Categories = categories

	// Load features
	features, err := s.listingRepo.GetFeatures(ctx, listing.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load features: %w", err)
	}
	listing.Features = features

	// Load performance metrics
	performance, err := s.performanceRepo.FindByListingID(ctx, listing.ID)
	if err == nil && performance != nil {
		listing.Performance = performance
	}

	return listing, nil
}

// IncrementViews increments the view count for a listing
func (s *ListingService) IncrementViews(ctx context.Context, id uuid.UUID) error {
	return s.listingRepo.IncrementViews(ctx, id)
}

// GetListings retrieves listings with optional filtering
func (s *ListingService) GetListings(
	ctx context.Context,
	category string,
	sort string,
	search string,
	priceMin string,
	priceMax string,
	page int,
	perPage int,
) ([]*model.Listing, int, error) {
	// Build filter map
	filters := make(map[string]interface{})

	// Add category filter if provided
	if category != "" {
		// Find category by name
		cat, err := s.categoryRepo.FindByName(ctx, category)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to find category: %w", err)
		}
		if cat != nil {
			filters["category_id"] = cat.ID
		}
	}

	// Add search filter
	if search != "" {
		filters["search"] = search
	}

	// Add price range filters
	if priceMin != "" {
		if min, err := strconv.ParseFloat(priceMin, 64); err == nil {
			filters["price_min"] = min
		}
	}

	if priceMax != "" {
		if max, err := strconv.ParseFloat(priceMax, 64); err == nil {
			filters["price_max"] = max
		}
	}

	// Add sort parameter
	if sort != "" {
		filters["sort"] = sort
	}

	// Only show active listings
	filters["is_active"] = true

	// Get listings with filters
	listings, total, err := s.listingRepo.FindAll(ctx, filters, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find listings: %w", err)
	}

	// Load additional data for each listing
	for _, listing := range listings {
		// Load categories
		categories, err := s.listingRepo.GetCategories(ctx, listing.ID)
		if err == nil {
			listing.Categories = categories
		}

		// Load performance
		performance, err := s.performanceRepo.FindByListingID(ctx, listing.ID)
		if err == nil && performance != nil {
			listing.Performance = performance
		}

		// Load features
		features, err := s.listingRepo.GetFeatures(ctx, listing.ID)
		if err == nil {
			listing.Features = features
		}
	}

	return listings, total, nil
}

// UpdateListing updates an existing listing
func (s *ListingService) UpdateListing(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	updates map[string]interface{},
) (*model.Listing, error) {
	// Get the existing listing
	listing, err := s.listingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find listing: %w", err)
	}
	if listing == nil {
		return nil, ErrListingNotFound
	}

	// Check if user owns the listing
	if listing.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Update fields
	if title, ok := updates["title"].(string); ok && title != "" {
		listing.Title = title
	}

	if description, ok := updates["description"].(string); ok {
		listing.Description = description
	}

	if price, ok := updates["price"].(float64); ok && price > 0 {
		listing.Price = price
	}

	if subscriptionType, ok := updates["subscription_type"].(string); ok && subscriptionType != "" {
		listing.SubscriptionType = model.SubscriptionType(subscriptionType)
	}

	if isActivePtr, ok := updates["is_active"].(*bool); ok && isActivePtr != nil {
		listing.IsActive = *isActivePtr
	}

	// Update timestamp
	listing.UpdatedAt = time.Now()

	// Save the updated listing
	if err := s.listingRepo.Update(ctx, listing); err != nil {
		return nil, fmt.Errorf("failed to update listing: %w", err)
	}

	// Update categories if provided
	if categoryIDs, ok := updates["category_ids"].([]uuid.UUID); ok && len(categoryIDs) > 0 {
		// First remove all existing categories
		existingCategories, err := s.listingRepo.GetCategories(ctx, listing.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing categories: %w", err)
		}

		var existingCategoryIDs []uuid.UUID
		for _, category := range existingCategories {
			existingCategoryIDs = append(existingCategoryIDs, category.ID)
		}

		if len(existingCategoryIDs) > 0 {
			if err := s.listingRepo.RemoveCategories(ctx, listing.ID, existingCategoryIDs); err != nil {
				return nil, fmt.Errorf("failed to remove existing categories: %w", err)
			}
		}

		// Then add the new categories
		if err := s.listingRepo.AddCategories(ctx, listing.ID, categoryIDs); err != nil {
			return nil, fmt.Errorf("failed to add new categories: %w", err)
		}
	}

	// Update features if provided
	if features, ok := updates["features"].([]string); ok {
		if err := s.listingRepo.SetFeatures(ctx, listing.ID, features); err != nil {
			return nil, fmt.Errorf("failed to update features: %w", err)
		}
	}

	// Reload the listing with all related data
	return s.GetListingByID(ctx, id)
}

// DeleteListing deletes a listing
func (s *ListingService) DeleteListing(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Get the listing
	listing, err := s.listingRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find listing: %w", err)
	}
	if listing == nil {
		return ErrListingNotFound
	}

	// Check if user owns the listing
	if listing.UserID != userID {
		return ErrUnauthorized
	}

	// Delete the listing
	if err := s.listingRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	return nil
}

// UpdatePerformanceMetrics updates performance metrics for a listing
func (s *ListingService) UpdatePerformanceMetrics(
	ctx context.Context,
	listingID uuid.UUID,
	userID uuid.UUID,
	metrics *model.Performance,
) error {
	// Get the listing
	listing, err := s.listingRepo.FindByID(ctx, listingID)
	if err != nil {
		return fmt.Errorf("failed to find listing: %w", err)
	}
	if listing == nil {
		return ErrListingNotFound
	}

	// Check if user owns the listing
	if listing.UserID != userID {
		return ErrUnauthorized
	}

	// Check if performance exists
	existing, err := s.performanceRepo.FindByListingID(ctx, listingID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("failed to find performance: %w", err)
	}

	metrics.ListingID = listingID
	metrics.UpdatedAt = time.Now()

	if existing == nil {
		// Create new performance entry
		if err := s.performanceRepo.Create(ctx, metrics); err != nil {
			return fmt.Errorf("failed to create performance metrics: %w", err)
		}
	} else {
		// Update existing performance entry
		if err := s.performanceRepo.Update(ctx, metrics); err != nil {
			return fmt.Errorf("failed to update performance metrics: %w", err)
		}
	}

	return nil
}

// GetListingsByUserID gets all listings by user ID
func (s *ListingService) GetListingsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Listing, error) {
	listings, err := s.listingRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find listings by user ID: %w", err)
	}

	// Load additional data for each listing
	for _, listing := range listings {
		// Load categories
		categories, err := s.listingRepo.GetCategories(ctx, listing.ID)
		if err == nil {
			listing.Categories = categories
		}

		// Load performance
		performance, err := s.performanceRepo.FindByListingID(ctx, listing.ID)
		if err == nil && performance != nil {
			listing.Performance = performance
		}

		// Load features
		features, err := s.listingRepo.GetFeatures(ctx, listing.ID)
		if err == nil {
			listing.Features = features
		}
	}

	return listings, nil
}
