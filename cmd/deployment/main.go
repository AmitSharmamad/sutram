// Deployment Scheduler Service
//
// This service provides both simulated and real-time deployment execution capabilities:
//
// DEVELOPMENT MODE (APP_ENV=development):
// - Uses time simulator for fast-forward testing
// - Simulates market time progression (1 minute per second by default)
// - Starts from May 5, 2025 at 9:15 AM and progresses through market hours
// - Configurable simulation speed via SIMULATION_SPEED environment variable
//
// PRODUCTION MODE (APP_ENV=production):
// - Uses cron-based scheduler for real-time execution
// - Executes deployments based on actual market hours and timeframes
// - Supports different cron expressions for each timeframe (1m, 5m, 1h, etc.)
// - Runs only during configured market hours on weekdays
//
// CONFIGURATION:
// Environment variables for customization:
// - APP_ENV: development|production|staging
// - SIMULATION_SPEED: Duration (e.g., "1s", "500ms") for simulation tick speed
// - FORCE_SIMULATOR: "true" to force simulator mode even in production
// - MARKET_START_HOUR: Market opening hour (default: 9)
// - MARKET_START_MIN: Market opening minute (default: 15)
// - MARKET_END_HOUR: Market closing hour (default: 15)
// - MARKET_END_MIN: Market closing minute (default: 30)
// - CUSTOM_CRON_1M: Custom cron expression for 1-minute timeframe
// - CUSTOM_CRON_5M: Custom cron expression for 5-minute timeframe
// - etc.
//
// FEATURES:
// - WebSocket server for real-time deployment events
// - Timeframe-specific execution scheduling
// - Market hours validation
// - Graceful shutdown handling
// - Comprehensive logging and monitoring
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	deploymentModule "algo/internal/modules/deployment"
	"algo/internal/modules/deployment/application/dto"
	deploymentService "algo/internal/modules/deployment/application/service"
	strategyModule "algo/internal/modules/strategy"
	userModule "algo/internal/modules/user"
)

// Configuration constants for different environments
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvStaging     = "staging"
)

// SchedulerConfig holds configuration for the scheduler
type SchedulerConfig struct {
	Environment     string
	UseSimulator    bool
	MarketStartHour int
	MarketStartMin  int
	MarketEndHour   int
	MarketEndMin    int
	SimulationSpeed time.Duration // How often to tick in simulation mode
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	if err := initDB(); err != nil {
		log.Printf("Error initializing database: %v", err)
		os.Exit(1)
	}

	// Initialize user and strategy modules
	strategyMod := strategyModule.New()
	userMod := userModule.New()
	userController := userMod.GetController()

	// Initialize deployment module
	deploymentMod := deploymentModule.New()
	deploymentMod.Initialize(userController, strategyMod)
	deploymentSvc := deploymentMod.GetDeploymentService()

	if err := deploymentSvc.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize deployment service: %v", err)
	}

	// Start WebSocket server for deployment events
	go startWebSocketServer(deploymentMod)

	// Get active deployments and group them by timeframe
	deploymentsByTimeframe, err := getActiveDeploymentsByTimeframe(ctx, deploymentSvc)
	if err != nil {
		log.Fatalf("Failed to get active deployments: %v", err)
	}

	// Create deployment event listener for dynamic updates
	eventListener := NewDeploymentEventListener(deploymentsByTimeframe, deploymentSvc)

	// Start WebSocket listener to handle deployment changes
	go startWebSocketListener(ctx, eventListener)

	// Get scheduler configuration
	config := getSchedulerConfig()
	log.Printf("Running in %s environment", config.Environment)

	if config.UseSimulator {
		// Start custom scheduler for simulated time progression
		log.Println("Starting SIMULATED TIME SCHEDULER")
		log.Printf("Features: Fast-forward time simulation, tick every %v", config.SimulationSpeed)
		log.Printf("Market hours: %02d:%02d to %02d:%02d",
			config.MarketStartHour, config.MarketStartMin,
			config.MarketEndHour, config.MarketEndMin)
		go startSimulatedTimeScheduler(ctx, eventListener, deploymentSvc, config)
	} else {
		// Start cron-based scheduler for production environment
		log.Println("Starting CRON-BASED SCHEDULER")
		log.Println("Features: Real-time execution based on actual market hours")
		log.Printf("Market hours: %02d:%02d to %02d:%02d",
			config.MarketStartHour, config.MarketStartMin,
			config.MarketEndHour, config.MarketEndMin)
		go startCronScheduler(ctx, eventListener, deploymentSvc, config)
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down deployment scheduler...")

	// Cancel context to signal shutdown
	cancel()
}

