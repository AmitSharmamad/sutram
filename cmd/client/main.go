package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"algo/internal/modules/instrument/domain/model"
	"algo/internal/modules/instrument/domain/repository"
	"algo/internal/modules/instrument/infrastructure/persistence/clickhouse"
)

// Feed codes to simulate different types of market data
const (
	FeedCodeTicker    = 2
	FeedCodeQuote     = 4
	FeedCodeOI        = 5
	FeedCodePrevClose = 6
	FeedCodeFull      = 8
)

// Market hours constants
const (
	MarketStartHour = 9
	MarketStartMin  = 15
	MarketEndHour   = 15
	MarketEndMin    = 29
)

// Request codes for subscription
const (
	RequestCodeSubscribe = 15
)

// InstrumentRequest represents the structure for subscribing to instruments
type ClientInstrumentRequest struct {
	RequestCode     int                    `json:"RequestCode"`
	InstrumentCount int                    `json:"InstrumentCount"`
	InstrumentList  []ClientInstrumentInfo `json:"InstrumentList"`
}

// InstrumentInfo represents information about a single instrument
type ClientInstrumentInfo struct {
	ExchangeSegment string `json:"ExchangeSegment"`
	SecurityId      string `json:"SecurityId"`
}

// Instrument represents a financial instrument with its market data
type Instrument struct {
	SecurityID      uint32
	ExchangeSegment byte
	Symbol          string
	LTP             float32 // Last Traded Price
	Open            float32
	High            float32
	Low             float32
	Close           float32
	Volume          uint32
	ATP             float32 // Average Traded Price
	LTQ             uint16  // Last Traded Quantity
	OI              uint32  // Open Interest
	PrevClose       float32
}

var (
	// Sample instruments to simulate (initial state)
	initialInstruments = []Instrument{
		{SecurityID: 12018, ExchangeSegment: 1, Symbol: "HDFCBANK", LTP: 1650.20, Open: 1645.50, High: 1658.75, Low: 1642.30, Close: 1650.20, Volume: 156000, ATP: 1649.85, LTQ: 75, OI: 0, PrevClose: 1643.80},
	}

	// Current instruments (working copy that gets modified during simulation)
	instruments []Instrument

	// Map to store active WebSocket connections
	clients = make(map[*websocket.Conn]bool)
	// Mutex to protect clients map
	clientsMutex sync.Mutex

	// Random number generator
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	// Track current simulation time - will be initialized in main()
	currentSimulationTime time.Time
	// Initial simulation time (for resets)
	initialSimulationTime time.Time

	// Channel to signal simulation restart
	restartSimulationChan = make(chan struct{}, 1)

	// Mutex to protect simulation state
	simulationMutex sync.Mutex

	// ClickHouse database connection
	chDB                 *sqlx.DB
	instrumentRepository repository.InstrumentRepository

	// WebSocket upgrader configuration
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Allow connections from any origin
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// initDB initializes the ClickHouse database connection
func initDB() error {
	// Get database connection details from environment variables
	conn := getEnv("CH_URL", "clickhouse://admin:admin@localhost:9000")

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
	instrumentRepository = clickhouse.NewInstrumentRepository(chDB)

	return nil
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	// Initialize ClickHouse database connection
	if err := initDB(); err != nil {
		log.Printf("Warning: Failed to connect to ClickHouse: %v", err)
		log.Println("Will use random data simulation instead of database data")
	}

	// Initialize simulation time to May 5th, 2025 at 9:15 AM IST
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Printf("Warning: Could not load IST timezone, using UTC: %v", err)
		istLocation = time.UTC
	}
	initialSimulationTime = time.Date(2025, 5, 5, 9, 30, 0, 0, istLocation)
	currentSimulationTime = initialSimulationTime

	// Initialize instruments with the initial state
	resetInstruments()

	// Get port from environment or use default
	port := getEnv("CLIENT_PORT", "9092")
	serverAddr := fmt.Sprintf(":%s", port)

	log.Printf("Starting market data simulator server...")
	log.Printf("Server will listen on port: %s", port)

	// Set up HTTP server for WebSocket connections
	http.HandleFunc("/ws", handleWebSocketConnection)

	// Create HTTP server
	server := &http.Server{
		Addr:    serverAddr,
		Handler: nil,
	}

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal...")
		cancel()

		// Shutdown server gracefully
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	// Start market data simulation in background
	go startMarketDataSimulation(ctx)

	// Start the server
	log.Printf("Market data simulator server listening on %s", serverAddr)
	log.Printf("Connect to: ws://localhost%s/ws", serverAddr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Market data simulator server stopped")
}

// handleWebSocketConnection handles new WebSocket connections
func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	log.Printf("New client connected from: %s", r.RemoteAddr)

	// Add connection to clients map
	clientsMutex.Lock()
	clients[conn] = true
	clientCount := len(clients)
	clientsMutex.Unlock()

	log.Printf("Total connected clients: %d", clientCount)

	// Reset simulation for new client (restart from beginning)
	resetSimulation()

	// Handle client disconnection
	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientCount := len(clients)
		clientsMutex.Unlock()

		conn.Close()
		log.Printf("Client disconnected. Remaining clients: %d", clientCount)
	}()

	// Listen for client messages (subscription requests, etc.)
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		if messageType == websocket.TextMessage {
			log.Printf("Received subscription request from client: %s", string(message))
			// Process subscription request if needed
			handleSubscriptionRequest(conn, message)
		}
	}
}

