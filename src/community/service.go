package community

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"

	"algo/src/shared"
)

var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput is returned when the input is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized is returned when the user is not authorized to perform an action
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when the user is forbidden from performing an action
	ErrForbidden = errors.New("forbidden")

	// ErrConflict is returned when there's a conflict with an existing resource
	ErrConflict = errors.New("conflict")
)

// Service defines the business logic for community features
type Service interface {
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
	CreateTopic(ctx context.Context, userID uuid.UUID, topic *shared.Topic, tagIDs []int) error
	UpdateTopic(ctx context.Context, userID uuid.UUID, topic *shared.Topic, tagIDs []int) error
	DeleteTopic(ctx context.Context, userID uuid.UUID, id int) error
	PinTopic(ctx context.Context, userID uuid.UUID, id int, pin bool) error
	CloseTopic(ctx context.Context, userID uuid.UUID, id int, close bool) error

	// Posts
	GetPostsByTopicID(ctx context.Context, topicID int, filter PostFilter) (*shared.PaginatedPosts, error)
	GetPostByID(ctx context.Context, id int) (*shared.Post, error)
	CreatePost(ctx context.Context, userID uuid.UUID, post *shared.Post) error
	UpdatePost(ctx context.Context, userID uuid.UUID, post *shared.Post) error
	DeletePost(ctx context.Context, userID uuid.UUID, id int) error
	MarkPostAsSolution(ctx context.Context, userID uuid.UUID, postID int) error

	// Tags
	GetTags(ctx context.Context) ([]shared.Tag, error)
	CreateTag(ctx context.Context, userID uuid.UUID, tag *shared.Tag) error
	UpdateTag(ctx context.Context, userID uuid.UUID, tag *shared.Tag) error
	DeleteTag(ctx context.Context, userID uuid.UUID, id int) error

	// Likes
	ToggleLike(ctx context.Context, userID uuid.UUID, entityType string, entityID int) (bool, error)
	GetLikesByEntity(ctx context.Context, entityType string, entityID int) ([]shared.Like, error)

	// Bookmarks
	GetBookmarksByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.Topic, int, error)
	ToggleBookmark(ctx context.Context, userID uuid.UUID, topicID int) (bool, error)

	// User Follows
	GetFollowersByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.User, int, error)
	GetFollowingByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.User, int, error)
	ToggleUserFollow(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)

	// Category Follows
	GetCategoryFollowsByUser(ctx context.Context, userID uuid.UUID) ([]shared.Category, error)
	ToggleCategoryFollow(ctx context.Context, userID uuid.UUID, categoryID int) (bool, error)

	// Topic Follows
	GetTopicFollowsByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.Topic, int, error)
	ToggleTopicFollow(ctx context.Context, userID uuid.UUID, topicID int) (bool, error)

	// Notifications
	GetNotificationsByUser(ctx context.Context, userID uuid.UUID, page, limit int, unreadOnly bool) ([]shared.Notification, int, error)
	MarkNotificationAsRead(ctx context.Context, userID uuid.UUID, id int) error
	MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error
	DeleteNotification(ctx context.Context, userID uuid.UUID, id int) error

	// Badges
	GetBadges(ctx context.Context) ([]shared.Badge, error)
	GetUserBadges(ctx context.Context, userID uuid.UUID) ([]shared.UserBadge, error)

	// Reports
	CreateReport(ctx context.Context, userID uuid.UUID, report *shared.Report) error
	RespondToReport(ctx context.Context, userID uuid.UUID, reportID int, status string) error

	// Polls
	CreatePoll(ctx context.Context, userID uuid.UUID, poll *shared.Poll, options []shared.PollOption) error
	GetPollsByTopic(ctx context.Context, topicID int) ([]shared.Poll, error)
	VoteInPoll(ctx context.Context, userID uuid.UUID, pollID int, optionID int) error
	GetPollResults(ctx context.Context, pollID int) ([]shared.PollOption, error)

	// User Reputation
	GetUserReputation(ctx context.Context, userID uuid.UUID) (*shared.UserReputation, error)
	GetTopUsers(ctx context.Context, limit int) ([]shared.UserReputation, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new instance of Service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetCategories retrieves all active categories
func (s *service) GetCategories(ctx context.Context) ([]shared.Category, error) {
	return s.repo.GetCategories(ctx)
}

// GetCategoryByID retrieves a category by ID
func (s *service) GetCategoryByID(ctx context.Context, id int) (*shared.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by ID: %w", err)
	}

	if category == nil {
		return nil, ErrNotFound
	}

	return category, nil
}

// GetCategoryBySlug retrieves a category by slug
func (s *service) GetCategoryBySlug(ctx context.Context, slug string) (*shared.Category, error) {
	category, err := s.repo.GetCategoryBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	if category == nil {
		return nil, ErrNotFound
	}

	return category, nil
}