// initDB initializes the ClickHouse database connection
func initDB() error {
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

	log.Println("Connected to ClickHouse database")

	return nil
}

// getActiveDeploymentsByTimeframe gets active deployments grouped by timeframe
func getActiveDeploymentsByTimeframe(ctx context.Context, deploymentSvc deploymentService.DeploymentService) (map[string][]*deploymentService.DeploymentWithStrategy, error) {
	// Get all active deployments
	activeDeployments, err := deploymentSvc.GetActiveDeploymentsWithStrategies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active deployments: %w", err)
	}

	// Group deployments by their strategy's chart timeframe
	deploymentsByTimeframe := make(map[string][]*deploymentService.DeploymentWithStrategy)

	for _, deployment := range activeDeployments {
		timeframe := string(deployment.Strategy.ChartTimeframe)
		deploymentsByTimeframe[timeframe] = append(deploymentsByTimeframe[timeframe], deployment)
	}

	return deploymentsByTimeframe, nil
}

// startSimulatedTimeScheduler runs deployments with simulated time progression
// Starts from May 5, 2025 9:15:00 and increments by 1 minute every second
// Runs from 9:15:00 to 15:30:00 (3:30 PM) each day, then moves to next day
// Executes deployments based on their specific timeframes (1m, 5m, 15m, 1h, etc.)
func startSimulatedTimeScheduler(ctx context.Context, eventListener *DeploymentEventListener, svc deploymentService.DeploymentService, config *SchedulerConfig) {
	// Start date: May 5, 2025 at 9:15:00 IST
	startDate := time.Date(2025, 5, 5, 9, 15, 0, 0, time.Local)
	currentTimestamp := startDate

	marketStartHour, marketStartMinute := config.MarketStartHour, config.MarketStartMin
	marketEndHour, marketEndMinute := config.MarketEndHour, config.MarketEndMin

	ticker := time.NewTicker(config.SimulationSpeed) // Execute based on simulation speed
	defer ticker.Stop()

	log.Printf("Starting simulated time scheduler from %v", startDate)
	log.Printf("Market hours: %02d:%02d to %02d:%02d", marketStartHour, marketStartMinute, marketEndHour, marketEndMinute)
	log.Printf("Timeframe-based execution: deployments will run only when their timeframe conditions are met")

	// Get initial deployments from event listener
	deploymentsByTimeframe := eventListener.GetDeploymentsByTimeframe()

	// Log deployment counts by timeframe with execution schedule
	log.Printf("📋 DEPLOYMENT SCHEDULE:")
	for timeframe, deployments := range deploymentsByTimeframe {
		schedule := getTimeframeScheduleDescription(timeframe)
		nextExec := getNextExecutionTime(timeframe, startDate)
		log.Printf("  • %s: %d deployments - %s", timeframe, len(deployments), schedule)
		log.Printf("    Next execution: %v", nextExec.Format("2006-01-02 15:04:05"))
	}

	if len(deploymentsByTimeframe) == 0 {
		log.Printf("⚠️ No deployments found! The scheduler will run but won't execute any strategies.")
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Simulated time scheduler stopped")
			return
		case <-ticker.C:
			// Check if current time is within market hours
			currentHour := currentTimestamp.Hour()
			currentMinute := currentTimestamp.Minute()

			// Check if we're within market hours
			marketStart := marketStartHour*60 + marketStartMinute
			marketEnd := marketEndHour*60 + marketEndMinute
			currentMinutes := currentHour*60 + currentMinute

			if currentMinutes >= marketStart && currentMinutes <= marketEnd {
				// Get the latest deployments from the event listener
				currentDeploymentsByTimeframe := eventListener.GetDeploymentsByTimeframe()

				// Execute deployments based on their timeframes
				var totalExecutedDeployments int
				var executedTimeframes []string

				for timeframe, deployments := range currentDeploymentsByTimeframe {
					if shouldExecuteTimeframe(timeframe, currentTimestamp) {
						log.Printf("✓ Executing %d deployments for timeframe %s at %v",
							len(deployments), timeframe, currentTimestamp.Format("2006-01-02 15:04:05"))

						for _, dep := range deployments {
							go func(d *deploymentService.DeploymentWithStrategy, timestamp time.Time, tf string) {
								if err := svc.Run(ctx, d, timestamp, nil); err != nil {
									log.Printf("Deployment %v (timeframe %s) run error at %v: %v",
										d.Deployment.ID, tf, timestamp.Format("2006-01-02 15:04:05"), err)
								} else {
									log.Printf("Deployment %v (timeframe %s) executed successfully at %v",
										d.Deployment.ID, tf, timestamp.Format("2006-01-02 15:04:05"))
								}
							}(dep, currentTimestamp, timeframe)
						}
						totalExecutedDeployments += len(deployments)
						executedTimeframes = append(executedTimeframes, timeframe)
					} else {
						// Log skipped timeframes for debugging (only occasionally to avoid spam)
						if currentTimestamp.Minute()%10 == 0 { // Log every 10 minutes
							log.Printf("⏸ Skipping %d deployments for timeframe %s at %v (not scheduled)",
								len(deployments), timeframe, currentTimestamp.Format("15:04:05"))
						}
					}
				}

				if totalExecutedDeployments > 0 {
					log.Printf("📊 EXECUTION SUMMARY at %v: %d deployments across timeframes [%s]",
						currentTimestamp.Format("2006-01-02 15:04:05"), totalExecutedDeployments,
						strings.Join(executedTimeframes, ", "))
				}

				// Increment timestamp by 1 minute
				currentTimestamp = currentTimestamp.Add(1 * time.Minute)
			} else {
				// Market is closed, move to next day at market start time
				nextDay := currentTimestamp.AddDate(0, 0, 1)
				currentTimestamp = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(),
					marketStartHour, marketStartMinute, 0, 0, nextDay.Location())
				log.Printf("Market closed, moving to next day: %v", currentTimestamp)
			}
		}
	}
}

