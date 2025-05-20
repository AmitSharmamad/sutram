// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/category_repo.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"algo/internal/modules/community/domain/model"
)

// CategoryRepository implements the repository.CategoryRepository interface with PostgreSQL
type CategoryRepository struct {
	db *sql.DB
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll(ctx context.Context) ([]*model.Category, error) {
	query := `
		SELECT id, name, slug, description, icon, color, parent_id, display_order, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = true
		ORDER BY display_order ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		var category model.Category
		var icon, color sql.NullString
		var parentID sql.NullInt32

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&icon,
			&color,
			&parentID,
			&category.DisplayOrder,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan category row: %w", err)
		}

		// Handle nullable fields
		if icon.Valid {
			category.Icon = &icon.String
		}

		if color.Valid {
			category.Color = &color.String
		}

		if parentID.Valid {
			id := int(parentID.Int32)
			category.ParentID = &id
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during category row iteration: %w", err)
	}

	return categories, nil
}

// GetByID retrieves a category by its ID
func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*model.Category, error) {
	query := `
		SELECT id, name, slug, description, icon, color, parent_id, display_order, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	var category model.Category
	var icon, color sql.NullString
	var parentID sql.NullInt32

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&icon,
		&color,
		&parentID,
		&category.DisplayOrder,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Handle nullable fields
	if icon.Valid {
		category.Icon = &icon.String
	}

	if color.Valid {
		category.Color = &color.String
	}

	if parentID.Valid {
		id := int(parentID.Int32)
		category.ParentID = &id
	}

	return &category, nil
}

// GetBySlug retrieves a category by its slug
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	query := `
		SELECT id, name, slug, description, icon, color, parent_id, display_order, is_active, created_at, updated_at
		FROM categories
		WHERE slug = $1
	`

	var category model.Category
	var icon, color sql.NullString
	var parentID sql.NullInt32

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&icon,
		&color,
		&parentID,
		&category.DisplayOrder,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	// Handle nullable fields
	if icon.Valid {
		category.Icon = &icon.String
	}

	if color.Valid {
		category.Color = &color.String
	}

	if parentID.Valid {
		id := int(parentID.Int32)
		category.ParentID = &id
	}

	return &category, nil
}

// Create creates a new category
func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) error {
	query := `
		INSERT INTO categories (name, slug, description, icon, color, parent_id, display_order, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.Icon,
		category.Color,
		category.ParentID,
		category.DisplayOrder,
		category.IsActive,
		category.CreatedAt,
		category.UpdatedAt,
	).Scan(&category.ID)

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

// Update updates a category
func (r *CategoryRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	// Build dynamic query
	setClause := ""
	args := []interface{}{id} // First arg is the ID
	argCount := 2

	for key, value := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += key + " = $" + fmt.Sprintf("%d", argCount)
		args = append(args, value)
		argCount++
	}

	if setClause == "" {
		return nil // Nothing to update
	}

	query := "UPDATE categories SET " + setClause + " WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	return nil
}

// Delete deletes a category
func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// GetTopicCount gets the number of topics in a category
func (r *CategoryRepository) GetTopicCount(ctx context.Context, categoryID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM topics WHERE category_id = $1 AND is_hidden = false",
		categoryID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get topic count: %w", err)
	}

	return count, nil
}
