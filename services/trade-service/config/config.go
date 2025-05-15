package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Config 交易服务配置
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Redis     RedisConfig     `yaml:"redis"`
	Payment   PaymentConfig   `yaml:"payment"`
	Telemetry TelemetryConfig `yaml:"telemetry"`
	JWT       JWTConfig       `yaml:"jwt"`
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

// PaymentConfig 支付配置
type PaymentConfig struct {
	AliPay    AliPayConfig    `yaml:"alipay"`
	WeChatPay WeChatPayConfig `yaml:"wechatpay"`
	PayPal    PayPalConfig    `yaml:"paypal"`
	Stripe    StripeConfig    `yaml:"stripe"`
}

// AliPayConfig 支付宝配置
type AliPayConfig struct {
	AppID      string `yaml:"appId"`
	PrivateKey string `yaml:"privateKey"`
	PublicKey  string `yaml:"publicKey"`
	NotifyURL  string `yaml:"notifyUrl"`
	ReturnURL  string `yaml:"returnUrl"`
}

// WeChatPayConfig 微信支付配置
type WeChatPayConfig struct {
	AppID     string `yaml:"appId"`
	MchID     string `yaml:"mchId"`
	MchKey    string `yaml:"mchKey"`
	NotifyURL string `yaml:"notifyUrl"`
	CertPath  string `yaml:"certPath"`
}

// PayPalConfig PayPal配置
type PayPalConfig struct {
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	Environment  string `yaml:"environment"`
}

// StripeConfig Stripe配置
type StripeConfig struct {
	SecretKey  string `yaml:"secretKey"`
	PublicKey  string `yaml:"publicKey"`
	WebhookKey string `yaml:"webhookKey"`
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

// Load 从文件加载配置
func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	return &config, err
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            8085,
			Environment:     "development",
			ShutdownTimeout: 5 * time.Second,
		},
		Database: DatabaseConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
			Database: "wz_trade",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		Payment: PaymentConfig{
			AliPay: AliPayConfig{
				AppID:      "example_appid",
				PrivateKey: "example_private_key",
				PublicKey:  "example_public_key",
				NotifyURL:  "http://localhost:8085/api/v1/payments/alipay/notify",
				ReturnURL:  "http://localhost:3000/payment-result",
			},
			WeChatPay: WeChatPayConfig{
				AppID:     "example_appid",
				MchID:     "example_mchid",
				MchKey:    "example_mchkey",
				NotifyURL: "http://localhost:8085/api/v1/payments/wechatpay/notify",
				CertPath:  "./cert/wechatpay.p12",
			},
		},
		Telemetry: TelemetryConfig{
			CollectorURL: "http://localhost:4317",
		},
		JWT: JWTConfig{
			Secret:     "your-secret-key",
			ExpireTime: 24 * time.Hour,
		},
	}
}
