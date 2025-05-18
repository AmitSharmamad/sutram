// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/topic_repo.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// TopicRepository implements the repository.TopicRepository interface with PostgreSQL
type TopicRepository struct {
	db *sql.DB
}

// NewTopicRepository is defined in repository_factory.go
// This implementation satisfies the repository.TopicRepository interface

// Create creates a new topic in the database
func (r *TopicRepository) Create(ctx context.Context, topic *model.Topic) error {
	query := `
		INSERT INTO topics (
			title, slug, content, category_id, author_id, is_pinned, is_closed, is_hidden, view_count, 
			created_at, updated_at, last_activity_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		topic.Title,
		topic.Slug,
		topic.Content,
		topic.CategoryID,
		topic.AuthorID,
		topic.IsPinned,
		topic.IsClosed,
		topic.IsHidden,
		topic.ViewCount,
		topic.CreatedAt,
		topic.UpdatedAt,
		topic.LastActivityAt,
	).Scan(&topic.ID)

	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}

// GetByID retrieves a topic by its ID
func (r *TopicRepository) GetByID(ctx context.Context, id int) (*model.Topic, error) {
	query := `
		SELECT t.id, t.title, t.slug, t.content, t.category_id, t.author_id, t.is_pinned, t.is_closed, 
			t.is_hidden, t.view_count, t.created_at, t.updated_at, t.last_activity_at,
			c.name AS category_name, c.slug AS category_slug, c.description AS category_description,
			u.name AS author_name
		FROM topics t
		JOIN categories c ON t.category_id = c.id
		JOIN "user" u ON t.author_id = u.id
		WHERE t.id = $1
	`

	var topic model.Topic
	var categoryName, categorySlug, categoryDesc, authorName string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		&categoryDesc,
		&authorName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	// Set related objects
	topic.Category = model.Category{
		ID:          topic.CategoryID,
		Name:        categoryName,
		Slug:        categorySlug,
		Description: categoryDesc,
	}

	topic.Author = model.User{
		ID:   topic.AuthorID,
		Name: authorName,
	}

	// Get topic stats
	stats, err := r.GetTopicStats(ctx, id)
	if err == nil {
		topic.Stats = stats
	}

	// Get tags
	tags, err := r.GetTags(ctx, id)
	if err == nil {
		topic.Tags = tags
	}

	return &topic, nil
}

// GetTopics retrieves topics with optional filtering
func (r *TopicRepository) GetTopics(
	ctx context.Context,
	categorySlug string,
	tagSlug string,
	sort string,
	search string,
	page int,
	perPage int,
) ([]*model.Topic, int, error) {
	// Build query
	baseQuery := `
		FROM topics t
		JOIN categories c ON t.category_id = c.id
		JOIN "user" u ON t.author_id = u.id
		LEFT JOIN topic_stats ts ON t.id = ts.topic_id
	`

	// Add WHERE clauses
	whereClause := "WHERE t.is_hidden = false"
	args := []interface{}{}
	argIndex := 1

	if categorySlug != "" {
		whereClause += fmt.Sprintf(" AND c.slug = $%d", argIndex)
		args = append(args, categorySlug)
		argIndex++
	}

	if tagSlug != "" {
		whereClause += fmt.Sprintf(`
			AND t.id IN (
				SELECT tt.topic_id 
				FROM topic_tags tt 
				JOIN tags tg ON tt.tag_id = tg.id 
				WHERE tg.slug = $%d
			)`, argIndex)
		args = append(args, tagSlug)
		argIndex++
	}

	if search != "" {
		whereClause += fmt.Sprintf(" AND (t.title ILIKE $%d OR t.content ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	// Add ORDER BY clause
	orderClause := "ORDER BY "
	switch sort {
	case "newest":
		orderClause += "t.created_at DESC"
	case "active":
		orderClause += "t.last_activity_at DESC"
	case "popular":
		orderClause += "t.view_count DESC"
	case "replies":
		orderClause += "ts.reply_count DESC"
	default:
		// Default to pinned first, then most recent activity
		orderClause += "t.is_pinned DESC, t.last_activity_at DESC"
	}

	// Count total matching rows
	countQuery := "SELECT COUNT(*) " + baseQuery + " " + whereClause
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count topics: %w", err)
	}

	// Calculate limit and offset
	offset := (page - 1) * perPage

	// Add pagination
	limitOffset := fmt.Sprintf("LIMIT %d OFFSET %d", perPage, offset)

	// Build final query
	query := `
		SELECT t.id, t.title, t.slug, t.content, t.category_id, t.author_id, t.is_pinned, t.is_closed, 
			t.is_hidden, t.view_count, t.created_at, t.updated_at, t.last_activity_at,
			c.name AS category_name, c.slug AS category_slug,
			u.name AS author_name,
			COALESCE(ts.reply_count, 0) AS reply_count,
			COALESCE(ts.like_count, 0) AS like_count,
			COALESCE(ts.bookmark_count, 0) AS bookmark_count
	` + baseQuery + " " + whereClause + " " + orderClause + " " + limitOffset

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query topics: %w", err)
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
		return nil, 0, fmt.Errorf("error during topic row iteration: %w", err)
	}

	// Load tags for each topic
	for _, topic := range topics {
		tags, err := r.GetTags(ctx, topic.ID)
		if err == nil {
			topic.Tags = tags
		}
	}

	return topics, total, nil
}

// GetTopicsByUserID retrieves topics created by a specific user
func (r *TopicRepository) GetTopicsByUserID(
	ctx context.Context,
	userID uuid.UUID,
	page int,
	perPage int,
) ([]*model.Topic, int, error) {
	// Count total matching rows
	countQuery := "SELECT COUNT(*) FROM topics WHERE author_id = $1 AND is_hidden = false"
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user topics: %w", err)
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
		JOIN categories c ON t.category_id = c.id
		JOIN "user" u ON t.author_id = u.id
		LEFT JOIN topic_stats ts ON t.id = ts.topic_id
		WHERE t.author_id = $1 AND t.is_hidden = false
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query user topics: %w", err)
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
		return nil, 0, fmt.Errorf("error during topic row iteration: %w", err)
	}

	// Load tags for each topic
	for _, topic := range topics {
		tags, err := r.GetTags(ctx, topic.ID)
		if err == nil {
			topic.Tags = tags
		}
	}

	return topics, total, nil
}

// Update updates a topic
func (r *TopicRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	// Build SET clause
	setClauses := []string{}
	args := []interface{}{id} // First argument is the ID for the WHERE clause
	argCount := 2             // Start from 2 since 1 is used for the ID

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argCount))
		args = append(args, value)
		argCount++
	}

	if len(setClauses) == 0 {
		return nil // Nothing to update
	}

	query := fmt.Sprintf(
		"UPDATE topics SET %s WHERE id = $1",
		strings.Join(setClauses, ", "),
	)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update topic: %w", err)
	}

	return nil
}

// Delete deletes a topic
func (r *TopicRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM topics WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	return nil
}

// IncrementViews increments the view count for a topic
func (r *TopicRepository) IncrementViews(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE topics SET view_count = view_count + 1 WHERE id = $1",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}

	return nil
}

// GetTopicStats retrieves statistics for a topic
func (r *TopicRepository) GetTopicStats(ctx context.Context, id int) (*model.TopicStats, error) {
	query := `
		SELECT ts.topic_id, ts.reply_count, ts.like_count, ts.bookmark_count, 
			ts.last_reply_id, ts.last_reply_at, ts.last_reply_user_id,
			u.name AS last_reply_user_name
		FROM topic_stats ts
		LEFT JOIN "user" u ON ts.last_reply_user_id = u.id
		WHERE ts.topic_id = $1
	`

	var stats model.TopicStats
	var lastReplyID sql.NullInt32
	var lastReplyAt sql.NullTime
	var lastReplyUserID sql.NullString
	var lastReplyUserName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&stats.TopicID,
		&stats.ReplyCount,
		&stats.LikeCount,
		&stats.BookmarkCount,
		&lastReplyID,
		&lastReplyAt,
		&lastReplyUserID,
		&lastReplyUserName,
	)

	if err == sql.ErrNoRows {
		// Create empty stats for this topic
		return &model.TopicStats{
			TopicID:       id,
			ReplyCount:    0,
			LikeCount:     0,
			BookmarkCount: 0,
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get topic stats: %w", err)
	}

	// Handle nullable fields
	if lastReplyID.Valid {
		id := int(lastReplyID.Int32)
		stats.LastReplyID = &id
	}

	if lastReplyAt.Valid {
		stats.LastReplyAt = &lastReplyAt.Time
	}

	if lastReplyUserID.Valid {
		userID, err := uuid.Parse(lastReplyUserID.String)
		if err == nil {
			stats.LastReplyUserID = &userID
		}
	}

	if lastReplyUserName.Valid && stats.LastReplyUserID != nil {
		stats.LastReplyUser = model.User{
			ID:   *stats.LastReplyUserID,
			Name: lastReplyUserName.String,
		}
	}

	return &stats, nil
}

// AddTags adds tags to a topic
func (r *TopicRepository) AddTags(ctx context.Context, topicID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO topic_tags (topic_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, tagID := range tagIDs {
		_, err := stmt.ExecContext(ctx, topicID, tagID)
		if err != nil {
			return fmt.Errorf("failed to add tag %d: %w", tagID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RemoveTags removes tags from a topic
func (r *TopicRepository) RemoveTags(ctx context.Context, topicID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(tagIDs))
	args := []interface{}{topicID}

	for i, tagID := range tagIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args = append(args, tagID)
	}

	query := fmt.Sprintf(
		"DELETE FROM topic_tags WHERE topic_id = $1 AND tag_id IN (%s)",
		strings.Join(placeholders, ", "),
	)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to remove tags: %w", err)
	}

	return nil
}

// SetTags replaces all tags for a topic
func (r *TopicRepository) SetTags(ctx context.Context, topicID int, tagIDs []int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Remove all existing tags
	_, err = tx.ExecContext(ctx, "DELETE FROM topic_tags WHERE topic_id = $1", topicID)
	if err != nil {
		return fmt.Errorf("failed to remove existing tags: %w", err)
	}

	// Add new tags if any
	if len(tagIDs) > 0 {
		stmt, err := tx.PrepareContext(ctx, "INSERT INTO topic_tags (topic_id, tag_id) VALUES ($1, $2)")
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for _, tagID := range tagIDs {
			_, err := stmt.ExecContext(ctx, topicID, tagID)
			if err != nil {
				return fmt.Errorf("failed to add tag %d: %w", tagID, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTags gets tags for a topic
func (r *TopicRepository) GetTags(ctx context.Context, topicID int) ([]*model.Tag, error) {
	query := `
		SELECT t.id, t.name, t.slug, t.description, t.color, t.is_active, t.created_at
		FROM tags t
		JOIN topic_tags tt ON t.id = tt.tag_id
		WHERE tt.topic_id = $1
		ORDER BY t.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, topicID)
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

// IsAuthor checks if a user is the author of a topic
func (r *TopicRepository) IsAuthor(ctx context.Context, topicID int, userID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM topics WHERE id = $1 AND author_id = $2",
		topicID,
		userID,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check topic authorship: %w", err)
	}

	return count > 0, nil
}
