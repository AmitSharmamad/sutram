// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/domain/repository/performance_repository.go
package repository

import (
	"context"

	"algo/internal/modules/marketplace/domain/model"

	"github.com/google/uuid"
)

// PerformanceRepository defines the interface for marketplace performance repository
type PerformanceRepository interface {
	// Create creates performance metrics for a listing
	Create(ctx context.Context, performance *model.Performance) error

	// FindByListingID returns performance metrics for a listing
	FindByListingID(ctx context.Context, listingID uuid.UUID) (*model.Performance, error)

	// Update updates performance metrics for a listing
	Update(ctx context.Context, performance *model.Performance) error

	// Delete deletes performance metrics for a listing
	Delete(ctx context.Context, listingID uuid.UUID) error
}
