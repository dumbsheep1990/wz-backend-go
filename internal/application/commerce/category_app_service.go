package commerce

import (
	"context"
	"time"

	"wz-backend-go/internal/domain/commerce/dto"
	"wz-backend-go/internal/domain/commerce/entity"
	"wz-backend-go/internal/domain/commerce/repository"
	"wz-backend-go/internal/domain/commerce/service"
)

// CategoryAppService handles category-related application use cases
type CategoryAppService struct {
	commerceService    *service.CommerceService
	categoryRepository repository.CategoryRepository
}

// NewCategoryAppService creates a new instance of CategoryAppService
func NewCategoryAppService(
	commerceService *service.CommerceService,
	categoryRepo repository.CategoryRepository,
) *CategoryAppService {
	return &CategoryAppService{
		commerceService:    commerceService,
		categoryRepository: categoryRepo,
	}
}

// CreateCategory creates a new product category
func (s *CategoryAppService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.commerceService.CreateCategory(
		ctx,
		req.Name,
		req.DisplayName,
		req.Description,
		req.ParentID,
		req.IconURL,
	)
	
	if err != nil {
		return nil, err
	}
	
	return s.categoryToResponse(ctx, category, false)
}

// UpdateCategory updates an existing category
func (s *CategoryAppService) UpdateCategory(ctx context.Context, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	// First check if category exists
	category, err := s.categoryRepository.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	
	if category == nil {
		return nil, service.ErrCategoryNotFound
	}
	
	// Update category details
	category, err = s.commerceService.UpdateCategory(
		ctx,
		req.ID,
		req.Name,
		req.DisplayName,
		req.Description,
		req.IconURL,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Handle optional boolean values
	if req.IsActive != nil {
		if *req.IsActive {
			category.Activate()
		} else {
			category.Deactivate()
		}
		
		// Save changes
		if err := s.categoryRepository.Update(ctx, category); err != nil {
			return nil, err
		}
	}
	
	return s.categoryToResponse(ctx, category, false)
}

// GetCategoryByID retrieves a category by its ID
func (s *CategoryAppService) GetCategoryByID(ctx context.Context, id string, includeChildren bool) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if category == nil {
		return nil, service.ErrCategoryNotFound
	}
	
	return s.categoryToResponse(ctx, category, includeChildren)
}

// GetRootCategories retrieves all root categories
func (s *CategoryAppService) GetRootCategories(ctx context.Context, activeOnly bool) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepository.FindRootCategories(ctx)
	if err != nil {
		return nil, err
	}
	
	// Filter active categories if requested
	if activeOnly {
		var activeCats []*entity.Category
		for _, cat := range categories {
			if cat.IsActive {
				activeCats = append(activeCats, cat)
			}
		}
		categories = activeCats
	}
	
	responses := make([]*dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		resp, err := s.categoryToResponse(ctx, category, true)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	
	return responses, nil
}

// GetCategories retrieves all categories with optional filtering
func (s *CategoryAppService) GetCategories(ctx context.Context, filter *dto.CategoryFilterRequest) ([]*dto.CategoryResponse, error) {
	var categories []*entity.Category
	var err error
	
	if filter.ParentID != "" {
		// Get categories by parent
		categories, err = s.categoryRepository.FindByParentID(ctx, filter.ParentID)
	} else if filter.Level > 0 {
		// Get categories by level
		categories, err = s.categoryRepository.FindByLevel(ctx, filter.Level)
	} else {
		// Get all categories
		categories, err = s.categoryRepository.FindAll(ctx, filter.ActiveOnly)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Filter active categories if requested
	if filter.ActiveOnly {
		var activeCats []*entity.Category
		for _, cat := range categories {
			if cat.IsActive {
				activeCats = append(activeCats, cat)
			}
		}
		categories = activeCats
	}
	
	responses := make([]*dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		resp, err := s.categoryToResponse(ctx, category, false)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	
	return responses, nil
}

// ReorderCategories updates the sort order of multiple categories
func (s *CategoryAppService) ReorderCategories(ctx context.Context, req *dto.ReorderCategoriesRequest) error {
	return s.commerceService.ReorderCategories(ctx, req.CategoryIDs)
}

// GetCategoryHierarchy retrieves the complete category hierarchy
func (s *CategoryAppService) GetCategoryHierarchy(ctx context.Context, activeOnly bool) ([]*dto.CategoryResponse, error) {
	// Get root categories
	rootCategories, err := s.categoryRepository.FindRootCategories(ctx)
	if err != nil {
		return nil, err
	}
	
	// Filter active categories if requested
	if activeOnly {
		var activeCats []*entity.Category
		for _, cat := range rootCategories {
			if cat.IsActive {
				activeCats = append(activeCats, cat)
			}
		}
		rootCategories = activeCats
	}
	
	// Build hierarchy
	responses := make([]*dto.CategoryResponse, 0, len(rootCategories))
	for _, rootCategory := range rootCategories {
		resp, err := s.buildCategoryHierarchy(ctx, rootCategory, activeOnly)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	
	return responses, nil
}

// buildCategoryHierarchy recursively builds the category hierarchy
func (s *CategoryAppService) buildCategoryHierarchy(ctx context.Context, category *entity.Category, activeOnly bool) (*dto.CategoryResponse, error) {
	// Convert category to response
	resp, err := s.categoryToResponse(ctx, category, false)
	if err != nil {
		return nil, err
	}
	
	// Get children
	children, err := s.categoryRepository.FindByParentID(ctx, category.ID)
	if err != nil {
		return nil, err
	}
	
	// Filter active children if requested
	if activeOnly {
		var activeChildren []*entity.Category
		for _, child := range children {
			if child.IsActive {
				activeChildren = append(activeChildren, child)
			}
		}
		children = activeChildren
	}
	
	// Recursively process children
	childResponses := make([]*dto.CategoryResponse, 0, len(children))
	for _, child := range children {
		childResp, err := s.buildCategoryHierarchy(ctx, child, activeOnly)
		if err != nil {
			return nil, err
		}
		childResponses = append(childResponses, childResp)
	}
	
	resp.Children = childResponses
	
	return resp, nil
}

// categoryToResponse converts a category entity to a category response DTO
func (s *CategoryAppService) categoryToResponse(ctx context.Context, category *entity.Category, includeChildren bool) (*dto.CategoryResponse, error) {
	if category == nil {
		return nil, nil
	}
	
	// Get product count
	productCount, err := s.categoryRepository.CountProducts(ctx, category.ID)
	if err != nil {
		productCount = 0 // default to 0 if there's an error
	}
	
	response := &dto.CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		DisplayName:  category.DisplayName,
		Description:  category.Description,
		ParentID:     category.ParentID,
		IconURL:      category.IconURL,
		SortOrder:    category.SortOrder,
		Level:        category.Level,
		IsActive:     category.IsActive,
		ProductCount: productCount,
		CreatedAt:    category.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    category.UpdatedAt.Format(time.RFC3339),
	}
	
	// Get children if requested
	if includeChildren {
		children, err := s.categoryRepository.FindByParentID(ctx, category.ID)
		if err != nil {
			return nil, err
		}
		
		if len(children) > 0 {
			childResponses := make([]*dto.CategoryResponse, 0, len(children))
			for _, child := range children {
				childResp, err := s.categoryToResponse(ctx, child, false) // Avoid deep recursion
				if err != nil {
					return nil, err
				}
				childResponses = append(childResponses, childResp)
			}
			response.Children = childResponses
		}
	}
	
	return response, nil
}
