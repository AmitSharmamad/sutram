package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"algo/internal/modules/backtest/domain/model"
	"algo/internal/modules/backtest/domain/repository"

	"github.com/google/uuid"
)

// PositionRepository implements the repository.PositionRepository interface

// PositionRepository implements the repository.PositionRepository interface
type PositionRepository struct {
	db *sql.DB
}

// NewPositionRepository creates a new PostgreSQL repository for positions
func NewPositionRepository(db *sql.DB) repository.PositionRepository {
	return &PositionRepository{
		db: db,
	}
}

// Create saves a new position in the database
func (r *PositionRepository) Create(ctx context.Context, position *model.Position) error {
	query := `
		INSERT INTO position (
			id, run_id, state, direction, entry_price_signal, entry_time_signal,
			entry_price_order, entry_time_order, entry_signal, exit_price_signal,
			exit_time_signal, exit_price_order, exit_time_order, exit_signal,
			exit_reason, profit_loss_absolute, profit_loss_percentage, 
			transaction_cost, entry_fees, exit_fees, net_profit_loss,
			position_size, entry_order_id, exit_order_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		position.ID,
		position.RunID,
		position.State,
		position.Direction,
		position.EntryPriceSignal,
		position.EntryTimeSignal,
		position.EntryPriceOrder,
		position.EntryTimeOrder,
		position.EntrySignal,
		position.ExitPriceSignal,
		position.ExitTimeSignal,
		position.ExitPriceOrder,
		position.ExitTimeOrder,
		position.ExitSignal,
		position.ExitReason,
		position.ProfitLossAbsolute,
		position.ProfitLossPercentage,
		position.TransactionCost,
		position.EntryFees,
		position.ExitFees,
		position.NetProfitLoss,
		position.PositionSize,
		position.EntryOrderID,
		position.ExitOrderID,
		position.CreatedAt,
		position.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetByID retrieves a position by its ID
func (r *PositionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Position, error) {
	query := `
		SELECT 
			id, run_id, state, direction, entry_price_signal, entry_time_signal,
			entry_price_order, entry_time_order, entry_signal, exit_price_signal,
			exit_time_signal, exit_price_order, exit_time_order, exit_signal,
			exit_reason, profit_loss_absolute, profit_loss_percentage, 
			transaction_cost, entry_fees, exit_fees, net_profit_loss,
			position_size, entry_order_id, exit_order_id, created_at, updated_at
		FROM position
		WHERE id = $1
	`

	var position model.Position
	var entryOrderID, exitOrderID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&position.ID,
		&position.RunID,
		&position.State,
		&position.Direction,
		&position.EntryPriceSignal,
		&position.EntryTimeSignal,
		&position.EntryPriceOrder,
		&position.EntryTimeOrder,
		&position.EntrySignal,
		&position.ExitPriceSignal,
		&position.ExitTimeSignal,
		&position.ExitPriceOrder,
		&position.ExitTimeOrder,
		&position.ExitSignal,
		&position.ExitReason,
		&position.ProfitLossAbsolute,
		&position.ProfitLossPercentage,
		&position.TransactionCost,
		&position.EntryFees,
		&position.ExitFees,
		&position.NetProfitLoss,
		&position.PositionSize,
		&entryOrderID,
		&exitOrderID,
		&position.CreatedAt,
		&position.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPositionNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	// Handle nullable UUIDs
	if entryOrderID.Valid {
		entryID, err := uuid.Parse(entryOrderID.String)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid entry order ID: %v", ErrDatabaseError, err)
		}
		position.EntryOrderID = &entryID
	}

	if exitOrderID.Valid {
		exitID, err := uuid.Parse(exitOrderID.String)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid exit order ID: %v", ErrDatabaseError, err)
		}
		position.ExitOrderID = &exitID
	}

	return &position, nil
}

// GetByRunID retrieves all positions for a specific run
func (r *PositionRepository) GetByRunID(ctx context.Context, runID uuid.UUID) ([]*model.Position, error) {
	query := `
		SELECT 
			id, run_id, state, direction, entry_price_signal, entry_time_signal,
			entry_price_order, entry_time_order, entry_signal, exit_price_signal,
			exit_time_signal, exit_price_order, exit_time_order, exit_signal,
			exit_reason, profit_loss_absolute, profit_loss_percentage, 
			transaction_cost, entry_fees, exit_fees, net_profit_loss,
			position_size, entry_order_id, exit_order_id, created_at, updated_at
		FROM position
		WHERE run_id = $1
		ORDER BY entry_time_signal ASC
	`

	rows, err := r.db.QueryContext(ctx, query, runID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var positions []*model.Position
	for rows.Next() {
		var position model.Position
		var entryOrderID, exitOrderID sql.NullString

		if err := rows.Scan(
			&position.ID,
			&position.RunID,
			&position.State,
			&position.Direction,
			&position.EntryPriceSignal,
			&position.EntryTimeSignal,
			&position.EntryPriceOrder,
			&position.EntryTimeOrder,
			&position.EntrySignal,
			&position.ExitPriceSignal,
			&position.ExitTimeSignal,
			&position.ExitPriceOrder,
			&position.ExitTimeOrder,
			&position.ExitSignal,
			&position.ExitReason,
			&position.ProfitLossAbsolute,
			&position.ProfitLossPercentage,
			&position.TransactionCost,
			&position.EntryFees,
			&position.ExitFees,
			&position.NetProfitLoss,
			&position.PositionSize,
			&entryOrderID,
			&exitOrderID,
			&position.CreatedAt,
			&position.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		// Handle nullable UUIDs
		if entryOrderID.Valid {
			entryID, err := uuid.Parse(entryOrderID.String)
			if err != nil {
				return nil, fmt.Errorf("%w: invalid entry order ID: %v", ErrDatabaseError, err)
			}
			position.EntryOrderID = &entryID
		}

		if exitOrderID.Valid {
			exitID, err := uuid.Parse(exitOrderID.String)
			if err != nil {
				return nil, fmt.Errorf("%w: invalid exit order ID: %v", ErrDatabaseError, err)
			}
			position.ExitOrderID = &exitID
		}

		positions = append(positions, &position)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return positions, nil
}

// Update updates an existing position in the database
func (r *PositionRepository) Update(ctx context.Context, position *model.Position) error {
	query := `
		UPDATE position
		SET 
			state = $2,
			direction = $3,
			entry_price_signal = $4,
			entry_time_signal = $5,
			entry_price_order = $6,
			entry_time_order = $7,
			entry_signal = $8,
			exit_price_signal = $9,
			exit_time_signal = $10,
			exit_price_order = $11,
			exit_time_order = $12,
			exit_signal = $13,
			exit_reason = $14,
			profit_loss_absolute = $15,
			profit_loss_percentage = $16,
			transaction_cost = $17,
			entry_fees = $18,
			exit_fees = $19,
			net_profit_loss = $20,
			position_size = $21,
			entry_order_id = $22,
			exit_order_id = $23,
			updated_at = $24
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		position.ID,
		position.State,
		position.Direction,
		position.EntryPriceSignal,
		position.EntryTimeSignal,
		position.EntryPriceOrder,
		position.EntryTimeOrder,
		position.EntrySignal,
		position.ExitPriceSignal,
		position.ExitTimeSignal,
		position.ExitPriceOrder,
		position.ExitTimeOrder,
		position.ExitSignal,
		position.ExitReason,
		position.ProfitLossAbsolute,
		position.ProfitLossPercentage,
		position.TransactionCost,
		position.EntryFees,
		position.ExitFees,
		position.NetProfitLoss,
		position.PositionSize,
		position.EntryOrderID,
		position.ExitOrderID,
		position.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrPositionNotFound
	}

	return nil
}

