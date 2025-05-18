// filepath: /Users/amit/workspace/sutram/internal/modules/community/application/service/category_service.go
package service

import (
	"context"
	"fmt"

	"algo/internal/modules/community/domain/model"
	"algo/internal/modules/community/domain/repository"
)

// CategoryService handles community category business logic
type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// GetAllCategories gets all categories
func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Get topic count for each category
	for _, category := range categories {
		count, err := s.categoryRepo.GetTopicCount(ctx, category.ID)
		if err == nil {
			category.TopicCount = count
		}
	}

	return categories, nil
}

// GetCategoryByID gets a category by ID
func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if category != nil {
		// Get topic count
		count, err := s.categoryRepo.GetTopicCount(ctx, category.ID)
		if err == nil {
			category.TopicCount = count
		}
	}

	return category, nil
}

// GetCategoryBySlug gets a category by slug
func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if category != nil {
		// Get topic count
		count, err := s.categoryRepo.GetTopicCount(ctx, category.ID)
		if err == nil {
			category.TopicCount = count
		}
	}

	return category, nil
}
