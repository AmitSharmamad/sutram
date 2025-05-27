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

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// Feed codes to simulate different types of market data
const (
	FeedCodeTicker    = 2
	FeedCodeQuote     = 4
	FeedCodeOI        = 5
	FeedCodePrevClose = 6
	FeedCodeFull      = 8
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
	// Sample instruments to simulate
	instruments = []Instrument{
		{SecurityID: 12018, ExchangeSegment: 1, Symbol: "HDFCBANK", LTP: 1650.20, Open: 1645.50, High: 1658.75, Low: 1642.30, Close: 1650.20, Volume: 156000, ATP: 1649.85, LTQ: 75, OI: 0, PrevClose: 1643.80},
	}

	// Map to store active WebSocket connections
	clients = make(map[*websocket.Conn]bool)
	// Mutex to protect clients map
	clientsMutex sync.Mutex

	// Random number generator
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

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

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

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
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	feedTypes := []byte{FeedCodeTicker, FeedCodeQuote, FeedCodeFull}

	log.Println("Market data simulation started...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Market data simulation stopped")
			return
		case <-ticker.C:
			// Check if there are any connected clients
			clientsMutex.Lock()
			clientCount := len(clients)
			clientsMutex.Unlock()

			if clientCount == 0 {
				continue // No clients connected, skip this iteration
			}

			// Send data for each instrument
			for i := range instruments {
				// Update instrument with random variations
				updateInstrumentData(&instruments[i])

				// Choose random feed type
				feedType := feedTypes[rng.Intn(len(feedTypes))]

				// Create market data packet
				packet := createMarketDataPacket(&instruments[i], feedType)

				// Broadcast to all connected clients
				broadcastMarketData(packet, &instruments[i], feedType)
			}
		}
	}
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
		log.Printf("Broadcasted %s data for %s (ID: %d) to %d clients: LTP=%.2f, Vol=%d",
			getFeedTypeName(feedType), instrument.Symbol, instrument.SecurityID,
			len(clients), instrument.LTP, instrument.Volume)
	}
}

// updateInstrumentData updates instrument data with realistic random variations
func updateInstrumentData(instrument *Instrument) {
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
	binary.LittleEndian.PutUint32(packet[8:12], math.Float32bits(instrument.LTP)) // LTP
	binary.LittleEndian.PutUint32(packet[12:16], uint32(time.Now().Unix()))       // LTT (timestamp)

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
	binary.LittleEndian.PutUint32(packet[8:12], math.Float32bits(instrument.LTP))  // LTP
	binary.LittleEndian.PutUint16(packet[12:14], instrument.LTQ)                   // LTQ
	binary.LittleEndian.PutUint32(packet[14:18], uint32(time.Now().Unix()))        // LTT
	binary.LittleEndian.PutUint32(packet[18:22], math.Float32bits(instrument.ATP)) // ATP
	binary.LittleEndian.PutUint32(packet[22:26], instrument.Volume)                // Volume

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
	binary.LittleEndian.PutUint32(packet[8:12], math.Float32bits(instrument.LTP))  // LTP
	binary.LittleEndian.PutUint16(packet[12:14], instrument.LTQ)                   // LTQ
	binary.LittleEndian.PutUint32(packet[14:18], uint32(time.Now().Unix()))        // LTT
	binary.LittleEndian.PutUint32(packet[18:22], math.Float32bits(instrument.ATP)) // ATP
	binary.LittleEndian.PutUint32(packet[22:26], instrument.Volume)                // Volume

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
