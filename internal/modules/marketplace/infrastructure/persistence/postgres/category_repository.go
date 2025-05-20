// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/infrastructure/persistence/postgres/category_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// CategoryRepository implements the repository.CategoryRepository interface
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// FindAll returns all categories
func (r *CategoryRepository) FindAll(ctx context.Context) ([]model.Category, error) {
	query := `
		SELECT 
			id, name, description, icon, created_at, updated_at
		FROM marketplace_category
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
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

// FindByID returns a category by its ID
func (r *CategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	query := `
		SELECT 
			id, name, description, icon, created_at, updated_at
		FROM marketplace_category
		WHERE id = $1
	`

	var category model.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Icon,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	return &category, nil
}

// FindByName returns a category by its name
func (r *CategoryRepository) FindByName(ctx context.Context, name string) (*model.Category, error) {
	query := `
		SELECT 
			id, name, description, icon, created_at, updated_at
		FROM marketplace_category
		WHERE LOWER(name) = LOWER($1)
	`

	var category model.Category
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Icon,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find category by name: %w", err)
	}

	return &category, nil
}

// FindByIDs returns categories by their IDs
func (r *CategoryRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Category, error) {
	if len(ids) == 0 {
		return []model.Category{}, nil
	}

	// Build placeholder string for the query
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT 
			id, name, description, icon, created_at, updated_at
		FROM marketplace_category
		WHERE id IN (%s)
		ORDER BY name
	`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories by IDs: %w", err)
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
