// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/infrastructure/persistence/postgres/subscription_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// SubscriptionRepository implements repository.SubscriptionRepository
type SubscriptionRepository struct {
	db *sql.DB
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db *sql.DB) repository.SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

// Helper function to scan subscription row
func scanSubscription(row *sql.Row) (*model.Subscription, error) {
	var subscription model.Subscription
	var endDate, renewalDate sql.NullTime

	err := row.Scan(
		&subscription.ID,
		&subscription.ListingID,
		&subscription.SubscriberID,
		&subscription.Status,
		&subscription.SubscriptionType,
		&subscription.PricePaid,
		&subscription.StartDate,
		&endDate,
		&renewalDate,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if endDate.Valid {
		subscription.EndDate = &endDate.Time
	}

	if renewalDate.Valid {
		subscription.RenewalDate = &renewalDate.Time
	}

	return &subscription, nil
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
	query := `
		INSERT INTO marketplace_subscription (
			id, listing_id, subscriber_id, status, subscription_type,
			price_paid, start_date, end_date, renewal_date, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		subscription.ID,
		subscription.ListingID,
		subscription.SubscriberID,
		subscription.Status,
		subscription.SubscriptionType,
		subscription.PricePaid,
		subscription.StartDate,
		subscription.EndDate,
		subscription.RenewalDate,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

// FindByID returns a subscription by ID
func (r *SubscriptionRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	query := `
		SELECT 
			id, listing_id, subscriber_id, status, subscription_type,
			price_paid, start_date, end_date, renewal_date, created_at, updated_at
		FROM marketplace_subscription
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	subscription, err := scanSubscription(row)
	if err != nil {
		return nil, fmt.Errorf("failed to find subscription: %w", err)
	}

	return subscription, nil
}

// Helper function to scan subscription rows
func scanSubscriptions(rows *sql.Rows) ([]*model.Subscription, error) {
	var subscriptions []*model.Subscription
	for rows.Next() {
		var subscription model.Subscription
		var endDate, renewalDate sql.NullTime

		err := rows.Scan(
			&subscription.ID,
			&subscription.ListingID,
			&subscription.SubscriberID,
			&subscription.Status,
			&subscription.SubscriptionType,
			&subscription.PricePaid,
			&subscription.StartDate,
			&endDate,
			&renewalDate,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			subscription.EndDate = &endDate.Time
		}

		if renewalDate.Valid {
			subscription.RenewalDate = &renewalDate.Time
		}

		subscriptions = append(subscriptions, &subscription)
	}

	return subscriptions, nil
}

// FindByUserID returns all subscriptions for a user
func (r *SubscriptionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error) {
	query := `
		SELECT 
			id, listing_id, subscriber_id, status, subscription_type,
			price_paid, start_date, end_date, renewal_date, created_at, updated_at
		FROM marketplace_subscription
		WHERE subscriber_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions, err := scanSubscriptions(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan subscriptions: %w", err)
	}

	return subscriptions, nil
}

// FindByListingID returns all subscriptions for a listing
func (r *SubscriptionRepository) FindByListingID(ctx context.Context, listingID uuid.UUID) ([]*model.Subscription, error) {
	query := `
		SELECT 
			id, listing_id, subscriber_id, status, subscription_type,
			price_paid, start_date, end_date, renewal_date, created_at, updated_at
		FROM marketplace_subscription
		WHERE listing_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions, err := scanSubscriptions(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan subscriptions: %w", err)
	}

	return subscriptions, nil
}

// Update updates a subscription
func (r *SubscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	query := `
		UPDATE marketplace_subscription SET
			status = $1,
			end_date = $2,
			renewal_date = $3,
			updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		subscription.Status,
		subscription.EndDate,
		subscription.RenewalDate,
		time.Now(),
		subscription.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// IsUserSubscribed checks if a user is subscribed to a listing
func (r *SubscriptionRepository) IsUserSubscribed(ctx context.Context, userID uuid.UUID, listingID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM marketplace_subscription
			WHERE subscriber_id = $1
			AND listing_id = $2
			AND status = 'active'
		)
	`

	var isSubscribed bool
	err := r.db.QueryRowContext(ctx, query, userID, listingID).Scan(&isSubscribed)
	if err != nil {
		return false, fmt.Errorf("failed to check subscription status: %w", err)
	}

	return isSubscribed, nil
}

// GetActiveSubscriptions returns all active subscriptions
func (r *SubscriptionRepository) GetActiveSubscriptions(ctx context.Context) ([]*model.Subscription, error) {
	query := `
		SELECT 
			id, listing_id, subscriber_id, status, subscription_type,
			price_paid, start_date, end_date, renewal_date, created_at, updated_at
		FROM marketplace_subscription
		WHERE status = 'active'
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions, err := scanSubscriptions(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan subscriptions: %w", err)
	}

	return subscriptions, nil
}

// GetExpiredSubscriptions returns all expired subscriptions
func (r *SubscriptionRepository) GetExpiredSubscriptions(ctx context.Context) ([]*model.Subscription, error) {
	query := `
		SELECT 
			id, listing_id, subscriber_id, status, subscription_type,
			price_paid, start_date, end_date, renewal_date, created_at, updated_at
		FROM marketplace_subscription
		WHERE status = 'expired'
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query expired subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions, err := scanSubscriptions(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan subscriptions: %w", err)
	}

	return subscriptions, nil
}
