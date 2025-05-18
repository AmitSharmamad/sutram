// filepath: /Users/amit/workspace/sutram/internal/modules/community/domain/model/topic.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// Topic represents a community discussion topic
type Topic struct {
	ID             int       `json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	Slug           string    `json:"slug" db:"slug"`
	Content        string    `json:"content" db:"content"`
	CategoryID     int       `json:"category_id" db:"category_id"`
	AuthorID       uuid.UUID `json:"author_id" db:"author_id"`
	IsPinned       bool      `json:"is_pinned" db:"is_pinned"`
	IsClosed       bool      `json:"is_closed" db:"is_closed"`
	IsHidden       bool      `json:"is_hidden" db:"is_hidden"`
	ViewCount      int       `json:"view_count" db:"view_count"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	LastActivityAt time.Time `json:"last_activity_at" db:"last_activity_at"`

	// Associated objects
	Category Category    `json:"category"` // Loaded separately
	Author   User        `json:"author"`   // Loaded separately
	Tags     []*Tag      `json:"tags"`     // Loaded separately
	Stats    *TopicStats `json:"stats"`    // Loaded separately
}

// TopicStats represents the statistics of a topic
type TopicStats struct {
	TopicID         int        `json:"topic_id" db:"topic_id"`
	ReplyCount      int        `json:"reply_count" db:"reply_count"`
	LikeCount       int        `json:"like_count" db:"like_count"`
	BookmarkCount   int        `json:"bookmark_count" db:"bookmark_count"`
	LastReplyID     *int       `json:"last_reply_id" db:"last_reply_id"`
	LastReplyAt     *time.Time `json:"last_reply_at" db:"last_reply_at"`
	LastReplyUserID *uuid.UUID `json:"last_reply_user_id" db:"last_reply_user_id"`

	// Associated objects
	LastReplyUser User `json:"last_reply_user"` // Loaded separately
}
