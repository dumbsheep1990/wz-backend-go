package learn

import (
	"context"

	"wz-backend-go/internal/domain/learn/dto"
	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/service"
)

// ChapterLessonAppService 章节和课时应用服务，处理章节和课时相关的应用场景
type ChapterLessonAppService struct {
	chapterLessonService *service.ChapterLessonService
}

// NewChapterLessonAppService 创建章节和课时应用服务
func NewChapterLessonAppService(
	chapterLessonService *service.ChapterLessonService,
) *ChapterLessonAppService {
	return &ChapterLessonAppService{
		chapterLessonService: chapterLessonService,
	}
}

// CreateChapter 创建章节
func (s *ChapterLessonAppService) CreateChapter(ctx context.Context, courseID, title, description string, order int) (*dto.ChapterDTO, error) {
	// 创建章节
	chapter, err := s.chapterLessonService.CreateChapter(ctx, courseID, title, description, order)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return &dto.ChapterDTO{
		ID:          chapter.ID,
		Title:       chapter.Title,
		Description: chapter.Description,
		Order:       chapter.Order,
		LessonCount: chapter.LessonCount,
		Duration:    chapter.Duration,
		CourseID:    chapter.CourseID,
		Lessons:     make([]*dto.LessonBasicDTO, 0),
	}, nil
}

// UpdateChapter 更新章节
func (s *ChapterLessonAppService) UpdateChapter(ctx context.Context, id, title, description string, order int) (*dto.ChapterDTO, error) {
	// 更新章节
	chapter, err := s.chapterLessonService.UpdateChapter(ctx, id, title, description, order)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	return &dto.ChapterDTO{
		ID:          chapter.ID,
		Title:       chapter.Title,
		Description: chapter.Description,
		Order:       chapter.Order,
		LessonCount: chapter.LessonCount,
		Duration:    chapter.Duration,
		CourseID:    chapter.CourseID,
		Lessons:     make([]*dto.LessonBasicDTO, 0),
	}, nil
}

// DeleteChapter 删除章节
func (s *ChapterLessonAppService) DeleteChapter(ctx context.Context, id string) error {
	return s.chapterLessonService.DeleteChapter(ctx, id)
}

// GetChapterByID 获取章节详情
func (s *ChapterLessonAppService) GetChapterByID(ctx context.Context, id string) (*dto.ChapterDTO, error) {
	// 获取章节信息
	chapter, err := s.chapterLessonService.GetChapterByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取章节下的课时
	lessons, err := s.chapterLessonService.GetLessonsByChapterID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 构建课时 DTO 列表
	lessonDTOs := make([]*dto.LessonBasicDTO, 0, len(lessons))
	for _, lesson := range lessons {
		lessonDTOs = append(lessonDTOs, &dto.LessonBasicDTO{
			ID:          lesson.ID,
			Title:       lesson.Title,
			Description: lesson.Description,
			Type:        string(lesson.Type),
			IsFree:      lesson.IsFree,
			Duration:    lesson.Duration,
			Order:       lesson.Order,
			Status:      string(lesson.Status),
		})
	}

	// 构建章节 DTO
	chapterDTO := &dto.ChapterDTO{
		ID:          chapter.ID,
		Title:       chapter.Title,
		Description: chapter.Description,
		Order:       chapter.Order,
		LessonCount: chapter.LessonCount,
		Duration:    chapter.Duration,
		CourseID:    chapter.CourseID,
		Lessons:     lessonDTOs,
	}

	return chapterDTO, nil
}

// GetChaptersByCourseID 获取课程的所有章节
func (s *ChapterLessonAppService) GetChaptersByCourseID(ctx context.Context, courseID string) ([]*dto.ChapterDTO, error) {
	// 获取课程的所有章节
	chapters, err := s.chapterLessonService.GetChaptersByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// 构建章节 DTO 列表
	chapterDTOs := make([]*dto.ChapterDTO, 0, len(chapters))
	for _, chapter := range chapters {
		// 获取章节下的所有课时
		lessons, err := s.chapterLessonService.GetLessonsByChapterID(ctx, chapter.ID)
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
			CourseID:    chapter.CourseID,
			Lessons:     lessonDTOs,
		})
	}

	return chapterDTOs, nil
}

// CreateLesson 创建课时
func (s *ChapterLessonAppService) CreateLesson(ctx context.Context, courseID, chapterID, title, description string, order int, lessonType entity.LessonType, isFree bool) (*dto.LessonBasicDTO, error) {
	// 创建课时
	lesson, err := s.chapterLessonService.CreateLesson(ctx, courseID, chapterID, title, description, order, lessonType, isFree)
	if err != nil {
		return nil, err
	}

	// 转换为基本 DTO
	return &dto.LessonBasicDTO{
		ID:          lesson.ID,
		Title:       lesson.Title,
		Description: lesson.Description,
		Type:        string(lesson.Type),
		IsFree:      lesson.IsFree,
		Duration:    lesson.Duration,
		Order:       lesson.Order,
		Status:      string(lesson.Status),
		CourseID:    lesson.CourseID,
		ChapterID:   lesson.ChapterID,
	}, nil
}

