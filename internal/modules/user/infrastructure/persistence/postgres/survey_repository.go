package postgres

import (
	"context"
	"database/sql"
	"errors"

	"algo/internal/modules/user/domain/model"
	"algo/internal/modules/user/domain/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SurveyRepository implements the survey repository interface
type SurveyRepository struct {
	db *sqlx.DB
}

// NewSurveyRepository creates a new survey repository
func NewSurveyRepository(db *sqlx.DB) repository.SurveyRepository {
	return &SurveyRepository{db: db}
}

// Create creates a new survey entry
func (r *SurveyRepository) Create(ctx context.Context, survey *model.Survey) error {
	query := `
		INSERT INTO survey (email, nature_of_work, trading_type, experience_level, work_type, age, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		survey.Email,
		survey.NatureOfWork,
		survey.TradingType,
		survey.ExperienceLevel,
		survey.WorkType,
		survey.Age,
		survey.CreatedAt,
		survey.UpdatedAt).Scan(&survey.ID)
	if err != nil {
		return err
	}

	return nil
}

// FindByEmail finds a survey entry by email
func (r *SurveyRepository) FindByEmail(ctx context.Context, email string) (*model.Survey, error) {
	query := `
		SELECT id, email, nature_of_work, trading_type, experience_level, work_type, age, created_at, updated_at
		FROM survey
		WHERE email = $1
	`

	var survey model.Survey
	err := r.db.GetContext(ctx, &survey, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &survey, nil
}

// FindByID finds a survey entry by ID
func (r *SurveyRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Survey, error) {
	query := `
		SELECT id, email, nature_of_work, trading_type, experience_level, work_type, age, created_at, updated_at
		FROM survey
		WHERE id = $1
	`

	var survey model.Survey
	err := r.db.GetContext(ctx, &survey, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &survey, nil
}

// FindAll retrieves all survey entries
func (r *SurveyRepository) FindAll(ctx context.Context) ([]model.Survey, error) {
	query := `
		SELECT id, email, nature_of_work, trading_type, experience_level, work_type, age, created_at, updated_at
		FROM survey
		ORDER BY created_at DESC
	`

	var surveys []model.Survey
	err := r.db.SelectContext(ctx, &surveys, query)
	if err != nil {
		return nil, err
	}

	return surveys, nil
}

// GetCount returns the total number of survey entries
func (r *SurveyRepository) GetCount(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM survey`

	var count int64
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}
