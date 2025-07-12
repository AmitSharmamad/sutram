package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"algo/internal/domain/entity"
	"algo/internal/domain/repository"

	"github.com/google/uuid"
)

// UsageLogRepository implements repository.UsageLogRepository for PostgreSQL
type UsageLogRepository struct {
	db *sql.DB
}

// NewUsageLogRepository creates a new PostgreSQL usage log repository
func NewUsageLogRepository(db *sql.DB) repository.UsageLogRepository {
	return &UsageLogRepository{db: db}
}

// Create creates a new usage log entry
func (r *UsageLogRepository) Create(ctx context.Context, log *entity.UsageLog) error {
	query := `
		INSERT INTO user_usage_log (user_id, limit_type, action_type, previous_value, new_value, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	log.CreatedAt = time.Now()

	var metadataJSON []byte
	var err error
	if log.Metadata != nil {
		metadataJSON, err = json.Marshal(log.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	err = r.db.QueryRowContext(ctx, query,
		log.UserID,
		string(log.LimitType),
		log.ActionType,
		log.PreviousValue,
		log.NewValue,
		metadataJSON,
		log.CreatedAt,
	).Scan(&log.ID)

	if err != nil {
		return fmt.Errorf("failed to create usage log: %w", err)
	}

	return nil
}

// GetByUserID gets usage logs for a user with pagination
func (r *UsageLogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*entity.UsageLog, error) {
	query := `
		SELECT id, user_id, limit_type, action_type, previous_value, new_value, metadata, created_at
		FROM user_usage_log
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query usage logs: %w", err)
	}
	defer rows.Close()

	var logs []*entity.UsageLog
	for rows.Next() {
		log := &entity.UsageLog{}
		var limitTypeStr string
		var metadataJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&limitTypeStr,
			&log.ActionType,
			&log.PreviousValue,
			&log.NewValue,
			&metadataJSON,
			&log.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan usage log: %w", err)
		}

		log.LimitType = entity.UsageLimitType(limitTypeStr)

		if metadataJSON != nil {
			if err := json.Unmarshal(metadataJSON, &log.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating usage log rows: %w", err)
	}

	return logs, nil
}

// GetByUserIDAndType gets usage logs for a user and specific limit type with pagination
func (r *UsageLogRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, limitType entity.UsageLimitType, limit int, offset int) ([]*entity.UsageLog, error) {
	query := `
		SELECT id, user_id, limit_type, action_type, previous_value, new_value, metadata, created_at
		FROM user_usage_log
		WHERE user_id = $1 AND limit_type = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, userID, string(limitType), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query usage logs: %w", err)
	}
	defer rows.Close()

	var logs []*entity.UsageLog
	for rows.Next() {
		log := &entity.UsageLog{}
		var limitTypeStr string
		var metadataJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&limitTypeStr,
			&log.ActionType,
			&log.PreviousValue,
			&log.NewValue,
			&metadataJSON,
			&log.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan usage log: %w", err)
		}

		log.LimitType = entity.UsageLimitType(limitTypeStr)

		if metadataJSON != nil {
			if err := json.Unmarshal(metadataJSON, &log.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating usage log rows: %w", err)
	}

	return logs, nil
}
