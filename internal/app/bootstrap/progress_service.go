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

// ProgressService 学习进度服务组装器
type ProgressService struct {
	db         *gorm.DB
	eventBus   event.EventBus
	unitOfWork database.UnitOfWork
}

// NewProgressService 创建学习进度服务组装器
func NewProgressService(db *gorm.DB, eventBus event.EventBus, unitOfWork database.UnitOfWork) *ProgressService {
	return &ProgressService{
		db:         db,
		eventBus:   eventBus,
		unitOfWork: unitOfWork,
	}
}

// Initialize 初始化学习进度服务
func (s *ProgressService) Initialize(router *gin.RouterGroup) {
	// 初始化仓储层
	progressRepo := sql.NewSQLProgressRepository(s.db)
	courseRepo := sql.NewSQLCourseRepository(s.db)
	lessonRepo := sql.NewSQLLessonRepository(s.db)
	enrollmentRepo := sql.NewSQLEnrollmentRepository(s.db)

	// 初始化应用服务层
	progressAppService := learn.NewProgressAppService(
		progressRepo,
		courseRepo,
		lessonRepo,
		enrollmentRepo,
		s.eventBus,
		s.unitOfWork,
	)

	// 初始化HTTP处理器
	progressHandler := handler.NewProgressHandler(progressAppService)

	// 注册路由
	s.registerRoutes(router, progressHandler)
}

// registerRoutes 注册学习进度相关路由
func (s *ProgressService) registerRoutes(router *gin.RouterGroup, handler *handler.ProgressHandler) {
	progressGroup := router.Group("/progress")
	{
		// 学习进度操作
		progressGroup.PUT("/lesson/:lessonId", handler.UpdateLessonProgress)    // 更新课时进度
		progressGroup.POST("/lesson/:lessonId/complete", handler.CompleteLesson) // 完成课时
		progressGroup.POST("/lesson/:lessonId/reset", handler.ResetLesson)       // 重置课时进度

		// 进度查询
		progressGroup.GET("/my", handler.GetMyProgress)                          // 获取我的学习进度
		progressGroup.GET("/course/:courseId", handler.GetCourseProgress)        // 获取课程学习进度
		progressGroup.GET("/recent", handler.GetRecentProgress)                  // 获取最近学习进度

		// 统计信息
		progressGroup.GET("/course/:courseId/stats", handler.GetCourseProgressStats) // 获取课程进度统计
		progressGroup.GET("/stats", handler.GetUserProgressStats)                    // 获取用户整体统计

		// 管理操作
		progressGroup.POST("/course/:courseId/initialize", handler.InitializeCourseProgress) // 初始化课程进度
	}
}
