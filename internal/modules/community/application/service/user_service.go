// filepath: /Users/amit/workspace/sutram/internal/modules/community/application/service/user_service.go
package service

import (
	"context"
	"errors"
	"fmt"

	"algo/internal/modules/community/domain/model"
	"algo/internal/modules/community/domain/repository"

	"github.com/google/uuid"
)

// Error definitions
var (
	ErrUserNotFound = errors.New("user not found")
)

// UserService handles user-related community business logic
type UserService struct {
	userReputationRepo repository.UserReputationRepository
	badgeRepo          repository.BadgeRepository
}

// NewUserService creates a new user service
func NewUserService(
	userReputationRepo repository.UserReputationRepository,
	badgeRepo repository.BadgeRepository,
) *UserService {
	return &UserService{
		userReputationRepo: userReputationRepo,
		badgeRepo:          badgeRepo,
	}
}

// GetUserProfile gets a user's community profile
func (s *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.UserProfile, error) {
	// Get user reputation
	userProfile, err := s.userReputationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	if userProfile == nil {
		return nil, ErrUserNotFound
	}

	// Get topic and post count
	topicCount, err := s.userReputationRepo.GetTopicCount(ctx, userID)
	if err == nil {
		userProfile.TopicCount = topicCount
	}

	postCount, err := s.userReputationRepo.GetPostCount(ctx, userID)
	if err == nil {
		userProfile.PostCount = postCount
	}

	// Get badges
	badges, err := s.badgeRepo.GetUserBadges(ctx, userID)
	if err == nil {
		userProfile.Badges = badges
	}

	// Update last active timestamp
	_ = s.userReputationRepo.UpdateLastActive(ctx, userID)

	return userProfile, nil
}

// UpdateUserReputationScore updates a user's reputation score
func (s *UserService) UpdateUserReputationScore(ctx context.Context, userID uuid.UUID, points int) error {
	return s.userReputationRepo.IncrementReputationScore(ctx, userID, points)
}

// IncrementUserContributionCount increments a user's contribution count
func (s *UserService) IncrementUserContributionCount(ctx context.Context, userID uuid.UUID) error {
	return s.userReputationRepo.IncrementContributionCount(ctx, userID)
}

// UpdateLastActive updates a user's last active timestamp
func (s *UserService) UpdateLastActive(ctx context.Context, userID uuid.UUID) error {
	return s.userReputationRepo.UpdateLastActive(ctx, userID)
}
