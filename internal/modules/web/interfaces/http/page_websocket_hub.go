package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"algo/internal/modules/deployment/application/dto"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// PageSpecificWebSocketConnection represents a WebSocket connection with page context
type PageSpecificWebSocketConnection struct {
	ConnectionID string
	UserID       string
	Connection   *websocket.Conn
	PageContext  dto.PageContext
	Send         chan []byte
	ConnectedAt  time.Time
	LastActiveAt time.Time
	mutex        sync.RWMutex
}

// UpdatePageContext updates the page context for this connection
func (conn *PageSpecificWebSocketConnection) UpdatePageContext(newContext dto.PageContext) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.PageContext = newContext
	conn.LastActiveAt = time.Now()
}

// GetPageContext safely returns the current page context
func (conn *PageSpecificWebSocketConnection) GetPageContext() dto.PageContext {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	return conn.PageContext
}

// PageSpecificWebSocketHub manages page-specific WebSocket connections
type PageSpecificWebSocketHub struct {
	// Map of user ID to list of their connections (supporting multiple tabs/devices)
	userConnections map[string][]*PageSpecificWebSocketConnection

	// Map of connection ID to connection for quick lookup
	connections map[string]*PageSpecificWebSocketConnection

	// Channels for managing connections
	register   chan *PageSpecificWebSocketConnection
	unregister chan *PageSpecificWebSocketConnection

	// Channel for broadcasting events
	broadcast chan *dto.WebEvent

	// Mutex to protect the maps
	mutex sync.RWMutex

	// WebSocket upgrader
	upgrader websocket.Upgrader
}

// NewPageSpecificWebSocketHub creates a new page-specific WebSocket hub
func NewPageSpecificWebSocketHub() *PageSpecificWebSocketHub {
	hub := &PageSpecificWebSocketHub{
		userConnections: make(map[string][]*PageSpecificWebSocketConnection),
		connections:     make(map[string]*PageSpecificWebSocketConnection),
		register:        make(chan *PageSpecificWebSocketConnection),
		unregister:      make(chan *PageSpecificWebSocketConnection),
		broadcast:       make(chan *dto.WebEvent, 256),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// In production, implement proper origin checking
				return true
			},
		},
	}

	// Start the hub
	go hub.run()

	return hub
}

// run handles the main event loop for the WebSocket hub
func (h *PageSpecificWebSocketHub) run() {
	ticker := time.NewTicker(54 * time.Second) // Ping interval
	defer ticker.Stop()

	for {
		select {
		case conn := <-h.register:
			h.registerConnection(conn)

		case conn := <-h.unregister:
			h.unregisterConnection(conn)

		case event := <-h.broadcast:
			h.broadcastEvent(event)

		case <-ticker.C:
			h.pingConnections()
		}
	}
}

// registerConnection registers a new WebSocket connection
func (h *PageSpecificWebSocketHub) registerConnection(conn *PageSpecificWebSocketConnection) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Add to connections map
	h.connections[conn.ConnectionID] = conn

	// Add to user connections list
	if userConns, exists := h.userConnections[conn.UserID]; exists {
		h.userConnections[conn.UserID] = append(userConns, conn)
	} else {
		h.userConnections[conn.UserID] = []*PageSpecificWebSocketConnection{conn}
	}

	log.Printf("User %s connected to %s page (connection: %s)",
		conn.UserID, conn.PageContext.PageType, conn.ConnectionID)
}

// unregisterConnection unregisters a WebSocket connection
func (h *PageSpecificWebSocketHub) unregisterConnection(conn *PageSpecificWebSocketConnection) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Remove from connections map
	delete(h.connections, conn.ConnectionID)

	// Remove from user connections list
	if userConns, exists := h.userConnections[conn.UserID]; exists {
		for i, userConn := range userConns {
			if userConn.ConnectionID == conn.ConnectionID {
				// Remove this connection from the slice
				h.userConnections[conn.UserID] = append(userConns[:i], userConns[i+1:]...)
				break
			}
		}

		// If no more connections for this user, remove the user entry
		if len(h.userConnections[conn.UserID]) == 0 {
			delete(h.userConnections, conn.UserID)
		}
	}

	// Close the send channel and connection
	close(conn.Send)
	conn.Connection.Close()

	log.Printf("User %s disconnected from %s page (connection: %s)",
		conn.UserID, conn.PageContext.PageType, conn.ConnectionID)
}

