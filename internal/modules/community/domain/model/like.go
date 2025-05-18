// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/model/like.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// Like represents a like on a topic or post
type Like struct {
	ID           int       `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	LikeableType string    `json:"likeable_type" db:"likeable_type"` // "topic" or "post"
	LikeableID   int       `json:"likeable_id" db:"likeable_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
