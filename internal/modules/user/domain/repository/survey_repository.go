package repository

import (
	"context"

	"algo/internal/modules/user/domain/model"

	"github.com/google/uuid"
)

// SurveyRepository defines the interface for survey operations
type SurveyRepository interface {
	// Create a new survey entry
	Create(ctx context.Context, survey *model.Survey) error

	// FindByEmail finds a survey entry by email
	FindByEmail(ctx context.Context, email string) (*model.Survey, error)

	// FindByID finds a survey entry by ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Survey, error)

	// FindAll retrieves all survey entries
	FindAll(ctx context.Context) ([]model.Survey, error)

	// GetCount returns the total number of survey entries
	GetCount(ctx context.Context) (int64, error)
}
