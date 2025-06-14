package sql

import (
	"context"
	"time"

	"gorm.io/gorm"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// LessonPO 课时持久化对象
type LessonPO struct {
	ID             string  `gorm:"primaryKey;column:id;size:36"`
	CourseID       string  `gorm:"index;column:course_id;size:36;not null"`
	ChapterID      string  `gorm:"index;column:chapter_id;size:36;not null"`
	Title          string  `gorm:"column:title;size:200;not null"`
	Description    string  `gorm:"column:description;type:text"`
	Type           string  `gorm:"column:type;size:20;not null"`
	Status         string  `gorm:"column:status;size:20;not null"`
	Order          int     `gorm:"column:order;not null"`
	Duration       int     `gorm:"column:duration;default:0"`
	VideoURL       string  `gorm:"column:video_url;size:500"`
	VideoSize      int64   `gorm:"column:video_size;default:0"`
	ArticleContent string  `gorm:"column:article_content;type:longtext"`
	AudioURL       string  `gorm:"column:audio_url;size:500"`
	AudioSize      int64   `gorm:"column:audio_size;default:0"`
	IsFree         bool    `gorm:"column:is_free;default:false"`
	CreatedAt      int64   `gorm:"column:created_at;not null"`
	UpdatedAt      int64   `gorm:"column:updated_at;not null"`
	PublishedAt    *int64  `gorm:"column:published_at"`
}

// TableName 表名
func (LessonPO) TableName() string {
	return "lessons"
}

// SQLLessonRepository 课时仓储SQL实现
type SQLLessonRepository struct {
	db *gorm.DB
}

// NewSQLLessonRepository 创建课时仓储实例
func NewSQLLessonRepository(database database.Database) repository.LessonRepository {
	if gormDB, ok := database.(*gorm.DB); ok {
		return &SQLLessonRepository{db: gormDB}
	}
	panic("database must be GORM instance")
}

// Create 创建课时
func (r *SQLLessonRepository) Create(ctx context.Context, lesson *entity.Lesson) error {
	po := r.toLessonPO(lesson)
	return r.db.WithContext(ctx).Create(po).Error
}

// GetByID 根据ID获取课时
func (r *SQLLessonRepository) GetByID(ctx context.Context, id string) (*entity.Lesson, error) {
	var po LessonPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// Update 更新课时
func (r *SQLLessonRepository) Update(ctx context.Context, lesson *entity.Lesson) error {
	po := r.toLessonPO(lesson)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除课时
func (r *SQLLessonRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&LessonPO{}).Error
}

// ListByChapterID 根据章节ID查询课时列表
func (r *SQLLessonRepository) ListByChapterID(ctx context.Context, chapterID string) ([]*entity.Lesson, error) {
	var pos []LessonPO
	err := r.db.WithContext(ctx).
		Where("chapter_id = ?", chapterID).
		Order("`order` ASC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	lessons := make([]*entity.Lesson, len(pos))
	for i, po := range pos {
		lessons[i] = r.toDomainEntity(&po)
	}
	
	return lessons, nil
}

// ListByCourseID 根据课程ID查询课时列表
func (r *SQLLessonRepository) ListByCourseID(ctx context.Context, courseID string) ([]*entity.Lesson, error) {
	var pos []LessonPO
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("`order` ASC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	lessons := make([]*entity.Lesson, len(pos))
	for i, po := range pos {
		lessons[i] = r.toDomainEntity(&po)
	}
	
	return lessons, nil
}

// ListFreeLessons 根据课程ID查询免费课时列表
func (r *SQLLessonRepository) ListFreeLessons(ctx context.Context, courseID string) ([]*entity.Lesson, error) {
	var pos []LessonPO
	err := r.db.WithContext(ctx).
		Where("course_id = ? AND is_free = ? AND status = ?", courseID, true, entity.LessonStatusPublished).
		Order("`order` ASC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	lessons := make([]*entity.Lesson, len(pos))
	for i, po := range pos {
		lessons[i] = r.toDomainEntity(&po)
	}
	
	return lessons, nil
}

// BatchCreate 批量创建课时
func (r *SQLLessonRepository) BatchCreate(ctx context.Context, lessons []*entity.Lesson) error {
	pos := make([]LessonPO, len(lessons))
	for i, lesson := range lessons {
		pos[i] = *r.toLessonPO(lesson)
	}
	return r.db.WithContext(ctx).Create(&pos).Error
}

// BatchUpdate 批量更新课时
func (r *SQLLessonRepository) BatchUpdate(ctx context.Context, lessons []*entity.Lesson) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, lesson := range lessons {
			po := r.toLessonPO(lesson)
			if err := tx.Save(po).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// CountByChapterID 根据章节ID统计课时数量
func (r *SQLLessonRepository) CountByChapterID(ctx context.Context, chapterID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&LessonPO{}).Where("chapter_id = ?", chapterID).Count(&count).Error
	return count, err
}

// CountByCourseID 根据课程ID统计课时数量
func (r *SQLLessonRepository) CountByCourseID(ctx context.Context, courseID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&LessonPO{}).Where("course_id = ?", courseID).Count(&count).Error
	return count, err
}

// CountPublishedByCourseID 根据课程ID统计已发布课时数量
func (r *SQLLessonRepository) CountPublishedByCourseID(ctx context.Context, courseID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&LessonPO{}).
		Where("course_id = ? AND status = ?", courseID, entity.LessonStatusPublished).
		Count(&count).Error
	return count, err
}

// 辅助方法：领域实体转持久化对象
func (r *SQLLessonRepository) toLessonPO(lesson *entity.Lesson) *LessonPO {
	po := &LessonPO{
		ID:             lesson.ID,
		CourseID:       lesson.CourseID,
		ChapterID:      lesson.ChapterID,
		Title:          lesson.Title,
		Description:    lesson.Description,
		Type:           string(lesson.Type),
		Status:         string(lesson.Status),
		Order:          lesson.Order,
		Duration:       lesson.Duration,
		VideoURL:       lesson.VideoURL,
		VideoSize:      lesson.VideoSize,
		ArticleContent: lesson.ArticleContent,
		AudioURL:       lesson.AudioURL,
		AudioSize:      lesson.AudioSize,
		IsFree:         lesson.IsFree,
		CreatedAt:      lesson.CreatedAt.Unix(),
		UpdatedAt:      lesson.UpdatedAt.Unix(),
	}
	
	if lesson.PublishedAt != nil {
		publishedAt := lesson.PublishedAt.Unix()
		po.PublishedAt = &publishedAt
	}
	
	return po
}

// 辅助方法：持久化对象转领域实体
func (r *SQLLessonRepository) toDomainEntity(po *LessonPO) *entity.Lesson {
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
