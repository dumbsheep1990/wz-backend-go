package bootstrap

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
	"wz-backend-go/internal/infrastructure/persistence/sql"
	"wz-backend-go/internal/interfaces/http/controller"
)

// InitLearningService 初始化学习模块主服务
func InitLearningService(
	router *gin.RouterGroup,
	db database.Database,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) {
	// 创建学习模块路由组
	learningGroup := router.Group("/learning")

	// 初始化各个子服务
	InitTeacherService(learningGroup, db, eventBus, unitOfWork)
	InitCourseService(learningGroup, db, eventBus, unitOfWork)
	InitCategoryService(learningGroup, db, eventBus, unitOfWork)
	InitEnrollmentService(learningGroup, db, eventBus, unitOfWork)
	InitCertificateService(learningGroup, db, eventBus, unitOfWork)
	InitReviewService(learningGroup, db, eventBus, unitOfWork)
	InitProgressService(learningGroup, db, eventBus, unitOfWork)
	InitChapterLessonService(learningGroup, db, eventBus, unitOfWork)
}

// InitReviewService 初始化评价服务
func InitReviewService(
	router *gin.RouterGroup,
	db database.Database,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) {
	reviewService := NewReviewService(db.GetDB(), eventBus, unitOfWork)
	reviewService.Initialize(router)
}

// InitProgressService 初始化学习进度服务
func InitProgressService(
	router *gin.RouterGroup,
	db database.Database,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) {
	progressService := NewProgressService(db.GetDB(), eventBus, unitOfWork)
	progressService.Initialize(router)
}
