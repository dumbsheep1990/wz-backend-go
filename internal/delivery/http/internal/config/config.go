package config

import "github.com/zeromicro/go-zero/rest"

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
}
