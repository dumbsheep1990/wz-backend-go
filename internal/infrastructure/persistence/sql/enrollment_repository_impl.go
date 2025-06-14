package sql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// EnrollmentPO 报名持久化对象
type EnrollmentPO struct {
	ID             string   `gorm:"primaryKey;column:id;size:36"`
	CourseID       string   `gorm:"index;column:course_id;size:36;not null"`
	UserID         string   `gorm:"index;column:user_id;size:36;not null"`
	OrderID        string   `gorm:"index;column:order_id;size:36;not null"`
	Status         string   `gorm:"column:status;size:20;not null"`
	Progress       float64  `gorm:"column:progress;default:0"`
	CompletedCount int      `gorm:"column:completed_count;default:0"`
	TotalCount     int      `gorm:"column:total_count;default:0"`
	LastLearnTime  *int64   `gorm:"column:last_learn_time"`
	Rating         *float64 `gorm:"column:rating"`
	Comment        string   `gorm:"column:comment;type:text"`
	ExpiresAt      *int64   `gorm:"column:expires_at"`
	CreatedAt      int64    `gorm:"column:created_at;not null"`
	UpdatedAt      int64    `gorm:"column:updated_at;not null"`
	CompletedAt    *int64   `gorm:"column:completed_at"`
}

// TableName 表名
func (EnrollmentPO) TableName() string {
	return "enrollments"
}

// SQLEnrollmentRepository 报名仓储SQL实现
type SQLEnrollmentRepository struct {
	db *gorm.DB
}

// NewSQLEnrollmentRepository 创建报名仓储实例
func NewSQLEnrollmentRepository(database database.Database) repository.EnrollmentRepository {
	if gormDB, ok := database.(*gorm.DB); ok {
		return &SQLEnrollmentRepository{db: gormDB}
	}
	panic("database must be GORM instance")
}

// Create 创建报名记录
func (r *SQLEnrollmentRepository) Create(ctx context.Context, enrollment *entity.Enrollment) error {
	po := r.toEnrollmentPO(enrollment)
	return r.db.WithContext(ctx).Create(po).Error
}

