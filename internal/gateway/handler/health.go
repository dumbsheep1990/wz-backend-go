package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// 健康状态信息
type HealthInfo struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
	Timestamp time.Time `json:"timestamp"`
	GoVersion string    `json:"go_version"`
	Memory    MemStats  `json:"memory"`
}

// 内存统计信息
type MemStats struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"num_gc"`
}

// 存储服务启动时间
var startTime = time.Now()

// HealthCheck 返回健康检查处理器
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取内存统计信息
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		
		// 计算运行时间
		uptime := time.Since(startTime).String()
		
		// 构建健康信息
		health := HealthInfo{
			Status:    "UP",
			Version:   "1.0.0", // 可从配置或构建信息中获取
			Uptime:    uptime,
			Timestamp: time.Now(),
			GoVersion: runtime.Version(),
			Memory: MemStats{
				Alloc:      memStats.Alloc,
				TotalAlloc: memStats.TotalAlloc,
				Sys:        memStats.Sys,
				NumGC:      memStats.NumGC,
			},
		}
		
		c.JSON(http.StatusOK, health)
	}
}