// CreateCategory creates a new category
func (s *service) CreateCategory(ctx context.Context, category *shared.Category) error {
	// Validate category
	if category.Name == "" {
		return fmt.Errorf("%w: category name is required", ErrInvalidInput)
	}

	// Generate slug if not provided
	if category.Slug == "" {
		category.Slug = slug.Make(category.Name)
	}

	// Check for duplicate slug
	existing, err := s.repo.GetCategoryBySlug(ctx, category.Slug)
	if err != nil {
		return fmt.Errorf("failed to check for duplicate slug: %w", err)
	}

	if existing != nil {
		return fmt.Errorf("%w: category with slug '%s' already exists", ErrConflict, category.Slug)
	}

	// Set defaults
	if category.DisplayOrder == 0 {
		category.DisplayOrder = 1
	}
	if !category.IsActive {
		category.IsActive = true
	}

	// Create category
	return s.repo.CreateCategory(ctx, category)
}

// UpdateCategory updates an existing category
func (s *service) UpdateCategory(ctx context.Context, category *shared.Category) error {
	// Validate category
	if category.ID <= 0 {
		return fmt.Errorf("%w: category ID is required", ErrInvalidInput)
	}
	if category.Name == "" {
		return fmt.Errorf("%w: category name is required", ErrInvalidInput)
	}

	// Check if category exists
	existing, err := s.repo.GetCategoryByID(ctx, category.ID)
	if err != nil {
		return fmt.Errorf("failed to check if category exists: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("%w: category with ID %d", ErrNotFound, category.ID)
	}

	// Generate slug if not provided
	if category.Slug == "" {
		category.Slug = slug.Make(category.Name)
	}

	// Check for duplicate slug (only if slug has changed)
	if category.Slug != existing.Slug {
		duplicateSlug, err := s.repo.GetCategoryBySlug(ctx, category.Slug)
		if err != nil {
			return fmt.Errorf("failed to check for duplicate slug: %w", err)
		}

		if duplicateSlug != nil && duplicateSlug.ID != category.ID {
			return fmt.Errorf("%w: category with slug '%s' already exists", ErrConflict, category.Slug)
		}
	}

	// Update category
	return s.repo.UpdateCategory(ctx, category)
}

// DeleteCategory deletes a category
func (s *service) DeleteCategory(ctx context.Context, id int) error {
	// Check if category exists
	existing, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if category exists: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("%w: category with ID %d", ErrNotFound, id)
	}

	// Delete category
	return s.repo.DeleteCategory(ctx, id)
}

// GetTopics retrieves topics with pagination and filtering
func (s *service) GetTopics(ctx context.Context, filter TopicFilter) (*shared.PaginatedTopics, error) {
	// Validate filter
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}

	// Map filter to repository filter
	repoFilter := TopicFilter{
		CategoryID:    filter.CategoryID,
		AuthorID:      filter.AuthorID,
		Tags:          filter.Tags,
		SearchQuery:   filter.SearchQuery,
		SortBy:        filter.SortBy,
		Page:          filter.Page,
		Limit:         filter.Limit,
		IncludeHidden: false, // Default to not including hidden topics
	}

	return s.repo.GetTopics(ctx, repoFilter)
}

// GetTopicByID retrieves a topic by ID
func (s *service) GetTopicByID(ctx context.Context, id int) (*shared.Topic, error) {
	topic, err := s.repo.GetTopicByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic by ID: %w", err)
	}

	if topic == nil {
		return nil, ErrNotFound
	}

	// Increment view count
	if err := s.repo.IncrementTopicViewCount(ctx, id); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to increment view count for topic %d: %s\n", id, err)
	}

	return topic, nil
}

// GetTopicBySlug retrieves a topic by slug
func (s *service) GetTopicBySlug(ctx context.Context, slug string) (*shared.Topic, error) {
	topic, err := s.repo.GetTopicBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic by slug: %w", err)
	}

	if topic == nil {
		return nil, ErrNotFound
	}

	// Increment view count
	if err := s.repo.IncrementTopicViewCount(ctx, topic.ID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to increment view count for topic %d: %s\n", topic.ID, err)
	}

	return topic, nil
}

