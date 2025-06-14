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

// CoursePO 课程持久化对象
type CoursePO struct {
	ID               string  `gorm:"primaryKey;column:id;size:36"`
	Title            string  `gorm:"column:title;size:200;not null"`
	Subtitle         string  `gorm:"column:subtitle;size:500"`
	Description      string  `gorm:"column:description;type:text"`
	TeacherID        string  `gorm:"index;column:teacher_id;size:36;not null"`
	Cover            string  `gorm:"column:cover;size:255"`
	Level            string  `gorm:"column:level;size:20;not null"`
	Duration         int     `gorm:"column:duration;default:0"`
	Price            float64 `gorm:"column:price;default:0"`
	DiscountPrice    float64 `gorm:"column:discount_price;default:0"`
	Status           string  `gorm:"column:status;size:20;not null"`
	CategoryIDs      string  `gorm:"column:category_ids;type:text"` // JSON字符串
	Tags             string  `gorm:"column:tags;type:text"`         // JSON字符串
	ChaptersCount    int     `gorm:"column:chapters_count;default:0"`
	LessonsCount     int     `gorm:"column:lessons_count;default:0"`
	EnrollmentsCount int     `gorm:"column:enrollments_count;default:0"`
	Rating           float64 `gorm:"column:rating;default:0"`
	RatingCount      int     `gorm:"column:rating_count;default:0"`
	CreatedAt        int64   `gorm:"column:created_at;not null"`
	UpdatedAt        int64   `gorm:"column:updated_at;not null"`
	PublishedAt      *int64  `gorm:"column:published_at"`
}

// TableName 表名
func (CoursePO) TableName() string {
	return "courses"
}

// SQLCourseRepository 课程仓储SQL实现
type SQLCourseRepository struct {
	db *gorm.DB
}

// NewSQLCourseRepository 创建课程仓储实例
func NewSQLCourseRepository(database database.Database) repository.CourseRepository {
	if gormDB, ok := database.(*gorm.DB); ok {
		return &SQLCourseRepository{db: gormDB}
	}
	panic("database must be GORM instance")
}

// Create 创建课程
func (r *SQLCourseRepository) Create(ctx context.Context, course *entity.Course) error {
	po := r.toCoursePO(course)
	return r.db.WithContext(ctx).Create(po).Error
}