// Delete removes a position from the database
func (r *PositionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM position
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrPositionNotFound
	}

	return nil
}

// BulkCreate saves multiple positions in the database in a single operation
func (r *PositionRepository) BulkCreate(ctx context.Context, positions []*model.Position) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback()

	// Prepare the statement
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO position (
			id, run_id, state, direction, entry_price_signal, entry_time_signal,
			entry_price_order, entry_time_order, entry_signal, exit_price_signal,
			exit_time_signal, exit_price_order, exit_time_order, exit_signal,
			exit_reason, profit_loss_absolute, profit_loss_percentage, 
			transaction_cost, entry_fees, exit_fees, net_profit_loss,
			position_size, entry_order_id, exit_order_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
		)
	`)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer stmt.Close()

	// Execute for each position
	for _, position := range positions {
		_, err := stmt.ExecContext(
			ctx,
			position.ID,
			position.RunID,
			position.State,
			position.Direction,
			position.EntryPriceSignal,
			position.EntryTimeSignal,
			position.EntryPriceOrder,
			position.EntryTimeOrder,
			position.EntrySignal,
			position.ExitPriceSignal,
			position.ExitTimeSignal,
			position.ExitPriceOrder,
			position.ExitTimeOrder,
			position.ExitSignal,
			position.ExitReason,
			position.ProfitLossAbsolute,
			position.ProfitLossPercentage,
			position.TransactionCost,
			position.EntryFees,
			position.ExitFees,
			position.NetProfitLoss,
			position.PositionSize,
			position.EntryOrderID,
			position.ExitOrderID,
			position.CreatedAt,
			position.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}
