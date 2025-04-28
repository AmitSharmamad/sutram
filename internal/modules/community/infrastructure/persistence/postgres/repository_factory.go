// filepath: /Users/amit/workspace/sutram/internal/modules/community/infrastructure/persistence/postgres/repository_factory.go
package postgres

import (
	"database/sql"

	"algo/internal/modules/community/domain/repository"
)

// NewTopicRepository creates a new topic repository
func NewTopicRepository(db *sql.DB) repository.TopicRepository {
	return &TopicRepository{db: db}
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *sql.DB) repository.PostRepository {
	return &PostRepository{db: db}
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &CategoryRepository{db: db}
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *sql.DB) repository.TagRepository {
	return &TagRepository{db: db}
}

// NewUserReputationRepository creates a new user reputation repository
func NewUserReputationRepository(db *sql.DB) repository.UserReputationRepository {
	return &UserReputationRepository{db: db}
}

// NewLikeRepository creates a new like repository
func NewLikeRepository(db *sql.DB) repository.LikeRepository {
	return &LikeRepository{db: db}
}

// NewBookmarkRepository creates a new bookmark repository
func NewBookmarkRepository(db *sql.DB) repository.BookmarkRepository {
	return &BookmarkRepository{db: db}
}

// NewPollRepository creates a new poll repository
func NewPollRepository(db *sql.DB) repository.PollRepository {
	return &PollRepository{db: db}
}

// NewBadgeRepository creates a new badge repository
func NewBadgeRepository(db *sql.DB) repository.BadgeRepository {
	return &BadgeRepository{db: db}
}
