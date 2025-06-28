package main

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"algo/internal/modules/instrument/domain/model"
	"algo/internal/modules/instrument/domain/repository"
	instrumentRepo "algo/internal/modules/instrument/infrastructure/persistence/clickhouse"

	// Deployment processing
	deploymentModule "algo/internal/modules/deployment"
	deploymentService "algo/internal/modules/deployment/application/service"
	strategyModule "algo/internal/modules/strategy"
	userModule "algo/internal/modules/user"
)

// JWT payload structure for basic validation
type JWTPayload struct {
	Exp          int64  `json:"exp"`
	Iss          string `json:"iss"`
	DhanClientId string `json:"dhanClientId"`
}

// validateJWTToken performs basic JWT token validation
func validateJWTToken(token string) error {
	// Split the JWT token
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid JWT token format")
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("failed to decode JWT payload: %v", err)
	}

	// Parse the payload
	var jwtPayload JWTPayload
	if err := json.Unmarshal(payload, &jwtPayload); err != nil {
		return fmt.Errorf("failed to parse JWT payload: %v", err)
	}

	// Check if token is expired
	now := time.Now().Unix()
	if jwtPayload.Exp < now {
		return fmt.Errorf("JWT token has expired. Expiry: %d, Current: %d", jwtPayload.Exp, now)
	}

	// Check issuer
	if jwtPayload.Iss != "dhan" {
		return fmt.Errorf("invalid JWT issuer: %s", jwtPayload.Iss)
	}

	log.Printf("JWT token is valid. Expires at: %s", time.Unix(jwtPayload.Exp, 0).Format(time.RFC3339))
	log.Printf("Dhan Client ID in token: %s", jwtPayload.DhanClientId)

	return nil
}

var (
	// Map to store active WebSocket connections
	clients = make(map[*websocket.Conn]bool)
	// Mutex to protect clients map
	clientsMutex sync.Mutex
	// Map to store subscribed instruments
	subscribedInstruments = make(map[string]bool)
	// ClickHouse database connection
	chDB *sqlx.DB
	// Instrument repository
	instrumentRepository repository.InstrumentRepository
	// Deployment service for processing deployments (TODO: initialize when needed)
	deploymentSvc deploymentService.DeploymentService
	// Application context and cancel function
	appCtx    context.Context
	cancelCtx context.CancelFunc
	// IST timezone location for timestamp conversion
	istLocation *time.Location
)

// Constants for request and feed codes
const (
	RequestCodeSubscribe   = 15
	RequestCodeUnsubscribe = 16
	RequestCodeDisconnect  = 12

	FeedCodeTicker     = 2
	FeedCodeQuote      = 4
	FeedCodeOI         = 5
	FeedCodePrevClose  = 6
	FeedCodeFull       = 8
	FeedCodeDisconnect = 50
)

// InstrumentRequest represents the structure for subscribing to instruments
type InstrumentRequest struct {
	RequestCode     int              `json:"RequestCode"`
	InstrumentCount int              `json:"InstrumentCount"`
	InstrumentList  []InstrumentInfo `json:"InstrumentList"`
}

// InstrumentInfo represents information about a single instrument
type InstrumentInfo struct {
	ExchangeSegment string `json:"ExchangeSegment"`
	SecurityId      string `json:"SecurityId"`
}

// DisconnectRequest represents the structure for disconnecting from the WebSocket
type DisconnectRequest struct {
	RequestCode int `json:"RequestCode"`
}

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

// Upstream connection for market data
var upstreamConn *websocket.Conn

