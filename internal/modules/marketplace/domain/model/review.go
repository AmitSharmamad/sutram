// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/domain/model/review.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// Review represents a user review of a marketplace listing
type Review struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ListingID uuid.UUID `json:"listing_id" db:"listing_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Rating    int       `json:"rating" db:"rating"`
	Comment   string    `json:"comment" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Optional relationships
	UserName string `json:"user_name,omitempty" db:"-"` // For display purposes
}
