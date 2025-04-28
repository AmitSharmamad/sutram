// filepath: /Users/amit/workspace/sutram/internal/modules/community/interfaces/http/router.go
package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the community module routes
func RegisterRoutes(router *gin.RouterGroup, controller *CommunityController) {
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
	communityAuth.Use(AuthMiddleware())
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

// AuthMiddleware implements authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from header
		tokenString := ExtractToken(ctx)
		if tokenString == "" {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		// In a real implementation, validate the token and extract user ID
		// For now, just set a placeholder user ID
		userID, err := ExtractUserIDFromToken(tokenString)
		if err != nil {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}

// ExtractToken extracts the token from the Authorization header
func ExtractToken(ctx *gin.Context) string {
	bearerToken := ctx.GetHeader("Authorization")
	if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
		return bearerToken[7:]
	}
	return ""
}

// ExtractUserIDFromToken extracts the user ID from the token
// This is a placeholder - in a real app, this would validate the JWT and extract the user ID
func ExtractUserIDFromToken(tokenString string) (interface{}, error) {
	// Placeholder - in a real app, this would parse and validate a JWT
	// Mock user ID for development
	return MockUserID(), nil
}

// MockUserID returns a mock user ID for development
func MockUserID() interface{} {
	// Return a mock UUID for testing
	return "00000000-0000-0000-0000-000000000000"
}