// handleSubscriptionRequest processes subscription requests from clients
func handleSubscriptionRequest(conn *websocket.Conn, message []byte) {
	var request ClientInstrumentRequest
	if err := json.Unmarshal(message, &request); err != nil {
		log.Printf("Failed to parse subscription request: %v", err)
		return
	}

	log.Printf("Client subscribed to %d instruments", request.InstrumentCount)

	// Send confirmation back to client
	response := map[string]any{
		"status":      "subscribed",
		"count":       request.InstrumentCount,
		"instruments": request.InstrumentList,
	}

	responseData, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, responseData)
}

// startMarketDataSimulation starts the market data simulation in background
func startMarketDataSimulation(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Always use quote feed type
	feedType := byte(FeedCodeQuote)

	// Ensure simulation starts during market hours
	if !isMarketOpen(currentSimulationTime) {
		currentSimulationTime = getNextMarketOpen(currentSimulationTime)
	}

	log.Printf("Market data simulation started from: %s", currentSimulationTime.Format("2006-01-02 15:04:05"))
	log.Printf("Market hours: %02d:%02d to %02d:%02d", MarketStartHour, MarketStartMin, MarketEndHour, MarketEndMin)
	log.Printf("Feed type set to: QUOTE (always)")

	for {
		select {
		case <-ctx.Done():
			log.Println("Market data simulation stopped")
			return
		case <-restartSimulationChan:
			// Simulation restart requested
			if !isMarketOpen(currentSimulationTime) {
				currentSimulationTime = getNextMarketOpen(currentSimulationTime)
			}
			log.Printf("Simulation restarted from: %s", currentSimulationTime.Format("2006-01-02 15:04:05"))
			continue
		case <-ticker.C:
			// Check if there are any connected clients
			clientsMutex.Lock()
			clientCount := len(clients)
			clientsMutex.Unlock()

			if clientCount == 0 {
				continue // No clients connected, skip this iteration
			}

			// Check if market is open before sending data
			if !isMarketOpen(currentSimulationTime) {
				// Market is closed, move to next day
				log.Printf("Market closed at: %s", currentSimulationTime.Format("2006-01-02 15:04:05"))
				simulationMutex.Lock()
				currentSimulationTime = getNextMarketOpen(currentSimulationTime)
				simulationMutex.Unlock()
				log.Printf("Moving to next trading session: %s",
					currentSimulationTime.Format("2006-01-02 15:04:05"))
				continue
			}

			// Increment simulation time by 1 minute for each broadcast
			simulationMutex.Lock()
			currentSimulationTime = currentSimulationTime.Add(1 * time.Minute)
			currentTime := currentSimulationTime
			simulationMutex.Unlock()

			// Check again if we've crossed market close time after incrementing
			if !isMarketOpen(currentTime) {
				log.Printf("Market session ended at: %s", currentTime.Format("2006-01-02 15:04:05"))
				continue
			}

			// Send data for each instrument
			for i := range instruments {
				// Update instrument with data from ClickHouse or random data
				updateInstrumentDataFromDB(&instruments[i])

				// Always use quote feed type

				// Create market data packet
				packet := createMarketDataPacket(&instruments[i], feedType)

				// Broadcast to all connected clients
				broadcastMarketData(packet, &instruments[i], feedType)
			}
		}
	}
}

