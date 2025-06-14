package learn

import (
	"context"
	"fmt"

	"wz-backend-go/internal/domain/learn/dto"
	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// ProgressAppService 学习进度应用服务
type ProgressAppService struct {
	progressRepo   repository.ProgressRepository
	courseRepo     repository.CourseRepository
	lessonRepo     repository.LessonRepository
	enrollmentRepo repository.EnrollmentRepository
	eventBus       event.EventBus
	unitOfWork     database.UnitOfWork
}

// NewProgressAppService 创建学习进度应用服务
func NewProgressAppService(
	progressRepo repository.ProgressRepository,
	courseRepo repository.CourseRepository,
	lessonRepo repository.LessonRepository,
	enrollmentRepo repository.EnrollmentRepository,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) *ProgressAppService {
	return &ProgressAppService{
		progressRepo:   progressRepo,
		courseRepo:     courseRepo,
		lessonRepo:     lessonRepo,
		enrollmentRepo: enrollmentRepo,
		eventBus:       eventBus,
		unitOfWork:     unitOfWork,
	}
}

// UpdateLessonProgress 更新课时学习进度
func (s *ProgressAppService) UpdateLessonProgress(ctx context.Context, userID, lessonID string, watchedDuration int) (*dto.ProgressDTO, error) {
	// 获取课时信息
	lesson, err := s.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("获取课时信息失败: %w", err)
	}
	if lesson == nil {
		return nil, fmt.Errorf("课时不存在")
	}

	// 检查用户是否已报名该课程
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, lesson.CourseID)
	if err != nil {
		return nil, fmt.Errorf("检查报名状态失败: %w", err)
	}
	if enrollment == nil {
		return nil, fmt.Errorf("用户未报名该课程")
	}

	// 获取或创建学习进度
	progress, err := s.progressRepo.GetByUserAndLesson(ctx, userID, lessonID)
	if err != nil {
		return nil, fmt.Errorf("获取学习进度失败: %w", err)
	}

	if progress == nil {
		// 创建新的学习进度
		progress = entity.NewProgress(userID, lesson.CourseID, lessonID, lesson.Duration*60) // 转换为秒
	}

	// 更新进度
	progress.UpdateProgress(watchedDuration)

	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		if progress.CreatedAt.IsZero() {
			// 新创建的进度
			if err := s.progressRepo.Create(ctx, progress); err != nil {
				return fmt.Errorf("创建学习进度失败: %w", err)
			}
		} else {
			// 更新现有进度
			if err := s.progressRepo.Update(ctx, progress); err != nil {
				return fmt.Errorf("更新学习进度失败: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.convertToDTO(ctx, progress), nil
}

// CompleteLessonProgress 标记课时完成
func (s *ProgressAppService) CompleteLessonProgress(ctx context.Context, userID, lessonID string) (*dto.ProgressDTO, error) {
	// 获取课时信息
	lesson, err := s.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("获取课时信息失败: %w", err)
	}
	if lesson == nil {
		return nil, fmt.Errorf("课时不存在")
	}

	// 检查用户是否已报名该课程
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, lesson.CourseID)
	if err != nil {
		return nil, fmt.Errorf("检查报名状态失败: %w", err)
	}
	if enrollment == nil {
		return nil, fmt.Errorf("用户未报名该课程")
	}

	// 获取或创建学习进度
	progress, err := s.progressRepo.GetByUserAndLesson(ctx, userID, lessonID)
	if err != nil {
		return nil, fmt.Errorf("获取学习进度失败: %w", err)
	}

	if progress == nil {
		// 创建新的学习进度并直接标记完成
		progress = entity.NewProgress(userID, lesson.CourseID, lessonID, lesson.Duration*60)
	}

	// 标记完成
	progress.Complete()

	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		if progress.CreatedAt.IsZero() {
			// 新创建的进度
			if err := s.progressRepo.Create(ctx, progress); err != nil {
				return fmt.Errorf("创建学习进度失败: %w", err)
			}
		} else {
			// 更新现有进度
			if err := s.progressRepo.Update(ctx, progress); err != nil {
				return fmt.Errorf("更新学习进度失败: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.convertToDTO(ctx, progress), nil
}

// ResetLessonProgress 重置课时进度
func (s *ProgressAppService) ResetLessonProgress(ctx context.Context, userID, lessonID string) (*dto.ProgressDTO, error) {
	// 获取学习进度
	progress, err := s.progressRepo.GetByUserAndLesson(ctx, userID, lessonID)
	if err != nil {
		return nil, fmt.Errorf("获取学习进度失败: %w", err)
	}
	if progress == nil {
		return nil, fmt.Errorf("学习进度不存在")
	}

	// 重置进度
	progress.Reset()

	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		if err := s.progressRepo.Update(ctx, progress); err != nil {
			return fmt.Errorf("重置学习进度失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.convertToDTO(ctx, progress), nil
}

// GetUserProgress 获取用户学习进度
func (s *ProgressAppService) GetUserProgress(ctx context.Context, userID string, page, pageSize int) ([]*dto.ProgressDTO, int64, error) {
	offset := (page - 1) * pageSize
	progresses, err := s.progressRepo.ListByUser(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户学习进度失败: %w", err)
	}

	total, err := s.progressRepo.CountByUser(ctx, userID, "")
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户学习进度总数失败: %w", err)
	}

	dtos := make([]*dto.ProgressDTO, len(progresses))
	for i, progress := range progresses {
		dtos[i] = s.convertToDTO(ctx, progress)
	}

	return dtos, total, nil
}

// GetCourseProgress 获取用户在某课程的学习进度
func (s *ProgressAppService) GetCourseProgress(ctx context.Context, userID, courseID string) ([]*dto.ProgressDTO, error) {
	progresses, err := s.progressRepo.ListByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("获取课程学习进度失败: %w", err)
	}

	dtos := make([]*dto.ProgressDTO, len(progresses))
	for i, progress := range progresses {
		dtos[i] = s.convertToDTO(ctx, progress)
	}

	return dtos, nil
}

// GetRecentProgress 获取用户最近学习进度
func (s *ProgressAppService) GetRecentProgress(ctx context.Context, userID string, limit int) ([]*dto.ProgressDTO, error) {
	progresses, err := s.progressRepo.ListRecentProgress(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取最近学习进度失败: %w", err)
	}

	dtos := make([]*dto.ProgressDTO, len(progresses))
	for i, progress := range progresses {
		dtos[i] = s.convertToDTO(ctx, progress)
	}

	return dtos, nil
}

// GetCourseProgressStats 获取课程学习进度统计
func (s *ProgressAppService) GetCourseProgressStats(ctx context.Context, userID, courseID string) (*dto.CourseProgressStatsDTO, error) {
	// 获取课程整体进度
	overallProgress, err := s.progressRepo.GetCourseProgressByUser(ctx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("获取课程整体进度失败: %w", err)
	}

	// 获取已完成课时数
	completedCount, err := s.progressRepo.CountByUserAndCourse(ctx, userID, courseID, entity.ProgressStatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("获取已完成课时数失败: %w", err)
	}

	// 获取总课时数
	totalCount, err := s.lessonRepo.CountByCourseID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("获取总课时数失败: %w", err)
	}

	// 获取学习中课时数
	inProgressCount, err := s.progressRepo.CountByUserAndCourse(ctx, userID, courseID, entity.ProgressStatusInProgress)
	if err != nil {
		return nil, fmt.Errorf("获取学习中课时数失败: %w", err)
	}

	return &dto.CourseProgressStatsDTO{
		CourseID:         courseID,
		UserID:           userID,
		OverallProgress:  overallProgress,
		CompletedCount:   completedCount,
		InProgressCount:  inProgressCount,
		TotalCount:       totalCount,
		CompletionRate:   float64(completedCount) / float64(totalCount),
	}, nil
}

// GetUserOverallStats 获取用户整体学习统计
func (s *ProgressAppService) GetUserOverallStats(ctx context.Context, userID string) (*dto.UserProgressStatsDTO, error) {
	// 获取整体学习进度
	overallProgress, err := s.progressRepo.GetOverallProgressByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取整体学习进度失败: %w", err)
	}

	// 获取已完成课时总数
	completedCount, err := s.progressRepo.CountByUser(ctx, userID, entity.ProgressStatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("获取已完成课时总数失败: %w", err)
	}

	// 获取学习中课时总数
	inProgressCount, err := s.progressRepo.CountByUser(ctx, userID, entity.ProgressStatusInProgress)
	if err != nil {
		return nil, fmt.Errorf("获取学习中课时总数失败: %w", err)
	}

	// 获取总学习课时数
	totalCount, err := s.progressRepo.CountByUser(ctx, userID, "")
	if err != nil {
		return nil, fmt.Errorf("获取总学习课时数失败: %w", err)
	}

	return &dto.UserProgressStatsDTO{
		UserID:          userID,
		OverallProgress: overallProgress,
		CompletedCount:  completedCount,
		InProgressCount: inProgressCount,
		TotalCount:      totalCount,
	}, nil
}

// InitializeCourseProgress 初始化课程学习进度
func (s *ProgressAppService) InitializeCourseProgress(ctx context.Context, userID, courseID string) error {
	// 检查用户是否已报名该课程
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return fmt.Errorf("检查报名状态失败: %w", err)
	}
	if enrollment == nil {
		return fmt.Errorf("用户未报名该课程")
	}

	// 获取课程所有课时
	lessons, err := s.lessonRepo.ListByCourseID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("获取课程课时失败: %w", err)
	}

	// 为每个课时创建进度记录
	progresses := make([]*entity.Progress, 0, len(lessons))
	for _, lesson := range lessons {
		// 检查是否已存在进度记录
		existingProgress, err := s.progressRepo.GetByUserAndLesson(ctx, userID, lesson.ID)
		if err != nil {
			return fmt.Errorf("检查课时进度失败: %w", err)
		}
		if existingProgress == nil {
			progress := entity.NewProgress(userID, courseID, lesson.ID, lesson.Duration*60)
			progresses = append(progresses, progress)
		}
	}

	if len(progresses) > 0 {
		err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
			if err := s.progressRepo.BatchCreate(ctx, progresses); err != nil {
				return fmt.Errorf("批量创建学习进度失败: %w", err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// convertToDTO 转换为DTO
func (s *ProgressAppService) convertToDTO(ctx context.Context, progress *entity.Progress) *dto.ProgressDTO {
	return &dto.ProgressDTO{
		ID:               progress.ID,
		UserID:           progress.UserID,
		CourseID:         progress.CourseID,
		LessonID:         progress.LessonID,
		Status:           string(progress.Status),
		WatchedDuration:  progress.WatchedDuration,
		TotalDuration:    progress.TotalDuration,
		CompletionRate:   progress.CompletionRate,
		ProgressPercent:  progress.GetProgressPercentage(),
		LastWatchedAt:    progress.LastWatchedAt,
		CompletedAt:      progress.CompletedAt,
		CreatedAt:        progress.CreatedAt,
		UpdatedAt:        progress.UpdatedAt,
	}
}
