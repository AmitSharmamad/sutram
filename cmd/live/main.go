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

	deploymentModule "algo/internal/modules/deployment"
	deploymentService "algo/internal/modules/deployment/application/service"
	deploymentHttp "algo/internal/modules/deployment/interfaces/http"
	"algo/internal/modules/instrument/domain/model"
	"algo/internal/modules/instrument/domain/repository"
	instrumentRepo "algo/internal/modules/instrument/infrastructure/persistence/clickhouse"
	strategyModule "algo/internal/modules/strategy"
	userModule "algo/internal/modules/user"
	userHttp "algo/internal/modules/user/interfaces/http"
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
	// Deployment service
	deploymentSvc deploymentService.DeploymentService
	// Deployment module for accessing event handlers
	deploymentMod *deploymentModule.Module
	// User controller for authentication
	userController *userHttp.UserController
	// Application context
	appCtx context.Context
	// Context cancel function
	cancelCtx context.CancelFunc
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

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	if err := initDB(); err != nil {
		log.Printf("Error initializing database: %v", err)
		// Continue without database connection
	}

	// Initialize deployment service
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

	// Initialize instrument repository
	instrumentRepository = instrumentRepo.NewInstrumentRepository(chDB)

	return nil
}

// initDeploymentService initializes the deployment service and all its dependencies
func initDeploymentService() error {
	// Check if database is available
	if chDB == nil {
		return fmt.Errorf("database connection not available")
	}

	log.Printf("Initializing deployment service...")

	// Initialize strategy module first (required dependency for deployment)
	strategyMod := strategyModule.New()

	// Initialize user module (required for auth, but not used in live trading)
	userMod := userModule.New()
	userController = userMod.GetController()

	// Initialize deployment module
	deploymentModInstance := deploymentModule.New()
	deploymentModInstance.Initialize(userController, strategyMod)

	// Store the deployment module for WebSocket route setup
	deploymentMod = deploymentModInstance

	// Get the deployment service from the module
	deploymentSvc = deploymentModInstance.GetDeploymentService()

	// Initialize the deployment service for live trading
	if err := deploymentSvc.Initialize(appCtx); err != nil {
		panic("something is going wrong with the deployment setup, exiting")
	}

	// Get active deployments and subscribe to their instruments
	if err := subscribeToActiveDeploymentInstruments(); err != nil {
		log.Printf("Warning: Failed to subscribe to deployment instruments: %v", err)
		// Don't fail the initialization, just log the warning
	}

	log.Printf("Deployment service initialized successfully")

	return nil
}

// subscribeToActiveDeploymentInstruments gets all active deployments and subscribes to their instruments
func subscribeToActiveDeploymentInstruments() error {
	if deploymentSvc == nil {
		return fmt.Errorf("deployment service not initialized")
	}

	log.Printf("Getting active deployments to subscribe to instruments...")

	// Get unique instrument IDs from the deployment service
	// We need to add a helper method to extract this information
	instrumentIDs := getInstrumentIDsFromDeployments()

	if len(instrumentIDs) == 0 {
		log.Printf("No active deployments found, no instruments to subscribe to")
		return nil
	}

	// Convert instrument IDs to InstrumentInfo objects for subscription
	instruments := make([]InstrumentInfo, 0, len(instrumentIDs))
	for _, instrumentID := range instrumentIDs {
		// For now, we'll assume these are NSE instruments
		// In a real implementation, you might want to store exchange information in deployments
		instruments = append(instruments, InstrumentInfo{
			ExchangeSegment: "NSE_EQ", // Default to NSE, can be made configurable
			SecurityId:      fmt.Sprintf("%d", instrumentID),
		})
	}

	log.Printf("Subscribing to %d instruments from active deployments", len(instruments))

	// Subscribe to instruments in a goroutine to avoid blocking initialization
	go func() {
		// Wait for upstream connection to be established
		for upstreamConn == nil {
			time.Sleep(100 * time.Millisecond)
		}

		// Subscribe to instruments
		if success := subscribeToInstruments(instruments); success {
			log.Printf("Successfully subscribed to %d deployment instruments", len(instruments))
		} else {
			log.Printf("Failed to subscribe to deployment instruments")
		}
	}()

	return nil
}

// getInstrumentIDsFromDeployments extracts unique instrument IDs from active deployments
// This function accesses the deployment service's internal state after Initialize() has been called
func getInstrumentIDsFromDeployments() []int {
	if deploymentSvc == nil {
		return []int{}
	}

	// Get instrument IDs from the deployment service
	return deploymentSvc.GetActiveInstrumentIDs()
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

	// Add deployment events WebSocket endpoint if deployment module is available
	if deploymentMod != nil && userController != nil {
		// Get the event handler from the deployment module
		eventHandler := deploymentMod.GetEventHandler()

		// Create API router group and register deployment events WebSocket routes
		apiRouter := router.Group("/api")

		if eventHandler != nil {
			// Register deployment events WebSocket routes
			// This will register the /api/deployment/events/ws endpoint
			deploymentHttp.RegisterWebSocketRoutes(apiRouter, eventHandler, userController)
			log.Printf("Registered deployment events WebSocket endpoint: /ws/api/deployment/events/ws")
		}
	}
}

