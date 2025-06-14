package bootstrap

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
	"wz-backend-go/internal/infrastructure/persistence/sql"
	"wz-backend-go/internal/interfaces/http/controller"
)

// InitCourseService 初始化课程服务
func InitCourseService(
	router *gin.RouterGroup,
	db database.Database,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) {
	// 创建仓储实现
	courseRepo := sql.NewSQLCourseRepository(db)
	categoryRepo := sql.NewSQLCategoryRepository(db)
	teacherRepo := sql.NewSQLTeacherRepository(db)

	// 创建应用服务
	courseAppService := learn.NewCourseAppService(
		courseRepo,
		categoryRepo,
		teacherRepo,
		eventBus,
		unitOfWork,
	)

	// 创建控制器
	courseController := controller.NewCourseController(courseAppService)

	// 注册路由
	courseController.RegisterRoutes(router.Group("/courses"))
}