// startCronScheduler runs deployments using cron expressions with actual time
// Used in production environment for real-time deployment execution
func startCronScheduler(ctx context.Context, eventListener *DeploymentEventListener, svc deploymentService.DeploymentService, config *SchedulerConfig) {
	log.Printf("Starting cron scheduler for deployments by timeframe")

	// Create a new cron scheduler
	c := cron.New(cron.WithLocation(time.Local))

	// Get initial deployments from event listener
	deploymentsByTimeframe := eventListener.GetDeploymentsByTimeframe()

	// Add cron jobs based on different timeframes
	for timeframe, deployments := range deploymentsByTimeframe {
		cronExpression := getCronExpressionForTimeframe(timeframe, config)
		timeframeCopy := timeframe // Create a copy for the closure

		log.Printf("Adding cron job for timeframe %s with %d deployments using expression: %s",
			timeframe, len(deployments), cronExpression)

		_, err := c.AddFunc(cronExpression, func() {
			currentTime := time.Now()
			// Reduce 1 second to normalize the timestamp (cron runs with 1-second delay)
			normalizedTime := currentTime.Add(-1 * time.Second)

			// Check if current time is within market hours using config
			if !isMarketHoursWithConfig(currentTime, config) {
				log.Printf("Outside market hours, skipping %s execution at %v", timeframeCopy, currentTime)
				return
			}

			// Get the latest deployments for this timeframe from the event listener
			currentDeploymentsByTimeframe := eventListener.GetDeploymentsByTimeframe()
			currentDeployments := currentDeploymentsByTimeframe[timeframeCopy]

			log.Printf("Running %d deployments for timeframe %s at actual time %v (normalized: %v)",
				len(currentDeployments), timeframeCopy, currentTime, normalizedTime)

			for _, dep := range currentDeployments {
				go func(d *deploymentService.DeploymentWithStrategy, timestamp time.Time) {
					if err := svc.Run(ctx, d, timestamp, nil); err != nil {
						log.Printf("Deployment %v run error at %v: %v", d.Deployment.ID, timestamp, err)
					}
				}(dep, normalizedTime)
			}
		})

		if err != nil {
			log.Printf("Failed to add cron job for timeframe %s: %v", timeframe, err)
			continue
		}
	}

	// Start the cron scheduler
	c.Start()
	log.Println("Cron scheduler started")
	log.Printf("Market hours: %02d:%02d to %02d:%02d, Monday to Friday",
		config.MarketStartHour, config.MarketStartMin,
		config.MarketEndHour, config.MarketEndMin)

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Cron scheduler stopped")
	c.Stop()
}

