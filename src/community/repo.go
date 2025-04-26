package community

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"algo/src/shared"
)

// Repository defines the data access layer for community features
type Repository interface {
	// Categories
	GetCategories(ctx context.Context) ([]shared.Category, error)
	GetCategoryByID(ctx context.Context, id int) (*shared.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*shared.Category, error)
	CreateCategory(ctx context.Context, category *shared.Category) error
	UpdateCategory(ctx context.Context, category *shared.Category) error
	DeleteCategory(ctx context.Context, id int) error

	// Topics
	GetTopics(ctx context.Context, filter TopicFilter) (*shared.PaginatedTopics, error)
	GetTopicByID(ctx context.Context, id int) (*shared.Topic, error)
	GetTopicBySlug(ctx context.Context, slug string) (*shared.Topic, error)
	CreateTopic(ctx context.Context, topic *shared.Topic) error
	UpdateTopic(ctx context.Context, topic *shared.Topic) error
	DeleteTopic(ctx context.Context, id int) error
	IncrementTopicViewCount(ctx context.Context, id int) error

	// Posts
	GetPostsByTopicID(ctx context.Context, topicID int, filter PostFilter) (*shared.PaginatedPosts, error)
	GetPostByID(ctx context.Context, id int) (*shared.Post, error)
	CreatePost(ctx context.Context, post *shared.Post) error
	UpdatePost(ctx context.Context, post *shared.Post) error
	DeletePost(ctx context.Context, id int) error
	MarkPostAsSolution(ctx context.Context, postID int) error

	// Tags
	GetTags(ctx context.Context) ([]shared.Tag, error)
	GetTagByID(ctx context.Context, id int) (*shared.Tag, error)
	GetTagBySlug(ctx context.Context, slug string) (*shared.Tag, error)
	CreateTag(ctx context.Context, tag *shared.Tag) error
	UpdateTag(ctx context.Context, tag *shared.Tag) error
	DeleteTag(ctx context.Context, id int) error
	AssignTagsToTopic(ctx context.Context, topicID int, tagIDs []int) error
	RemoveTagFromTopic(ctx context.Context, topicID, tagID int) error

	// Likes
	GetLikesByEntityID(ctx context.Context, entityType string, entityID int) ([]shared.Like, error)
	GetLikeByUserAndEntity(ctx context.Context, userID uuid.UUID, entityType string, entityID int) (*shared.Like, error)
	CreateLike(ctx context.Context, like *shared.Like) error
	DeleteLike(ctx context.Context, id int) error

	// Bookmarks
	GetBookmarksByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.Topic, int, error)
	GetBookmarkByUserAndTopic(ctx context.Context, userID uuid.UUID, topicID int) (*shared.Bookmark, error)
	CreateBookmark(ctx context.Context, bookmark *shared.Bookmark) error
	DeleteBookmark(ctx context.Context, id int) error

	// User Follows
	GetFollowersByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.User, int, error)
	GetFollowingByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.User, int, error)
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	CreateUserFollow(ctx context.Context, follow *shared.UserFollow) error
	DeleteUserFollow(ctx context.Context, followerID, followingID uuid.UUID) error

	// Category Follows
	GetCategoryFollowsByUserID(ctx context.Context, userID uuid.UUID) ([]shared.Category, error)
	IsCategoryFollowed(ctx context.Context, userID uuid.UUID, categoryID int) (bool, error)
	CreateCategoryFollow(ctx context.Context, follow *shared.CategoryFollow) error
	DeleteCategoryFollow(ctx context.Context, userID uuid.UUID, categoryID int) error

	// Topic Follows
	GetTopicFollowsByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.Topic, int, error)
	IsTopicFollowed(ctx context.Context, userID uuid.UUID, topicID int) (bool, error)
	CreateTopicFollow(ctx context.Context, follow *shared.TopicFollow) error
	DeleteTopicFollow(ctx context.Context, userID uuid.UUID, topicID int) error

	// Notifications
	GetNotificationsByUserID(ctx context.Context, userID uuid.UUID, page, limit int, unreadOnly bool) ([]shared.Notification, int, error)
	GetNotificationByID(ctx context.Context, id int) (*shared.Notification, error)
	CreateNotification(ctx context.Context, notification *shared.Notification) error
	MarkNotificationAsRead(ctx context.Context, id int) error
	MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error
	DeleteNotification(ctx context.Context, id int) error

	// Badges
	GetBadges(ctx context.Context) ([]shared.Badge, error)
	GetBadgeByID(ctx context.Context, id int) (*shared.Badge, error)
	GetUserBadges(ctx context.Context, userID uuid.UUID) ([]shared.UserBadge, error)
	AwardBadgeToUser(ctx context.Context, userID uuid.UUID, badgeID int) error

	// Reports
	CreateReport(ctx context.Context, report *shared.Report) error
	GetReports(ctx context.Context, status string, page, limit int) ([]shared.Report, int, error)
	GetReportByID(ctx context.Context, id int) (*shared.Report, error)
	UpdateReportStatus(ctx context.Context, id int, status string, reviewedByID uuid.UUID) error

	// Polls
	CreatePoll(ctx context.Context, poll *shared.Poll, options []shared.PollOption) error
	GetPollByID(ctx context.Context, id int) (*shared.Poll, error)
	GetPollsByTopicID(ctx context.Context, topicID int) ([]shared.Poll, error)
	AddVoteToPoll(ctx context.Context, vote *shared.PollVote) error
	RemoveVoteFromPoll(ctx context.Context, userID uuid.UUID, pollID int) error
	GetPollResults(ctx context.Context, pollID int) ([]shared.PollOption, error)
	ClosePoll(ctx context.Context, pollID int) error

	// User Reputation
	GetUserReputation(ctx context.Context, userID uuid.UUID) (*shared.UserReputation, error)
	UpdateUserReputation(ctx context.Context, userID uuid.UUID, points int) error
	GetTopUsers(ctx context.Context, limit int) ([]shared.UserReputation, error)

	// Activity Logs
	LogActivity(ctx context.Context, log *shared.ActivityLog) error
	GetUserActivities(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.ActivityLog, int, error)
}

