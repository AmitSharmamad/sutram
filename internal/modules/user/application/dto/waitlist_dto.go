package dto

import (
	"time"

	"algo/internal/modules/user/domain/model"

	"github.com/google/uuid"
)

// JoinWaitlistRequest represents the request to join waitlist with survey data
type JoinWaitlistRequest struct {
	Email           string `json:"email" binding:"required,email"`
	NatureOfWork    string `json:"nature_of_work" binding:"required"`
	TradingType     string `json:"trading_type" binding:"required"`
	ExperienceLevel string `json:"experience_level" binding:"required"`
	WorkType        string `json:"work_type" binding:"required"`
	Age             int    `json:"age" binding:"required,min=1,max=120"`
}

// WaitlistResponse represents the response for waitlist operations
type WaitlistResponse struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	SurveyID  *uuid.UUID `json:"survey_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// WaitlistCountResponse represents the response for waitlist count
type WaitlistCountResponse struct {
	Count int64 `json:"count"`
}

// ToWaitlistResponse converts a waitlist model to response DTO
func ToWaitlistResponse(waitlist *model.Waitlist) *WaitlistResponse {
	return &WaitlistResponse{
		ID:        waitlist.ID,
		Email:     waitlist.Email,
		SurveyID:  waitlist.SurveyID,
		CreatedAt: waitlist.CreatedAt,
		UpdatedAt: waitlist.UpdatedAt,
	}
}

// ToWaitlistResponses converts a slice of waitlist models to response DTOs
func ToWaitlistResponses(waitlists []model.Waitlist) []WaitlistResponse {
	responses := make([]WaitlistResponse, len(waitlists))
	for i, waitlist := range waitlists {
		responses[i] = *ToWaitlistResponse(&waitlist)
	}
	return responses
}
