package marketplace

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/google/wire"

	"algo/src/core/db"
)

// Root represents the marketplace root component
type Root struct{}

// NewRootProv creates a new marketplace root component
func NewRootProv(router *gin.RouterGroup) *Root {
	return &Root{}
}

// MockAuthMiddleware is a temporary middleware to bypass authentication
func MockAuthMiddleware() gin.HandlerFunc {
	// Create a fixed mock user ID for testing
	mockUserID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")

	return func(c *gin.Context) {
		// Set a mock user ID in the context
		c.Set("userID", mockUserID)
		c.Next()
	}
}

// NewRepoProv provides a new marketplace repository
func NewRepoProv(factory *db.DBFactory) Repository {
	// Use the factory to get a database connection
	dbConn := factory.Get("postgres").GetDB()
	return NewRepository(dbConn)
}

// NewServiceProv provides a new marketplace service
func NewServiceProv(repo Repository) Service {
	return NewService(repo)
}

// NewControllerProv provides a new marketplace controller
func NewControllerProv(service Service) *Controller {
	return NewController(service)
}

// NewRoutesProv registers marketplace routes with the engine
func NewRoutesProv(engine *gin.Engine, controller *Controller) *gin.RouterGroup {
	// Create a router group for the API
	apiRouter := engine.Group("")

	// Apply mock auth middleware to the API router
	apiRouter.Use(MockAuthMiddleware())

	// Register routes
	RegisterRoutes(apiRouter, controller)

	return apiRouter
}

// Set is the marketplace dependency injection set
var Set = wire.NewSet(
	NewRepoProv,
	NewServiceProv,
	NewControllerProv,
	NewRoutesProv,
	NewRootProv,
)
