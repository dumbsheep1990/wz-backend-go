package sql

import (
	"context"
	"time"

	"gorm.io/gorm"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// CategoryPO 分类持久化对象
type CategoryPO struct {
	ID           string  `gorm:"primaryKey;column:id;size:36"`
	Name         string  `gorm:"column:name;size:100;not null"`
	Description  string  `gorm:"column:description;type:text"`
	Icon         string  `gorm:"column:icon;size:255"`
	ParentID     *string `gorm:"index;column:parent_id;size:36"`
	Level        int     `gorm:"column:level;not null"`
	Order        int     `gorm:"column:order;default:0"`
	CoursesCount int     `gorm:"column:courses_count;default:0"`
	IsActive     bool    `gorm:"column:is_active;default:true"`
	CreatedAt    int64   `gorm:"column:created_at;not null"`
	UpdatedAt    int64   `gorm:"column:updated_at;not null"`
}

// TableName 表名
func (CategoryPO) TableName() string {
	return "categories"
}

// SQLCategoryRepository 分类仓储SQL实现
type SQLCategoryRepository struct {
	db *gorm.DB
}

// NewSQLCategoryRepository 创建分类仓储实例
func NewSQLCategoryRepository(database database.Database) repository.CategoryRepository {
	if gormDB, ok := database.(*gorm.DB); ok {
		return &SQLCategoryRepository{db: gormDB}
	}
	panic("database must be GORM instance")
}

// Create 创建分类
func (r *SQLCategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	po := r.toCategoryPO(category)
	return r.db.WithContext(ctx).Create(po).Error
}

// GetByID 根据ID获取分类
func (r *SQLCategoryRepository) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	var po CategoryPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// Update 更新分类
func (r *SQLCategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	po := r.toCategoryPO(category)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除分类
func (r *SQLCategoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&CategoryPO{}).Error
}

// List 查询所有分类
func (r *SQLCategoryRepository) List(ctx context.Context) ([]*entity.Category, error) {
	var pos []CategoryPO
	err := r.db.WithContext(ctx).Order("level ASC, `order` ASC").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	categories := make([]*entity.Category, len(pos))
	for i, po := range pos {
		categories[i] = r.toDomainEntity(&po)
	}
	
	return categories, nil
}

// ListActive 查询活跃分类
func (r *SQLCategoryRepository) ListActive(ctx context.Context) ([]*entity.Category, error) {
	var pos []CategoryPO
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("level ASC, `order` ASC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	categories := make([]*entity.Category, len(pos))
	for i, po := range pos {
		categories[i] = r.toDomainEntity(&po)
	}
	
	return categories, nil
}

// ListByParentID 根据父级ID查询分类
func (r *SQLCategoryRepository) ListByParentID(ctx context.Context, parentID *string) ([]*entity.Category, error) {
	query := r.db.WithContext(ctx).Order("`order` ASC")
	
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	
	var pos []CategoryPO
	err := query.Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	categories := make([]*entity.Category, len(pos))
	for i, po := range pos {
		categories[i] = r.toDomainEntity(&po)
	}
	
	return categories, nil
}

// ListByLevel 根据层级查询分类
func (r *SQLCategoryRepository) ListByLevel(ctx context.Context, level int) ([]*entity.Category, error) {
	var pos []CategoryPO
	err := r.db.WithContext(ctx).
		Where("level = ?", level).
		Order("`order` ASC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	categories := make([]*entity.Category, len(pos))
	for i, po := range pos {
		categories[i] = r.toDomainEntity(&po)
	}
	
	return categories, nil
}

// ListWithCourseCount 查询带课程数量的分类
func (r *SQLCategoryRepository) ListWithCourseCount(ctx context.Context) ([]*entity.Category, error) {
	// 这里可以使用子查询来获取实际的课程数量，或者依赖于定期更新的courses_count字段
	var pos []CategoryPO
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("courses_count DESC, level ASC, `order` ASC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	categories := make([]*entity.Category, len(pos))
	for i, po := range pos {
		categories[i] = r.toDomainEntity(&po)
	}
	
	return categories, nil
}

// GetTree 获取分类树形结构
func (r *SQLCategoryRepository) GetTree(ctx context.Context) ([]*entity.Category, error) {
	// 获取所有活跃分类
	categories, err := r.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	
	// 构建树形结构（这里返回扁平列表，前端可以根据ParentID构建树）
	return categories, nil
}

// CountAll 统计所有分类数量
func (r *SQLCategoryRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&CategoryPO{}).Count(&count).Error
	return count, err
}

// CountByLevel 根据层级统计分类数量
func (r *SQLCategoryRepository) CountByLevel(ctx context.Context, level int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&CategoryPO{}).Where("level = ?", level).Count(&count).Error
	return count, err
}

// 辅助方法：领域实体转持久化对象
func (r *SQLCategoryRepository) toCategoryPO(category *entity.Category) *CategoryPO {
	return &CategoryPO{
		ID:           category.ID,
		Name:         category.Name,
		Description:  category.Description,
		Icon:         category.Icon,
		ParentID:     category.ParentID,
		Level:        category.Level,
		Order:        category.Order,
		CoursesCount: category.CoursesCount,
		IsActive:     category.IsActive,
		CreatedAt:    category.CreatedAt.Unix(),
		UpdatedAt:    category.UpdatedAt.Unix(),
	}
}

// 辅助方法：持久化对象转领域实体
func (r *SQLCategoryRepository) toDomainEntity(po *CategoryPO) *entity.Category {
	return &entity.Category{
		ID:           po.ID,
		Name:         po.Name,
		Description:  po.Description,
		Icon:         po.Icon,
		ParentID:     po.ParentID,
		Level:        po.Level,
		Order:        po.Order,
		CoursesCount: po.CoursesCount,
		IsActive:     po.IsActive,
		CreatedAt:    time.Unix(po.CreatedAt, 0),
		UpdatedAt:    time.Unix(po.UpdatedAt, 0),
	}
}
