package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Config 用户服务配置
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	Telemetry  TelemetryConfig  `yaml:"telemetry"`
	JWT        JWTConfig        `yaml:"jwt"`
	Validation ValidationConfig `yaml:"validation"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            int           `yaml:"port"`
	Environment     string        `yaml:"environment"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// TelemetryConfig 遥测配置
type TelemetryConfig struct {
	CollectorURL string `yaml:"collectorUrl"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	ExpireTime time.Duration `yaml:"expireTime"`
}

// ValidationConfig 验证配置
type ValidationConfig struct {
	EmailVerificationTimeout time.Duration `yaml:"emailVerificationTimeout"`
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
			Port:            8081,
			Environment:     "development",
			ShutdownTimeout: 5 * time.Second,
		},
		Database: DatabaseConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
			Database: "wz_user",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		Telemetry: TelemetryConfig{
			CollectorURL: "http://localhost:4317",
		},
		JWT: JWTConfig{
			Secret:     "your-secret-key",
			ExpireTime: 24 * time.Hour,
		},
		Validation: ValidationConfig{
			EmailVerificationTimeout: 24 * time.Hour,
		},
	}
}
