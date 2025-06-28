package tourism

import (
	"context"
	"time"

	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/dto"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/entity"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/service"
)

// CategoryAppService defines the application service for tourism category operations
type CategoryAppService struct {
	categoryService *service.CategoryService
}

// NewCategoryAppService creates a new CategoryAppService instance
func NewCategoryAppService(
	categoryService *service.CategoryService,
) *CategoryAppService {
	return &CategoryAppService{
		categoryService: categoryService,
	}
}

// CreateCategory creates a new category
func (s *CategoryAppService) CreateCategory(ctx context.Context, req *dto.CategoryCreateRequest) (*dto.CategoryResponse, error) {
	// Create category entity from request
	category := entity.NewCategory(
		req.Name,
		req.Description,
		req.ParentID,
		req.Icon,
		req.SortOrder,
	)

	// Create category using domain service
	if err := s.categoryService.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	// Return response
	return s.entityToResponse(ctx, category)
}

// GetCategory retrieves a category by ID
func (s *CategoryAppService) GetCategory(ctx context.Context, id string) (*dto.CategoryResponse, error) {
	// Get category using domain service
	category, err := s.categoryService.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Return response
	return s.entityToResponse(ctx, category)
}

// UpdateCategory updates a category
func (s *CategoryAppService) UpdateCategory(ctx context.Context, id string, req *dto.CategoryUpdateRequest) (*dto.CategoryResponse, error) {
	// Get existing category
	category, err := s.categoryService.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update category entity from request
	category.Update(
		req.Name,
		req.Description,
		req.ParentID,
		req.Icon,
		req.SortOrder,
	)

	// Update category using domain service
	if err := s.categoryService.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}

	// Return response
	return s.entityToResponse(ctx, category)
}

// DeleteCategory deletes a category
func (s *CategoryAppService) DeleteCategory(ctx context.Context, id string) error {
	return s.categoryService.DeleteCategory(ctx, id)
}

// ListCategories lists all categories with pagination
func (s *CategoryAppService) ListCategories(ctx context.Context, offset, limit int) (*dto.CategoryListResponse, error) {
	// List categories using domain service
	categories, total, err := s.categoryService.ListCategories(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to responses
	responses := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		response, err := s.entityToResponse(ctx, category)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}

	// Return response
	return &dto.CategoryListResponse{
		Total: total,
		Items: responses,
	}, nil
}

// GetCategoryTree retrieves the category tree structure
func (s *CategoryAppService) GetCategoryTree(ctx context.Context) ([]dto.CategoryTreeNode, error) {
	// Get root categories
	rootCategories, err := s.categoryService.ListRootCategories(ctx)
	if err != nil {
		return nil, err
	}

	// Build category tree
	tree := make([]dto.CategoryTreeNode, 0, len(rootCategories))
	for _, rootCategory := range rootCategories {
		node, err := s.buildCategoryTreeNode(ctx, rootCategory)
		if err != nil {
			return nil, err
		}
		tree = append(tree, *node)
	}

	return tree, nil
}

// Helper methods

// entityToResponse converts a category entity to response
func (s *CategoryAppService) entityToResponse(ctx context.Context, category *entity.Category) (*dto.CategoryResponse, error) {
	// Get parent name if parent ID is provided
	var parentName string
	if category.ParentID != "" {
		parentCategory, err := s.categoryService.GetCategoryByID(ctx, category.ParentID)
		if err == nil && parentCategory != nil {
			parentName = parentCategory.Name
		}
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ParentID:    category.ParentID,
		ParentName:  parentName,
		Icon:        category.Icon,
		SortOrder:   category.SortOrder,
		CreatedAt:   category.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   category.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// buildCategoryTreeNode recursively builds the category tree structure
func (s *CategoryAppService) buildCategoryTreeNode(ctx context.Context, category *entity.Category) (*dto.CategoryTreeNode, error) {
	// Create node
	node := &dto.CategoryTreeNode{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Icon:        category.Icon,
		SortOrder:   category.SortOrder,
		Children:    []dto.CategoryTreeNode{},
	}

	// Get children categories
	children, err := s.categoryService.ListCategoriesByParent(ctx, category.ID)
	if err != nil {
		return nil, err
	}

	// Recursively build children nodes
	for _, child := range children {
		childNode, err := s.buildCategoryTreeNode(ctx, child)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, *childNode)
	}

	return node, nil
}
