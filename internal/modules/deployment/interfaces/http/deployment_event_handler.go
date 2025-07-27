package http

import (
	"log"
	"net/http"

	"algo/internal/infra/http/middleware"
	"algo/internal/modules/deployment/application/service"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// DeploymentEventHandler handles WebSocket connections for deployment events
type DeploymentEventHandler struct {
	eventPublisher service.DeploymentEventPublisher
	upgrader       websocket.Upgrader
}

// NewDeploymentEventHandler creates a new deployment event handler
func NewDeploymentEventHandler(eventPublisher service.DeploymentEventPublisher) *DeploymentEventHandler {
	return &DeploymentEventHandler{
		eventPublisher: eventPublisher,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// In production, implement proper origin checking
				return true
			},
		},
	}
}

// HandleWebSocket handles WebSocket upgrade and connection for deployment events
func (h *DeploymentEventHandler) HandleWebSocket(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authentication required",
		})
		return
	}

	// Convert userID to string - it could be either uuid.UUID or string
	var userID string
	switch v := userIDInterface.(type) {
	case string:
		userID = v
	case uuid.UUID:
		userID = v.String()
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid user ID type",
		})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection for user %s: %v", userID, err)
		return
	}

	log.Printf("New deployment events WebSocket connection for user: %s", userID)

	// Subscribe the user to deployment events
	if err := h.eventPublisher.SubscribeUser(c.Request.Context(), userID, conn); err != nil {
		log.Printf("Failed to subscribe user %s to deployment events: %v", userID, err)
		conn.Close()
		return
	}

	// The connection handling is now managed by the event publisher
}

// HandleInternalWebSocket handles WebSocket upgrade for internal services without authentication
func (h *DeploymentEventHandler) HandleInternalWebSocket(c *gin.Context) {
	// Check if this is an internal service call (optional security check)
	clientIP := c.ClientIP()
	if clientIP != "127.0.0.1" && clientIP != "::1" && clientIP != "localhost" {
		log.Printf("Internal WebSocket endpoint accessed from non-localhost IP: %s", clientIP)
		// In production, you might want to be more strict about this
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade internal WebSocket connection: %v", err)
		return
	}

	log.Printf("New internal deployment events WebSocket connection from %s", clientIP)

	// For internal services, we use a special internal user ID and subscribe to all events
	internalUserID := "internal-service"
	if err := h.eventPublisher.SubscribeUser(c.Request.Context(), internalUserID, conn); err != nil {
		log.Printf("Failed to subscribe internal service to deployment events: %v", err)
		conn.Close()
		return
	}

	// The connection handling is now managed by the event publisher
}

// RegisterWebSocketRoutes registers WebSocket routes for deployment events
func RegisterWebSocketRoutes(router *gin.RouterGroup, eventHandler *DeploymentEventHandler, userController *userHttp.UserController) {
	// Create auth service adapter
	authServiceAdapter := userHttp.NewAuthServiceAdapter(userController.GetAuthService())

	// Create middleware
	authMiddleware := middleware.AuthMiddleware(authServiceAdapter)

	// WebSocket endpoint for deployment events (authenticated)
	router.GET("/deployment/events/ws", authMiddleware, eventHandler.HandleWebSocket)

	// Internal WebSocket endpoint for deployment events (no authentication required)
	router.GET("/deployment/events/internal/ws", eventHandler.HandleInternalWebSocket)
}
