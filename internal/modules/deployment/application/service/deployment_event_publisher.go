package service

import (
	"context"

	"algo/internal/modules/deployment/application/dto"
)

// DeploymentEventPublisher defines the interface for publishing deployment events
type DeploymentEventPublisher interface {
	// PublishEvent publishes a deployment event to all subscribed clients
	PublishEvent(ctx context.Context, event *dto.DeploymentEvent) error

	// PublishEventToUser publishes a deployment event to a specific user
	PublishEventToUser(ctx context.Context, userID string, event *dto.DeploymentEvent) error

	// SubscribeUser subscribes a user to deployment events
	SubscribeUser(ctx context.Context, userID string, conn interface{}) error

	// UnsubscribeUser unsubscribes a user from deployment events
	UnsubscribeUser(ctx context.Context, userID string) error

	// GetSubscribedUsers returns a list of currently subscribed user IDs
	GetSubscribedUsers() []string
}
