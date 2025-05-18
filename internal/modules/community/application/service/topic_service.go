// filepath: /Users/amit/workspace/sutram/internal/modules/community/application/service/topic_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/community/domain/model"
	"algo/internal/modules/community/domain/repository"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

// Error definitions
var (
	ErrTopicNotFound    = errors.New("topic not found")
	ErrUnauthorized     = errors.New("unauthorized action")
	ErrCategoryNotFound = errors.New("category not found")
	ErrTopicClosed      = errors.New("topic is closed")
	ErrInvalidRequest   = errors.New("invalid request")
)

// TopicService handles community topic business logic
type TopicService struct {
	topicRepo    repository.TopicRepository
	tagRepo      repository.TagRepository
	categoryRepo repository.CategoryRepository
}

// NewTopicService creates a new topic service
func NewTopicService(
	topicRepo repository.TopicRepository,
	tagRepo repository.TagRepository,
	categoryRepo repository.CategoryRepository,
) *TopicService {
	return &TopicService{
		topicRepo:    topicRepo,
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
	}
}

// GetTopics gets topics with optional filtering
func (s *TopicService) GetTopics(
	ctx context.Context,
	categorySlug string,
	tagSlug string,
	sort string,
	search string,
	page int,
	perPage int,
) ([]*model.Topic, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	return s.topicRepo.GetTopics(ctx, categorySlug, tagSlug, sort, search, page, perPage)
}

// GetTopicByID gets a topic by ID
func (s *TopicService) GetTopicByID(ctx context.Context, id int) (*model.Topic, error) {
	topic, err := s.topicRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	if topic == nil {
		return nil, ErrTopicNotFound
	}

	return topic, nil
}

// CreateTopic creates a new topic
func (s *TopicService) CreateTopic(
	ctx context.Context,
	userID uuid.UUID,
	title string,
	content string,
	categoryID int,
	tagIDs []int,
	isPinned bool,
) (*model.Topic, error) {
	// Check if category exists
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil || category == nil {
		return nil, ErrCategoryNotFound
	}

	// Generate slug
	slugStr := slug.Make(title)

	// Create topic
	topic := &model.Topic{
		Title:          title,
		Slug:           slugStr,
		Content:        content,
		CategoryID:     categoryID,
		AuthorID:       userID,
		IsPinned:       isPinned,
		IsClosed:       false,
		IsHidden:       false,
		ViewCount:      0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastActivityAt: time.Now(),
	}

	// Save to database
	if err := s.topicRepo.Create(ctx, topic); err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	// Add tags if provided
	if len(tagIDs) > 0 {
		if err := s.topicRepo.AddTags(ctx, topic.ID, tagIDs); err != nil {
			return nil, fmt.Errorf("failed to add tags: %w", err)
		}

		// Load tags
		tags, err := s.topicRepo.GetTags(ctx, topic.ID)
		if err == nil {
			topic.Tags = tags
		}
	}

	// Load category
	topic.Category = *category

	return topic, nil
}

// UpdateTopic updates an existing topic
func (s *TopicService) UpdateTopic(
	ctx context.Context,
	id int,
	userID uuid.UUID,
	updates map[string]interface{},
) (*model.Topic, error) {
	// Check if topic exists
	topic, err := s.topicRepo.GetByID(ctx, id)
	if err != nil || topic == nil {
		return nil, ErrTopicNotFound
	}

	// Check if user is the author
	isAuthor, err := s.topicRepo.IsAuthor(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check authorship: %w", err)
	}

	if !isAuthor {
		return nil, ErrUnauthorized
	}

	// If updating category, check if it exists
	if categoryID, ok := updates["category_id"].(int); ok {
		category, err := s.categoryRepo.GetByID(ctx, categoryID)
		if err != nil || category == nil {
			return nil, ErrCategoryNotFound
		}
	}

	// If updating title, generate new slug
	if title, ok := updates["title"].(string); ok {
		updates["slug"] = slug.Make(title)
	}

	// Add updated timestamp
	updates["updated_at"] = time.Now()

	// Update topic
	if err := s.topicRepo.Update(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update topic: %w", err)
	}

	// Update tags if provided
	if tagIDs, ok := updates["tag_ids"].([]int); ok {
		if err := s.topicRepo.SetTags(ctx, id, tagIDs); err != nil {
			return nil, fmt.Errorf("failed to update tags: %w", err)
		}
	}

	// Get updated topic
	updatedTopic, err := s.topicRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated topic: %w", err)
	}

	return updatedTopic, nil
}

// DeleteTopic deletes a topic
func (s *TopicService) DeleteTopic(ctx context.Context, id int, userID uuid.UUID) error {
	// Check if topic exists
	topic, err := s.topicRepo.GetByID(ctx, id)
	if err != nil || topic == nil {
		return ErrTopicNotFound
	}

	// Check if user is the author
	isAuthor, err := s.topicRepo.IsAuthor(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to check authorship: %w", err)
	}

	if !isAuthor {
		return ErrUnauthorized
	}

	// Delete topic
	if err := s.topicRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	return nil
}

// IncrementViews increments the view count for a topic
func (s *TopicService) IncrementViews(ctx context.Context, id int) error {
	return s.topicRepo.IncrementViews(ctx, id)
}

// GetUserTopics gets topics created by a specific user
func (s *TopicService) GetUserTopics(
	ctx context.Context,
	userID uuid.UUID,
	page int,
	perPage int,
) ([]*model.Topic, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	return s.topicRepo.GetTopicsByUserID(ctx, userID, page, perPage)
}
