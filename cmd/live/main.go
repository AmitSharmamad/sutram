package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var (
	// Map to store active WebSocket connections
	clients = make(map[*websocket.Conn]bool)
)

// WebSocket upgrader configuration
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connections from any origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket client configuration
var wsDialer = websocket.DefaultDialer

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Set Gin mode
	gin.SetMode(getEnv("GIN_MODE", "debug"))

	// Initialize router
	router := gin.Default()

	// Setup CORS
	setupCORS(router)

	// Setup middleware
	setupMiddleware(router)

	// Setup WebSocket routes
	wsRouter := router.Group("/ws")
	setupWebSocketRoutes(wsRouter)

	// Setup API routes for live data
	apiRouter := router.Group("/api/live")
	setupAPIRoutes(apiRouter)

	// Get port from environment or use default
	port := getEnv("WS_PORT", "9091")
	serverAddr := fmt.Sprintf(":%s", port)

	// Create a server with graceful shutdown
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Connect to upstream WebSocket server
	go connectToUpstreamWebSocket()

	// Start the server in a goroutine
	go func() {
		log.Printf("WebSocket server started on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down WebSocket server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("WebSocket server exited properly")
}

// setupCORS configures CORS middleware
func setupCORS(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

// setupMiddleware adds global middleware to the router
func setupMiddleware(router *gin.Engine) {
	// Add request ID middleware
	router.Use(requestIDMiddleware())

	// Add logging middleware
	router.Use(loggingMiddleware())

	// Add recovery middleware (already included in gin.Default())
}

// requestIDMiddleware generates a unique request ID for each request
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-ID", generateRequestID())
		c.Next()
	}
}

// generateRequestID creates a unique ID for request tracing
func generateRequestID() string {
	// Simple implementation - in production, use a more robust method
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// loggingMiddleware logs request information
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		log.Printf(
			"[%s] %s %s %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			latency,
		)
	}
}

// setupWebSocketRoutes sets up routes for WebSocket connections
func setupWebSocketRoutes(router *gin.RouterGroup) {
	// Endpoint for client WebSocket connections
	router.GET("/feed", handleClientWebSocket)
}

// setupAPIRoutes sets up REST API routes for the live data module
func setupAPIRoutes(router *gin.RouterGroup) {
	// Endpoint for status check
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "live",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}

// handleClientWebSocket handles WebSocket connections from clients
func handleClientWebSocket(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Add connection to clients map
	clients[conn] = true
	defer delete(clients, conn)

	log.Printf("New client connected: %s", conn.RemoteAddr())

	// Listen for client messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message from client: %s", string(message))
		// Process client messages if needed
	}
}

// connectToUpstreamWebSocket establishes connection to the upstream WebSocket server
func connectToUpstreamWebSocket() {
	// Get WebSocket URL and authentication token from environment
	wsURL := getEnv("WS_UPSTREAM_URL", "wss://api-feed.dhan.co")
	wsToken := getEnv("WS_TOKEN", "")
	clientID := getEnv("WS_CLIENT_ID", "")

	// Construct full WebSocket URL with parameters
	fullWSURL := fmt.Sprintf("%s?version=2&token=%s&clientId=%s&authType=2", wsURL, wsToken, clientID)

	log.Printf("Connecting to upstream WebSocket: %s", wsURL)

	// Connect to the upstream WebSocket server
	conn, _, err := wsDialer.Dial(fullWSURL, nil)
	if err != nil {
		log.Printf("Error connecting to upstream WebSocket: %v", err)
		// Retry connection after a delay
		time.Sleep(5 * time.Second)
		go connectToUpstreamWebSocket()
		return
	}
	defer conn.Close()

	log.Printf("Connected to upstream WebSocket")

	// Listen for messages from the upstream WebSocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading from upstream WebSocket: %v", err)
			// Reconnect after a brief delay
			time.Sleep(5 * time.Second)
			go connectToUpstreamWebSocket()
			return
		}

		// Process and broadcast message to all connected clients
		broadcastMessage(message)
	}
}

// broadcastMessage sends a message to all connected clients
func broadcastMessage(message []byte) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
