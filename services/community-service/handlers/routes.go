package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/wz-project/wz-backend-go/services/community-service/models"
	"github.com/wz-project/wz-backend-go/services/community-service/service"
)

// RegisterSimilarApplicationRoutes 注册申请建同相关的路由
func RegisterSimilarApplicationRoutes(router *gin.Engine, db *gorm.DB) {
	// 创建服务和处理器
	similarAppService := service.NewSimilarApplicationService(db)
	similarAppHandler := NewSimilarApplicationHandler(similarAppService)

	// API v1 路由组
	v1 := router.Group("/api/v1")

	// 身份验证中间件
	authMiddleware := func(c *gin.Context) {
		// 检查用户ID是否存在
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(401, gin.H{"error": "需要身份验证"})
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}

	// 同乡/同事/同学申请路由
	similarApps := v1.Group("/similar-applications")
	{
		// 公开路由 - 不需要身份验证
		similarApps.GET("", func(c *gin.Context) {
			page := c.DefaultQuery("page", "1")
			pageSize := c.DefaultQuery("page_size", "10")
			appType := c.DefaultQuery("type", "")
			status := c.DefaultQuery("status", "")

			// 构建过滤条件
			filters := make(map[string]interface{})
			if appType != "" {
				filters["application_type"] = appType
			}
			if status != "" {
				filters["status"] = status
			}

			// 转换为整数
			pageInt := 1
			if p, err := strconv.Atoi(page); err == nil && p > 0 {
				pageInt = p
			}

			pageSizeInt := 10
			if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
				pageSizeInt = ps
			}

			applications, total, err := similarAppService.ListApplications(pageInt, pageSizeInt, filters)
			if err != nil {
				c.JSON(500, gin.H{"error": "获取申请列表失败: " + err.Error()})
				return
			}

			c.JSON(200, gin.H{
				"data":       applications,
				"total":      total,
				"page":       pageInt,
				"page_size":  pageSizeInt,
				"total_page": (total + int64(pageSizeInt) - 1) / int64(pageSizeInt),
			})
		})

		similarApps.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			application, err := similarAppService.GetApplicationByID(id)
			if err != nil {
				c.JSON(404, gin.H{"error": "找不到申请: " + err.Error()})
				return
			}
			c.JSON(200, application)
		})

		// 需要身份验证的路由
		authSimilarApps := similarApps.Group("")
		authSimilarApps.Use(authMiddleware)
		{
			authSimilarApps.POST("", func(c *gin.Context) {
				var application models.SimilarApplication
				if err := c.ShouldBindJSON(&application); err != nil {
					c.JSON(400, gin.H{"error": "无效的请求数据: " + err.Error()})
					return
				}

				// 设置用户ID
				userID, _ := c.Get("user_id")
				application.UserID = userID.(string)

				if err := similarAppService.CreateApplication(&application); err != nil {
					c.JSON(500, gin.H{"error": "创建申请失败: " + err.Error()})
					return
				}

				c.JSON(201, gin.H{
					"message": "申请创建成功",
					"data":    application,
				})
			})

			authSimilarApps.PUT("/:id", func(c *gin.Context) {
				id := c.Param("id")
				var updates map[string]interface{}
				if err := c.ShouldBindJSON(&updates); err != nil {
					c.JSON(400, gin.H{"error": "无效的请求数据: " + err.Error()})
					return
				}

				if err := similarAppService.UpdateApplication(id, updates); err != nil {
					c.JSON(500, gin.H{"error": "更新申请失败: " + err.Error()})
					return
				}

				c.JSON(200, gin.H{"message": "申请更新成功"})
			})

			authSimilarApps.PATCH("/:id/status", func(c *gin.Context) {
				id := c.Param("id")
				var statusUpdate struct {
					Status string `json:"status"`
				}
				if err := c.ShouldBindJSON(&statusUpdate); err != nil {
					c.JSON(400, gin.H{"error": "无效的请求数据: " + err.Error()})
					return
				}

				if err := similarAppService.UpdateApplicationStatus(id, statusUpdate.Status); err != nil {
					c.JSON(500, gin.H{"error": "更新状态失败: " + err.Error()})
					return
				}

				c.JSON(200, gin.H{"message": "申请状态更新成功"})
			})

			authSimilarApps.DELETE("/:id", func(c *gin.Context) {
				id := c.Param("id")
				if err := similarAppService.DeleteApplication(id); err != nil {
					c.JSON(500, gin.H{"error": "删除申请失败: " + err.Error()})
					return
				}

				c.JSON(200, gin.H{"message": "申请删除成功"})
			})
		}
	}
}
