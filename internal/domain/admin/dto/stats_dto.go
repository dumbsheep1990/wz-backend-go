package dto

import (
	"time"
)

// UserStats 用户服务统计数据DTO
type UserStats struct {
	TotalUsers     int64     `json:"totalUsers"`
	ActiveUsers    int64     `json:"activeUsers"`
	NewUsers       int64     `json:"newUsers"`
	VerifiedUsers  int64     `json:"verifiedUsers"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}

// ContentStats 内容服务统计数据DTO
type ContentStats struct {
	TotalContent     int64     `json:"totalContent"`
	PublishedContent int64     `json:"publishedContent"`
	DraftContent     int64     `json:"draftContent"`
	Categories       int64     `json:"categories"`
	Tags             int64     `json:"tags"`
	LastUpdateTime   time.Time `json:"lastUpdateTime"`
}

// TradeStats 交易服务统计数据DTO
type TradeStats struct {
	TotalOrders     int64     `json:"totalOrders"`
	PendingOrders   int64     `json:"pendingOrders"`
	CompletedOrders int64     `json:"completedOrders"`
	TotalRevenue    float64   `json:"totalRevenue"`
	DailyRevenue    float64   `json:"dailyRevenue"`
	LastUpdateTime  time.Time `json:"lastUpdateTime"`
}

// CommunityStats 社区服务统计数据DTO
type CommunityStats struct {
	TotalCommunities  int64     `json:"totalCommunities"`
	ActiveCommunities int64     `json:"activeCommunities"`
	TotalGroups       int64     `json:"totalGroups"`
	TotalPosts        int64     `json:"totalPosts"`
	TotalComments     int64     `json:"totalComments"`
	LastUpdateTime    time.Time `json:"lastUpdateTime"`
}

// SiteStats 站点服务统计数据DTO
type SiteStats struct {
	TotalSites     int64     `json:"totalSites"`
	ActiveSites    int64     `json:"activeSites"`
	TotalPages     int64     `json:"totalPages"`
	PublishedPages int64     `json:"publishedPages"`
	TotalTemplates int64     `json:"totalTemplates"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}

// ComponentStats 组件服务统计数据DTO
type ComponentStats struct {
	TotalComponents     int64     `json:"totalComponents"`
	PublishedComponents int64     `json:"publishedComponents"`
	Categories          int64     `json:"categories"`
	LastUpdateTime      time.Time `json:"lastUpdateTime"`
}

// RenderStats 渲染服务统计数据DTO
type RenderStats struct {
	TotalRenders      int64     `json:"totalRenders"`
	CacheHitRate      float64   `json:"cacheHitRate"`
	AverageRenderTime float64   `json:"averageRenderTime"`
	ErrorRate         float64   `json:"errorRate"`
	LastUpdateTime    time.Time `json:"lastUpdateTime"`
}
