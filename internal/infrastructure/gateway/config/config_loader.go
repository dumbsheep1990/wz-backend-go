package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// GatewayConfig 网关配置
type GatewayConfig struct {
	Server    ServerConfig    `yaml:"server"`
	Logging   LoggingConfig   `yaml:"logging"`
	Telemetry TelemetryConfig `yaml:"telemetry"`
	Security  SecurityConfig  `yaml:"security"`
	CORS      CORSConfig      `yaml:"cors"`
	Rate      RateConfig      `yaml:"rate"`
	Services  []ServiceConfig `yaml:"services"`
	// 数据库和Redis配置将从环境变量或单独的配置中获取
	Database  DatabaseConfig `yaml:"database,omitempty"`
	Redis     RedisConfig    `yaml:"redis,omitempty"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            int    `yaml:"port"`
	Host            string `yaml:"host"`
	ReadTimeout     int    `yaml:"readTimeout"`
	WriteTimeout    int    `yaml:"writeTimeout"`
	ShutdownTimeout int    `yaml:"shutdownTimeout"`
	Environment     string `yaml:"environment"`
	Version         string `yaml:"version,omitempty"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level   string `yaml:"level"`
	File    string `yaml:"file"`
	Format  string `yaml:"format"`
	Console bool   `yaml:"console"`
}

// TelemetryConfig 遥测配置
type TelemetryConfig struct {
	Enabled         bool          `yaml:"enabled"`
	CollectorURL    string        `yaml:"collectorURL"`
	ServiceName     string        `yaml:"serviceName"`
	ServiceVersion  string        `yaml:"serviceVersion"`
	SamplingRatio   float64       `yaml:"samplingRatio"`
	ExportBatchSize int           `yaml:"exportBatchSize"`
	ExportTimeout   time.Duration `yaml:"exportTimeout"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	JWTSecret      string   `yaml:"jwtSecret"`
	JWTExpiration  int      `yaml:"jwtExpiration"` // 分钟
	AllowedOrigins []string `yaml:"allowedOrigins"`
	AllowedHeaders []string `yaml:"allowedHeaders"`
	TrustedProxies []string `yaml:"trustedProxies"`
	RateLimit      bool     `yaml:"rateLimit"`
	XSSProtection  bool     `yaml:"xssProtection"`
	CSRFProtection bool     `yaml:"csrfProtection"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	Enabled         bool     `yaml:"enabled"`
	AllowAllOrigins bool     `yaml:"allowAllOrigins"`
	AllowOrigins    []string `yaml:"allowOrigins"`
	AllowMethods    []string `yaml:"allowMethods"`
	AllowHeaders    []string `yaml:"allowHeaders"`
	ExposeHeaders   []string `yaml:"exposeHeaders"`
	MaxAge          int      `yaml:"maxAge"`
}

// RateConfig 限流配置
type RateConfig struct {
	Enabled      bool   `yaml:"enabled"`
	MaxRequests  int    `yaml:"maxRequests"`
	IntervalSecs int    `yaml:"intervalSecs"`
	Strategy     string `yaml:"strategy"`
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	Name           string           `yaml:"name"`
	Prefix         string           `yaml:"prefix"`
	Target         string           `yaml:"target"`
	Type           string           `yaml:"type"`
	Authentication bool             `yaml:"authentication"`
	Timeout        int              `yaml:"timeout"`
	Methods        []string         `yaml:"methods"`
	Routes         []RouteDetail    `yaml:"routes,omitempty"`
	LoadBalancer   *LoadBalancer    `yaml:"loadBalancer,omitempty"`
	GRPCOptions    *GRPCOptions     `yaml:"grpcOptions,omitempty"`
}

// RouteDetail 路由详情配置
type RouteDetail struct {
	Path           string `yaml:"path"`
	Method         string `yaml:"method"`
	Authentication *bool  `yaml:"authentication,omitempty"`
	StripPath      bool   `yaml:"stripPath"`
}

// LoadBalancer 负载均衡配置
type LoadBalancer struct {
	Type        string      `yaml:"type"`
	HealthCheck HealthCheck `yaml:"healthCheck"`
}

// HealthCheck 健康检查配置
type HealthCheck struct {
	Enabled  bool   `yaml:"enabled"`
	Interval int    `yaml:"interval"`
	Path     string `yaml:"path"`
	Timeout  int    `yaml:"timeout"`
}

// GRPCOptions gRPC服务配置
type GRPCOptions struct {
	MaxRecvMsgSize int              `yaml:"maxRecvMsgSize,omitempty"`
	MaxSendMsgSize int              `yaml:"maxSendMsgSize,omitempty"`
	PackageName    string           `yaml:"packageName"`
	ServiceName    string           `yaml:"serviceName"`
	Compression    bool             `yaml:"compression,omitempty"`
	TLS            bool             `yaml:"tls,omitempty"`
	Methods        []GRPCMethodInfo `yaml:"methods"`
}

// GRPCMethodInfo gRPC方法映射信息
type GRPCMethodInfo struct {
	Name       string `yaml:"name"`
	HTTPMethod string `yaml:"httpMethod"`
	Path       string `yaml:"path"`
}

// ConfigLoader 配置加载器
type ConfigLoader struct {
	configPath string
	config     *GatewayConfig
}

// NewConfigLoader 创建新的配置加载器
func NewConfigLoader(configPath string) *ConfigLoader {
	return &ConfigLoader{
		configPath: configPath,
	}
}

// Load 加载配置
func (l *ConfigLoader) Load() (*GatewayConfig, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(l.configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", l.configPath)
	}

	// 读取配置文件
	configData, err := ioutil.ReadFile(l.configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析YAML配置
	var config GatewayConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 缓存配置
	l.config = &config
	return &config, nil
}

// GetConfig 获取配置
func (l *ConfigLoader) GetConfig() *GatewayConfig {
	return l.config
}

// LoadDefault 加载默认配置
func LoadDefault() (*GatewayConfig, error) {
	// 尝试从环境变量获取配置路径
	configPath := os.Getenv("GATEWAY_CONFIG_PATH")
	if configPath == "" {
		// 使用默认配置路径
		configPath = "configs/gateway.yaml"
	}

	// 确保配置路径是绝对路径
	if !filepath.IsAbs(configPath) {
		// 获取工作目录
		workDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("获取工作目录失败: %w", err)
		}
		configPath = filepath.Join(workDir, configPath)
	}

	log.Printf("加载配置文件: %s", configPath)
	loader := NewConfigLoader(configPath)
	return loader.Load()
}