// getCronExpressionForTimeframe returns appropriate cron expression based on timeframe
func getCronExpressionForTimeframe(timeframe string, config *SchedulerConfig) string {
	// Check for custom cron expression in environment variables first
	customCronKey := fmt.Sprintf("CUSTOM_CRON_%s", strings.ToUpper(timeframe))
	if customCron := getEnv(customCronKey, ""); customCron != "" {
		log.Printf("Using custom cron expression for %s: %s", timeframe, customCron)
		return customCron
	}

	startHour := config.MarketStartHour
	endHour := config.MarketEndHour
	startMin := config.MarketStartMin

	switch timeframe {
	case "1m":
		// Every minute during market hours (configurable hours, Mon-Fri) with 1 second delay
		return fmt.Sprintf("1 */%d %d-%d * * 1-5", 1, startHour, endHour)
	case "3m":
		// Every 3 minutes during market hours with 1 second delay
		return fmt.Sprintf("1 */%d %d-%d * * 1-5", 3, startHour, endHour)
	case "5m":
		// Every 5 minutes during market hours with 1 second delay
		return fmt.Sprintf("1 */%d %d-%d * * 1-5", 5, startHour, endHour)
	case "10m":
		// Every 10 minutes during market hours with 1 second delay
		return fmt.Sprintf("1 */%d %d-%d * * 1-5", 10, startHour, endHour)
	case "15m":
		// Every 15 minutes during market hours with 1 second delay
		return fmt.Sprintf("1 */%d %d-%d * * 1-5", 15, startHour, endHour)
	case "30m":
		// Every 30 minutes during market hours with 1 second delay
		return fmt.Sprintf("1 */%d %d-%d * * 1-5", 30, startHour, endHour)
	case "1h":
		// Every hour during market hours with 1 second delay
		return fmt.Sprintf("1 0 %d-%d * * 1-5", startHour, endHour)
	case "1d":
		// Once per day at market open with 1 second delay
		return fmt.Sprintf("1 %d %d * * 1-5", startMin, startHour)
	default:
		// Default to every minute for unknown timeframes with 1 second delay
		log.Printf("Unknown timeframe %s, defaulting to 1-minute interval", timeframe)
		return fmt.Sprintf("1 */%d %d-%d * * 1-5", 1, startHour, endHour)
	}
}

// isMarketHoursWithConfig checks if the current time is within market trading hours using config
func isMarketHoursWithConfig(currentTime time.Time, config *SchedulerConfig) bool {
	// Skip weekends
	weekday := currentTime.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	currentHour := currentTime.Hour()
	currentMinute := currentTime.Minute()
	currentMinutes := currentHour*60 + currentMinute

	marketStart := config.MarketStartHour*60 + config.MarketStartMin
	marketEnd := config.MarketEndHour*60 + config.MarketEndMin

	return currentMinutes >= marketStart && currentMinutes <= marketEnd
}

