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

// ReviewAppService 评价应用服务
type ReviewAppService struct {
	reviewRepo     repository.ReviewRepository
	courseRepo     repository.CourseRepository
	enrollmentRepo repository.EnrollmentRepository
	eventBus       event.EventBus
	unitOfWork     database.UnitOfWork
}

// NewReviewAppService 创建评价应用服务
func NewReviewAppService(
	reviewRepo repository.ReviewRepository,
	courseRepo repository.CourseRepository,
	enrollmentRepo repository.EnrollmentRepository,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) *ReviewAppService {
	return &ReviewAppService{
		reviewRepo:     reviewRepo,
		courseRepo:     courseRepo,
		enrollmentRepo: enrollmentRepo,
		eventBus:       eventBus,
		unitOfWork:     unitOfWork,
	}
}

// CreateReview 创建课程评价
func (s *ReviewAppService) CreateReview(ctx context.Context, userID, courseID string, rating int, content string) (*dto.ReviewDTO, error) {
	// 验证评分范围
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("评分必须在1-5之间")
	}

	// 检查课程是否存在
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("获取课程信息失败: %w", err)
	}
	if course == nil {
		return nil, fmt.Errorf("课程不存在")
	}

	// 检查用户是否已报名该课程
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("检查报名状态失败: %w", err)
	}
	if enrollment == nil {
		return nil, fmt.Errorf("用户未报名该课程，无法评价")
	}

	// 检查是否已经评价过
	existingReview, err := s.reviewRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("检查评价状态失败: %w", err)
	}
	if existingReview != nil {
		return nil, fmt.Errorf("用户已评价过该课程")
	}

	// 创建评价
	review := entity.NewReview(userID, courseID, rating, content)

	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		if err := s.reviewRepo.Create(ctx, review); err != nil {
			return fmt.Errorf("创建评价失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.convertToDTO(ctx, review), nil
}

// UpdateReview 更新评价
func (s *ReviewAppService) UpdateReview(ctx context.Context, userID, reviewID string, rating int, content string) (*dto.ReviewDTO, error) {
	// 验证评分范围
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("评分必须在1-5之间")
	}

	// 获取评价
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("获取评价失败: %w", err)
	}
	if review == nil {
		return nil, fmt.Errorf("评价不存在")
	}

	// 检查权限
	if review.UserID != userID {
		return nil, fmt.Errorf("无权限修改该评价")
	}

	// 只有待审核状态的评价才能修改
	if review.Status != entity.ReviewStatusPending {
		return nil, fmt.Errorf("只有待审核的评价才能修改")
	}

	// 更新评价
	review.Update(rating, content)

	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		if err := s.reviewRepo.Update(ctx, review); err != nil {
			return fmt.Errorf("更新评价失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.convertToDTO(ctx, review), nil
}

// DeleteReview 删除评价
func (s *ReviewAppService) DeleteReview(ctx context.Context, userID, reviewID string) error {
	// 获取评价
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("获取评价失败: %w", err)
	}
	if review == nil {
		return fmt.Errorf("评价不存在")
	}

	// 检查权限
	if review.UserID != userID {
		return fmt.Errorf("无权限删除该评价")
	}

	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		if err := s.reviewRepo.Delete(ctx, reviewID); err != nil {
			return fmt.Errorf("删除评价失败: %w", err)
		}
		return nil
	})

	return err
}

