package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// AdminConfig 后台管理服务配置
type AdminConfig struct {
	rest.RestConf // HTTP服务配置

	Log logx.LogConf // 日志配置

	Auth struct {
		AccessSecret string // JWT密钥
		AccessExpire int64  // JWT过期时间
	}

	Database struct {
		Driver     string // 数据库驱动
		DataSource string // 数据库连接字符串
	}

	Redis struct {
		Host     string
		Pass     string
		DB       int
		PoolSize int
	}

	// 微服务客户端配置
	UserRpc    zrpc.RpcClientConf // 用户服务
	ContentRpc zrpc.RpcClientConf // 内容服务
	TradeRpc   zrpc.RpcClientConf // 交易服务
	SearchRpc  zrpc.RpcClientConf // 搜索服务
}