// CreateTopic creates a new topic
func (s *service) CreateTopic(ctx context.Context, userID uuid.UUID, topic *shared.Topic, tagIDs []int) error {
	// Set the author ID
	topic.AuthorID = userID

	// Validate topic
	if topic.Title == "" {
		return fmt.Errorf("%w: topic title is required", ErrInvalidInput)
	}
	if topic.Content == "" {
		return fmt.Errorf("%w: topic content is required", ErrInvalidInput)
	}
	if topic.CategoryID <= 0 {
		return fmt.Errorf("%w: category ID is required", ErrInvalidInput)
	}

	// Check if category exists
	category, err := s.repo.GetCategoryByID(ctx, topic.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to check if category exists: %w", err)
	}

	if category == nil {
		return fmt.Errorf("%w: category with ID %d", ErrNotFound, topic.CategoryID)
	}

	// Generate slug if not provided
	if topic.Slug == "" {
		topic.Slug = slug.Make(topic.Title)
	}

	// Check for duplicate slug
	existing, err := s.repo.GetTopicBySlug(ctx, topic.Slug)
	if err != nil {
		return fmt.Errorf("failed to check for duplicate slug: %w", err)
	}

	if existing != nil {
		// Append a unique identifier to make the slug unique
		topic.Slug = fmt.Sprintf("%s-%d", topic.Slug, time.Now().Unix())
	}

	// Set default values
	topic.CreatedAt = time.Now()
	topic.UpdatedAt = topic.CreatedAt
	topic.LastActivityAt = topic.CreatedAt

	// Create topic
	if err := s.repo.CreateTopic(ctx, topic); err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	// Assign tags if provided
	if len(tagIDs) > 0 {
		if err := s.repo.AssignTagsToTopic(ctx, topic.ID, tagIDs); err != nil {
			return fmt.Errorf("failed to assign tags to topic: %w", err)
		}
	}

	// Log activity
	topicIDStr := fmt.Sprintf("%d", topic.ID)
	topicIDUUID, _ := uuid.Parse(topicIDStr)
	activityLog := &shared.ActivityLog{
		UserID:        userID,
		ActivityType:  "created_topic",
		Description:   fmt.Sprintf("Created topic: %s", topic.Title),
		ReferenceType: "topic",
		ReferenceID:   &topicIDUUID,
	}

	if err := s.repo.LogActivity(ctx, activityLog); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to log activity for topic creation: %s\n", err)
	}

	return nil
}

// UpdateTopic updates an existing topic
func (s *service) UpdateTopic(ctx context.Context, userID uuid.UUID, topic *shared.Topic, tagIDs []int) error {
	// Validate topic
	if topic.ID <= 0 {
		return fmt.Errorf("%w: topic ID is required", ErrInvalidInput)
	}
	if topic.Title == "" {
		return fmt.Errorf("%w: topic title is required", ErrInvalidInput)
	}
	if topic.Content == "" {
		return fmt.Errorf("%w: topic content is required", ErrInvalidInput)
	}
	if topic.CategoryID <= 0 {
		return fmt.Errorf("%w: category ID is required", ErrInvalidInput)
	}

	// Check if topic exists
	existing, err := s.repo.GetTopicByID(ctx, topic.ID)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("%w: topic with ID %d", ErrNotFound, topic.ID)
	}

	// Check if user is the author
	if existing.AuthorID != userID {
		return ErrUnauthorized
	}

	// Update timestamp
	topic.UpdatedAt = time.Now()
	topic.LastActivityAt = topic.UpdatedAt

	// Update topic
	if err := s.repo.UpdateTopic(ctx, topic); err != nil {
		return fmt.Errorf("failed to update topic: %w", err)
	}

	// Update tags if provided
	if tagIDs != nil {
		// First, remove all existing tags
		// Then assign new tags
		// This is a simplification - a more efficient approach would be to compare and only make necessary changes
		if err := s.repo.AssignTagsToTopic(ctx, topic.ID, tagIDs); err != nil {
			return fmt.Errorf("failed to update tags for topic: %w", err)
		}
	}

	return nil
}

// DeleteTopic deletes a topic
func (s *service) DeleteTopic(ctx context.Context, userID uuid.UUID, id int) error {
	// Check if topic exists
	topic, err := s.repo.GetTopicByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if topic == nil {
		return fmt.Errorf("%w: topic with ID %d", ErrNotFound, id)
	}

	// Check if user is the author
	if topic.AuthorID != userID {
		return ErrUnauthorized
	}

	// Delete topic
	return s.repo.DeleteTopic(ctx, id)
}

// PinTopic pins or unpins a topic
func (s *service) PinTopic(ctx context.Context, userID uuid.UUID, id int, pin bool) error {
	// Check if topic exists
	topic, err := s.repo.GetTopicByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if topic == nil {
		return fmt.Errorf("%w: topic with ID %d", ErrNotFound, id)
	}

	// Only allow author to pin/unpin
	if topic.AuthorID != userID {
		return ErrUnauthorized
	}

	// Update pin status
	topic.IsPinned = pin
	topic.UpdatedAt = time.Now()

	return s.repo.UpdateTopic(ctx, topic)
}

