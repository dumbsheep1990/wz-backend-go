package bootstrap

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
	"wz-backend-go/internal/infrastructure/persistence/sql"
	"wz-backend-go/internal/interfaces/http/controller"
)

// InitChapterLessonService 初始化章节课程服务
func InitChapterLessonService(
	router *gin.RouterGroup,
	db database.Database,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) {
	// 创建仓储实现
	chapterRepo := sql.NewSQLChapterRepository(db)
	lessonRepo := sql.NewSQLLessonRepository(db)
	courseRepo := sql.NewSQLCourseRepository(db)

	// 创建应用服务
	chapterLessonAppService := learn.NewChapterLessonAppService(
		chapterRepo,
		lessonRepo,
		courseRepo,
		eventBus,
		unitOfWork,
	)

	// 创建控制器
	chapterLessonController := controller.NewChapterLessonController(chapterLessonAppService)

	// 注册路由
	chapterLessonController.RegisterRoutes(router.Group("/chapters"))
}
