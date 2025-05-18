// filepath: /Users/amit/workspace/sutram/internal/modules/community/application/service/poll_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/community/domain/model"
	"algo/internal/modules/community/domain/repository"

	"github.com/google/uuid"
)

// Error definitions
var (
	ErrPollNotFound      = errors.New("poll not found")
	ErrPollClosed        = errors.New("poll is closed")
	ErrAlreadyVoted      = errors.New("already voted in this poll")
	ErrInvalidPollOption = errors.New("invalid poll option")
)

// PollService handles poll business logic
type PollService struct {
	pollRepo  repository.PollRepository
	topicRepo repository.TopicRepository
}

// NewPollService creates a new poll service
func NewPollService(
	pollRepo repository.PollRepository,
	topicRepo repository.TopicRepository,
) *PollService {
	return &PollService{
		pollRepo:  pollRepo,
		topicRepo: topicRepo,
	}
}

// CreatePoll creates a new poll
func (s *PollService) CreatePoll(
	ctx context.Context,
	topicID int,
	userID uuid.UUID,
	title string,
	description string,
	options []string,
	isMultipleChoice bool,
	closesAt *time.Time,
) (*model.Poll, error) {
	// Check if topic exists
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to check topic: %w", err)
	}

	if topic == nil {
		return nil, ErrTopicNotFound
	}

	// Create poll
	poll := &model.Poll{
		TopicID:          topicID,
		Title:            title,
		Description:      description,
		CreatedByID:      userID,
		IsMultipleChoice: isMultipleChoice,
		IsClosed:         false,
		ClosesAt:         closesAt,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Save poll
	if err := s.pollRepo.CreatePoll(ctx, poll); err != nil {
		return nil, fmt.Errorf("failed to create poll: %w", err)
	}

	// Create options
	poll.Options = make([]model.PollOption, len(options))
	for i, optionText := range options {
		option := &model.PollOption{
			PollID:       poll.ID,
			OptionText:   optionText,
			DisplayOrder: i,
		}

		if err := s.pollRepo.CreateOption(ctx, option); err != nil {
			return nil, fmt.Errorf("failed to create poll option: %w", err)
		}

		poll.Options[i] = *option
	}

	return poll, nil
}

// GetPollsByTopicID gets polls for a specific topic
func (s *PollService) GetPollsByTopicID(ctx context.Context, topicID int) ([]*model.Poll, error) {
	// Check if topic exists
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to check topic: %w", err)
	}

	if topic == nil {
		return nil, ErrTopicNotFound
	}

	// Get polls
	polls, err := s.pollRepo.GetPollsByTopicID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get polls: %w", err)
	}

	// Load options and vote counts for each poll
	for _, poll := range polls {
		options, err := s.pollRepo.GetOptionsByPollID(ctx, poll.ID)
		if err == nil {
			poll.Options = make([]model.PollOption, len(options))
			for i, option := range options {
				// Get vote count for each option
				count, err := s.pollRepo.GetVoteCount(ctx, option.ID)
				if err == nil {
					option.VoteCount = count
				}
				poll.Options[i] = *option
			}
		}
	}

	return polls, nil
}

// GetPollByID gets a poll by ID
func (s *PollService) GetPollByID(ctx context.Context, id int) (*model.Poll, error) {
	poll, err := s.pollRepo.GetPollByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}

	if poll == nil {
		return nil, ErrPollNotFound
	}

	// Load options and vote counts
	options, err := s.pollRepo.GetOptionsByPollID(ctx, poll.ID)
	if err == nil {
		poll.Options = make([]model.PollOption, len(options))
		for i, option := range options {
			// Get vote count
			count, err := s.pollRepo.GetVoteCount(ctx, option.ID)
			if err == nil {
				option.VoteCount = count
			}
			poll.Options[i] = *option
		}
	}

	return poll, nil
}

// VoteOnPoll allows a user to vote on a poll
func (s *PollService) VoteOnPoll(ctx context.Context, pollID int, optionID int, userID uuid.UUID) error {
	// Check if poll exists
	poll, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return fmt.Errorf("failed to check poll: %w", err)
	}

	if poll == nil {
		return ErrPollNotFound
	}

	// Check if poll is closed
	if poll.IsClosed {
		return ErrPollClosed
	}

	// Check if poll has timed out
	if poll.ClosesAt != nil && time.Now().After(*poll.ClosesAt) {
		// Auto-close the poll
		updates := map[string]interface{}{
			"is_closed": true,
		}
		if err := s.pollRepo.UpdatePoll(ctx, pollID, updates); err != nil {
			// Log error but continue
			fmt.Printf("failed to auto-close poll: %v\n", err)
		}
		return ErrPollClosed
	}

	// Check if option belongs to this poll
	options, err := s.pollRepo.GetOptionsByPollID(ctx, pollID)
	if err != nil {
		return fmt.Errorf("failed to get poll options: %w", err)
	}

	optionFound := false
	for _, option := range options {
		if option.ID == optionID {
			optionFound = true
			break
		}
	}

	if !optionFound {
		return ErrInvalidPollOption
	}

	// Check if user has already voted (unless it's a multiple choice poll)
	if !poll.IsMultipleChoice {
		hasVoted, err := s.pollRepo.HasUserVoted(ctx, pollID, userID)
		if err != nil {
			return fmt.Errorf("failed to check if user voted: %w", err)
		}

		if hasVoted {
			return ErrAlreadyVoted
		}
	}

	// Create vote
	vote := &model.PollVote{
		PollID:    pollID,
		OptionID:  optionID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	return s.pollRepo.CreateVote(ctx, vote)
}

// RemoveVote removes a user's vote from a poll option
func (s *PollService) RemoveVote(ctx context.Context, pollID int, optionID int, userID uuid.UUID) error {
	// Check if poll exists
	poll, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return fmt.Errorf("failed to check poll: %w", err)
	}

	if poll == nil {
		return ErrPollNotFound
	}

	// Check if poll is closed
	if poll.IsClosed {
		return ErrPollClosed
	}

	return s.pollRepo.DeleteVote(ctx, userID, pollID, optionID)
}
