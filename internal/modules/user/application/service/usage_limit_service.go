package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"algo/internal/domain/entity"
	"algo/internal/domain/repository"

	"github.com/google/uuid"
)

// UsageLimitService defines the interface for usage limit operations
type UsageLimitService interface {
	// Limit validation
	CanPerformAction(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) (bool, error)
	GetRemainingUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) (int, error)

	// Usage tracking
	IncrementUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, metadata map[string]interface{}) error
	DecrementUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, metadata map[string]interface{}) error

	// Limit management
	SetUserLimit(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, limitValue int) error
	GetUserLimits(ctx context.Context, userID uuid.UUID) ([]*entity.UsageLimit, error)
	ResetUsageForUser(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) error

	// Periodic resets
	ResetMonthlyUsage(ctx context.Context) error
	ResetDailyUsage(ctx context.Context) error
	ResetHourlyUsage(ctx context.Context) error

	// Usage history
	GetUsageHistory(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, limit int, offset int) ([]*entity.UsageLog, error)
}

// usageLimitService implements UsageLimitService
type usageLimitService struct {
	usageLimitRepo repository.UsageLimitRepository
	usageLogRepo   repository.UsageLogRepository
}

// NewUsageLimitService creates a new usage limit service
func NewUsageLimitService(
	usageLimitRepo repository.UsageLimitRepository,
	usageLogRepo repository.UsageLogRepository,
) UsageLimitService {
	return &usageLimitService{
		usageLimitRepo: usageLimitRepo,
		usageLogRepo:   usageLogRepo,
	}
}

// CanPerformAction checks if a user can perform an action based on their limits
func (s *usageLimitService) CanPerformAction(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) (bool, error) {
	limit, err := s.usageLimitRepo.GetByUserIDAndType(ctx, userID, limitType)
	if err != nil {
		return false, fmt.Errorf("failed to get usage limit: %w", err)
	}

	if limit == nil {
		// If no limit exists, create a default one based on user role
		// For now, allow the action but log the missing limit
		return true, nil
	}

	// Check if usage should be reset
	now := time.Now()
	if limit.ShouldReset(now) {
		if err := s.resetUsageInternal(ctx, limit, now); err != nil {
			return false, fmt.Errorf("failed to reset usage: %w", err)
		}
		// Update the limit object with reset values
		limit.CurrentUsage = 0
		limit.LastResetAt = now
	}

	return limit.CanPerformAction(), nil
}

// GetRemainingUsage returns the remaining usage for a user and limit type
func (s *usageLimitService) GetRemainingUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) (int, error) {
	limit, err := s.usageLimitRepo.GetByUserIDAndType(ctx, userID, limitType)
	if err != nil {
		return 0, fmt.Errorf("failed to get usage limit: %w", err)
	}

	if limit == nil {
		return 0, fmt.Errorf("usage limit not found")
	}

	// Check if usage should be reset
	now := time.Now()
	if limit.ShouldReset(now) {
		if err := s.resetUsageInternal(ctx, limit, now); err != nil {
			return 0, fmt.Errorf("failed to reset usage: %w", err)
		}
		// Update the limit object with reset values
		limit.CurrentUsage = 0
		limit.LastResetAt = now
	}

	return limit.GetRemainingUsage(), nil
}

// IncrementUsage increments the usage for a user and limit type
func (s *usageLimitService) IncrementUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, metadata map[string]interface{}) error {
	// Get current limit
	limit, err := s.usageLimitRepo.GetByUserIDAndType(ctx, userID, limitType)
	if err != nil {
		return fmt.Errorf("failed to get usage limit: %w", err)
	}

	if limit == nil {
		return fmt.Errorf("usage limit not found for user %s and type %s", userID, limitType)
	}

	// Check if usage should be reset
	now := time.Now()
	if limit.ShouldReset(now) {
		if err := s.resetUsageInternal(ctx, limit, now); err != nil {
			return fmt.Errorf("failed to reset usage: %w", err)
		}
		limit.CurrentUsage = 0
	}

	previousValue := limit.CurrentUsage

	// Increment usage
	if err := s.usageLimitRepo.IncrementUsage(ctx, userID, limitType); err != nil {
		return fmt.Errorf("failed to increment usage: %w", err)
	}

	// Log the usage change
	if err := s.logUsageChange(ctx, userID, limitType, "increment", previousValue, previousValue+1, metadata); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to log usage change: %v\n", err)
	}

	return nil
}

