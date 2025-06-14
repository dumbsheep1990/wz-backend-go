package sql

import (
	"context"
	"time"

	"gorm.io/gorm"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// ProgressPO 学习进度持久化对象
type ProgressPO struct {
	ID              string   `gorm:"primaryKey;column:id;size:36"`
	UserID          string   `gorm:"index;column:user_id;size:36;not null"`
	CourseID        string   `gorm:"index;column:course_id;size:36;not null"`
	LessonID        string   `gorm:"index;column:lesson_id;size:36;not null"`
	Status          string   `gorm:"column:status;size:20;not null"`
	WatchedDuration int      `gorm:"column:watched_duration;default:0"`
	TotalDuration   int      `gorm:"column:total_duration;default:0"`
	CompletionRate  float64  `gorm:"column:completion_rate;default:0"`
	LastWatchedAt   *int64   `gorm:"column:last_watched_at"`
	CompletedAt     *int64   `gorm:"column:completed_at"`
	CreatedAt       int64    `gorm:"column:created_at;not null"`
	UpdatedAt       int64    `gorm:"column:updated_at;not null"`
}

// TableName 表名
func (ProgressPO) TableName() string {
	return "progresses"
}

// SQLProgressRepository 学习进度仓储SQL实现
type SQLProgressRepository struct {
	db *gorm.DB
}

// NewSQLProgressRepository 创建学习进度仓储实例
func NewSQLProgressRepository(database database.Database) repository.ProgressRepository {
	if gormDB, ok := database.(*gorm.DB); ok {
		return &SQLProgressRepository{db: gormDB}
	}
	panic("database must be GORM instance")
}

// Create 创建学习进度
func (r *SQLProgressRepository) Create(ctx context.Context, progress *entity.Progress) error {
	po := r.toProgressPO(progress)
	return r.db.WithContext(ctx).Create(po).Error
}

// GetByID 根据ID获取学习进度
func (r *SQLProgressRepository) GetByID(ctx context.Context, id string) (*entity.Progress, error) {
	var po ProgressPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// Update 更新学习进度
func (r *SQLProgressRepository) Update(ctx context.Context, progress *entity.Progress) error {
	po := r.toProgressPO(progress)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除学习进度
func (r *SQLProgressRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&ProgressPO{}).Error
}

// GetByUserAndLesson 根据用户和课时获取学习进度
func (r *SQLProgressRepository) GetByUserAndLesson(ctx context.Context, userID, lessonID string) (*entity.Progress, error) {
	var po ProgressPO
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND lesson_id = ?", userID, lessonID).
		First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// ListByUserAndCourse 根据用户和课程查询学习进度列表
func (r *SQLProgressRepository) ListByUserAndCourse(ctx context.Context, userID, courseID string) ([]*entity.Progress, error) {
	var pos []ProgressPO
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Order("created_at ASC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	progresses := make([]*entity.Progress, len(pos))
	for i, po := range pos {
		progresses[i] = r.toDomainEntity(&po)
	}
	
	return progresses, nil
}

// ListByUser 根据用户查询学习进度列表
func (r *SQLProgressRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*entity.Progress, error) {
	var pos []ProgressPO
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	progresses := make([]*entity.Progress, len(pos))
	for i, po := range pos {
		progresses[i] = r.toDomainEntity(&po)
	}
	
	return progresses, nil
}

// ListByCourse 根据课程查询学习进度列表
func (r *SQLProgressRepository) ListByCourse(ctx context.Context, courseID string, limit, offset int) ([]*entity.Progress, error) {
	var pos []ProgressPO
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	progresses := make([]*entity.Progress, len(pos))
	for i, po := range pos {
		progresses[i] = r.toDomainEntity(&po)
	}
	
	return progresses, nil
}

// ListRecentProgress 查询用户最近学习进度
func (r *SQLProgressRepository) ListRecentProgress(ctx context.Context, userID string, limit int) ([]*entity.Progress, error) {
	var pos []ProgressPO
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND last_watched_at IS NOT NULL", userID).
		Order("last_watched_at DESC").
		Limit(limit).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	progresses := make([]*entity.Progress, len(pos))
	for i, po := range pos {
		progresses[i] = r.toDomainEntity(&po)
	}
	
	return progresses, nil
}

// BatchCreate 批量创建学习进度
func (r *SQLProgressRepository) BatchCreate(ctx context.Context, progresses []*entity.Progress) error {
	pos := make([]ProgressPO, len(progresses))
	for i, progress := range progresses {
		pos[i] = *r.toProgressPO(progress)
	}
	return r.db.WithContext(ctx).Create(&pos).Error
}

// BatchUpdate 批量更新学习进度
func (r *SQLProgressRepository) BatchUpdate(ctx context.Context, progresses []*entity.Progress) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, progress := range progresses {
			po := r.toProgressPO(progress)
			if err := tx.Save(po).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// CountByUserAndCourse 根据用户和课程统计学习进度数量
func (r *SQLProgressRepository) CountByUserAndCourse(ctx context.Context, userID, courseID string, status entity.ProgressStatus) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&ProgressPO{}).
		Where("user_id = ? AND course_id = ?", userID, courseID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Count(&count).Error
	return count, err
}

// CountByUser 根据用户统计学习进度数量
func (r *SQLProgressRepository) CountByUser(ctx context.Context, userID string, status entity.ProgressStatus) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&ProgressPO{}).Where("user_id = ?", userID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Count(&count).Error
	return count, err
}

// GetCourseProgressByUser 获取用户在某课程的整体进度
func (r *SQLProgressRepository) GetCourseProgressByUser(ctx context.Context, userID, courseID string) (float64, error) {
	var result struct {
		AvgProgress float64
	}
	
	err := r.db.WithContext(ctx).Model(&ProgressPO{}).
		Select("AVG(completion_rate) as avg_progress").
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Scan(&result).Error
	
	return result.AvgProgress, err
}

// GetOverallProgressByUser 获取用户的整体学习进度
func (r *SQLProgressRepository) GetOverallProgressByUser(ctx context.Context, userID string) (float64, error) {
	var result struct {
		AvgProgress float64
	}
	
	err := r.db.WithContext(ctx).Model(&ProgressPO{}).
		Select("AVG(completion_rate) as avg_progress").
		Where("user_id = ?", userID).
		Scan(&result).Error
	
	return result.AvgProgress, err
}

// 辅助方法：领域实体转持久化对象
func (r *SQLProgressRepository) toProgressPO(progress *entity.Progress) *ProgressPO {
	po := &ProgressPO{
		ID:              progress.ID,
		UserID:          progress.UserID,
		CourseID:        progress.CourseID,
		LessonID:        progress.LessonID,
		Status:          string(progress.Status),
		WatchedDuration: progress.WatchedDuration,
		TotalDuration:   progress.TotalDuration,
		CompletionRate:  progress.CompletionRate,
		CreatedAt:       progress.CreatedAt.Unix(),
		UpdatedAt:       progress.UpdatedAt.Unix(),
	}
	
	if progress.LastWatchedAt != nil {
		lastWatchedAt := progress.LastWatchedAt.Unix()
		po.LastWatchedAt = &lastWatchedAt
	}
	
	if progress.CompletedAt != nil {
		completedAt := progress.CompletedAt.Unix()
		po.CompletedAt = &completedAt
	}
	
	return po
}

// 辅助方法：持久化对象转领域实体
func (r *SQLProgressRepository) toDomainEntity(po *ProgressPO) *entity.Progress {
	progress := &entity.Progress{
		ID:              po.ID,
		UserID:          po.UserID,
		CourseID:        po.CourseID,
		LessonID:        po.LessonID,
		Status:          entity.ProgressStatus(po.Status),
		WatchedDuration: po.WatchedDuration,
		TotalDuration:   po.TotalDuration,
		CompletionRate:  po.CompletionRate,
		CreatedAt:       time.Unix(po.CreatedAt, 0),
		UpdatedAt:       time.Unix(po.UpdatedAt, 0),
	}
	
	if po.LastWatchedAt != nil {
		lastWatchedAt := time.Unix(*po.LastWatchedAt, 0)
		progress.LastWatchedAt = &lastWatchedAt
	}
	
	if po.CompletedAt != nil {
		completedAt := time.Unix(*po.CompletedAt, 0)
		progress.CompletedAt = &completedAt
	}
	
	return progress
}
