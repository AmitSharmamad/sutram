// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/repository/topic_repository.go
package repository

import (
	"context"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// TopicRepository defines the operations for topic persistence
type TopicRepository interface {
	// Create a new topic
	Create(ctx context.Context, topic *model.Topic) error

	// Get a topic by ID
	GetByID(ctx context.Context, id int) (*model.Topic, error)

	// Get topics with optional filtering
	GetTopics(ctx context.Context, categorySlug, tagSlug, sort, search string, page, perPage int) ([]*model.Topic, int, error)

	// Get topics created by a specific user
	GetTopicsByUserID(ctx context.Context, userID uuid.UUID, page, perPage int) ([]*model.Topic, int, error)

	// Update a topic
	Update(ctx context.Context, id int, updates map[string]interface{}) error

	// Delete a topic
	Delete(ctx context.Context, id int) error

	// Increment view count
	IncrementViews(ctx context.Context, id int) error

	// Get topic stats
	GetTopicStats(ctx context.Context, id int) (*model.TopicStats, error)

	// Add tags to a topic
	AddTags(ctx context.Context, topicID int, tagIDs []int) error

	// Remove tags from a topic
	RemoveTags(ctx context.Context, topicID int, tagIDs []int) error

	// Set tags for a topic (replaces existing tags)
	SetTags(ctx context.Context, topicID int, tagIDs []int) error

	// Get tags for a topic
	GetTags(ctx context.Context, topicID int) ([]*model.Tag, error)

	// Check if a user is the author of a topic
	IsAuthor(ctx context.Context, topicID int, userID uuid.UUID) (bool, error)
}
