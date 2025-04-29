package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics 返回Prometheus指标监控处理器
func Metrics() gin.HandlerFunc {
	h := promhttp.Handler()
	
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
