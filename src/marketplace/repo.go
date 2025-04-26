package marketplace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"algo/src/shared"
)

// Repository defines the interface for marketplace data operations
type Repository interface {
	// Listing operations
	GetListings(ctx context.Context, params ListingQueryParams) (shared.MarketplaceListings, error)
	GetListingByID(ctx context.Context, id uuid.UUID) (*shared.MarketplaceListing, error)
	CreateListing(ctx context.Context, listing *shared.MarketplaceListing, categoryIDs []string, features []string) (*shared.MarketplaceListing, error)
	UpdateListing(ctx context.Context, id uuid.UUID, listing *shared.UpdateListingRequest) (*shared.MarketplaceListing, error)
	DeleteListing(ctx context.Context, id uuid.UUID) error
	IncrementListingViews(ctx context.Context, id uuid.UUID) error

	// Category operations
	GetCategories(ctx context.Context) (shared.MarketplaceCategories, error)

	// Performance metrics operations
	GetPerformance(ctx context.Context, listingID uuid.UUID) (*shared.Performance, error)
	UpdatePerformance(ctx context.Context, listingID uuid.UUID, performance *shared.UpdatePerformanceRequest) error

	// Subscription operations
	CreateSubscription(ctx context.Context, subscription *shared.Subscription, transaction *shared.Transaction) (*shared.Subscription, error)
	GetUserSubscriptions(ctx context.Context, userID uuid.UUID) (shared.Subscriptions, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, status string, reason string) error

	// Review operations
	GetReviews(ctx context.Context, listingID uuid.UUID) (shared.Reviews, error)
	CreateReview(ctx context.Context, review *shared.Review) (*shared.Review, error)
	UpdateReview(ctx context.Context, id uuid.UUID, userID uuid.UUID, request *shared.UpdateReviewRequest) (*shared.Review, error)
	DeleteReview(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// Transaction operations
	GetTransactions(ctx context.Context, userID uuid.UUID) (shared.Transactions, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*shared.Transaction, error)

	// Seller operations
	GetSellerAnalytics(ctx context.Context, userID uuid.UUID) (*shared.SellerAnalytics, error)
	GetSellerListings(ctx context.Context, userID uuid.UUID) (shared.MarketplaceListings, error)
}

// ListingQueryParams represents the filters for querying marketplace listings
type ListingQueryParams struct {
	Category string
	Sort     string
	Search   string
	Page     int
	PageSize int
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new instance of the marketplace repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

// GetListings retrieves filtered marketplace listings
func (r *repository) GetListings(ctx context.Context, params ListingQueryParams) (shared.MarketplaceListings, error) {
	listings := shared.MarketplaceListings{}

	query := `
		SELECT ml.* FROM marketplace_listing ml
		WHERE ml.is_active = true
	`
	args := []interface{}{}

	// Apply category filter if provided
	if params.Category != "" {
		query += `
			AND ml.id IN (
				SELECT listing_id FROM marketplace_listing_category mlc
				JOIN marketplace_category mc ON mlc.category_id = mc.id
				WHERE mc.name ILIKE $1
			)
		`
		args = append(args, "%"+params.Category+"%")
	}

	// Apply search filter if provided
	if params.Search != "" {
		if len(args) == 0 {
			query += ` AND (ml.title ILIKE $1 OR ml.description ILIKE $1)`
			args = append(args, "%"+params.Search+"%")
		} else {
			query += ` AND (ml.title ILIKE $2 OR ml.description ILIKE $2)`
			args = append(args, "%"+params.Search+"%")
		}
	}

	// Apply sorting
	switch params.Sort {
	case "price_asc":
		query += " ORDER BY ml.price ASC"
	case "price_desc":
		query += " ORDER BY ml.price DESC"
	case "rating":
		query += " ORDER BY ml.ratings DESC"
	case "newest":
		query += " ORDER BY ml.created_at DESC"
	default:
		query += " ORDER BY ml.is_featured DESC, ml.ratings DESC"
	}

	// Apply pagination
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.PageSize, offset)
	}

	err := r.db.SelectContext(ctx, &listings, query, args...)
	if err != nil {
		return nil, err
	}

	// Fetch categories and features for each listing
	for i := range listings {
		// Load categories
		marketplaceCategories := []shared.MarketplaceCategory{}
		err = r.db.SelectContext(ctx, &marketplaceCategories, `
			SELECT mc.* FROM marketplace_category mc
			JOIN marketplace_listing_category mlc ON mc.id = mlc.category_id
			WHERE mlc.listing_id = $1
		`, listings[i].ID)
		if err != nil {
			return nil, err
		}
		listings[i].Categories = marketplaceCategories

		// Get performance metrics
		performance, err := r.GetPerformance(ctx, listings[i].ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if performance != nil {
			listings[i].Performance = performance
		}
	}

	return listings, nil
}

// GetListingByID retrieves a specific listing by its ID
func (r *repository) GetListingByID(ctx context.Context, id uuid.UUID) (*shared.MarketplaceListing, error) {
	listing := shared.MarketplaceListing{}

	err := r.db.GetContext(ctx, &listing, "SELECT * FROM marketplace_listing WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("listing not found")
		}
		return nil, err
	}

	// Load categories
	marketplaceCategories := []shared.MarketplaceCategory{}
	err = r.db.SelectContext(ctx, &marketplaceCategories, `
		SELECT mc.* FROM marketplace_category mc
		JOIN marketplace_listing_category mlc ON mc.id = mlc.category_id
		WHERE mlc.listing_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	listing.Categories = marketplaceCategories

	// Get performance metrics
	performance, err := r.GetPerformance(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if performance != nil {
		listing.Performance = performance
	}

	// Increment view count in a non-blocking way
	go func(repo *repository, listingID uuid.UUID) {
		bgCtx := context.Background()
		_ = repo.IncrementListingViews(bgCtx, listingID)
	}(r, id)

	return &listing, nil
}

// CreateListing creates a new marketplace listing
func (r *repository) CreateListing(ctx context.Context, listing *shared.MarketplaceListing, categoryIDs []string, features []string) (*shared.MarketplaceListing, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Generate new UUID if not provided
	if listing.ID == uuid.Nil {
		listing.ID = uuid.New()
	}

	// Set current time for timestamps
	now := time.Now()
	listing.CreatedAt = now
	listing.UpdatedAt = now

	// Insert the listing
	_, err = tx.ExecContext(ctx, `
		INSERT INTO marketplace_listing (
			id, strategy_id, user_id, title, description, price, 
			subscription_type, is_featured, is_active, views, ratings, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`,
		listing.ID, listing.StrategyID, listing.UserID, listing.Title,
		listing.Description, listing.Price, listing.SubscriptionType,
		listing.IsFeatured, listing.IsActive, 0, 0, listing.CreatedAt, listing.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Associate the listing with categories
	for _, categoryID := range categoryIDs {
		catID, err := uuid.Parse(categoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %s", categoryID)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO marketplace_listing_category (listing_id, category_id)
			VALUES ($1, $2)
		`, listing.ID, catID)
		if err != nil {
			return nil, err
		}
	}

	// Store features as a JSON array in a settings table or similar
	if len(features) > 0 {
		featuresJSON, err := json.Marshal(features)
		if err != nil {
			return nil, err
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO marketplace_settings (key, value, description, updated_at, updated_by)
			VALUES ($1, $2, $3, $4, $5)
		`, "listing_features_"+listing.ID.String(), featuresJSON, "Features for listing "+listing.ID.String(), now, listing.UserID)
		if err != nil {
			return nil, err
		}

		listing.Features = features
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return listing, nil
}

// UpdateListing updates an existing marketplace listing
func (r *repository) UpdateListing(ctx context.Context, id uuid.UUID, updateRequest *shared.UpdateListingRequest) (*shared.MarketplaceListing, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// First, check if the listing exists
	var listing shared.MarketplaceListing
	err = tx.GetContext(ctx, &listing, "SELECT * FROM marketplace_listing WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("listing not found")
		}
		return nil, err
	}

	// Build update query dynamically
	var updates []string
	var args []interface{}
	argCount := 1

	if updateRequest.Title != "" {
		updates = append(updates, fmt.Sprintf("title = $%d", argCount))
		args = append(args, updateRequest.Title)
		listing.Title = updateRequest.Title
		argCount++
	}

	if updateRequest.Description != "" {
		updates = append(updates, fmt.Sprintf("description = $%d", argCount))
		args = append(args, updateRequest.Description)
		listing.Description = updateRequest.Description
		argCount++
	}

	if updateRequest.Price > 0 {
		updates = append(updates, fmt.Sprintf("price = $%d", argCount))
		args = append(args, updateRequest.Price)
		listing.Price = updateRequest.Price
		argCount++
	}

	if updateRequest.SubscriptionType != "" {
		updates = append(updates, fmt.Sprintf("subscription_type = $%d", argCount))
		args = append(args, updateRequest.SubscriptionType)
		listing.SubscriptionType = updateRequest.SubscriptionType
		argCount++
	}

	// Always update the is_active flag and updated_at timestamp
	updates = append(updates, fmt.Sprintf("is_active = $%d", argCount))
	args = append(args, updateRequest.IsActive)
	listing.IsActive = updateRequest.IsActive
	argCount++

	now := time.Now()
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, now)
	listing.UpdatedAt = now
	argCount++

	// Add listing ID to args
	args = append(args, id)

	if len(updates) > 0 {
		query := fmt.Sprintf("UPDATE marketplace_listing SET %s WHERE id = $%d",
			strings.Join(updates, ", "), argCount)

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}
	}

	// Update category associations if provided
	if len(updateRequest.CategoryIDs) > 0 {
		// Remove existing associations
		_, err = tx.ExecContext(ctx, `
			DELETE FROM marketplace_listing_category
			WHERE listing_id = $1
		`, id)
		if err != nil {
			return nil, err
		}

		// Add new associations
		for _, categoryID := range updateRequest.CategoryIDs {
			catID, err := uuid.Parse(categoryID)
			if err != nil {
				return nil, fmt.Errorf("invalid category ID: %s", categoryID)
			}

			_, err = tx.ExecContext(ctx, `
				INSERT INTO marketplace_listing_category (listing_id, category_id)
				VALUES ($1, $2)
			`, id, catID)
			if err != nil {
				return nil, err
			}
		}
	}

	// Update features if provided
	if len(updateRequest.Features) > 0 {
		featuresJSON, err := json.Marshal(updateRequest.Features)
		if err != nil {
			return nil, err
		}

		// Check if features record already exists
		var count int
		err = tx.GetContext(ctx, &count, `
			SELECT COUNT(*) FROM marketplace_settings 
			WHERE key = $1
		`, "listing_features_"+id.String())
		if err != nil {
			return nil, err
		}

		if count > 0 {
			_, err = tx.ExecContext(ctx, `
				UPDATE marketplace_settings
				SET value = $1, updated_at = $2
				WHERE key = $3
			`, featuresJSON, now, "listing_features_"+id.String())
		} else {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO marketplace_settings (key, value, description, updated_at)
				VALUES ($1, $2, $3, $4)
			`, "listing_features_"+id.String(), featuresJSON, "Features for listing "+id.String(), now)
		}

		if err != nil {
			return nil, err
		}

		listing.Features = updateRequest.Features
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Load categories for the updated listing
	marketplaceCategories := []shared.MarketplaceCategory{}
	err = r.db.SelectContext(ctx, &marketplaceCategories, `
		SELECT mc.* FROM marketplace_category mc
		JOIN marketplace_listing_category mlc ON mc.id = mlc.category_id
		WHERE mlc.listing_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	listing.Categories = marketplaceCategories

	return &listing, nil
}

// DeleteListing removes a marketplace listing
func (r *repository) DeleteListing(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM marketplace_listing WHERE id = $1", id)
	return err
}

// IncrementListingViews increments the view count for a listing
func (r *repository) IncrementListingViews(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE marketplace_listing SET views = views + 1 WHERE id = $1", id)
	return err
}

// GetCategories retrieves all marketplace categories
func (r *repository) GetCategories(ctx context.Context) (shared.MarketplaceCategories, error) {
	categories := shared.MarketplaceCategories{}
	err := r.db.SelectContext(ctx, &categories, "SELECT * FROM marketplace_category ORDER BY name")
	return categories, err
}

// GetPerformance retrieves performance metrics for a listing
func (r *repository) GetPerformance(ctx context.Context, listingID uuid.UUID) (*shared.Performance, error) {
	performance := shared.Performance{}

	err := r.db.GetContext(ctx, &performance, "SELECT * FROM marketplace_performance WHERE listing_id = $1", listingID)
	if err != nil {
		return nil, err
	}

	return &performance, nil
}

// UpdatePerformance updates performance metrics for a listing
func (r *repository) UpdatePerformance(ctx context.Context, listingID uuid.UUID, perf *shared.UpdatePerformanceRequest) error {
	// Check if performance metrics already exist
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM marketplace_performance WHERE listing_id = $1", listingID)
	if err != nil {
		return err
	}

	now := time.Now()

	if count > 0 {
		_, err = r.db.ExecContext(ctx, `
			UPDATE marketplace_performance SET
			sharpe_ratio = $1,
			sortino_ratio = $2,
			total_return_pct = $3,
			annualized_return = $4,
			max_drawdown = $5,
			win_rate = $6,
			profit_factor = $7,
			backtest_period_days = $8,
			updated_at = $9
			WHERE listing_id = $10
		`,
			perf.SharpeRatio, perf.SortinoRatio, perf.TotalReturnPct,
			perf.AnnualizedReturn, perf.MaxDrawdown, perf.WinRate,
			perf.ProfitFactor, perf.BacktestPeriodDays, now, listingID)
	} else {
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO marketplace_performance (
				listing_id, sharpe_ratio, sortino_ratio, total_return_pct,
				annualized_return, max_drawdown, win_rate, profit_factor,
				backtest_period_days, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
			)
		`,
			listingID, perf.SharpeRatio, perf.SortinoRatio, perf.TotalReturnPct,
			perf.AnnualizedReturn, perf.MaxDrawdown, perf.WinRate,
			perf.ProfitFactor, perf.BacktestPeriodDays, now)
	}

	return err
}

