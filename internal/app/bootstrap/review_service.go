package bootstrap

import (
	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/delivery/http/handler"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
	"wz-backend-go/internal/infrastructure/persistence/sql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ReviewService 评价服务组装器
type ReviewService struct {
	db         *gorm.DB
	eventBus   event.EventBus
	unitOfWork database.UnitOfWork
}

// NewReviewService 创建评价服务组装器
func NewReviewService(db *gorm.DB, eventBus event.EventBus, unitOfWork database.UnitOfWork) *ReviewService {
	return &ReviewService{
		db:         db,
		eventBus:   eventBus,
		unitOfWork: unitOfWork,
	}
}

// Initialize 初始化评价服务
func (s *ReviewService) Initialize(router *gin.RouterGroup) {
	// 初始化仓储层
	reviewRepo := sql.NewSQLReviewRepository(s.db)
	courseRepo := sql.NewSQLCourseRepository(s.db)
	enrollmentRepo := sql.NewSQLEnrollmentRepository(s.db)

	// 初始化应用服务层
	reviewAppService := learn.NewReviewAppService(
		reviewRepo,
		courseRepo,
		enrollmentRepo,
		s.eventBus,
		s.unitOfWork,
	)

	// 初始化HTTP处理器
	reviewHandler := handler.NewReviewHandler(reviewAppService)

	// 注册路由
	s.registerRoutes(router, reviewHandler)
}

// registerRoutes 注册评价相关路由
func (s *ReviewService) registerRoutes(router *gin.RouterGroup, handler *handler.ReviewHandler) {
	reviewGroup := router.Group("/reviews")
	{
		// 用户评价操作
		reviewGroup.POST("", handler.CreateReview)           // 创建评价
		reviewGroup.PUT("/:id", handler.UpdateReview)        // 更新评价
		reviewGroup.DELETE("/:id", handler.DeleteReview)     // 删除评价
		reviewGroup.GET("/my", handler.GetMyReviews)         // 获取我的评价

		// 课程评价查询
		reviewGroup.GET("/course/:courseId", handler.GetCourseReviews)     // 获取课程评价
		reviewGroup.GET("/course/:courseId/stats", handler.GetRatingStats) // 获取课程评分统计

		// 管理员操作
		adminGroup := reviewGroup.Group("/admin")
		{
			adminGroup.GET("/pending", handler.GetPendingReviews) // 获取待审核评价
			adminGroup.PUT("/:id/approve", handler.ApproveReview) // 审核通过
			adminGroup.PUT("/:id/reject", handler.RejectReview)   // 审核拒绝
		}
	}
}
