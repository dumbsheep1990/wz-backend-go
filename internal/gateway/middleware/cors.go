package middleware

import (
	"time"
	
	"wz-backend-go/internal/gateway/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 返回一个CORS处理中间件
func Cors(corsConfig config.CorsConfig) gin.HandlerFunc {
	corsOptions := cors.Config{
		AllowMethods:     corsConfig.AllowMethods,
		AllowHeaders:     corsConfig.AllowHeaders,
		ExposeHeaders:    corsConfig.ExposeHeaders,
		AllowCredentials: true,
		MaxAge:           time.Duration(corsConfig.MaxAge) * time.Second,
	}

	// 设置允许的源
	if corsConfig.AllowAllOrigins {
		corsOptions.AllowAllOrigins = true
	} else {
		corsOptions.AllowOrigins = corsConfig.AllowOrigins
	}

	return cors.New(corsOptions)
}