// isMarketOpen checks if the current simulation time is within market hours
func isMarketOpen(t time.Time) bool {
	hour := t.Hour()
	min := t.Minute()

	// Market is open from 9:15 AM to 3:29 PM
	marketStart := hour*60+min >= MarketStartHour*60+MarketStartMin
	marketEnd := hour*60+min <= MarketEndHour*60+MarketEndMin

	return marketStart && marketEnd
}

// getNextMarketOpen returns the next market opening time
func getNextMarketOpen(t time.Time) time.Time {
	// If it's past market close, move to next day
	hour := t.Hour()
	min := t.Minute()

	if hour*60+min > MarketEndHour*60+MarketEndMin {
		// Move to next day's market open
		nextDay := t.AddDate(0, 0, 1)
		return time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(),
			MarketStartHour, MarketStartMin, 0, 0, t.Location())
	}

	// Same day market open
	return time.Date(t.Year(), t.Month(), t.Day(),
		MarketStartHour, MarketStartMin, 0, 0, t.Location())
}

// broadcastMarketData sends market data to all connected clients
func broadcastMarketData(packet []byte, instrument *Instrument, feedType byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for conn := range clients {
		if err := conn.WriteMessage(websocket.BinaryMessage, packet); err != nil {
			log.Printf("Error sending market data to client: %v", err)
			// Remove disconnected client
			delete(clients, conn)
			conn.Close()
		}
	}

	if len(clients) > 0 {
		log.Printf("Broadcasted %s data for %s (ID: %d) to %d clients: LTP=%.2f, Vol=%d, Time=%s",
			getFeedTypeName(feedType), instrument.Symbol, instrument.SecurityID,
			len(clients), instrument.LTP, instrument.Volume, currentSimulationTime.Format("2006-01-02 15:04:05"))
	}
}

// createMarketDataPacket creates a binary market data packet based on feed type
func createMarketDataPacket(instrument *Instrument, feedType byte) []byte {
	switch feedType {
	case FeedCodeTicker:
		return createTickerPacket(instrument)
	case FeedCodeQuote:
		return createQuotePacket(instrument)
	case FeedCodeFull:
		return createFullPacket(instrument)
	default:
		return createTickerPacket(instrument)
	}
}

// createTickerPacket creates a ticker packet (16 bytes minimum)
func createTickerPacket(instrument *Instrument) []byte {
	packet := make([]byte, 16)

	// Header
	packet[0] = FeedCodeTicker                                        // Feed code
	binary.LittleEndian.PutUint16(packet[1:3], 16)                    // Message length
	packet[3] = instrument.ExchangeSegment                            // Exchange segment
	binary.LittleEndian.PutUint32(packet[4:8], instrument.SecurityID) // Security ID

	// Ticker data
	binary.LittleEndian.PutUint32(packet[8:12], math.Float32bits(instrument.LTP))      // LTP
	binary.LittleEndian.PutUint32(packet[12:16], uint32(currentSimulationTime.Unix())) // LTT (timestamp)

	return packet
}

// createQuotePacket creates a quote packet (50+ bytes)
func createQuotePacket(instrument *Instrument) []byte {
	packet := make([]byte, 63) // Extended size to include OHLC

	// Header
	packet[0] = FeedCodeQuote                                         // Feed code
	binary.LittleEndian.PutUint16(packet[1:3], 63)                    // Message length
	packet[3] = instrument.ExchangeSegment                            // Exchange segment
	binary.LittleEndian.PutUint32(packet[4:8], instrument.SecurityID) // Security ID

	// Quote data
	binary.LittleEndian.PutUint32(packet[8:12], math.Float32bits(instrument.LTP))      // LTP
	binary.LittleEndian.PutUint16(packet[12:14], instrument.LTQ)                       // LTQ
	binary.LittleEndian.PutUint32(packet[14:18], uint32(currentSimulationTime.Unix())) // LTT
	binary.LittleEndian.PutUint32(packet[18:22], math.Float32bits(instrument.ATP))     // ATP
	binary.LittleEndian.PutUint32(packet[22:26], instrument.Volume)                    // Volume

	// Market depth placeholders (sell/buy quantities)
	binary.LittleEndian.PutUint32(packet[26:30], uint32(rng.Intn(10000)+1000)) // Sell quantity
	binary.LittleEndian.PutUint32(packet[30:34], uint32(rng.Intn(10000)+1000)) // Buy quantity

	// OHLC data
	binary.LittleEndian.PutUint32(packet[34:38], math.Float32bits(instrument.Open))  // Open
	binary.LittleEndian.PutUint32(packet[38:42], math.Float32bits(instrument.High))  // High
	binary.LittleEndian.PutUint32(packet[42:46], math.Float32bits(instrument.Low))   // Low
	binary.LittleEndian.PutUint32(packet[46:50], math.Float32bits(instrument.Close)) // Close

	// Additional data
	binary.LittleEndian.PutUint32(packet[50:54], math.Float32bits(instrument.PrevClose)) // Previous close
	binary.LittleEndian.PutUint32(packet[54:58], instrument.OI)                          // Open Interest
	binary.LittleEndian.PutUint32(packet[58:62], uint32(rng.Intn(100000)))               // Total traded value
	packet[62] = 1                                                                       // Trading status (1 = active)

	return packet
}