// setupAPIRoutes sets up REST API routes for the live data module
func setupAPIRoutes(router *gin.RouterGroup) {
	// Endpoint for status check
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "live",
			"time":   time.Now().Format(time.RFC3339),
			"websocket_endpoints": gin.H{
				"market_data":       "/ws/feed",
				"deployment_events": "/ws/api/deployment/events/ws",
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
		log.Printf("Received invalid message: too short")
		return
	}

	// Parse message header
	feedCode := message[0]
	// messageLength := binary.LittleEndian.Uint16(message[1:3])
	// exchangeSegment := message[3]
	// securityID := binary.LittleEndian.Uint32(message[4:8])

	// log.Printf("Received market data - Code: %d, Exchange: %d, SecurityID: %d",
	// 	feedCode, exchangeSegment, securityID)

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
	}

	// Convert market data to instrument model
	instrument := toInstrument(message)

	// Process the instrument data with the deployment service if available
	if deploymentSvc != nil {
		if err := deploymentSvc.ProcessInstrumentData(appCtx, instrument); err != nil {
			log.Printf("Error processing instrument data with deployment service: %v", err)
		} else {
			// log.Printf("Processed instrument %d with deployment service", instrument.ExchangeSecurityID)
		}
	}

	// Broadcast the message to all connected clients
	broadcastMessage(message)
}

// processTicker handles ticker packet data
func processTicker(data []byte) {
	if len(data) < 16 {
		return
	}

	// Extract data from ticker packet according to Dhan documentation
	ltp := math.Float32frombits(binary.LittleEndian.Uint32(data[8:12])) // Last Traded Price
	ltt := binary.LittleEndian.Uint32(data[12:16])                      // Last Trade Time

	// Extract header information
	exchangeSegment := data[3]
	securityID := binary.LittleEndian.Uint32(data[4:8])

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
		return
	}

	// Extract basic quote data according to Dhan documentation
	ltp := math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))  // Latest Traded Price
	ltq := binary.LittleEndian.Uint16(data[12:14])                       // Last Traded Quantity
	ltt := binary.LittleEndian.Uint32(data[14:18])                       // Last Trade Time (LTT)
	atp := math.Float32frombits(binary.LittleEndian.Uint32(data[18:22])) // Average Trade Price (ATP)
	volume := binary.LittleEndian.Uint32(data[22:26])                    // Volume
	// sellQty := binary.LittleEndian.Uint32(data[26:30])               // Total Sell Quantity
	// buyQty := binary.LittleEndian.Uint32(data[30:34])                // Total Buy Quantity

	// Extract OHLC data according to documentation
	var open, high, low, close float32
	if len(data) >= 50 {
		open = math.Float32frombits(binary.LittleEndian.Uint32(data[34:38]))  // Day Open Value
		close = math.Float32frombits(binary.LittleEndian.Uint32(data[38:42])) // Day Close Value
		high = math.Float32frombits(binary.LittleEndian.Uint32(data[42:46]))  // Day High Value
		low = math.Float32frombits(binary.LittleEndian.Uint32(data[46:50]))   // Day Low Value
	}

	// If any value is 0, use ltp as fallback
	if open == 0 {
		open = ltp
	}
	if high == 0 {
		high = ltp
	}
	if low == 0 {
		low = ltp
	}
	if close == 0 {
		close = ltp
	}

	log.Printf("Quote: LTP: %.2f, LTQ: %d, LTT: %d, ATP: %.2f, Vol: %d, O: %.2f, H: %.2f, L: %.2f, C: %.2f",
		ltp, ltq, ltt, atp, volume, open, high, low, close)

	// Extract header information
	exchangeSegment := data[3]
	securityID := binary.LittleEndian.Uint32(data[4:8])

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
	if len(data) < 163 {
		return
	}

	// Extract key data from full packet according to Dhan documentation
	ltp := math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))  // Latest Traded Price
	ltq := binary.LittleEndian.Uint16(data[12:14])                       // Last Traded Quantity
	ltt := binary.LittleEndian.Uint32(data[14:18])                       // Last Trade Time (LTT)
	atp := math.Float32frombits(binary.LittleEndian.Uint32(data[18:22])) // Average Trade Price (ATP)
	volume := binary.LittleEndian.Uint32(data[22:26])                    // Volume
	// sellQty := binary.LittleEndian.Uint32(data[26:30])               // Total Sell Quantity
	// buyQty := binary.LittleEndian.Uint32(data[30:34])                // Total Buy Quantity
	// oi := binary.LittleEndian.Uint32(data[34:38])                    // Open Interest
	// highOI := binary.LittleEndian.Uint32(data[38:42])                // Highest Open Interest
	// lowOI := binary.LittleEndian.Uint32(data[42:46])                 // Lowest Open Interest

	// Extract OHLC data according to documentation
	var open, high, low, close float32
	if len(data) >= 63 {
		open = math.Float32frombits(binary.LittleEndian.Uint32(data[46:50]))  // Day Open Value
		close = math.Float32frombits(binary.LittleEndian.Uint32(data[50:54])) // Day Close Value
		high = math.Float32frombits(binary.LittleEndian.Uint32(data[54:58]))  // Day High Value
		low = math.Float32frombits(binary.LittleEndian.Uint32(data[58:62]))   // Day Low Value
	}

	// Market Depth data starts at byte 62 (5 packets of 20 bytes each)
	// We're not processing market depth for now, but it's available in data[62:162]

	// If any value is 0, use ltp as fallback
	if open == 0 {
		open = ltp
	}
	if high == 0 {
		high = ltp
	}
	if low == 0 {
		low = ltp
	}
	if close == 0 {
		close = ltp
	}

	log.Printf("Full Packet: LTP: %.2f, LTQ: %d, LTT: %d, ATP: %.2f, Vol: %d, O: %.2f, H: %.2f, L: %.2f, C: %.2f",
		ltp, ltq, ltt, atp, volume, open, high, low, close)

	// Extract header information
	exchangeSegment := data[3]
	securityID := binary.LittleEndian.Uint32(data[4:8])

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

