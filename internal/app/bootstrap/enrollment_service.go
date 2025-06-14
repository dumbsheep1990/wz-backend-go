package bootstrap

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
	"wz-backend-go/internal/infrastructure/persistence/sql"
	"wz-backend-go/internal/interfaces/http/controller"
)

// InitEnrollmentService 初始化报名服务
func InitEnrollmentService(
	router *gin.RouterGroup,
	db database.Database,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) {
	// 创建仓储实现
	enrollmentRepo := sql.NewSQLEnrollmentRepository(db)
	courseRepo := sql.NewSQLCourseRepository(db)

	// 创建应用服务
	enrollmentAppService := learn.NewEnrollmentAppService(
		enrollmentRepo,
		courseRepo,
		eventBus,
		unitOfWork,
	)

	// 创建控制器
	enrollmentController := controller.NewEnrollmentController(enrollmentAppService)

	// 注册路由
	enrollmentController.RegisterRoutes(router.Group("/enrollments"))
}