func main() {
	// Create application context with cancel function
	appCtx, cancelCtx = context.WithCancel(context.Background())
	defer cancelCtx()

	// Initialize IST timezone location
	var err error
	istLocation, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatalf("Failed to load IST timezone: %v", err)
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	if err := initDB(); err != nil {
		log.Printf("Error initializing database: %v", err)
		// Continue without database connection
	}

	// Initialize deployment service for automatic subscription
	if err := initDeploymentService(); err != nil {
		log.Printf("Error initializing deployment service: %v", err)
		// Continue without deployment service
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
	port := getEnv("WS_PORT", "9099")
	serverAddr := fmt.Sprintf(":%s", port)

	// Create a server with graceful shutdown
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Connect to upstream WebSocket server
	go connectToUpstreamWebSocket()

	// Start periodic subscription refresh for deployment instruments
	go startPeriodicSubscriptionRefresh()

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

	// Cancel the application context to signal shutdown
	cancelCtx()

	// Send disconnect request to upstream server if connected
	if upstreamConn != nil {
		disconnectRequest := DisconnectRequest{
			RequestCode: RequestCodeDisconnect,
		}
		if err := upstreamConn.WriteJSON(disconnectRequest); err != nil {
			log.Printf("Error sending disconnect request: %v", err)
		}
	}

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("WebSocket server exited properly")
}

// initDB initializes the ClickHouse database connection
func initDB() error {
	// Get database connection details from environment variables

	// Create DSN for ClickHouse
	conn := "clickhouse://admin:admin@localhost:9000"
	db, err := sqlx.Connect("clickhouse", conn)

	if err != nil {
		return fmt.Errorf("error connecting to ClickHouse: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging ClickHouse: %w", err)
	}

	chDB = db
	log.Println("Connected to ClickHouse database")

	// Test table schema
	var tableExists int
	err = db.QueryRow("SELECT 1 FROM system.tables WHERE database = 'sutram' AND name = 'instrument'").Scan(&tableExists)
	if err != nil {
		log.Printf("Warning: Could not verify instrument table exists: %v", err)
	} else {
		log.Println("Verified instrument table exists in ClickHouse")

		// Test table structure
		rows, err := db.Query("DESCRIBE sutram.instrument")
		if err != nil {
			log.Printf("Warning: Could not describe instrument table: %v", err)
		} else {
			log.Println("Instrument table schema:")
			for rows.Next() {
				var name, typeStr, defaultValue, comment string
				if err := rows.Scan(&name, &typeStr, &defaultValue, &comment); err == nil {
					log.Printf("  %s: %s", name, typeStr)
				}
			}
			rows.Close()
		}
	}

	// Initialize instrument repository
	instrumentRepository = instrumentRepo.NewInstrumentRepository(chDB)

	return nil
}

// initDeploymentService initializes the deployment service for automatic subscription
func initDeploymentService() error {
	// Initialize user and strategy modules
	strategyMod := strategyModule.New()
	userMod := userModule.New()
	userController := userMod.GetController()

	// Initialize deployment module
	deploymentMod := deploymentModule.New()
	deploymentMod.Initialize(userController, strategyMod)
	deploymentSvc = deploymentMod.GetDeploymentService()

	// Initialize the deployment service
	ctx := context.Background()
	if err := deploymentSvc.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize deployment service: %w", err)
	}

	log.Println("Deployment service initialized successfully")
	return nil
}

// subscribeToActiveDeploymentInstruments gets all active deployments and subscribes to their instruments
func subscribeToActiveDeploymentInstruments() error {
	if deploymentSvc == nil {
		log.Printf("Deployment service not initialized, skipping automatic subscription")
		return nil
	}

	ctx := context.Background()

	// Get all active deployments with their strategies
	activeDeployments, err := deploymentSvc.GetActiveDeploymentsWithStrategies(ctx)
	if err != nil {
		log.Printf("Error getting active deployments: %v", err)
		return err
	}

	if len(activeDeployments) == 0 {
		log.Printf("No active deployments found, no instruments to subscribe to")
		return nil
	}

	// Extract unique instrument IDs from deployments
	instrumentIDsMap := make(map[int]bool)
	for _, deployment := range activeDeployments {
		instrumentIDsMap[deployment.Deployment.InstrumentID] = true
	}

	// Convert to slice
	instrumentIDs := make([]int, 0, len(instrumentIDsMap))
	for id := range instrumentIDsMap {
		instrumentIDs = append(instrumentIDs, id)
	}

	// Convert instrument IDs to InstrumentInfo format required by Dhan API
	instruments := make([]InstrumentInfo, len(instrumentIDs))
	for i, id := range instrumentIDs {
		instruments[i] = InstrumentInfo{
			ExchangeSegment: "NSE_EQ", // Default to NSE Equity, you may want to make this configurable
			SecurityId:      fmt.Sprintf("%d", id),
		}
	}

	log.Printf("Found %d active deployments with %d unique instruments to subscribe to",
		len(activeDeployments), len(instruments))

	// Subscribe to the instruments
	if len(instruments) > 0 {
		success := subscribeToInstruments(instruments)
		if !success {
			return fmt.Errorf("failed to subscribe to deployment instruments")
		}
		log.Printf("Successfully subscribed to %d instruments from active deployments", len(instruments))
	}

	return nil
}

// getInstrumentIDsFromDeployments extracts unique instrument IDs from active deployments
func getInstrumentIDsFromDeployments() []int {
	if deploymentSvc == nil {
		return []int{}
	}

	ctx := context.Background()

	// Get all active deployments with their strategies
	activeDeployments, err := deploymentSvc.GetActiveDeploymentsWithStrategies(ctx)
	if err != nil {
		log.Printf("Error getting active deployments: %v", err)
		return []int{}
	}

	// Extract unique instrument IDs
	instrumentIDsMap := make(map[int]bool)
	for _, deployment := range activeDeployments {
		instrumentIDsMap[deployment.Deployment.InstrumentID] = true
	}

	// Convert to slice
	instrumentIDs := make([]int, 0, len(instrumentIDsMap))
	for id := range instrumentIDsMap {
		instrumentIDs = append(instrumentIDs, id)
	}

	return instrumentIDs
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

	// Endpoint for instrument subscription
	router.POST("/subscribe", handleSubscribeInstruments)
}

