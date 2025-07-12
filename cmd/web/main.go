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

	"algo/internal/modules/backtest"
	"algo/internal/modules/community"
	"algo/internal/modules/deployment"
	"algo/internal/modules/instrument"
	"algo/internal/modules/marketplace"
	"algo/internal/modules/oauth"
	"algo/internal/modules/strategy"
	"algo/internal/modules/user"
	userService "algo/internal/modules/user/application/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	// Setup API routes
	apiRouter := router.Group("/api")

	// Initialize modules
	initializeModules(apiRouter)

	// Get port from environment or use default
	port := getEnv("WEB_HTTP_PORT", "")
	if port == "" {
		log.Fatal("WEB_HTTP_PORT environment variable is required")
	}
	serverAddr := fmt.Sprintf(":%s", port)

	// Create a server with graceful shutdown
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server started on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
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

// initializeModules initializes and connects all application modules
func initializeModules(router *gin.RouterGroup) {
	// Initialize user module first as other modules depend on its controllers for auth
	userModule := user.New()

	// Get the user controller to pass to other modules for auth
	userController := userModule.GetController()

	// Initialize instrument module
	instrumentModule := instrument.New()
	instrumentModule.RegisterRoutes(router)

	// Initialize strategy module
	strategyModule := strategy.New()
	// Set user controller to enable auth middleware for all strategy endpoints
	// This ensures that all strategy endpoints are protected by authentication
	strategyModule.SetUserController(userController)
	strategyModule.RegisterRoutes(router)

	// Initialize deployment module
	deploymentModule := deployment.New()
	deploymentModule.Initialize(userController, strategyModule)
	deploymentModule.RegisterRoutes(router)

	// Initialize backtest module
	backtestModule := backtest.New()
	backtestModule.SetUserController(userController)
	backtestModule.RegisterRoutes(router)

	// Initialize marketplace module
	marketplaceModule := marketplace.New()
	marketplaceModule.SetUserController(userController)
	marketplaceModule.RegisterRoutes(router)

	// Initialize community module
	communityModule := community.New()
	communityModule.SetUserController(userController)
	communityModule.RegisterRoutes(router)

	// Initialize oauth module
	oauthModule := oauth.New()
	oauthModule.Initialize(userController)
	oauthModule.RegisterRoutes(router)

	// Create and set dashboard service after all modules are initialized but before user routes
	dashboardService := userService.NewDashboardService(
		deploymentModule.GetDeploymentService(),
		strategyModule.GetStrategyService(),
		backtestModule.GetBacktestService(),
		marketplaceModule.GetSellerService(),
		communityModule.GetUserService(),
	)
	userModule.SetDashboardService(dashboardService)

	// Register user routes after dashboard service is set
	userModule.RegisterRoutes(router)

	// Additional modules will be initialized here
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
