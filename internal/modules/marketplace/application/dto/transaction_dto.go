// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/dto/transaction_dto.go
package dto

import (
	"time"

	"algo/internal/modules/marketplace/domain/model"

	"github.com/google/uuid"
)

// TransactionResponse represents a transaction response
type TransactionResponse struct {
	ID             uuid.UUID `json:"id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`
	PaymentMethod  string    `json:"payment_method"`
	PaymentID      string    `json:"payment_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// ToTransactionResponse converts a model.Transaction to a TransactionResponse
func ToTransactionResponse(transaction *model.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:             transaction.ID,
		SubscriptionID: transaction.SubscriptionID,
		Amount:         transaction.Amount,
		Status:         string(transaction.Status),
		PaymentMethod:  transaction.PaymentMethod,
		PaymentID:      transaction.PaymentID,
		CreatedAt:      transaction.CreatedAt,
	}
}
