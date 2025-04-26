package shared

import (
	"time"

	"github.com/google/uuid"
)

// Community models

// UserReputation represents a user's reputation in the community
type UserReputation struct {
	UserID            uuid.UUID `json:"user_id"`
	ReputationScore   int       `json:"reputation_score"`
	ContributionCount int       `json:"contribution_count"`
	JoinedDate        time.Time `json:"joined_date"`
	LastActive        time.Time `json:"last_active"`
}

// Category represents a discussion category
type Category struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description,omitempty"`
	Icon         string    `json:"icon,omitempty"`
	Color        string    `json:"color,omitempty"`
	ParentID     *int      `json:"parent_id,omitempty"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TopicCount   int       `json:"topic_count,omitempty"` // Calculated field
}

// Topic represents a discussion topic
type Topic struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Content        string    `json:"content"`
	CategoryID     int       `json:"category_id"`
	AuthorID       uuid.UUID `json:"author_id"`
	Author         *User     `json:"author,omitempty"`
	IsPinned       bool      `json:"is_pinned"`
	IsClosed       bool      `json:"is_closed"`
	IsHidden       bool      `json:"is_hidden"`
	ViewCount      int       `json:"view_count"`
	ReplyCount     int       `json:"reply_count,omitempty"`    // From topic_stats
	LikeCount      int       `json:"like_count,omitempty"`     // From topic_stats
	BookmarkCount  int       `json:"bookmark_count,omitempty"` // From topic_stats
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastActivityAt time.Time `json:"last_activity_at"`
	Tags           []Tag     `json:"tags,omitempty"`
	Category       *Category `json:"category,omitempty"`
}

// TopicStats represents statistics for a topic
type TopicStats struct {
	TopicID         int        `json:"topic_id"`
	ReplyCount      int        `json:"reply_count"`
	LikeCount       int        `json:"like_count"`
	BookmarkCount   int        `json:"bookmark_count"`
	LastReplyID     *int       `json:"last_reply_id,omitempty"`
	LastReplyAt     *time.Time `json:"last_reply_at,omitempty"`
	LastReplyUserID *uuid.UUID `json:"last_reply_user_id,omitempty"`
}

// Post represents a reply to a topic
type Post struct {
	ID           int        `json:"id"`
	TopicID      int        `json:"topic_id"`
	AuthorID     uuid.UUID  `json:"author_id"`
	Author       *User      `json:"author,omitempty"`
	Content      string     `json:"content"`
	IsSolution   bool       `json:"is_solution"`
	ParentPostID *int       `json:"parent_post_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	EditedByID   *uuid.UUID `json:"edited_by_id,omitempty"`
	EditCount    int        `json:"edit_count"`
	LikeCount    int        `json:"like_count,omitempty"` // Calculated field
}

// Tag represents a topic tag
type Tag struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	Color       string    `json:"color,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// Like represents a like on a topic or post
type Like struct {
	ID           int       `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	LikeableType string    `json:"likeable_type"` // 'topic' or 'post'
	LikeableID   int       `json:"likeable_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// Bookmark represents a user bookmarking a topic
type Bookmark struct {
	ID        int       `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TopicID   int       `json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
}

// UserFollow represents a user following another user
type UserFollow struct {
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// CategoryFollow represents a user following a category
type CategoryFollow struct {
	UserID     uuid.UUID `json:"user_id"`
	CategoryID int       `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// TopicFollow represents a user following a topic
type TopicFollow struct {
	UserID    uuid.UUID `json:"user_id"`
	TopicID   int       `json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Notification represents a notification for a user
type Notification struct {
	ID               int        `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	NotificationType string     `json:"notification_type"` // 'mention', 'reply', 'like', etc.
	ActorID          *uuid.UUID `json:"actor_id,omitempty"`
	ReferenceType    string     `json:"reference_type"` // 'topic', 'post', 'badge', etc.
	ReferenceID      int        `json:"reference_id"`
	Message          string     `json:"message"`
	IsRead           bool       `json:"is_read"`
	CreatedAt        time.Time  `json:"created_at"`
	Actor            *User      `json:"actor,omitempty"` // The user who triggered the notification
}

// BadgeType represents a type of badge
type BadgeType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon,omitempty"`
	Color       string    `json:"color,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Badge represents a badge that can be earned by users
type Badge struct {
	ID          int        `json:"id"`
	TypeID      int        `json:"type_id"`
	Type        *BadgeType `json:"type,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Requirement string     `json:"requirement"`
	Tier        int        `json:"tier"`
	Points      int        `json:"points"`
	Icon        string     `json:"icon,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
}

// UserBadge represents a badge earned by a user
type UserBadge struct {
	ID        int       `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	BadgeID   int       `json:"badge_id"`
	Badge     *Badge    `json:"badge,omitempty"`
	AwardedAt time.Time `json:"awarded_at"`
}

// Report represents a content report
type Report struct {
	ID           int        `json:"id"`
	ReporterID   uuid.UUID  `json:"reporter_id"`
	Reporter     *User      `json:"reporter,omitempty"`
	ContentType  string     `json:"content_type"` // 'topic', 'post', 'user'
	ContentID    int        `json:"content_id"`
	Reason       string     `json:"reason"`
	Details      string     `json:"details,omitempty"`
	Status       string     `json:"status"` // 'pending', 'resolved', 'dismissed'
	ReviewedByID *uuid.UUID `json:"reviewed_by_id,omitempty"`
	ReviewedBy   *User      `json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ActivityLog represents a log of user activity
type ActivityLog struct {
	ID            int        `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	ActivityType  string     `json:"activity_type"`
	Description   string     `json:"description"`
	ReferenceType string     `json:"reference_type,omitempty"`
	ReferenceID   *uuid.UUID `json:"reference_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// Poll represents a community poll
type Poll struct {
	ID               int          `json:"id"`
	TopicID          int          `json:"topic_id"`
	Title            string       `json:"title"`
	Description      string       `json:"description,omitempty"`
	CreatedByID      uuid.UUID    `json:"created_by_id"`
	CreatedBy        *User        `json:"created_by,omitempty"`
	IsMultipleChoice bool         `json:"is_multiple_choice"`
	IsClosed         bool         `json:"is_closed"`
	ClosesAt         *time.Time   `json:"closes_at,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
	Options          []PollOption `json:"options,omitempty"`
}

// PollOption represents an option in a poll
type PollOption struct {
	ID           int    `json:"id"`
	PollID       int    `json:"poll_id"`
	OptionText   string `json:"option_text"`
	DisplayOrder int    `json:"display_order"`
	VoteCount    int    `json:"vote_count,omitempty"` // Calculated field
}

// PollVote represents a vote in a poll
type PollVote struct {
	ID        int       `json:"id"`
	PollID    int       `json:"poll_id"`
	OptionID  int       `json:"option_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// PaginatedTopics represents a paginated list of topics
type PaginatedTopics struct {
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
	Topics []Topic `json:"topics"`
}

// PaginatedPosts represents a paginated list of posts
type PaginatedPosts struct {
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Posts []Post `json:"posts"`
}
