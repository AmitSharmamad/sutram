// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/model/badge.go
package model

import (
	"time"
)

// Badge represents a badge that can be earned by users
type Badge struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	Icon        *string   `json:"icon" db:"icon"`
	Color       *string   `json:"color" db:"color"`
	Points      int       `json:"points" db:"points"`
	BadgeType   string    `json:"badge_type" db:"badge_type"`
	Criteria    *string   `json:"criteria" db:"criteria"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	// Additional fields
	TypeName  string     `json:"type_name,omitempty" db:"-"`  // Loaded from badge_types
	AwardedAt *time.Time `json:"awarded_at,omitempty" db:"-"` // Loaded from user_badges
}
