package sql

import (
	"context"
	"time"

	"gorm.io/gorm"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// ReviewPO 评价持久化对象
type ReviewPO struct {
	ID         string  `gorm:"primaryKey;column:id;size:36"`
	UserID     string  `gorm:"index;column:user_id;size:36;not null"`
	CourseID   string  `gorm:"index;column:course_id;size:36;not null"`
	Rating     int     `gorm:"column:rating;not null"`
	Content    string  `gorm:"column:content;type:text"`
	Status     string  `gorm:"column:status;size:20;not null"`
	CreatedAt  int64   `gorm:"column:created_at;not null"`
	UpdatedAt  int64   `gorm:"column:updated_at;not null"`
	ApprovedAt *int64  `gorm:"column:approved_at"`
}

// TableName 表名
func (ReviewPO) TableName() string {
	return "reviews"
}

// SQLReviewRepository 评价仓储SQL实现
type SQLReviewRepository struct {
	db *gorm.DB
}

// NewSQLReviewRepository 创建评价仓储实例
func NewSQLReviewRepository(database database.Database) repository.ReviewRepository {
	if gormDB, ok := database.(*gorm.DB); ok {
		return &SQLReviewRepository{db: gormDB}
	}
	panic("database must be GORM instance")
}

// Create 创建评价
func (r *SQLReviewRepository) Create(ctx context.Context, review *entity.Review) error {
	po := r.toReviewPO(review)
	return r.db.WithContext(ctx).Create(po).Error
}

// GetByID 根据ID获取评价
func (r *SQLReviewRepository) GetByID(ctx context.Context, id string) (*entity.Review, error) {
	var po ReviewPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// Update 更新评价
func (r *SQLReviewRepository) Update(ctx context.Context, review *entity.Review) error {
	po := r.toReviewPO(review)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除评价
func (r *SQLReviewRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&ReviewPO{}).Error
}

// ListByCourseID 根据课程ID查询评价列表
func (r *SQLReviewRepository) ListByCourseID(ctx context.Context, courseID string, status entity.ReviewStatus, limit, offset int) ([]*entity.Review, error) {
	var pos []ReviewPO
	query := r.db.WithContext(ctx).Where("course_id = ?", courseID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	reviews := make([]*entity.Review, len(pos))
	for i, po := range pos {
		reviews[i] = r.toDomainEntity(&po)
	}
	
	return reviews, nil
}

// ListByUserID 根据用户ID查询评价列表
func (r *SQLReviewRepository) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.Review, error) {
	var pos []ReviewPO
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	reviews := make([]*entity.Review, len(pos))
	for i, po := range pos {
		reviews[i] = r.toDomainEntity(&po)
	}
	
	return reviews, nil
}

// GetByUserAndCourse 根据用户和课程获取评价
func (r *SQLReviewRepository) GetByUserAndCourse(ctx context.Context, userID, courseID string) (*entity.Review, error) {
	var po ReviewPO
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

// ListPendingReviews 查询待审核评价列表
func (r *SQLReviewRepository) ListPendingReviews(ctx context.Context, limit, offset int) ([]*entity.Review, error) {
	var pos []ReviewPO
	err := r.db.WithContext(ctx).
		Where("status = ?", entity.ReviewStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	reviews := make([]*entity.Review, len(pos))
	for i, po := range pos {
		reviews[i] = r.toDomainEntity(&po)
	}
	
	return reviews, nil
}

// CountByCourseID 根据课程ID统计评价数量
func (r *SQLReviewRepository) CountByCourseID(ctx context.Context, courseID string, status entity.ReviewStatus) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&ReviewPO{}).Where("course_id = ?", courseID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Count(&count).Error
	return count, err
}

// CountByUserID 根据用户ID统计评价数量
func (r *SQLReviewRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ReviewPO{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetAverageRatingByCourseID 根据课程ID获取平均评分
func (r *SQLReviewRepository) GetAverageRatingByCourseID(ctx context.Context, courseID string) (float64, error) {
	var result struct {
		AvgRating float64
	}
	
	err := r.db.WithContext(ctx).Model(&ReviewPO{}).
		Select("AVG(rating) as avg_rating").
		Where("course_id = ? AND status = ?", courseID, entity.ReviewStatusApproved).
		Scan(&result).Error
	
	return result.AvgRating, err
}

// GetRatingDistributionByCourseID 根据课程ID获取评分分布
func (r *SQLReviewRepository) GetRatingDistributionByCourseID(ctx context.Context, courseID string) (map[int]int64, error) {
	var results []struct {
		Rating int   `gorm:"column:rating"`
		Count  int64 `gorm:"column:count"`
	}
	
	err := r.db.WithContext(ctx).Model(&ReviewPO{}).
		Select("rating, COUNT(*) as count").
		Where("course_id = ? AND status = ?", courseID, entity.ReviewStatusApproved).
		Group("rating").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	
	distribution := make(map[int]int64)
	for _, result := range results {
		distribution[result.Rating] = result.Count
	}
	
	return distribution, nil
}

// 辅助方法：领域实体转持久化对象
func (r *SQLReviewRepository) toReviewPO(review *entity.Review) *ReviewPO {
	po := &ReviewPO{
		ID:        review.ID,
		UserID:    review.UserID,
		CourseID:  review.CourseID,
		Rating:    review.Rating,
		Content:   review.Content,
		Status:    string(review.Status),
		CreatedAt: review.CreatedAt.Unix(),
		UpdatedAt: review.UpdatedAt.Unix(),
	}
	
	if review.ApprovedAt != nil {
		approvedAt := review.ApprovedAt.Unix()
		po.ApprovedAt = &approvedAt
	}
	
	return po
}

// 辅助方法：持久化对象转领域实体
func (r *SQLReviewRepository) toDomainEntity(po *ReviewPO) *entity.Review {
	review := &entity.Review{
		ID:        po.ID,
		UserID:    po.UserID,
		CourseID:  po.CourseID,
		Rating:    po.Rating,
		Content:   po.Content,
		Status:    entity.ReviewStatus(po.Status),
		CreatedAt: time.Unix(po.CreatedAt, 0),
		UpdatedAt: time.Unix(po.UpdatedAt, 0),
	}
	
	if po.ApprovedAt != nil {
		approvedAt := time.Unix(*po.ApprovedAt, 0)
		review.ApprovedAt = &approvedAt
	}
	
	return review
}
