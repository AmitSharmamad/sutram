package model

import (
	"time"

	"github.com/google/uuid"
)

// Waitlist represents a waitlist entry
type Waitlist struct {
	ID        int64      `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	SurveyID  *uuid.UUID `json:"survey_id" db:"survey_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// NewWaitlist creates a new waitlist entry
func NewWaitlist(email string, surveyID *uuid.UUID) *Waitlist {
	now := time.Now()
	return &Waitlist{
		Email:     email,
		SurveyID:  surveyID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
