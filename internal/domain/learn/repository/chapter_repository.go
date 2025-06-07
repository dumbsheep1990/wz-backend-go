package repository

import (
	"context"

	"wz-backend-go/internal/domain/learn/entity"
)

// ChapterRepository 定义章节仓储接口
type ChapterRepository interface {
	// 基本CRUD操作
	Create(ctx context.Context, chapter *entity.Chapter) error
	GetByID(ctx context.Context, id string) (*entity.Chapter, error)
	Update(ctx context.Context, chapter *entity.Chapter) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	ListByCourseID(ctx context.Context, courseID string) ([]*entity.Chapter, error)
	GetChapterWithLessons(ctx context.Context, id string) (*entity.Chapter, []*entity.Lesson, error)
	
	// 批量操作
	BatchCreate(ctx context.Context, chapters []*entity.Chapter) error
	BatchUpdate(ctx context.Context, chapters []*entity.Chapter) error
	
	// 统计操作
	CountByCourseID(ctx context.Context, courseID string) (int64, error)
}