// CloseTopic closes or reopens a topic
func (s *service) CloseTopic(ctx context.Context, userID uuid.UUID, id int, close bool) error {
	// Check if topic exists
	topic, err := s.repo.GetTopicByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if topic == nil {
		return fmt.Errorf("%w: topic with ID %d", ErrNotFound, id)
	}

	// Only allow author to close/reopen
	if topic.AuthorID != userID {
		return ErrUnauthorized
	}

	// Update closed status
	topic.IsClosed = close
	topic.UpdatedAt = time.Now()

	return s.repo.UpdateTopic(ctx, topic)
}

// GetPostsByTopicID retrieves posts for a topic
func (s *service) GetPostsByTopicID(ctx context.Context, topicID int, filter PostFilter) (*shared.PaginatedPosts, error) {
	// Validate topic exists
	topic, err := s.repo.GetTopicByID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if topic == nil {
		return nil, fmt.Errorf("%w: topic with ID %d", ErrNotFound, topicID)
	}

	// Validate filter
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}

	return s.repo.GetPostsByTopicID(ctx, topicID, filter)
}

// GetPostByID retrieves a post by ID
func (s *service) GetPostByID(ctx context.Context, id int) (*shared.Post, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}

	if post == nil {
		return nil, ErrNotFound
	}

	return post, nil
}

// CreatePost creates a new post
func (s *service) CreatePost(ctx context.Context, userID uuid.UUID, post *shared.Post) error {
	// Set the author ID
	post.AuthorID = userID

	// Validate post
	if post.Content == "" {
		return fmt.Errorf("%w: post content is required", ErrInvalidInput)
	}
	if post.TopicID <= 0 {
		return fmt.Errorf("%w: topic ID is required", ErrInvalidInput)
	}

	// Check if topic exists and is not closed
	topic, err := s.repo.GetTopicByID(ctx, post.TopicID)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if topic == nil {
		return fmt.Errorf("%w: topic with ID %d", ErrNotFound, post.TopicID)
	}

	if topic.IsClosed {
		return fmt.Errorf("%w: topic is closed", ErrForbidden)
	}

	// If this is a reply to another post, check if that post exists
	if post.ParentPostID != nil && *post.ParentPostID > 0 {
		parentPost, err := s.repo.GetPostByID(ctx, *post.ParentPostID)
		if err != nil {
			return fmt.Errorf("failed to check if parent post exists: %w", err)
		}

		if parentPost == nil {
			return fmt.Errorf("%w: parent post with ID %d", ErrNotFound, *post.ParentPostID)
		}

		// Make sure parent post is in the same topic
		if parentPost.TopicID != post.TopicID {
			return fmt.Errorf("%w: parent post is not in the same topic", ErrInvalidInput)
		}
	}

	// Set timestamps
	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt

	// Create post
	if err := s.repo.CreatePost(ctx, post); err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	// Update topic's last activity
	topic.LastActivityAt = post.CreatedAt
	if err := s.repo.UpdateTopic(ctx, topic); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to update topic's last activity: %s\n", err)
	}

	// Create notification for topic author (if not the same as post author)
	if topic.AuthorID != userID {
		notification := &shared.Notification{
			UserID:           topic.AuthorID,
			NotificationType: "new_reply",
			ActorID:          &userID,
			ReferenceType:    "post",
			ReferenceID:      post.ID,
			Message:          fmt.Sprintf("New reply to your topic '%s'", topic.Title),
			IsRead:           false,
			CreatedAt:        post.CreatedAt,
		}

		if err := s.repo.CreateNotification(ctx, notification); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to create notification: %s\n", err)
		}
	}

	// Create notification for parent post author (if applicable and not the same as post author)
	if post.ParentPostID != nil && *post.ParentPostID > 0 {
		parentPost, _ := s.repo.GetPostByID(ctx, *post.ParentPostID)
		if parentPost != nil && parentPost.AuthorID != userID && parentPost.AuthorID != topic.AuthorID {
			notification := &shared.Notification{
				UserID:           parentPost.AuthorID,
				NotificationType: "reply_to_post",
				ActorID:          &userID,
				ReferenceType:    "post",
				ReferenceID:      post.ID,
				Message:          fmt.Sprintf("New reply to your post in '%s'", topic.Title),
				IsRead:           false,
				CreatedAt:        post.CreatedAt,
			}

			if err := s.repo.CreateNotification(ctx, notification); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Failed to create notification: %s\n", err)
			}
		}
	}

	return nil
}

