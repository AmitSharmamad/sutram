// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/model/badge.go
package model

import (
	"time"
)

// Badge represents a badge that can be earned by users
type Badge struct {
	ID          int       `json:"id" db:"id"`
	TypeID      int       `json:"type_id" db:"type_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Requirement string    `json:"requirement" db:"requirement"`
	Tier        int       `json:"tier" db:"tier"`
	Points      int       `json:"points" db:"points"`
	Icon        *string   `json:"icon" db:"icon"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	// Additional fields
	TypeName  string    `json:"type_name"`  // Loaded from badge_types
	AwardedAt time.Time `json:"awarded_at"` // Loaded from user_badges
}
