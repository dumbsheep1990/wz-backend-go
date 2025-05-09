package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

// Config Admin服务配置结构
type Config struct {
	rest.RestConf          // REST服务配置
	Auth          struct { // 认证相关配置
		AccessSecret string // JWT密钥
		AccessExpire int64  // JWT过期时间
	}
	DB struct { // 数据库配置
		DataSource string // 数据源
	}
	Cache cache.CacheConf // 缓存配置

	// 微服务客户端配置
	UserServiceRPC         string // 用户服务RPC地址
	ContentServiceRPC      string // 内容服务RPC地址
	InteractionServiceRPC  string // 交互服务RPC地址
	AIServiceRPC           string // AI服务RPC地址
	SearchServiceRPC       string // 搜索服务RPC地址
	NotificationServiceRPC string // 通知服务RPC地址
	TradeServiceRPC        string // 交易服务RPC地址
	FileServiceRPC         string // 文件服务RPC地址
	StatisticsServiceRPC   string // 统计服务RPC地址
	AdServiceRPC           string // 广告服务RPC地址
	RecommendServiceRPC    string // 推荐服务RPC地址
}