// UpdatePost updates an existing post
func (s *service) UpdatePost(ctx context.Context, userID uuid.UUID, post *shared.Post) error {
	// Validate post
	if post.ID <= 0 {
		return fmt.Errorf("%w: post ID is required", ErrInvalidInput)
	}
	if post.Content == "" {
		return fmt.Errorf("%w: post content is required", ErrInvalidInput)
	}

	// Check if post exists
	existing, err := s.repo.GetPostByID(ctx, post.ID)
	if err != nil {
		return fmt.Errorf("failed to check if post exists: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("%w: post with ID %d", ErrNotFound, post.ID)
	}

	// Check if user is the author
	if existing.AuthorID != userID {
		return ErrUnauthorized
	}

	// Update post
	existing.Content = post.Content
	existing.UpdatedAt = time.Now()
	existing.EditedByID = &userID
	existing.EditCount++

	return s.repo.UpdatePost(ctx, existing)
}

// DeletePost deletes a post
func (s *service) DeletePost(ctx context.Context, userID uuid.UUID, id int) error {
	// Check if post exists
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if post exists: %w", err)
	}

	if post == nil {
		return fmt.Errorf("%w: post with ID %d", ErrNotFound, id)
	}

	// Check if user is the author
	if post.AuthorID != userID {
		return ErrUnauthorized
	}

	// Delete post
	return s.repo.DeletePost(ctx, id)
}

// MarkPostAsSolution marks a post as the solution for its topic
func (s *service) MarkPostAsSolution(ctx context.Context, userID uuid.UUID, postID int) error {
	// Check if post exists
	post, err := s.repo.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to check if post exists: %w", err)
	}

	if post == nil {
		return fmt.Errorf("%w: post with ID %d", ErrNotFound, postID)
	}

	// Check if topic exists
	topic, err := s.repo.GetTopicByID(ctx, post.TopicID)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if topic == nil {
		return fmt.Errorf("%w: topic with ID %d", ErrNotFound, post.TopicID)
	}

	// Check if user is the topic author (only topic author can mark solution)
	if topic.AuthorID != userID {
		return ErrUnauthorized
	}

	// Mark post as solution
	if err := s.repo.MarkPostAsSolution(ctx, postID); err != nil {
		return fmt.Errorf("failed to mark post as solution: %w", err)
	}

	// Award reputation to post author
	if err := s.repo.UpdateUserReputation(ctx, post.AuthorID, 15); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to update user reputation: %s\n", err)
	}

	// Create notification for post author
	notification := &shared.Notification{
		UserID:           post.AuthorID,
		NotificationType: "solution_marked",
		ActorID:          &userID,
		ReferenceType:    "post",
		ReferenceID:      postID,
		Message:          fmt.Sprintf("Your post was marked as the solution for '%s'", topic.Title),
		IsRead:           false,
		CreatedAt:        time.Now(),
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to create notification: %s\n", err)
	}

	return nil
}

// Tags methods

// GetTags returns all tags
func (s *service) GetTags(ctx context.Context) ([]shared.Tag, error) {
	return s.repo.GetTags(ctx)
}

// CreateTag creates a new tag
func (s *service) CreateTag(ctx context.Context, userID uuid.UUID, tag *shared.Tag) error {
	// Validate tag
	if tag.Name == "" {
		return fmt.Errorf("%w: tag name is required", ErrInvalidInput)
	}

	// Generate slug if not provided
	if tag.Slug == "" {
		tag.Slug = slug.Make(tag.Name)
	}

	// Set defaults
	tag.IsActive = true
	tag.CreatedAt = time.Now()

	return s.repo.CreateTag(ctx, tag)
}

// UpdateTag updates an existing tag
func (s *service) UpdateTag(ctx context.Context, userID uuid.UUID, tag *shared.Tag) error {
	// Validate tag
	if tag.ID <= 0 {
		return fmt.Errorf("%w: tag ID is required", ErrInvalidInput)
	}
	if tag.Name == "" {
		return fmt.Errorf("%w: tag name is required", ErrInvalidInput)
	}

	// Check if tag exists
	existing, err := s.repo.GetTagByID(ctx, tag.ID)
	if err != nil {
		return fmt.Errorf("failed to check if tag exists: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("%w: tag with ID %d", ErrNotFound, tag.ID)
	}

	// Generate slug if not provided
	if tag.Slug == "" {
		tag.Slug = slug.Make(tag.Name)
	}

	return s.repo.UpdateTag(ctx, tag)
}

// DeleteTag deletes a tag
func (s *service) DeleteTag(ctx context.Context, userID uuid.UUID, id int) error {
	// Check if tag exists
	tag, err := s.repo.GetTagByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if tag exists: %w", err)
	}

	if tag == nil {
		return fmt.Errorf("%w: tag with ID %d", ErrNotFound, id)
	}

	return s.repo.DeleteTag(ctx, id)
}

// Likes methods

