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
	db               *sql.DB
	jwtSecret        string
	controller       *http.UserController
	dashboardService *service.DashboardService
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
	// Get sqlx connection for waitlist repository
	sqlxDB, err := infraRepo.ProvidePostgresSQLXDB()
	if err != nil {
		panic(err) // In a real app, handle this more gracefully
	}

	// Create repositories
	userRepo := postgres.NewUserRepository(m.db)
	roleRepo := postgres.NewRoleRepository(m.db)
	waitlistRepo := postgres.NewWaitlistRepository(sqlxDB)
	surveyRepo := postgres.NewSurveyRepository(sqlxDB)
	usageLimitRepo := postgres.NewUsageLimitRepository(m.db)
	usageLogRepo := postgres.NewUsageLogRepository(m.db)

	// Create services
	userService := service.NewUserService(userRepo, roleRepo, usageLimitRepo, usageLogRepo)
	surveyService := service.NewSurveyService(surveyRepo)
	waitlistService := service.NewWaitlistService(waitlistRepo, surveyRepo)
	usageLimitService := service.NewUsageLimitService(usageLimitRepo, usageLogRepo)
	jwtService := auth.NewJWTService(m.jwtSecret, 24*time.Hour) // 24-hour token validity
	authService := service.NewAuthService(userService, jwtService)

	// Create controller (dashboard service will be set later)
	m.controller = http.NewUserController(authService, userService, nil)
	m.controller.SetWaitlistService(waitlistService)
	m.controller.SetSurveyService(surveyService)
	m.controller.SetUsageLimitService(usageLimitService)
}

// RegisterRoutes registers the user module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	http.RegisterRoutes(router, m.controller)
}

// GetController returns the user controller
func (m *Module) GetController() *http.UserController {
	return m.controller
}

// GetUsageLimitService returns the usage limit service through the controller
func (m *Module) GetUsageLimitService() service.UsageLimitService {
	return m.controller.GetUsageLimitService()
}

// SetDashboardService sets the dashboard service and updates the controller
func (m *Module) SetDashboardService(dashboardService *service.DashboardService) {
	m.dashboardService = dashboardService
	// Recreate the controller with the dashboard service
	authService := m.controller.GetAuthService()
	userService := m.controller.GetUserService()
	usageLimitService := m.controller.GetUsageLimitService()

	// Get sqlx connection for waitlist repository
	sqlxDB, err := infraRepo.ProvidePostgresSQLXDB()
	if err != nil {
		panic(err) // In a real app, handle this more gracefully
	}

	// Recreate waitlist service
	waitlistRepo := postgres.NewWaitlistRepository(sqlxDB)
	surveyRepo := postgres.NewSurveyRepository(sqlxDB)
	waitlistService := service.NewWaitlistService(waitlistRepo, surveyRepo)
	surveyService := service.NewSurveyService(surveyRepo)

	m.controller = http.NewUserController(authService, userService, dashboardService)
	m.controller.SetWaitlistService(waitlistService)
	m.controller.SetSurveyService(surveyService)
	m.controller.SetUsageLimitService(usageLimitService)
}