// startWebSocketServer starts the HTTP server with WebSocket support for deployment events
func startWebSocketServer(deploymentMod *deploymentModule.Module) {
	// Set Gin mode
	gin.SetMode(getEnv("GIN_MODE", "debug"))

	// Initialize router
	router := gin.Default()

	// Setup CORS
	setupCORS(router)

	// Setup middleware
	setupMiddleware(router)

	// Setup API routes
	apiRouter := router.Group("/api")

	// Register deployment module routes (including WebSocket routes)
	deploymentMod.RegisterRoutes(apiRouter)

	// Setup WebSocket routes with /ws prefix for compatibility
	wsRouter := router.Group("/ws")
	wsApiRouter := wsRouter.Group("/api")
	deploymentMod.RegisterRoutes(wsApiRouter)

	// Add status endpoint
	apiRouter.GET("/deployment/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "running",
			"time":   time.Now().Format(time.RFC3339),
			"websocket_endpoints": gin.H{
				"deployment_events": "/ws/api/deployment/events/ws",
			},
		})
	})

	// Get port from environment or use default
	port := getEnv("DEPLOYMENT_WEBSOCKET_PORT", "")
	if port == "" {
		log.Fatal("DEPLOYMENT_WEBSOCKET_PORT environment variable is required")
	}
	serverAddr := fmt.Sprintf(":%s", port)

	// Create a server with graceful shutdown
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Deployment WebSocket server started on port %s", port)
		log.Printf("WebSocket endpoint available at: ws://localhost:%s/api/deployment/events/ws", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s", err)
		}
	}()

	// The server will be shut down when the main context is cancelled
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

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getSchedulerConfig returns scheduler configuration based on environment
func getSchedulerConfig() *SchedulerConfig {
	environment := getEnv("APP_ENV", EnvDevelopment)

	// Parse market hours from environment or use defaults
	marketStartHour := parseIntEnv("MARKET_START_HOUR", 9)
	marketStartMin := parseIntEnv("MARKET_START_MIN", 15)
	marketEndHour := parseIntEnv("MARKET_END_HOUR", 15)
	marketEndMin := parseIntEnv("MARKET_END_MIN", 30)

	config := &SchedulerConfig{
		Environment:     environment,
		MarketStartHour: marketStartHour,
		MarketStartMin:  marketStartMin,
		MarketEndHour:   marketEndHour,
		MarketEndMin:    marketEndMin,
		SimulationSpeed: 1 * time.Second, // Default: 1 minute per second
	}

	// Override configuration based on environment
	switch environment {
	case EnvDevelopment:
		config.UseSimulator = true
		// Allow overriding simulation speed via environment
		if speedStr := getEnv("SIMULATION_SPEED", ""); speedStr != "" {
			if speed, err := time.ParseDuration(speedStr); err == nil {
				config.SimulationSpeed = speed
				log.Printf("Using custom simulation speed: %v", speed)
			} else {
				log.Printf("Invalid SIMULATION_SPEED format: %s, using default", speedStr)
			}
		}
	case EnvProduction, EnvStaging:
		config.UseSimulator = false
		// In production/staging, we can still override to use simulator if needed
		if getEnv("FORCE_SIMULATOR", "false") == "true" {
			config.UseSimulator = true
			log.Println("FORCE_SIMULATOR enabled - using simulator in production environment")
		}
	}

	log.Printf("Scheduler config: Environment=%s, UseSimulator=%t, MarketHours=%02d:%02d-%02d:%02d, SimulationSpeed=%v",
		config.Environment, config.UseSimulator,
		config.MarketStartHour, config.MarketStartMin,
		config.MarketEndHour, config.MarketEndMin,
		config.SimulationSpeed)

	return config
}

// parseIntEnv parses an integer from environment variable with fallback to default
func parseIntEnv(key string, defaultValue int) int {
	if strVal := getEnv(key, ""); strVal != "" {
		if intVal, err := strconv.Atoi(strVal); err == nil {
			return intVal
		} else {
			log.Printf("Invalid integer value for %s: %s, using default %d", key, strVal, defaultValue)
		}
	}
	return defaultValue
}

// DeploymentEventListener listens for deployment events and updates the scheduler
type DeploymentEventListener struct {
	deploymentsByTimeframe map[string][]*deploymentService.DeploymentWithStrategy
	deploymentService      deploymentService.DeploymentService
	mutex                  sync.RWMutex
}

// NewDeploymentEventListener creates a new deployment event listener
func NewDeploymentEventListener(
	deploymentsByTimeframe map[string][]*deploymentService.DeploymentWithStrategy,
	deploymentSvc deploymentService.DeploymentService,
) *DeploymentEventListener {
	return &DeploymentEventListener{
		deploymentsByTimeframe: deploymentsByTimeframe,
		deploymentService:      deploymentSvc,
	}
}

