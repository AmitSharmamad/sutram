// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/poll_repo.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// PollRepository implements the repository.PollRepository interface with PostgreSQL
type PollRepository struct {
	db *sql.DB
}

// CreatePoll creates a new poll
func (r *PollRepository) CreatePoll(ctx context.Context, poll *model.Poll) error {
	query := `
		INSERT INTO polls (
			topic_id, title, description, created_by, is_multiple_choice, is_closed, closes_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		poll.TopicID,
		poll.Title,
		poll.Description,
		poll.CreatedByID,
		poll.IsMultipleChoice,
		poll.IsClosed,
		poll.ClosesAt,
		poll.CreatedAt,
		poll.UpdatedAt,
	).Scan(&poll.ID)

	if err != nil {
		return fmt.Errorf("failed to create poll: %w", err)
	}

	return nil
}

// GetPollByID retrieves a poll by its ID
func (r *PollRepository) GetPollByID(ctx context.Context, id int) (*model.Poll, error) {
	query := `
		SELECT p.id, p.topic_id, p.title, p.description, p.created_by, p.is_multiple_choice,
			p.is_closed, p.closes_at, p.created_at, p.updated_at,
			u.name AS created_by_name
		FROM polls p
		JOIN "user" u ON p.created_by = u.id
		WHERE p.id = $1
	`

	var poll model.Poll
	var closesAt sql.NullTime
	var createdByName string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&poll.ID,
		&poll.TopicID,
		&poll.Title,
		&poll.Description,
		&poll.CreatedByID,
		&poll.IsMultipleChoice,
		&poll.IsClosed,
		&closesAt,
		&poll.CreatedAt,
		&poll.UpdatedAt,
		&createdByName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}

	// Handle nullable fields
	if closesAt.Valid {
		poll.ClosesAt = &closesAt.Time
	}

	// Set created by user
	poll.CreatedBy = model.User{
		ID:   poll.CreatedByID,
		Name: createdByName,
	}

	// Get options
	options, err := r.GetOptionsByPollID(ctx, id)
	if err == nil {
		poll.Options = make([]model.PollOption, len(options))
		for i, opt := range options {
			poll.Options[i] = *opt
		}
	}

	return &poll, nil
}

// GetPollsByTopicID retrieves polls for a specific topic
func (r *PollRepository) GetPollsByTopicID(ctx context.Context, topicID int) ([]*model.Poll, error) {
	query := `
		SELECT p.id, p.topic_id, p.title, p.description, p.created_by, p.is_multiple_choice,
			p.is_closed, p.closes_at, p.created_at, p.updated_at,
			u.name AS created_by_name
		FROM polls p
		JOIN "user" u ON p.created_by = u.id
		WHERE p.topic_id = $1
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to query topic polls: %w", err)
	}
	defer rows.Close()

	var polls []*model.Poll
	for rows.Next() {
		var poll model.Poll
		var closesAt sql.NullTime
		var createdByName string

		err := rows.Scan(
			&poll.ID,
			&poll.TopicID,
			&poll.Title,
			&poll.Description,
			&poll.CreatedByID,
			&poll.IsMultipleChoice,
			&poll.IsClosed,
			&closesAt,
			&poll.CreatedAt,
			&poll.UpdatedAt,
			&createdByName,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan poll row: %w", err)
		}

		// Handle nullable fields
		if closesAt.Valid {
			poll.ClosesAt = &closesAt.Time
		}

		// Set created by user
		poll.CreatedBy = model.User{
			ID:   poll.CreatedByID,
			Name: createdByName,
		}

		polls = append(polls, &poll)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during poll row iteration: %w", err)
	}

	return polls, nil
}

// UpdatePoll updates a poll
func (r *PollRepository) UpdatePoll(ctx context.Context, id int, updates map[string]interface{}) error {
	// Build dynamic query
	setClause := ""
	args := []interface{}{id} // First arg is the ID
	argCount := 2

	for key, value := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += key + " = $" + fmt.Sprintf("%d", argCount)
		args = append(args, value)
		argCount++
	}

	if setClause == "" {
		return nil // Nothing to update
	}

	// Add updated_at timestamp
	if _, ok := updates["updated_at"]; !ok {
		if setClause != "" {
			setClause += ", "
		}
		setClause += "updated_at = $" + fmt.Sprintf("%d", argCount)
		args = append(args, time.Now())
	}

	query := "UPDATE polls SET " + setClause + " WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update poll: %w", err)
	}

	return nil
}

// DeletePoll deletes a poll
func (r *PollRepository) DeletePoll(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM polls WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete poll: %w", err)
	}

	return nil
}