// TopicFilter defines filters for querying topics
type TopicFilter struct {
	CategoryID    *int
	AuthorID      *uuid.UUID
	Tags          []string
	SearchQuery   string
	SortBy        string // "newest", "popular", "activity"
	Page          int
	Limit         int
	IncludeHidden bool
}

// PostFilter defines filters for querying posts
type PostFilter struct {
	Page   int
	Limit  int
	SortBy string // "oldest", "newest"
}

// repository implements Repository interface
type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new instance of Repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

// GetCategories retrieves all categories
func (r *repository) GetCategories(ctx context.Context) ([]shared.Category, error) {
	var categories []shared.Category

	query := `
		SELECT 
			c.*,
			COUNT(DISTINCT t.id) AS topic_count
		FROM 
			categories c
		LEFT JOIN 
			topics t ON c.id = t.category_id AND t.is_hidden = false
		WHERE 
			c.is_active = true
		GROUP BY 
			c.id
		ORDER BY 
			c.display_order, c.name`

	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by ID
func (r *repository) GetCategoryByID(ctx context.Context, id int) (*shared.Category, error) {
	var category shared.Category

	query := `
		SELECT 
			c.*,
			COUNT(DISTINCT t.id) AS topic_count
		FROM 
			categories c
		LEFT JOIN 
			topics t ON c.id = t.category_id AND t.is_hidden = false
		WHERE 
			c.id = $1
		GROUP BY 
			c.id`

	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category by ID: %w", err)
	}

	return &category, nil
}

// GetCategoryBySlug retrieves a category by slug
func (r *repository) GetCategoryBySlug(ctx context.Context, slug string) (*shared.Category, error) {
	var category shared.Category

	query := `
		SELECT 
			c.*,
			COUNT(DISTINCT t.id) AS topic_count
		FROM 
			categories c
		LEFT JOIN 
			topics t ON c.id = t.category_id AND t.is_hidden = false
		WHERE 
			c.slug = $1
		GROUP BY 
			c.id`

	err := r.db.GetContext(ctx, &category, query, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	return &category, nil
}

// CreateCategory creates a new category
func (r *repository) CreateCategory(ctx context.Context, category *shared.Category) error {
	query := `
		INSERT INTO categories (
			name, slug, description, icon, color, parent_id, display_order, is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id, created_at, updated_at`

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
	).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

// UpdateCategory updates an existing category
func (r *repository) UpdateCategory(ctx context.Context, category *shared.Category) error {
	query := `
		UPDATE categories
		SET 
			name = $1,
			slug = $2,
			description = $3,
			icon = $4,
			color = $5,
			parent_id = $6,
			display_order = $7,
			is_active = $8,
			updated_at = NOW()
		WHERE 
			id = $9
		RETURNING updated_at`

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
		category.ID,
	).Scan(&category.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("category not found: %d", category.ID)
		}
		return fmt.Errorf("failed to update category: %w", err)
	}

	return nil
}