// UpdateLesson 更新课时基本信息
func (s *ChapterLessonAppService) UpdateLesson(ctx context.Context, id, title, description string, order int, isFree bool) (*dto.LessonBasicDTO, error) {
	// 更新课时
	lesson, err := s.chapterLessonService.UpdateLesson(ctx, id, title, description, order, isFree)
	if err != nil {
		return nil, err
	}

	// 转换为基本 DTO
	return &dto.LessonBasicDTO{
		ID:          lesson.ID,
		Title:       lesson.Title,
		Description: lesson.Description,
		Type:        string(lesson.Type),
		IsFree:      lesson.IsFree,
		Duration:    lesson.Duration,
		Order:       lesson.Order,
		Status:      string(lesson.Status),
		CourseID:    lesson.CourseID,
		ChapterID:   lesson.ChapterID,
	}, nil
}

// PublishLesson 发布课时
func (s *ChapterLessonAppService) PublishLesson(ctx context.Context, id string) (*dto.LessonBasicDTO, error) {
	// 发布课时
	lesson, err := s.chapterLessonService.PublishLesson(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为基本 DTO
	return &dto.LessonBasicDTO{
		ID:          lesson.ID,
		Title:       lesson.Title,
		Description: lesson.Description,
		Type:        string(lesson.Type),
		IsFree:      lesson.IsFree,
		Duration:    lesson.Duration,
		Order:       lesson.Order,
		Status:      string(lesson.Status),
		CourseID:    lesson.CourseID,
		ChapterID:   lesson.ChapterID,
	}, nil
}

// DeleteLesson 删除课时
func (s *ChapterLessonAppService) DeleteLesson(ctx context.Context, id string) error {
	return s.chapterLessonService.DeleteLesson(ctx, id)
}

// GetLessonByID 获取课时详情
func (s *ChapterLessonAppService) GetLessonByID(ctx context.Context, id string) (*dto.LessonDetailDTO, error) {
	// 获取课时
	lesson, err := s.chapterLessonService.GetLessonByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为详细 DTO
	return &dto.LessonDetailDTO{
		ID:          lesson.ID,
		Title:       lesson.Title,
		Description: lesson.Description,
		Type:        string(lesson.Type),
		IsFree:      lesson.IsFree,
		Duration:    lesson.Duration,
		Size:        lesson.Size,
		Order:       lesson.Order,
		Status:      string(lesson.Status),
		CourseID:    lesson.CourseID,
		ChapterID:   lesson.ChapterID,
		VideoURL:    lesson.VideoURL,
		ArticleContent: lesson.ArticleContent,
		AudioURL:    lesson.AudioURL,
		CreatedAt:   lesson.CreatedAt,
		UpdatedAt:   lesson.UpdatedAt,
	}, nil
}

// SetVideoContent 设置视频课时内容
func (s *ChapterLessonAppService) SetVideoContent(ctx context.Context, id, videoURL string, duration int, size int64) (*dto.LessonDetailDTO, error) {
	// 设置视频内容
	lesson, err := s.chapterLessonService.SetVideoContent(ctx, id, videoURL, duration, size)
	if err != nil {
		return nil, err
	}

	// 转换为详细 DTO
	return &dto.LessonDetailDTO{
		ID:          lesson.ID,
		Title:       lesson.Title,
		Description: lesson.Description,
		Type:        string(lesson.Type),
		IsFree:      lesson.IsFree,
		Duration:    lesson.Duration,
		Size:        lesson.Size,
		Order:       lesson.Order,
		Status:      string(lesson.Status),
		CourseID:    lesson.CourseID,
		ChapterID:   lesson.ChapterID,
		VideoURL:    lesson.VideoURL,
		CreatedAt:   lesson.CreatedAt,
		UpdatedAt:   lesson.UpdatedAt,
	}, nil
}

// SetArticleContent 设置文章课时内容
func (s *ChapterLessonAppService) SetArticleContent(ctx context.Context, id, content string) (*dto.LessonDetailDTO, error) {
	// 设置文章内容
	lesson, err := s.chapterLessonService.SetArticleContent(ctx, id, content)
	if err != nil {
		return nil, err
	}

	// 转换为详细 DTO
	return &dto.LessonDetailDTO{
		ID:            lesson.ID,
		Title:         lesson.Title,
		Description:   lesson.Description,
		Type:          string(lesson.Type),
		IsFree:        lesson.IsFree,
		Duration:      lesson.Duration,
		Order:         lesson.Order,
		Status:        string(lesson.Status),
		CourseID:      lesson.CourseID,
		ChapterID:     lesson.ChapterID,
		ArticleContent: lesson.ArticleContent,
		CreatedAt:     lesson.CreatedAt,
		UpdatedAt:     lesson.UpdatedAt,
	}, nil
}

// SetAudioContent 设置音频课时内容
func (s *ChapterLessonAppService) SetAudioContent(ctx context.Context, id, audioURL string, duration int, size int64) (*dto.LessonDetailDTO, error) {
	// 设置音频内容
	lesson, err := s.chapterLessonService.SetAudioContent(ctx, id, audioURL, duration, size)
	if err != nil {
		return nil, err
	}

	// 转换为详细 DTO
	return &dto.LessonDetailDTO{
		ID:          lesson.ID,
		Title:       lesson.Title,
		Description: lesson.Description,
		Type:        string(lesson.Type),
		IsFree:      lesson.IsFree,
		Duration:    lesson.Duration,
		Size:        lesson.Size,
		Order:       lesson.Order,
		Status:      string(lesson.Status),
		CourseID:    lesson.CourseID,
		ChapterID:   lesson.ChapterID,
		AudioURL:    lesson.AudioURL,
		CreatedAt:   lesson.CreatedAt,
		UpdatedAt:   lesson.UpdatedAt,
	}, nil
}
