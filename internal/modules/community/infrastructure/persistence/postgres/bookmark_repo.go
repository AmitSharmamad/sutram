// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/bookmark_repo.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// BookmarkRepository implements the repository.BookmarkRepository interface with PostgreSQL
type BookmarkRepository struct {
	db *sql.DB
}

// Create creates a new bookmark
func (r *BookmarkRepository) Create(ctx context.Context, bookmark *model.Bookmark) error {
	query := `
		INSERT INTO bookmarks (user_id, topic_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		bookmark.UserID,
		bookmark.TopicID,
		bookmark.CreatedAt,
	).Scan(&bookmark.ID)

	if err != nil {
		return fmt.Errorf("failed to create bookmark: %w", err)
	}

	return nil
}

// Delete deletes a bookmark
func (r *BookmarkRepository) Delete(ctx context.Context, userID uuid.UUID, topicID int) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM bookmarks WHERE user_id = $1 AND topic_id = $2",
		userID,
		topicID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	return nil
}

// Exists checks if a bookmark exists
func (r *BookmarkRepository) Exists(ctx context.Context, userID uuid.UUID, topicID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = $1 AND topic_id = $2)",
		userID,
		topicID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check bookmark existence: %w", err)
	}

	return exists, nil
}

// GetByUserID gets bookmarks for a user
func (r *BookmarkRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page int, perPage int) ([]*model.Topic, int, error) {
	// Count total matching rows
	countQuery := "SELECT COUNT(*) FROM bookmarks WHERE user_id = $1"
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user bookmarks: %w", err)
	}

	// Calculate limit and offset
	offset := (page - 1) * perPage

	// Build query
	query := `
		SELECT t.id, t.title, t.slug, t.content, t.category_id, t.author_id, t.is_pinned, t.is_closed, 
			t.is_hidden, t.view_count, t.created_at, t.updated_at, t.last_activity_at,
			c.name AS category_name, c.slug AS category_slug,
			u.name AS author_name,
			COALESCE(ts.reply_count, 0) AS reply_count,
			COALESCE(ts.like_count, 0) AS like_count,
			COALESCE(ts.bookmark_count, 0) AS bookmark_count
		FROM topics t
		JOIN bookmarks b ON t.id = b.topic_id AND b.user_id = $1
		JOIN categories c ON t.category_id = c.id
		JOIN "user" u ON t.author_id = u.id
		LEFT JOIN topic_stats ts ON t.id = ts.topic_id
		WHERE t.is_hidden = false
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query user bookmarks: %w", err)
	}
	defer rows.Close()

	// Parse results
	var topics []*model.Topic
	for rows.Next() {
		var topic model.Topic
		var categoryName, categorySlug, authorName string
		var replyCount, likeCount, bookmarkCount int

		err := rows.Scan(
			&topic.ID,
			&topic.Title,
			&topic.Slug,
			&topic.Content,
			&topic.CategoryID,
			&topic.AuthorID,
			&topic.IsPinned,
			&topic.IsClosed,
			&topic.IsHidden,
			&topic.ViewCount,
			&topic.CreatedAt,
			&topic.UpdatedAt,
			&topic.LastActivityAt,
			&categoryName,
			&categorySlug,
			&authorName,
			&replyCount,
			&likeCount,
			&bookmarkCount,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan topic row: %w", err)
		}

		// Set related objects
		topic.Category = model.Category{
			ID:   topic.CategoryID,
			Name: categoryName,
			Slug: categorySlug,
		}

		topic.Author = model.User{
			ID:   topic.AuthorID,
			Name: authorName,
		}

		topic.Stats = &model.TopicStats{
			TopicID:       topic.ID,
			ReplyCount:    replyCount,
			LikeCount:     likeCount,
			BookmarkCount: bookmarkCount,
		}

		topics = append(topics, &topic)
	}

	// Check for errors during row iteration
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during bookmarked topic row iteration: %w", err)
	}

	return topics, total, nil
}

// GetBookmarkCount gets the bookmark count for a topic
func (r *BookmarkRepository) GetBookmarkCount(ctx context.Context, topicID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM bookmarks WHERE topic_id = $1",
		topicID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get bookmark count: %w", err)
	}

	return count, nil
}