// broadcastEvent sends an event to relevant connections based on page context
func (h *PageSpecificWebSocketHub) broadcastEvent(event *dto.WebEvent) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// Marshal the event once
	eventData, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	// Track how many connections received the event
	sentCount := 0

	// Send to all relevant connections
	for userID, userConns := range h.userConnections {
		// Check if this event is relevant for this user
		if !h.isEventRelevantForUser(event, userID) {
			continue
		}

		for _, conn := range userConns {
			// Check if this event is relevant for this page context
			if !h.isEventRelevantForPageContext(event, conn.GetPageContext()) {
				continue
			}

			// Try to send the event (non-blocking)
			select {
			case conn.Send <- eventData:
				sentCount++
			default:
				// Connection send channel is full, log and remove connection
				log.Printf("Send channel full for connection %s, removing", conn.ConnectionID)
				go func(c *PageSpecificWebSocketConnection) {
					h.unregister <- c
				}(conn)
			}
		}
	}

	if sentCount > 0 {
		log.Printf("Broadcasted %s event to %d connections", event.EventType, sentCount)
	}
}

// isEventRelevantForUser checks if an event is relevant for a specific user
func (h *PageSpecificWebSocketHub) isEventRelevantForUser(event *dto.WebEvent, userID string) bool {
	// Check if event has user-specific metadata
	if eventUserID, exists := event.Metadata["user_id"]; exists {
		if eventUserIDStr, ok := eventUserID.(string); ok {
			return eventUserIDStr == userID
		}
	}

	// For system-wide events (like market status), send to all users
	if event.EventType == "system.market_status" || event.EventType == "system.maintenance" {
		return true
	}

	return false
}

// isEventRelevantForPageContext checks if an event is relevant for a specific page context
func (h *PageSpecificWebSocketHub) isEventRelevantForPageContext(event *dto.WebEvent, pageContext dto.PageContext) bool {
	// If event is page-specific, check if it matches
	if event.PageType != "" && event.PageType != pageContext.PageType {
		return false
	}

	// For deployment detail pages, check if deployment ID matches
	if pageContext.PageType == dto.PageTypeDeploymentDetail {
		if deploymentID, exists := pageContext.PageParams["deployment_id"]; exists {
			if eventDeploymentID, exists := event.Metadata["deployment_id"]; exists {
				return deploymentID == eventDeploymentID
			}
		}
	}

	// High priority events go to all pages
	if event.Priority == dto.EventPriorityHigh {
		return true
	}

	return true
}

// pingConnections sends ping messages to all connections to keep them alive
func (h *PageSpecificWebSocketHub) pingConnections() {
	h.mutex.RLock()
	connections := make([]*PageSpecificWebSocketConnection, 0, len(h.connections))
	for _, conn := range h.connections {
		connections = append(connections, conn)
	}
	h.mutex.RUnlock()

	for _, conn := range connections {
		if err := conn.Connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			log.Printf("Error sending ping to connection %s: %v", conn.ConnectionID, err)
			h.unregister <- conn
		}
	}
}

// PublishEvent publishes an event to the hub for broadcasting
func (h *PageSpecificWebSocketHub) PublishEvent(event *dto.WebEvent) {
	select {
	case h.broadcast <- event:
	default:
		log.Printf("Broadcast channel full, dropping event: %s", event.EventType)
	}
}

// UpdateConnectionContext updates the page context for a specific connection
func (h *PageSpecificWebSocketHub) UpdateConnectionContext(connectionID string, newContext dto.PageContext) bool {
	h.mutex.RLock()
	conn, exists := h.connections[connectionID]
	h.mutex.RUnlock()

	if !exists {
		return false
	}

	conn.UpdatePageContext(newContext)
	log.Printf("Updated context for connection %s to %s page", connectionID, newContext.PageType)
	return true
}

