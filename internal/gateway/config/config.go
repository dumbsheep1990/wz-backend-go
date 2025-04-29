package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config 网关配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Logging  LoggingConfig  `yaml:"logging"`
	Services []ServiceConfig `yaml:"services"`
	Security SecurityConfig `yaml:"security"`
	Cors     CorsConfig     `yaml:"cors"`
	Rate     RateConfig     `yaml:"rate"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"readTimeout"`  // 秒
	WriteTimeout int    `yaml:"writeTimeout"` // 秒
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level   string `yaml:"level"`
	File    string `yaml:"file"`
	Format  string `yaml:"format"`
	Console bool   `yaml:"console"`
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	Name          string               `yaml:"name"`
	Prefix        string               `yaml:"prefix"`
	Target        string               `yaml:"target"`
	Type          string               `yaml:"type"` // http, grpc
	Authentication bool                 `yaml:"authentication"`
	Routes        []RouteConfig        `yaml:"routes"`
	Methods       []string             `yaml:"methods"`
	Timeout       int                  `yaml:"timeout"` // 秒
	LoadBalancer  LoadBalancerConfig   `yaml:"loadBalancer"`
	GrpcOptions   GrpcServiceOptions   `yaml:"grpcOptions,omitempty"`
}

// RouteConfig 路由配置
type RouteConfig struct {
	Path       string            `yaml:"path"`
	Method     string            `yaml:"method"`
	Target     string            `yaml:"target"`
	StripPath  bool              `yaml:"stripPath"`
	Headers    map[string]string `yaml:"headers"`
	QueryParams map[string]string `yaml:"queryParams"`
}

// LoadBalancerConfig 负载均衡配置
type LoadBalancerConfig struct {
	Type      string `yaml:"type"` // round-robin, random, least-conn
	HealthCheck HealthCheckConfig `yaml:"healthCheck"`
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Interval int    `yaml:"interval"` // 秒
	Path     string `yaml:"path"`
	Timeout  int    `yaml:"timeout"`  // 秒
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	JwtSecret      string            `yaml:"jwtSecret"`
	JwtExpiration  int               `yaml:"jwtExpiration"` // 分钟
	AllowedOrigins []string          `yaml:"allowedOrigins"`
	AllowedHeaders []string          `yaml:"allowedHeaders"`
	TrustedProxies []string          `yaml:"trustedProxies"`
	RateLimit      bool              `yaml:"rateLimit"`
	XSSProtection bool              `yaml:"xssProtection"`
	CSRFProtection bool              `yaml:"csrfProtection"`
}

// CorsConfig CORS配置
type CorsConfig struct {
	Enabled         bool     `yaml:"enabled"`
	AllowAllOrigins bool     `yaml:"allowAllOrigins"`
	AllowOrigins    []string `yaml:"allowOrigins"`
	AllowMethods    []string `yaml:"allowMethods"`
	AllowHeaders    []string `yaml:"allowHeaders"`
	ExposeHeaders   []string `yaml:"exposeHeaders"`
	MaxAge          int      `yaml:"maxAge"` // 秒
}

// RateConfig 限流配置
type RateConfig struct {
	Enabled      bool   `yaml:"enabled"`
	MaxRequests  int    `yaml:"maxRequests"`
	IntervalSecs int    `yaml:"intervalSecs"`
	Strategy     string `yaml:"strategy"` // ip, user, custom
}

// GrpcServiceOptions gRPC服务特定配置
type GrpcServiceOptions struct {
	MaxRecvMsgSize int               `yaml:"maxRecvMsgSize"`
	MaxSendMsgSize int               `yaml:"maxSendMsgSize"`
	PackageName    string            `yaml:"packageName"`
	ServiceName    string            `yaml:"serviceName"`
	Compression    bool              `yaml:"compression"`
	TLS            bool              `yaml:"tls"`
	CertFile       string            `yaml:"certFile"`
	KeyFile        string            `yaml:"keyFile"`
	Methods        []GrpcMethodConfig `yaml:"methods"`
}

// GrpcMethodConfig gRPC方法配置
type GrpcMethodConfig struct {
	Name       string `yaml:"name"`
	HTTPMethod string `yaml:"httpMethod"`
	Path       string `yaml:"path"`
}

// Load 从文件加载配置
func (c *Config) Load(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
