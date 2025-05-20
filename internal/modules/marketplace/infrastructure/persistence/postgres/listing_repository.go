// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/infrastructure/persistence/postgres/listing_repository.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// ListingRepository implements the repository.ListingRepository interface
type ListingRepository struct {
	db *sql.DB
}

// NewListingRepository creates a new listing repository
func NewListingRepository(db *sql.DB) repository.ListingRepository {
	return &ListingRepository{
		db: db,
	}
}

// Create creates a new listing
func (r *ListingRepository) Create(ctx context.Context, listing *model.Listing) error {
	query := `
		INSERT INTO marketplace_listing (
			id, strategy_id, user_id, title, description, price,
			subscription_type, is_featured, is_active, views, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		listing.ID,
		listing.StrategyID,
		listing.UserID,
		listing.Title,
		listing.Description,
		listing.Price,
		listing.SubscriptionType,
		listing.IsFeatured,
		listing.IsActive,
		listing.Views,
		listing.CreatedAt,
		listing.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create listing: %w", err)
	}

	return nil
}

// FindByID returns a listing by ID
func (r *ListingRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Listing, error) {
	query := `
		SELECT 
			l.id, l.strategy_id, l.user_id, l.title, l.description, l.price,
			l.subscription_type, l.is_featured, l.is_active, l.views, l.ratings, l.created_at, l.updated_at,
			u.name as seller_name
		FROM marketplace_listing l
		LEFT JOIN "user" u ON l.user_id = u.id
		WHERE l.id = $1
	`

	var listing model.Listing
	var sellerName sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		sqlx.Rebind(sqlx.DOLLAR, query),
		id,
	).Scan(
		&listing.ID,
		&listing.StrategyID,
		&listing.UserID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.SubscriptionType,
		&listing.IsFeatured,
		&listing.IsActive,
		&listing.Views,
		&listing.Ratings,
		&listing.CreatedAt,
		&listing.UpdatedAt,
		&sellerName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find listing: %w", err)
	}

	return &listing, nil
}

// FindAll returns all listings with optional filtering and pagination
func (r *ListingRepository) FindAll(ctx context.Context, filters map[string]interface{}, page, perPage int) ([]*model.Listing, int, error) {
	// Build the query with filters
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT 
			l.id, l.strategy_id, l.user_id, l.title, l.description, l.price,
			l.subscription_type, l.is_featured, l.is_active, l.views, l.ratings, l.created_at, l.updated_at,
			u.name as seller_name
		FROM marketplace_listing l
		LEFT JOIN "user" u ON l.user_id = u.id
	`)

	// Base WHERE clause
	whereClause := []string{}
	args := []interface{}{}
	argIndex := 1

	// Add filters
	if filters != nil {
		// Filter by category
		if categoryID, ok := filters["category_id"].(uuid.UUID); ok {
			whereClause = append(whereClause, fmt.Sprintf("EXISTS (SELECT 1 FROM marketplace_listing_category mlc WHERE mlc.listing_id = l.id AND mlc.category_id = $%d)", argIndex))
			args = append(args, categoryID)
			argIndex++
		}

		// Filter by active status
		if isActive, ok := filters["is_active"].(bool); ok {
			whereClause = append(whereClause, fmt.Sprintf("l.is_active = $%d", argIndex))
			args = append(args, isActive)
			argIndex++
		}

		// Filter by search term (search in title and description)
		if search, ok := filters["search"].(string); ok && search != "" {
			whereClause = append(whereClause, fmt.Sprintf("(l.title ILIKE $%d OR l.description ILIKE $%d)", argIndex, argIndex))
			args = append(args, "%"+search+"%")
			argIndex++
		}

		// Filter by price range
		if priceMin, ok := filters["price_min"].(float64); ok {
			whereClause = append(whereClause, fmt.Sprintf("l.price >= $%d", argIndex))
			args = append(args, priceMin)
			argIndex++
		}

		if priceMax, ok := filters["price_max"].(float64); ok {
			whereClause = append(whereClause, fmt.Sprintf("l.price <= $%d", argIndex))
			args = append(args, priceMax)
			argIndex++
		}
	}

	// Add WHERE clause to query if any filters were applied
	if len(whereClause) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(whereClause, " AND "))
	}

	// Add ORDER BY clause
	orderBy := "l.created_at DESC" // default sorting
	if filters != nil {
		if sort, ok := filters["sort"].(string); ok {
			switch sort {
			case "price":
				orderBy = "l.price ASC"
			case "rating":
				orderBy = "l.ratings DESC"
			case "popularity":
				orderBy = "l.views DESC"
			}
		}
	}
	queryBuilder.WriteString(" ORDER BY " + orderBy)

	// Add pagination
	offset := (page - 1) * perPage
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1))
	args = append(args, perPage, offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query listings: %w", err)
	}
	defer rows.Close()

	// Parse results
	var listings []*model.Listing
	for rows.Next() {
		var listing model.Listing
		var sellerName sql.NullString

		err := rows.Scan(
			&listing.ID,
			&listing.StrategyID,
			&listing.UserID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.SubscriptionType,
			&listing.IsFeatured,
			&listing.IsActive,
			&listing.Views,
			&listing.Ratings,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&sellerName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan listing row: %w", err)
		}

		listings = append(listings, &listing)
	}

	// Count total results for pagination
	countQuery := strings.Builder{}
	countQuery.WriteString("SELECT COUNT(*) FROM marketplace_listing l")
	if len(whereClause) > 0 {
		countQuery.WriteString(" WHERE ")
		countQuery.WriteString(strings.Join(whereClause, " AND "))
	}

	var total int
	err = r.db.QueryRowContext(ctx, countQuery.String(), args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count total listings: %w", err)
	}

	return listings, total, nil
}

