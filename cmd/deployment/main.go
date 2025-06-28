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
	"syscall"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	deploymentModule "algo/internal/modules/deployment"
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
		go startSimulatedTimeScheduler(ctx, deploymentsByTimeframe, deploymentSvc, config)
	} else {
		// Start cron-based scheduler for production environment
		log.Println("Starting CRON-BASED SCHEDULER")
		log.Println("Features: Real-time execution based on actual market hours")
		log.Printf("Market hours: %02d:%02d to %02d:%02d",
			config.MarketStartHour, config.MarketStartMin,
			config.MarketEndHour, config.MarketEndMin)
		go startCronScheduler(ctx, deploymentsByTimeframe, deploymentSvc, config)
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
func startSimulatedTimeScheduler(ctx context.Context, deploymentsByTimeframe map[string][]*deploymentService.DeploymentWithStrategy, svc deploymentService.DeploymentService, config *SchedulerConfig) {
	// Start date: May 5, 2025 at 9:15:00 IST
	startDate := time.Date(2025, 5, 5, 9, 15, 0, 0, time.Local)
	currentTimestamp := startDate

	marketStartHour, marketStartMinute := config.MarketStartHour, config.MarketStartMin
	marketEndHour, marketEndMinute := config.MarketEndHour, config.MarketEndMin

	ticker := time.NewTicker(config.SimulationSpeed) // Execute based on simulation speed
	defer ticker.Stop()

	log.Printf("Starting simulated time scheduler from %v", startDate)
	log.Printf("Market hours: %02d:%02d to %02d:%02d", marketStartHour, marketStartMinute, marketEndHour, marketEndMinute)

	// Collect all deployments from all timeframes
	var allDeployments []*deploymentService.DeploymentWithStrategy
	for _, deployments := range deploymentsByTimeframe {
		allDeployments = append(allDeployments, deployments...)
	}

	log.Printf("Total deployments to run: %d", len(allDeployments))

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
				// Run deployments for current timestamp
				log.Printf("Running %d deployments at simulated time %v", len(allDeployments), currentTimestamp)

				for _, dep := range allDeployments {
					go func(d *deploymentService.DeploymentWithStrategy, timestamp time.Time) {
						if err := svc.Run(ctx, d, timestamp, nil); err != nil {
							log.Printf("Deployment %v run error at %v: %v", d.Deployment.ID, timestamp, err)
						}
					}(dep, currentTimestamp)
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
func startCronScheduler(ctx context.Context, deploymentsByTimeframe map[string][]*deploymentService.DeploymentWithStrategy, svc deploymentService.DeploymentService, config *SchedulerConfig) {
	log.Printf("Starting cron scheduler for deployments by timeframe")

	// Create a new cron scheduler
	c := cron.New(cron.WithLocation(time.Local))

	// Add cron jobs based on different timeframes
	for timeframe, deployments := range deploymentsByTimeframe {
		cronExpression := getCronExpressionForTimeframe(timeframe, config)
		deploymentsCopy := deployments // Create a copy for the closure
		timeframeCopy := timeframe     // Create a copy for the closure

		log.Printf("Adding cron job for timeframe %s with %d deployments using expression: %s",
			timeframe, len(deployments), cronExpression)

		_, err := c.AddFunc(cronExpression, func() {
			currentTime := time.Now()

			// Check if current time is within market hours using config
			if !isMarketHoursWithConfig(currentTime, config) {
				log.Printf("Outside market hours, skipping %s execution at %v", timeframeCopy, currentTime)
				return
			}

			log.Printf("Running %d deployments for timeframe %s at actual time %v",
				len(deploymentsCopy), timeframeCopy, currentTime)

			for _, dep := range deploymentsCopy {
				go func(d *deploymentService.DeploymentWithStrategy, timestamp time.Time) {
					if err := svc.Run(ctx, d, timestamp, nil); err != nil {
						log.Printf("Deployment %v run error at %v: %v", d.Deployment.ID, timestamp, err)
					}
				}(dep, currentTime)
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
		// Every minute during market hours (configurable hours, Mon-Fri)
		return fmt.Sprintf("*/%d %d-%d * * 1-5", 1, startHour, endHour)
	case "3m":
		// Every 3 minutes during market hours
		return fmt.Sprintf("*/%d %d-%d * * 1-5", 3, startHour, endHour)
	case "5m":
		// Every 5 minutes during market hours
		return fmt.Sprintf("*/%d %d-%d * * 1-5", 5, startHour, endHour)
	case "10m":
		// Every 10 minutes during market hours
		return fmt.Sprintf("*/%d %d-%d * * 1-5", 10, startHour, endHour)
	case "15m":
		// Every 15 minutes during market hours
		return fmt.Sprintf("*/%d %d-%d * * 1-5", 15, startHour, endHour)
	case "30m":
		// Every 30 minutes during market hours
		return fmt.Sprintf("*/%d %d-%d * * 1-5", 30, startHour, endHour)
	case "1h":
		// Every hour during market hours
		return fmt.Sprintf("0 %d-%d * * 1-5", startHour, endHour)
	case "1d":
		// Once per day at market open
		return fmt.Sprintf("%d %d * * 1-5", startMin, startHour)
	default:
		// Default to every minute for unknown timeframes
		log.Printf("Unknown timeframe %s, defaulting to 1-minute interval", timeframe)
		return fmt.Sprintf("*/%d %d-%d * * 1-5", 1, startHour, endHour)
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
	port := getEnv("DEPLOYMENT_PORT", "9091")
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

// Additional environment configuration options that can be added to .env file:
//
// # Basic Configuration
// APP_ENV=development
// DEPLOYMENT_PORT=9091
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
