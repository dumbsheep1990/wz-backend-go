package config

import (
	"time"
	
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	MainDomain string // 主域名，用于解析租户子域名，例如：example.com
	DB struct {
		DataSource string // 数据库连接字符串
	}
	// 服务注册与发现配置
	Registry struct {
		ServerAddr      string        // Nacos服务器地址
		ServerPort      uint64        // Nacos服务器端口
		Namespace       string        // 命名空间ID
		Group           string        // 分组
		LogDir          string        // 日志目录
		CacheDir        string        // 缓存目录
		LogLevel        string        // 日志级别
		Username        string        // 用户名(可选)
		Password        string        // 密码(可选)
		HealthCheckPort int           // 健康检查端口
		CheckInterval   time.Duration // 健康检查间隔
	}
}
