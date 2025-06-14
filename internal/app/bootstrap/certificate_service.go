package bootstrap

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
	"wz-backend-go/internal/infrastructure/persistence/sql"
	"wz-backend-go/internal/interfaces/http/controller"
)

// InitCertificateService 初始化证书服务
func InitCertificateService(
	router *gin.RouterGroup,
	db database.Database,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) {
	// 创建仓储实现
	certificateRepo := sql.NewSQLCertificateRepository(db)
	enrollmentRepo := sql.NewSQLEnrollmentRepository(db)
	courseRepo := sql.NewSQLCourseRepository(db)

	// 创建应用服务
	certificateAppService := learn.NewCertificateAppService(
		certificateRepo,
		enrollmentRepo,
		courseRepo,
		eventBus,
		unitOfWork,
	)

	// 创建控制器
	certificateController := controller.NewCertificateController(certificateAppService)

	// 注册路由
	certificateController.RegisterRoutes(router.Group("/certificates"))
}
