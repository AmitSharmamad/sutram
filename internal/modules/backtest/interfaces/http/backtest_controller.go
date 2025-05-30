package http

import (
	"algo/internal/modules/backtest/application/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BacktestEvent represents an event sent via Server-Sent Events during backtest execution
type BacktestEvent struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// BacktestController handles HTTP requests for backtests
type BacktestController struct {
	backtestService service.BacktestService
	engineService   service.EngineService
}

// NewBacktestController creates a new backtest controller
func NewBacktestController(backtestService service.BacktestService, engineService service.EngineService) *BacktestController {
	return &BacktestController{
		backtestService: backtestService,
		engineService:   engineService,
	}
}

// HealthCheck handles the health check endpoint
func (c *BacktestController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "Backtest module is running",
	})
}

// RunBacktest handles the run backtest endpoint with Server-Sent Events
func (c *BacktestController) RunBacktest(ctx *gin.Context) {
	// Get strategy ID from query parameter
	strategyIDStr := ctx.Query("strategy_id")
	if strategyIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "strategy_id query parameter is required",
		})
		return
	}

	// Parse the strategy ID
	strategyID, err := uuid.Parse(strategyIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid strategy ID format",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Set headers for Server-Sent Events
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Cache-Control")

	// Create a channel for sending updates
	eventChan := make(chan BacktestEvent, 3) // Increase capacity for 3 events: started, results, done

	// Send immediate "started" event to establish 200 OK
	// startedEvent := BacktestEvent{
	// 	Type: "started",
	// 	Data: map[string]interface{}{
	// 		"message":     "Backtest request received and processing started",
	// 		"strategy_id": strategyID.String(),
	// 		"timestamp":   time.Now(),
	// 	},
	// }

	// Write the initial event immediately to establish 200 OK
	// data, err := json.Marshal(startedEvent.Data)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal started event"})
	// 	return
	// }

	ctx.Writer.WriteHeader(http.StatusOK)
	// fmt.Fprintf(ctx.Writer, "event: %s\n", startedEvent.Type)
	// fmt.Fprintf(ctx.Writer, "data: %s\n\n", data)
	ctx.Writer.Flush()

	// Start the backtest in a goroutine
	go func() {
		defer close(eventChan)

		// Call the engine service to start the backtest (use the original Start method)
		response, err := c.engineService.Start(ctx, strategyID, userID)
		if err != nil {
			// If there's an error, send done event with error message instead of separate error event
			eventChan <- BacktestEvent{
				Type: "done",
				Data: map[string]interface{}{
					"message":   "Backtest failed: " + err.Error(),
					"error":     true,
					"timestamp": time.Now(),
				},
			}
			return
		}

		// Send results event
		eventChan <- BacktestEvent{
			Type: "results",
			Data: map[string]interface{}{
				"results":   response,
				"timestamp": time.Now(),
			},
		}

		// Send completion event
		eventChan <- BacktestEvent{
			Type: "done",
			Data: map[string]interface{}{
				"message":   "Backtest completed successfully",
				"error":     false,
				"timestamp": time.Now(),
			},
		}
	}()

	// Stream remaining events to the client
	ctx.Stream(func(w io.Writer) bool {
		select {
		case event, ok := <-eventChan:
			if !ok {
				return false // Channel closed, stop streaming
			}

			// Marshal event data to JSON
			data, err := json.Marshal(event.Data)
			if err != nil {
				return false
			}

			// Write SSE formatted data
			fmt.Fprintf(w, "event: %s\n", event.Type)
			fmt.Fprintf(w, "data: %s\n\n", data)

			// Flush immediately to ensure data is sent
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

			return true
		case <-ctx.Request.Context().Done():
			return false // Client disconnected
		}
	})
}
