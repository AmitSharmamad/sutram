// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/model/user.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system with minimal necessary fields for community
type User struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

// UserProfile represents a user's community profile
type UserProfile struct {
	ID                uuid.UUID `json:"id" db:"user_id"`
	Name              string    `json:"name"`
	ReputationScore   int       `json:"reputation_score" db:"reputation_score"`
	ContributionCount int       `json:"contribution_count" db:"contribution_count"`
	JoinedDate        time.Time `json:"joined_date" db:"joined_date"`
	LastActive        time.Time `json:"last_active" db:"last_active"`

	// Additional fields not in database table
	TopicCount int      `json:"topic_count"`
	PostCount  int      `json:"post_count"`
	Badges     []*Badge `json:"badges"` // Loaded separately
}
