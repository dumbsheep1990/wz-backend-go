package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 后台管理服务配置结构
type Config struct {
	rest.RestConf              // REST服务配置
	Log           logx.LogConf // 日志配置
	Auth          struct {     // 认证相关配置
		AccessSecret string // JWT密钥
		AccessExpire int64  // JWT过期时间
	}
	DB struct { // 数据库配置
		Driver     string // 数据库驱动
		DataSource string // 数据源
	}
	Cache cache.CacheConf // 缓存配置

	// 微服务客户端配置
	UserRPC         zrpc.RpcClientConf // 用户服务RPC
	ContentRPC      zrpc.RpcClientConf // 内容服务RPC
	InteractionRPC  zrpc.RpcClientConf // 交互服务RPC
	AIRPC           zrpc.RpcClientConf // AI服务RPC
	SearchRPC       zrpc.RpcClientConf // 搜索服务RPC
	NotificationRPC zrpc.RpcClientConf // 通知服务RPC
	TradeRPC        zrpc.RpcClientConf // 交易服务RPC
	FileRPC         zrpc.RpcClientConf // 文件服务RPC
	StatisticsRPC   zrpc.RpcClientConf // 统计服务RPC
	AdRPC           zrpc.RpcClientConf // 广告服务RPC
	RecommendRPC    zrpc.RpcClientConf // 推荐服务RPC
}