// DecrementUsage decrements the usage for a user and limit type
func (s *usageLimitService) DecrementUsage(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, metadata map[string]interface{}) error {
	// Get current limit
	limit, err := s.usageLimitRepo.GetByUserIDAndType(ctx, userID, limitType)
	if err != nil {
		return fmt.Errorf("failed to get usage limit: %w", err)
	}

	if limit == nil {
		return fmt.Errorf("usage limit not found for user %s and type %s", userID, limitType)
	}

	if limit.CurrentUsage <= 0 {
		return nil // Cannot decrement below 0
	}

	previousValue := limit.CurrentUsage

	// Decrement usage
	if err := s.usageLimitRepo.DecrementUsage(ctx, userID, limitType); err != nil {
		return fmt.Errorf("failed to decrement usage: %w", err)
	}

	// Log the usage change
	if err := s.logUsageChange(ctx, userID, limitType, "decrement", previousValue, previousValue-1, metadata); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to log usage change: %v\n", err)
	}

	return nil
}

// SetUserLimit sets a custom limit for a user
func (s *usageLimitService) SetUserLimit(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, limitValue int) error {
	return s.usageLimitRepo.SetUserLimit(ctx, userID, limitType, limitValue)
}

// GetUserLimits gets all limits for a user
func (s *usageLimitService) GetUserLimits(ctx context.Context, userID uuid.UUID) ([]*entity.UsageLimit, error) {
	return s.usageLimitRepo.GetByUserID(ctx, userID)
}

// ResetUsageForUser resets usage for a specific user and limit type
func (s *usageLimitService) ResetUsageForUser(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType) error {
	// Get current limit
	limit, err := s.usageLimitRepo.GetByUserIDAndType(ctx, userID, limitType)
	if err != nil {
		return fmt.Errorf("failed to get usage limit: %w", err)
	}

	if limit == nil {
		return fmt.Errorf("usage limit not found")
	}

	return s.resetUsageInternal(ctx, limit, time.Now())
}

// ResetMonthlyUsage resets all monthly usage limits
func (s *usageLimitService) ResetMonthlyUsage(ctx context.Context) error {
	return s.usageLimitRepo.BulkResetUsage(ctx, entity.BacktestPerMonth)
}

// ResetDailyUsage resets all daily usage limits
func (s *usageLimitService) ResetDailyUsage(ctx context.Context) error {
	return s.usageLimitRepo.BulkResetUsage(ctx, entity.LiveTradesPerDay)
}

// ResetHourlyUsage resets all hourly usage limits
func (s *usageLimitService) ResetHourlyUsage(ctx context.Context) error {
	return s.usageLimitRepo.BulkResetUsage(ctx, entity.APICallsPerHour)
}

// GetUsageHistory gets usage history for a user and limit type
func (s *usageLimitService) GetUsageHistory(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, limit int, offset int) ([]*entity.UsageLog, error) {
	return s.usageLogRepo.GetByUserIDAndType(ctx, userID, limitType, limit, offset)
}

// resetUsageInternal resets usage for a specific limit
func (s *usageLimitService) resetUsageInternal(ctx context.Context, limit *entity.UsageLimit, resetTime time.Time) error {
	previousValue := limit.CurrentUsage

	// Reset usage in database
	if err := s.usageLimitRepo.ResetUsage(ctx, limit.UserID, limit.LimitType); err != nil {
		return fmt.Errorf("failed to reset usage in database: %w", err)
	}

	// Log the reset
	metadata := map[string]interface{}{
		"reset_reason": "periodic_reset",
		"reset_period": string(limit.ResetPeriod),
		"reset_time":   resetTime.Format(time.RFC3339),
	}

	if err := s.logUsageChange(ctx, limit.UserID, limit.LimitType, "reset", previousValue, 0, metadata); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to log usage reset: %v\n", err)
	}

	return nil
}

// logUsageChange logs a usage change to the usage log
func (s *usageLimitService) logUsageChange(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, actionType string, previousValue, newValue int, metadata map[string]interface{}) error {
	// Convert metadata to ensure it's serializable
	var metadataJSON map[string]interface{}
	if metadata != nil {
		// Ensure all values are JSON serializable
		jsonBytes, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &metadataJSON); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	log := &entity.UsageLog{
		UserID:        userID,
		LimitType:     limitType,
		ActionType:    actionType,
		PreviousValue: previousValue,
		NewValue:      newValue,
		Metadata:      metadataJSON,
		CreatedAt:     time.Now(),
	}

	return s.usageLogRepo.Create(ctx, log)
}
