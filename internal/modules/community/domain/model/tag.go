// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/model/tag.go
package model

import (
	"time"
)

// Tag represents a topic tag
type Tag struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	Color       *string   `json:"color" db:"color"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
