package sql

import (
	"context"
	"time"

	"gorm.io/gorm"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// ChapterPO 章节持久化对象
type ChapterPO struct {
	ID          string `gorm:"primaryKey;column:id;size:36"`
	CourseID    string `gorm:"index;column:course_id;size:36;not null"`
	Title       string `gorm:"column:title;size:200;not null"`
	Description string `gorm:"column:description;type:text"`
	Order       int    `gorm:"column:order;not null"`
	LessonCount int    `gorm:"column:lesson_count;default:0"`
	Duration    int    `gorm:"column:duration;default:0"`
	CreatedAt   int64  `gorm:"column:created_at;not null"`
	UpdatedAt   int64  `gorm:"column:updated_at;not null"`
}

// TableName 表名
func (ChapterPO) TableName() string {
	return "chapters"
}

// SQLChapterRepository 章节仓储SQL实现
type SQLChapterRepository struct {
	db *gorm.DB
}

// NewSQLChapterRepository 创建章节仓储实例
func NewSQLChapterRepository(database database.Database) repository.ChapterRepository {
	if gormDB, ok := database.(*gorm.DB); ok {
		return &SQLChapterRepository{db: gormDB}
	}
	panic("database must be GORM instance")
}

// Create 创建章节
func (r *SQLChapterRepository) Create(ctx context.Context, chapter *entity.Chapter) error {
	po := r.toChapterPO(chapter)
	return r.db.WithContext(ctx).Create(po).Error
}

// GetByID 根据ID获取章节
func (r *SQLChapterRepository) GetByID(ctx context.Context, id string) (*entity.Chapter, error) {
	var po ChapterPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// Update 更新章节
func (r *SQLChapterRepository) Update(ctx context.Context, chapter *entity.Chapter) error {
	po := r.toChapterPO(chapter)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除章节
func (r *SQLChapterRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&ChapterPO{}).Error
}

// ListByCourseID 根据课程ID查询章节列表
func (r *SQLChapterRepository) ListByCourseID(ctx context.Context, courseID string) ([]*entity.Chapter, error) {
	var pos []ChapterPO
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("`order` ASC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	chapters := make([]*entity.Chapter, len(pos))
	for i, po := range pos {
		chapters[i] = r.toDomainEntity(&po)
	}
	
	return chapters, nil
}

// GetChapterWithLessons 获取章节及其课时列表
func (r *SQLChapterRepository) GetChapterWithLessons(ctx context.Context, id string) (*entity.Chapter, []*entity.Lesson, error) {
	// 获取章节
	chapter, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if chapter == nil {
		return nil, nil, nil
	}
	
	// 获取课时列表
	var lessonPOs []LessonPO
	err = r.db.WithContext(ctx).
		Where("chapter_id = ?", id).
		Order("`order` ASC").
		Find(&lessonPOs).Error
	if err != nil {
		return nil, nil, err
	}
	
	lessons := make([]*entity.Lesson, len(lessonPOs))
	for i, po := range lessonPOs {
		lessons[i] = r.toLessonDomainEntity(&po)
	}
	
	return chapter, lessons, nil
}

// BatchCreate 批量创建章节
func (r *SQLChapterRepository) BatchCreate(ctx context.Context, chapters []*entity.Chapter) error {
	pos := make([]ChapterPO, len(chapters))
	for i, chapter := range chapters {
		pos[i] = *r.toChapterPO(chapter)
	}
	return r.db.WithContext(ctx).Create(&pos).Error
}

// BatchUpdate 批量更新章节
func (r *SQLChapterRepository) BatchUpdate(ctx context.Context, chapters []*entity.Chapter) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, chapter := range chapters {
			po := r.toChapterPO(chapter)
			if err := tx.Save(po).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// CountByCourseID 根据课程ID统计章节数量
func (r *SQLChapterRepository) CountByCourseID(ctx context.Context, courseID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ChapterPO{}).Where("course_id = ?", courseID).Count(&count).Error
	return count, err
}

// 辅助方法：领域实体转持久化对象
func (r *SQLChapterRepository) toChapterPO(chapter *entity.Chapter) *ChapterPO {
	return &ChapterPO{
		ID:          chapter.ID,
		CourseID:    chapter.CourseID,
		Title:       chapter.Title,
		Description: chapter.Description,
		Order:       chapter.Order,
		LessonCount: chapter.LessonCount,
		Duration:    chapter.Duration,
		CreatedAt:   chapter.CreatedAt.Unix(),
		UpdatedAt:   chapter.UpdatedAt.Unix(),
	}
}

// 辅助方法：持久化对象转领域实体
func (r *SQLChapterRepository) toDomainEntity(po *ChapterPO) *entity.Chapter {
	return &entity.Chapter{
		ID:          po.ID,
		CourseID:    po.CourseID,
		Title:       po.Title,
		Description: po.Description,
		Order:       po.Order,
		LessonCount: po.LessonCount,
		Duration:    po.Duration,
		CreatedAt:   time.Unix(po.CreatedAt, 0),
		UpdatedAt:   time.Unix(po.UpdatedAt, 0),
	}
}

// 辅助方法：课时持久化对象转领域实体（用于GetChapterWithLessons方法）
func (r *SQLChapterRepository) toLessonDomainEntity(po *LessonPO) *entity.Lesson {
	lesson := &entity.Lesson{
		ID:             po.ID,
		CourseID:       po.CourseID,
		ChapterID:      po.ChapterID,
		Title:          po.Title,
		Description:    po.Description,
		Type:           entity.LessonType(po.Type),
		Status:         entity.LessonStatus(po.Status),
		Order:          po.Order,
		Duration:       po.Duration,
		VideoURL:       po.VideoURL,
		VideoSize:      po.VideoSize,
		ArticleContent: po.ArticleContent,
		AudioURL:       po.AudioURL,
		AudioSize:      po.AudioSize,
		IsFree:         po.IsFree,
		CreatedAt:      time.Unix(po.CreatedAt, 0),
		UpdatedAt:      time.Unix(po.UpdatedAt, 0),
	}
	
	if po.PublishedAt != nil {
		publishedAt := time.Unix(*po.PublishedAt, 0)
		lesson.PublishedAt = &publishedAt
	}
	
	return lesson
}