// setupAPIRoutes sets up REST API routes for the live data module
func setupAPIRoutes(router *gin.RouterGroup) {
	// Endpoint for status check
	router.GET("/status", func(c *gin.Context) {
		subscribedCount := len(subscribedInstruments)
		deploymentInstruments := getInstrumentIDsFromDeployments()

		c.JSON(http.StatusOK, gin.H{
			"status": "live",
			"time":   time.Now().Format(time.RFC3339),
			"websocket_endpoints": gin.H{
				"market_data":       "/ws/feed",
				"deployment_events": "/ws/api/deployment/events/ws",
			},
			"subscriptions": gin.H{
				"total_subscribed_instruments": subscribedCount,
				"deployment_instruments":       len(deploymentInstruments),
				"deployment_instrument_ids":    deploymentInstruments,
				"upstream_connected":           upstreamConn != nil,
			},
			"deployment_service": gin.H{
				"initialized": deploymentSvc != nil,
			},
		})
	})

	// Endpoint to add instrument subscriptions
	router.POST("/subscribe", func(c *gin.Context) {
		var req struct {
			Instruments []InstrumentInfo `json:"instruments"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		if len(req.Instruments) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No instruments provided"})
			return
		}

		// Subscribe to instruments in chunks of 100 as per API limitation
		success := subscribeToInstruments(req.Instruments)

		if success {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"count":  len(req.Instruments),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to subscribe to instruments",
			})
		}
	})

	// Endpoint to manually refresh deployment subscriptions
	router.POST("/refresh-subscriptions", func(c *gin.Context) {
		if err := subscribeToActiveDeploymentInstruments(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to refresh subscriptions",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Subscriptions refreshed successfully",
			"count":   len(getInstrumentIDsFromDeployments()),
		})
	})
}

// handleSubscribeInstruments handles WebSocket requests for instrument subscription
func handleSubscribeInstruments(c *gin.Context) {
	var req struct {
		Instruments []InstrumentInfo `json:"instruments"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(req.Instruments) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No instruments provided"})
		return
	}

	success := subscribeToInstruments(req.Instruments)

	if success {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"count":  len(req.Instruments),
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to subscribe to instruments",
		})
	}
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

	// Add connection to clients map with mutex lock
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
	}()

	log.Printf("New client connected: %s", conn.RemoteAddr())

	// Listen for client messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message from client: %s", string(message))

		// Check if it's a subscription request
		if messageType == websocket.TextMessage {
			var req InstrumentRequest
			if err := json.Unmarshal(message, &req); err == nil {
				if req.RequestCode == RequestCodeSubscribe && len(req.InstrumentList) > 0 {
					// Forward subscription request to upstream WebSocket
					if upstreamConn != nil {
						if err := upstreamConn.WriteMessage(websocket.TextMessage, message); err != nil {
							log.Printf("Error forwarding subscription to upstream: %v", err)
						}
					}
				}
			}
		}
	}
}

// connectToUpstreamWebSocket establishes connection to the upstream WebSocket server
func connectToUpstreamWebSocket() {
	// Get WebSocket URL and authentication token from environment
	wsURL := getEnv("WS_UPSTREAM_URL", "wss://api-feed.dhan.co")
	wsToken := getEnv("WS_TOKEN", "")
	clientID := getEnv("WS_CLIENT_ID", "")

	// Validate required parameters
	if wsToken == "" {
		log.Printf("Error: WS_TOKEN is not set in environment variables")
		time.Sleep(10 * time.Second)
		go connectToUpstreamWebSocket()
		return
	}

	if clientID == "" {
		log.Printf("Error: WS_CLIENT_ID is not set in environment variables")
		time.Sleep(10 * time.Second)
		go connectToUpstreamWebSocket()
		return
	}

	// Validate JWT token
	if err := validateJWTToken(wsToken); err != nil {
		log.Printf("JWT Token validation failed: %v", err)
		log.Printf("Please update your WS_TOKEN with a valid Dhan access token")
		time.Sleep(30 * time.Second) // Wait longer for token issues
		go connectToUpstreamWebSocket()
		return
	}

	// Construct full WebSocket URL with parameters according to Dhan documentation:
	// wss://api-feed.dhan.co?version=2&token=eyxxxxx&clientId=100xxxxxxx&authType=2
	fullWSURL := fmt.Sprintf("%s?version=2&token=%s&clientId=%s&authType=2", wsURL, wsToken, clientID)

	log.Printf("Connecting to upstream WebSocket: %s", wsURL)
	log.Printf("Using Client ID: %s", clientID)
	log.Printf("Token length: %d characters", len(wsToken))
	log.Printf("Full WebSocket URL: %s?version=2&token=%s...&clientId=%s&authType=2",
		wsURL, wsToken[:20], clientID) // Only show first 20 chars of token for security

	// Create request headers for better debugging
	headers := make(map[string][]string)
	headers["User-Agent"] = []string{"Sutram-Live-Feed/1.0"}

	// Connect to the upstream WebSocket server
	conn, resp, err := wsDialer.Dial(fullWSURL, headers)
	if err != nil {
		log.Printf("Error connecting to upstream WebSocket: %v", err)
		if resp != nil {
			log.Printf("HTTP Response Status: %s", resp.Status)
			log.Printf("HTTP Response Headers: %v", resp.Header)
			// Try to read response body for more details
			if resp.Body != nil {
				body := make([]byte, 1024)
				n, _ := resp.Body.Read(body)
				if n > 0 {
					log.Printf("Response Body: %s", string(body[:n]))
				}
				resp.Body.Close()
			}
		}
		// Retry connection after a delay
		time.Sleep(5 * time.Second)
		go connectToUpstreamWebSocket()
		return
	}

	// Store the connection globally
	upstreamConn = conn

	defer func() {
		conn.Close()
		upstreamConn = nil
		// Reconnect after a brief delay
		time.Sleep(5 * time.Second)
		go connectToUpstreamWebSocket()
	}()

	log.Printf("Connected to upstream WebSocket")

	// Set up ping handler to keep connection alive
	conn.SetPingHandler(func(data string) error {
		err := conn.WriteMessage(websocket.PongMessage, []byte{})
		if err != nil {
			log.Printf("Error sending pong: %v", err)
			return err
		}
		return nil
	})

	// Re-subscribe to any existing instruments after reconnection
	resubscribeToInstruments()

	// Subscribe to instruments from active deployments
	if err := subscribeToActiveDeploymentInstruments(); err != nil {
		log.Printf("Error subscribing to deployment instruments: %v", err)
	}

	// Listen for messages from the upstream WebSocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading from upstream WebSocket: %v", err)
			return
		}

		// Process binary message from upstream server
		processMarketDataMessage(message)
	}
}

