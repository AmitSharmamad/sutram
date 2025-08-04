package dto

import (
	"time"

	"algo/internal/modules/user/domain/model"

	"github.com/google/uuid"
)

// CreateSurveyRequest represents the request to create a survey
type CreateSurveyRequest struct {
	Email           string `json:"email" binding:"required,email"`
	NatureOfWork    string `json:"nature_of_work" binding:"required"`
	TradingType     string `json:"trading_type" binding:"required"`
	ExperienceLevel string `json:"experience_level" binding:"required"`
	WorkType        string `json:"work_type" binding:"required"`
	Age             int    `json:"age" binding:"required,min=1,max=120"`
}

// SurveyResponse represents the response for survey operations
type SurveyResponse struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	NatureOfWork    string    `json:"nature_of_work"`
	TradingType     string    `json:"trading_type"`
	ExperienceLevel string    `json:"experience_level"`
	WorkType        string    `json:"work_type"`
	Age             int       `json:"age"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SurveyCountResponse represents the response for survey count
type SurveyCountResponse struct {
	Count int64 `json:"count"`
}

// ToSurveyResponse converts a survey model to response DTO
func ToSurveyResponse(survey *model.Survey) *SurveyResponse {
	return &SurveyResponse{
		ID:              survey.ID,
		Email:           survey.Email,
		NatureOfWork:    string(survey.NatureOfWork),
		TradingType:     string(survey.TradingType),
		ExperienceLevel: string(survey.ExperienceLevel),
		WorkType:        string(survey.WorkType),
		Age:             survey.Age,
		CreatedAt:       survey.CreatedAt,
		UpdatedAt:       survey.UpdatedAt,
	}
}

// ToSurveyResponses converts a slice of survey models to response DTOs
func ToSurveyResponses(surveys []model.Survey) []SurveyResponse {
	responses := make([]SurveyResponse, len(surveys))
	for i, survey := range surveys {
		responses[i] = *ToSurveyResponse(&survey)
	}
	return responses
}
