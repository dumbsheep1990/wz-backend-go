package service

import (
	"context"
	"errors"
	"time"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
)

var (
	ErrEnrollmentExists    = errors.New("已经报名该课程")
	ErrEnrollmentNotFound  = errors.New("报名记录不存在")
	ErrInvalidRating       = errors.New("无效的评分")
	ErrInvalidEnrollmentStatus = errors.New("无效的报名状态")
)

// EnrollmentService 课程报名领域服务
type EnrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	courseRepo     repository.CourseRepository
	teacherRepo    repository.TeacherRepository
	lessonRepo     repository.LessonRepository
}

// NewEnrollmentService 创建报名服务
func NewEnrollmentService(
	enrollmentRepo repository.EnrollmentRepository,
	courseRepo repository.CourseRepository,
	teacherRepo repository.TeacherRepository,
	lessonRepo repository.LessonRepository,
) *EnrollmentService {
	return &EnrollmentService{
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
		teacherRepo:    teacherRepo,
		lessonRepo:     lessonRepo,
	}
}

// EnrollCourse 用户报名课程
func (s *EnrollmentService) EnrollCourse(ctx context.Context, courseID, userID, orderID string, expiresAt *time.Time) (*entity.Enrollment, error) {
	// 检查是否已经报名
	existingEnrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err == nil && existingEnrollment != nil {
		return nil, ErrEnrollmentExists
	}

	// 获取课程
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	// 获取课时总数
	lessonsCount, err := s.lessonRepo.CountPublishedByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// 创建报名记录
	enrollment := entity.NewEnrollment(courseID, userID, orderID, int(lessonsCount), expiresAt)

	// 保存报名记录
	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	// 更新课程报名人数
	course.AddEnrollment()
	if err := s.courseRepo.Update(ctx, course); err != nil {
		return nil, err
	}

	// 更新讲师的学生数
	teacher, err := s.teacherRepo.GetByID(ctx, course.TeacherID)
	if err == nil {
		teacher.AddStudentsCount(1)
		if err := s.teacherRepo.Update(ctx, teacher); err != nil {
			return nil, err
		}
	}

	return enrollment, nil
}

// UpdateProgress 更新学习进度
func (s *EnrollmentService) UpdateProgress(ctx context.Context, id string, completedCount int) (*entity.Enrollment, error) {
	// 获取报名记录
	enrollment, err := s.enrollmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrEnrollmentNotFound
	}

	// 更新进度
	enrollment.UpdateProgress(completedCount)

	// 保存报名记录
	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

// CompleteCourse 标记课程为已完成
func (s *EnrollmentService) CompleteCourse(ctx context.Context, id string) (*entity.Enrollment, error) {
	// 获取报名记录
	enrollment, err := s.enrollmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrEnrollmentNotFound
	}

	// 标记为已完成
	enrollment.Complete()

	// 保存报名记录
	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

// AddRating 添加课程评分
func (s *EnrollmentService) AddRating(ctx context.Context, id string, score float64, comment string) (*entity.Enrollment, error) {
	// 验证评分
	if score < 0 || score > 5 {
		return nil, ErrInvalidRating
	}

	// 获取报名记录
	enrollment, err := s.enrollmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrEnrollmentNotFound
	}

	// 添加评分
	enrollment.AddRating(score, comment)

	// 保存报名记录
	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return nil, err
	}

	// 更新课程的评分
	course, err := s.courseRepo.GetByID(ctx, enrollment.CourseID)
	if err == nil {
		course.AddRating(score)
		if err := s.courseRepo.Update(ctx, course); err != nil {
			return nil, err
		}
	}

	// 更新讲师的评分
	if course != nil {
		teacher, err := s.teacherRepo.GetByID(ctx, course.TeacherID)
		if err == nil {
			teacher.AddRating(score)
			if err := s.teacherRepo.Update(ctx, teacher); err != nil {
				return nil, err
			}
		}
	}

	return enrollment, nil
}

// ProcessExpiredEnrollments 处理过期的报名记录
func (s *EnrollmentService) ProcessExpiredEnrollments(ctx context.Context) (int, error) {
	now := time.Now()
	
	// 获取已过期但未标记为过期的报名记录
	expiredEnrollments, err := s.enrollmentRepo.ListExpired(ctx, now)
	if err != nil {
		return 0, err
	}

	// 处理每条过期记录
	count := 0
	for _, enrollment := range expiredEnrollments {
		if enrollment.Status != entity.EnrollmentStatusExpired {
			enrollment.Expire()
			if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
				continue
			}
			count++
		}
	}

	return count, nil
}

// RefundEnrollment 退款处理
func (s *EnrollmentService) RefundEnrollment(ctx context.Context, id string) (*entity.Enrollment, error) {
	// 获取报名记录
	enrollment, err := s.enrollmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrEnrollmentNotFound
	}

	// 标记为已退款
	enrollment.Refund()

	// 保存报名记录
	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

// GetEnrollmentStats 获取报名统计数据
func (s *EnrollmentService) GetEnrollmentStats(ctx context.Context) (totalCount, activeCount, completedCount, recentCount int64, err error) {
	totalCount, err = s.enrollmentRepo.CountActiveEnrollments(ctx)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	activeCount, err = s.enrollmentRepo.CountByStatus(ctx, entity.EnrollmentStatusActive)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	completedCount, err = s.enrollmentRepo.CountByStatus(ctx, entity.EnrollmentStatusCompleted)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// 获取最近30天的报名数
	thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour)
	recentCount, err = s.enrollmentRepo.CountRecentEnrollments(ctx, thirtyDaysAgo)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return totalCount, activeCount, completedCount, recentCount, nil
}
