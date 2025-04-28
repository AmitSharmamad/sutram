// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/infrastructure/persistence/postgres/transaction_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// TransactionRepository implements repository.TransactionRepository
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	query := `
		INSERT INTO marketplace_transaction (
			id, subscription_id, amount, status, payment_method, payment_id, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		transaction.ID,
		transaction.SubscriptionID,
		transaction.Amount,
		transaction.Status,
		transaction.PaymentMethod,
		transaction.PaymentID,
		transaction.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// FindByID returns a transaction by ID
func (r *TransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	query := `
		SELECT 
			id, subscription_id, amount, status, payment_method, payment_id, created_at
		FROM marketplace_transaction
		WHERE id = $1
	`

	var transaction model.Transaction
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.SubscriptionID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.PaymentMethod,
		&transaction.PaymentID,
		&transaction.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	return &transaction, nil
}

// FindBySubscriptionID returns all transactions for a subscription
func (r *TransactionRepository) FindBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]*model.Transaction, error) {
	query := `
		SELECT 
			id, subscription_id, amount, status, payment_method, payment_id, created_at
		FROM marketplace_transaction
		WHERE subscription_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.SubscriptionID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.PaymentMethod,
			&transaction.PaymentID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// FindByUserID returns all transactions for a user
func (r *TransactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Transaction, error) {
	query := `
		SELECT 
			t.id, t.subscription_id, t.amount, t.status, t.payment_method, t.payment_id, t.created_at
		FROM marketplace_transaction t
		JOIN marketplace_subscription s ON t.subscription_id = s.id
		WHERE s.subscriber_id = $1
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.SubscriptionID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.PaymentMethod,
			&transaction.PaymentID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// FindBySellerID returns all transactions for a seller
func (r *TransactionRepository) FindBySellerID(ctx context.Context, sellerID uuid.UUID) ([]*model.Transaction, error) {
	query := `
		SELECT 
			t.id, t.subscription_id, t.amount, t.status, t.payment_method, t.payment_id, t.created_at
		FROM marketplace_transaction t
		JOIN marketplace_subscription s ON t.subscription_id = s.id
		JOIN marketplace_listing l ON s.listing_id = l.id
		WHERE l.user_id = $1
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.SubscriptionID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.PaymentMethod,
			&transaction.PaymentID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// Update updates a transaction
func (r *TransactionRepository) Update(ctx context.Context, transaction *model.Transaction) error {
	query := `
		UPDATE marketplace_transaction SET
			status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		transaction.Status,
		transaction.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

// GetRevenueBySeller gets the revenue generated by a seller
func (r *TransactionRepository) GetRevenueBySeller(ctx context.Context, sellerID uuid.UUID, startDate, endDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(t.amount), 0)
		FROM marketplace_transaction t
		JOIN marketplace_subscription s ON t.subscription_id = s.id
		JOIN marketplace_listing l ON s.listing_id = l.id
		WHERE l.user_id = $1
		AND t.status = 'completed'
		AND t.created_at BETWEEN $2 AND $3
	`

	var revenue float64
	err := r.db.QueryRowContext(ctx, query, sellerID, startDate, endDate).Scan(&revenue)
	if err != nil {
		return 0, fmt.Errorf("failed to get revenue: %w", err)
	}

	return revenue, nil
}
