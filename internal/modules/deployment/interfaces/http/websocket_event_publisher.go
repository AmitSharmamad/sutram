package http

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"algo/internal/modules/deployment/application/dto"
	"algo/internal/modules/deployment/application/service"

	"github.com/gorilla/websocket"
)

// WebSocketConnection represents a WebSocket connection with user context
type WebSocketConnection struct {
	UserID     string
	Connection *websocket.Conn
	Send       chan []byte
}

// WebSocketEventPublisher implements the DeploymentEventPublisher interface using WebSocket
type WebSocketEventPublisher struct {
	// Map of user ID to WebSocket connection
	connections map[string]*WebSocketConnection
	// Mutex to protect the connections map
	mutex sync.RWMutex
	// Channel to handle connection registration
	register chan *WebSocketConnection
	// Channel to handle connection unregistration
	unregister chan *WebSocketConnection
	// Channel to broadcast events to all connections
	broadcast chan []byte
}

// NewWebSocketEventPublisher creates a new WebSocket-based event publisher
func NewWebSocketEventPublisher() service.DeploymentEventPublisher {
	publisher := &WebSocketEventPublisher{
		connections: make(map[string]*WebSocketConnection),
		register:    make(chan *WebSocketConnection),
		unregister:  make(chan *WebSocketConnection),
		broadcast:   make(chan []byte),
	}

	// Start the hub goroutine
	go publisher.run()

	return publisher
}

// run handles the main event loop for the WebSocket publisher
func (p *WebSocketEventPublisher) run() {
	for {
		select {
		case conn := <-p.register:
			p.mutex.Lock()
			p.connections[conn.UserID] = conn
			p.mutex.Unlock()
			log.Printf("User %s subscribed to deployment events", conn.UserID)

		case conn := <-p.unregister:
			p.mutex.Lock()
			if _, ok := p.connections[conn.UserID]; ok {
				delete(p.connections, conn.UserID)
				close(conn.Send)
				p.mutex.Unlock()
				log.Printf("User %s unsubscribed from deployment events", conn.UserID)
			} else {
				p.mutex.Unlock()
			}

		case message := <-p.broadcast:
			p.mutex.RLock()
			for userID, conn := range p.connections {
				select {
				case conn.Send <- message:
				default:
					// Client's send channel is full, remove the connection
					delete(p.connections, userID)
					close(conn.Send)
					log.Printf("Removed slow client: %s", userID)
				}
			}
			p.mutex.RUnlock()
		}
	}
}

// PublishEvent publishes a deployment event to all subscribed clients
func (p *WebSocketEventPublisher) PublishEvent(ctx context.Context, event *dto.DeploymentEvent) error {
	// Marshal the event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Send to broadcast channel (non-blocking)
	select {
	case p.broadcast <- data:
	default:
		log.Printf("Broadcast channel full, dropping event")
	}

	return nil
}

// PublishEventToUser publishes a deployment event to a specific user
func (p *WebSocketEventPublisher) PublishEventToUser(ctx context.Context, userID string, event *dto.DeploymentEvent) error {
	// Marshal the event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	p.mutex.RLock()
	conn, exists := p.connections[userID]
	p.mutex.RUnlock()

	if !exists {
		// User is not subscribed, silently ignore
		return nil
	}

	// Send to user's channel (non-blocking)
	select {
	case conn.Send <- data:
	default:
		log.Printf("User %s send channel full, dropping event", userID)
	}

	return nil
}

// SubscribeUser subscribes a user to deployment events
func (p *WebSocketEventPublisher) SubscribeUser(ctx context.Context, userID string, conn interface{}) error {
	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		log.Printf("Invalid connection type for user %s", userID)
		return nil
	}

	connection := &WebSocketConnection{
		UserID:     userID,
		Connection: wsConn,
		Send:       make(chan []byte, 256),
	}

	// Register the connection
	select {
	case p.register <- connection:
	default:
		log.Printf("Register channel full for user %s", userID)
	}

	// Start goroutines to handle this connection
	go p.writePump(connection)
	go p.readPump(connection)

	return nil
}

// UnsubscribeUser unsubscribes a user from deployment events
func (p *WebSocketEventPublisher) UnsubscribeUser(ctx context.Context, userID string) error {
	p.mutex.RLock()
	conn, exists := p.connections[userID]
	p.mutex.RUnlock()

	if exists {
		select {
		case p.unregister <- conn:
		default:
			log.Printf("Unregister channel full for user %s", userID)
		}
	}

	return nil
}

// GetSubscribedUsers returns a list of currently subscribed user IDs
func (p *WebSocketEventPublisher) GetSubscribedUsers() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	users := make([]string, 0, len(p.connections))
	for userID := range p.connections {
		users = append(users, userID)
	}

	return users
}

// writePump pumps messages from the hub to the websocket connection
func (p *WebSocketEventPublisher) writePump(conn *WebSocketConnection) {
	defer func() {
		conn.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-conn.Send:
			if !ok {
				// The hub closed the channel
				conn.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			conn.Connection.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// readPump pumps messages from the websocket connection to the hub
func (p *WebSocketEventPublisher) readPump(conn *WebSocketConnection) {
	defer func() {
		select {
		case p.unregister <- conn:
		default:
		}
		conn.Connection.Close()
	}()

	for {
		_, _, err := conn.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %s: %v", conn.UserID, err)
			}
			break
		}
		// For now, we don't handle incoming messages from clients
		// This could be extended to handle subscription preferences, etc.
	}
}
