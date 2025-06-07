package repository

import (
	"context"

	"wz-backend-go/internal/domain/learn/entity"
)

// CategoryRepository 定义分类仓储接口
type CategoryRepository interface {
	// 基本CRUD操作
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id string) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	List(ctx context.Context) ([]*entity.Category, error)
	ListActive(ctx context.Context) ([]*entity.Category, error)
	ListByParentID(ctx context.Context, parentID *string) ([]*entity.Category, error)
	ListByLevel(ctx context.Context, level int) ([]*entity.Category, error)
	ListWithCourseCount(ctx context.Context) ([]*entity.Category, error)
	
	// 树形结构操作
	GetTree(ctx context.Context) ([]*entity.Category, error)
	
	// 统计操作
	CountAll(ctx context.Context) (int64, error)
	CountByLevel(ctx context.Context, level int) (int64, error)
}