// HandleDeploymentCreated handles deployment.created events
func (l *DeploymentEventListener) HandleDeploymentCreated(ctx context.Context, event *dto.DeploymentEvent) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Extract deployment created data
	deploymentData, ok := event.Data.(*dto.DeploymentCreatedData)
	if !ok {
		log.Printf("Invalid deployment created event data")
		return fmt.Errorf("invalid deployment created event data")
	}

	// Get the deployment details to validate it exists
	_, err := l.deploymentService.GetDeploymentByID(ctx, deploymentData.DeploymentID, event.Metadata.UserID)
	if err != nil {
		log.Printf("Failed to get deployment: %v", err)
		return err
	}

	// We need to create a DeploymentWithStrategy manually since we don't have a direct method
	// For now, we'll re-fetch all active deployments to get the proper structure
	// This is not the most efficient, but ensures we get the correct data structure
	allDeployments, err := l.deploymentService.GetActiveDeploymentsWithStrategies(ctx)
	if err != nil {
		log.Printf("Failed to get active deployments: %v", err)
		return err
	}

	// Find the newly created deployment in the list
	var newDeployment *deploymentService.DeploymentWithStrategy
	for _, dep := range allDeployments {
		if dep.Deployment.ID == deploymentData.DeploymentID {
			newDeployment = dep
			break
		}
	}

	if newDeployment == nil {
		log.Printf("Newly created deployment not found in active deployments list")
		return fmt.Errorf("newly created deployment not found")
	}

	// Add to the appropriate timeframe list
	timeframe := deploymentData.Timeframe
	if timeframe != "" {
		// Check if deployment already exists in the list to avoid duplicates
		exists := false
		for _, existing := range l.deploymentsByTimeframe[timeframe] {
			if existing.Deployment.ID == deploymentData.DeploymentID {
				exists = true
				break
			}
		}

		if !exists {
			l.deploymentsByTimeframe[timeframe] = append(l.deploymentsByTimeframe[timeframe], newDeployment)
			log.Printf("Added new deployment %v to timeframe %s (total: %d deployments)",
				deploymentData.DeploymentID, timeframe, len(l.deploymentsByTimeframe[timeframe]))
		} else {
			log.Printf("Deployment %v already exists in timeframe %s", deploymentData.DeploymentID, timeframe)
		}
	}

	return nil
}

// GetDeploymentsByTimeframe returns a copy of the current deployments by timeframe
func (l *DeploymentEventListener) GetDeploymentsByTimeframe() map[string][]*deploymentService.DeploymentWithStrategy {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[string][]*deploymentService.DeploymentWithStrategy)
	for timeframe, deployments := range l.deploymentsByTimeframe {
		result[timeframe] = make([]*deploymentService.DeploymentWithStrategy, len(deployments))
		copy(result[timeframe], deployments)
	}
	return result
}

// startWebSocketListener starts a WebSocket client to listen for deployment events
func startWebSocketListener(ctx context.Context, eventListener *DeploymentEventListener) { // Get WebSocket server details
	wsPort := getEnv("DEPLOYMENT_WEBSOCKET_PORT", "")
	if wsPort == "" {
		log.Printf("DEPLOYMENT_WEBSOCKET_PORT not set, WebSocket listener disabled")
		return
	}

	// TODO: Implement actual WebSocket client connection to ws://localhost:{wsPort}/api/deployment/events/ws
	// For now, we'll implement a simple polling mechanism to check for new deployments
	// In a full implementation, you would establish a WebSocket connection and listen for events
	// This is a simplified approach that periodically refreshes the deployment list

	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	log.Printf("WebSocket listener started, will poll for deployment changes every 30 seconds")

	for {
		select {
		case <-ctx.Done():
			log.Println("WebSocket listener stopped")
			return
		case <-ticker.C:
			// Periodically refresh the deployment list to pick up new deployments
			log.Printf("Refreshing deployment list from database...")

			// Get fresh deployments from the database
			freshDeployments, err := getActiveDeploymentsByTimeframe(ctx, eventListener.deploymentService)
			if err != nil {
				log.Printf("Failed to refresh deployments: %v", err)
				continue
			}

			// Update the event listener's deployment map
			eventListener.mutex.Lock()
			eventListener.deploymentsByTimeframe = freshDeployments
			eventListener.mutex.Unlock()

			// Log the current deployment counts
			for timeframe, deployments := range freshDeployments {
				log.Printf("Timeframe %s: %d deployments", timeframe, len(deployments))
			}
		}
	}
}

