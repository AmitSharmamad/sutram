// filepath: /Users/amit/workspace/sutram/internal/modules/community/interfaces/http/community_controller.go
package http

import (
	"errors"
	"net/http"
	"strconv"

	"algo/internal/modules/community/application/dto"
	"algo/internal/modules/community/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CommunityController handles community HTTP requests
type CommunityController struct {
	topicService    *service.TopicService
	categoryService *service.CategoryService
	postService     *service.PostService
	userService     *service.UserService
	likeService     *service.LikeService
	bookmarkService *service.BookmarkService
	pollService     *service.PollService
	badgeService    *service.BadgeService
}

// NewCommunityController creates a new community controller
func NewCommunityController(
	topicService *service.TopicService,
	categoryService *service.CategoryService,
	postService *service.PostService,
	userService *service.UserService,
	likeService *service.LikeService,
	bookmarkService *service.BookmarkService,
	pollService *service.PollService,
	badgeService *service.BadgeService,
) *CommunityController {
	return &CommunityController{
		topicService:    topicService,
		categoryService: categoryService,
		postService:     postService,
		userService:     userService,
		likeService:     likeService,
		bookmarkService: bookmarkService,
		pollService:     pollService,
		badgeService:    badgeService,
	}
}

// GetTopics handles the request to get all community topics with optional filtering
func (c *CommunityController) GetTopics(ctx *gin.Context) {
	var filters dto.TopicFilterRequest
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid filter parameters: " + err.Error(),
		})
		return
	}

	// Get topics with filters
	topics, total, err := c.topicService.GetTopics(
		ctx.Request.Context(),
		filters.Category,
		filters.Tag,
		filters.Sort,
		filters.Search,
		filters.Page,
		filters.PerPage,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get topics: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	topicResponses := make([]dto.TopicResponse, len(topics))
	for i, topic := range topics {
		topicResponses[i] = dto.ToTopicResponse(topic)
	}

	// Return paginated response
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data: dto.NewPaginatedResponse(
			topicResponses,
			total,
			filters.Page,
			filters.PerPage,
		),
	})
}

// GetTopicByID handles the request to get a specific topic
func (c *CommunityController) GetTopicByID(ctx *gin.Context) {
	idStr := ctx.Param("topic_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid topic ID",
		})
		return
	}

	// Get topic
	topic, err := c.topicService.GetTopicByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTopicNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Topic not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get topic: " + err.Error(),
		})
		return
	}

	// Increment view count
	go c.topicService.IncrementViews(ctx.Request.Context(), id)

	// Convert to response DTO
	response := dto.ToTopicDetailResponse(topic)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    response,
	})
}

// CreateTopic handles the request to create a new topic
func (c *CommunityController) CreateTopic(ctx *gin.Context) {
	var req dto.CreateTopicRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Parse tag IDs if provided
	var tagIDs []int
	if req.TagIDs != nil && len(req.TagIDs) > 0 {
		tagIDs = req.TagIDs
	}

	// Create topic
	topic, err := c.topicService.CreateTopic(
		ctx.Request.Context(),
		userID,
		req.Title,
		req.Content,
		req.CategoryID,
		tagIDs,
		req.IsPinned,
	)

	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
				Success: false,
				Error:   "Invalid category",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to create topic: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToTopicResponse(topic)

	ctx.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Topic created successfully",
		Data:    response,
	})
}

// UpdateTopic handles the request to update an existing topic
func (c *CommunityController) UpdateTopic(ctx *gin.Context) {
	idStr := ctx.Param("topic_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid topic ID",
		})
		return
	}

	var req dto.UpdateTopicRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Prepare updates
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.IsPinned != nil {
		updates["is_pinned"] = *req.IsPinned
	}
	if req.IsClosed != nil {
		updates["is_closed"] = *req.IsClosed
	}

	// Update topic
	updatedTopic, err := c.topicService.UpdateTopic(
		ctx.Request.Context(),
		id,
		userID,
		updates,
	)

	if err != nil {
		if errors.Is(err, service.ErrTopicNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Topic not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to update this topic",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to update topic: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToTopicResponse(updatedTopic)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Topic updated successfully",
		Data:    response,
	})
}