// GetByID 根据ID获取报名记录
func (r *SQLEnrollmentRepository) GetByID(ctx context.Context, id string) (*entity.Enrollment, error) {
	var po EnrollmentPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// GetByUserAndCourse 根据用户和课程获取报名记录
func (r *SQLEnrollmentRepository) GetByUserAndCourse(ctx context.Context, userID, courseID string) (*entity.Enrollment, error) {
	var po EnrollmentPO
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// Update 更新报名记录
func (r *SQLEnrollmentRepository) Update(ctx context.Context, enrollment *entity.Enrollment) error {
	po := r.toEnrollmentPO(enrollment)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除报名记录
func (r *SQLEnrollmentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&EnrollmentPO{}).Error
}

// ListByUserID 根据用户ID查询报名列表
func (r *SQLEnrollmentRepository) ListByUserID(ctx context.Context, userID string, params repository.EnrollmentQueryParams) ([]*entity.Enrollment, int64, error) {
	query := r.db.WithContext(ctx).Model(&EnrollmentPO{}).Where("user_id = ?", userID)
	
	// 构建查询条件
	query = r.buildWhereClause(query, params)
	
	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 构建排序和分页
	query = r.buildOrderClause(query, params)
	query = r.buildLimitClause(query, params)
	
	var pos []EnrollmentPO
	if err := query.Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	
	enrollments := make([]*entity.Enrollment, len(pos))
	for i, po := range pos {
		enrollments[i] = r.toDomainEntity(&po)
	}
	
	return enrollments, total, nil
}

// ListByCourseID 根据课程ID查询报名列表
func (r *SQLEnrollmentRepository) ListByCourseID(ctx context.Context, courseID string, params repository.EnrollmentQueryParams) ([]*entity.Enrollment, int64, error) {
	query := r.db.WithContext(ctx).Model(&EnrollmentPO{}).Where("course_id = ?", courseID)
	
	// 构建查询条件
	query = r.buildWhereClause(query, params)
	
	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 构建排序和分页
	query = r.buildOrderClause(query, params)
	query = r.buildLimitClause(query, params)
	
	var pos []EnrollmentPO
	if err := query.Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	
	enrollments := make([]*entity.Enrollment, len(pos))
	for i, po := range pos {
		enrollments[i] = r.toDomainEntity(&po)
	}
	
	return enrollments, total, nil
}

// ListExpired 查询过期的报名记录
func (r *SQLEnrollmentRepository) ListExpired(ctx context.Context, before time.Time) ([]*entity.Enrollment, error) {
	beforeUnix := before.Unix()
	var pos []EnrollmentPO
	err := r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ? AND status = ?", beforeUnix, entity.EnrollmentStatusActive).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	enrollments := make([]*entity.Enrollment, len(pos))
	for i, po := range pos {
		enrollments[i] = r.toDomainEntity(&po)
	}
	
	return enrollments, nil
}

// CountByUserID 根据用户ID统计报名数量
func (r *SQLEnrollmentRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&EnrollmentPO{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountByCourseID 根据课程ID统计报名数量
func (r *SQLEnrollmentRepository) CountByCourseID(ctx context.Context, courseID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&EnrollmentPO{}).Where("course_id = ?", courseID).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态统计报名数量
func (r *SQLEnrollmentRepository) CountByStatus(ctx context.Context, status entity.EnrollmentStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&EnrollmentPO{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// CountActiveEnrollments 统计活跃报名数量
func (r *SQLEnrollmentRepository) CountActiveEnrollments(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&EnrollmentPO{}).Where("status = ?", entity.EnrollmentStatusActive).Count(&count).Error
	return count, err
}

// CountRecentEnrollments 统计最近的报名数量
func (r *SQLEnrollmentRepository) CountRecentEnrollments(ctx context.Context, since time.Time) (int64, error) {
	sinceUnix := since.Unix()
	var count int64
	err := r.db.WithContext(ctx).Model(&EnrollmentPO{}).Where("created_at >= ?", sinceUnix).Count(&count).Error
	return count, err
}

// 辅助方法：构建WHERE子句
func (r *SQLEnrollmentRepository) buildWhereClause(query *gorm.DB, params repository.EnrollmentQueryParams) *gorm.DB {
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	
	if params.Since != nil {
		sinceUnix := params.Since.Unix()
		query = query.Where("created_at >= ?", sinceUnix)
	}
	
	if params.Until != nil {
		untilUnix := params.Until.Unix()
		query = query.Where("created_at <= ?", untilUnix)
	}
	
	return query
}

// 辅助方法：构建ORDER BY子句
func (r *SQLEnrollmentRepository) buildOrderClause(query *gorm.DB, params repository.EnrollmentQueryParams) *gorm.DB {
	if params.SortBy == "" {
		return query.Order("created_at DESC")
	}
	
	order := "ASC"
	if strings.ToUpper(params.SortOrder) == "DESC" {
		order = "DESC"
	}
	
	return query.Order(fmt.Sprintf("%s %s", params.SortBy, order))
}

// 辅助方法：构建LIMIT子句
func (r *SQLEnrollmentRepository) buildLimitClause(query *gorm.DB, params repository.EnrollmentQueryParams) *gorm.DB {
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	
	offset := (params.Page - 1) * params.PageSize
	return query.Limit(params.PageSize).Offset(offset)
}

// 辅助方法：领域实体转持久化对象
func (r *SQLEnrollmentRepository) toEnrollmentPO(enrollment *entity.Enrollment) *EnrollmentPO {
	po := &EnrollmentPO{
		ID:             enrollment.ID,
		CourseID:       enrollment.CourseID,
		UserID:         enrollment.UserID,
		OrderID:        enrollment.OrderID,
		Status:         string(enrollment.Status),
		Progress:       enrollment.Progress,
		CompletedCount: enrollment.CompletedCount,
		TotalCount:     enrollment.TotalCount,
		Rating:         enrollment.Rating,
		Comment:        enrollment.Comment,
		CreatedAt:      enrollment.CreatedAt.Unix(),
		UpdatedAt:      enrollment.UpdatedAt.Unix(),
	}
	
	if enrollment.LastLearnTime != nil {
		lastLearnTime := enrollment.LastLearnTime.Unix()
		po.LastLearnTime = &lastLearnTime
	}
	
	if enrollment.ExpiresAt != nil {
		expiresAt := enrollment.ExpiresAt.Unix()
		po.ExpiresAt = &expiresAt
	}
	
	if enrollment.CompletedAt != nil {
		completedAt := enrollment.CompletedAt.Unix()
		po.CompletedAt = &completedAt
	}
	
	return po
}

// 辅助方法：持久化对象转领域实体
func (r *SQLEnrollmentRepository) toDomainEntity(po *EnrollmentPO) *entity.Enrollment {
	enrollment := &entity.Enrollment{
		ID:             po.ID,
		CourseID:       po.CourseID,
		UserID:         po.UserID,
		OrderID:        po.OrderID,
		Status:         entity.EnrollmentStatus(po.Status),
		Progress:       po.Progress,
		CompletedCount: po.CompletedCount,
		TotalCount:     po.TotalCount,
		Rating:         po.Rating,
		Comment:        po.Comment,
		CreatedAt:      time.Unix(po.CreatedAt, 0),
		UpdatedAt:      time.Unix(po.UpdatedAt, 0),
	}
	
	if po.LastLearnTime != nil {
		lastLearnTime := time.Unix(*po.LastLearnTime, 0)
		enrollment.LastLearnTime = &lastLearnTime
	}
	
	if po.ExpiresAt != nil {
		expiresAt := time.Unix(*po.ExpiresAt, 0)
		enrollment.ExpiresAt = &expiresAt
	}
	
	if po.CompletedAt != nil {
		completedAt := time.Unix(*po.CompletedAt, 0)
		enrollment.CompletedAt = &completedAt
	}
	
	return enrollment
}
