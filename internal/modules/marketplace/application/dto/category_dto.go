// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/dto/category_dto.go
package dto

import (
	"github.com/google/uuid"
)

// CategoryResponse represents a marketplace category
type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Icon        string    `json:"icon,omitempty"`
}