// createFullPacket creates a full market data packet (63+ bytes)
func createFullPacket(instrument *Instrument) []byte {
	packet := make([]byte, 70) // Extended full packet

	// Header
	packet[0] = FeedCodeFull                                          // Feed code
	binary.LittleEndian.PutUint16(packet[1:3], 70)                    // Message length
	packet[3] = instrument.ExchangeSegment                            // Exchange segment
	binary.LittleEndian.PutUint32(packet[4:8], instrument.SecurityID) // Security ID

	// Full market data
	binary.LittleEndian.PutUint32(packet[8:12], math.Float32bits(instrument.LTP))      // LTP
	binary.LittleEndian.PutUint16(packet[12:14], instrument.LTQ)                       // LTQ
	binary.LittleEndian.PutUint32(packet[14:18], uint32(currentSimulationTime.Unix())) // LTT
	binary.LittleEndian.PutUint32(packet[18:22], math.Float32bits(instrument.ATP))     // ATP
	binary.LittleEndian.PutUint32(packet[22:26], instrument.Volume)                    // Volume

	// Market depth and additional fields
	binary.LittleEndian.PutUint32(packet[26:30], uint32(rng.Intn(10000)+1000))     // Total sell quantity
	binary.LittleEndian.PutUint32(packet[30:34], uint32(rng.Intn(10000)+1000))     // Total buy quantity
	binary.LittleEndian.PutUint32(packet[34:38], math.Float32bits(instrument.ATP)) // VWAP
	binary.LittleEndian.PutUint32(packet[38:42], uint32(rng.Intn(1000000)))        // Total traded value
	binary.LittleEndian.PutUint32(packet[42:46], instrument.OI)                    // Open Interest

	// OHLC data (matching the format expected by processFullPacket)
	binary.LittleEndian.PutUint32(packet[46:50], math.Float32bits(instrument.Open))  // Open
	binary.LittleEndian.PutUint32(packet[50:54], math.Float32bits(instrument.Close)) // Close
	binary.LittleEndian.PutUint32(packet[54:58], math.Float32bits(instrument.High))  // High
	binary.LittleEndian.PutUint32(packet[58:62], math.Float32bits(instrument.Low))   // Low

	// Additional metadata
	binary.LittleEndian.PutUint32(packet[62:66], math.Float32bits(instrument.PrevClose)) // Previous close
	binary.LittleEndian.PutUint32(packet[66:70], uint32(rng.Intn(50000)))                // Day change value

	return packet
}

