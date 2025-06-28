package service

import (
	"context"
	
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/entity"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/repository"
)

// CategoryService defines the domain service for tourism category operations
type CategoryService struct {
	categoryRepo   repository.CategoryRepository
	scenicSpotRepo repository.ScenicSpotRepository
}

// NewCategoryService creates a new CategoryService instance
func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	scenicSpotRepo repository.ScenicSpotRepository,
) *CategoryService {
	return &CategoryService{
		categoryRepo:   categoryRepo,
		scenicSpotRepo: scenicSpotRepo,
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(ctx context.Context, category *entity.Category) error {
	// Validate parent category if specified
	if category.ParentID != "" {
		parentCategory, err := s.categoryRepo.GetByID(ctx, category.ParentID)
		if err != nil {
			return err
		}
		if parentCategory == nil {
			return ErrParentCategoryNotFound
		}
	}

	return s.categoryRepo.Create(ctx, category)
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(ctx context.Context, id string) (*entity.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(ctx context.Context, category *entity.Category) error {
	// Validate parent category if specified
	if category.ParentID != "" {
		if category.ID == category.ParentID {
			return ErrCategoryCannotBeItsOwnParent
		}
		
		parentCategory, err := s.categoryRepo.GetByID(ctx, category.ParentID)
		if err != nil {
			return err
		}
		if parentCategory == nil {
			return ErrParentCategoryNotFound
		}
	}

	return s.categoryRepo.Update(ctx, category)
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	// Check if there are any subcategories
	subcategories, err := s.categoryRepo.ListByParent(ctx, id)
	if err != nil {
		return err
	}
	if len(subcategories) > 0 {
		return ErrCategoryHasSubcategories
	}

	// Check if there are any scenic spots in this category
	spots, _, err := s.scenicSpotRepo.ListByCategory(ctx, id, 0, 1)
	if err != nil {
		return err
	}
	if len(spots) > 0 {
		return ErrCategoryHasScenicSpots
	}

	return s.categoryRepo.Delete(ctx, id)
}

// ListCategories lists all categories
func (s *CategoryService) ListCategories(ctx context.Context, offset, limit int) ([]*entity.Category, int, error) {
	return s.categoryRepo.List(ctx, offset, limit)
}

// ListCategoriesByParent lists categories by parent ID
func (s *CategoryService) ListCategoriesByParent(ctx context.Context, parentID string) ([]*entity.Category, error) {
	return s.categoryRepo.ListByParent(ctx, parentID)
}

// ListRootCategories lists top-level categories
func (s *CategoryService) ListRootCategories(ctx context.Context) ([]*entity.Category, error) {
	return s.categoryRepo.ListRootCategories(ctx)
}
