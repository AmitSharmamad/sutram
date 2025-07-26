package http

import (
	"log"
	"net/http"

	"algo/internal/modules/deployment/application/dto"
	webService "algo/internal/modules/web/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WebEventHandler handles page-specific WebSocket connections for the web server
type WebEventHandler struct {
	hub         *PageSpecificWebSocketHub
	transformer *webService.DeploymentEventTransformer
}

// NewWebEventHandler creates a new web event handler
func NewWebEventHandler() *WebEventHandler {
	return &WebEventHandler{
		hub:         NewPageSpecificWebSocketHub(),
		transformer: webService.NewDeploymentEventTransformer(),
	}
}

// HandleWebSocket handles WebSocket upgrade for page-specific connections
func (h *WebEventHandler) HandleWebSocket(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authentication required",
		})
		return
	}

	// Convert userID to string
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

	// Upgrade the connection with page context
	conn, err := h.hub.UpgradeConnection(c.Writer, c.Request, userID)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upgrade connection",
		})
		return
	}

	log.Printf("User %s connected to %s page via WebSocket", userID, conn.PageContext.PageType)
}

// HandleDeploymentEvent processes deployment events from the deployment server
func (h *WebEventHandler) HandleDeploymentEvent(deploymentEvent *dto.DeploymentEvent) {
	// Transform deployment event into page-specific web events
	webEvents := h.transformer.TransformToWebEvents(deploymentEvent)

	// Publish each web event to the hub
	for _, webEvent := range webEvents {
		h.hub.PublishEvent(webEvent)
	}

	if len(webEvents) > 0 {
		log.Printf("Processed deployment event %s into %d web events",
			deploymentEvent.EventType, len(webEvents))
	}
}

// PublishWebEvent publishes a custom web event
func (h *WebEventHandler) PublishWebEvent(event *dto.WebEvent) {
	h.hub.PublishEvent(event)
}

// GetConnectionStats returns connection statistics
func (h *WebEventHandler) GetConnectionStats() gin.H {
	return gin.H{
		"active_connections": h.hub.GetActiveConnections(),
		"timestamp":          h.hub.GetActiveConnections(), // placeholder
	}
}

// UpdateConnectionContext updates the page context for a connection
func (h *WebEventHandler) UpdateConnectionContext(connectionID string, newContext dto.PageContext) bool {
	return h.hub.UpdateConnectionContext(connectionID, newContext)
}

// RegisterWebEventRoutes registers WebSocket routes for page-specific events
func RegisterWebEventRoutes(router *gin.RouterGroup, handler *WebEventHandler, authMiddleware gin.HandlerFunc) {
	// Main WebSocket endpoint for page-specific events
	router.GET("/events/ws", authMiddleware, handler.HandleWebSocket)

	// Status endpoint for monitoring
	router.GET("/events/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, handler.GetConnectionStats())
	})
}