// storeInstrumentData stores market data in ClickHouse
func storeInstrumentData(securityID uint32, exchangeSegment byte, symbol string, price float32, timestamp uint32, volume float64, open, high, low, close float32) {
	// Skip if ClickHouse is not connected
	if chDB == nil || instrumentRepository == nil {
		log.Println("Skipping data storage: database not connected")
		return
	}

	// Convert timestamp from Unix to time.Time
	timestampTime := time.Unix(int64(timestamp), 0)

	// Create instrument model
	instrument := &model.Instrument{
		ExchangeSecurityID: int32(securityID),
		Open:               float64(open),
		High:               float64(high),
		Low:                float64(low),
		Close:              float64(close),
		Volume:             volume,
		Timestamp:          timestampTime,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Store in ClickHouse
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := instrumentRepository.Create(ctx, instrument); err != nil {
		log.Printf("Error storing instrument data: %v", err)
		return
	}

	// log.Printf("Stored ticker data for %s (ExchangeSecurityID: %d) at %s",
	// 	symbol, securityID, timestampTime.Format(time.RFC3339))
}

// toInstrument converts raw market data to an Instrument model
func toInstrument(data []byte) *model.Instrument {
	// Extract header information
	// exchangeSegment := data[3] - not used currently
	securityID := binary.LittleEndian.Uint32(data[4:8])

	// Extract market data based on the feed type
	feedCode := data[0]

	var (
		open, high, low, close float64
		volume                 float64
		ltp                    float32 // Last traded price
		ltt                    uint32  // Last trade time
	)

	// Process based on feed code according to Dhan documentation
	if feedCode == FeedCodeQuote {
		// Quote packet structure
		if len(data) >= 12 {
			ltp = math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))
		}
		if len(data) >= 18 {
			ltt = binary.LittleEndian.Uint32(data[14:18])
		}
		if len(data) >= 26 {
			volume = float64(binary.LittleEndian.Uint32(data[22:26]))
		}

		// Extract OHLC data from quote packet
		if len(data) >= 50 {
			open = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[34:38])))
			close = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[38:42])))
			high = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[42:46])))
			low = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[46:50])))
		}
	} else if feedCode == FeedCodeFull {
		// Full packet structure
		if len(data) >= 12 {
			ltp = math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))
		}
		if len(data) >= 18 {
			ltt = binary.LittleEndian.Uint32(data[14:18])
		}
		if len(data) >= 26 {
			volume = float64(binary.LittleEndian.Uint32(data[22:26]))
		}

		// Extract OHLC data from full packet
		if len(data) >= 63 {
			open = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[46:50])))
			close = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[50:54])))
			high = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[54:58])))
			low = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[58:62])))
		}
	} else if feedCode == FeedCodeTicker {
		// Ticker packet structure
		if len(data) >= 12 {
			ltp = math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))
		}
		if len(data) >= 16 {
			ltt = binary.LittleEndian.Uint32(data[12:16])
		}
	}

	// If any value is 0, use ltp as fallback
	ltpFloat := float64(ltp)
	if open == 0 {
		open = ltpFloat
	}
	if high == 0 {
		high = ltpFloat
	}
	if low == 0 {
		low = ltpFloat
	}
	if close == 0 {
		close = ltpFloat
	}

	// Convert LTT to time.Time, fallback to current time if LTT is 0
	var timestamp time.Time
	if ltt > 0 {
		timestamp = time.Unix(int64(ltt), 0)
	} else {
		timestamp = time.Now()
	}

	// Create and return the instrument object
	return &model.Instrument{
		ExchangeSecurityID: int32(securityID),
		Open:               open,
		High:               high,
		Low:                low,
		Close:              close,
		Volume:             volume,
		Timestamp:          timestamp, // Use LTT from packet
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
