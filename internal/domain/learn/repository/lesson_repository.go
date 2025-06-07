package repository

import (
	"context"

	"wz-backend-go/internal/domain/learn/entity"
)

// LessonRepository 定义课时仓储接口
type LessonRepository interface {
	// 基本CRUD操作
	Create(ctx context.Context, lesson *entity.Lesson) error
	GetByID(ctx context.Context, id string) (*entity.Lesson, error)
	Update(ctx context.Context, lesson *entity.Lesson) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	ListByChapterID(ctx context.Context, chapterID string) ([]*entity.Lesson, error)
	ListByCourseID(ctx context.Context, courseID string) ([]*entity.Lesson, error)
	ListFreeLessons(ctx context.Context, courseID string) ([]*entity.Lesson, error)
	
	// 批量操作
	BatchCreate(ctx context.Context, lessons []*entity.Lesson) error
	BatchUpdate(ctx context.Context, lessons []*entity.Lesson) error
	
	// 统计操作
	CountByChapterID(ctx context.Context, chapterID string) (int64, error)
	CountByCourseID(ctx context.Context, courseID string) (int64, error)
	CountPublishedByCourseID(ctx context.Context, courseID string) (int64, error)
}
