package bootstrap

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
	"wz-backend-go/internal/infrastructure/persistence/sql"
	"wz-backend-go/internal/interfaces/http/controller"
)

// InitTeacherService 初始化讲师服务
func InitTeacherService(
	router *gin.RouterGroup,
	db database.Database,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) {
	// 创建仓储实现
	teacherRepo := sql.NewSQLTeacherRepository(db)

	// 创建应用服务
	teacherAppService := learn.NewTeacherAppService(
		teacherRepo,
		eventBus,
		unitOfWork,
	)

	// 创建控制器
	teacherController := controller.NewTeacherController(teacherAppService)

	// 注册路由
	teacherController.RegisterRoutes(router.Group("/teachers"))
}
