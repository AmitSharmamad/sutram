// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/model/post.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a reply in a community topic
type Post struct {
	ID           int        `json:"id" db:"id"`
	TopicID      int        `json:"topic_id" db:"topic_id"`
	AuthorID     uuid.UUID  `json:"author_id" db:"author_id"`
	Content      string     `json:"content" db:"content"`
	IsSolution   bool       `json:"is_solution" db:"is_solution"`
	ParentPostID *int       `json:"parent_post_id" db:"parent_post_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	EditedByID   *uuid.UUID `json:"edited_by_id" db:"edited_by"`
	EditCount    int        `json:"edit_count" db:"edit_count"`

	// Additional fields not in database table
	LikeCount int `json:"like_count"`

	// Associated objects
	Topic    *Topic `json:"topic"`     // Loaded separately
	Author   User   `json:"author"`    // Loaded separately
	EditedBy *User  `json:"edited_by"` // Loaded separately if not nil
}
