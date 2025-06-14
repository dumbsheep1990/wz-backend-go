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

// TeacherPO 讲师持久化对象
type TeacherPO struct {
	ID          string `gorm:"primaryKey;column:id;size:36"`
	UserID      string `gorm:"index;column:user_id;size:36;not null"`
	Name        string `gorm:"column:name;size:100;not null"`
	Title       string `gorm:"column:title;size:100"`
	Bio         string `gorm:"column:bio;type:text"`
	Avatar      string `gorm:"column:avatar;size:255"`
	Email       string `gorm:"column:email;size:100"`
	Phone       string `gorm:"column:phone;size:20"`
	Expertise   string `gorm:"column:expertise;type:text"` // JSON字符串
	Experience  int    `gorm:"column:experience;default:0"`
	Rating      float64 `gorm:"column:rating;default:0"`
	TotalCourses int   `gorm:"column:total_courses;default:0"`
	TotalStudents int  `gorm:"column:total_students;default:0"`
	IsActive    bool   `gorm:"column:is_active;default:true"`
	CreatedAt   int64  `gorm:"column:created_at;not null"`
	UpdatedAt   int64  `gorm:"column:updated_at;not null"`
}

// TableName 表名
func (TeacherPO) TableName() string {
	return "teachers"
}

// SQLTeacherRepository 讲师仓储SQL实现
type SQLTeacherRepository struct {
	db *gorm.DB
}

// NewSQLTeacherRepository 创建讲师仓储实例
func NewSQLTeacherRepository(database database.Database) repository.TeacherRepository {
	// 假设database.Database接口有GetGormDB方法
	if gormDB, ok := database.(*gorm.DB); ok {
		return &SQLTeacherRepository{db: gormDB}
	}
	// 如果不是GORM，需要适配
	panic("database must be GORM instance")
}

// Create 创建讲师
func (r *SQLTeacherRepository) Create(ctx context.Context, teacher *entity.Teacher) error {
	po := r.toTeacherPO(teacher)
	return r.db.WithContext(ctx).Create(po).Error
}

// GetByID 根据ID获取讲师
func (r *SQLTeacherRepository) GetByID(ctx context.Context, id string) (*entity.Teacher, error) {
	var po TeacherPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// GetByUserID 根据用户ID获取讲师
func (r *SQLTeacherRepository) GetByUserID(ctx context.Context, userID string) (*entity.Teacher, error) {
	var po TeacherPO
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEntity(&po), nil
}

// Update 更新讲师
func (r *SQLTeacherRepository) Update(ctx context.Context, teacher *entity.Teacher) error {
	po := r.toTeacherPO(teacher)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除讲师
func (r *SQLTeacherRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&TeacherPO{}).Error
}

// List 查询讲师列表
func (r *SQLTeacherRepository) List(ctx context.Context, params repository.TeacherQueryParams) ([]*entity.Teacher, int64, error) {
	query := r.db.WithContext(ctx).Model(&TeacherPO{})
	
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
	
	var pos []TeacherPO
	if err := query.Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	
	teachers := make([]*entity.Teacher, len(pos))
	for i, po := range pos {
		teachers[i] = r.toDomainEntity(&po)
	}
	
	return teachers, total, nil
}

// Search 搜索讲师
func (r *SQLTeacherRepository) Search(ctx context.Context, keyword string, params repository.TeacherQueryParams) ([]*entity.Teacher, int64, error) {
	query := r.db.WithContext(ctx).Model(&TeacherPO{})
	
	// 搜索条件
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR title LIKE ? OR bio LIKE ?", searchPattern, searchPattern, searchPattern)
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
	
	var pos []TeacherPO
	if err := query.Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	
	teachers := make([]*entity.Teacher, len(pos))
	for i, po := range pos {
		teachers[i] = r.toDomainEntity(&po)
	}
	
	return teachers, total, nil
}

// CountAll 统计所有讲师数量
func (r *SQLTeacherRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&TeacherPO{}).Count(&count).Error
	return count, err
}

// CountActive 统计活跃讲师数量
func (r *SQLTeacherRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&TeacherPO{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

// 辅助方法：构建WHERE子句
func (r *SQLTeacherRepository) buildWhereClause(query *gorm.DB, params repository.TeacherQueryParams) *gorm.DB {
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}
	
	if params.MinRating != nil {
		query = query.Where("rating >= ?", *params.MinRating)
	}
	
	if params.MinExperience != nil {
		query = query.Where("experience >= ?", *params.MinExperience)
	}
	
	if params.HasCourses != nil && *params.HasCourses {
		query = query.Where("total_courses > 0")
	}
	
	return query
}

// 辅助方法：构建ORDER BY子句
func (r *SQLTeacherRepository) buildOrderClause(query *gorm.DB, params repository.TeacherQueryParams) *gorm.DB {
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
func (r *SQLTeacherRepository) buildLimitClause(query *gorm.DB, params repository.TeacherQueryParams) *gorm.DB {
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
func (r *SQLTeacherRepository) toTeacherPO(teacher *entity.Teacher) *TeacherPO {
	// 将expertise数组转为JSON字符串
	expertiseStr := strings.Join(teacher.Expertise, ",")
	
	return &TeacherPO{
		ID:            teacher.ID,
		UserID:        teacher.UserID,
		Name:          teacher.Name,
		Title:         teacher.Title,
		Bio:           teacher.Bio,
		Avatar:        teacher.Avatar,
		Email:         teacher.Email,
		Phone:         teacher.Phone,
		Expertise:     expertiseStr,
		Experience:    teacher.Experience,
		Rating:        teacher.Rating,
		TotalCourses:  teacher.TotalCourses,
		TotalStudents: teacher.TotalStudents,
		IsActive:      teacher.IsActive,
		CreatedAt:     teacher.CreatedAt.Unix(),
		UpdatedAt:     teacher.UpdatedAt.Unix(),
	}
}

// 辅助方法：持久化对象转领域实体
func (r *SQLTeacherRepository) toDomainEntity(po *TeacherPO) *entity.Teacher {
	// 将expertise字符串转为数组
	var expertise []string
	if po.Expertise != "" {
		expertise = strings.Split(po.Expertise, ",")
	}
	
	return &entity.Teacher{
		ID:            po.ID,
		UserID:        po.UserID,
		Name:          po.Name,
		Title:         po.Title,
		Bio:           po.Bio,
		Avatar:        po.Avatar,
		Email:         po.Email,
		Phone:         po.Phone,
		Expertise:     expertise,
		Experience:    po.Experience,
		Rating:        po.Rating,
		TotalCourses:  po.TotalCourses,
		TotalStudents: po.TotalStudents,
		IsActive:      po.IsActive,
		CreatedAt:     time.Unix(po.CreatedAt, 0),
		UpdatedAt:     time.Unix(po.UpdatedAt, 0),
	}
}
