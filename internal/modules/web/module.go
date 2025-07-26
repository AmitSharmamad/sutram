package web

import (
	"fmt"
	"os"

	"algo/internal/infra/http/middleware"
	userHttp "algo/internal/modules/user/interfaces/http"
	"algo/internal/modules/web/application/service"
	"algo/internal/modules/web/interfaces/http"

	"github.com/gin-gonic/gin"
)

// Module represents the web events module
type Module struct {
	eventHandler              *http.WebEventHandler
	deploymentEventSubscriber *service.DeploymentEventSubscriber
	userController            *userHttp.UserController
}

// New creates a new web module instance
func New() *Module {
	return &Module{}
}

// Initialize initializes the web module with dependencies
func (m *Module) Initialize(userController *userHttp.UserController) error {
	m.userController = userController

	// Create event handler
	m.eventHandler = http.NewWebEventHandler()

	// Get deployment server URL from environment
	deploymentServerURL := getEnv("DEPLOYMENT_SERVER_URL", "http://localhost:9091")

	// Create deployment event subscriber
	m.deploymentEventSubscriber = service.NewDeploymentEventSubscriber(
		deploymentServerURL,
		m.eventHandler, // The event handler implements EventHandler interface
	)

	// Start subscribing to deployment events
	if err := m.deploymentEventSubscriber.Start(); err != nil {
		return fmt.Errorf("failed to start deployment event subscriber: %w", err)
	}

	return nil
}

// RegisterRoutes registers the web module routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	if m.eventHandler == nil {
		panic("Web module not initialized. Call Initialize() first.")
	}

	// Create auth middleware
	authServiceAdapter := userHttp.NewAuthServiceAdapter(m.userController.GetAuthService())
	authMiddleware := middleware.AuthMiddleware(authServiceAdapter)

	// Register WebSocket routes
	http.RegisterWebEventRoutes(router, m.eventHandler, authMiddleware)
}

// GetEventHandler returns the event handler for external use
func (m *Module) GetEventHandler() *http.WebEventHandler {
	return m.eventHandler
}

// Shutdown gracefully shuts down the web module
func (m *Module) Shutdown() {
	if m.deploymentEventSubscriber != nil {
		m.deploymentEventSubscriber.Stop()
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
