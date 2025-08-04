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
	surveyRepo   repository.SurveyRepository
}

// NewWaitlistService creates a new waitlist service
func NewWaitlistService(waitlistRepo repository.WaitlistRepository, surveyRepo repository.SurveyRepository) *WaitlistService {
	return &WaitlistService{
		waitlistRepo: waitlistRepo,
		surveyRepo:   surveyRepo,
	}
}

// JoinWaitlist adds a new email to the waitlist with survey data
func (s *WaitlistService) JoinWaitlist(ctx context.Context, email string, natureOfWork model.NatureOfWork, tradingType model.TradingType, experienceLevel model.ExperienceLevel, workType model.WorkType, age int) (*model.Waitlist, error) {
	// Validate email format
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	// Validate age
	if age <= 0 || age > 120 {
		return nil, errors.New("age must be between 1 and 120")
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

	// Create survey entry first
	survey := model.NewSurvey(email, natureOfWork, tradingType, experienceLevel, workType, age)
	err = s.surveyRepo.Create(ctx, survey)
	if err != nil {
		return nil, fmt.Errorf("error creating survey entry: %w", err)
	}

	// Create waitlist entry with survey reference
	waitlist := model.NewWaitlist(email, &survey.ID)

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