// startPeriodicSubscriptionRefresh periodically checks for new deployments and subscribes to their instruments
func startPeriodicSubscriptionRefresh() {
	ticker := time.NewTicker(5 * time.Minute) // Check every 5 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("Performing periodic subscription refresh...")
			if err := subscribeToActiveDeploymentInstruments(); err != nil {
				log.Printf("Error during periodic subscription refresh: %v", err)
			}
		case <-appCtx.Done():
			log.Printf("Stopping periodic subscription refresh")
			return
		}
	}
}

// subscribeToInstruments sends subscription requests to the upstream WebSocket
// Handles chunking of instruments in batches of 100 as per API limitation
func subscribeToInstruments(instruments []InstrumentInfo) bool {
	if upstreamConn == nil {
		log.Println("Cannot subscribe - not connected to upstream WebSocket")
		return false
	}

	// Process instruments in chunks of 100
	for i := 0; i < len(instruments); i += 100 {
		end := i + 100
		if end > len(instruments) {
			end = len(instruments)
		}

		chunk := instruments[i:end]

		// Create subscription request
		req := InstrumentRequest{
			RequestCode:     RequestCodeSubscribe,
			InstrumentCount: len(chunk),
			InstrumentList:  chunk,
		}

		// Send subscription request
		if err := upstreamConn.WriteJSON(req); err != nil {
			log.Printf("Error subscribing to instruments: %v", err)
			return false
		}

		// Store subscribed instruments in the map for reconnection management
		for _, inst := range chunk {
			key := fmt.Sprintf("%s:%s", inst.ExchangeSegment, inst.SecurityId)
			subscribedInstruments[key] = true
		}

		log.Printf("Subscribed to %d instruments", len(chunk))
	}

	return true
}

// unsubscribeFromInstruments sends unsubscribe requests to the upstream WebSocket
func unsubscribeFromInstruments(instruments []InstrumentInfo) bool {
	if upstreamConn == nil {
		return false
	}

	// Process instruments in chunks of 100
	for i := 0; i < len(instruments); i += 100 {
		end := i + 100
		if end > len(instruments) {
			end = len(instruments)
		}

		chunk := instruments[i:end]

		// Create unsubscription request
		req := InstrumentRequest{
			RequestCode:     RequestCodeUnsubscribe,
			InstrumentCount: len(chunk),
			InstrumentList:  chunk,
		}

		// Send unsubscription request
		if err := upstreamConn.WriteJSON(req); err != nil {
			log.Printf("Error unsubscribing from instruments: %v", err)
			return false
		}

		// Remove from the subscribed instruments map
		for _, inst := range chunk {
			key := fmt.Sprintf("%s:%s", inst.ExchangeSegment, inst.SecurityId)
			delete(subscribedInstruments, key)
		}

		log.Printf("Unsubscribed from %d instruments", len(chunk))
	}

	return true
}

// resubscribeToInstruments resubscribes to all previously subscribed instruments
// Called after reconnection to upstream WebSocket
func resubscribeToInstruments() {
	if len(subscribedInstruments) == 0 {
		return
	}

	instruments := []InstrumentInfo{}

	// Convert map keys back to InstrumentInfo objects
	for key := range subscribedInstruments {
		var segment, securityId string
		fmt.Sscanf(key, "%s:%s", &segment, &securityId)

		instruments = append(instruments, InstrumentInfo{
			ExchangeSegment: segment,
			SecurityId:      securityId,
		})
	}

	// Resubscribe to all instruments
	subscribeToInstruments(instruments)
}

