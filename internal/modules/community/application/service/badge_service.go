// filepath: /Users/amit/workspace/sutram/internal/modules/community/application/service/badge_service.go
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
	ErrBadgeNotFound  = errors.New("badge not found")
	ErrAlreadyAwarded = errors.New("badge already awarded to user")
)

// BadgeService handles badge business logic
type BadgeService struct {
	badgeRepo          repository.BadgeRepository
	userReputationRepo repository.UserReputationRepository
}

// NewBadgeService creates a new badge service
func NewBadgeService(
	badgeRepo repository.BadgeRepository,
	userReputationRepo repository.UserReputationRepository,
) *BadgeService {
	return &BadgeService{
		badgeRepo:          badgeRepo,
		userReputationRepo: userReputationRepo,
	}
}

// GetAllBadges gets all available badges
func (s *BadgeService) GetAllBadges(ctx context.Context) ([]*model.Badge, error) {
	return s.badgeRepo.GetAll(ctx)
}

// GetBadgeByID gets a badge by ID
func (s *BadgeService) GetBadgeByID(ctx context.Context, id int) (*model.Badge, error) {
	badge, err := s.badgeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get badge: %w", err)
	}

	if badge == nil {
		return nil, ErrBadgeNotFound
	}

	return badge, nil
}

// GetUserBadges gets badges earned by a specific user
func (s *BadgeService) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]*model.Badge, error) {
	// Check if user exists
	userProfile, err := s.userReputationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	if userProfile == nil {
		return nil, ErrUserNotFound
	}

	return s.badgeRepo.GetUserBadges(ctx, userID)
}

// AwardBadge awards a badge to a user
func (s *BadgeService) AwardBadge(ctx context.Context, userID uuid.UUID, badgeID int) error {
	// Check if badge exists
	badge, err := s.badgeRepo.GetByID(ctx, badgeID)
	if err != nil {
		return fmt.Errorf("failed to check badge: %w", err)
	}

	if badge == nil {
		return ErrBadgeNotFound
	}

	// Check if user exists
	userProfile, err := s.userReputationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user: %w", err)
	}

	if userProfile == nil {
		return ErrUserNotFound
	}

	// Check if already awarded
	hasIt, err := s.badgeRepo.HasBadge(ctx, userID, badgeID)
	if err != nil {
		return fmt.Errorf("failed to check existing award: %w", err)
	}

	if hasIt {
		return ErrAlreadyAwarded
	}

	// Award badge
	if err := s.badgeRepo.AwardBadge(ctx, userID, badgeID); err != nil {
		return fmt.Errorf("failed to award badge: %w", err)
	}

	// Update user's reputation score
	if badge.Points > 0 {
		if err := s.userReputationRepo.IncrementReputationScore(ctx, userID, badge.Points); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to update reputation: %v\n", err)
		}
	}

	return nil
}

// RevokeBadge revokes a badge from a user
func (s *BadgeService) RevokeBadge(ctx context.Context, userID uuid.UUID, badgeID int) error {
	// Check if badge exists
	badge, err := s.badgeRepo.GetByID(ctx, badgeID)
	if err != nil {
		return fmt.Errorf("failed to check badge: %w", err)
	}

	if badge == nil {
		return ErrBadgeNotFound
	}

	// Check if user exists
	userProfile, err := s.userReputationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user: %w", err)
	}

	if userProfile == nil {
		return ErrUserNotFound
	}

	// Revoke the badge
	if err := s.badgeRepo.RevokeBadge(ctx, userID, badgeID); err != nil {
		return fmt.Errorf("failed to revoke badge: %w", err)
	}

	// Update user's reputation score (deduct points)
	if badge.Points > 0 {
		if err := s.userReputationRepo.IncrementReputationScore(ctx, userID, -badge.Points); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to update reputation: %v\n", err)
		}
	}

	return nil
}
