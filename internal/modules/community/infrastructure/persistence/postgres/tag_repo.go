// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/tag_repo.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"algo/internal/modules/community/domain/model"
)

// TagRepository implements the repository.TagRepository interface with PostgreSQL
type TagRepository struct {
	db *sql.DB
}

// GetAll retrieves all tags
func (r *TagRepository) GetAll(ctx context.Context) ([]*model.Tag, error) {
	query := `
		SELECT id, name, slug, description, color, is_active, created_at
		FROM tags
		WHERE is_active = true
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	var tags []*model.Tag
	for rows.Next() {
		var tag model.Tag
		var color sql.NullString

		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Slug,
			&tag.Description,
			&color,
			&tag.IsActive,
			&tag.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan tag row: %w", err)
		}

		// Handle nullable fields
		if color.Valid {
			tag.Color = &color.String
		}

		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during tag row iteration: %w", err)
	}

	return tags, nil
}

// GetByID retrieves a tag by its ID
func (r *TagRepository) GetByID(ctx context.Context, id int) (*model.Tag, error) {
	query := `
		SELECT id, name, slug, description, color, is_active, created_at
		FROM tags
		WHERE id = $1
	`

	var tag model.Tag
	var color sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Slug,
		&tag.Description,
		&color,
		&tag.IsActive,
		&tag.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	// Handle nullable fields
	if color.Valid {
		tag.Color = &color.String
	}

	return &tag, nil
}

// GetBySlug retrieves a tag by its slug
func (r *TagRepository) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	query := `
		SELECT id, name, slug, description, color, is_active, created_at
		FROM tags
		WHERE slug = $1
	`

	var tag model.Tag
	var color sql.NullString

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Slug,
		&tag.Description,
		&color,
		&tag.IsActive,
		&tag.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get tag by slug: %w", err)
	}

	// Handle nullable fields
	if color.Valid {
		tag.Color = &color.String
	}

	return &tag, nil
}

// Create creates a new tag
func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) error {
	query := `
		INSERT INTO tags (name, slug, description, color, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		tag.Name,
		tag.Slug,
		tag.Description,
		tag.Color,
		tag.IsActive,
		tag.CreatedAt,
	).Scan(&tag.ID)

	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

// Update updates a tag
func (r *TagRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
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

	query := "UPDATE tags SET " + setClause + " WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	return nil
}

// Delete deletes a tag
func (r *TagRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tags WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// GetTagsByTopicID retrieves tags for a topic
func (r *TagRepository) GetTagsByTopicID(ctx context.Context, topicID int) ([]*model.Tag, error) {
	query := `
		SELECT t.id, t.name, t.slug, t.description, t.color, t.is_active, t.created_at
		FROM tags t
		JOIN topic_tags tt ON t.id = tt.tag_id
		WHERE tt.topic_id = $1
		ORDER BY t.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to query topic tags: %w", err)
	}
	defer rows.Close()

	var tags []*model.Tag
	for rows.Next() {
		var tag model.Tag
		var color sql.NullString

		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Slug,
			&tag.Description,
			&color,
			&tag.IsActive,
			&tag.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan tag row: %w", err)
		}

		// Handle nullable fields
		if color.Valid {
			tag.Color = &color.String
		}

		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during tag row iteration: %w", err)
	}

	return tags, nil
}
