package repository

import (
	"context"

	"algo/internal/domain/entity"

	"github.com/google/uuid"
)

// UsageLimitRepository defines the interface for usage limit operations
type UsageLimitRepository interface {
	// Usage limit operations
	Create(ctx context.Context, limit *entity.UsageLimit) error
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) (*entity.UsageLimit, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UsageLimit, error)
	Update(ctx context.Context, limit *entity.UsageLimit) error
	IncrementUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) error
	DecrementUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) error
	ResetUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) error
	ResetAllUsageByType(ctx context.Context, limitType entity.UsageLimitType) error
	SetUserLimit(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, limitValue int) error

	// Batch operations
	GetUsersRequiringReset(ctx context.Context, limitType entity.UsageLimitType) ([]*entity.UsageLimit, error)
	BulkResetUsage(ctx context.Context, limitType entity.UsageLimitType) error
}

// UsageLogRepository defines the interface for usage log operations
type UsageLogRepository interface {
	Create(ctx context.Context, log *entity.UsageLog) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*entity.UsageLog, error)
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, limit int, offset int) ([]*entity.UsageLog, error)
}
