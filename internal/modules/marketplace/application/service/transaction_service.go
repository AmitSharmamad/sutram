// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/service/transaction_service.go
package service

import (
	"context"
	"errors"
	"fmt"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// Additional errors
var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

// TransactionService handles marketplace transaction business logic
type TransactionService struct {
	transactionRepo repository.TransactionRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(transactionRepo repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

// GetTransactionsByUserID gets transactions for a user
func (s *TransactionService) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Transaction, error) {
	transactions, err := s.transactionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

// GetTransactionByID gets a transaction by ID
func (s *TransactionService) GetTransactionByID(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) (*model.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}
	if transaction == nil {
		return nil, ErrTransactionNotFound
	}

	// TODO: Check if user is authorized to view this transaction
	// This would require a join query to check if the transaction is associated
	// with the user's subscription or if the user is the seller

	return transaction, nil
}
