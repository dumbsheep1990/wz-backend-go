package learn

import (
	"context"

	"wz-backend-go/internal/domain/learn/dto"
	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/service"
)

// CategoryAppService 分类应用服务，处理分类相关的应用场景
type CategoryAppService struct {
	categoryService *service.CategoryService
}

// NewCategoryAppService 创建分类应用服务
func NewCategoryAppService(
	categoryService *service.CategoryService,
) *CategoryAppService {
	return &CategoryAppService{
		categoryService: categoryService,
	}
}

// CreateCategory 创建分类
func (s *CategoryAppService) CreateCategory(ctx context.Context, name, description, icon string, parentID *string, order int) (*dto.CategoryDTO, error) {
	// 创建分类
	category, err := s.categoryService.CreateCategory(ctx, name, description, icon, parentID, order)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(category), nil
}

// UpdateCategory 更新分类
func (s *CategoryAppService) UpdateCategory(ctx context.Context, id, name, description, icon string, order int) (*dto.CategoryDTO, error) {
	// 更新分类
	category, err := s.categoryService.UpdateCategory(ctx, id, name, description, icon, order)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(category), nil
}

// ActivateCategory 激活分类
func (s *CategoryAppService) ActivateCategory(ctx context.Context, id string) (*dto.CategoryDTO, error) {
	// 激活分类
	category, err := s.categoryService.ActivateCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(category), nil
}

// DeactivateCategory 停用分类
func (s *CategoryAppService) DeactivateCategory(ctx context.Context, id string) (*dto.CategoryDTO, error) {
	// 停用分类
	category, err := s.categoryService.DeactivateCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(category), nil
}

// DeleteCategory 删除分类
func (s *CategoryAppService) DeleteCategory(ctx context.Context, id string) error {
	return s.categoryService.DeleteCategory(ctx, id)
}

// GetCategoryByID 获取分类详情
func (s *CategoryAppService) GetCategoryByID(ctx context.Context, id string) (*dto.CategoryDTO, error) {
	// 获取分类
	category, err := s.categoryService.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(category), nil
}

// GetCategoryTree 获取分类树
func (s *CategoryAppService) GetCategoryTree(ctx context.Context) ([]*dto.CategoryTreeDTO, error) {
	// 获取分类树
	categories, err := s.categoryService.GetCategoryTree(ctx)
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	treeMap := make(map[string]*dto.CategoryTreeDTO)
	topLevelNodes := make([]*dto.CategoryTreeDTO, 0)

	// 首先创建所有节点
	for _, cat := range categories {
		treeDTO := &dto.CategoryTreeDTO{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			Icon:        cat.Icon,
			Level:       cat.Level,
			Order:       cat.Order,
			IsActive:    cat.IsActive,
			Children:    make([]*dto.CategoryTreeDTO, 0),
		}
		treeMap[cat.ID] = treeDTO

		// 如果是顶级分类，直接添加到结果列表
		if cat.ParentID == nil {
			topLevelNodes = append(topLevelNodes, treeDTO)
		}
	}

	// 然后构建父子关系
	for _, cat := range categories {
		if cat.ParentID != nil {
			if parent, exists := treeMap[*cat.ParentID]; exists {
				parent.Children = append(parent.Children, treeMap[cat.ID])
			}
		}
	}

	return topLevelNodes, nil
}

// GetActiveCategories 获取所有激活状态的分类
func (s *CategoryAppService) GetActiveCategories(ctx context.Context) ([]*dto.CategoryDTO, error) {
	// 获取所有激活状态的分类
	categories, err := s.categoryService.GetActiveCategories(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO 列表
	dtos := make([]*dto.CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, s.convertToDTO(category))
	}

	return dtos, nil
}

// GetCategoriesWithCourseCount 获取带课程数的分类列表
func (s *CategoryAppService) GetCategoriesWithCourseCount(ctx context.Context) ([]*dto.CategoryDTO, error) {
	// 获取带课程数的分类
	categories, err := s.categoryService.GetCategoriesWithCourseCount(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO 列表
	dtos := make([]*dto.CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, s.convertToDTO(category))
	}

	return dtos, nil
}

// GetCategoryStats 获取分类统计信息
func (s *CategoryAppService) GetCategoryStats(ctx context.Context) (*dto.CategoryStats, error) {
	// 获取分类统计
	totalCount, level1Count, level2Count, err := s.categoryService.GetCategoryStats(ctx)
	if err != nil {
		return nil, err
	}

	// 构建统计 DTO
	stats := &dto.CategoryStats{
		TotalCount: int(totalCount),
		Level1Count: int(level1Count),
		Level2Count: int(level2Count),
	}

	return stats, nil
}

// 辅助函数：转换分类实体到 DTO
func (s *CategoryAppService) convertToDTO(category *entity.Category) *dto.CategoryDTO {
	return &dto.CategoryDTO{
		ID:           category.ID,
		Name:         category.Name,
		Description:  category.Description,
		Icon:         category.Icon,
		Level:        category.Level,
		Order:        category.Order,
		ParentID:     category.ParentID,
		IsActive:     category.IsActive,
		CoursesCount: category.CoursesCount,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}
}