// shouldExecuteTimeframe determines if deployments for a specific timeframe should execute at the given timestamp
// This function implements the timeframe logic to ensure deployments run at appropriate intervals
//
// Timeframe execution logic:
// - 1m: Every minute (always true)
// - 3m: At minutes 0, 3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 33, 36, 39, 42, 45, 48, 51, 54, 57
// - 5m: At minutes 0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55
// - 10m: At minutes 0, 10, 20, 30, 40, 50
// - 15m: At minutes 0, 15, 30, 45
// - 30m: At minutes 0, 30
// - 1h: At minute 0 of each hour
// - 3h: At 9:00, 12:00, 15:00 (market hours aligned)
// - 1d: Once per day at market open (9:15)
func shouldExecuteTimeframe(timeframe string, currentTimestamp time.Time) bool {
	currentMinute := currentTimestamp.Minute()
	currentHour := currentTimestamp.Hour()

	switch timeframe {
	case "1m":
		// Execute every minute
		return true

	case "3m":
		// Execute every 3 minutes - aligned to market start (9:15)
		// Minutes since market start: current_minute - 15
		// For hours after 9, we need to account for the full hour cycles
		minutesSinceMarketStart := ((currentHour-9)*60 + currentMinute - 15)
		if currentHour == 9 {
			// First hour: start from minute 15, then 18, 21, etc.
			minutesSinceMarketStart = currentMinute - 15
		}
		// Execute if it's exactly at market start or at 3-minute intervals thereafter
		return minutesSinceMarketStart >= 0 && minutesSinceMarketStart%3 == 0

	case "5m":
		// Execute every 5 minutes - aligned to market start (9:15)
		minutesSinceMarketStart := ((currentHour-9)*60 + currentMinute - 15)
		if currentHour == 9 {
			minutesSinceMarketStart = currentMinute - 15
		}
		return minutesSinceMarketStart >= 0 && minutesSinceMarketStart%5 == 0

	case "10m":
		// Execute every 10 minutes - aligned to market start (9:15)
		minutesSinceMarketStart := ((currentHour-9)*60 + currentMinute - 15)
		if currentHour == 9 {
			minutesSinceMarketStart = currentMinute - 15
		}
		return minutesSinceMarketStart >= 0 && minutesSinceMarketStart%10 == 0

	case "15m":
		// Execute every 15 minutes - aligned to market start (9:15)
		minutesSinceMarketStart := ((currentHour-9)*60 + currentMinute - 15)
		if currentHour == 9 {
			minutesSinceMarketStart = currentMinute - 15
		}
		return minutesSinceMarketStart >= 0 && minutesSinceMarketStart%15 == 0

	case "30m":
		// Execute every 30 minutes - aligned to market start (9:15)
		minutesSinceMarketStart := ((currentHour-9)*60 + currentMinute - 15)
		if currentHour == 9 {
			minutesSinceMarketStart = currentMinute - 15
		}
		return minutesSinceMarketStart >= 0 && minutesSinceMarketStart%30 == 0

	case "1h":
		// Execute every hour at minute 15 (aligned to market pattern)
		// So at 9:15, 10:15, 11:15, 12:15, 13:15, 14:15, 15:15
		return currentMinute == 15

	case "3h":
		// Execute every 3 hours during market session
		// At 9:15 (market open), 12:15 (midday), 15:15 (near close)
		return currentMinute == 15 && (currentHour == 9 || currentHour == 12 || currentHour == 15)

	case "1d":
		// Execute once per day at market open (9:15)
		return currentHour == 9 && currentMinute == 15

	default:
		// For unknown timeframes, default to every minute execution
		log.Printf("⚠️ Unknown timeframe %s, defaulting to 1-minute execution", timeframe)
		return true
	}
}

// getTimeframeScheduleDescription returns a human-readable description of when a timeframe executes
func getTimeframeScheduleDescription(timeframe string) string {
	switch timeframe {
	case "1m":
		return "executes every minute"
	case "3m":
		return "executes every 3 minutes (9:15, 9:18, 9:21, ...)"
	case "5m":
		return "executes every 5 minutes (9:15, 9:20, 9:25, ...)"
	case "10m":
		return "executes every 10 minutes (9:15, 9:25, 9:35, ...)"
	case "15m":
		return "executes every 15 minutes (9:15, 9:30, 9:45, ...)"
	case "30m":
		return "executes every 30 minutes (9:15, 9:45, 10:15, ...)"
	case "1h":
		return "executes every hour (9:15, 10:15, 11:15, ...)"
	case "3h":
		return "executes every 3 hours (9:15, 12:15, 15:15)"
	case "1d":
		return "executes once per day at market open (9:15)"
	default:
		return "executes every minute (unknown timeframe)"
	}
}