// ToggleLike toggles a like on an entity (topic or post)
func (s *service) ToggleLike(ctx context.Context, userID uuid.UUID, entityType string, entityID int) (bool, error) {
	// Validate entity type
	if entityType != "topic" && entityType != "post" {
		return false, fmt.Errorf("%w: entity type must be 'topic' or 'post'", ErrInvalidInput)
	}

	// Check if entity exists
	var exists bool
	var err error
	if entityType == "topic" {
		var topic *shared.Topic
		topic, err = s.repo.GetTopicByID(ctx, entityID)
		exists = topic != nil
	} else {
		var post *shared.Post
		post, err = s.repo.GetPostByID(ctx, entityID)
		exists = post != nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if %s exists: %w", entityType, err)
	}

	if !exists {
		return false, fmt.Errorf("%w: %s with ID %d", ErrNotFound, entityType, entityID)
	}

	// Check if user has already liked the entity
	like, err := s.repo.GetLikeByUserAndEntity(ctx, userID, entityType, entityID)
	if err != nil {
		return false, fmt.Errorf("failed to check if user has liked entity: %w", err)
	}

	if like != nil {
		// User has already liked, so unlike
		if err := s.repo.DeleteLike(ctx, like.ID); err != nil {
			return false, fmt.Errorf("failed to unlike entity: %w", err)
		}
		return false, nil
	}

	// User hasn't liked, so add like
	like = &shared.Like{
		UserID:       userID,
		LikeableType: entityType,
		LikeableID:   entityID,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateLike(ctx, like); err != nil {
		return false, fmt.Errorf("failed to like entity: %w", err)
	}

	return true, nil
}

// GetLikesByEntity returns all likes for an entity
func (s *service) GetLikesByEntity(ctx context.Context, entityType string, entityID int) ([]shared.Like, error) {
	// Validate entity type
	if entityType != "topic" && entityType != "post" {
		return nil, fmt.Errorf("%w: entity type must be 'topic' or 'post'", ErrInvalidInput)
	}

	return s.repo.GetLikesByEntityID(ctx, entityType, entityID)
}

// Bookmarks methods

// GetBookmarksByUser returns all bookmarks for a user
func (s *service) GetBookmarksByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.Topic, int, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.repo.GetBookmarksByUserID(ctx, userID, page, limit)
}

// ToggleBookmark toggles a bookmark on a topic
func (s *service) ToggleBookmark(ctx context.Context, userID uuid.UUID, topicID int) (bool, error) {
	// Check if topic exists
	topic, err := s.repo.GetTopicByID(ctx, topicID)
	if err != nil {
		return false, fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if topic == nil {
		return false, fmt.Errorf("%w: topic with ID %d", ErrNotFound, topicID)
	}

	// Check if user has already bookmarked the topic
	bookmark, err := s.repo.GetBookmarkByUserAndTopic(ctx, userID, topicID)
	if err != nil {
		return false, fmt.Errorf("failed to check if user has bookmarked topic: %w", err)
	}

	if bookmark != nil {
		// User has already bookmarked, so unbookmark
		if err := s.repo.DeleteBookmark(ctx, bookmark.ID); err != nil {
			return false, fmt.Errorf("failed to unbookmark topic: %w", err)
		}
		return false, nil
	}

	// User hasn't bookmarked, so add bookmark
	bookmark = &shared.Bookmark{
		UserID:    userID,
		TopicID:   topicID,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateBookmark(ctx, bookmark); err != nil {
		return false, fmt.Errorf("failed to bookmark topic: %w", err)
	}

	return true, nil
}

// User Follows methods

// GetFollowersByUser returns all followers of a user
func (s *service) GetFollowersByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.User, int, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.repo.GetFollowersByUserID(ctx, userID, page, limit)
}

// GetFollowingByUser returns all users followed by a user
func (s *service) GetFollowingByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.User, int, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.repo.GetFollowingByUserID(ctx, userID, page, limit)
}

// ToggleUserFollow toggles following a user
func (s *service) ToggleUserFollow(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	// User cannot follow themselves
	if followerID == followingID {
		return false, fmt.Errorf("%w: user cannot follow themselves", ErrInvalidInput)
	}

	// Check if already following
	isFollowing, err := s.repo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is following: %w", err)
	}

	if isFollowing {
		// Already following, so unfollow
		if err := s.repo.DeleteUserFollow(ctx, followerID, followingID); err != nil {
			return false, fmt.Errorf("failed to unfollow user: %w", err)
		}
		return false, nil
	}

	// Not following, so follow
	follow := &shared.UserFollow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateUserFollow(ctx, follow); err != nil {
		return false, fmt.Errorf("failed to follow user: %w", err)
	}

	return true, nil
}

// Category Follows methods

// GetCategoryFollowsByUser returns all categories followed by a user
func (s *service) GetCategoryFollowsByUser(ctx context.Context, userID uuid.UUID) ([]shared.Category, error) {
	return s.repo.GetCategoryFollowsByUserID(ctx, userID)
}

