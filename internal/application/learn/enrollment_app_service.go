package learn

import (
	"context"
	"time"

	"wz-backend-go/internal/domain/learn/dto"
	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/domain/learn/service"
)

// EnrollmentAppService 课程报名应用服务，处理报名相关的应用场景
type EnrollmentAppService struct {
	enrollmentService *service.EnrollmentService
	courseService     *service.CourseService
	teacherService    *service.TeacherService
}

// NewEnrollmentAppService 创建报名应用服务
func NewEnrollmentAppService(
	enrollmentService *service.EnrollmentService,
	courseService *service.CourseService,
	teacherService *service.TeacherService,
) *EnrollmentAppService {
	return &EnrollmentAppService{
		enrollmentService: enrollmentService,
		courseService:     courseService,
		teacherService:    teacherService,
	}
}

// EnrollCourse 用户报名课程
func (s *EnrollmentAppService) EnrollCourse(ctx context.Context, courseID, userID, orderID string, expiresAt *time.Time) (*dto.EnrollmentDTO, error) {
	// 报名课程
	enrollment, err := s.enrollmentService.EnrollCourse(ctx, courseID, userID, orderID, expiresAt)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(ctx, enrollment)
}

// UpdateProgress 更新学习进度
func (s *EnrollmentAppService) UpdateProgress(ctx context.Context, id string, completedCount int) (*dto.EnrollmentDTO, error) {
	// 更新进度
	enrollment, err := s.enrollmentService.UpdateProgress(ctx, id, completedCount)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(ctx, enrollment)
}

// CompleteCourse 标记课程为已完成
func (s *EnrollmentAppService) CompleteCourse(ctx context.Context, id string) (*dto.EnrollmentDTO, error) {
	// 标记为已完成
	enrollment, err := s.enrollmentService.CompleteCourse(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(ctx, enrollment)
}

// AddRating 添加课程评分
func (s *EnrollmentAppService) AddRating(ctx context.Context, id string, score float64, comment string) (*dto.EnrollmentDTO, error) {
	// 添加评分
	enrollment, err := s.enrollmentService.AddRating(ctx, id, score, comment)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(ctx, enrollment)
}

// GetEnrollmentByID 获取报名记录详情
func (s *EnrollmentAppService) GetEnrollmentByID(ctx context.Context, id string) (*dto.EnrollmentDTO, error) {
	// 获取报名记录
	enrollment, err := s.enrollmentService.GetEnrollmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(ctx, enrollment)
}

// ListUserEnrollments 获取用户的所有报名记录
func (s *EnrollmentAppService) ListUserEnrollments(ctx context.Context, userID string, page, pageSize int) ([]*dto.EnrollmentDTO, int64, error) {
	// 查询参数
	params := repository.EnrollmentQueryParams{
		UserID: userID,
	}

	// 获取用户的报名记录
	enrollments, total, err := s.enrollmentService.ListEnrollments(ctx, params, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为 DTO 列表
	dtos := make([]*dto.EnrollmentDTO, 0, len(enrollments))
	for _, enrollment := range enrollments {
		dto, err := s.convertToDTO(ctx, enrollment)
		if err != nil {
			continue
		}
		dtos = append(dtos, dto)
	}

	return dtos, total, nil
}

// ListCourseEnrollments 获取课程的所有报名记录
func (s *EnrollmentAppService) ListCourseEnrollments(ctx context.Context, courseID string, page, pageSize int) ([]*dto.EnrollmentDTO, int64, error) {
	// 查询参数
	params := repository.EnrollmentQueryParams{
		CourseID: courseID,
	}

	// 获取课程的报名记录
	enrollments, total, err := s.enrollmentService.ListEnrollments(ctx, params, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为 DTO 列表
	dtos := make([]*dto.EnrollmentDTO, 0, len(enrollments))
	for _, enrollment := range enrollments {
		dto, err := s.convertToDTO(ctx, enrollment)
		if err != nil {
			continue
		}
		dtos = append(dtos, dto)
	}

	return dtos, total, nil
}

// RefundEnrollment 退款处理
func (s *EnrollmentAppService) RefundEnrollment(ctx context.Context, id string) (*dto.EnrollmentDTO, error) {
	// 标记为已退款
	enrollment, err := s.enrollmentService.RefundEnrollment(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDTO(ctx, enrollment)
}

// ProcessExpiredEnrollments 处理过期的报名记录
func (s *EnrollmentAppService) ProcessExpiredEnrollments(ctx context.Context) (int, error) {
	// 处理过期记录
	return s.enrollmentService.ProcessExpiredEnrollments(ctx)
}

// GetEnrollmentStats 获取报名统计信息
func (s *EnrollmentAppService) GetEnrollmentStats(ctx context.Context) (*dto.EnrollmentStats, error) {
	// 获取报名统计
	totalCount, activeCount, completedCount, recentCount, err := s.enrollmentService.GetEnrollmentStats(ctx)
	if err != nil {
		return nil, err
	}

	// 构建统计 DTO
	stats := &dto.EnrollmentStats{
		TotalCount:     int(totalCount),
		ActiveCount:    int(activeCount),
		CompletedCount: int(completedCount),
		RecentCount:    int(recentCount),
	}

	return stats, nil
}

// 辅助函数：转换报名实体到 DTO
func (s *EnrollmentAppService) convertToDTO(ctx context.Context, enrollment *entity.Enrollment) (*dto.EnrollmentDTO, error) {
	// 获取课程信息
	course, err := s.courseService.GetCourseByID(ctx, enrollment.CourseID)
	if err != nil {
		return nil, err
	}

	// 构建 DTO
	dto := &dto.EnrollmentDTO{
		ID:               enrollment.ID,
		UserID:           enrollment.UserID,
		CourseID:         enrollment.CourseID,
		CourseTitle:      course.Title,
		CourseCoverImage: course.CoverImage,
		Status:           string(enrollment.Status),
		Progress:         enrollment.Progress,
		CompletedCount:   enrollment.CompletedCount,
		TotalCount:       enrollment.TotalCount,
		Rating:           enrollment.Rating,
		Comment:          enrollment.Comment,
		EnrolledAt:       enrollment.CreatedAt,
		CompletedAt:      enrollment.CompletedAt,
		ExpiredAt:        enrollment.ExpiredAt,
	}

	// 获取教师信息
	if teacher, err := s.teacherService.GetTeacherByID(ctx, course.TeacherID); err == nil {
		dto.TeacherID = teacher.ID
		dto.TeacherName = teacher.Name
	}

	return dto, nil
}
