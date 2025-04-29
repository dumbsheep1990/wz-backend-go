package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger 返回Swagger文档处理器
func Swagger() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
