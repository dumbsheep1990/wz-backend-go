package types

import (
	"time"
)

// 服务列表请求
type ListServicesReq struct {
}

// 服务列表响应
type ListServicesResp struct {
	Services []string `json:"services"`
}

// 获取服务实例请求
type GetServiceInstancesReq struct {
	ServiceName string `path:"serviceName"`
}

// 服务实例信息
type ServiceInstanceInfo struct {
	InstanceID    string            `json:"instance_id"`
	ServiceName   string            `json:"service_name"`
	IP            string            `json:"ip"`
	Port          int               `json:"port"`
	Status        string            `json:"status"`
	Metadata      map[string]string `json:"metadata"`
	StartTime     time.Time         `json:"start_time"`
	LastHeartbeat time.Time         `json:"last_heartbeat"`
	Healthy       bool              `json:"healthy"`
}

// 获取服务实例响应
type GetServiceInstancesResp struct {
	Instances []ServiceInstanceInfo `json:"instances"`
}

// 注册服务实例请求
type RegisterInstanceReq struct {
	ServiceName     string            `path:"serviceName"`
	IP              string            `json:"ip"`
	Port            int               `json:"port"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	HealthCheckPath string            `json:"health_check_path,omitempty"`
}

// 注册服务实例响应
type RegisterInstanceResp struct {
	InstanceID string `json:"instance_id"`
}

// 注销服务实例请求
type DeregisterInstanceReq struct {
	ServiceName string `path:"serviceName"`
	InstanceID  string `path:"instanceId"`
}

// 注销服务实例响应
type DeregisterInstanceResp struct {
	Success bool `json:"success"`
}

// 更新实例状态请求
type UpdateInstanceStatusReq struct {
	ServiceName string `path:"serviceName"`
	InstanceID  string `path:"instanceId"`
	Status      string `json:"status"`
}

// 更新实例状态响应
type UpdateInstanceStatusResp struct {
	Success bool `json:"success"`
}

// 健康状态响应
type HealthStatusResp struct {
	Services map[string]ServiceHealthInfo `json:"services"`
}

// 服务健康信息
type ServiceHealthInfo struct {
	ServiceName string                 `json:"service_name"`
	Status      string                 `json:"status"`
	Instances   int                    `json:"instances"`
	HealthyRate float64                `json:"healthy_rate"`
	Checks      map[string]CheckResult `json:"checks,omitempty"`
}

// 健康检查结果
type CheckResult struct {
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// 触发健康检查请求
type TriggerHealthCheckReq struct {
	ServiceName string `path:"serviceName,optional"`
	InstanceID  string `json:"instance_id,optional"`
}

// 触发健康检查响应
type TriggerHealthCheckResp struct {
	Success bool                   `json:"success"`
	Results map[string]CheckResult `json:"results,omitempty"`
}

// 获取服务依赖请求
type GetServiceDependenciesReq struct {
	ServiceName string `path:"serviceName"`
}

// 服务依赖关系
type ServiceDependency struct {
	Source      string `json:"source"`
	Target      string `json:"target"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// 获取服务依赖响应
type GetServiceDependenciesResp struct {
	ServiceName  string              `json:"service_name"`
	Dependencies []ServiceDependency `json:"dependencies"`
}
