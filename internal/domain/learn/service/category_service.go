package service

import (
	"context"
	"errors"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
)

var (
	ErrCategoryNotFound = errors.New("分类不存在")
	ErrCategoryExists   = errors.New("分类已存在")
	ErrInvalidLevel     = errors.New("无效的分类级别")
)

// CategoryService 分类领域服务
type CategoryService struct {
	categoryRepo repository.CategoryRepository
	courseRepo   repository.CourseRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	courseRepo repository.CourseRepository,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		courseRepo:   courseRepo,
	}
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(ctx context.Context, name, description, icon string, parentID *string, order int) (*entity.Category, error) {
	// 确定分类级别
	level := 1
	if parentID != nil {
		// 获取父分类
		parent, err := s.categoryRepo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, ErrCategoryNotFound
		}
		level = parent.Level + 1
	}

	// 创建分类
	category := entity.NewCategory(name, parentID, level, order)
	category.Description = description
	category.Icon = icon

	// 保存分类
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(ctx context.Context, id, name, description, icon string, order int) (*entity.Category, error) {
	// 获取分类
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	// 更新分类信息
	category.Update(name, description, icon, order)

	// 保存分类
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// ActivateCategory 激活分类
func (s *CategoryService) ActivateCategory(ctx context.Context, id string) (*entity.Category, error) {
	// 获取分类
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	// 激活分类
	category.Activate()

	// 保存分类
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeactivateCategory 停用分类
func (s *CategoryService) DeactivateCategory(ctx context.Context, id string) (*entity.Category, error) {
	// 获取分类
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	// 停用分类
	category.Deactivate()

	// 保存分类
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	// 获取分类
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return ErrCategoryNotFound
	}

	// 检查是否有子分类
	children, err := s.categoryRepo.ListByParentID(ctx, &id)
	if err == nil && len(children) > 0 {
		return errors.New("存在子分类，无法删除")
	}

	// 检查是否有关联课程
	if category.CoursesCount > 0 {
		return errors.New("分类下存在课程，无法删除")
	}

	// 删除分类
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetCategoryTree 获取分类树
func (s *CategoryService) GetCategoryTree(ctx context.Context) ([]*entity.Category, error) {
	// 获取分类树
	return s.categoryRepo.GetTree(ctx)
}

// GetActiveCategories 获取所有激活状态的分类
func (s *CategoryService) GetActiveCategories(ctx context.Context) ([]*entity.Category, error) {
	// 获取所有激活状态的分类
	return s.categoryRepo.ListActive(ctx)
}

// GetCategoriesWithCourseCount 获取带课程数的分类列表
func (s *CategoryService) GetCategoriesWithCourseCount(ctx context.Context) ([]*entity.Category, error) {
	// 获取带课程数的分类
	return s.categoryRepo.ListWithCourseCount(ctx)
}

// GetCategoryStats 获取分类统计数据
func (s *CategoryService) GetCategoryStats(ctx context.Context) (totalCount, level1Count, level2Count int64, err error) {
	totalCount, err = s.categoryRepo.CountAll(ctx)
	if err != nil {
		return 0, 0, 0, err
	}

	level1Count, err = s.categoryRepo.CountByLevel(ctx, 1)
	if err != nil {
		return 0, 0, 0, err
	}

	level2Count, err = s.categoryRepo.CountByLevel(ctx, 2)
	if err != nil {
		return 0, 0, 0, err
	}

	return totalCount, level1Count, level2Count, nil
}