// ToggleCategoryFollow toggles following a category
func (s *service) ToggleCategoryFollow(ctx context.Context, userID uuid.UUID, categoryID int) (bool, error) {
	// Check if category exists
	category, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return false, fmt.Errorf("failed to check if category exists: %w", err)
	}

	if category == nil {
		return false, fmt.Errorf("%w: category with ID %d", ErrNotFound, categoryID)
	}

	// Check if already following
	isFollowing, err := s.repo.IsCategoryFollowed(ctx, userID, categoryID)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is following category: %w", err)
	}

	if isFollowing {
		// Already following, so unfollow
		if err := s.repo.DeleteCategoryFollow(ctx, userID, categoryID); err != nil {
			return false, fmt.Errorf("failed to unfollow category: %w", err)
		}
		return false, nil
	}

	// Not following, so follow
	follow := &shared.CategoryFollow{
		UserID:     userID,
		CategoryID: categoryID,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.CreateCategoryFollow(ctx, follow); err != nil {
		return false, fmt.Errorf("failed to follow category: %w", err)
	}

	return true, nil
}

// Topic Follows methods

// GetTopicFollowsByUser returns all topics followed by a user
func (s *service) GetTopicFollowsByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]shared.Topic, int, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.repo.GetTopicFollowsByUserID(ctx, userID, page, limit)
}

// ToggleTopicFollow toggles following a topic
func (s *service) ToggleTopicFollow(ctx context.Context, userID uuid.UUID, topicID int) (bool, error) {
	// Check if topic exists
	topic, err := s.repo.GetTopicByID(ctx, topicID)
	if err != nil {
		return false, fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if topic == nil {
		return false, fmt.Errorf("%w: topic with ID %d", ErrNotFound, topicID)
	}

	// Check if already following
	isFollowing, err := s.repo.IsTopicFollowed(ctx, userID, topicID)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is following topic: %w", err)
	}

	if isFollowing {
		// Already following, so unfollow
		if err := s.repo.DeleteTopicFollow(ctx, userID, topicID); err != nil {
			return false, fmt.Errorf("failed to unfollow topic: %w", err)
		}
		return false, nil
	}

	// Not following, so follow
	follow := &shared.TopicFollow{
		UserID:    userID,
		TopicID:   topicID,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateTopicFollow(ctx, follow); err != nil {
		return false, fmt.Errorf("failed to follow topic: %w", err)
	}

	return true, nil
}

// Notifications methods

// GetNotificationsByUser returns all notifications for a user
func (s *service) GetNotificationsByUser(ctx context.Context, userID uuid.UUID, page, limit int, unreadOnly bool) ([]shared.Notification, int, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.repo.GetNotificationsByUserID(ctx, userID, page, limit, unreadOnly)
}

// MarkNotificationAsRead marks a notification as read
func (s *service) MarkNotificationAsRead(ctx context.Context, userID uuid.UUID, id int) error {
	// Check if notification exists and belongs to user
	notification, err := s.repo.GetNotificationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if notification exists: %w", err)
	}

	if notification == nil {
		return fmt.Errorf("%w: notification with ID %d", ErrNotFound, id)
	}

	if notification.UserID != userID {
		return ErrUnauthorized
	}

	return s.repo.MarkNotificationAsRead(ctx, id)
}

// MarkAllNotificationsAsRead marks all notifications for a user as read
func (s *service) MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllNotificationsAsRead(ctx, userID)
}

// DeleteNotification deletes a notification
func (s *service) DeleteNotification(ctx context.Context, userID uuid.UUID, id int) error {
	// Check if notification exists and belongs to user
	notification, err := s.repo.GetNotificationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if notification exists: %w", err)
	}

	if notification == nil {
		return fmt.Errorf("%w: notification with ID %d", ErrNotFound, id)
	}

	if notification.UserID != userID {
		return ErrUnauthorized
	}

	return s.repo.DeleteNotification(ctx, id)
}

// Badges methods

// GetBadges returns all badges
func (s *service) GetBadges(ctx context.Context) ([]shared.Badge, error) {
	return s.repo.GetBadges(ctx)
}

// GetUserBadges returns all badges for a user
func (s *service) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]shared.UserBadge, error) {
	return s.repo.GetUserBadges(ctx, userID)
}

// Reports methods