// processMarketDataMessage handles different types of market data messages
func processMarketDataMessage(message []byte) {
	// Ensure we have at least header size data (8 bytes)
	if len(message) < 8 {
		log.Printf("Received invalid message: too short (%d bytes)", len(message))
		return
	}

	// Parse message header
	feedCode := message[0]
	// messageLength := binary.LittleEndian.Uint16(message[1:3])
	// exchangeSegment := message[3]
	// securityID := binary.LittleEndian.Uint32(message[4:8])

	// Process message based on feed code
	switch feedCode {
	case FeedCodeTicker:
		processTicker(message)
	case FeedCodeQuote:
		processQuote(message)
	case FeedCodeOI:
		processOpenInterest(message)
	case FeedCodePrevClose:
		processPrevClose(message)
	case FeedCodeFull:
		processFullPacket(message)
	case FeedCodeDisconnect:
		processDisconnection(message)
	default:
		// Unknown feed code
		log.Printf("Received unknown feed code: %d", feedCode)
		return
	}

	// Convert market data to instrument model for additional processing if needed
	instrument := toInstrument(message)
	if instrument == nil {
		log.Printf("Failed to convert message to instrument model")
		return
	}

	// Store the instrument data in the database (already handled in individual processors)
	// But we can use this for deployment processing
	if deploymentSvc != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := processInstrumentDataForDeployments(ctx, instrument); err != nil {
			log.Printf("Error processing deployments for instrument %d: %v", instrument.ExchangeSecurityID, err)
		}
	}

	// Broadcast the message to all connected clients
	broadcastMessage(message)
}

// processTicker handles ticker packet data
func processTicker(data []byte) {
	if len(data) < 16 {
		log.Printf("Ticker packet too short: %d bytes", len(data))
		return
	}

	// Extract header information
	exchangeSegment := data[3]
	securityID := binary.LittleEndian.Uint32(data[4:8])

	// Extract data from ticker packet according to Dhan documentation
	ltp := math.Float32frombits(binary.LittleEndian.Uint32(data[8:12])) // Last Traded Price
	ltt := binary.LittleEndian.Uint32(data[12:16])                      // Last Trade Time

	// Validate extracted data
	if ltp <= 0 {
		log.Printf("Invalid LTP for security %d: %.2f", securityID, ltp)
		return
	}

	if ltt == 0 {
		log.Printf("Invalid timestamp for security %d, using current time", securityID)
		ltt = uint32(time.Now().Unix())
	}

	log.Printf("Ticker: Security ID: %d, Exchange: %d, LTP: %.2f, LTT: %d",
		securityID, exchangeSegment, ltp, ltt)

	// Get symbol from the map of subscribed instruments
	var symbol string
	for key := range subscribedInstruments {
		var segment string
		var id string
		fmt.Sscanf(key, "%s:%s", &segment, &id)
		if id == fmt.Sprintf("%d", securityID) {
			symbol = id // Use ID as symbol if we don't have a better identifier
			break
		}
	}

	// If symbol is empty, use security ID as the symbol
	if symbol == "" {
		symbol = fmt.Sprintf("%d", securityID)
	}

	// Store in ClickHouse - use LTP for all OHLC values since this is a ticker update
	storeInstrumentData(securityID, exchangeSegment, symbol, ltp, ltt, 0.0, ltp, ltp, ltp, ltp)
}

// processQuote handles quote packet data
func processQuote(data []byte) {
	if len(data) < 50 {
		log.Printf("Quote packet too short: %d bytes, expected at least 50", len(data))
		return
	}

	// Extract header information
	exchangeSegment := data[3]
	securityID := binary.LittleEndian.Uint32(data[4:8])

	// Extract basic quote data according to Dhan documentation
	ltp := math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))  // Latest Traded Price
	ltq := binary.LittleEndian.Uint16(data[12:14])                       // Last Traded Quantity
	ltt := binary.LittleEndian.Uint32(data[14:18])                       // Last Trade Time (LTT)
	atp := math.Float32frombits(binary.LittleEndian.Uint32(data[18:22])) // Average Trade Price (ATP)
	volume := binary.LittleEndian.Uint32(data[22:26])                    // Volume

	// Extract OHLC data according to documentation
	open := math.Float32frombits(binary.LittleEndian.Uint32(data[34:38]))  // Day Open Value
	close := math.Float32frombits(binary.LittleEndian.Uint32(data[38:42])) // Day Close Value
	high := math.Float32frombits(binary.LittleEndian.Uint32(data[42:46]))  // Day High Value
	low := math.Float32frombits(binary.LittleEndian.Uint32(data[46:50]))   // Day Low Value

	// Validate extracted data
	if ltp <= 0 {
		log.Printf("Invalid LTP for security %d: %.2f", securityID, ltp)
		return
	}

	if ltt == 0 {
		log.Printf("Invalid timestamp for security %d, using current time", securityID)
		ltt = uint32(time.Now().Unix())
	}

	// If any OHLC value is 0 or invalid, use ltp as fallback
	if open <= 0 {
		open = ltp
	}
	if high <= 0 {
		high = ltp
	}
	if low <= 0 {
		low = ltp
	}
	if close <= 0 {
		close = ltp
	}

	// Validate OHLC relationships
	if high < low {
		log.Printf("Invalid OHLC data for security %d: high %.2f < low %.2f, using LTP", securityID, high, low)
		open, high, low, close = ltp, ltp, ltp, ltp
	}

	log.Printf("Quote: Security ID: %d, LTP: %.2f, LTQ: %d, LTT: %d, ATP: %.2f, Vol: %d, O: %.2f, H: %.2f, L: %.2f, C: %.2f",
		securityID, ltp, ltq, ltt, atp, volume, open, high, low, close)

	// Get symbol from the map of subscribed instruments
	var symbol string
	for key := range subscribedInstruments {
		var segment string
		var id string
		fmt.Sscanf(key, "%s:%s", &segment, &id)
		if id == fmt.Sprintf("%d", securityID) {
			symbol = id // Use ID as symbol if we don't have a better identifier
			break
		}
	}

	// If symbol is empty, use security ID as the symbol
	if symbol == "" {
		symbol = fmt.Sprintf("%d", securityID)
	}

	// Store complete OHLCV data in ClickHouse
	storeInstrumentData(securityID, exchangeSegment, symbol, ltp, ltt, float64(volume), open, high, low, close)
}

