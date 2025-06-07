package repository

import (
	"context"

	"wz-backend-go/internal/domain/learn/entity"
)

// CourseRepository 定义课程仓储接口
type CourseRepository interface {
	// 基本CRUD操作
	Create(ctx context.Context, course *entity.Course) error
	GetByID(ctx context.Context, id string) (*entity.Course, error)
	Update(ctx context.Context, course *entity.Course) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	List(ctx context.Context, params CourseQueryParams) ([]*entity.Course, int64, error)
	ListByTeacherID(ctx context.Context, teacherID string, params CourseQueryParams) ([]*entity.Course, int64, error)
	ListByIDs(ctx context.Context, ids []string) ([]*entity.Course, error)
	ListByCategoryID(ctx context.Context, categoryID string, params CourseQueryParams) ([]*entity.Course, int64, error)
	ListPopular(ctx context.Context, limit int) ([]*entity.Course, error)
	ListRecent(ctx context.Context, limit int) ([]*entity.Course, error)
	
	// 统计操作
	CountAll(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status entity.CourseStatus) (int64, error)
	CountByTeacherID(ctx context.Context, teacherID string) (int64, error)
	CountByCategoryID(ctx context.Context, categoryID string) (int64, error)
	
	// 搜索操作
	Search(ctx context.Context, keyword string, params CourseQueryParams) ([]*entity.Course, int64, error)
}

// CourseQueryParams 课程查询参数
type CourseQueryParams struct {
	Page      int               // 页码
	PageSize  int               // 每页数量
	SortBy    string            // 排序字段
	SortOrder string            // 排序顺序 (asc,desc)
	Status    *entity.CourseStatus // 状态过滤
	Level     *entity.CourseLevel  // 难度级别过滤
	PriceRange [2]float64       // 价格范围 [min, max]
	Tags      []string          // 标签过滤
	FreeCourses *bool           // 是否只查询免费课程
}
