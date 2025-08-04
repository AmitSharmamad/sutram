package postgres

import (
	"context"
	"database/sql"
	"errors"

	"algo/internal/modules/user/domain/model"
	"algo/internal/modules/user/domain/repository"

	"github.com/jmoiron/sqlx"
)

// WaitlistRepository implements the waitlist repository interface
type WaitlistRepository struct {
	db *sqlx.DB
}

// NewWaitlistRepository creates a new waitlist repository
func NewWaitlistRepository(db *sqlx.DB) repository.WaitlistRepository {
	return &WaitlistRepository{db: db}
}

// Create creates a new waitlist entry
func (r *WaitlistRepository) Create(ctx context.Context, waitlist *model.Waitlist) error {
	query := `
		INSERT INTO waitlist (email, survey_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, waitlist.Email, waitlist.SurveyID, waitlist.CreatedAt, waitlist.UpdatedAt).Scan(&waitlist.ID)
	if err != nil {
		return err
	}

	return nil
}

// FindByEmail finds a waitlist entry by email
func (r *WaitlistRepository) FindByEmail(ctx context.Context, email string) (*model.Waitlist, error) {
	query := `
		SELECT id, email, survey_id, created_at, updated_at
		FROM waitlist
		WHERE email = $1
	`

	var waitlist model.Waitlist
	err := r.db.GetContext(ctx, &waitlist, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &waitlist, nil
}

// FindAll retrieves all waitlist entries
func (r *WaitlistRepository) FindAll(ctx context.Context) ([]model.Waitlist, error) {
	query := `
		SELECT id, email, survey_id, created_at, updated_at
		FROM waitlist
		ORDER BY created_at DESC
	`

	var waitlists []model.Waitlist
	err := r.db.SelectContext(ctx, &waitlists, query)
	if err != nil {
		return nil, err
	}

	return waitlists, nil
}

// GetCount returns the total number of waitlist entries
func (r *WaitlistRepository) GetCount(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM waitlist`

	var count int64
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}
