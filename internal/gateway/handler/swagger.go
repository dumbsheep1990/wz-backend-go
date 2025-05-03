package handler

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	
	"wz-backend-go/internal/gateway/config"
	"wz-backend-go/internal/gateway/docs"
)

// SwaggerHandler Swagger文档处理器
type SwaggerHandler struct {
	docGenerator *docs.APIDocGenerator
	swaggerPath  string
	lastGenTime  time.Time
	config       config.Config
}

// NewSwaggerHandler 创建新的Swagger处理器
func NewSwaggerHandler(conf config.Config, outputPath string) *SwaggerHandler {
	return &SwaggerHandler{
		docGenerator: docs.NewAPIDocGenerator(conf, outputPath),
		swaggerPath:  outputPath,
		config:       conf,
	}
}

// RegisterSwaggerRoutes 注册Swagger文档路由
func (h *SwaggerHandler) RegisterSwaggerRoutes(r *gin.Engine) {
	// 立即生成一次文档
	if err := h.docGenerator.GenerateFromServices(); err != nil {
		// 记录错误但不中断启动
		r.Use(func(c *gin.Context) {
			c.Next()
		})
	}
	
	// 注册Swagger UI路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// 注册API文档JSON端点
	r.GET("/api/docs/swagger.json", h.GetSwaggerJSON)
	
	// 注册手动生成文档端点
	r.POST("/api/docs/generate", h.GenerateSwagger)
}

// GetSwaggerJSON 获取Swagger JSON
func (h *SwaggerHandler) GetSwaggerJSON(c *gin.Context) {
	// 自动重新生成文档（如果超过一定时间）
	if time.Since(h.lastGenTime) > 10*time.Minute {
		if err := h.docGenerator.GenerateFromServices(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "生成Swagger文档失败",
				"error":   err.Error(),
			})
			return
		}
		h.lastGenTime = time.Now()
	}
	
	// 获取Swagger JSON
	jsonData, err := h.docGenerator.GetSwaggerJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取Swagger JSON失败",
			"error":   err.Error(),
		})
		return
	}
	
	c.Data(http.StatusOK, "application/json", jsonData)
}

// GenerateSwagger 手动生成Swagger文档
func (h *SwaggerHandler) GenerateSwagger(c *gin.Context) {
	// 生成文档
	if err := h.docGenerator.GenerateFromServices(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成Swagger文档失败",
			"error":   err.Error(),
		})
		return
	}
	
	h.lastGenTime = time.Now()
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Swagger文档生成成功",
		"path":    filepath.Join(h.swaggerPath, "swagger.json"),
	})
}

// Legacy Swagger handler for backward compatibility
func Swagger() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
