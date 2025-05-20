// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/model/poll.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// Poll represents a community poll
type Poll struct {
	ID               int        `json:"id" db:"id"`
	TopicID          int        `json:"topic_id" db:"topic_id"`
	Title            string     `json:"title" db:"title"`
	Description      string     `json:"description" db:"description"`
	CreatedByID      uuid.UUID  `json:"created_by_id" db:"created_by"`
	IsMultipleChoice bool       `json:"is_multiple_choice" db:"is_multiple_choice"`
	IsClosed         bool       `json:"is_closed" db:"is_closed"`
	ClosesAt         *time.Time `json:"closes_at" db:"closes_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`

	// Associated objects
	CreatedBy User         `json:"created_by" db:"-"` // User who created the poll
	Options   []PollOption `json:"options" db:"-"`    // Loaded separately
	Votes     []PollVote   `json:"votes" db:"-"`      // Loaded separately
}

// PollOption represents an option in a poll
type PollOption struct {
	ID           int    `json:"id" db:"id"`
	PollID       int    `json:"poll_id" db:"poll_id"`
	OptionText   string `json:"option_text" db:"option_text"`
	DisplayOrder int    `json:"display_order" db:"display_order"`

	// Additional fields not in database table
	VoteCount int `json:"vote_count"`
}

// PollVote represents a vote on a poll option
type PollVote struct {
	ID        int       `json:"id" db:"id"`
	PollID    int       `json:"poll_id" db:"poll_id"`
	OptionID  int       `json:"option_id" db:"option_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Associated objects
	User   User       `json:"user" db:"-"`   // User who voted
	Option PollOption `json:"option" db:"-"` // The option that was voted for
}