// getFeedTypeName returns human-readable feed type name
func getFeedTypeName(feedType byte) string {
	switch feedType {
	case FeedCodeTicker:
		return "TICKER"
	case FeedCodeQuote:
		return "QUOTE"
	case FeedCodeFull:
		return "FULL"
	case FeedCodeOI:
		return "OI"
	case FeedCodePrevClose:
		return "PREV_CLOSE"
	default:
		return "UNKNOWN"
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getMarketDataFromDB fetches real market data from ClickHouse for the given timestamp
func getMarketDataFromDB(securityID int32, timestamp time.Time) (*model.Instrument, error) {
	if chDB == nil || instrumentRepository == nil {
		return nil, fmt.Errorf("database not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the data for the exact timestamp or closest to it
	// We'll use a small time window (1 minute) around the target timestamp
	startTime := timestamp.Add(-30 * time.Second)
	endTime := timestamp.Add(30 * time.Second)

	instruments, err := instrumentRepository.FindByExchangeSecurityIDAndTimeRange(
		ctx, securityID, startTime, endTime, "1m")
	if err != nil {
		return nil, fmt.Errorf("error querying market data: %w", err)
	}

	if len(instruments) == 0 {
		return nil, fmt.Errorf("no market data found for timestamp: %s", timestamp.Format("2006-01-02 15:04:05"))
	}

	// Return the first (closest) match
	return instruments[0], nil
}

// resetInstruments resets instruments to their initial state
func resetInstruments() {
	instruments = make([]Instrument, len(initialInstruments))
	copy(instruments, initialInstruments)
	log.Println("Instruments reset to initial state")
}

// resetSimulation resets the simulation to its initial state
func resetSimulation() {
	simulationMutex.Lock()
	defer simulationMutex.Unlock()

	currentSimulationTime = initialSimulationTime
	resetInstruments()

	// Send restart signal (non-blocking)
	select {
	case restartSimulationChan <- struct{}{}:
		log.Printf("Simulation reset to initial time: %s", currentSimulationTime.Format("2006-01-02 15:04:05"))
	default:
		// Channel is full, restart signal already pending
	}
}

// updateInstrumentDataFromDB updates instrument data with real data from ClickHouse
func updateInstrumentDataFromDB(instrument *Instrument) {
	// Try to get real data from database
	marketData, err := getMarketDataFromDB(int32(instrument.SecurityID), currentSimulationTime)
	if err != nil {
		// If database query fails, log and fall back to random data
		log.Printf("Failed to get market data from DB for %s: %v. Using random data.", instrument.Symbol, err)
		updateInstrumentDataRandom(instrument)
		return
	}

	// Map database fields to local instrument struct
	// Use Close price as LTP (Last Traded Price)
	instrument.LTP = float32(marketData.Close)
	instrument.Open = float32(marketData.Open)
	instrument.High = float32(marketData.High)
	instrument.Low = float32(marketData.Low)
	instrument.Close = float32(marketData.Close)
	instrument.Volume = uint32(marketData.Volume)

	// Calculate ATP as average of OHLC
	instrument.ATP = float32((marketData.Open + marketData.High + marketData.Low + marketData.Close) / 4.0)

	// Simulate missing fields with realistic values
	instrument.LTQ = uint16(rng.Intn(200) + 10)

	// For OI (Open Interest), simulate based on volume (only for derivatives)
	if instrument.Symbol != "NIFTY" && instrument.Symbol != "BANKNIFTY" {
		instrument.OI = uint32(float64(instrument.Volume) * 0.8) // OI is typically less than volume
	}

	log.Printf("Updated %s with real data: LTP=%.2f, Vol=%d, Time=%s",
		instrument.Symbol, instrument.LTP, instrument.Volume, currentSimulationTime.Format("2006-01-02 15:04:05"))
}

// updateInstrumentDataRandom updates instrument data with random variations (fallback method)
func updateInstrumentDataRandom(instrument *Instrument) {
	// Generate price movements (±2% max change)
	priceChange := (rng.Float32() - 0.5) * 0.04 * instrument.LTP
	newLTP := instrument.LTP + priceChange

	// Ensure positive price
	if newLTP <= 0 {
		newLTP = instrument.LTP * 0.99
	}

	instrument.LTP = newLTP

	// Update high/low if necessary
	if instrument.LTP > instrument.High {
		instrument.High = instrument.LTP
	}
	if instrument.LTP < instrument.Low {
		instrument.Low = instrument.LTP
	}

	// Update close price (same as LTP for live data)
	instrument.Close = instrument.LTP

	// Update volume (add random amount)
	volumeIncrease := uint32(rng.Intn(1000) + 100)
	instrument.Volume += volumeIncrease

	// Update ATP (weighted average)
	instrument.ATP = (instrument.ATP*0.9 + instrument.LTP*0.1)

	// Update LTQ (random quantity)
	instrument.LTQ = uint16(rng.Intn(200) + 10)

	// Update OI for derivatives (if applicable)
	if rng.Float32() < 0.3 { // 30% chance to update OI
		oiChange := int32(rng.Intn(2000) - 1000)
		newOI := int32(instrument.OI) + oiChange
		if newOI >= 0 {
			instrument.OI = uint32(newOI)
		}
	}
}