// CreateOption creates a poll option
func (r *PollRepository) CreateOption(ctx context.Context, option *model.PollOption) error {
	query := `
		INSERT INTO poll_options (poll_id, option_text, display_order)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		option.PollID,
		option.OptionText,
		option.DisplayOrder,
	).Scan(&option.ID)

	if err != nil {
		return fmt.Errorf("failed to create poll option: %w", err)
	}

	return nil
}

// GetOptionsByPollID gets options for a poll
func (r *PollRepository) GetOptionsByPollID(ctx context.Context, pollID int) ([]*model.PollOption, error) {
	query := `
		SELECT id, poll_id, option_text, display_order
		FROM poll_options
		WHERE poll_id = $1
		ORDER BY display_order ASC
	`

	rows, err := r.db.QueryContext(ctx, query, pollID)
	if err != nil {
		return nil, fmt.Errorf("failed to query poll options: %w", err)
	}
	defer rows.Close()

	var options []*model.PollOption
	for rows.Next() {
		var option model.PollOption

		err := rows.Scan(
			&option.ID,
			&option.PollID,
			&option.OptionText,
			&option.DisplayOrder,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan option row: %w", err)
		}

		// Get vote count
		count, err := r.GetVoteCount(ctx, option.ID)
		if err == nil {
			option.VoteCount = count
		}

		options = append(options, &option)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during option row iteration: %w", err)
	}

	return options, nil
}

// UpdateOption updates a poll option
func (r *PollRepository) UpdateOption(ctx context.Context, id int, updates map[string]interface{}) error {
	// Build dynamic query
	setClause := ""
	args := []interface{}{id} // First arg is the ID
	argCount := 2

	for key, value := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += key + " = $" + fmt.Sprintf("%d", argCount)
		args = append(args, value)
		argCount++
	}

	if setClause == "" {
		return nil // Nothing to update
	}

	query := "UPDATE poll_options SET " + setClause + " WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update poll option: %w", err)
	}

	return nil
}

// DeleteOption deletes a poll option
func (r *PollRepository) DeleteOption(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM poll_options WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete poll option: %w", err)
	}

	return nil
}

// CreateVote creates a poll vote
func (r *PollRepository) CreateVote(ctx context.Context, vote *model.PollVote) error {
	query := `
		INSERT INTO poll_votes (poll_id, option_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		vote.PollID,
		vote.OptionID,
		vote.UserID,
		vote.CreatedAt,
	).Scan(&vote.ID)

	if err != nil {
		return fmt.Errorf("failed to create vote: %w", err)
	}

	return nil
}

// GetVotesByPollID gets votes for a poll
func (r *PollRepository) GetVotesByPollID(ctx context.Context, pollID int) ([]*model.PollVote, error) {
	query := `
		SELECT v.id, v.poll_id, v.option_id, v.user_id, v.created_at,
			u.name AS user_name,
			o.option_text
		FROM poll_votes v
		JOIN "user" u ON v.user_id = u.id
		JOIN poll_options o ON v.option_id = o.id
		WHERE v.poll_id = $1
		ORDER BY v.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, pollID)
	if err != nil {
		return nil, fmt.Errorf("failed to query poll votes: %w", err)
	}
	defer rows.Close()

	var votes []*model.PollVote
	for rows.Next() {
		var vote model.PollVote
		var userName string
		var optionText string

		err := rows.Scan(
			&vote.ID,
			&vote.PollID,
			&vote.OptionID,
			&vote.UserID,
			&vote.CreatedAt,
			&userName,
			&optionText,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan vote row: %w", err)
		}

		// Set user and option
		vote.User = model.User{
			ID:   vote.UserID,
			Name: userName,
		}

		vote.Option = model.PollOption{
			ID:         vote.OptionID,
			PollID:     vote.PollID,
			OptionText: optionText,
		}

		votes = append(votes, &vote)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during vote row iteration: %w", err)
	}

	return votes, nil
}

// HasUserVoted checks if a user has voted on a poll
func (r *PollRepository) HasUserVoted(ctx context.Context, pollID int, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM poll_votes WHERE poll_id = $1 AND user_id = $2)",
		pollID,
		userID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check vote existence: %w", err)
	}

	return exists, nil
}

// GetVoteCount gets the vote count for an option
func (r *PollRepository) GetVoteCount(ctx context.Context, optionID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM poll_votes WHERE option_id = $1",
		optionID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get vote count: %w", err)
	}

	return count, nil
}

// DeleteVote deletes a vote
func (r *PollRepository) DeleteVote(ctx context.Context, userID uuid.UUID, pollID int, optionID int) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM poll_votes WHERE user_id = $1 AND poll_id = $2 AND option_id = $3",
		userID,
		pollID,
		optionID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete vote: %w", err)
	}

	return nil
}
