package http

import (
	"algo/internal/infra/http/middleware"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the community module routes
func RegisterRoutes(router *gin.RouterGroup, controller *CommunityController, userController *userHttp.UserController) {
	// Create auth service adapter
	authServiceAdapter := userHttp.NewAuthServiceAdapter(userController.GetAuthService())

	// Public routes (no authentication required)
	router.GET("/community/topics", controller.GetTopics)
	router.GET("/community/topics/:topic_id", controller.GetTopicByID)
	router.GET("/community/categories", controller.GetCategories)
	router.GET("/community/topics/:topic_id/posts", controller.GetPosts)
	router.GET("/community/users/:user_id/profile", controller.GetUserProfile)
	router.GET("/community/users/:user_id/topics", controller.GetUserTopics)
	router.GET("/community/users/:user_id/badges", controller.GetUserBadges)

	// Authenticated routes
	communityAuth := router.Group("/community")
	communityAuth.Use(middleware.AuthMiddleware(authServiceAdapter))
	{
		// Topic management
		communityAuth.POST("/topics", controller.CreateTopic)
		communityAuth.PUT("/topics/:topic_id", controller.UpdateTopic)
		communityAuth.DELETE("/topics/:topic_id", controller.DeleteTopic)

		// Post management
		communityAuth.POST("/topics/:topic_id/posts", controller.CreatePost)
		communityAuth.PUT("/posts/:post_id", controller.UpdatePost)
		communityAuth.DELETE("/posts/:post_id", controller.DeletePost)

		// Like/Unlike functionality
		communityAuth.POST("/like", controller.LikeContent)
		communityAuth.DELETE("/unlike", controller.UnlikeContent)

		// Bookmark functionality
		communityAuth.POST("/topics/:topic_id/bookmark", controller.BookmarkTopic)
		communityAuth.DELETE("/topics/:topic_id/bookmark", controller.UnbookmarkTopic)
		communityAuth.GET("/bookmarks", controller.GetUserBookmarks)
	}
}
