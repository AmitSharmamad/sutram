// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/repository/poll_repository.go
package repository

import (
	"context"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// PollRepository defines the operations for poll persistence
type PollRepository interface {
	// Create a new poll
	CreatePoll(ctx context.Context, poll *model.Poll) error

	// Get a poll by ID
	GetPollByID(ctx context.Context, id int) (*model.Poll, error)

	// Get polls for a topic
	GetPollsByTopicID(ctx context.Context, topicID int) ([]*model.Poll, error)

	// Update a poll
	UpdatePoll(ctx context.Context, id int, updates map[string]interface{}) error

	// Delete a poll
	DeletePoll(ctx context.Context, id int) error

	// Create a poll option
	CreateOption(ctx context.Context, option *model.PollOption) error

	// Get options for a poll
	GetOptionsByPollID(ctx context.Context, pollID int) ([]*model.PollOption, error)

	// Update a poll option
	UpdateOption(ctx context.Context, id int, updates map[string]interface{}) error

	// Delete a poll option
	DeleteOption(ctx context.Context, id int) error

	// Create a vote
	CreateVote(ctx context.Context, vote *model.PollVote) error

	// Get votes for a poll
	GetVotesByPollID(ctx context.Context, pollID int) ([]*model.PollVote, error)

	// Check if a user has voted on a poll
	HasUserVoted(ctx context.Context, pollID int, userID uuid.UUID) (bool, error)

	// Get vote count for a poll option
	GetVoteCount(ctx context.Context, optionID int) (int, error)

	// Delete a vote
	DeleteVote(ctx context.Context, userID uuid.UUID, pollID int, optionID int) error
}
