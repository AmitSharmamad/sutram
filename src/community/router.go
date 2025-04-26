package community

import (
	"github.com/gin-gonic/gin"

	"algo/src/web"
)

// Router defines routes for community features
type Router struct {
	controller *Controller
}

// NewRouter creates a new instance of Router
func NewRouter(controller *Controller) *Router {
	return &Router{
		controller: controller,
	}
}

// RegisterRoutes registers community routes
func (r *Router) RegisterRoutes(router *gin.Engine) {
	communityGroup := router.Group("/api/v1/community")

	// Public routes (no authentication required)
	publicGroup := communityGroup.Group("")
	{
		// Categories
		publicGroup.GET("/categories", r.controller.GetCategories)
		publicGroup.GET("/categories/:id", r.controller.GetCategoryByID)
		publicGroup.GET("/categories/slug/:slug", r.controller.GetCategoryBySlug)

		// Topics
		publicGroup.GET("/topics", r.controller.GetTopics)
		publicGroup.GET("/topics/:id", r.controller.GetTopicByID)
		publicGroup.GET("/topics/slug/:slug", r.controller.GetTopicBySlug)

		// Posts
		publicGroup.GET("/topics/:topicId/posts", r.controller.GetPosts)
		publicGroup.GET("/posts/:id", r.controller.GetPostByID)

		// Tags
		publicGroup.GET("/tags", r.controller.GetTags)

		// Likes
		publicGroup.GET("/likes", r.controller.GetLikesByEntity)
	}

	// Protected routes (authentication required)
	protectedGroup := communityGroup.Group("")
	protectedGroup.Use(web.AuthMiddleware())
	{
		// Categories (admin only)
		protectedGroup.POST("/categories", web.RoleRequiredMiddleware("admin"), r.controller.CreateCategory)
		protectedGroup.PUT("/categories/:id", web.RoleRequiredMiddleware("admin"), r.controller.UpdateCategory)
		protectedGroup.DELETE("/categories/:id", web.RoleRequiredMiddleware("admin"), r.controller.DeleteCategory)

		// Topics
		protectedGroup.POST("/topics", r.controller.CreateTopic)
		protectedGroup.PUT("/topics/:id", r.controller.UpdateTopic)
		protectedGroup.DELETE("/topics/:id", r.controller.DeleteTopic)

		// Additional routes for topics
		// These routes should be uncommented once the controller methods are implemented
		// protectedGroup.POST("/topics/:id/pin", r.controller.PinTopic)
		// protectedGroup.POST("/topics/:id/close", r.controller.CloseTopic)

		// Posts
		// These routes should be uncommented once the controller methods are implemented
		// protectedGroup.POST("/posts", r.controller.CreatePost)
		// protectedGroup.PUT("/posts/:id", r.controller.UpdatePost)
		// protectedGroup.DELETE("/posts/:id", r.controller.DeletePost)
		// protectedGroup.POST("/posts/:id/solution", r.controller.MarkPostAsSolution)

		// Likes
		// This route should be uncommented once the controller method is implemented
		// protectedGroup.POST("/likes", r.controller.ToggleLike)

		/*
			// The following routes are commented out because the controller methods are not yet implemented
			// Uncomment and implement each feature as needed:

			// Bookmarks
			// protectedGroup.GET("/bookmarks", r.controller.GetBookmarks)
			// protectedGroup.POST("/bookmarks", r.controller.ToggleBookmark)

			// User Follows
			// protectedGroup.GET("/users/:userId/followers", r.controller.GetFollowers)
			// protectedGroup.GET("/users/:userId/following", r.controller.GetFollowing)
			// protectedGroup.POST("/users/:userId/follow", r.controller.ToggleUserFollow)

			// Category Follows
			// protectedGroup.GET("/users/me/category-follows", r.controller.GetCategoryFollows)
			// protectedGroup.POST("/categories/:id/follow", r.controller.ToggleCategoryFollow)

			// Topic Follows
			// protectedGroup.GET("/users/me/topic-follows", r.controller.GetTopicFollows)
			// protectedGroup.POST("/topics/:id/follow", r.controller.ToggleTopicFollow)

			// Notifications
			// protectedGroup.GET("/notifications", r.controller.GetNotifications)
			// protectedGroup.POST("/notifications/:id/read", r.controller.MarkNotificationAsRead)
			// protectedGroup.POST("/notifications/read-all", r.controller.MarkAllNotificationsAsRead)
			// protectedGroup.DELETE("/notifications/:id", r.controller.DeleteNotification)

			// Badges
			// protectedGroup.GET("/badges", r.controller.GetBadges)
			// protectedGroup.GET("/users/me/badges", r.controller.GetUserBadges)

			// Reports
			// protectedGroup.POST("/reports", r.controller.CreateReport)
			// protectedGroup.PUT("/reports/:id", web.RoleRequiredMiddleware("moderator"), r.controller.RespondToReport)

			// Polls
			// protectedGroup.POST("/polls", r.controller.CreatePoll)
			// protectedGroup.GET("/topics/:topicId/polls", r.controller.GetPollsByTopic)
			// protectedGroup.POST("/polls/:id/vote", r.controller.VoteInPoll)
			// protectedGroup.GET("/polls/:id/results", r.controller.GetPollResults)

			// User Reputation
			// protectedGroup.GET("/users/:userId/reputation", r.controller.GetUserReputation)
			// protectedGroup.GET("/leaderboard", r.controller.GetTopUsers)
		*/
	}
}

// Additional functions and helpers would be implemented here as needed
