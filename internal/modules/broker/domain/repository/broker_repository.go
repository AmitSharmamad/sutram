package repository

import (
	"context"

	"algo/internal/modules/broker/domain/model"

	"github.com/google/uuid"
)

// BrokerRepository defines the interface for broker data operations
type BrokerRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Broker, error)
	FindByName(ctx context.Context, name string) (*model.Broker, error)
	Create(ctx context.Context, broker *model.Broker) error
	Update(ctx context.Context, broker *model.Broker) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*model.Broker, error)
}

// BrokerUserRepository defines the interface for broker user data operations
type BrokerUserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.BrokerUser, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.BrokerUser, error)
	FindByUserIDAndBrokerID(ctx context.Context, userID, brokerID uuid.UUID) (*model.BrokerUser, error)
	Create(ctx context.Context, brokerUser *model.BrokerUser) error
	Update(ctx context.Context, brokerUser *model.BrokerUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*model.BrokerUser, error)
}

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error)
	FindByPositionID(ctx context.Context, positionID uuid.UUID) ([]*model.Order, error)
	FindByBrokerOrderID(ctx context.Context, brokerOrderID string) (*model.Order, error)
	Create(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*model.Order, error)
	ListPaginated(ctx context.Context, offset, limit int) ([]*model.Order, int64, error)
}