// GetByID 根据ID获取课程
func (r *SQLCourseRepository) GetByID(ctx context.Context, id string) (*entity.Course, error) {
	var po CoursePO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// Update 更新课程
func (r *SQLCourseRepository) Update(ctx context.Context, course *entity.Course) error {
	po := r.toCoursePO(course)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除课程
func (r *SQLCourseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&CoursePO{}).Error
}

// List 查询课程列表
func (r *SQLCourseRepository) List(ctx context.Context, params repository.CourseQueryParams) ([]*entity.Course, int64, error) {
	query := r.db.WithContext(ctx).Model(&CoursePO{})
	
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
	
	var pos []CoursePO
	if err := query.Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	
	courses := make([]*entity.Course, len(pos))
	for i, po := range pos {
		courses[i] = r.toDomainEntity(&po)
	}
	
	return courses, total, nil
}

// ListByTeacherID 根据讲师ID查询课程列表
func (r *SQLCourseRepository) ListByTeacherID(ctx context.Context, teacherID string, params repository.CourseQueryParams) ([]*entity.Course, int64, error) {
	query := r.db.WithContext(ctx).Model(&CoursePO{}).Where("teacher_id = ?", teacherID)
	
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
	
	var pos []CoursePO
	if err := query.Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	
	courses := make([]*entity.Course, len(pos))
	for i, po := range pos {
		courses[i] = r.toDomainEntity(&po)
	}
	
	return courses, total, nil
}

// ListByIDs 根据ID列表查询课程
func (r *SQLCourseRepository) ListByIDs(ctx context.Context, ids []string) ([]*entity.Course, error) {
	var pos []CoursePO
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	courses := make([]*entity.Course, len(pos))
	for i, po := range pos {
		courses[i] = r.toDomainEntity(&po)
	}
	
	return courses, nil
}

// ListByCategoryID 根据分类ID查询课程列表
func (r *SQLCourseRepository) ListByCategoryID(ctx context.Context, categoryID string, params repository.CourseQueryParams) ([]*entity.Course, int64, error) {
	query := r.db.WithContext(ctx).Model(&CoursePO{}).Where("category_ids LIKE ?", "%"+categoryID+"%")
	
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
	
	var pos []CoursePO
	if err := query.Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	
	courses := make([]*entity.Course, len(pos))
	for i, po := range pos {
		courses[i] = r.toDomainEntity(&po)
	}
	
	return courses, total, nil
}

// ListPopular 查询热门课程
func (r *SQLCourseRepository) ListPopular(ctx context.Context, limit int) ([]*entity.Course, error) {
	var pos []CoursePO
	err := r.db.WithContext(ctx).
		Where("status = ?", entity.CourseStatusPublished).
		Order("enrollments_count DESC, rating DESC").
		Limit(limit).
		Find(&pos).Error
	
	if err != nil {
		return nil, err
	}
	
	courses := make([]*entity.Course, len(pos))
	for i, po := range pos {
		courses[i] = r.toDomainEntity(&po)
	}
	
	return courses, nil
}

// ListRecent 查询最新课程
func (r *SQLCourseRepository) ListRecent(ctx context.Context, limit int) ([]*entity.Course, error) {
	var pos []CoursePO
	err := r.db.WithContext(ctx).
		Where("status = ?", entity.CourseStatusPublished).
		Order("published_at DESC").
		Limit(limit).
		Find(&pos).Error
	
	if err != nil {
		return nil, err
	}
	
	courses := make([]*entity.Course, len(pos))
	for i, po := range pos {
		courses[i] = r.toDomainEntity(&po)
	}
	
	return courses, nil
}

// CountAll 统计所有课程数量
func (r *SQLCourseRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&CoursePO{}).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态统计课程数量
func (r *SQLCourseRepository) CountByStatus(ctx context.Context, status entity.CourseStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&CoursePO{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// CountByTeacherID 根据讲师ID统计课程数量
func (r *SQLCourseRepository) CountByTeacherID(ctx context.Context, teacherID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&CoursePO{}).Where("teacher_id = ?", teacherID).Count(&count).Error
	return count, err
}

// CountByCategoryID 根据分类ID统计课程数量
func (r *SQLCourseRepository) CountByCategoryID(ctx context.Context, categoryID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&CoursePO{}).Where("category_ids LIKE ?", "%"+categoryID+"%").Count(&count).Error
	return count, err
}

// Search 搜索课程
func (r *SQLCourseRepository) Search(ctx context.Context, keyword string, params repository.CourseQueryParams) ([]*entity.Course, int64, error) {
	query := r.db.WithContext(ctx).Model(&CoursePO{})
	
	// 搜索条件
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR subtitle LIKE ? OR description LIKE ? OR tags LIKE ?", 
			searchPattern, searchPattern, searchPattern, searchPattern)
	}
	
	// 其他查询条件
	query = r.buildWhereClause(query, params)
	
	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 构建排序和分页
	query = r.buildOrderClause(query, params)
	query = r.buildLimitClause(query, params)
	
	var pos []CoursePO
	if err := query.Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	
	courses := make([]*entity.Course, len(pos))
	for i, po := range pos {
		courses[i] = r.toDomainEntity(&po)
	}
	
	return courses, total, nil
}

