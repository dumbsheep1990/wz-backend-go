package admin

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/admin"
	"wz-backend-go/internal/application/admin/service"
	"wz-backend-go/internal/delivery/http/internal/middleware"
)

// RegisterRoutes 注册所有admin相关路由
func RegisterRoutes(
	r *gin.Engine,
	adminService *service.AdminApplicationService,
	roleService *service.RoleApplicationService,
	menuService *service.MenuApplicationService,
	systemService *service.SystemApplicationService,
	dictionaryService *service.DictionaryApplicationService,
	operationRecordService *service.OperationRecordApplicationService,
	paramsService *service.ParamsApplicationService,
	dashboardService *admin.DashboardService,
	jwtMiddleware *middleware.JWTMiddleware,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	// 创建处理程序
	loginHandler := NewLoginHandler(adminService)
	userHandler := NewUserHandler(adminService, roleService)
	authorityHandler := NewAuthorityHandler(roleService)
	menuHandler := NewMenuHandler(menuService)
	systemHandler := NewSystemHandler(systemService)
	dictionaryHandler := NewDictionaryHandler(dictionaryService)
	operationRecordHandler := NewOperationRecordHandler(operationRecordService)
	paramsHandler := NewParamsHandler(paramsService)
	dashboardHandler := NewDashboardHandler(dashboardService)

	// API路径前缀
	apiV1 := r.Group("/api/v1")

	// 公共接口（不需要认证）
	{
		// 登录相关
		apiV1.POST("/admin/login", loginHandler.Login)
		apiV1.POST("/admin/refresh_token", loginHandler.RefreshToken)
		apiV1.GET("/admin/captcha", loginHandler.Captcha)
	}

	// 需要JWT认证的接口
	adminGroup := apiV1.Group("/admin")
	adminGroup.Use(jwtMiddleware.Handler())
	{
		// 需要额外权限控制的接口
		authGroup := adminGroup.Group("")
		authGroup.Use(casbinMiddleware.Handler())
		{
			// 用户管理
			userRoutes := authGroup.Group("/user")
			{
				userRoutes.POST("/getUserList", userHandler.GetUserList)
				userRoutes.GET("/getUserInfo", userHandler.GetUserInfo)
				userRoutes.POST("/admin_register", userHandler.CreateUser)
				userRoutes.PUT("/setUserInfo", userHandler.UpdateUser)
				userRoutes.PUT("/setSelfInfo", userHandler.UpdateSelfInfo)
				userRoutes.DELETE("/deleteUser", userHandler.DeleteUser)
				userRoutes.POST("/changePassword", userHandler.ChangePassword)
				userRoutes.POST("/resetPassword", userHandler.ResetPassword)
				userRoutes.POST("/setUserAuthority", userHandler.SetUserAuthority)
				userRoutes.POST("/setUserAuthorities", userHandler.SetUserAuthorities)
			}

			// 角色权限管理
			authorityRoutes := authGroup.Group("/authority")
			{
				authorityRoutes.POST("/getAuthorityList", authorityHandler.GetAuthorityList)
				authorityRoutes.POST("/createAuthority", authorityHandler.CreateAuthority)
				authorityRoutes.PUT("/updateAuthority", authorityHandler.UpdateAuthority)
				authorityRoutes.POST("/deleteAuthority", authorityHandler.DeleteAuthority)
				authorityRoutes.POST("/setDataAuthority", authorityHandler.SetDataAuthority)
			}

			// 按钮权限管理
			authorityBtnRoutes := authGroup.Group("/authorityBtn")
			{
				// 在原API中有以下端点，但在应用服务中缺少实现
				// authorityBtnRoutes.POST("/getAuthorityBtn", authorityHandler.GetAuthorityBtnList)
				// authorityBtnRoutes.POST("/setAuthorityBtn", authorityHandler.SetAuthorityBtn)
			}

			// Casbin策略管理
			casbinRoutes := authGroup.Group("/casbin")
			{
				casbinRoutes.POST("/getPolicyPathByAuthorityId", authorityHandler.GetCasbinPolicy)
				casbinRoutes.POST("/updateCasbin", authorityHandler.UpdateCasbinPolicy)
			}

			// 菜单管理
			menuRoutes := authGroup.Group("/menu")
			{
				menuRoutes.POST("/getMenuList", menuHandler.GetMenuList)
				menuRoutes.POST("/addMenu", menuHandler.AddMenu)
				menuRoutes.POST("/updateMenu", menuHandler.UpdateMenu)
				menuRoutes.POST("/deleteMenu", menuHandler.DeleteMenu)
				menuRoutes.POST("/getMenuAuthority", menuHandler.GetMenusByAuthority)
				menuRoutes.POST("/addMenuAuthority", menuHandler.AddMenuAuthority)
			}

			// 系统管理
			systemRoutes := authGroup.Group("/system")
			{
				systemRoutes.POST("/getSystemConfig", systemHandler.GetSystemConfig)
				systemRoutes.POST("/setSystemConfig", systemHandler.SetSystemConfig)
				systemRoutes.POST("/getServerInfo", systemHandler.GetServerInfo)
				systemRoutes.POST("/reloadSystem", systemHandler.ReloadSystem)
			}

			// 系统字典管理
			dictionaryRoutes := authGroup.Group("/sysDictionary")
			{
				dictionaryRoutes.POST("/getSysDictionaryList", dictionaryHandler.GetDictionaryList)
				dictionaryRoutes.POST("/createSysDictionary", dictionaryHandler.CreateDictionary)
				dictionaryRoutes.PUT("/updateSysDictionary", dictionaryHandler.UpdateDictionary)
				dictionaryRoutes.DELETE("/deleteSysDictionary", dictionaryHandler.DeleteDictionary)
			}

			// 系统字典详情管理
			dictionaryDetailRoutes := authGroup.Group("/sysDictionaryDetail")
			{
				dictionaryDetailRoutes.POST("/getSysDictionaryDetailList", dictionaryHandler.GetDictionaryDetailList)
				dictionaryDetailRoutes.POST("/createSysDictionaryDetail", dictionaryHandler.CreateDictionaryDetail)
				dictionaryDetailRoutes.PUT("/updateSysDictionaryDetail", dictionaryHandler.UpdateDictionaryDetail)
				dictionaryDetailRoutes.DELETE("/deleteSysDictionaryDetail", dictionaryHandler.DeleteDictionaryDetail)
			}

			// 操作日志管理
			operationRecordRoutes := authGroup.Group("/sysOperationRecord")
			{
				operationRecordRoutes.POST("/getSysOperationRecordList", operationRecordHandler.GetOperationRecordList)
				operationRecordRoutes.DELETE("/deleteSysOperationRecord", operationRecordHandler.DeleteOperationRecord)
				operationRecordRoutes.DELETE("/deleteSysOperationRecordByIds", operationRecordHandler.DeleteOperationRecordsByIds)
			}

			// 系统参数管理
			paramsRoutes := authGroup.Group("/params")
			{
				paramsRoutes.POST("/getParamsList", paramsHandler.GetParamsList)
				paramsRoutes.POST("/createParams", paramsHandler.CreateParams)
				paramsRoutes.PUT("/updateParams", paramsHandler.UpdateParams)
				paramsRoutes.DELETE("/deleteParams", paramsHandler.DeleteParams)
			}

			// 仪表盘数据
			dashboardRoutes := authGroup.Group("/dashboard")
			{
				dashboardRoutes.GET("/stats", dashboardHandler.GetDashboardStats)
				dashboardRoutes.GET("/userStats", dashboardHandler.GetUserStats)
				dashboardRoutes.GET("/contentStats", dashboardHandler.GetContentStats)
				dashboardRoutes.GET("/tradeStats", dashboardHandler.GetTradeStats)
				dashboardRoutes.GET("/siteStats", dashboardHandler.GetSiteStats)
				dashboardRoutes.GET("/communityStats", dashboardHandler.GetCommunityStats)
				dashboardRoutes.GET("/componentStats", dashboardHandler.GetComponentStats)
				dashboardRoutes.GET("/renderStats", dashboardHandler.GetRenderStats)
			}
		}
	}
}