// GetActiveConnections returns the number of active connections
func (h *PageSpecificWebSocketHub) GetActiveConnections() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.connections)
}

// GetUserConnectionCount returns the number of connections for a specific user
func (h *PageSpecificWebSocketHub) GetUserConnectionCount(userID string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if userConns, exists := h.userConnections[userID]; exists {
		return len(userConns)
	}
	return 0
}

// UpgradeConnection upgrades an HTTP connection to WebSocket with page context
func (h *PageSpecificWebSocketHub) UpgradeConnection(w http.ResponseWriter, r *http.Request, userID string) (*PageSpecificWebSocketConnection, error) {
	// Parse page context from query parameters
	pageContext, err := h.parsePageContextFromQuery(r.URL.Query())
	if err != nil {
		return nil, err
	}

	// Upgrade the connection
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	// Create connection wrapper
	wsConn := &PageSpecificWebSocketConnection{
		ConnectionID: uuid.New().String(),
		UserID:       userID,
		Connection:   conn,
		PageContext:  pageContext,
		Send:         make(chan []byte, 256),
		ConnectedAt:  time.Now(),
		LastActiveAt: time.Now(),
	}

	// Register the connection
	h.register <- wsConn

	// Start read and write pumps
	go h.writePump(wsConn)
	go h.readPump(wsConn)

	return wsConn, nil
}

// parsePageContextFromQuery parses page context from URL query parameters
func (h *PageSpecificWebSocketHub) parsePageContextFromQuery(values url.Values) (dto.PageContext, error) {
	pageType := dto.PageType(values.Get("page"))
	if pageType == "" {
		pageType = dto.PageTypeDashboard // Default to dashboard
	}

	pageParams := make(map[string]interface{})

	// Parse page-specific parameters
	switch pageType {
	case dto.PageTypeDeploymentDetail:
		if deploymentID := values.Get("deployment_id"); deploymentID != "" {
			pageParams["deployment_id"] = deploymentID
		}
	}

	return dto.PageContext{
		PageType:   pageType,
		PageParams: pageParams,
	}, nil
}

// writePump pumps messages from the hub to the websocket connection
func (h *PageSpecificWebSocketHub) writePump(conn *PageSpecificWebSocketConnection) {
	defer func() {
		h.unregister <- conn
	}()

	for message := range conn.Send {
		conn.Connection.SetWriteDeadline(time.Now().Add(10 * time.Second))

		if err := conn.Connection.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message to connection %s: %v", conn.ConnectionID, err)
			return
		}
	}
}

// readPump pumps messages from the websocket connection to the hub
func (h *PageSpecificWebSocketHub) readPump(conn *PageSpecificWebSocketConnection) {
	defer func() {
		h.unregister <- conn
	}()

	conn.Connection.SetReadLimit(512)
	conn.Connection.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.Connection.SetPongHandler(func(string) error {
		conn.Connection.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := conn.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for connection %s: %v", conn.ConnectionID, err)
			}
			break
		}

		// Handle incoming messages (context updates, subscriptions, etc.)
		h.handleIncomingMessage(conn, message)
	}
}

// handleIncomingMessage handles messages received from the client
func (h *PageSpecificWebSocketHub) handleIncomingMessage(conn *PageSpecificWebSocketConnection, message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling message from connection %s: %v", conn.ConnectionID, err)
		return
	}

	// Handle context update messages
	if msgType, exists := msg["type"]; exists && msgType == "update_context" {
		var contextUpdate dto.ContextUpdateMessage
		if err := json.Unmarshal(message, &contextUpdate); err != nil {
			log.Printf("Error unmarshaling context update: %v", err)
			return
		}

		conn.UpdatePageContext(contextUpdate.PageContext)
		log.Printf("Updated page context for connection %s to %s",
			conn.ConnectionID, contextUpdate.PageContext.PageType)
	}
}
