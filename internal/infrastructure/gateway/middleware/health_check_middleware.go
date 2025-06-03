package middleware

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthStatus 健康状态
type HealthStatus struct {
	Status     string        `json:"status"`
	Version    string        `json:"version"`
	Uptime     string        `json:"uptime"`
	GoVersion  string        `json:"go_version"`
	NumGoroutine int         `json:"num_goroutine"`
	MemoryStats MemoryStats  `json:"memory_stats"`
	Services   []ServiceHealth `json:"services,omitempty"`
}

// MemoryStats 内存使用统计
type MemoryStats struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"num_gc"`
}

// ServiceHealth 服务健康状态
type ServiceHealth struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// HealthCheckMiddleware 健康检查中间件
func HealthCheckMiddleware(version string, startTime time.Time) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只处理健康检查路径
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/healthz" {
			// 简单健康检查
			if c.Request.Method == "HEAD" {
				c.Status(http.StatusOK)
				c.Abort()
				return
			}
			
			// 详细健康检查
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			
			// 获取服务注册表的健康状态
			services := []ServiceHealth{}
			
			// 尝试从上下文中获取服务注册表
			if registry, exists := c.Get("service_registry"); exists {
				if sr, ok := registry.(interface {
					GetHealthStatus() []ServiceHealth
				}); ok {
					services = sr.GetHealthStatus()
				}
			}
			
			// 构建健康状态响应
			status := HealthStatus{
				Status:    "UP",
				Version:   version,
				Uptime:    time.Since(startTime).String(),
				GoVersion: runtime.Version(),
				NumGoroutine: runtime.NumGoroutine(),
				MemoryStats: MemoryStats{
					Alloc:      memStats.Alloc,
					TotalAlloc: memStats.TotalAlloc,
					Sys:        memStats.Sys,
					NumGC:      memStats.NumGC,
				},
				Services: services,
			}
			
			// 如果有服务不健康，则整体状态为DOWN
			for _, service := range services {
				if service.Status != "UP" {
					status.Status = "DEGRADED"
					break
				}
			}
			
			c.JSON(http.StatusOK, status)
			c.Abort()
			return
		}
		
		// 不是健康检查路径，继续处理
		c.Next()
	}
}