// DeleteTopic handles the request to delete a topic
func (c *CommunityController) DeleteTopic(ctx *gin.Context) {
	idStr := ctx.Param("topic_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid topic ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Delete topic
	err = c.topicService.DeleteTopic(ctx.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, service.ErrTopicNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Topic not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to delete this topic",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to delete topic: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Topic deleted successfully",
	})
}

// GetPosts handles the request to get posts for a topic
func (c *CommunityController) GetPosts(ctx *gin.Context) {
	topicIDStr := ctx.Param("topic_id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid topic ID",
		})
		return
	}

	// Get page parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "20"))

	// Get posts
	posts, total, err := c.postService.GetPostsByTopicID(ctx.Request.Context(), topicID, page, perPage)
	if err != nil {
		if errors.Is(err, service.ErrTopicNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Topic not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get posts: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	postResponses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = dto.ToPostResponse(post)
	}

	// Return paginated response
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data: dto.NewPaginatedResponse(
			postResponses,
			total,
			page,
			perPage,
		),
	})
}

// CreatePost handles the request to create a new post
func (c *CommunityController) CreatePost(ctx *gin.Context) {
	topicIDStr := ctx.Param("topic_id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid topic ID",
		})
		return
	}

	var req dto.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Create post
	post, err := c.postService.CreatePost(
		ctx.Request.Context(),
		topicID,
		userID,
		req.Content,
		req.ParentPostID,
	)

	if err != nil {
		if errors.Is(err, service.ErrTopicNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Topic not found",
			})
			return
		}
		if errors.Is(err, service.ErrTopicClosed) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "Topic is closed for new posts",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to create post: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToPostResponse(post)

	ctx.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Post created successfully",
		Data:    response,
	})
}

// UpdatePost handles the request to update a post
func (c *CommunityController) UpdatePost(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid post ID",
		})
		return
	}

	var req dto.UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Prepare updates
	updates := make(map[string]interface{})
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.IsSolution != nil {
		updates["is_solution"] = *req.IsSolution
	}

	// Update post
	updatedPost, err := c.postService.UpdatePost(
		ctx.Request.Context(),
		postID,
		userID,
		updates,
	)

	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Post not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to update this post",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to update post: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToPostResponse(updatedPost)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Post updated successfully",
		Data:    response,
	})
}

// DeletePost handles the request to delete a post
func (c *CommunityController) DeletePost(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid post ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Delete post
	err = c.postService.DeletePost(ctx.Request.Context(), postID, userID)
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Post not found",
			})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, dto.ApiResponse{
				Success: false,
				Error:   "You are not authorized to delete this post",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to delete post: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Post deleted successfully",
	})
}

// GetCategories handles the request to get all community categories
func (c *CommunityController) GetCategories(ctx *gin.Context) {
	categories, err := c.categoryService.GetAllCategories(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get categories: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	categoryResponses := make([]dto.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = dto.ToCategoryResponse(category)
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    categoryResponses,
	})
}

// LikeContent handles the request to like a topic or post
func (c *CommunityController) LikeContent(ctx *gin.Context) {
	var req dto.LikeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Create like
	err := c.likeService.Like(
		ctx.Request.Context(),
		userID,
		req.LikeableType,
		req.LikeableID,
	)

	if err != nil {
		if errors.Is(err, service.ErrContentNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Content not found",
			})
			return
		}
		if errors.Is(err, service.ErrAlreadyLiked) {
			ctx.JSON(http.StatusConflict, dto.ApiResponse{
				Success: false,
				Error:   "Already liked",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to like content: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Content liked successfully",
	})
}

