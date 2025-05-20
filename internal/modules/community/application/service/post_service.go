// filepath: /Users/amit/workspace/sutram/internal/modules/community/application/service/post_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/community/domain/model"
	"algo/internal/modules/community/domain/repository"

	"github.com/google/uuid"
)

// Error definitions
var (
	ErrPostNotFound = errors.New("post not found")
)

// PostService handles community post business logic
type PostService struct {
	postRepo  repository.PostRepository
	topicRepo repository.TopicRepository
}

// NewPostService creates a new post service
func NewPostService(
	postRepo repository.PostRepository,
	topicRepo repository.TopicRepository,
) *PostService {
	return &PostService{
		postRepo:  postRepo,
		topicRepo: topicRepo,
	}
}

// GetPostsByTopicID gets posts for a specific topic
func (s *PostService) GetPostsByTopicID(
	ctx context.Context,
	topicID int,
	page int,
	perPage int,
) ([]*model.Post, int, error) {
	// Check if topic exists
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil || topic == nil {
		return nil, 0, ErrTopicNotFound
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	return s.postRepo.GetByTopicID(ctx, topicID, page, perPage)
}

// CreatePost creates a new post
func (s *PostService) CreatePost(
	ctx context.Context,
	topicID int,
	userID uuid.UUID,
	content string,
	parentPostID *int,
) (*model.Post, error) {
	// Check if topic exists and is open
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil || topic == nil {
		return nil, ErrTopicNotFound
	}

	if topic.IsClosed {
		return nil, ErrTopicClosed
	}

	// If parentPostID is provided, check if it exists
	if parentPostID != nil {
		parentPost, err := s.postRepo.GetByID(ctx, *parentPostID)
		if err != nil || parentPost == nil {
			return nil, ErrPostNotFound
		}

		// Ensure parent post is in the same topic
		if parentPost.TopicID != topicID {
			return nil, errors.New("parent post is not in the same topic")
		}
	}

	// Create post
	post := &model.Post{
		TopicID:      topicID,
		AuthorID:     userID,
		Content:      content,
		IsSolution:   false,
		ParentPostID: parentPostID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		EditCount:    0,
	}

	// Save to database
	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Update topic's last activity timestamp
	updates := map[string]interface{}{
		"last_activity_at": time.Now(),
	}
	if err := s.topicRepo.Update(ctx, topicID, updates); err != nil {
		// Log but don't fail the operation
		fmt.Printf("failed to update topic last activity: %v\n", err)
	}

	// Get like count
	likeCount, err := s.postRepo.GetLikeCount(ctx, post.ID)
	if err == nil {
		post.LikeCount = likeCount
	}

	return post, nil
}

// UpdatePost updates an existing post
func (s *PostService) UpdatePost(
	ctx context.Context,
	id int,
	userID uuid.UUID,
	updates map[string]interface{},
) (*model.Post, error) {
	// Check if post exists
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil || post == nil {
		return nil, ErrPostNotFound
	}

	// Check if user is the author
	isAuthor, err := s.postRepo.IsAuthor(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check authorship: %w", err)
	}

	if !isAuthor {
		return nil, ErrUnauthorized
	}

	// Add updated timestamp and increment edit count
	updates["updated_at"] = time.Now()
	updates["edit_count"] = post.EditCount + 1
	updates["edited_by"] = userID

	// Update post
	if err := s.postRepo.Update(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Get updated post
	updatedPost, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated post: %w", err)
	}

	return updatedPost, nil
}

// DeletePost deletes a post
func (s *PostService) DeletePost(ctx context.Context, id int, userID uuid.UUID) error {
	// Check if post exists
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil || post == nil {
		return ErrPostNotFound
	}

	// Check if user is the author
	isAuthor, err := s.postRepo.IsAuthor(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to check authorship: %w", err)
	}

	if !isAuthor {
		return ErrUnauthorized
	}

	// Delete post
	if err := s.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// GetPostByID gets a post by ID
func (s *PostService) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if post == nil {
		return nil, ErrPostNotFound
	}

	return post, nil
}
