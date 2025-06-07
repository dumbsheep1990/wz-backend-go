package learn

import (
	"context"

	"wz-backend-go/internal/domain/learn/dto"
	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/domain/learn/service"
)

// TeacherAppService 讲师应用服务，处理讲师相关的应用场景
type TeacherAppService struct {
	teacherService *service.TeacherService
	courseService  *service.CourseService
}

// NewTeacherAppService 创建讲师应用服务
func NewTeacherAppService(
	teacherService *service.TeacherService,
	courseService *service.CourseService,
) *TeacherAppService {
	return &TeacherAppService{
		teacherService: teacherService,
		courseService:  courseService,
	}
}

// CreateTeacher 创建新讲师
func (s *TeacherAppService) CreateTeacher(ctx context.Context, userID, name, title, introduction, avatar string) (*dto.TeacherDetailDTO, error) {
	// 创建讲师
	teacher, err := s.teacherService.CreateTeacher(ctx, userID, name, title, introduction, avatar)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, teacher)
}

// UpdateTeacherProfile 更新讲师基本资料
func (s *TeacherAppService) UpdateTeacherProfile(ctx context.Context, id, name, avatar, title, introduction string) (*dto.TeacherDetailDTO, error) {
	// 更新讲师信息
	teacher, err := s.teacherService.UpdateTeacherProfile(ctx, id, name, avatar, title, introduction)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, teacher)
}

// UpdateTeacherContact 更新讲师联系信息
func (s *TeacherAppService) UpdateTeacherContact(ctx context.Context, id, email, phone string) (*dto.TeacherDetailDTO, error) {
	// 更新联系信息
	teacher, err := s.teacherService.UpdateTeacherContact(ctx, id, email, phone)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, teacher)
}

// SetTeacherSpecialties 设置讲师专长领域
func (s *TeacherAppService) SetTeacherSpecialties(ctx context.Context, id string, specialties []string) (*dto.TeacherDetailDTO, error) {
	// 更新讲师专长
	teacher, err := s.teacherService.SetTeacherSpecialties(ctx, id, specialties)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, teacher)
}

// SetTeacherSocialProfiles 设置讲师社交档案
func (s *TeacherAppService) SetTeacherSocialProfiles(ctx context.Context, id string, profiles []string) (*dto.TeacherDetailDTO, error) {
	// 更新社交档案
	teacher, err := s.teacherService.SetTeacherSocialProfiles(ctx, id, profiles)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, teacher)
}

// ActivateTeacher 激活讲师
func (s *TeacherAppService) ActivateTeacher(ctx context.Context, id string) (*dto.TeacherDetailDTO, error) {
	// 激活讲师
	teacher, err := s.teacherService.ActivateTeacher(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, teacher)
}

// DeactivateTeacher 停用讲师
func (s *TeacherAppService) DeactivateTeacher(ctx context.Context, id string) (*dto.TeacherDetailDTO, error) {
	// 停用讲师
	teacher, err := s.teacherService.DeactivateTeacher(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, teacher)
}

// GetTeacherDetail 获取讲师详情
func (s *TeacherAppService) GetTeacherDetail(ctx context.Context, id string) (*dto.TeacherDetailDTO, error) {
	// 获取讲师
	teacher, err := s.teacherService.GetTeacherByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, teacher)
}

// GetTeacherByUserID 通过用户ID获取讲师
func (s *TeacherAppService) GetTeacherByUserID(ctx context.Context, userID string) (*dto.TeacherDetailDTO, error) {
	// 获取讲师
	teacher, err := s.teacherService.GetTeacherByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, teacher)
}

