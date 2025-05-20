// filepath: /Users/amit/workspace/sutram/internal/modules/community/application/dto/dto.go
package dto

import (
	"time"

	"algo/internal/modules/community/domain/model"

	"github.com/google/uuid"
)

// ApiResponse is the standard API response format
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse wraps data with pagination information
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, total int, page int, perPage int) PaginatedResponse {
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}
	return PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
}

// TopicFilterRequest contains parameters for filtering topics
type TopicFilterRequest struct {
	Category string `form:"category"`
	Tag      string `form:"tag"`
	Sort     string `form:"sort"`
	Search   string `form:"search"`
	Page     int    `form:"page,default=1"`
	PerPage  int    `form:"per_page,default=20"`
}

// CreateTopicRequest is the request for creating a topic
type CreateTopicRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	CategoryID int    `json:"category_id" binding:"required"`
	TagIDs     []int  `json:"tag_ids"`
	IsPinned   bool   `json:"is_pinned"`
}

// UpdateTopicRequest is the request for updating a topic
type UpdateTopicRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID int    `json:"category_id"`
	TagIDs     []int  `json:"tag_ids"`
	IsPinned   *bool  `json:"is_pinned"`
	IsClosed   *bool  `json:"is_closed"`
}

// CreatePostRequest is the request for creating a post
type CreatePostRequest struct {
	Content      string `json:"content" binding:"required"`
	ParentPostID *int   `json:"parent_post_id"`
}

// UpdatePostRequest is the request for updating a post
type UpdatePostRequest struct {
	Content    string `json:"content"`
	IsSolution *bool  `json:"is_solution"`
}

// LikeRequest is the request for liking content
type LikeRequest struct {
	LikeableType string `json:"likeable_type" binding:"required,oneof=topic post"` // "topic" or "post"
	LikeableID   int    `json:"likeable_id" binding:"required"`
}

// TopicResponse represents a topic in API responses
type TopicResponse struct {
	ID             int           `json:"id"`
	Title          string        `json:"title"`
	Slug           string        `json:"slug"`
	Content        string        `json:"content"`
	CategoryID     int           `json:"category_id"`
	Category       string        `json:"category"`
	AuthorID       uuid.UUID     `json:"author_id"`
	AuthorName     string        `json:"author_name"`
	IsPinned       bool          `json:"is_pinned"`
	IsClosed       bool          `json:"is_closed"`
	ViewCount      int           `json:"view_count"`
	ReplyCount     int           `json:"reply_count"`
	LikeCount      int           `json:"like_count"`
	BookmarkCount  int           `json:"bookmark_count"`
	Tags           []TagResponse `json:"tags"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	LastActivityAt time.Time     `json:"last_activity_at"`
}

// ToTopicResponse converts a topic model to a topic response
func ToTopicResponse(topic *model.Topic) TopicResponse {
	tagResponses := make([]TagResponse, len(topic.Tags))
	for i, tag := range topic.Tags {
		tagResponses[i] = TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		}
	}

	return TopicResponse{
		ID:             topic.ID,
		Title:          topic.Title,
		Slug:           topic.Slug,
		Content:        topic.Content,
		CategoryID:     topic.CategoryID,
		Category:       topic.Category.Name,
		AuthorID:       topic.AuthorID,
		AuthorName:     topic.Author.Name,
		IsPinned:       topic.IsPinned,
		IsClosed:       topic.IsClosed,
		ViewCount:      topic.ViewCount,
		ReplyCount:     topic.Stats.ReplyCount,
		LikeCount:      topic.Stats.LikeCount,
		BookmarkCount:  topic.Stats.BookmarkCount,
		Tags:           tagResponses,
		CreatedAt:      topic.CreatedAt,
		UpdatedAt:      topic.UpdatedAt,
		LastActivityAt: topic.LastActivityAt,
	}
}

// TopicDetailResponse represents a detailed topic in API responses
type TopicDetailResponse struct {
	TopicResponse
	LastReplyAt       *time.Time `json:"last_reply_at"`
	LastReplyUserID   *uuid.UUID `json:"last_reply_user_id"`
	LastReplyUserName *string    `json:"last_reply_user_name"`
}

// ToTopicDetailResponse converts a topic model to a detailed topic response
func ToTopicDetailResponse(topic *model.Topic) TopicDetailResponse {
	base := ToTopicResponse(topic)

	var lastReplyAt *time.Time
	var lastReplyUserID *uuid.UUID
	var lastReplyUserName *string

	if topic.Stats.LastReplyAt != nil {
		lastReplyAt = topic.Stats.LastReplyAt
	}

	if topic.Stats.LastReplyUserID != nil {
		lastReplyUserID = topic.Stats.LastReplyUserID
		lastReplyUserName = &topic.Stats.LastReplyUser.Name
	}

	return TopicDetailResponse{
		TopicResponse:     base,
		LastReplyAt:       lastReplyAt,
		LastReplyUserID:   lastReplyUserID,
		LastReplyUserName: lastReplyUserName,
	}
}

// PostResponse represents a post in API responses
type PostResponse struct {
	ID           int        `json:"id"`
	TopicID      int        `json:"topic_id"`
	AuthorID     uuid.UUID  `json:"author_id"`
	AuthorName   string     `json:"author_name"`
	Content      string     `json:"content"`
	IsSolution   bool       `json:"is_solution"`
	ParentPostID *int       `json:"parent_post_id,omitempty"`
	LikeCount    int        `json:"like_count"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	EditCount    int        `json:"edit_count"`
	EditedByID   *uuid.UUID `json:"edited_by_id,omitempty"`
	EditedByName *string    `json:"edited_by_name,omitempty"`
}

