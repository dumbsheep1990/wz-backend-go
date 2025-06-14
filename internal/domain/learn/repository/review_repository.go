package repository

import (
	"context"

	"wz-backend-go/internal/domain/learn/entity"
)

// ReviewRepository 定义评价仓储接口
type ReviewRepository interface {
	// 基本CRUD操作
	Create(ctx context.Context, review *entity.Review) error
	GetByID(ctx context.Context, id string) (*entity.Review, error)
	Update(ctx context.Context, review *entity.Review) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	ListByCourseID(ctx context.Context, courseID string, status entity.ReviewStatus, limit, offset int) ([]*entity.Review, error)
	ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.Review, error)
	GetByUserAndCourse(ctx context.Context, userID, courseID string) (*entity.Review, error)
	ListPendingReviews(ctx context.Context, limit, offset int) ([]*entity.Review, error)
	
	// 统计操作
	CountByCourseID(ctx context.Context, courseID string, status entity.ReviewStatus) (int64, error)
	CountByUserID(ctx context.Context, userID string) (int64, error)
	GetAverageRatingByCourseID(ctx context.Context, courseID string) (float64, error)
	GetRatingDistributionByCourseID(ctx context.Context, courseID string) (map[int]int64, error)
}
