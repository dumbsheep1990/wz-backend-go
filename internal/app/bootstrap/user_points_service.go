package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/application/user/service"
	"wz-backend-go/internal/infrastructure/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
	"wz-backend-go/internal/infrastructure/persistence/sql"
	"wz-backend-go/internal/interfaces/http/controller"
)

// InitUserPointsService 初始化用户积分服务
func InitUserPointsService(db *sqlx.DB, router *gin.RouterGroup, userRepo service.UserRepository) {
	// 创建事件总线
	eventBus := event.NewSimpleEventBus()

	// 创建工作单元
	unitOfWork := database.NewSQLUnitOfWork(db)

	// 创建仓储实现
	userPointsRepo := sql.NewUserPointsRepository(db)
	pointsRulesRepo := sql.NewPointsRulesRepository(db)

	// 创建应用服务
	userPointsService := service.NewUserPointsApplicationService(
		userPointsRepo,
		pointsRulesRepo,
		userRepo,
		eventBus,
		unitOfWork,
	)

	// 创建控制器
	userPointsController := controller.NewUserPointsController(userPointsService)

	// 注册路由
	userPointsController.Register(router)
}