// ListTeachers 获取讲师列表
func (s *TeacherAppService) ListTeachers(ctx context.Context, params repository.TeacherQueryParams, page, pageSize int) ([]*dto.TeacherBasicDTO, int64, error) {
	// 获取讲师列表
	teachers, total, err := s.teacherService.ListTeachers(ctx, params, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为 DTO 列表
	dtos := make([]*dto.TeacherBasicDTO, 0, len(teachers))
	for _, teacher := range teachers {
		basicDTO, err := s.convertToBasicDTO(ctx, teacher)
		if err != nil {
			continue
		}
		dtos = append(dtos, basicDTO)
	}

	return dtos, total, nil
}

// DeleteTeacher 删除讲师
func (s *TeacherAppService) DeleteTeacher(ctx context.Context, id string) error {
	return s.teacherService.DeleteTeacher(ctx, id)
}

// GetTeacherWithCourses 获取讲师及其课程
func (s *TeacherAppService) GetTeacherWithCourses(ctx context.Context, id string) (*dto.TeacherDetailDTO, []*dto.CourseBasicDTO, error) {
	// 获取讲师详情
	teacherDTO, err := s.GetTeacherDetail(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	// 查询讲师的课程
	params := repository.CourseQueryParams{
		TeacherID: id,
	}
	courses, _, err := s.courseService.ListCourses(ctx, params, 1, 100)
	if err != nil {
		return teacherDTO, nil, err
	}

	// 转换课程为 DTO 列表
	courseDTOs := make([]*dto.CourseBasicDTO, 0, len(courses))
	for _, course := range courses {
		courseDTO := &dto.CourseBasicDTO{
			ID:          course.ID,
			Title:       course.Title,
			Subtitle:    course.Subtitle,
			CoverImage:  course.CoverImage,
			TeacherID:   course.TeacherID,
			TeacherName: teacherDTO.Name,
			Level:       string(course.Level),
			Status:      string(course.Status),
			EnrollCount: course.EnrollCount,
			LessonCount: course.LessonCount,
			Duration:    course.Duration,
			Rating:      course.Rating,
			CreatedAt:   course.CreatedAt,
		}
		courseDTOs = append(courseDTOs, courseDTO)
	}

	return teacherDTO, courseDTOs, nil
}

// GetTeacherStats 获取讲师统计信息
func (s *TeacherAppService) GetTeacherStats(ctx context.Context) (*dto.TeacherStats, error) {
	// 获取讲师统计
	totalCount, activeCount, inactiveCount, err := s.teacherService.GetTeacherStats(ctx)
	if err != nil {
		return nil, err
	}

	// 构建统计 DTO
	stats := &dto.TeacherStats{
		TotalCount:   int(totalCount),
		ActiveCount:  int(activeCount),
		InactiveCount: int(inactiveCount),
	}

	return stats, nil
}

// 辅助函数：转换讲师实体到基本 DTO
func (s *TeacherAppService) convertToBasicDTO(ctx context.Context, teacher *entity.Teacher) (*dto.TeacherBasicDTO, error) {
	return &dto.TeacherBasicDTO{
		ID:           teacher.ID,
		Name:         teacher.Name,
		Title:        teacher.Title,
		Avatar:       teacher.Avatar,
		Status:       string(teacher.Status),
		Rating:       teacher.Rating,
		CoursesCount: teacher.CoursesCount,
		StudentsCount: teacher.StudentsCount,
	}, nil
}

// 辅助函数：转换讲师实体到详细 DTO
func (s *TeacherAppService) convertToDetailDTO(ctx context.Context, teacher *entity.Teacher) (*dto.TeacherDetailDTO, error) {
	return &dto.TeacherDetailDTO{
		ID:             teacher.ID,
		UserID:         teacher.UserID,
		Name:           teacher.Name,
		Title:          teacher.Title,
		Avatar:         teacher.Avatar,
		Introduction:   teacher.Introduction,
		Email:          teacher.Email,
		Phone:          teacher.Phone,
		Specialties:    teacher.Specialties,
		SocialProfiles: teacher.SocialProfiles,
		Status:         string(teacher.Status),
		Rating:         teacher.Rating,
		CoursesCount:   teacher.CoursesCount,
		StudentsCount:  teacher.StudentsCount,
		CreatedAt:      teacher.CreatedAt,
		UpdatedAt:      teacher.UpdatedAt,
	}, nil
}
