// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/domain/model/transaction.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// TransactionStatus defines the status of a transaction
type TransactionStatus string

const (
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusRefunded  TransactionStatus = "refunded"
)

// Transaction represents a marketplace transaction
type Transaction struct {
	ID             uuid.UUID         `json:"id" db:"id"`
	SubscriptionID uuid.UUID         `json:"subscription_id" db:"subscription_id"`
	Amount         float64           `json:"amount" db:"amount"`
	Status         TransactionStatus `json:"status" db:"status"`
	PaymentMethod  string            `json:"payment_method" db:"payment_method"`
	PaymentID      string            `json:"payment_id,omitempty" db:"payment_id"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
}
