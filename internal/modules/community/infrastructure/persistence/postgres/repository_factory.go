// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/repository_factory.go
package postgres

import (
	"database/sql"
)

// NewTopicRepository creates a new topic repository
func NewTopicRepository(db *sql.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

// NewUserReputationRepository creates a new user reputation repository
func NewUserReputationRepository(db *sql.DB) *UserReputationRepository {
	return &UserReputationRepository{db: db}
}

// NewLikeRepository creates a new like repository
func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

// NewBookmarkRepository creates a new bookmark repository
func NewBookmarkRepository(db *sql.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

// NewPollRepository creates a new poll repository
func NewPollRepository(db *sql.DB) *PollRepository {
	return &PollRepository{db: db}
}

// NewBadgeRepository creates a new badge repository
func NewBadgeRepository(db *sql.DB) *BadgeRepository {
	return &BadgeRepository{db: db}
}
