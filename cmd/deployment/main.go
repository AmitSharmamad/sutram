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

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	deploymentModule "algo/internal/modules/deployment"
	deploymentService "algo/internal/modules/deployment/application/service"
	strategyModule "algo/internal/modules/strategy"
	userModule "algo/internal/modules/user"
)

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

	// Start custom scheduler for simulated time progression
	go startSimulatedTimeScheduler(ctx, deploymentsByTimeframe, deploymentSvc)

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
func startSimulatedTimeScheduler(ctx context.Context, deploymentsByTimeframe map[string][]*deploymentService.DeploymentWithStrategy, svc deploymentService.DeploymentService) {
	// Start date: May 5, 2025 at 9:15:00 IST
	startDate := time.Date(2025, 5, 5, 9, 15, 0, 0, time.Local)
	currentTimestamp := startDate

	// Market hours: 9:15:00 to 15:30:00 (3:30 PM)
	marketStartHour, marketStartMinute := 9, 15
	marketEndHour, marketEndMinute := 15, 30

	ticker := time.NewTicker(1 * time.Second) // Execute every second
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
						if err := svc.Run(ctx, d, timestamp); err != nil {
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