// DeleteCategory deletes a category
func (r *repository) DeleteCategory(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// GetTopics retrieves topics with pagination and filtering
func (r *repository) GetTopics(ctx context.Context, filter TopicFilter) (*shared.PaginatedTopics, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	offset := (filter.Page - 1) * filter.Limit

	// Build base queries
	baseCountQuery := `
		SELECT 
			COUNT(DISTINCT t.id)
		FROM 
			topics t
		LEFT JOIN 
			topic_tags tt ON t.id = tt.topic_id
		LEFT JOIN 
			tags tg ON tt.tag_id = tg.id`

	baseSelectQuery := `
		SELECT 
			t.*,
			ts.reply_count,
			ts.like_count,
			ts.bookmark_count,
			u.id as "author.id",
			u.name as "author.name",
			u.email as "author.email",
			u.profile_image as "author.profile_image",
			c.id as "category.id",
			c.name as "category.name",
			c.slug as "category.slug"
		FROM 
			topics t
		LEFT JOIN 
			topic_stats ts ON t.id = ts.topic_id
		LEFT JOIN 
			"user" u ON t.author_id = u.id
		LEFT JOIN 
			categories c ON t.category_id = c.id
		LEFT JOIN 
			topic_tags tt ON t.id = tt.topic_id
		LEFT JOIN 
			tags tg ON tt.tag_id = tg.id`

	// Build WHERE conditions
	conditions := []string{}
	args := []interface{}{}
	argCount := 1

	if !filter.IncludeHidden {
		conditions = append(conditions, "t.is_hidden = false")
	}

	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("t.category_id = $%d", argCount))
		args = append(args, *filter.CategoryID)
		argCount++
	}

	if filter.AuthorID != nil {
		conditions = append(conditions, fmt.Sprintf("t.author_id = $%d", argCount))
		args = append(args, *filter.AuthorID)
		argCount++
	}

	if len(filter.Tags) > 0 {
		placeholders := make([]string, len(filter.Tags))
		for i, tag := range filter.Tags {
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, tag)
			argCount++
		}
		conditions = append(conditions, fmt.Sprintf("tg.slug IN (%s)", strings.Join(placeholders, ",")))
	}

	if filter.SearchQuery != "" {
		conditions = append(conditions, fmt.Sprintf("(t.title ILIKE $%d OR t.content ILIKE $%d)", argCount, argCount))
		args = append(args, "%"+filter.SearchQuery+"%")
		argCount++
	}

	// Add WHERE clause if needed
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Add GROUP BY clause for topics
	groupByClause := "GROUP BY t.id, ts.reply_count, ts.like_count, ts.bookmark_count, u.id, u.name, u.email, u.profile_image, c.id, c.name, c.slug"

	// Add ORDER BY clause
	orderByClause := "ORDER BY "
	switch filter.SortBy {
	case "popular":
		orderByClause += "ts.like_count DESC, t.created_at DESC"
	case "activity":
		orderByClause += "t.last_activity_at DESC"
	default:
		orderByClause += "t.created_at DESC"
	}

	// Add pagination
	paginationClause := fmt.Sprintf("LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filter.Limit, offset)

	// Final queries
	countQuery := baseCountQuery + " " + whereClause
	selectQuery := baseSelectQuery + " " + whereClause + " " + groupByClause + " " + orderByClause + " " + paginationClause

	// Execute count query
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args[:len(args)-2]...); err != nil {
		return nil, fmt.Errorf("failed to count topics: %w", err)
	}

	// Execute select query
	rows, err := r.db.NamedQueryContext(ctx, selectQuery, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query topics: %w", err)
	}
	defer rows.Close()

	// Extract topics
	var topics []shared.Topic
	for rows.Next() {
		var topic shared.Topic
		var author shared.User
		var category shared.Category

		dest := map[string]interface{}{
			"id":                   &topic.ID,
			"title":                &topic.Title,
			"slug":                 &topic.Slug,
			"content":              &topic.Content,
			"category_id":          &topic.CategoryID,
			"author_id":            &topic.AuthorID,
			"is_pinned":            &topic.IsPinned,
			"is_closed":            &topic.IsClosed,
			"is_hidden":            &topic.IsHidden,
			"view_count":           &topic.ViewCount,
			"created_at":           &topic.CreatedAt,
			"updated_at":           &topic.UpdatedAt,
			"last_activity_at":     &topic.LastActivityAt,
			"reply_count":          &topic.ReplyCount,
			"like_count":           &topic.LikeCount,
			"bookmark_count":       &topic.BookmarkCount,
			"author.id":            &author.ID,
			"author.name":          &author.Name,
			"author.email":         &author.Email,
			"author.profile_image": &author.ProfileImage,
			"category.id":          &category.ID,
			"category.name":        &category.Name,
			"category.slug":        &category.Slug,
		}

		if err := rows.MapScan(dest); err != nil {
			return nil, fmt.Errorf("failed to scan topic: %w", err)
		}

		// Set relationship
		topic.Author = &author
		topic.Category = &category

		// Get tags for this topic
		var tags []shared.Tag
		tagQuery := `
			SELECT 
				t.id, t.name, t.slug, t.color
			FROM 
				tags t
			JOIN 
				topic_tags tt ON t.id = tt.tag_id
			WHERE 
				tt.topic_id = $1`

		if err := r.db.SelectContext(ctx, &tags, tagQuery, topic.ID); err != nil {
			return nil, fmt.Errorf("failed to get tags for topic %d: %w", topic.ID, err)
		}

		topic.Tags = tags
		topics = append(topics, topic)
	}

	return &shared.PaginatedTopics{
		Total:  total,
		Page:   filter.Page,
		Limit:  filter.Limit,
		Topics: topics,
	}, nil
}

