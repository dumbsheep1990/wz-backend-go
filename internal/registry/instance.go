package registry

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	// ErrServiceNotRegistered 表示服务未注册错误
	ErrServiceNotRegistered = errors.New("服务未注册")

	// ErrInstanceNotFound 表示未找到实例错误
	ErrInstanceNotFound = errors.New("未找到服务实例")
)

// InstanceManager 服务实例管理器接口
type InstanceManager interface {
	// RegisterInstance 注册服务实例
	RegisterInstance(serviceName string, instanceInfo *InstanceInfo) error

	// DeregisterInstance 注销服务实例
	DeregisterInstance(serviceName string, instanceID string) error

	// GetInstances 获取服务所有实例
	GetInstances(serviceName string) ([]*InstanceInfo, error)

	// GetInstance 获取特定实例
	GetInstance(serviceName, instanceID string) (*InstanceInfo, error)

	// UpdateInstanceStatus 更新实例状态
	UpdateInstanceStatus(serviceName, instanceID string, status InstanceStatus) error

	// StartHeartbeat 开始实例心跳
	StartHeartbeat(serviceName string, instanceID string, interval time.Duration) error

	// StopHeartbeat 停止实例心跳
	StopHeartbeat(serviceName string, instanceID string) error
}

// InstanceStatus 服务实例状态
type InstanceStatus string

const (
	// StatusUp 表示实例正常运行
	StatusUp InstanceStatus = "UP"

	// StatusDown 表示实例不可用
	StatusDown InstanceStatus = "DOWN"

	// StatusStarting 表示实例正在启动
	StatusStarting InstanceStatus = "STARTING"

	// StatusOutOfService 表示实例暂时不可用
	StatusOutOfService InstanceStatus = "OUT_OF_SERVICE"

	// StatusUnknown 表示实例状态未知
	StatusUnknown InstanceStatus = "UNKNOWN"
)

// InstanceInfo 服务实例信息
type InstanceInfo struct {
	InstanceID      string            `json:"instance_id"`
	ServiceName     string            `json:"service_name"`
	IP              string            `json:"ip"`
	Port            int               `json:"port"`
	Status          InstanceStatus    `json:"status"`
	Metadata        map[string]string `json:"metadata"`
	StartTime       time.Time         `json:"start_time"`
	LastHeartbeat   time.Time         `json:"last_heartbeat"`
	HealthCheckPath string            `json:"health_check_path"`
}

// NacosInstanceManager 基于Nacos的实例管理器
type NacosInstanceManager struct {
	registry     *NacosRegistry
	heartbeats   map[string]chan struct{}
	heartbeatsMu sync.Mutex
}

// NewNacosInstanceManager 创建新的Nacos实例管理器
func NewNacosInstanceManager(registry *NacosRegistry) *NacosInstanceManager {
	return &NacosInstanceManager{
		registry:   registry,
		heartbeats: make(map[string]chan struct{}),
	}
}

// RegisterInstance 注册服务实例
func (m *NacosInstanceManager) RegisterInstance(serviceName string, instance *InstanceInfo) error {
	if instance == nil {
		return errors.New("实例信息不能为空")
	}

	// 确保实例ID存在
	if instance.InstanceID == "" {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown-host"
		}
		instance.InstanceID = fmt.Sprintf("%s-%s-%d", serviceName, hostname, instance.Port)
	}

	// 确保状态有效
	if instance.Status == "" {
		instance.Status = StatusUp
	}

	// 初始化元数据
	if instance.Metadata == nil {
		instance.Metadata = make(map[string]string)
	}

	// 添加实例信息到元数据
	instance.Metadata["instance_id"] = instance.InstanceID
	instance.Metadata["start_time"] = instance.StartTime.Format(time.RFC3339)
	instance.Metadata["status"] = string(instance.Status)
	
	if instance.HealthCheckPath != "" {
		instance.Metadata["health_check_path"] = instance.HealthCheckPath
	}

	// 更新开始时间和心跳时间
	now := time.Now()
	if instance.StartTime.IsZero() {
		instance.StartTime = now
	}
	instance.LastHeartbeat = now

	// 调用Nacos注册服务
	return m.registry.Register(serviceName, instance.IP, instance.Port, instance.Metadata)
}

// DeregisterInstance 注销服务实例
func (m *NacosInstanceManager) DeregisterInstance(serviceName string, instanceID string) error {
	// 先获取实例信息
	instance, err := m.GetInstance(serviceName, instanceID)
	if err != nil {
		return err
	}

	// 停止心跳
	m.StopHeartbeat(serviceName, instanceID)

	// 从Nacos注销服务
	return m.registry.Deregister(serviceName, instance.IP, instance.Port)
}

