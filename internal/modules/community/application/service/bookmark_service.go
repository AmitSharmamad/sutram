// filepath: /Users/amit/workspace/sutram/internal/modules/community/application/service/bookmark_service.go
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
	ErrAlreadyBookmarked = errors.New("topic already bookmarked")
	ErrNotBookmarked     = errors.New("topic not bookmarked")
)

// BookmarkService handles bookmark business logic
type BookmarkService struct {
	bookmarkRepo repository.BookmarkRepository
	topicRepo    repository.TopicRepository
}

// NewBookmarkService creates a new bookmark service
func NewBookmarkService(
	bookmarkRepo repository.BookmarkRepository,
	topicRepo repository.TopicRepository,
) *BookmarkService {
	return &BookmarkService{
		bookmarkRepo: bookmarkRepo,
		topicRepo:    topicRepo,
	}
}

// AddBookmark adds a bookmark to a topic
func (s *BookmarkService) AddBookmark(ctx context.Context, userID uuid.UUID, topicID int) error {
	// Check if topic exists
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return fmt.Errorf("failed to check topic: %w", err)
	}

	if topic == nil {
		return ErrTopicNotFound
	}

	// Check if already bookmarked
	exists, err := s.bookmarkRepo.Exists(ctx, userID, topicID)
	if err != nil {
		return fmt.Errorf("failed to check bookmark: %w", err)
	}

	if exists {
		return ErrAlreadyBookmarked
	}

	// Create bookmark
	bookmark := &model.Bookmark{
		UserID:    userID,
		TopicID:   topicID,
		CreatedAt: time.Now(),
	}

	return s.bookmarkRepo.Create(ctx, bookmark)
}

// RemoveBookmark removes a bookmark from a topic
func (s *BookmarkService) RemoveBookmark(ctx context.Context, userID uuid.UUID, topicID int) error {
	// Check if topic exists
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return fmt.Errorf("failed to check topic: %w", err)
	}

	if topic == nil {
		return ErrTopicNotFound
	}

	// Check if bookmarked
	exists, err := s.bookmarkRepo.Exists(ctx, userID, topicID)
	if err != nil {
		return fmt.Errorf("failed to check bookmark: %w", err)
	}

	if !exists {
		return ErrNotBookmarked
	}

	return s.bookmarkRepo.Delete(ctx, userID, topicID)
}

// GetUserBookmarks gets bookmarks for a user
func (s *BookmarkService) GetUserBookmarks(ctx context.Context, userID uuid.UUID, page int, perPage int) ([]*model.Topic, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	return s.bookmarkRepo.GetByUserID(ctx, userID, page, perPage)
}

// GetBookmarkCount gets the bookmark count for a topic
func (s *BookmarkService) GetBookmarkCount(ctx context.Context, topicID int) (int, error) {
	return s.bookmarkRepo.GetBookmarkCount(ctx, topicID)
}