// GetTopicByID retrieves a topic by ID
func (r *repository) GetTopicByID(ctx context.Context, id int) (*shared.Topic, error) {
	// First query for the topic and its relationships
	query := `
		SELECT 
			t.*,
			ts.reply_count,
			ts.like_count,
			ts.bookmark_count,
			u.id as "author_id",
			u.name as "author_name",
			u.email as "author_email",
			u.profile_image as "author_profile_image",
			c.id as "category_id",
			c.name as "category_name",
			c.slug as "category_slug"
		FROM 
			topics t
		LEFT JOIN 
			topic_stats ts ON t.id = ts.topic_id
		LEFT JOIN 
			"user" u ON t.author_id = u.id
		LEFT JOIN 
			categories c ON t.category_id = c.id
		WHERE 
			t.id = $1`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic by ID: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var topic shared.Topic
	var author shared.User
	var category shared.Category
	var authorName, authorEmail, authorProfileImage sql.NullString
	var categoryName, categorySlug sql.NullString

	err = rows.Scan(
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
		&topic.ReplyCount,
		&topic.LikeCount,
		&topic.BookmarkCount,
		&author.ID,
		&authorName,
		&authorEmail,
		&authorProfileImage,
		&category.ID,
		&categoryName,
		&categorySlug,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan topic: %w", err)
	}

	// Set nullable fields
	if authorName.Valid {
		author.Name = authorName.String
	}
	if authorEmail.Valid {
		author.Email = authorEmail.String
	}
	if authorProfileImage.Valid {
		author.ProfileImage = authorProfileImage.String
	}
	if categoryName.Valid {
		category.Name = categoryName.String
	}
	if categorySlug.Valid {
		category.Slug = categorySlug.String
	}

	// Set relationships
	topic.Author = &author
	topic.Category = &category

	// Get tags for this topic
	var tags []shared.Tag
	tagQuery := `
		SELECT 
			t.id, t.name, t.slug, t.color
		FROM 
			tags t
		JOIN 
			topic_tags tt ON t.id = tt.tag_id
		WHERE 
			tt.topic_id = $1`

	if err := r.db.SelectContext(ctx, &tags, tagQuery, id); err != nil {
		return nil, fmt.Errorf("failed to get tags for topic %d: %w", id, err)
	}

	topic.Tags = tags

	return &topic, nil
}

// Implement remaining repository methods (additional methods omitted for brevity)
// The implementation would follow similar patterns as shown above

// Methods from the interface that are not implemented yet would need to be
// implemented in a similar fashion to the methods above, following the repository
// pattern and properly handling database operations, error handling, and mapping
// between database records and domain entities.
