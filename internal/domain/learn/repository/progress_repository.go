package repository

import (
	"context"

	"wz-backend-go/internal/domain/learn/entity"
)

// ProgressRepository 定义学习进度仓储接口
type ProgressRepository interface {
	// 基本CRUD操作
	Create(ctx context.Context, progress *entity.Progress) error
	GetByID(ctx context.Context, id string) (*entity.Progress, error)
	Update(ctx context.Context, progress *entity.Progress) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	GetByUserAndLesson(ctx context.Context, userID, lessonID string) (*entity.Progress, error)
	ListByUserAndCourse(ctx context.Context, userID, courseID string) ([]*entity.Progress, error)
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]*entity.Progress, error)
	ListByCourse(ctx context.Context, courseID string, limit, offset int) ([]*entity.Progress, error)
	ListRecentProgress(ctx context.Context, userID string, limit int) ([]*entity.Progress, error)
	
	// 批量操作
	BatchCreate(ctx context.Context, progresses []*entity.Progress) error
	BatchUpdate(ctx context.Context, progresses []*entity.Progress) error
	
	// 统计操作
	CountByUserAndCourse(ctx context.Context, userID, courseID string, status entity.ProgressStatus) (int64, error)
	CountByUser(ctx context.Context, userID string, status entity.ProgressStatus) (int64, error)
	GetCourseProgressByUser(ctx context.Context, userID, courseID string) (float64, error)
	GetOverallProgressByUser(ctx context.Context, userID string) (float64, error)
}
