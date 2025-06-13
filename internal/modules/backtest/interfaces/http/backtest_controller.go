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
	"github.com/xuri/excelize/v2"
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
				Data: map[string]any{
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
			Data: map[string]any{
				"results":   response,
				"timestamp": time.Now(),
			},
		}

		// Send completion event
		eventChan <- BacktestEvent{
			Type: "done",
			Data: map[string]any{
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

// ExportToExcel handles the export backtest data to Excel endpoint
func (c *BacktestController) ExportToExcel(ctx *gin.Context) {
	// Get strategy ID from query parameter
	strategyIDStr := ctx.Query("strategy_id")
	if strategyIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "strategy_id query parameter is required",
		})
		return
	}

	// Parse strategy ID
	strategyID, err := uuid.Parse(strategyIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid strategy_id format",
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

	// Get backtest data from engine service
	response, err := c.engineService.Start(ctx, strategyID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve backtest data: " + err.Error(),
		})
		return
	}

	// Generate Excel file
	excelFile, err := c.generateExcelFile(response)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Excel file: " + err.Error(),
		})
		return
	}

	// Generate filename with timestamp
	filename := fmt.Sprintf("backtest_results_%s_%s.xlsx",
		strategyIDStr,
		time.Now().Format("20060102_150405"))

	// Set response headers for file download
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Write Excel file to response
	if err := excelFile.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write Excel file: " + err.Error(),
		})
		return
	}
}

// ExportIndicatorsToExcel handles the export indicator values to Excel endpoint
func (c *BacktestController) ExportIndicatorsToExcel(ctx *gin.Context) {
	// Get strategy ID from query parameter
	strategyIDStr := ctx.Query("strategy_id")
	if strategyIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "strategy_id query parameter is required",
		})
		return
	}

	// Parse strategy ID
	strategyID, err := uuid.Parse(strategyIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid strategy_id format",
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

	// Get indicator data from engine service
	indicatorData, err := c.engineService.GetIndicatorData(ctx, strategyID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve indicator data: " + err.Error(),
		})
		return
	}

	// Generate Excel file with indicator data
	excelFile, err := c.generateIndicatorExcelFile(indicatorData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Excel file: " + err.Error(),
		})
		return
	}

	// Generate filename with timestamp
	filename := fmt.Sprintf("indicator_data_%s_%s.xlsx",
		strategyIDStr,
		time.Now().Format("20060102_150405"))

	// Set response headers for file download
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Write Excel file to response
	if err := excelFile.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write Excel file: " + err.Error(),
		})
		return
	}
}

// generateExcelFile creates an Excel file from backtest data
func (c *BacktestController) generateExcelFile(responses []service.Response) (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create style for headers
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D3D3D3"},
			Pattern: 1,
		},
	})
	if err != nil {
		return nil, err
	}

	for i, response := range responses {
		sheetName := fmt.Sprintf("Backtest_%d", i+1)

		// Create new sheet for each response (if multiple instruments)
		if i > 0 {
			index, err := f.NewSheet(sheetName)
			if err != nil {
				return nil, err
			}
			f.SetActiveSheet(index)
		} else {
			// Rename the default sheet
			f.SetSheetName("Sheet1", sheetName)
		}

		// Add instrument information
		if response.InstrumentMaster != nil {
			f.SetCellValue(sheetName, "A1", "Instrument Information")
			f.SetCellStyle(sheetName, "A1", "A1", headerStyle)
			f.SetCellValue(sheetName, "A2", "Symbol")
			f.SetCellValue(sheetName, "B2", response.InstrumentMaster.Symbol)
			f.SetCellValue(sheetName, "A3", "Exchange")
			f.SetCellValue(sheetName, "B3", response.InstrumentMaster.Exchange)
			f.SetCellValue(sheetName, "A4", "Instrument ID")
			f.SetCellValue(sheetName, "B4", response.InstrumentID)
		}

		// Add positions table
		if len(response.Positions) > 0 {
			startRow := 6
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), "Positions")
			f.SetCellStyle(sheetName, fmt.Sprintf("A%d", startRow), fmt.Sprintf("A%d", startRow), headerStyle)

			// Position headers
			headerRow := startRow + 1
			headers := []string{
				"Position ID", "State", "Direction", "Entry Signal", "Entry Time Signal",
				"Entry Price Signal", "Exit Signal", "Exit Time Signal", "Exit Price Signal",
				"P&L Absolute", "P&L Percentage", "Position Size", "Transaction Cost",
				"Entry Fees", "Exit Fees", "Net P&L",
			}

			for j, header := range headers {
				col := string(rune('A' + j))
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, headerRow), header)
				f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, headerRow), fmt.Sprintf("%s%d", col, headerRow), headerStyle)
			}

			// Position data
			for k, position := range response.Positions {
				row := headerRow + 1 + k
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), position.ID.String())
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), string(position.State))
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), string(position.Direction))
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), position.EntrySignal)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), position.EntryTimeSignal.Format("2006-01-02 15:04:05"))
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), position.EntryPriceSignal)
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), position.ExitSignal)
				if !position.ExitTimeSignal.IsZero() {
					f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), position.ExitTimeSignal.Format("2006-01-02 15:04:05"))
				}
				f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), position.ExitPriceSignal)
				f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), position.ProfitLossAbsolute)
				f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), position.ProfitLossPercentage)
				f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), position.PositionSize)
				f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), position.TransactionCost)
				f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), position.EntryFees)
				f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), position.ExitFees)
				f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), position.NetProfitLoss)
			}
		}

		// Add metrics table
		if response.Metrics != nil {
			metricsStartRow := 6 + len(response.Positions) + 4
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", metricsStartRow), "Performance Metrics")
			f.SetCellStyle(sheetName, fmt.Sprintf("A%d", metricsStartRow), fmt.Sprintf("A%d", metricsStartRow), headerStyle)

			metricsData := [][]interface{}{
				{"Total Return %", response.Metrics.TotalReturnPct},
				{"Annualized Return", response.Metrics.AnnualizedReturn},
				{"Risk Adjusted Return", response.Metrics.RiskAdjustedReturn},
				{"Sharpe Ratio", response.Metrics.SharpeRatio},
				{"Sortino Ratio", response.Metrics.SortinoRatio},
				{"Max Drawdown", response.Metrics.MaxDrawdown},
				{"Max Drawdown Duration", response.Metrics.MaxDrawdownDuration},
				{"Value at Risk", response.Metrics.ValueAtRisk},
				{"Standard Deviation", response.Metrics.StdDeviation},
				{"Downside Deviation", response.Metrics.DownsideDeviation},
				{"Total Trades", response.Metrics.TotalTrades},
				{"Win Rate", response.Metrics.WinRate},
				{"Average Win Size", response.Metrics.AvgWinSize},
				{"Average Loss Size", response.Metrics.AvgLossSize},
				{"Average Trade Duration", response.Metrics.AvgTradeDuration},
				{"Profit Factor", response.Metrics.ProfitFactor},
				{"Average Position Size", response.Metrics.AvgPositionSize},
				{"Max Position Size", response.Metrics.MaxPositionSize},
				{"Position Turnover", response.Metrics.PositionTurnover},
				{"Average Exposure", response.Metrics.AvgExposure},
				{"Total Commission", response.Metrics.TotalCommission},
				{"Total Slippage", response.Metrics.TotalSlippage},
				{"Transaction Costs %", response.Metrics.TransactionCostsPct},
			}

			for j, metric := range metricsData {
				row := metricsStartRow + 1 + j
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), metric[0])
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), metric[1])
			}
		}

		// Auto-fit columns
		cols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}
		for _, col := range cols {
			if err := f.SetColWidth(sheetName, col, col, 15); err != nil {
				return nil, err
			}
		}
	}

	return f, nil
}

