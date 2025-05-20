package user

import (
	"database/sql"
	"os"
	"time"

	infraRepo "algo/internal/infra/persistance/repository"
	"algo/internal/modules/user/application/service"
	"algo/internal/modules/user/application/service/auth"
	"algo/internal/modules/user/infrastructure/persistence/postgres"
	"algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Module represents the user module
type Module struct {
	db         *sql.DB
	jwtSecret  string
	controller *http.UserController
}

// ProvideModule provides a wire provider set for the user module
var ProvideModule = wire.NewSet(
	ProvideUserModule,
)

// ProvideJWTSecret provides the JWT secret key
func ProvideJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-jwt-secret-for-development-only"
	}
	return secret
}

// ProvideUserModule provides a user module instance
func ProvideUserModule() (*Module, error) {
	return New(), nil
}

// New creates a new user module
func New() *Module {
	// Get database connection
	db, err := infraRepo.ProvidePostgresDB()
	if err != nil {
		panic(err) // In a real app, handle this more gracefully
	}

	// Get JWT secret
	jwtSecret := ProvideJWTSecret()

	module := &Module{
		db:        db,
		jwtSecret: jwtSecret,
	}

	// Initialize module components
	module.initialize()

	return module
}

// initialize sets up the user module
func (m *Module) initialize() {
	// Create repositories
	userRepo := postgres.NewUserRepository(m.db)
	roleRepo := postgres.NewRoleRepository(m.db)

	// Create services
	userService := service.NewUserService(userRepo, roleRepo)
	jwtService := auth.NewJWTService(m.jwtSecret, 24*time.Hour) // 24-hour token validity
	authService := service.NewAuthService(userService, jwtService)

	// Create controller
	m.controller = http.NewUserController(authService, userService)
}

// RegisterRoutes registers the user module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	http.RegisterRoutes(router, m.controller)
}

// GetController returns the user controller
func (m *Module) GetController() *http.UserController {
	return m.controller
}