// getNextExecutionTime calculates the next execution time for a given timeframe from the start time
func getNextExecutionTime(timeframe string, startTime time.Time) time.Time {
	switch timeframe {
	case "1m":
		return startTime // Always executes immediately
	case "3m":
		// Find next 3-minute interval from market start
		minutes := startTime.Minute()
		if minutes == 15 {
			return startTime // Exactly at market start
		}
		// Calculate minutes since market start and find next interval
		minutesSinceStart := minutes - 15
		if minutesSinceStart < 0 {
			minutesSinceStart += 60 // Handle hour boundary
		}
		nextInterval := ((minutesSinceStart / 3) + 1) * 3
		return startTime.Add(time.Duration(nextInterval-minutesSinceStart) * time.Minute)
	case "5m", "10m", "15m", "30m":
		// Similar logic for other minute-based intervals
		interval := 5
		if timeframe == "10m" {
			interval = 10
		} else if timeframe == "15m" {
			interval = 15
		} else if timeframe == "30m" {
			interval = 30
		}

		minutes := startTime.Minute()
		if minutes == 15 {
			return startTime
		}
		minutesSinceStart := minutes - 15
		if minutesSinceStart < 0 {
			minutesSinceStart += 60
		}
		nextInterval := ((minutesSinceStart / interval) + 1) * interval
		return startTime.Add(time.Duration(nextInterval-minutesSinceStart) * time.Minute)
	case "1h":
		// Next hour at minute 15
		if startTime.Minute() <= 15 {
			return time.Date(startTime.Year(), startTime.Month(), startTime.Day(),
				startTime.Hour(), 15, 0, 0, startTime.Location())
		}
		return time.Date(startTime.Year(), startTime.Month(), startTime.Day(),
			startTime.Hour()+1, 15, 0, 0, startTime.Location())
	case "3h":
		// Next 3-hour interval (9:15, 12:15, 15:15)
		hour := startTime.Hour()
		if hour < 9 || (hour == 9 && startTime.Minute() <= 15) {
			return time.Date(startTime.Year(), startTime.Month(), startTime.Day(),
				9, 15, 0, 0, startTime.Location())
		} else if hour < 12 || (hour == 12 && startTime.Minute() <= 15) {
			return time.Date(startTime.Year(), startTime.Month(), startTime.Day(),
				12, 15, 0, 0, startTime.Location())
		} else if hour < 15 || (hour == 15 && startTime.Minute() <= 15) {
			return time.Date(startTime.Year(), startTime.Month(), startTime.Day(),
				15, 15, 0, 0, startTime.Location())
		}
		// Next day
		nextDay := startTime.AddDate(0, 0, 1)
		return time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(),
			9, 15, 0, 0, nextDay.Location())
	case "1d":
		// Next market open
		if startTime.Hour() < 9 || (startTime.Hour() == 9 && startTime.Minute() <= 15) {
			return time.Date(startTime.Year(), startTime.Month(), startTime.Day(),
				9, 15, 0, 0, startTime.Location())
		}
		// Next day
		nextDay := startTime.AddDate(0, 0, 1)
		return time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(),
			9, 15, 0, 0, nextDay.Location())
	default:
		return startTime
	}
}

// Additional environment configuration options that can be added to .env file:
//
// # Basic Configuration
// APP_ENV=development
// DEPLOYMENT_WEBSOCKET_PORT=9091
//
// # Simulator Configuration (Development)
// SIMULATION_SPEED=1s
// FORCE_SIMULATOR=false
//
// # Market Hours Configuration
// MARKET_START_HOUR=9
// MARKET_START_MIN=15
// MARKET_END_HOUR=15
// MARKET_END_MIN=30
//
// # Custom Cron Expressions (Production)
// CUSTOM_CRON_1M="*/1 9-15 * * 1-5"
// CUSTOM_CRON_5M="*/5 9-15 * * 1-5"
// CUSTOM_CRON_15M="*/15 9-15 * * 1-5"
// CUSTOM_CRON_1H="0 9-15 * * 1-5"
// CUSTOM_CRON_1D="15 9 * * 1-5"
//
// Example usage:
// - Development: Fast simulation for testing strategies
// - Production: Real-time execution during market hours
// - Staging: Real-time but with FORCE_SIMULATOR=true for testing
//
// WebSocket Integration:
// - When a deployment is created via API, it publishes a 'deployment.created' event
// - The deployment scheduler listens for these events and dynamically updates its deployment list
// - This enables real-time addition of new deployments without requiring a restart
// - The scheduler polls every 30 seconds to sync with the database for reliability
