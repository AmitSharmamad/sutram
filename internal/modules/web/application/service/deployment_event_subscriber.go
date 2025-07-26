package service

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"time"

	"algo/internal/modules/deployment/application/dto"

	"github.com/gorilla/websocket"
)

// DeploymentEventSubscriber subscribes to deployment events from the deployment server
type DeploymentEventSubscriber struct {
	deploymentServerURL string
	eventHandler        EventHandler
	conn                *websocket.Conn
	ctx                 context.Context
	cancel              context.CancelFunc
}

// EventHandler interface for handling deployment events
type EventHandler interface {
	HandleDeploymentEvent(event *dto.DeploymentEvent)
}

// NewDeploymentEventSubscriber creates a new deployment event subscriber
func NewDeploymentEventSubscriber(deploymentServerURL string, eventHandler EventHandler) *DeploymentEventSubscriber {
	ctx, cancel := context.WithCancel(context.Background())

	return &DeploymentEventSubscriber{
		deploymentServerURL: deploymentServerURL,
		eventHandler:        eventHandler,
		ctx:                 ctx,
		cancel:              cancel,
	}
}

// Start starts the subscription to deployment events
func (s *DeploymentEventSubscriber) Start() error {
	go s.connectAndListen()
	return nil
}

// Stop stops the subscription
func (s *DeploymentEventSubscriber) Stop() {
	s.cancel()
	if s.conn != nil {
		s.conn.Close()
	}
}

// connectAndListen establishes connection and listens for events
func (s *DeploymentEventSubscriber) connectAndListen() {
	for {
		select {
		case <-s.ctx.Done():
			log.Println("Deployment event subscriber stopped")
			return
		default:
			if err := s.connect(); err != nil {
				log.Printf("Failed to connect to deployment server: %v", err)
				// Wait before retrying
				select {
				case <-s.ctx.Done():
					return
				case <-time.After(30 * time.Second):
					continue
				}
			}
		}
	}
}

// connect establishes WebSocket connection to deployment server
func (s *DeploymentEventSubscriber) connect() error {
	// Parse and construct the WebSocket URL
	u, err := url.Parse(s.deploymentServerURL)
	if err != nil {
		return err
	}

	// Use internal WebSocket endpoint (no auth required for internal services)
	u.Path = "/api/deployment/events/internal/ws"
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	}

	log.Printf("Connecting to deployment server WebSocket: %s", u.String())

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	s.conn = conn
	log.Printf("Connected to deployment server WebSocket")

	// Set up ping/pong to keep connection alive
	conn.SetPongHandler(func(string) error {
		log.Printf("Received pong from deployment server")
		return nil
	})

	// Start ping ticker
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	// Channels for handling messages
	messageChan := make(chan []byte, 100)
	errorChan := make(chan error, 1)

	// Start message reader goroutine
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				errorChan <- err
				return
			}
			messageChan <- message
		}
	}()

	// Main event loop
	for {
		select {
		case <-s.ctx.Done():
			log.Println("Deployment event subscriber context cancelled")
			return nil

		case <-pingTicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				log.Printf("Failed to send ping to deployment server: %v", err)
				return err
			}

		case message := <-messageChan:
			if err := s.processMessage(message); err != nil {
				log.Printf("Failed to process message from deployment server: %v", err)
			}

		case err := <-errorChan:
			log.Printf("WebSocket read error from deployment server: %v", err)
			return err
		}
	}
}

// processMessage processes incoming deployment events
func (s *DeploymentEventSubscriber) processMessage(message []byte) error {
	var event dto.DeploymentEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return err
	}

	log.Printf("Received deployment event: %s for deployment %s",
		event.EventType, event.DeploymentID.String())

	// Forward the event to the handler
	s.eventHandler.HandleDeploymentEvent(&event)

	return nil
}

// IsConnected returns whether the subscriber is connected
func (s *DeploymentEventSubscriber) IsConnected() bool {
	return s.conn != nil
}