// processOpenInterest handles OI packet data
func processOpenInterest(data []byte) {
	if len(data) < 12 {
		return
	}

	// Extract OI data according to Dhan documentation
	oi := binary.LittleEndian.Uint32(data[8:12]) // Open Interest of the contract

	// Extract header information
	exchangeSegment := data[3]
	securityID := binary.LittleEndian.Uint32(data[4:8])

	log.Printf("OI: Security ID: %d, Exchange: %d, Open Interest: %d", securityID, exchangeSegment, oi)
}

// processPrevClose handles previous day close packet data
func processPrevClose(data []byte) {
	if len(data) < 16 {
		return
	}

	// Extract prev close data according to Dhan documentation
	prevClose := math.Float32frombits(binary.LittleEndian.Uint32(data[8:12])) // Previous day closing price
	prevOI := binary.LittleEndian.Uint32(data[12:16])                         // Open Interest - previous day

	// Extract header information
	exchangeSegment := data[3]
	securityID := binary.LittleEndian.Uint32(data[4:8])

	log.Printf("Prev Close: Security ID: %d, Exchange: %d, Prev Close: %.2f, Prev OI: %d",
		securityID, exchangeSegment, prevClose, prevOI)
}

// processFullPacket handles full market data packet
func processFullPacket(data []byte) {
	if len(data) < 63 {
		log.Printf("Full packet too short: %d bytes, expected at least 63", len(data))
		return
	}

	// Extract header information
	exchangeSegment := data[3]
	securityID := binary.LittleEndian.Uint32(data[4:8])

	// Extract key data from full packet according to Dhan documentation
	ltp := math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))  // Latest Traded Price
	ltq := binary.LittleEndian.Uint16(data[12:14])                       // Last Traded Quantity
	ltt := binary.LittleEndian.Uint32(data[14:18])                       // Last Trade Time (LTT)
	atp := math.Float32frombits(binary.LittleEndian.Uint32(data[18:22])) // Average Trade Price (ATP)
	volume := binary.LittleEndian.Uint32(data[22:26])                    // Volume

	// Extract OHLC data according to documentation
	open := math.Float32frombits(binary.LittleEndian.Uint32(data[46:50]))  // Day Open Value
	close := math.Float32frombits(binary.LittleEndian.Uint32(data[50:54])) // Day Close Value
	high := math.Float32frombits(binary.LittleEndian.Uint32(data[54:58]))  // Day High Value
	low := math.Float32frombits(binary.LittleEndian.Uint32(data[58:62]))   // Day Low Value

	// Validate extracted data
	if ltp <= 0 {
		log.Printf("Invalid LTP for security %d: %.2f", securityID, ltp)
		return
	}

	if ltt == 0 {
		log.Printf("Invalid timestamp for security %d, using current time", securityID)
		ltt = uint32(time.Now().Unix())
	}

	// If any OHLC value is 0 or invalid, use ltp as fallback
	if open <= 0 {
		open = ltp
	}
	if high <= 0 {
		high = ltp
	}
	if low <= 0 {
		low = ltp
	}
	if close <= 0 {
		close = ltp
	}

	// Validate OHLC relationships
	if high < low {
		log.Printf("Invalid OHLC data for security %d: high %.2f < low %.2f, using LTP", securityID, high, low)
		open, high, low, close = ltp, ltp, ltp, ltp
	}

	log.Printf("Full Packet: Security ID: %d, LTP: %.2f, LTQ: %d, LTT: %d, ATP: %.2f, Vol: %d, O: %.2f, H: %.2f, L: %.2f, C: %.2f",
		securityID, ltp, ltq, ltt, atp, volume, open, high, low, close)

	// Get symbol from the map of subscribed instruments
	var symbol string
	for key := range subscribedInstruments {
		var segment string
		var id string
		fmt.Sscanf(key, "%s:%s", &segment, &id)
		if id == fmt.Sprintf("%d", securityID) {
			symbol = id // Use ID as symbol if we don't have a better identifier
			break
		}
	}

	// If symbol is empty, use security ID as the symbol
	if symbol == "" {
		symbol = fmt.Sprintf("%d", securityID)
	}

	// Store complete OHLCV data in ClickHouse
	storeInstrumentData(securityID, exchangeSegment, symbol, ltp, ltt, float64(volume), open, high, low, close)
}

