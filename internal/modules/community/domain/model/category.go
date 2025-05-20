// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/model/category.go
package model

import (
	"time"
)

// Category represents a community topic category
type Category struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	Description  string    `json:"description" db:"description"`
	Icon         *string   `json:"icon" db:"icon"`
	Color        *string   `json:"color" db:"color"`
	ParentID     *int      `json:"parent_id" db:"parent_id"`
	DisplayOrder int       `json:"display_order" db:"display_order"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	// Additional fields not in database table
	TopicCount int `json:"topic_count"`
}