// CreateSubscription creates a new subscription with an associated transaction
func (r *repository) CreateSubscription(ctx context.Context, subscription *shared.Subscription, transaction *shared.Transaction) (*shared.Subscription, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Generate UUID if not provided
	if subscription.ID == uuid.Nil {
		subscription.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now
	subscription.StartDate = now

	// Calculate end date based on subscription type
	switch subscription.SubscriptionType {
	case "monthly":
		subscription.EndDate = now.AddDate(0, 1, 0)
		subscription.RenewalDate = subscription.EndDate
	case "yearly":
		subscription.EndDate = now.AddDate(1, 0, 0)
		subscription.RenewalDate = subscription.EndDate
	case "one-time":
		// For one-time subscriptions, set a far future end date
		subscription.EndDate = now.AddDate(100, 0, 0)
		subscription.RenewalDate = time.Time{}
	}

	// Insert the subscription
	_, err = tx.ExecContext(ctx, `
		INSERT INTO marketplace_subscription (
			id, listing_id, subscriber_id, status, subscription_type, 
			price_paid, start_date, end_date, renewal_date, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`,
		subscription.ID, subscription.ListingID, subscription.SubscriberID,
		subscription.Status, subscription.SubscriptionType, subscription.PricePaid,
		subscription.StartDate, subscription.EndDate, subscription.RenewalDate,
		subscription.CreatedAt, subscription.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Create associated transaction
	if transaction != nil {
		if transaction.ID == uuid.Nil {
			transaction.ID = uuid.New()
		}

		transaction.SubscriptionID = subscription.ID
		transaction.CreatedAt = now

		_, err = tx.ExecContext(ctx, `
			INSERT INTO marketplace_transaction (
				id, subscription_id, amount, status, payment_method, payment_id, created_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7
			)
		`,
			transaction.ID, transaction.SubscriptionID, transaction.Amount,
			transaction.Status, transaction.PaymentMethod, transaction.PaymentID,
			transaction.CreatedAt)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// GetUserSubscriptions retrieves all subscriptions for a user
func (r *repository) GetUserSubscriptions(ctx context.Context, userID uuid.UUID) (shared.Subscriptions, error) {
	subs := shared.Subscriptions{}

	err := r.db.SelectContext(ctx, &subs, `
		SELECT s.* FROM marketplace_subscription s
		WHERE s.subscriber_id = $1
		ORDER BY s.created_at DESC
	`, userID)

	return subs, err
}

// UpdateSubscription updates a subscription's status
func (r *repository) UpdateSubscription(ctx context.Context, id uuid.UUID, status string, reason string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE marketplace_subscription
		SET status = $1, updated_at = $2
		WHERE id = $3
	`, status, time.Now(), id)

	// Record cancellation reason if provided
	if status == "cancelled" && reason != "" {
		_, _ = r.db.ExecContext(ctx, `
			INSERT INTO marketplace_settings (key, value, updated_at)
			VALUES ($1, $2, $3)
		`, "cancellation_reason_"+id.String(), reason, time.Now())
	}

	return err
}

// GetReviews retrieves all reviews for a listing
func (r *repository) GetReviews(ctx context.Context, listingID uuid.UUID) (shared.Reviews, error) {
	reviews := shared.Reviews{}

	err := r.db.SelectContext(ctx, &reviews, `
		SELECT r.*, u.name as username
		FROM marketplace_review r
		JOIN "user" u ON r.user_id = u.id
		WHERE r.listing_id = $1
		ORDER BY r.created_at DESC
	`, listingID)

	return reviews, err
}

// CreateReview creates a new review
func (r *repository) CreateReview(ctx context.Context, review *shared.Review) (*shared.Review, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Check for an existing review by this user
	var existingCount int
	err = tx.GetContext(ctx, &existingCount, `
		SELECT COUNT(*) FROM marketplace_review
		WHERE listing_id = $1 AND user_id = $2
	`, review.ListingID, review.UserID)
	if err != nil {
		return nil, err
	}

	if existingCount > 0 {
		return nil, fmt.Errorf("user already submitted a review for this listing")
	}

	// Generate UUID if not provided
	if review.ID == uuid.Nil {
		review.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now()
	review.CreatedAt = now
	review.UpdatedAt = now

	// Insert the review
	_, err = tx.ExecContext(ctx, `
		INSERT INTO marketplace_review (
			id, listing_id, user_id, rating, comment, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`, review.ID, review.ListingID, review.UserID, review.Rating, review.Comment, review.CreatedAt, review.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Update the listing's overall rating
	_, err = tx.ExecContext(ctx, `
		UPDATE marketplace_listing
		SET ratings = (
			SELECT AVG(rating)::DECIMAL(3,2)
			FROM marketplace_review
			WHERE listing_id = $1
		)
		WHERE id = $1
	`, review.ListingID)
	if err != nil {
		return nil, err
	}

	// Get the username
	var username string
	err = tx.GetContext(ctx, &username, `
		SELECT name FROM "user" WHERE id = $1
	`, review.UserID)
	if err != nil {
		return nil, err
	}
	review.Username = username

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return review, nil
}

// UpdateReview updates an existing review
func (r *repository) UpdateReview(ctx context.Context, id uuid.UUID, userID uuid.UUID, request *shared.UpdateReviewRequest) (*shared.Review, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Check that the review exists and belongs to this user
	var review shared.Review
	err = tx.GetContext(ctx, &review, `
		SELECT * FROM marketplace_review
		WHERE id = $1 AND user_id = $2
	`, id, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("review not found or access denied")
		}
		return nil, err
	}

	// Update the review
	now := time.Now()
	_, err = tx.ExecContext(ctx, `
		UPDATE marketplace_review
		SET rating = $1, comment = $2, updated_at = $3
		WHERE id = $4
	`, request.Rating, request.Comment, now, id)
	if err != nil {
		return nil, err
	}

	// Update the listing's overall rating
	_, err = tx.ExecContext(ctx, `
		UPDATE marketplace_listing
		SET ratings = (
			SELECT AVG(rating)::DECIMAL(3,2)
			FROM marketplace_review
			WHERE listing_id = $1
		)
		WHERE id = $1
	`, review.ListingID)
	if err != nil {
		return nil, err
	}

	// Get the username
	var username string
	err = tx.GetContext(ctx, &username, `
		SELECT name FROM "user" WHERE id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Update the review object
	review.Rating = request.Rating
	review.Comment = request.Comment
	review.UpdatedAt = now
	review.Username = username

	return &review, nil
}

// DeleteReview removes a review
func (r *repository) DeleteReview(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Verify ownership and get the listing ID
	var listingID uuid.UUID
	err = tx.GetContext(ctx, &listingID, `
		SELECT listing_id FROM marketplace_review
		WHERE id = $1 AND user_id = $2
	`, id, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("review not found or access denied")
		}
		return err
	}

	// Delete the review
	_, err = tx.ExecContext(ctx, `DELETE FROM marketplace_review WHERE id = $1`, id)
	if err != nil {
		return err
	}

	// Update the listing's overall rating
	_, err = tx.ExecContext(ctx, `
		UPDATE marketplace_listing
		SET ratings = (
			SELECT COALESCE(AVG(rating)::DECIMAL(3,2), 0)
			FROM marketplace_review
			WHERE listing_id = $1
		)
		WHERE id = $1
	`, listingID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetTransactions retrieves all transactions for a user
func (r *repository) GetTransactions(ctx context.Context, userID uuid.UUID) (shared.Transactions, error) {
	transactions := shared.Transactions{}

	err := r.db.SelectContext(ctx, &transactions, `
		SELECT t.* FROM marketplace_transaction t
		JOIN marketplace_subscription s ON t.subscription_id = s.id
		WHERE s.subscriber_id = $1
		ORDER BY t.created_at DESC
	`, userID)

	return transactions, err
}

// GetTransactionByID retrieves a specific transaction
func (r *repository) GetTransactionByID(ctx context.Context, id uuid.UUID) (*shared.Transaction, error) {
	transaction := shared.Transaction{}

	err := r.db.GetContext(ctx, &transaction, "SELECT * FROM marketplace_transaction WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}

	return &transaction, nil
}

// GetSellerAnalytics retrieves analytics for a seller
func (r *repository) GetSellerAnalytics(ctx context.Context, userID uuid.UUID) (*shared.SellerAnalytics, error) {
	analytics := shared.SellerAnalytics{}

	// Get active listings count
	err := r.db.GetContext(ctx, &analytics.ActiveListings, `
		SELECT COUNT(*) FROM marketplace_listing
		WHERE user_id = $1 AND is_active = true
	`, userID)
	if err != nil {
		return nil, err
	}

	// Get total sales
	err = r.db.GetContext(ctx, &analytics.TotalSales, `
		SELECT COUNT(t.id) FROM marketplace_transaction t
		JOIN marketplace_subscription s ON t.subscription_id = s.id
		JOIN marketplace_listing l ON s.listing_id = l.id
		WHERE l.user_id = $1 AND t.status = 'completed'
	`, userID)
	if err != nil {
		return nil, err
	}

	// Get total earnings
	err = r.db.GetContext(ctx, &analytics.TotalEarnings, `
		SELECT COALESCE(SUM(t.amount), 0) FROM marketplace_transaction t
		JOIN marketplace_subscription s ON t.subscription_id = s.id
		JOIN marketplace_listing l ON s.listing_id = l.id
		WHERE l.user_id = $1 AND t.status = 'completed'
	`, userID)
	if err != nil {
		return nil, err
	}

	// Get average rating
	err = r.db.GetContext(ctx, &analytics.AverageRating, `
		SELECT COALESCE(AVG(r.rating), 0) FROM marketplace_review r
		JOIN marketplace_listing l ON r.listing_id = l.id
		WHERE l.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	// Get total subscribers
	err = r.db.GetContext(ctx, &analytics.TotalSubscribers, `
		SELECT COUNT(DISTINCT s.subscriber_id) FROM marketplace_subscription s
		JOIN marketplace_listing l ON s.listing_id = l.id
		WHERE l.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	return &analytics, nil
}

// GetSellerListings retrieves all listings created by a seller
func (r *repository) GetSellerListings(ctx context.Context, userID uuid.UUID) (shared.MarketplaceListings, error) {
	listings := shared.MarketplaceListings{}

	err := r.db.SelectContext(ctx, &listings, `
		SELECT * FROM marketplace_listing
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)

	return listings, err
}

// Helper function to convert []shared.Category to []shared.MarketplaceCategory
// The convertToMarketplaceCategories function has been removed as we now directly query
// the marketplace_category table and use the correct types.
