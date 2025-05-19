package oauth

import (
	"database/sql"

	infraRepo "algo/internal/infra/persistance/repository"
	"algo/internal/modules/oauth/application/service"
	"algo/internal/modules/oauth/infrastructure/persistence/postgres"
	"algo/internal/modules/oauth/interfaces/http"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Module represents the OAuth module
type Module struct {
	db             *sql.DB
	controller     *http.OAuthController
	userController *userHttp.UserController
}

// ProvideModule provides a wire provider set for the OAuth module
var ProvideModule = wire.NewSet(
	ProvideOAuthModule,
)

// ProvideOAuthModule provides an OAuth module instance
func ProvideOAuthModule() (*Module, error) {
	return New(), nil
}

// New creates a new OAuth module
func New() *Module {
	// Get database connection
	db, err := infraRepo.ProvidePostgresDB()
	if err != nil {
		panic(err) // In a real app, handle this more gracefully
	}

	module := &Module{
		db: db,
	}

	return module
}

// Initialize initializes the OAuth module with the user controller
func (m *Module) Initialize(userController *userHttp.UserController) {
	m.userController = userController

	// Initialize module components
	m.initializeComponents()
}

// initializeComponents sets up the OAuth module components
func (m *Module) initializeComponents() {
	// Create repositories
	providerRepo := postgres.NewOAuthProviderRepository(m.db)
	userOAuthRepo := postgres.NewUserOAuthRepository(m.db)

	// Get user service and auth service from user controller
	userService := m.userController.GetUserService()
	authService := m.userController.GetAuthService()

	// Create services
	oauthService := service.NewOAuthService(providerRepo, userOAuthRepo, userService, authService)

	// Create controller
	m.controller = http.NewOAuthController(oauthService)
}

// RegisterRoutes registers the OAuth module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	http.RegisterRoutes(router, m.controller, m.userController)
}

// GetController returns the OAuth controller
func (m *Module) GetController() *http.OAuthController {
	return m.controller
}
