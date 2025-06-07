package service

import (
	"context"
	"errors"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
)

var (
	ErrChapterNotFound = errors.New("章节不存在")
	ErrLessonNotFound  = errors.New("课时不存在")
	ErrInvalidOrder    = errors.New("无效的排序值")
)

// ChapterLessonService 章节和课时领域服务
type ChapterLessonService struct {
	chapterRepo repository.ChapterRepository
	lessonRepo  repository.LessonRepository
	courseRepo  repository.CourseRepository
}

// NewChapterLessonService 创建章节和课时服务
func NewChapterLessonService(
	chapterRepo repository.ChapterRepository,
	lessonRepo repository.LessonRepository,
	courseRepo repository.CourseRepository,
) *ChapterLessonService {
	return &ChapterLessonService{
		chapterRepo: chapterRepo,
		lessonRepo:  lessonRepo,
		courseRepo:  courseRepo,
	}
}

// CreateChapter 创建章节
func (s *ChapterLessonService) CreateChapter(ctx context.Context, courseID, title, description string, order int) (*entity.Chapter, error) {
	// 验证课程是否存在
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	// 创建章节
	chapter := entity.NewChapter(courseID, title, order)
	chapter.Description = description

	// 保存章节
	if err := s.chapterRepo.Create(ctx, chapter); err != nil {
		return nil, err
	}

	// 更新课程的章节数量
	chapterCount, _ := s.chapterRepo.CountByCourseID(ctx, courseID)
	lessonCount, _ := s.lessonRepo.CountByCourseID(ctx, courseID)
	course.UpdateLessonCounts(int(chapterCount), int(lessonCount))
	
	if err := s.courseRepo.Update(ctx, course); err != nil {
		return nil, err
	}

	return chapter, nil
}

// UpdateChapter 更新章节
func (s *ChapterLessonService) UpdateChapter(ctx context.Context, id, title, description string, order int) (*entity.Chapter, error) {
	// 获取章节
	chapter, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrChapterNotFound
	}

	// 更新章节信息
	chapter.Update(title, description, order)

	// 保存章节
	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		return nil, err
	}

	return chapter, nil
}

// DeleteChapter 删除章节
func (s *ChapterLessonService) DeleteChapter(ctx context.Context, id string) error {
	// 获取章节
	chapter, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		return ErrChapterNotFound
	}

	// 获取课程ID，用于后续更新课程章节数
	courseID := chapter.CourseID

	// 删除章节下的所有课时
	lessons, err := s.lessonRepo.ListByChapterID(ctx, id)
	if err == nil {
		for _, lesson := range lessons {
			if err := s.lessonRepo.Delete(ctx, lesson.ID); err != nil {
				return err
			}
		}
	}

	// 删除章节
	if err := s.chapterRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 更新课程的章节和课时数量
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err == nil {
		chapterCount, _ := s.chapterRepo.CountByCourseID(ctx, courseID)
		lessonCount, _ := s.lessonRepo.CountByCourseID(ctx, courseID)
		course.UpdateLessonCounts(int(chapterCount), int(lessonCount))
		if err := s.courseRepo.Update(ctx, course); err != nil {
			return err
		}
	}

	return nil
}

// CreateLesson 创建课时
func (s *ChapterLessonService) CreateLesson(ctx context.Context, courseID, chapterID, title, description string, order int, lessonType entity.LessonType, isFree bool) (*entity.Lesson, error) {
	// 验证课程是否存在
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	// 验证章节是否存在
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, ErrChapterNotFound
	}

	// 创建课时
	lesson := entity.NewLesson(courseID, chapterID, title, order, lessonType)
	lesson.Description = description
	lesson.IsFree = isFree

	// 保存课时
	if err := s.lessonRepo.Create(ctx, lesson); err != nil {
		return nil, err
	}

	// 更新章节的课时统计
	lessonCount, err := s.lessonRepo.CountByChapterID(ctx, chapterID)
	if err == nil {
		// 计算章节总时长
		var totalDuration int
		lessons, _ := s.lessonRepo.ListByChapterID(ctx, chapterID)
		for _, l := range lessons {
			totalDuration += l.Duration
		}
		
		chapter.UpdateLessonStats(int(lessonCount), totalDuration)
		if err := s.chapterRepo.Update(ctx, chapter); err != nil {
			return nil, err
		}
	}

	// 更新课程的总课时数
	totalLessonCount, _ := s.lessonRepo.CountByCourseID(ctx, courseID)
	chapterCount, _ := s.chapterRepo.CountByCourseID(ctx, courseID)
	course.UpdateLessonCounts(int(chapterCount), int(totalLessonCount))
	
	// 计算课程总时长
	var courseDuration int
	lessons, _ := s.lessonRepo.ListByCourseID(ctx, courseID)
	for _, l := range lessons {
		courseDuration += l.Duration
	}
	course.Duration = courseDuration
	
	if err := s.courseRepo.Update(ctx, course); err != nil {
		return nil, err
	}

	return lesson, nil
}

// UpdateLesson 更新课时
func (s *ChapterLessonService) UpdateLesson(ctx context.Context, id, title, description string, order int, isFree bool) (*entity.Lesson, error) {
	// 获取课时
	lesson, err := s.lessonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrLessonNotFound
	}

	// 更新课时信息
	lesson.Update(title, description, order, isFree)

	// 保存课时
	if err := s.lessonRepo.Update(ctx, lesson); err != nil {
		return nil, err
	}

	return lesson, nil
}

