package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
	"log"

	"wz-backend-go/internal/gateway/config"
	"wz-backend-go/internal/telemetry"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
)

var (
	// 熔断器实例映射
	circuitBreakers   = make(map[string]*gobreaker.CircuitBreaker)
	circuitBreakersMu sync.RWMutex
)

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	// 熔断器名称
	Name string
	// 最大连续失败次数，超过此值触发熔断
	MaxRequests uint32
	// 熔断后恢复尝试的时间间隔
	Interval time.Duration
	// 熔断超时时间，超过此时间后进入半开状态
	Timeout time.Duration
}

// CircuitBreaker 返回一个熔断中间件
func CircuitBreaker(serviceConfig config.ServiceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 根据服务名称或路径确定熔断器名称
		cbName := getCircuitBreakerName(c)
		
		// 获取或创建熔断器
		cb := getCircuitBreaker(cbName, CircuitBreakerConfig{
			Name:        cbName,
			MaxRequests: serviceConfig.CircuitBreaker.MaxRequests,
			Interval:    time.Duration(serviceConfig.CircuitBreaker.IntervalSecs) * time.Second,
			Timeout:     time.Duration(serviceConfig.CircuitBreaker.TimeoutSecs) * time.Second,
		})
		
		// 使用熔断器执行请求
		_, err := cb.Execute(func() (interface{}, error) {
			// 记录原始ResponseWriter，以便稍后检查状态码
			writer := c.Writer
			
			// 执行请求链
			c.Next()
			
			// 如果响应状态码表示服务器错误，返回错误
			if writer.Status() >= 500 {
				return nil, &ServiceError{Code: writer.Status(), Message: "服务器错误"}
			}
			
			return nil, nil
		})
		
		// 如果熔断器打开或执行出错
		if err != nil {
			// 如果之前已经写入了响应，不做额外处理
			if c.Writer.Written() {
				return
			}
			
			// 根据错误类型返回合适的响应
			if err == gobreaker.ErrOpenState {
				// 记录熔断事件
				telemetry.RecordEvent("circuit_breaker_open", map[string]string{
					"service": cbName,
				})
				
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"code":    503,
					"message": "服务暂时不可用，请稍后再试",
				})
			} else {
				// 其他错误，如超时等
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务内部错误",
				})
			}
		}
	}
}

// ServiceError 服务错误定义
type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

// getCircuitBreakerName 获取熔断器名称
func getCircuitBreakerName(c *gin.Context) string {
	// 优先使用设置的服务名称
	service, exists := c.Get("targetService")
	if exists && service != nil {
		return service.(string)
	}
	
	// 回退到使用请求路径前缀
	path := c.Request.URL.Path
	// 简化路径为前两级作为服务标识
	parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 3)
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	
	return path
}

// getCircuitBreaker 获取或创建熔断器
func getCircuitBreaker(name string, config CircuitBreakerConfig) *gobreaker.CircuitBreaker {
	circuitBreakersMu.RLock()
	cb, exists := circuitBreakers[name]
	circuitBreakersMu.RUnlock()
	
	if !exists {
		// 创建新的熔断器
		settings := gobreaker.Settings{
			Name:        config.Name,
			MaxRequests: config.MaxRequests,
			Interval:    config.Interval,
			Timeout:     config.Timeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				// 连续失败次数超过阈值，或者失败率高于50%且至少有5次请求
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.ConsecutiveFailures >= int64(config.MaxRequests) ||
					(counts.Requests >= 5 && failureRatio >= 0.5)
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				// 记录状态变化
				telemetry.RecordEvent("circuit_breaker_state_change", map[string]string{
					"service": name,
					"from":    from.String(),
					"to":      to.String(),
				})
				
				// 熔断器状态变化时记录日志
				log.Printf("熔断器 %s 从 %s 变为 %s", name, from.String(), to.String())
			},
		}
		
		cb = gobreaker.NewCircuitBreaker(settings)
		
		circuitBreakersMu.Lock()
		circuitBreakers[name] = cb
		circuitBreakersMu.Unlock()
	}
	
	return cb
}
