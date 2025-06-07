package repository

import (
	"context"
	"time"

	"wz-backend-go/internal/domain/learn/entity"
)

// EnrollmentRepository 定义课程报名仓储接口
type EnrollmentRepository interface {
	// 基本CRUD操作
	Create(ctx context.Context, enrollment *entity.Enrollment) error
	GetByID(ctx context.Context, id string) (*entity.Enrollment, error)
	GetByUserAndCourse(ctx context.Context, userID, courseID string) (*entity.Enrollment, error)
	Update(ctx context.Context, enrollment *entity.Enrollment) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	ListByUserID(ctx context.Context, userID string, params EnrollmentQueryParams) ([]*entity.Enrollment, int64, error)
	ListByCourseID(ctx context.Context, courseID string, params EnrollmentQueryParams) ([]*entity.Enrollment, int64, error)
	ListExpired(ctx context.Context, before time.Time) ([]*entity.Enrollment, error)
	
	// 统计操作
	CountByUserID(ctx context.Context, userID string) (int64, error)
	CountByCourseID(ctx context.Context, courseID string) (int64, error)
	CountByStatus(ctx context.Context, status entity.EnrollmentStatus) (int64, error)
	CountActiveEnrollments(ctx context.Context) (int64, error)
	CountRecentEnrollments(ctx context.Context, since time.Time) (int64, error)
}

// EnrollmentQueryParams 报名查询参数
type EnrollmentQueryParams struct {
	Page      int                     // 页码
	PageSize  int                     // 每页数量
	SortBy    string                  // 排序字段
	SortOrder string                  // 排序顺序 (asc,desc)
	Status    *entity.EnrollmentStatus   // 状态过滤
	Since     *time.Time              // 开始时间
	Until     *time.Time              // 结束时间
}