// CreateReport creates a new report
func (s *service) CreateReport(ctx context.Context, userID uuid.UUID, report *shared.Report) error {
	// Validate report
	if report.ContentType != "topic" && report.ContentType != "post" && report.ContentType != "user" {
		return fmt.Errorf("%w: content type must be 'topic', 'post', or 'user'", ErrInvalidInput)
	}
	if report.ContentID <= 0 {
		return fmt.Errorf("%w: content ID is required", ErrInvalidInput)
	}
	if report.Reason == "" {
		return fmt.Errorf("%w: reason is required", ErrInvalidInput)
	}

	// Set reporter ID and default values
	report.ReporterID = userID
	report.Status = "pending"
	report.CreatedAt = time.Now()
	report.UpdatedAt = report.CreatedAt

	return s.repo.CreateReport(ctx, report)
}

// RespondToReport updates the status of a report
func (s *service) RespondToReport(ctx context.Context, userID uuid.UUID, reportID int, status string) error {
	// Validate status
	if status != "resolved" && status != "dismissed" {
		return fmt.Errorf("%w: status must be 'resolved' or 'dismissed'", ErrInvalidInput)
	}

	// Get report
	report, err := s.repo.GetReportByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}

	if report == nil {
		return fmt.Errorf("%w: report with ID %d", ErrNotFound, reportID)
	}

	return s.repo.UpdateReportStatus(ctx, reportID, status, userID)
}

// Polls methods

// CreatePoll creates a new poll
func (s *service) CreatePoll(ctx context.Context, userID uuid.UUID, poll *shared.Poll, options []shared.PollOption) error {
	// Validate poll
	if poll.TopicID <= 0 {
		return fmt.Errorf("%w: topic ID is required", ErrInvalidInput)
	}
	if poll.Title == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if len(options) < 2 {
		return fmt.Errorf("%w: at least two options are required", ErrInvalidInput)
	}

	// Check if topic exists
	topic, err := s.repo.GetTopicByID(ctx, poll.TopicID)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if topic == nil {
		return fmt.Errorf("%w: topic with ID %d", ErrNotFound, poll.TopicID)
	}

	// Check if user is the topic author
	if topic.AuthorID != userID {
		return ErrUnauthorized
	}

	// Set default values
	poll.CreatedByID = userID
	poll.IsClosed = false
	poll.CreatedAt = time.Now()
	poll.UpdatedAt = poll.CreatedAt

	// Set display order for options
	for i := range options {
		options[i].DisplayOrder = i + 1
	}

	return s.repo.CreatePoll(ctx, poll, options)
}

// GetPollsByTopic returns all polls for a topic
func (s *service) GetPollsByTopic(ctx context.Context, topicID int) ([]shared.Poll, error) {
	// Check if topic exists
	topic, err := s.repo.GetTopicByID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if topic == nil {
		return nil, fmt.Errorf("%w: topic with ID %d", ErrNotFound, topicID)
	}

	return s.repo.GetPollsByTopicID(ctx, topicID)
}

// VoteInPoll adds a vote to a poll option
func (s *service) VoteInPoll(ctx context.Context, userID uuid.UUID, pollID, optionID int) error {
	// Check if poll exists
	poll, err := s.repo.GetPollByID(ctx, pollID)
	if err != nil {
		return fmt.Errorf("failed to get poll: %w", err)
	}

	if poll == nil {
		return fmt.Errorf("%w: poll with ID %d", ErrNotFound, pollID)
	}

	if poll.IsClosed {
		return fmt.Errorf("%w: poll is closed", ErrForbidden)
	}

	// Check if option exists for this poll
	optionExists := false
	for _, option := range poll.Options {
		if option.ID == optionID {
			optionExists = true
			break
		}
	}

	if !optionExists {
		return fmt.Errorf("%w: option with ID %d does not exist for this poll", ErrNotFound, optionID)
	}

	// Remove any existing votes if not a multiple choice poll
	if !poll.IsMultipleChoice {
		if err := s.repo.RemoveVoteFromPoll(ctx, userID, pollID); err != nil {
			return fmt.Errorf("failed to remove existing vote: %w", err)
		}
	}

	// Add vote
	vote := &shared.PollVote{
		PollID:    pollID,
		OptionID:  optionID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	return s.repo.AddVoteToPoll(ctx, vote)
}

// GetPollResults returns the results of a poll
func (s *service) GetPollResults(ctx context.Context, pollID int) ([]shared.PollOption, error) {
	// Check if poll exists
	poll, err := s.repo.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}

	if poll == nil {
		return nil, fmt.Errorf("%w: poll with ID %d", ErrNotFound, pollID)
	}

	return s.repo.GetPollResults(ctx, pollID)
}

// User Reputation methods

// GetUserReputation returns the reputation for a user
func (s *service) GetUserReputation(ctx context.Context, userID uuid.UUID) (*shared.UserReputation, error) {
	return s.repo.GetUserReputation(ctx, userID)
}

// GetTopUsers returns the top users by reputation
func (s *service) GetTopUsers(ctx context.Context, limit int) ([]shared.UserReputation, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetTopUsers(ctx, limit)
}