// ToPostResponse converts a post model to a post response
func ToPostResponse(post *model.Post) PostResponse {
	var editedByID *uuid.UUID
	var editedByName *string

	if post.EditedBy != nil {
		editedByID = &post.EditedBy.ID
		editedByName = &post.EditedBy.Name
	}

	return PostResponse{
		ID:           post.ID,
		TopicID:      post.TopicID,
		AuthorID:     post.AuthorID,
		AuthorName:   post.Author.Name,
		Content:      post.Content,
		IsSolution:   post.IsSolution,
		ParentPostID: post.ParentPostID,
		LikeCount:    post.LikeCount,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		EditCount:    post.EditCount,
		EditedByID:   editedByID,
		EditedByName: editedByName,
	}
}

// CategoryResponse represents a category in API responses
type CategoryResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Icon        *string `json:"icon"`
	Color       *string `json:"color"`
	ParentID    *int    `json:"parent_id"`
	TopicCount  int     `json:"topic_count"`
}

// ToCategoryResponse converts a category model to a category response
func ToCategoryResponse(category *model.Category) CategoryResponse {
	return CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Icon:        category.Icon,
		Color:       category.Color,
		ParentID:    category.ParentID,
		TopicCount:  category.TopicCount,
	}
}

// TagResponse represents a tag in API responses
type TagResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// UserProfileResponse represents a user's community profile in API responses
type UserProfileResponse struct {
	ID                uuid.UUID       `json:"id"`
	Name              string          `json:"name"`
	ReputationScore   int             `json:"reputation_score"`
	ContributionCount int             `json:"contribution_count"`
	JoinedDate        time.Time       `json:"joined_date"`
	LastActive        time.Time       `json:"last_active"`
	TopicCount        int             `json:"topic_count"`
	PostCount         int             `json:"post_count"`
	BadgeCount        int             `json:"badge_count"`
	Badges            []BadgeResponse `json:"badges"`
}

// ToUserProfileResponse converts a user profile model to a user profile response
func ToUserProfileResponse(userProfile *model.UserProfile) UserProfileResponse {
	badgeResponses := make([]BadgeResponse, len(userProfile.Badges))
	for i, badge := range userProfile.Badges {
		badgeResponses[i] = ToBadgeResponse(badge)
	}

	return UserProfileResponse{
		ID:                userProfile.ID,
		Name:              userProfile.Name,
		ReputationScore:   userProfile.ReputationScore,
		ContributionCount: userProfile.ContributionCount,
		JoinedDate:        userProfile.JoinedDate,
		LastActive:        userProfile.LastActive,
		TopicCount:        userProfile.TopicCount,
		PostCount:         userProfile.PostCount,
		BadgeCount:        len(userProfile.Badges),
		Badges:            badgeResponses,
	}
}

// BadgeResponse represents a badge in API responses
type BadgeResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	BadgeType   string     `json:"badge_type"`
	TypeName    string     `json:"type_name,omitempty"`
	Points      int        `json:"points"`
	Icon        *string    `json:"icon,omitempty"`
	Color       *string    `json:"color,omitempty"`
	AwardedAt   *time.Time `json:"awarded_at,omitempty"`
}

// ToBadgeResponse converts a badge model to a badge response
func ToBadgeResponse(badge *model.Badge) BadgeResponse {
	return BadgeResponse{
		ID:          badge.ID,
		Name:        badge.Name,
		Slug:        badge.Slug,
		Description: badge.Description,
		BadgeType:   badge.BadgeType,
		TypeName:    badge.TypeName,
		Points:      badge.Points,
		Icon:        badge.Icon,
		Color:       badge.Color,
		AwardedAt:   badge.AwardedAt,
	}
}
