// filepath: /Users/amit/workspace/sutram/internal/modules/community/application/service/like_service.go
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
	ErrContentNotFound = errors.New("content not found")
	ErrAlreadyLiked    = errors.New("content already liked")
	ErrNotLiked        = errors.New("content not liked")
)

// LikeService handles like business logic
type LikeService struct {
	likeRepo  repository.LikeRepository
	topicRepo repository.TopicRepository
	postRepo  repository.PostRepository
}

// NewLikeService creates a new like service
func NewLikeService(
	likeRepo repository.LikeRepository,
	topicRepo repository.TopicRepository,
	postRepo repository.PostRepository,
) *LikeService {
	return &LikeService{
		likeRepo:  likeRepo,
		topicRepo: topicRepo,
		postRepo:  postRepo,
	}
}

// Like adds a like to a topic or post
func (s *LikeService) Like(ctx context.Context, userID uuid.UUID, likeableType string, likeableID int) error {
	// Check if content exists
	exists, err := s.contentExists(ctx, likeableType, likeableID)
	if err != nil {
		return err
	}

	if !exists {
		return ErrContentNotFound
	}

	// Check if already liked
	alreadyLiked, err := s.likeRepo.Exists(ctx, userID, likeableType, likeableID)
	if err != nil {
		return fmt.Errorf("failed to check existing like: %w", err)
	}

	if alreadyLiked {
		return ErrAlreadyLiked
	}

	// Create like
	like := &model.Like{
		UserID:       userID,
		LikeableType: likeableType,
		LikeableID:   likeableID,
		CreatedAt:    time.Now(),
	}

	return s.likeRepo.Create(ctx, like)
}

// Unlike removes a like from a topic or post
func (s *LikeService) Unlike(ctx context.Context, userID uuid.UUID, likeableType string, likeableID int) error {
	// Check if content exists
	exists, err := s.contentExists(ctx, likeableType, likeableID)
	if err != nil {
		return err
	}

	if !exists {
		return ErrContentNotFound
	}

	// Check if liked
	liked, err := s.likeRepo.Exists(ctx, userID, likeableType, likeableID)
	if err != nil {
		return fmt.Errorf("failed to check existing like: %w", err)
	}

	if !liked {
		return ErrNotLiked
	}

	return s.likeRepo.Delete(ctx, userID, likeableType, likeableID)
}

// GetLikeCount gets the like count for a topic or post
func (s *LikeService) GetLikeCount(ctx context.Context, likeableType string, likeableID int) (int, error) {
	return s.likeRepo.GetLikeCount(ctx, likeableType, likeableID)
}

// contentExists checks if the content exists
func (s *LikeService) contentExists(ctx context.Context, likeableType string, likeableID int) (bool, error) {
	if likeableType == "topic" {
		topic, err := s.topicRepo.GetByID(ctx, likeableID)
		if err != nil {
			return false, fmt.Errorf("failed to check topic: %w", err)
		}
		return topic != nil, nil
	} else if likeableType == "post" {
		post, err := s.postRepo.GetByID(ctx, likeableID)
		if err != nil {
			return false, fmt.Errorf("failed to check post: %w", err)
		}
		return post != nil, nil
	}

	return false, errors.New("invalid likeable type")
}