// PublishLesson 发布课时
func (s *ChapterLessonService) PublishLesson(ctx context.Context, id string) (*entity.Lesson, error) {
	// 获取课时
	lesson, err := s.lessonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrLessonNotFound
	}

	// 发布课时
	lesson.Publish()

	// 保存课时
	if err := s.lessonRepo.Update(ctx, lesson); err != nil {
		return nil, err
	}

	return lesson, nil
}

// SetVideoContent 设置视频课时内容
func (s *ChapterLessonService) SetVideoContent(ctx context.Context, id, videoURL string, duration int, size int64) (*entity.Lesson, error) {
	// 获取课时
	lesson, err := s.lessonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrLessonNotFound
	}

	// 确保课时类型正确
	if lesson.Type != entity.LessonTypeVideo {
		return nil, errors.New("课时类型不是视频")
	}

	// 设置视频内容
	lesson.SetVideo(videoURL, duration, size)

	// 保存课时
	if err := s.lessonRepo.Update(ctx, lesson); err != nil {
		return nil, err
	}

	// 更新章节的课时统计
	lessons, _ := s.lessonRepo.ListByChapterID(ctx, lesson.ChapterID)
	var totalDuration int
	for _, l := range lessons {
		totalDuration += l.Duration
	}
	
	chapter, err := s.chapterRepo.GetByID(ctx, lesson.ChapterID)
	if err == nil {
		chapter.UpdateLessonStats(len(lessons), totalDuration)
		if err := s.chapterRepo.Update(ctx, chapter); err != nil {
			return nil, err
		}
	}

	// 更新课程总时长
	allLessons, _ := s.lessonRepo.ListByCourseID(ctx, lesson.CourseID)
	var courseDuration int
	for _, l := range allLessons {
		courseDuration += l.Duration
	}
	
	course, err := s.courseRepo.GetByID(ctx, lesson.CourseID)
	if err == nil {
		course.Duration = courseDuration
		if err := s.courseRepo.Update(ctx, course); err != nil {
			return nil, err
		}
	}

	return lesson, nil
}

// SetArticleContent 设置文章课时内容
func (s *ChapterLessonService) SetArticleContent(ctx context.Context, id, content string) (*entity.Lesson, error) {
	// 获取课时
	lesson, err := s.lessonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrLessonNotFound
	}

	// 确保课时类型正确
	if lesson.Type != entity.LessonTypeArticle {
		return nil, errors.New("课时类型不是文章")
	}

	// 设置文章内容
	lesson.SetArticle(content)

	// 保存课时
	if err := s.lessonRepo.Update(ctx, lesson); err != nil {
		return nil, err
	}

	return lesson, nil
}

// SetAudioContent 设置音频课时内容
func (s *ChapterLessonService) SetAudioContent(ctx context.Context, id, audioURL string, duration int, size int64) (*entity.Lesson, error) {
	// 获取课时
	lesson, err := s.lessonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrLessonNotFound
	}

	// 确保课时类型正确
	if lesson.Type != entity.LessonTypeAudio {
		return nil, errors.New("课时类型不是音频")
	}

	// 设置音频内容
	lesson.SetAudio(audioURL, duration, size)

	// 保存课时
	if err := s.lessonRepo.Update(ctx, lesson); err != nil {
		return nil, err
	}

	// 更新章节的课时统计
	lessons, _ := s.lessonRepo.ListByChapterID(ctx, lesson.ChapterID)
	var totalDuration int
	for _, l := range lessons {
		totalDuration += l.Duration
	}
	
	chapter, err := s.chapterRepo.GetByID(ctx, lesson.ChapterID)
	if err == nil {
		chapter.UpdateLessonStats(len(lessons), totalDuration)
		if err := s.chapterRepo.Update(ctx, chapter); err != nil {
			return nil, err
		}
	}

	// 更新课程总时长
	allLessons, _ := s.lessonRepo.ListByCourseID(ctx, lesson.CourseID)
	var courseDuration int
	for _, l := range allLessons {
		courseDuration += l.Duration
	}
	
	course, err := s.courseRepo.GetByID(ctx, lesson.CourseID)
	if err == nil {
		course.Duration = courseDuration
		if err := s.courseRepo.Update(ctx, course); err != nil {
			return nil, err
		}
	}

	return lesson, nil
}

// DeleteLesson 删除课时
func (s *ChapterLessonService) DeleteLesson(ctx context.Context, id string) error {
	// 获取课时
	lesson, err := s.lessonRepo.GetByID(ctx, id)
	if err != nil {
		return ErrLessonNotFound
	}

	// 保存课程和章节ID，用于后续更新
	courseID := lesson.CourseID
	chapterID := lesson.ChapterID

	// 删除课时
	if err := s.lessonRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 更新章节的课时统计
	lessons, _ := s.lessonRepo.ListByChapterID(ctx, chapterID)
	var totalDuration int
	for _, l := range lessons {
		totalDuration += l.Duration
	}
	
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err == nil {
		chapter.UpdateLessonStats(len(lessons), totalDuration)
		if err := s.chapterRepo.Update(ctx, chapter); err != nil {
			return err
		}
	}

	// 更新课程的总课时数和时长
	totalLessonCount, _ := s.lessonRepo.CountByCourseID(ctx, courseID)
	chapterCount, _ := s.chapterRepo.CountByCourseID(ctx, courseID)
	
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err == nil {
		course.UpdateLessonCounts(int(chapterCount), int(totalLessonCount))
		
		// 更新课程总时长
		allLessons, _ := s.lessonRepo.ListByCourseID(ctx, courseID)
		var courseDuration int
		for _, l := range allLessons {
			courseDuration += l.Duration
		}
		course.Duration = courseDuration
		
		if err := s.courseRepo.Update(ctx, course); err != nil {
			return err
		}
	}

	return nil
}