// generateIndicatorExcelFile creates an Excel file from indicator data
func (c *BacktestController) generateIndicatorExcelFile(indicatorData map[string]map[string][]service.IndicatorValue) (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create style for headers
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D3D3D3"},
			Pattern: 1,
		},
	})
	if err != nil {
		return nil, err
	}

	sheetIndex := 0
	for instrumentID, indicators := range indicatorData {
		sheetName := fmt.Sprintf("Instrument_%s", instrumentID)

		// Create new sheet for each instrument
		if sheetIndex > 0 {
			index, err := f.NewSheet(sheetName)
			if err != nil {
				return nil, err
			}
			f.SetActiveSheet(index)
		} else {
			// Rename the default sheet
			f.SetSheetName("Sheet1", sheetName)
		}

		// Collect all unique timestamps across all indicators for this instrument
		timestampSet := make(map[time.Time]bool)
		for _, values := range indicators {
			for _, val := range values {
				timestampSet[val.Timestamp] = true
			}
		}

		// Convert to sorted slice
		timestamps := make([]time.Time, 0, len(timestampSet))
		for ts := range timestampSet {
			timestamps = append(timestamps, ts)
		}

		// Sort timestamps
		for i := 0; i < len(timestamps)-1; i++ {
			for j := i + 1; j < len(timestamps); j++ {
				if timestamps[i].After(timestamps[j]) {
					timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
				}
			}
		}

		// Create headers: Timestamp + indicator names
		headers := []string{"Timestamp"}
		indicatorNames := make([]string, 0, len(indicators))
		for indicatorName := range indicators {
			headers = append(headers, indicatorName)
			indicatorNames = append(indicatorNames, indicatorName)
		}

		// Write headers
		for i, header := range headers {
			col := string(rune('A' + i))
			f.SetCellValue(sheetName, fmt.Sprintf("%s1", col), header)
			f.SetCellStyle(sheetName, fmt.Sprintf("%s1", col), fmt.Sprintf("%s1", col), headerStyle)
		}

		// Write data rows
		for rowIndex, timestamp := range timestamps {
			row := rowIndex + 2 // Start from row 2 (after header)

			// Write timestamp
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), timestamp.Format("2006-01-02 15:04:05"))

			// Write indicator values
			for colIndex, indicatorName := range indicatorNames {
				col := string(rune('B' + colIndex)) // Start from column B

				// Find the value for this timestamp and indicator
				var value interface{} = ""
				for _, val := range indicators[indicatorName] {
					if val.Timestamp.Equal(timestamp) {
						value = val.Value
						break
					}
				}

				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, row), value)
			}
		}

		// Auto-fit columns
		for i := 0; i < len(headers); i++ {
			col := string(rune('A' + i))
			if err := f.SetColWidth(sheetName, col, col, 20); err != nil {
				return nil, err
			}
		}

		sheetIndex++
	}

	return f, nil
}
