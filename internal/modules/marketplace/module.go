// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/module.go
package marketplace

import (
	"database/sql"

	infraRepo "algo/internal/infra/persistance/repository"
	"algo/internal/modules/marketplace/application/service"
	"algo/internal/modules/marketplace/infrastructure/persistence/postgres"
	"algo/internal/modules/marketplace/interfaces/http"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Module represents the marketplace module
type Module struct {
	db             *sql.DB
	controller     *http.MarketplaceController
	userController interface{}
}

// ProvideModule provides a wire provider set for the marketplace module
var ProvideModule = wire.NewSet(
	ProvideMarketplaceModule,
)

// ProvideMarketplaceModule provides a marketplace module instance
func ProvideMarketplaceModule() (*Module, error) {
	return New(), nil
}

// New creates a new marketplace module
func New() *Module {
	// Get database connection
	db, err := infraRepo.ProvidePostgresDB()
	if err != nil {
		panic(err) // In a real app, handle this more gracefully
	}

	module := &Module{
		db: db,
	}

	// Initialize module components
	module.initialize()

	return module
}

// initialize sets up the marketplace module
func (m *Module) initialize() {
	// Create repositories
	listingRepo := postgres.NewListingRepository(m.db)
	categoryRepo := postgres.NewCategoryRepository(m.db)
	subscriptionRepo := postgres.NewSubscriptionRepository(m.db)
	reviewRepo := postgres.NewReviewRepository(m.db)
	transactionRepo := postgres.NewTransactionRepository(m.db)
	performanceRepo := postgres.NewPerformanceRepository(m.db)

	// Create services
	listingService := service.NewListingService(listingRepo, categoryRepo, performanceRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, listingRepo, transactionRepo)
	reviewService := service.NewReviewService(reviewRepo, listingRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	sellerService := service.NewSellerService(listingRepo, subscriptionRepo, transactionRepo)

	// Create controller
	m.controller = http.NewMarketplaceController(
		listingService,
		categoryService,
		subscriptionService,
		reviewService,
		transactionService,
		sellerService,
	)
}

// RegisterRoutes registers the marketplace module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	// Cast userController to the expected type
	userController, ok := m.userController.(*userHttp.UserController)
	if !ok {
		// If userController is not set or is of wrong type, we can't proceed with auth
		// In a production app, you'd want to log this error properly
		panic("Marketplace module: UserController not set or invalid type")
	}

	http.RegisterRoutes(router, m.controller, userController)
}

// GetController returns the marketplace controller
func (m *Module) GetController() *http.MarketplaceController {
	return m.controller
}

// SetUserController sets the user controller for auth middleware
func (m *Module) SetUserController(userController interface{}) {
	m.userController = userController
}