// ApproveReview 审核通过评价
func (s *ReviewAppService) ApproveReview(ctx context.Context, reviewID string) (*dto.ReviewDTO, error) {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("获取评价失败: %w", err)
	}
	if review == nil {
		return nil, fmt.Errorf("评价不存在")
	}

	if review.Status != entity.ReviewStatusPending {
		return nil, fmt.Errorf("只有待审核的评价才能审核")
	}

	review.Approve()

	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		if err := s.reviewRepo.Update(ctx, review); err != nil {
			return fmt.Errorf("更新评价状态失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.convertToDTO(ctx, review), nil
}

// RejectReview 审核拒绝评价
func (s *ReviewAppService) RejectReview(ctx context.Context, reviewID string) (*dto.ReviewDTO, error) {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("获取评价失败: %w", err)
	}
	if review == nil {
		return nil, fmt.Errorf("评价不存在")
	}

	if review.Status != entity.ReviewStatusPending {
		return nil, fmt.Errorf("只有待审核的评价才能审核")
	}

	review.Reject()

	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		if err := s.reviewRepo.Update(ctx, review); err != nil {
			return fmt.Errorf("更新评价状态失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.convertToDTO(ctx, review), nil
}

// GetReviewsByUser 获取用户的评价列表
func (s *ReviewAppService) GetReviewsByUser(ctx context.Context, userID string, page, pageSize int) ([]*dto.ReviewDTO, int64, error) {
	offset := (page - 1) * pageSize
	reviews, err := s.reviewRepo.ListByUserID(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户评价列表失败: %w", err)
	}

	total, err := s.reviewRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户评价总数失败: %w", err)
	}

	dtos := make([]*dto.ReviewDTO, len(reviews))
	for i, review := range reviews {
		dtos[i] = s.convertToDTO(ctx, review)
	}

	return dtos, total, nil
}

// GetReviewsByCourse 获取课程的评价列表
func (s *ReviewAppService) GetReviewsByCourse(ctx context.Context, courseID string, status entity.ReviewStatus, page, pageSize int) ([]*dto.ReviewDTO, int64, error) {
	offset := (page - 1) * pageSize
	reviews, err := s.reviewRepo.ListByCourseID(ctx, courseID, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取课程评价列表失败: %w", err)
	}

	total, err := s.reviewRepo.CountByCourseID(ctx, courseID, status)
	if err != nil {
		return nil, 0, fmt.Errorf("获取课程评价总数失败: %w", err)
	}

	dtos := make([]*dto.ReviewDTO, len(reviews))
	for i, review := range reviews {
		dtos[i] = s.convertToDTO(ctx, review)
	}

	return dtos, total, nil
}

// GetPendingReviews 获取待审核评价列表
func (s *ReviewAppService) GetPendingReviews(ctx context.Context, page, pageSize int) ([]*dto.ReviewDTO, int64, error) {
	offset := (page - 1) * pageSize
	reviews, err := s.reviewRepo.ListPendingReviews(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取待审核评价列表失败: %w", err)
	}

	total, err := s.reviewRepo.CountByCourseID(ctx, "", entity.ReviewStatusPending)
	if err != nil {
		return nil, 0, fmt.Errorf("获取待审核评价总数失败: %w", err)
	}

	dtos := make([]*dto.ReviewDTO, len(reviews))
	for i, review := range reviews {
		dtos[i] = s.convertToDTO(ctx, review)
	}

	return dtos, total, nil
}

// GetCourseRatingStats 获取课程评分统计
func (s *ReviewAppService) GetCourseRatingStats(ctx context.Context, courseID string) (*dto.CourseRatingStatsDTO, error) {
	// 获取平均评分
	avgRating, err := s.reviewRepo.GetAverageRatingByCourseID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("获取平均评分失败: %w", err)
	}

	// 获取评分分布
	distribution, err := s.reviewRepo.GetRatingDistributionByCourseID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("获取评分分布失败: %w", err)
	}

	// 获取总评价数
	totalCount, err := s.reviewRepo.CountByCourseID(ctx, courseID, entity.ReviewStatusApproved)
	if err != nil {
		return nil, fmt.Errorf("获取评价总数失败: %w", err)
	}

	return &dto.CourseRatingStatsDTO{
		CourseID:      courseID,
		AverageRating: avgRating,
		TotalCount:    totalCount,
		Distribution:  distribution,
	}, nil
}

// convertToDTO 转换为DTO
func (s *ReviewAppService) convertToDTO(ctx context.Context, review *entity.Review) *dto.ReviewDTO {
	return &dto.ReviewDTO{
		ID:         review.ID,
		UserID:     review.UserID,
		CourseID:   review.CourseID,
		Rating:     review.Rating,
		Content:    review.Content,
		Status:     string(review.Status),
		CreatedAt:  review.CreatedAt,
		UpdatedAt:  review.UpdatedAt,
		ApprovedAt: review.ApprovedAt,
	}
}
