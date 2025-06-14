package bootstrap

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/database"
	"wz-backend-go/internal/infrastructure/persistence/sql"
	"wz-backend-go/internal/interfaces/http/controller"
)

// InitCategoryService 初始化分类服务
func InitCategoryService(
	router *gin.RouterGroup,
	db database.Database,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) {
	// 创建仓储实现
	categoryRepo := sql.NewSQLCategoryRepository(db)

	// 创建应用服务
	categoryAppService := learn.NewCategoryAppService(
		categoryRepo,
		eventBus,
		unitOfWork,
	)

	// 创建控制器
	categoryController := controller.NewCategoryController(categoryAppService)

	// 注册路由
	categoryController.RegisterRoutes(router.Group("/categories"))
}