// 辅助方法：构建WHERE子句
func (r *SQLCourseRepository) buildWhereClause(query *gorm.DB, params repository.CourseQueryParams) *gorm.DB {
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	
	if params.Level != nil {
		query = query.Where("level = ?", *params.Level)
	}
	
	if params.PriceRange[0] > 0 || params.PriceRange[1] > 0 {
		if params.PriceRange[1] > 0 {
			query = query.Where("price BETWEEN ? AND ?", params.PriceRange[0], params.PriceRange[1])
		} else {
			query = query.Where("price >= ?", params.PriceRange[0])
		}
	}
	
	if params.FreeCourses != nil && *params.FreeCourses {
		query = query.Where("price = 0")
	}
	
	if len(params.Tags) > 0 {
		for _, tag := range params.Tags {
			query = query.Where("tags LIKE ?", "%"+tag+"%")
		}
	}
	
	return query
}

// 辅助方法：构建ORDER BY子句
func (r *SQLCourseRepository) buildOrderClause(query *gorm.DB, params repository.CourseQueryParams) *gorm.DB {
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
func (r *SQLCourseRepository) buildLimitClause(query *gorm.DB, params repository.CourseQueryParams) *gorm.DB {
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
func (r *SQLCourseRepository) toCoursePO(course *entity.Course) *CoursePO {
	// 将分类ID数组转为字符串
	categoryIDsStr := strings.Join(course.CategoryIDs, ",")
	
	// 将标签数组转为字符串
	tagsStr := strings.Join(course.Tags, ",")
	
	po := &CoursePO{
		ID:               course.ID,
		Title:            course.Title,
		Subtitle:         course.Subtitle,
		Description:      course.Description,
		TeacherID:        course.TeacherID,
		Cover:            course.Cover,
		Level:            string(course.Level),
		Duration:         course.Duration,
		Price:            course.Price,
		DiscountPrice:    course.DiscountPrice,
		Status:           string(course.Status),
		CategoryIDs:      categoryIDsStr,
		Tags:             tagsStr,
		ChaptersCount:    course.ChaptersCount,
		LessonsCount:     course.LessonsCount,
		EnrollmentsCount: course.EnrollmentsCount,
		Rating:           course.Rating,
		RatingCount:      course.RatingCount,
		CreatedAt:        course.CreatedAt.Unix(),
		UpdatedAt:        course.UpdatedAt.Unix(),
	}
	
	if course.PublishedAt != nil {
		publishedAt := course.PublishedAt.Unix()
		po.PublishedAt = &publishedAt
	}
	
	return po
}

// 辅助方法：持久化对象转领域实体
func (r *SQLCourseRepository) toDomainEntity(po *CoursePO) *entity.Course {
	// 将分类ID字符串转为数组
	var categoryIDs []string
	if po.CategoryIDs != "" {
		categoryIDs = strings.Split(po.CategoryIDs, ",")
	}
	
	// 将标签字符串转为数组
	var tags []string
	if po.Tags != "" {
		tags = strings.Split(po.Tags, ",")
	}
	
	course := &entity.Course{
		ID:               po.ID,
		Title:            po.Title,
		Subtitle:         po.Subtitle,
		Description:      po.Description,
		TeacherID:        po.TeacherID,
		Cover:            po.Cover,
		Level:            entity.CourseLevel(po.Level),
		Duration:         po.Duration,
		Price:            po.Price,
		DiscountPrice:    po.DiscountPrice,
		Status:           entity.CourseStatus(po.Status),
		CategoryIDs:      categoryIDs,
		Tags:             tags,
		ChaptersCount:    po.ChaptersCount,
		LessonsCount:     po.LessonsCount,
		EnrollmentsCount: po.EnrollmentsCount,
		Rating:           po.Rating,
		RatingCount:      po.RatingCount,
		CreatedAt:        time.Unix(po.CreatedAt, 0),
		UpdatedAt:        time.Unix(po.UpdatedAt, 0),
	}
	
	if po.PublishedAt != nil {
		publishedAt := time.Unix(*po.PublishedAt, 0)
		course.PublishedAt = &publishedAt
	}
	
	return course
}