// GetInstances 获取服务所有实例
func (m *NacosInstanceManager) GetInstances(serviceName string) ([]*InstanceInfo, error) {
	// 从Nacos获取服务实例
	instances, err := m.registry.GetService(serviceName)
	if err != nil {
		return nil, err
	}

	// 转换为InstanceInfo
	result := make([]*InstanceInfo, 0, len(instances))
	for _, inst := range instances {
		instanceInfo := &InstanceInfo{
			ServiceName: serviceName,
			IP:          inst.IP,
			Port:        inst.Port,
			Metadata:    inst.Metadata,
		}

		// 从元数据提取信息
		if id, ok := inst.Metadata["instance_id"]; ok {
			instanceInfo.InstanceID = id
		} else {
			instanceInfo.InstanceID = fmt.Sprintf("%s-%s-%d", serviceName, inst.IP, inst.Port)
		}

		if status, ok := inst.Metadata["status"]; ok {
			instanceInfo.Status = InstanceStatus(status)
		} else {
			if inst.Healthy {
				instanceInfo.Status = StatusUp
			} else {
				instanceInfo.Status = StatusDown
			}
		}

		if startTime, ok := inst.Metadata["start_time"]; ok {
			if t, err := time.Parse(time.RFC3339, startTime); err == nil {
				instanceInfo.StartTime = t
			}
		}

		if healthCheck, ok := inst.Metadata["health_check_path"]; ok {
			instanceInfo.HealthCheckPath = healthCheck
		}

		result = append(result, instanceInfo)
	}

	return result, nil
}

// GetInstance 获取特定实例
func (m *NacosInstanceManager) GetInstance(serviceName, instanceID string) (*InstanceInfo, error) {
	instances, err := m.GetInstances(serviceName)
	if err != nil {
		return nil, err
	}

	for _, instance := range instances {
		if instance.InstanceID == instanceID {
			return instance, nil
		}
	}

	return nil, ErrInstanceNotFound
}

// UpdateInstanceStatus 更新实例状态
func (m *NacosInstanceManager) UpdateInstanceStatus(serviceName, instanceID string, status InstanceStatus) error {
	// 先获取实例信息
	instance, err := m.GetInstance(serviceName, instanceID)
	if err != nil {
		return err
	}

	// 更新状态
	instance.Status = status
	instance.Metadata["status"] = string(status)

	// 重新注册以更新信息
	return m.RegisterInstance(serviceName, instance)
}

// StartHeartbeat 开始实例心跳
func (m *NacosInstanceManager) StartHeartbeat(serviceName string, instanceID string, interval time.Duration) error {
	// 先获取实例信息
	instance, err := m.GetInstance(serviceName, instanceID)
	if err != nil {
		return err
	}

	// 确保我们不会启动重复的心跳
	heartbeatKey := fmt.Sprintf("%s:%s", serviceName, instanceID)
	m.heartbeatsMu.Lock()
	if _, exists := m.heartbeats[heartbeatKey]; exists {
		m.heartbeatsMu.Unlock()
		return nil
	}

	// 创建停止通道
	stopChan := make(chan struct{})
	m.heartbeats[heartbeatKey] = stopChan
	m.heartbeatsMu.Unlock()

	// 启动心跳协程
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// 更新心跳时间
				instance.LastHeartbeat = time.Now()
				instance.Metadata["last_heartbeat"] = instance.LastHeartbeat.Format(time.RFC3339)

				// 重新注册以刷新心跳
				if err := m.RegisterInstance(serviceName, instance); err != nil {
					fmt.Printf("心跳更新失败 [%s-%s]: %v\n", serviceName, instanceID, err)
				}
			case <-stopChan:
				return
			}
		}
	}()

	return nil
}

// StopHeartbeat 停止实例心跳
func (m *NacosInstanceManager) StopHeartbeat(serviceName string, instanceID string) error {
	heartbeatKey := fmt.Sprintf("%s:%s", serviceName, instanceID)
	
	m.heartbeatsMu.Lock()
	defer m.heartbeatsMu.Unlock()
	
	if stopChan, exists := m.heartbeats[heartbeatKey]; exists {
		close(stopChan)
		delete(m.heartbeats, heartbeatKey)
	}
	
	return nil
}

// 工具函数

// GetLocalIP 获取本机IP地址
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}
	
	return "", errors.New("无法获取本地IP地址")
}

// GetFreePort 获取一个空闲端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	
	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetEnvOrDefault 获取环境变量，如果不存在则返回默认值
func GetEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetEnvAsInt 将环境变量解析为整数
func GetEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvAsBool 将环境变量解析为布尔值
func GetEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
