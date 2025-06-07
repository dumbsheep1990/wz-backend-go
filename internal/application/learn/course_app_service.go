package learn

import (
	"context"

	"wz-backend-go/internal/domain/learn/dto"
	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/domain/learn/service"
)

// CourseAppService 课程应用服务，处理课程相关的应用场景
type CourseAppService struct {
	courseService      *service.CourseService
	categoryService    *service.CategoryService
	chapterLessonService *service.ChapterLessonService
	teacherService     *service.TeacherService
	enrollmentService  *service.EnrollmentService
}

// NewCourseAppService 创建课程应用服务
func NewCourseAppService(
	courseService *service.CourseService,
	categoryService *service.CategoryService,
	chapterLessonService *service.ChapterLessonService,
	teacherService *service.TeacherService,
	enrollmentService *service.EnrollmentService,
) *CourseAppService {
	return &CourseAppService{
		courseService:      courseService,
		categoryService:    categoryService,
		chapterLessonService: chapterLessonService,
		teacherService:     teacherService,
		enrollmentService:  enrollmentService,
	}
}

// CreateCourse 创建课程
func (s *CourseAppService) CreateCourse(ctx context.Context, teacherID, title, subtitle, description string, level entity.CourseLevel) (*dto.CourseDetailDTO, error) {
	// 创建课程
	course, err := s.courseService.CreateCourse(ctx, teacherID, title, subtitle, description, level)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, course)
}

// UpdateCourse 更新课程基本信息
func (s *CourseAppService) UpdateCourse(ctx context.Context, id, title, subtitle, description string, level entity.CourseLevel) (*dto.CourseDetailDTO, error) {
	// 更新课程
	course, err := s.courseService.UpdateCourse(ctx, id, title, subtitle, description, level)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, course)
}

// UpdateCourseMedia 更新课程媒体信息
func (s *CourseAppService) UpdateCourseMedia(ctx context.Context, id, coverImage, previewVideo string) (*dto.CourseDetailDTO, error) {
	// 更新课程媒体
	course, err := s.courseService.UpdateCourseMedia(ctx, id, coverImage, previewVideo)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, course)
}

// SetCourseCategories 设置课程分类
func (s *CourseAppService) SetCourseCategories(ctx context.Context, id string, categoryIDs []string) (*dto.CourseDetailDTO, error) {
	// 设置课程分类
	course, err := s.courseService.SetCourseCategories(ctx, id, categoryIDs)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, course)
}

// SetCourseTags 设置课程标签
func (s *CourseAppService) SetCourseTags(ctx context.Context, id string, tags []string) (*dto.CourseDetailDTO, error) {
	// 设置课程标签
	course, err := s.courseService.SetCourseTags(ctx, id, tags)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, course)
}

// PublishCourse 发布课程
func (s *CourseAppService) PublishCourse(ctx context.Context, id string) (*dto.CourseDetailDTO, error) {
	// 发布课程
	course, err := s.courseService.PublishCourse(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, course)
}

// ArchiveCourse 归档课程
func (s *CourseAppService) ArchiveCourse(ctx context.Context, id string) (*dto.CourseDetailDTO, error) {
	// 归档课程
	course, err := s.courseService.ArchiveCourse(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, course)
}

// GetCourseDetail 获取课程详情
func (s *CourseAppService) GetCourseDetail(ctx context.Context, id string) (*dto.CourseDetailDTO, error) {
	// 获取课程
	course, err := s.courseService.GetCourseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return s.convertToDetailDTO(ctx, course)
}

