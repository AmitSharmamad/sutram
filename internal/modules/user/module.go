package user

import (
	"database/sql"
	"time"

	"algo/internal/modules/user/application/service"
	"algo/internal/modules/user/application/service/auth"
	"algo/internal/modules/user/infrastructure/persistence/postgres"
	"algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
)

// Module represents the user module
type Module struct {
	db         *sql.DB
	jwtSecret  string
	router     *gin.RouterGroup
	controller *http.UserController
}

// NewModule creates a new user module
func NewModule(db *sql.DB, jwtSecret string, router *gin.RouterGroup) *Module {
	return &Module{
		db:        db,
		jwtSecret: jwtSecret,
		router:    router,
	}
}

// Initialize sets up the user module
func (m *Module) Initialize() {
	// Create repositories
	userRepo := postgres.NewUserRepository(m.db)
	roleRepo := postgres.NewRoleRepository(m.db)

	// Create services
	userService := service.NewUserService(userRepo, roleRepo)
	jwtService := auth.NewJWTService(m.jwtSecret, 24*time.Hour) // 24-hour token validity
	authService := service.NewAuthService(userService, jwtService)

	// Create controller
	m.controller = http.NewUserController(authService, userService)

	// Register routes
	http.RegisterRoutes(m.router, m.controller)
}

// GetController returns the user controller
func (m *Module) GetController() *http.UserController {
	return m.controller
}
