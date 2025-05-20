// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/post_repo.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// PostRepository implements the repository.PostRepository interface with PostgreSQL
type PostRepository struct {
	db *sql.DB
}

// Create creates a new post in the database
func (r *PostRepository) Create(ctx context.Context, post *model.Post) error {
	query := `
		INSERT INTO posts (
			topic_id, author_id, content, is_solution, parent_post_id,
			created_at, updated_at, edited_by, edit_count
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		post.TopicID,
		post.AuthorID,
		post.Content,
		post.IsSolution,
		post.ParentPostID,
		post.CreatedAt,
		post.UpdatedAt,
		post.EditedByID,
		post.EditCount,
	).Scan(&post.ID)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// GetByID retrieves a post by its ID
func (r *PostRepository) GetByID(ctx context.Context, id int) (*model.Post, error) {
	query := `
		SELECT p.id, p.topic_id, p.author_id, p.content, p.is_solution, p.parent_post_id,
			p.created_at, p.updated_at, p.edited_by, p.edit_count,
			u.name AS author_name,
			e.name AS edited_by_name
		FROM posts p
		JOIN "user" u ON p.author_id = u.id
		LEFT JOIN "user" e ON p.edited_by = e.id
		WHERE p.id = $1
	`

	var post model.Post
	var parentPostID sql.NullInt32
	var editedByID sql.NullString
	var editedByName sql.NullString
	var authorName string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.TopicID,
		&post.AuthorID,
		&post.Content,
		&post.IsSolution,
		&parentPostID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&editedByID,
		&post.EditCount,
		&authorName,
		&editedByName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Set author
	post.Author = model.User{
		ID:   post.AuthorID,
		Name: authorName,
	}

	// Handle nullable fields
	if parentPostID.Valid {
		id := int(parentPostID.Int32)
		post.ParentPostID = &id
	}

	if editedByID.Valid {
		userID, err := uuid.Parse(editedByID.String)
		if err == nil {
			post.EditedByID = &userID

			if editedByName.Valid {
				post.EditedBy = &model.User{
					ID:   userID,
					Name: editedByName.String,
				}
			}
		}
	}

	// Get like count
	likeCount, err := r.GetLikeCount(ctx, id)
	if err == nil {
		post.LikeCount = likeCount
	}

	return &post, nil
}

// GetByTopicID retrieves posts for a specific topic
func (r *PostRepository) GetByTopicID(ctx context.Context, topicID int, page, perPage int) ([]*model.Post, int, error) {
	// Count total matching rows
	countQuery := "SELECT COUNT(*) FROM posts WHERE topic_id = $1"
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, topicID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Calculate limit and offset
	offset := (page - 1) * perPage

	// Build query
	query := `
		SELECT p.id, p.topic_id, p.author_id, p.content, p.is_solution, p.parent_post_id,
			p.created_at, p.updated_at, p.edited_by, p.edit_count,
			u.name AS author_name,
			e.name AS edited_by_name
		FROM posts p
		JOIN "user" u ON p.author_id = u.id
		LEFT JOIN "user" e ON p.edited_by = e.id
		WHERE p.topic_id = $1
		ORDER BY 
			CASE WHEN p.parent_post_id IS NULL THEN 0 ELSE 1 END, -- Show parent posts first
			p.created_at ASC -- Then in chronological order
		LIMIT $2 OFFSET $3
	`

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, topicID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	// Parse results
	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		var parentPostID sql.NullInt32
		var editedByID sql.NullString
		var editedByName sql.NullString
		var authorName string

		err := rows.Scan(
			&post.ID,
			&post.TopicID,
			&post.AuthorID,
			&post.Content,
			&post.IsSolution,
			&parentPostID,
			&post.CreatedAt,
			&post.UpdatedAt,
			&editedByID,
			&post.EditCount,
			&authorName,
			&editedByName,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post row: %w", err)
		}

		// Set author
		post.Author = model.User{
			ID:   post.AuthorID,
			Name: authorName,
		}

		// Handle nullable fields
		if parentPostID.Valid {
			id := int(parentPostID.Int32)
			post.ParentPostID = &id
		}

		if editedByID.Valid {
			userID, err := uuid.Parse(editedByID.String)
			if err == nil {
				post.EditedByID = &userID

				if editedByName.Valid {
					post.EditedBy = &model.User{
						ID:   userID,
						Name: editedByName.String,
					}
				}
			}
		}

		posts = append(posts, &post)
	}

	// Check for errors during row iteration
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during post row iteration: %w", err)
	}

	// Get like counts for all posts
	for _, post := range posts {
		likeCount, err := r.GetLikeCount(ctx, post.ID)
		if err == nil {
			post.LikeCount = likeCount
		}
	}

	return posts, total, nil
}

// GetByUserID retrieves posts by a specific user
func (r *PostRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, perPage int) ([]*model.Post, int, error) {
	// Count total matching rows
	countQuery := "SELECT COUNT(*) FROM posts WHERE author_id = $1"
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Calculate limit and offset
	offset := (page - 1) * perPage

	// Build query
	query := `
		SELECT p.id, p.topic_id, p.author_id, p.content, p.is_solution, p.parent_post_id,
			p.created_at, p.updated_at, p.edited_by, p.edit_count,
			u.name AS author_name,
			e.name AS edited_by_name,
			t.title AS topic_title
		FROM posts p
		JOIN "user" u ON p.author_id = u.id
		JOIN topics t ON p.topic_id = t.id
		LEFT JOIN "user" e ON p.edited_by = e.id
		WHERE p.author_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	// Parse results
	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		var parentPostID sql.NullInt32
		var editedByID sql.NullString
		var editedByName sql.NullString
		var authorName string
		var topicTitle string

		err := rows.Scan(
			&post.ID,
			&post.TopicID,
			&post.AuthorID,
			&post.Content,
			&post.IsSolution,
			&parentPostID,
			&post.CreatedAt,
			&post.UpdatedAt,
			&editedByID,
			&post.EditCount,
			&authorName,
			&editedByName,
			&topicTitle,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post row: %w", err)
		}

		// Set author
		post.Author = model.User{
			ID:   post.AuthorID,
			Name: authorName,
		}

		// Handle nullable fields
		if parentPostID.Valid {
			id := int(parentPostID.Int32)
			post.ParentPostID = &id
		}

		if editedByID.Valid {
			userID, err := uuid.Parse(editedByID.String)
			if err == nil {
				post.EditedByID = &userID

				if editedByName.Valid {
					post.EditedBy = &model.User{
						ID:   userID,
						Name: editedByName.String,
					}
				}
			}
		}

		// Set minimal topic info
		post.Topic = &model.Topic{
			ID:    post.TopicID,
			Title: topicTitle,
		}

		posts = append(posts, &post)
	}

	// Check for errors during row iteration
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during post row iteration: %w", err)
	}

	// Get like counts for all posts
	for _, post := range posts {
		likeCount, err := r.GetLikeCount(ctx, post.ID)
		if err == nil {
			post.LikeCount = likeCount
		}
	}

	return posts, total, nil
}

// Update updates a post
func (r *PostRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
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

	query := "UPDATE posts SET " + setClause + " WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

// Delete deletes a post
func (r *PostRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// GetLikeCount gets the like count for a post
func (r *PostRepository) GetLikeCount(ctx context.Context, postID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM likes WHERE likeable_type = 'post' AND likeable_id = $1",
		postID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get like count: %w", err)
	}

	return count, nil
}

// IsAuthor checks if a user is the author of a post
func (r *PostRepository) IsAuthor(ctx context.Context, postID int, userID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM posts WHERE id = $1 AND author_id = $2",
		postID,
		userID,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check post authorship: %w", err)
	}

	return count > 0, nil
}
