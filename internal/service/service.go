package service

import (
	"time"

	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/types"
)

// LinkService 友情链接服务接口
type LinkService interface {
	CreateLink(link *domain.Link) (int64, error)
	GetLinkById(id int64) (*domain.Link, error)
	ListLinks(page, pageSize int, query map[string]interface{}) ([]*domain.Link, int64, error)
	UpdateLink(link *domain.Link) error
	DeleteLink(id int64) error
}

// SiteConfigService 站点配置服务接口
type SiteConfigService interface {
	GetSiteConfig(tenantID int64) (*domain.SiteConfig, error)
	UpdateSiteConfig(config *domain.SiteConfig) error
}

// ThemeService 主题服务接口
type ThemeService interface {
	CreateTheme(theme *domain.Theme) (int64, error)
	GetThemeById(id int64) (*domain.Theme, error)
	ListThemes(page, pageSize int, query map[string]interface{}) ([]*domain.Theme, int64, error)
	UpdateTheme(theme *domain.Theme) error
	DeleteTheme(id int64) error
	SetDefaultTheme(id int64, tenantID int64) error
	GetDefaultTheme(tenantID int64) (*domain.Theme, error)
}

// UserMessageService 用户消息服务接口
type UserMessageService interface {
	CreateMessage(message *domain.UserMessage) (int64, error)
	GetMessageById(id int64) (*domain.UserMessage, error)
	ListMessagesByUser(userID int64, page, pageSize int, query map[string]interface{}) ([]*domain.UserMessage, int64, error)
	MarkAsRead(id int64, userID int64) error
	MarkAllAsRead(userID int64) error
	DeleteMessage(id int64, userID int64) error
	CountUnread(userID int64) (int64, error)
}

// UserPointsService 用户积分服务接口
type UserPointsService interface {
	CreatePoints(points *domain.UserPoints) (int64, error)
	GetPointByID(id int64) (*domain.UserPoints, error)
	GetTotalPointsByUserID(userID int64) (int, error)
	ListPointsByUserID(userID int64, page, pageSize int64) ([]*domain.UserPoints, error)
	CountPointsByUserID(userID int64) (int64, error)
	ListPointsWithTotal(req *types.ListPointsRequest) ([]*domain.UserPoints, int64, error)
	DeletePoint(id int64) error
	ExportPointsData(req *types.ListPointsRequest) ([]byte, error)
	GetPointsStatistics() (*types.PointsStatisticsResponse, error)
	GetPointsRules() (*types.PointsRulesResponse, error)
	UpdatePointsRules(req *types.PointsRulesRequest) error
}

// UserFavoriteService 用户收藏服务接口
type UserFavoriteService interface {
	CreateFavorite(favorite *domain.UserFavorite) (int64, error)
	GetFavoriteByID(id int64) (*domain.UserFavorite, error)
	ListFavoritesWithTotal(req *types.ListFavoritesRequest) ([]*domain.UserFavorite, int64, error)
	DeleteFavorite(id int64) error
	BatchDeleteFavorites(ids []int64) error
	CheckFavorite(userID int64, itemID int64, itemType string) (bool, error)
	ExportFavoritesData(req *types.ListFavoritesRequest) ([]byte, error)

	// 统计相关方法
	GetFavoritesStatistics() (*types.FavoritesStatisticsResponse, error)
	GetHotContent() ([]*types.HotContentResponse, error)
	GetFavoritesTrend(period string) ([]*types.TrendDataResponse, error)
}

// StatisticsService 统计服务接口
type StatisticsService interface {
	RecordStatistics(data *domain.StatisticsData) error
	GetOverview(tenantID int64) (*domain.OverviewData, error)
	GetStatisticsByType(tenantID int64, type_ string, startDate, endDate time.Time) ([]*domain.StatisticsData, error)
	GetContentRanking(tenantID int64, limit int) ([]*domain.ContentRankingItem, error)
}