// ListCourses 获取课程列表
func (s *CourseAppService) ListCourses(ctx context.Context, params repository.CourseQueryParams, page, pageSize int) ([]*dto.CourseBasicDTO, int64, error) {
	// 获取课程列表
	courses, total, err := s.courseService.ListCourses(ctx, params, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为 DTO 列表
	dtos := make([]*dto.CourseBasicDTO, 0, len(courses))
	for _, course := range courses {
		basicDTO, err := s.convertToBasicDTO(ctx, course)
		if err != nil {
			continue
		}
		dtos = append(dtos, basicDTO)
	}

	return dtos, total, nil
}

// GetCourseFull 获取完整课程信息（包含章节和课时）
func (s *CourseAppService) GetCourseFull(ctx context.Context, id string) (*dto.CourseFullDTO, error) {
	// 获取课程详情
	courseDetail, err := s.GetCourseDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取章节列表
	chapters, err := s.courseService.GetCourseChapters(ctx, id)
	if err != nil {
		return nil, err
	}

	// 构建章节 DTO 列表
	chapterDTOs := make([]*dto.ChapterDTO, 0, len(chapters))
	for _, chapter := range chapters {
		// 获取章节下的课时
		lessons, err := s.courseService.GetChapterLessons(ctx, chapter.ID)
		if err != nil {
			continue
		}

		// 构建课时 DTO 列表
		lessonDTOs := make([]*dto.LessonBasicDTO, 0, len(lessons))
		for _, lesson := range lessons {
			lessonDTOs = append(lessonDTOs, &dto.LessonBasicDTO{
				ID:          lesson.ID,
				Title:       lesson.Title,
				Type:        string(lesson.Type),
				IsFree:      lesson.IsFree,
				Duration:    lesson.Duration,
				Order:       lesson.Order,
				Status:      string(lesson.Status),
			})
		}

		// 构建章节 DTO
		chapterDTOs = append(chapterDTOs, &dto.ChapterDTO{
			ID:          chapter.ID,
			Title:       chapter.Title,
			Description: chapter.Description,
			Order:       chapter.Order,
			LessonCount: chapter.LessonCount,
			Duration:    chapter.Duration,
			Lessons:     lessonDTOs,
		})
	}

	// 构建完整课程 DTO
	fullDTO := &dto.CourseFullDTO{
		CourseDetailDTO: *courseDetail,
		Chapters:        chapterDTOs,
	}

	return fullDTO, nil
}

// GetCourseStats 获取课程统计信息
func (s *CourseAppService) GetCourseStats(ctx context.Context) (*dto.LearnStats, error) {
	// 获取课程统计
	totalCourses, publishedCount, draftCount, err := s.courseService.GetCourseStats(ctx)
	if err != nil {
		return nil, err
	}

	// 获取教师统计
	totalTeachers, activeTeachers, inactiveTeachers, err := s.teacherService.GetTeacherStats(ctx)
	if err != nil {
		return nil, err
	}

	// 获取分类统计
	totalCategories, level1Categories, level2Categories, err := s.categoryService.GetCategoryStats(ctx)
	if err != nil {
		return nil, err
	}

	// 获取报名统计
	totalEnrollments, activeEnrollments, completedEnrollments, recentEnrollments, err := s.enrollmentService.GetEnrollmentStats(ctx)
	if err != nil {
		return nil, err
	}

	// 构建统计 DTO
	stats := &dto.LearnStats{
		CoursesCount: int(totalCourses),
		CourseStats: dto.CourseStats{
			PublishedCount: int(publishedCount),
			DraftCount:     int(draftCount),
		},
		TeachersCount: int(totalTeachers),
		TeacherStats: dto.TeacherStats{
			ActiveCount:   int(activeTeachers),
			InactiveCount: int(inactiveTeachers),
		},
		CategoriesCount: int(totalCategories),
		CategoryStats: dto.CategoryStats{
			Level1Count: int(level1Categories),
			Level2Count: int(level2Categories),
		},
		EnrollmentsCount: int(totalEnrollments),
		EnrollmentStats: dto.EnrollmentStats{
			ActiveCount:    int(activeEnrollments),
			CompletedCount: int(completedEnrollments),
			RecentCount:    int(recentEnrollments),
		},
	}

	return stats, nil
}

// 辅助函数：转换课程实体到基本 DTO
func (s *CourseAppService) convertToBasicDTO(ctx context.Context, course *entity.Course) (*dto.CourseBasicDTO, error) {
	// 获取教师信息
	teacher, err := s.teacherService.GetTeacherByID(ctx, course.TeacherID)
	if err != nil {
		return nil, err
	}

	return &dto.CourseBasicDTO{
		ID:           course.ID,
		Title:        course.Title,
		Subtitle:     course.Subtitle,
		CoverImage:   course.CoverImage,
		TeacherID:    course.TeacherID,
		TeacherName:  teacher.Name,
		Level:        string(course.Level),
		Status:       string(course.Status),
		EnrollCount:  course.EnrollCount,
		LessonCount:  course.LessonCount,
		Duration:     course.Duration,
		Rating:       course.Rating,
		CreatedAt:    course.CreatedAt,
	}, nil
}

// 辅助函数：转换课程实体到详细 DTO
func (s *CourseAppService) convertToDetailDTO(ctx context.Context, course *entity.Course) (*dto.CourseDetailDTO, error) {
	// 获取教师信息
	teacher, err := s.teacherService.GetTeacherByID(ctx, course.TeacherID)
	if err != nil {
		return nil, err
	}

	// 获取分类信息
	var categoryDTOs []*dto.CategoryDTO
	if len(course.CategoryIDs) > 0 {
		categories, err := s.courseService.GetCourseCategories(ctx, course.ID)
		if err == nil {
			categoryDTOs = make([]*dto.CategoryDTO, 0, len(categories))
			for _, cat := range categories {
				categoryDTOs = append(categoryDTOs, &dto.CategoryDTO{
					ID:          cat.ID,
					Name:        cat.Name,
					Description: cat.Description,
					Icon:        cat.Icon,
					Level:       cat.Level,
					IsActive:    cat.IsActive,
				})
			}
		}
	}

	// 构建详细 DTO
	detailDTO := &dto.CourseDetailDTO{
		ID:           course.ID,
		Title:        course.Title,
		Subtitle:     course.Subtitle,
		Description:  course.Description,
		CoverImage:   course.CoverImage,
		PreviewVideo: course.PreviewVideo,
		TeacherID:    course.TeacherID,
		TeacherName:  teacher.Name,
		TeacherTitle: teacher.Title,
		TeacherAvatar: teacher.Avatar,
		Level:        string(course.Level),
		Status:       string(course.Status),
		Tags:         course.Tags,
		Categories:   categoryDTOs,
		EnrollCount:  course.EnrollCount,
		ChapterCount: course.ChapterCount,
		LessonCount:  course.LessonCount,
		Duration:     course.Duration,
		Rating:       course.Rating,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
		PublishedAt:  course.PublishedAt,
	}

	return detailDTO, nil
}
