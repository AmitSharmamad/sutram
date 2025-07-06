package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"algo/internal/modules/user/domain/model"
	"algo/internal/modules/user/domain/repository"
)

// Common errors
var (
	ErrInvalidEmail           = errors.New("invalid email format")
	ErrEmailAlreadyInWaitlist = errors.New("email already in waitlist")
)

// WaitlistService handles waitlist-related business logic
type WaitlistService struct {
	waitlistRepo repository.WaitlistRepository
}

// NewWaitlistService creates a new waitlist service
func NewWaitlistService(waitlistRepo repository.WaitlistRepository) *WaitlistService {
	return &WaitlistService{
		waitlistRepo: waitlistRepo,
	}
}

// JoinWaitlist adds a new email to the waitlist
func (s *WaitlistService) JoinWaitlist(ctx context.Context, email string) (*model.Waitlist, error) {
	// Validate email format
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	// Normalize email (lowercase and trim)
	email = strings.ToLower(strings.TrimSpace(email))

	// Check if email already exists in waitlist
	existingEntry, err := s.waitlistRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing waitlist entry: %w", err)
	}
	if existingEntry != nil {
		return nil, ErrEmailAlreadyInWaitlist
	}

	// Create new waitlist entry
	waitlist := model.NewWaitlist(email)

	// Save to database
	err = s.waitlistRepo.Create(ctx, waitlist)
	if err != nil {
		return nil, fmt.Errorf("error creating waitlist entry: %w", err)
	}

	return waitlist, nil
}

// GetWaitlistCount returns the total number of waitlist entries
func (s *WaitlistService) GetWaitlistCount(ctx context.Context) (int64, error) {
	return s.waitlistRepo.GetCount(ctx)
}

// GetAllWaitlistEntries returns all waitlist entries (admin only)
func (s *WaitlistService) GetAllWaitlistEntries(ctx context.Context) ([]model.Waitlist, error) {
	return s.waitlistRepo.FindAll(ctx)
}

// isValidEmail validates email format using regex
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
