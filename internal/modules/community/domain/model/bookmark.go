// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/model/bookmark.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// Bookmark represents a bookmark on a topic
type Bookmark struct {
	ID        int       `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TopicID   int       `json:"topic_id" db:"topic_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
