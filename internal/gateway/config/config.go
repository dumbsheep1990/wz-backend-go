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
	GRPC     GRPCConfig     `yaml:"grpc"`
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

	// 认证机制
	AuthMethods    []string           `yaml:"authMethods"`   // 支持的认证方法：jwt, apikey, oauth2, basic
	DefaultAuth    string             `yaml:"defaultAuth"`   // 默认认证方法

	// API Key配置
	APIKey         APIKeyConfig       `yaml:"apiKey"`

	// OAuth2配置
	OAuth2         OAuth2Config        `yaml:"oauth2"`

	// RBAC配置
	RBAC           RBACConfig          `yaml:"rbac"`
}

// APIKeyConfig API Key配置
type APIKeyConfig struct {
	Enabled       bool                `yaml:"enabled"`
	HeaderName    string              `yaml:"headerName"`    // 自定义头名称，默认为X-API-Key
	QueryParamName string             `yaml:"queryParamName"` // 自定义查询参数名称，默认为api_key
	KeysFile      string              `yaml:"keysFile"`      // API Keys配置文件路径
}

// OAuth2Config OAuth2配置
type OAuth2Config struct {
	Enabled          bool              `yaml:"enabled"`
	ClientID         string            `yaml:"clientId"`
	ClientSecret     string            `yaml:"clientSecret"`
	AuthorizationURL string            `yaml:"authorizationUrl"`
	TokenURL         string            `yaml:"tokenUrl"`
	RedirectURI      string            `yaml:"redirectUri"`
	Scope            string            `yaml:"scope"`
	UserInfoURL      string            `yaml:"userInfoUrl"`
	ProviderType     string            `yaml:"providerType"`  // 如 'google', 'github' 等
}

// RBACConfig RBAC配置
type RBACConfig struct {
	Enabled          bool              `yaml:"enabled"`
	PolicyFile       string            `yaml:"policyFile"`    // RBAC策略文件路径
	DefaultRole      string            `yaml:"defaultRole"`   // 默认角色
	RoleClaim        string            `yaml:"roleClaim"`     // 角色声明字段名
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
	MaxRecvMsgSize    int               `yaml:"maxRecvMsgSize"`
	MaxSendMsgSize    int               `yaml:"maxSendMsgSize"`
	PackageName       string            `yaml:"packageName"`
	ServiceName       string            `yaml:"serviceName"`
	Compression       bool              `yaml:"compression"`
	TLS               bool              `yaml:"tls"`
	CertFile          string            `yaml:"certFile"`
	KeyFile           string            `yaml:"keyFile"`
	Methods           []GrpcMethodConfig `yaml:"methods"`
	DialTimeout       int               `yaml:"dialTimeout"`        // 连接超时时间(秒)
	KeepAliveTime     int               `yaml:"keepAliveTime"`      // keepalive探针间隔时间(秒)
	KeepAliveTimeout  int               `yaml:"keepAliveTimeout"`   // keepalive超时时间(秒)
	EnableReflection  bool              `yaml:"enableReflection"`   // 是否启用gRPC反射
	AutoDiscovery     bool              `yaml:"autoDiscovery"`      // 是否启用自动服务发现
	ServiceDiscovery  string            `yaml:"serviceDiscovery"`   // 服务发现类型（nacos, consul, etcd等）
}

// GrpcMethodConfig gRPC方法配置
type GrpcMethodConfig struct {
	Name       string `yaml:"name"`
	HTTPMethod string `yaml:"httpMethod"`
	Path       string `yaml:"path"`
}

// GRPCConfig 全局gRPC客户端配置
type GRPCConfig struct {
	MaxRecvMsgSize       int    `yaml:"maxRecvMsgSize"`       // 最大接收消息大小
	MaxSendMsgSize       int    `yaml:"maxSendMsgSize"`       // 最大发送消息大小
	DialTimeout          int    `yaml:"dialTimeout"`          // 连接超时(秒)，默认5
	KeepAliveTime        int    `yaml:"keepAliveTime"`        // keepalive探针间隔(秒)，默认60
	KeepAliveTimeout     int    `yaml:"keepAliveTimeout"`     // keepalive超时(秒)，默认20
	HealthCheckInterval  int    `yaml:"healthCheckInterval"`  // 健康检查间隔(秒)，默认30
	ConnectionPoolSize   int    `yaml:"connectionPoolSize"`   // 连接池大小，默认5
	EnableCompression    bool   `yaml:"enableCompression"`    // 是否启用压缩
	EnableReflection     bool   `yaml:"enableReflection"`     // 是否启用反射
	RetryCount           int    `yaml:"retryCount"`           // 重试次数，默认3
	RetryWaitTime        int    `yaml:"retryWaitTime"`        // 重试等待时间(秒)，默认1
}

// Load 从文件加载配置
func (c *Config) Load(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
