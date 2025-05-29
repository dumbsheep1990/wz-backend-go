package bootstrap

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/user/service"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
	"wz-backend-go/internal/infrastructure/persistence/sql"
	"wz-backend-go/internal/interfaces/http/controller"
)

// InitUserFavoriteService 初始化用户收藏服务
func InitUserFavoriteService(
	router *gin.RouterGroup,
	db database.Database,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) {
	// 创建仓储实现
	favoriteRepo := sql.NewSQLUserFavoriteRepository(db)

	// 创建应用服务
	favoriteAppService := service.NewUserFavoriteApplicationService(
		favoriteRepo,
		eventBus,
		unitOfWork,
	)

	// 创建控制器
	favoriteController := controller.NewUserFavoriteController(favoriteAppService)

	// 注册路由
	favoriteController.RegisterRoutes(router)
}
