package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"algo/internal/modules/user/domain/model"
	"algo/internal/modules/user/domain/repository"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrSurveyInvalidEmail  = errors.New("invalid email format")
	ErrSurveyInvalidAge    = errors.New("age must be between 1 and 120")
	ErrSurveyAlreadyExists = errors.New("survey already exists for this email")
	ErrSurveyInvalidData   = errors.New("invalid survey data")
)

// SurveyService handles survey-related business logic
type SurveyService struct {
	surveyRepo repository.SurveyRepository
}

// NewSurveyService creates a new survey service
func NewSurveyService(surveyRepo repository.SurveyRepository) *SurveyService {
	return &SurveyService{
		surveyRepo: surveyRepo,
	}
}

// CreateSurvey creates a new survey entry
func (s *SurveyService) CreateSurvey(ctx context.Context, email string, natureOfWork model.NatureOfWork, tradingType model.TradingType, experienceLevel model.ExperienceLevel, workType model.WorkType, age int) (*model.Survey, error) {
	// Validate email format
	if !isValidSurveyEmail(email) {
		return nil, ErrSurveyInvalidEmail
	}

	// Validate age
	if age <= 0 || age > 120 {
		return nil, ErrSurveyInvalidAge
	}

	// Validate enum values
	if !isValidNatureOfWork(natureOfWork) || !isValidTradingType(tradingType) || !isValidExperienceLevel(experienceLevel) || !isValidWorkType(workType) {
		return nil, ErrSurveyInvalidData
	}

	// Normalize email (lowercase and trim)
	email = strings.ToLower(strings.TrimSpace(email))

	// Check if survey already exists for this email
	existingSurvey, err := s.surveyRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing survey: %w", err)
	}
	if existingSurvey != nil {
		return nil, ErrSurveyAlreadyExists
	}

	// Create new survey entry
	survey := model.NewSurvey(email, natureOfWork, tradingType, experienceLevel, workType, age)

	// Save to database
	err = s.surveyRepo.Create(ctx, survey)
	if err != nil {
		return nil, fmt.Errorf("error creating survey: %w", err)
	}

	return survey, nil
}

// GetSurveyByEmail retrieves a survey by email
func (s *SurveyService) GetSurveyByEmail(ctx context.Context, email string) (*model.Survey, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	return s.surveyRepo.FindByEmail(ctx, email)
}

// GetSurveyByID retrieves a survey by ID
func (s *SurveyService) GetSurveyByID(ctx context.Context, id uuid.UUID) (*model.Survey, error) {
	return s.surveyRepo.FindByID(ctx, id)
}

// GetAllSurveys returns all survey entries (admin only)
func (s *SurveyService) GetAllSurveys(ctx context.Context) ([]model.Survey, error) {
	return s.surveyRepo.FindAll(ctx)
}

// GetSurveyCount returns the total number of survey entries
func (s *SurveyService) GetSurveyCount(ctx context.Context) (int64, error) {
	return s.surveyRepo.GetCount(ctx)
}

// isValidSurveyEmail validates email format using regex
func isValidSurveyEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isValidNatureOfWork validates nature of work enum
func isValidNatureOfWork(natureOfWork model.NatureOfWork) bool {
	switch natureOfWork {
	case model.NatureOfWorkIndividual, model.NatureOfWorkCompany, model.NatureOfWorkProprietaryTradingFirm, model.NatureOfWorkAssetManagement, model.NatureOfWorkFamilyOffice, model.NatureOfWorkOther:
		return true
	default:
		return false
	}
}

// isValidTradingType validates trading type enum
func isValidTradingType(tradingType model.TradingType) bool {
	switch tradingType {
	case model.TradingTypeEquity, model.TradingTypeFNO, model.TradingTypeBoth:
		return true
	default:
		return false
	}
}

// isValidExperienceLevel validates experience level enum
func isValidExperienceLevel(experienceLevel model.ExperienceLevel) bool {
	switch experienceLevel {
	case model.ExperienceLevelBeginner, model.ExperienceLevelIntermediate, model.ExperienceLevelPro:
		return true
	default:
		return false
	}
}

// isValidWorkType validates work type enum
func isValidWorkType(workType model.WorkType) bool {
	switch workType {
	case model.WorkTypePersonal, model.WorkTypeFirm:
		return true
	default:
		return false
	}
}
