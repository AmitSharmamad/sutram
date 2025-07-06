package dto

import (
	"time"

	"algo/internal/modules/user/domain/model"
)

// JoinWaitlistRequest represents the request to join waitlist
type JoinWaitlistRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// WaitlistResponse represents the response for waitlist operations
type WaitlistResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