// UnlikeContent handles the request to unlike a topic or post
func (c *CommunityController) UnlikeContent(ctx *gin.Context) {
	var req dto.LikeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Unlike content
	err := c.likeService.Unlike(
		ctx.Request.Context(),
		userID,
		req.LikeableType,
		req.LikeableID,
	)

	if err != nil {
		if errors.Is(err, service.ErrContentNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Content not found",
			})
			return
		}
		if errors.Is(err, service.ErrNotLiked) {
			ctx.JSON(http.StatusConflict, dto.ApiResponse{
				Success: false,
				Error:   "Content not liked",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to unlike content: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Content unliked successfully",
	})
}

// BookmarkTopic handles the request to bookmark a topic
func (c *CommunityController) BookmarkTopic(ctx *gin.Context) {
	topicIDStr := ctx.Param("topic_id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid topic ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Add bookmark
	err = c.bookmarkService.AddBookmark(ctx.Request.Context(), userID, topicID)
	if err != nil {
		if errors.Is(err, service.ErrTopicNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Topic not found",
			})
			return
		}
		if errors.Is(err, service.ErrAlreadyBookmarked) {
			ctx.JSON(http.StatusConflict, dto.ApiResponse{
				Success: false,
				Error:   "Topic already bookmarked",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to bookmark topic: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Topic bookmarked successfully",
	})
}

// UnbookmarkTopic handles the request to remove a bookmark
func (c *CommunityController) UnbookmarkTopic(ctx *gin.Context) {
	topicIDStr := ctx.Param("topic_id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid topic ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Remove bookmark
	err = c.bookmarkService.RemoveBookmark(ctx.Request.Context(), userID, topicID)
	if err != nil {
		if errors.Is(err, service.ErrTopicNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "Topic not found",
			})
			return
		}
		if errors.Is(err, service.ErrNotBookmarked) {
			ctx.JSON(http.StatusConflict, dto.ApiResponse{
				Success: false,
				Error:   "Topic not bookmarked",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to unbookmark topic: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Topic unbookmarked successfully",
	})
}

// GetUserBookmarks handles the request to get all bookmarks for a user
func (c *CommunityController) GetUserBookmarks(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	// Get page parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "20"))

	// Get bookmarks
	topics, total, err := c.bookmarkService.GetUserBookmarks(ctx.Request.Context(), userID, page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get bookmarks: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	topicResponses := make([]dto.TopicResponse, len(topics))
	for i, topic := range topics {
		topicResponses[i] = dto.ToTopicResponse(topic)
	}

	// Return paginated response
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data: dto.NewPaginatedResponse(
			topicResponses,
			total,
			page,
			perPage,
		),
	})
}

// GetUserProfile handles the request to get a user's community profile
func (c *CommunityController) GetUserProfile(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get user profile
	userProfile, err := c.userService.GetUserProfile(ctx.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "User not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get user profile: " + err.Error(),
		})
		return
	}

	// Convert to response DTO
	response := dto.ToUserProfileResponse(userProfile)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    response,
	})
}

// GetUserTopics handles the request to get all topics created by a user
func (c *CommunityController) GetUserTopics(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get page parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "20"))

	// Get topics
	topics, total, err := c.topicService.GetUserTopics(ctx.Request.Context(), userID, page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get user topics: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	topicResponses := make([]dto.TopicResponse, len(topics))
	for i, topic := range topics {
		topicResponses[i] = dto.ToTopicResponse(topic)
	}

	// Return paginated response
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data: dto.NewPaginatedResponse(
			topicResponses,
			total,
			page,
			perPage,
		),
	})
}

// GetUserBadges handles the request to get all badges for a user
func (c *CommunityController) GetUserBadges(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get badges
	badges, err := c.badgeService.GetUserBadges(ctx.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ApiResponse{
				Success: false,
				Error:   "User not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get user badges: " + err.Error(),
		})
		return
	}

	// Convert to response DTOs
	badgeResponses := make([]dto.BadgeResponse, len(badges))
	for i, badge := range badges {
		badgeResponses[i] = dto.ToBadgeResponse(badge)
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    badgeResponses,
	})
}

// Helper function to extract user ID from context
func GetUserIDFromContext(ctx *gin.Context) (uuid.UUID, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return uuid.UUID{}, false
	}
	return userID.(uuid.UUID), true
}