// Update updates a listing
func (r *ListingRepository) Update(ctx context.Context, listing *model.Listing) error {
	query := `
		UPDATE marketplace_listing SET
			title = $1,
			description = $2,
			price = $3,
			subscription_type = $4,
			is_featured = $5,
			is_active = $6,
			views = $7,
			ratings = $8,
			updated_at = $9
		WHERE id = $10
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		listing.Title,
		listing.Description,
		listing.Price,
		listing.SubscriptionType,
		listing.IsFeatured,
		listing.IsActive,
		listing.Views,
		listing.Ratings,
		time.Now(),
		listing.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update listing: %w", err)
	}

	return nil
}

// Delete deletes a listing
func (r *ListingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM marketplace_listing WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	return nil
}

// IncrementViews increments the view count for a listing
func (r *ListingRepository) IncrementViews(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE marketplace_listing SET views = views + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}

	return nil
}

// FindByUserID returns all listings created by a user
func (r *ListingRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Listing, error) {
	query := `
		SELECT 
			id, strategy_id, user_id, title, description, price,
			subscription_type, is_featured, is_active, views, ratings, created_at, updated_at
		FROM marketplace_listing
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query listings by user ID: %w", err)
	}
	defer rows.Close()

	var listings []*model.Listing
	for rows.Next() {
		var listing model.Listing
		err := rows.Scan(
			&listing.ID,
			&listing.StrategyID,
			&listing.UserID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.SubscriptionType,
			&listing.IsFeatured,
			&listing.IsActive,
			&listing.Views,
			&listing.Ratings,
			&listing.CreatedAt,
			&listing.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing row: %w", err)
		}

		listings = append(listings, &listing)
	}

	return listings, nil
}

// AddCategories adds categories to a listing
func (r *ListingRepository) AddCategories(ctx context.Context, listingID uuid.UUID, categoryIDs []uuid.UUID) error {
	if len(categoryIDs) == 0 {
		return nil
	}

	// Build the query with multiple values
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("INSERT INTO marketplace_listing_category (listing_id, category_id) VALUES ")

	args := []interface{}{listingID}
	placeholders := make([]string, len(categoryIDs))
	for i, categoryID := range categoryIDs {
		placeholders[i] = fmt.Sprintf("($1, $%d)", i+2)
		args = append(args, categoryID)
	}

	queryBuilder.WriteString(strings.Join(placeholders, ", "))
	queryBuilder.WriteString(" ON CONFLICT DO NOTHING")

	_, err := r.db.ExecContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return fmt.Errorf("failed to add categories: %w", err)
	}

	return nil
}

// RemoveCategories removes categories from a listing
func (r *ListingRepository) RemoveCategories(ctx context.Context, listingID uuid.UUID, categoryIDs []uuid.UUID) error {
	if len(categoryIDs) == 0 {
		return nil
	}

	query := `DELETE FROM marketplace_listing_category WHERE listing_id = $1 AND category_id = ANY($2)`
	_, err := r.db.ExecContext(ctx, query, listingID, pq.Array(categoryIDs))
	if err != nil {
		return fmt.Errorf("failed to remove categories: %w", err)
	}

	return nil
}

// GetCategories gets all categories for a listing
func (r *ListingRepository) GetCategories(ctx context.Context, listingID uuid.UUID) ([]model.Category, error) {
	query := `
		SELECT 
			c.id, c.name, c.description, c.icon, c.created_at, c.updated_at
		FROM marketplace_category c
		JOIN marketplace_listing_category lc ON c.id = lc.category_id
		WHERE lc.listing_id = $1
		ORDER BY c.name
	`

	rows, err := r.db.QueryContext(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Icon,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category row: %w", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// SetFeatures sets the feature list for a listing
func (r *ListingRepository) SetFeatures(ctx context.Context, listingID uuid.UUID, features []string) error {
	// Convert features to JSON
	featuresJSON, err := json.Marshal(features)
	if err != nil {
		return fmt.Errorf("failed to marshal features: %w", err)
	}

	// Check if features exist
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM marketplace_listing_features WHERE listing_id = $1)`
	err = r.db.QueryRowContext(ctx, query, listingID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if features exist: %w", err)
	}

	// Insert or update features
	if exists {
		query = `UPDATE marketplace_listing_features SET features = $1 WHERE listing_id = $2`
	} else {
		query = `INSERT INTO marketplace_listing_features (features, listing_id) VALUES ($1, $2)`
	}

	_, err = r.db.ExecContext(ctx, query, featuresJSON, listingID)
	if err != nil {
		return fmt.Errorf("failed to set features: %w", err)
	}

	return nil
}

// GetFeatures gets the feature list for a listing
func (r *ListingRepository) GetFeatures(ctx context.Context, listingID uuid.UUID) ([]string, error) {
	query := `SELECT features FROM marketplace_listing_features WHERE listing_id = $1`

	var featuresJSON []byte
	err := r.db.QueryRowContext(ctx, query, listingID).Scan(&featuresJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No features found, return empty slice
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to get features: %w", err)
	}

	// Parse features JSON
	var features []string
	err = json.Unmarshal(featuresJSON, &features)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal features: %w", err)
	}

	return features, nil
}
