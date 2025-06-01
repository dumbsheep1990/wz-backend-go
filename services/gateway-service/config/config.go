package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Config 是网关服务的配置结构
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Services  []ServiceConfig `yaml:"services"`
	Telemetry TelemetryConfig `yaml:"telemetry"`
	Database  DatabaseConfig  `yaml:"database"`
	Redis     RedisConfig     `yaml:"redis"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            int           `yaml:"port"`
	Environment     string        `yaml:"environment"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

// ServiceConfig 微服务配置
type ServiceConfig struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	RequireAuth bool   `yaml:"requireAuth"`
}

// TelemetryConfig 遥测配置
type TelemetryConfig struct {
	CollectorURL string `yaml:"collectorUrl"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Name            string        `yaml:"name"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"poolSize"`
}

// Load 从文件加载配置
func (c *Config) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Port:            8080,
			Environment:     "development",
			ShutdownTimeout: 5 * time.Second,
		},
		Services: []ServiceConfig{
			{
				Name:        "user-service",
				URL:         "http://localhost:8081",
				RequireAuth: true,
			},
			{
				Name:        "content-service",
				URL:         "http://localhost:8082",
				RequireAuth: true,
			},
			{
				Name:        "file-service",
				URL:         "http://localhost:8083",
				RequireAuth: true,
			},
			{
				Name:        "interaction-service",
				URL:         "http://localhost:8085",
				RequireAuth: true,
			},
			{
				Name:        "community-service",
				URL:         "http://localhost:8084",
				RequireAuth: true,
			},
			{
				Name:        "admin-service",
				URL:         "http://localhost:8086",
				RequireAuth: true,
			},
			{
				Name:        "render-service",
				URL:         "http://localhost:8087",
				RequireAuth: false,
			},
		},
		Telemetry: TelemetryConfig{
			CollectorURL: "http://localhost:4317",
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            3306,
			User:            "wanzhiuser",
			Password:        "wanzhipass",
			Name:            "wanzhi_gateway",
			MaxOpenConns:    50,
			MaxIdleConns:    10,
			ConnMaxLifetime: 5 * time.Minute,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
			PoolSize: 100,
		},
	}
} 