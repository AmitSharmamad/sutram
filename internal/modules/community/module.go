// filepath: /Users/amit/workspace/sutram/internal/modules/community/module.go
package community

import (
	"database/sql"

	infraRepo "algo/internal/infra/persistance/repository"
	"algo/internal/modules/community/application/service"
	"algo/internal/modules/community/infrastructure/persistence/postgres"
	"algo/internal/modules/community/interfaces/http"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Module represents the community module
type Module struct {
	db             *sql.DB
	controller     *http.CommunityController
	userController interface{}
}

// ProvideModule provides a wire provider set for the community module
var ProvideModule = wire.NewSet(
	ProvideCommunityModule,
)

// ProvideCommunityModule provides a community module instance
func ProvideCommunityModule() (*Module, error) {
	return New(), nil
}

// New creates a new community module
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

// initialize sets up the community module
func (m *Module) initialize() {
	// Create repositories
	topicRepo := postgres.NewTopicRepository(m.db)
	categoryRepo := postgres.NewCategoryRepository(m.db)
	postRepo := postgres.NewPostRepository(m.db)
	tagRepo := postgres.NewTagRepository(m.db)
	userReputationRepo := postgres.NewUserReputationRepository(m.db)
	likeRepo := postgres.NewLikeRepository(m.db)
	bookmarkRepo := postgres.NewBookmarkRepository(m.db)
	pollRepo := postgres.NewPollRepository(m.db)
	badgeRepo := postgres.NewBadgeRepository(m.db)

	// Create services
	topicService := service.NewTopicService(topicRepo, tagRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	postService := service.NewPostService(postRepo, topicRepo)
	userService := service.NewUserService(userReputationRepo, badgeRepo)
	likeService := service.NewLikeService(likeRepo, topicRepo, postRepo)
	bookmarkService := service.NewBookmarkService(bookmarkRepo, topicRepo)
	pollService := service.NewPollService(pollRepo, topicRepo)
	badgeService := service.NewBadgeService(badgeRepo, userReputationRepo)

	// Create controller
	m.controller = http.NewCommunityController(
		topicService,
		categoryService,
		postService,
		userService,
		likeService,
		bookmarkService,
		pollService,
		badgeService,
	)
}

// RegisterRoutes registers the community module routes with the given router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	// Cast userController to the expected type
	userController, ok := m.userController.(*userHttp.UserController)
	if !ok {
		// If userController is not set or is of wrong type, we can't proceed with auth
		// In a production app, you'd want to log this error properly
		panic("Community module: UserController not set or invalid type")
	}

	http.RegisterRoutes(router, m.controller, userController)
}

// GetController returns the community controller
func (m *Module) GetController() *http.CommunityController {
	return m.controller
}

// SetUserController sets the user controller for auth middleware
func (m *Module) SetUserController(userController interface{}) {
	m.userController = userController
}
