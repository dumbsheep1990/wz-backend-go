package repository

import (
	"context"

	"wz-backend-go/internal/domain/learn/entity"
)

// TeacherRepository 定义讲师仓储接口
type TeacherRepository interface {
	// 基本CRUD操作
	Create(ctx context.Context, teacher *entity.Teacher) error
	GetByID(ctx context.Context, id string) (*entity.Teacher, error)
	GetByUserID(ctx context.Context, userID string) (*entity.Teacher, error)
	Update(ctx context.Context, teacher *entity.Teacher) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	List(ctx context.Context, params TeacherQueryParams) ([]*entity.Teacher, int64, error)
	ListByIDs(ctx context.Context, ids []string) ([]*entity.Teacher, error)
	ListPopular(ctx context.Context, limit int) ([]*entity.Teacher, error)
	
	// 统计操作
	CountAll(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status entity.TeacherStatus) (int64, error)
	
	// 搜索操作
	Search(ctx context.Context, keyword string, params TeacherQueryParams) ([]*entity.Teacher, int64, error)
}

// TeacherQueryParams 讲师查询参数
type TeacherQueryParams struct {
	Page      int                 // 页码
	PageSize  int                 // 每页数量
	SortBy    string              // 排序字段
	SortOrder string              // 排序顺序 (asc,desc)
	Status    *entity.TeacherStatus  // 状态过滤
	Specialty *string             // 专长领域过滤
}
