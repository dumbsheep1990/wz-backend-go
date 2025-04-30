package registry

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HealthChecker 定义健康检查器接口
type HealthChecker interface {
	// Start 启动健康检查
	Start(ctx context.Context) error

	// Stop 停止健康检查
	Stop() error

	// RegisterCheck 注册健康检查函数
	RegisterCheck(name string, check CheckFunc)

	// DeregisterCheck 注销健康检查函数
	DeregisterCheck(name string)

	// GetStatus 获取健康状态
	GetStatus() map[string]CheckResult
}

// CheckFunc 定义健康检查函数类型
type CheckFunc func() error

// CheckResult 健康检查结果
type CheckResult struct {
	Status    CheckStatus `json:"status"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// CheckStatus 健康检查状态
type CheckStatus string

const (
	// StatusOK 表示健康状态正常
	StatusOK CheckStatus = "OK"

	// StatusFailing 表示健康状态异常
	StatusFailing CheckStatus = "FAILING"

	// StatusUnknown 表示健康状态未知
	StatusUnknown CheckStatus = "UNKNOWN"

	// 默认检查间隔
	defaultCheckInterval = 30 * time.Second
)

// HealthCheckServer 健康检查服务
type HealthCheckServer struct {
	server         *http.Server
	checks         map[string]CheckFunc
	results        map[string]CheckResult
	checkInterval  time.Duration
	mutex          sync.RWMutex
	stopChan       chan struct{}
	httpHandler    http.Handler
	checkInProcess bool
}

// NewHealthCheckServer 创建新的健康检查服务
func NewHealthCheckServer(port int, interval time.Duration) *HealthCheckServer {
	if interval <= 0 {
		interval = defaultCheckInterval
	}

	h := &HealthCheckServer{
		checks:        make(map[string]CheckFunc),
		results:       make(map[string]CheckResult),
		checkInterval: interval,
		stopChan:      make(chan struct{}),
		checkInProcess: false,
	}

	// 创建HTTP服务器
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.healthHandler)
	mux.HandleFunc("/health/live", h.livenessHandler)
	mux.HandleFunc("/health/ready", h.readinessHandler)

	h.httpHandler = mux
	h.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return h
}

// Start 开始健康检查服务
func (h *HealthCheckServer) Start(ctx context.Context) error {
	// 启动HTTP服务器
	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("健康检查服务启动失败: %v\n", err)
		}
	}()

	// 启动定时健康检查
	go h.runChecks(ctx)

	return nil
}

// Stop 停止健康检查服务
func (h *HealthCheckServer) Stop() error {
	// 发送停止信号
	close(h.stopChan)

	// 等待服务器关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.server.Shutdown(ctx)
}

// RegisterCheck 注册健康检查函数
func (h *HealthCheckServer) RegisterCheck(name string, check CheckFunc) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.checks[name] = check
	h.results[name] = CheckResult{
		Status:    StatusUnknown,
		Message:   "尚未进行检查",
		Timestamp: time.Now(),
	}
}

// DeregisterCheck 注销健康检查函数
func (h *HealthCheckServer) DeregisterCheck(name string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.checks, name)
	delete(h.results, name)
}

// GetStatus 获取健康状态
func (h *HealthCheckServer) GetStatus() map[string]CheckResult {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// 复制一份结果
	results := make(map[string]CheckResult, len(h.results))
	for k, v := range h.results {
		results[k] = v
	}

	return results
}

// runChecks 周期性执行健康检查
func (h *HealthCheckServer) runChecks(ctx context.Context) {
	ticker := time.NewTicker(h.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.performChecks()
		case <-h.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

// performChecks 执行所有健康检查
func (h *HealthCheckServer) performChecks() {
	h.mutex.Lock()
	if h.checkInProcess {
		h.mutex.Unlock()
		return
	}
	h.checkInProcess = true
	h.mutex.Unlock()

	defer func() {
		h.mutex.Lock()
		h.checkInProcess = false
		h.mutex.Unlock()
	}()

	h.mutex.RLock()
	checkFuncs := make(map[string]CheckFunc, len(h.checks))
	for name, check := range h.checks {
		checkFuncs[name] = check
	}
	h.mutex.RUnlock()

	// 执行检查
	results := make(map[string]CheckResult, len(checkFuncs))
	for name, check := range checkFuncs {
		result := CheckResult{
			Timestamp: time.Now(),
		}

		err := check()
		if err != nil {
			result.Status = StatusFailing
			result.Message = err.Error()
		} else {
			result.Status = StatusOK
			result.Message = "健康检查通过"
		}

		results[name] = result
	}

	// 更新结果
	h.mutex.Lock()
	defer h.mutex.Unlock()
	for name, result := range results {
		h.results[name] = result
	}
}

// isSystemHealthy 检查整体系统健康状态
func (h *HealthCheckServer) isSystemHealthy() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for _, result := range h.results {
		if result.Status == StatusFailing {
			return false
		}
	}
	return true
}

// healthHandler 处理健康检查请求
func (h *HealthCheckServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	
	healthy := true
	for _, result := range h.results {
		if result.Status == StatusFailing {
			healthy = false
			break
		}
	}

	if healthy {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"UP","checks":%v}`, formatChecksToJSON(h.results))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, `{"status":"DOWN","checks":%v}`, formatChecksToJSON(h.results))
	}
}

// livenessHandler 处理存活检查请求
func (h *HealthCheckServer) livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"UP"}`)
}

// readinessHandler 处理就绪检查请求
func (h *HealthCheckServer) readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if h.isSystemHealthy() {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"UP"}`)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `{"status":"DOWN"}`)
	}
}

// formatChecksToJSON 将检查结果格式化为JSON
func formatChecksToJSON(results map[string]CheckResult) string {
	jsonStr := "{"
	i := 0
	for name, result := range results {
		if i > 0 {
			jsonStr += ","
		}
		jsonStr += fmt.Sprintf(`"%s":{"status":"%s","message":"%s","timestamp":"%s"}`,
			name, result.Status, result.Message, result.Timestamp.Format(time.RFC3339))
		i++
	}
	jsonStr += "}"
	return jsonStr
}

// DefaultChecks 提供一些默认的健康检查函数

// DBHealthCheck 数据库健康检查
func DBHealthCheck(db interface{}) CheckFunc {
	return func() error {
		// 这里是示例代码，实际应根据使用的数据库类型实现相应的健康检查
		// 例如对于GORM: db.(*gorm.DB).Raw("SELECT 1").Scan(nil)
		return nil
	}
}

// RedisHealthCheck Redis健康检查
func RedisHealthCheck(client interface{}) CheckFunc {
	return func() error {
		// 这里是示例代码，实际应根据使用的Redis客户端实现相应的健康检查
		// 例如对于go-redis: client.(*redis.Client).Ping(context.Background()).Err()
		return nil
	}
}

// NacosHealthCheck Nacos健康检查
func NacosHealthCheck(registry *NacosRegistry) CheckFunc {
	return func() error {
		if registry == nil || registry.client == nil {
			return fmt.Errorf("Nacos客户端未初始化")
		}
		
		// 可以通过获取一个已知服务来检查Nacos连接
		_, err := registry.client.GetService(vo.GetServiceParam{
			ServiceName: "health-check-probe",
			GroupName:   registry.config.Group,
		})
		
		if err != nil {
			return fmt.Errorf("Nacos连接异常: %w", err)
		}
		
		return nil
	}
}
