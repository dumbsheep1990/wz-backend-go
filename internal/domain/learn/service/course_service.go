package service

import (
	"context"
	"errors"
	"time"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
)

var (
	ErrCourseNotFound        = errors.New("课程不存在")
	ErrTeacherNotFound       = errors.New("讲师不存在")
	ErrInvalidCourseStatus   = errors.New("无效的课程状态")
	ErrInvalidCourseDuration = errors.New("无效的课程时长")
)

// CourseService 课程领域服务，处理课程相关的业务逻辑
type CourseService struct {
	courseRepo    repository.CourseRepository
	teacherRepo   repository.TeacherRepository
	chapterRepo   repository.ChapterRepository
	lessonRepo    repository.LessonRepository
	categoryRepo  repository.CategoryRepository
	enrollmentRepo repository.EnrollmentRepository
}

// NewCourseService 创建课程服务
func NewCourseService(
	courseRepo repository.CourseRepository,
	teacherRepo repository.TeacherRepository,
	chapterRepo repository.ChapterRepository,
	lessonRepo repository.LessonRepository,
	categoryRepo repository.CategoryRepository,
	enrollmentRepo repository.EnrollmentRepository,
) *CourseService {
	return &CourseService{
		courseRepo:    courseRepo,
		teacherRepo:   teacherRepo,
		chapterRepo:   chapterRepo,
		lessonRepo:    lessonRepo,
		categoryRepo:  categoryRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

// CreateCourse 创建新课程
func (s *CourseService) CreateCourse(ctx context.Context, title, description string, teacherID string) (*entity.Course, error) {
	// 验证讲师是否存在
	teacher, err := s.teacherRepo.GetByID(ctx, teacherID)
	if err != nil {
		return nil, ErrTeacherNotFound
	}

	// 创建课程
	course := entity.NewCourse(title, teacherID)
	course.Description = description

	// 保存课程
	if err := s.courseRepo.Create(ctx, course); err != nil {
		return nil, err
	}

	// 增加讲师的课程数
	teacher.IncrementCourseCount()
	if err := s.teacherRepo.Update(ctx, teacher); err != nil {
		return nil, err
	}

	return course, nil
}

// UpdateCourse 更新课程信息
func (s *CourseService) UpdateCourse(ctx context.Context, id, title, subtitle, description, cover string, level entity.CourseLevel, price, discountPrice float64) (*entity.Course, error) {
	// 获取课程
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	// 更新课程信息
	course.Update(title, subtitle, description, cover, level, price, discountPrice)

	// 保存课程
	if err := s.courseRepo.Update(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// PublishCourse 发布课程
func (s *CourseService) PublishCourse(ctx context.Context, id string) (*entity.Course, error) {
	// 获取课程
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	// 检查课程是否有内容
	chaptersCount, err := s.chapterRepo.CountByCourseID(ctx, id)
	if err != nil {
		return nil, err
	}

	lessonsCount, err := s.lessonRepo.CountByCourseID(ctx, id)
	if err != nil {
		return nil, err
	}

	if chaptersCount == 0 || lessonsCount == 0 {
		return nil, errors.New("课程没有内容，无法发布")
	}

	// 更新课程状态
	course.Publish()
	course.UpdateLessonCounts(int(chaptersCount), int(lessonsCount))

	// 保存课程
	if err := s.courseRepo.Update(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// ArchiveCourse 归档课程
func (s *CourseService) ArchiveCourse(ctx context.Context, id string) (*entity.Course, error) {
	// 获取课程
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	// 更新课程状态
	course.Archive()

	// 保存课程
	if err := s.courseRepo.Update(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// SetCourseCategories 设置课程分类
func (s *CourseService) SetCourseCategories(ctx context.Context, courseID string, categoryIDs []string) (*entity.Course, error) {
	// 获取课程
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	// 移除旧分类关联
	oldCategoryIDs := course.CategoryIDs
	for _, oldID := range oldCategoryIDs {
		category, err := s.categoryRepo.GetByID(ctx, oldID)
		if err == nil {
			category.DecrementCourseCount()
			if err := s.categoryRepo.Update(ctx, category); err != nil {
				return nil, err
			}
		}
	}

	// 添加新分类关联
	for _, newID := range categoryIDs {
		category, err := s.categoryRepo.GetByID(ctx, newID)
		if err == nil {
			category.IncrementCourseCount()
			if err := s.categoryRepo.Update(ctx, category); err != nil {
				return nil, err
			}
		}
	}

	// 更新课程分类
	course.SetCategories(categoryIDs)
	if err := s.courseRepo.Update(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// SetCourseTags 设置课程标签
func (s *CourseService) SetCourseTags(ctx context.Context, courseID string, tags []string) (*entity.Course, error) {
	// 获取课程
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	// 更新课程标签
	course.SetTags(tags)
	if err := s.courseRepo.Update(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// DeleteCourse 删除课程
func (s *CourseService) DeleteCourse(ctx context.Context, id string) error {
	// 获取课程
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return ErrCourseNotFound
	}

	// 处理讲师的课程数
	teacher, err := s.teacherRepo.GetByID(ctx, course.TeacherID)
	if err == nil {
		teacher.DecrementCourseCount()
		if err := s.teacherRepo.Update(ctx, teacher); err != nil {
			return err
		}
	}

	// 处理分类课程数
	for _, categoryID := range course.CategoryIDs {
		category, err := s.categoryRepo.GetByID(ctx, categoryID)
		if err == nil {
			category.DecrementCourseCount()
			if err := s.categoryRepo.Update(ctx, category); err != nil {
				return err
			}
		}
	}

	// 删除所有关联的章节和课时
	chapters, err := s.chapterRepo.ListByCourseID(ctx, id)
	if err == nil {
		for _, chapter := range chapters {
			lessons, err := s.lessonRepo.ListByChapterID(ctx, chapter.ID)
			if err == nil {
				for _, lesson := range lessons {
					if err := s.lessonRepo.Delete(ctx, lesson.ID); err != nil {
						return err
					}
				}
			}
			if err := s.chapterRepo.Delete(ctx, chapter.ID); err != nil {
				return err
			}
		}
	}

	// 删除课程
	if err := s.courseRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetCourseStats 获取课程统计数据
func (s *CourseService) GetCourseStats(ctx context.Context) (totalCount, publishedCount, draftCount int64, err error) {
	totalCount, err = s.courseRepo.CountAll(ctx)
	if err != nil {
		return 0, 0, 0, err
	}

	publishedCount, err = s.courseRepo.CountByStatus(ctx, entity.CourseStatusPublished)
	if err != nil {
		return 0, 0, 0, err
	}

	draftCount, err = s.courseRepo.CountByStatus(ctx, entity.CourseStatusDraft)
	if err != nil {
		return 0, 0, 0, err
	}

	return totalCount, publishedCount, draftCount, nil
}