// processDisconnection handles disconnection messages
func processDisconnection(data []byte) {
	if len(data) < 10 {
		return
	}

	// Extract disconnection code
	disconnectCode := binary.LittleEndian.Uint16(data[8:10])

	log.Printf("Disconnection received with code: %d", disconnectCode)
}

// broadcastMessage sends a message to all connected clients
func broadcastMessage(message []byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.BinaryMessage, message)
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

// storeInstrumentData stores market data in ClickHouse and processes deployments
func storeInstrumentData(securityID uint32, exchangeSegment byte, symbol string, price float32, timestamp uint32, volume float64, open, high, low, close float32) {
	// Skip if ClickHouse is not connected
	if chDB == nil || instrumentRepository == nil {
		log.Println("Skipping data storage: database not connected")
		return
	}

	// Convert timestamp from Unix to time.Time
	// The LTT (Last Trade Time) from Dhan API is Unix timestamp in UTC
	// Format it consistently with migrate.go approach
	var timestampTime time.Time
	if timestamp > 0 {
		// Parse and format timestamp using the requested format
		// The timestamp from Dhan API is in UTC, so we convert it directly
		utcTime := time.Unix(int64(timestamp), 0).UTC()
		t := utcTime.Format("2006-01-02 15:04:05")
		parsedTime, err := time.Parse("2006-01-02 15:04:05", t)
		if err != nil {
			log.Printf("Error parsing formatted timestamp: %v, using raw Unix timestamp", err)
			timestampTime = time.Unix(int64(timestamp), 0).UTC()
		} else {
			// Keep the parsed time in UTC as that's what the API provides
			timestampTime = parsedTime.UTC()
		}
		// Convert to IST for storage and logging
		timestampTime = timestampTime.In(istLocation)
	} else {
		// Use current time in IST if timestamp is invalid
		timestampTime = time.Now().In(istLocation)
	}

	// Validate OHLC data consistency
	openFloat := float64(open)
	highFloat := float64(high)
	lowFloat := float64(low)
	closeFloat := float64(close)

	// Ensure high is the highest and low is the lowest
	if highFloat < lowFloat {
		log.Printf("Warning: Invalid OHLC data for security %d: high %.2f < low %.2f, swapping values", securityID, highFloat, lowFloat)
		highFloat, lowFloat = lowFloat, highFloat
	}
	if openFloat > highFloat {
		highFloat = openFloat
	}
	if openFloat < lowFloat {
		lowFloat = openFloat
	}
	if closeFloat > highFloat {
		highFloat = closeFloat
	}
	if closeFloat < lowFloat {
		lowFloat = closeFloat
	}

	// Create instrument model with validated data (consistent timezone handling)
	now := time.Now().In(istLocation)
	instrument := &model.Instrument{
		ExchangeSecurityID: int32(securityID),
		Open:               openFloat,
		High:               highFloat,
		Low:                lowFloat,
		Close:              closeFloat,
		Volume:             volume,
		Timestamp:          timestampTime,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// Store in ClickHouse
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := instrumentRepository.Create(ctx, instrument); err != nil {
		log.Printf("Error storing instrument data for security %d: %v", securityID, err)
		log.Printf("Failed data - Exchange: %d, Symbol: %s, OHLCV: %.2f/%.2f/%.2f/%.2f/%.2f, Time: %s",
			exchangeSegment, symbol, openFloat, highFloat, lowFloat, closeFloat, volume, timestampTime.Format(time.RFC3339))
		return
	}

	// Process deployments if deployment service is available
	if deploymentSvc != nil {
		if err := processInstrumentDataForDeployments(ctx, instrument); err != nil {
			log.Printf("Error processing deployments for instrument %d: %v", securityID, err)
		}
	}

	// Log successful storage (only for debugging, can be disabled in production)
	// Display the original UTC time from the API (before IST conversion)
	log.Printf("Successfully stored data for %s (SecurityID: %d): LTP=%.2f, Vol=%.0f at %s",
		symbol, securityID, closeFloat, volume, timestampTime.UTC().Format("15:04:05"))
}

// processInstrumentDataForDeployments processes incoming instrument data and runs the appropriate deployments
// This function was moved from internal/modules/deployment/application/service/deployment_service.go
// to enable real-time deployment processing as live market data arrives.
//
// The original ProcessInstrumentData method in the deployment service handled:
// 1. Storing instrument data in the database
// 2. Finding deployments for the instrument
// 3. Running each deployment's strategy evaluation
// 4. Executing trades when conditions are met
//
// Since the live service already handles data storage, this function focuses on
// the deployment processing aspects.
func processInstrumentDataForDeployments(ctx context.Context, instrumentData *model.Instrument) error {
	// Skip if deployment service is not initialized
	if deploymentSvc == nil {
		return nil
	}

	instrumentID := int(instrumentData.ExchangeSecurityID)

	// Get all active deployments with strategies
	activeDeployments, err := deploymentSvc.GetActiveDeploymentsWithStrategies(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active deployments: %w", err)
	}

	// Filter deployments for this specific instrument
	var relevantDeployments []*deploymentService.DeploymentWithStrategy
	for _, deployment := range activeDeployments {
		if deployment.Deployment.InstrumentID == instrumentID {
			relevantDeployments = append(relevantDeployments, deployment)
		}
	}

	if len(relevantDeployments) == 0 {
		// No deployments for this instrument, skip processing
		return nil
	}

	// // Process each relevant deployment
	// currentTimestamp := instrumentData.Timestamp
	// for _, deploymentWithStrategy := range relevantDeployments {
	// 	if err := deploymentSvc.Run(ctx, deploymentWithStrategy, currentTimestamp); err != nil {
	// 		log.Printf("Error running deployment %s for instrument %d: %v",
	// 			deploymentWithStrategy.Deployment.ID, instrumentID, err)
	// 		// Continue processing other deployments even if one fails
	// 	}
	// }

	return nil
}

// toInstrument converts raw market data to an Instrument model
func toInstrument(data []byte) *model.Instrument {
	// Validate minimum data length
	if len(data) < 8 {
		log.Printf("Invalid data length: %d bytes", len(data))
		return nil
	}

	// Extract header information
	feedCode := data[0]
	securityID := binary.LittleEndian.Uint32(data[4:8])

	var (
		open, high, low, close float64
		volume                 float64
		ltp                    float32 // Last traded price
		ltt                    uint32  // Last trade time
	)

	// Process based on feed code according to Dhan documentation
	switch feedCode {
	case FeedCodeQuote:
		// Quote packet structure
		if len(data) < 50 {
			log.Printf("Quote packet too short: %d bytes", len(data))
			return nil
		}

		ltp = math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))
		ltt = binary.LittleEndian.Uint32(data[14:18])
		volume = float64(binary.LittleEndian.Uint32(data[22:26]))

		// Extract OHLC data from quote packet
		open = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[34:38])))
		close = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[38:42])))
		high = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[42:46])))
		low = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[46:50])))

	case FeedCodeFull:
		// Full packet structure
		if len(data) < 63 {
			log.Printf("Full packet too short: %d bytes", len(data))
			return nil
		}

		ltp = math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))
		ltt = binary.LittleEndian.Uint32(data[14:18])
		volume = float64(binary.LittleEndian.Uint32(data[22:26]))

		// Extract OHLC data from full packet
		open = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[46:50])))
		close = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[50:54])))
		high = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[54:58])))
		low = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[58:62])))

	case FeedCodeTicker:
		// Ticker packet structure
		if len(data) < 16 {
			log.Printf("Ticker packet too short: %d bytes", len(data))
			return nil
		}

		ltp = math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))
		ltt = binary.LittleEndian.Uint32(data[12:16])
		// Volume is not available in ticker packets
		volume = 0

	default:
		log.Printf("Unsupported feed code: %d", feedCode)
		return nil
	}

	// Validate extracted data
	if ltp <= 0 {
		log.Printf("Invalid LTP for security %d: %.2f", securityID, ltp)
		return nil
	}

	// If any OHLC value is 0 or invalid, use ltp as fallback
	ltpFloat := float64(ltp)
	if open <= 0 {
		open = ltpFloat
	}
	if high <= 0 {
		high = ltpFloat
	}
	if low <= 0 {
		low = ltpFloat
	}
	if close <= 0 {
		close = ltpFloat
	}

	// Validate OHLC relationships
	if high < low {
		log.Printf("Invalid OHLC data for security %d: high %.2f < low %.2f, using LTP", securityID, high, low)
		open, high, low, close = ltpFloat, ltpFloat, ltpFloat, ltpFloat
	}

	// Convert LTT to time.Time (consistent with migrate.go approach)
	var timestamp time.Time
	if ltt > 0 {
		// Parse and format timestamp using the requested format
		// The timestamp from Dhan API is in UTC, so we convert it directly
		utcTime := time.Unix(int64(ltt), 0).UTC()
		t := utcTime.Format("2006-01-02 15:04:05")
		parsedTime, err := time.Parse("2006-01-02 15:04:05", t)
		if err != nil {
			log.Printf("Error parsing formatted timestamp: %v, using raw Unix timestamp", err)
			timestamp = time.Unix(int64(ltt), 0).UTC()
		} else {
			// Keep the parsed time in UTC as that's what the API provides
			timestamp = parsedTime.UTC()
		}
		// Convert to IST for storage and logging
		timestamp = timestamp.In(istLocation)
	} else {
		// Use current time in IST if timestamp is invalid
		timestamp = time.Now().In(istLocation)
	}

	// Create and return the instrument object (consistent timezone handling)
	now := time.Now().In(istLocation)
	return &model.Instrument{
		ExchangeSecurityID: int32(securityID),
		Open:               open,
		High:               high,
		Low:                low,
		Close:              close,
		Volume:             volume,
		Timestamp:          timestamp,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}
