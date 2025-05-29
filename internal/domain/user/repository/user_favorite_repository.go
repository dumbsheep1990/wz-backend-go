package repository

import (
	"context"

	"wz-backend-go/internal/domain/user/entity"
	"wz-backend-go/internal/domain/user/valueobject"
)

// TypeDistributionItem 收藏类型分布项
type TypeDistributionItem struct {
	Type  string
	Count int64
}

// HotContentItem 热门收藏内容项
type HotContentItem struct {
	ItemID     int64
	ItemType   string
	Title      string
	Cover      string
	Count      int64
	CreateDate string
}

// TrendItem 趋势数据项
type TrendItem struct {
	Date  string
	Count int64
}

// UserFavoriteRepository 用户收藏仓储接口
type UserFavoriteRepository interface {
	// 基本操作
	Create(ctx context.Context, favorite *entity.UserFavorite) (valueobject.FavoriteID, error)
	GetByID(ctx context.Context, id valueobject.FavoriteID) (*entity.UserFavorite, error)
	Update(ctx context.Context, favorite *entity.UserFavorite) error
	DeleteByID(ctx context.Context, id valueobject.FavoriteID) error
	BatchDelete(ctx context.Context, ids []valueobject.FavoriteID) error

	// 查询操作
	ListByUserID(ctx context.Context, userID valueobject.FavoriteUserID, offset, limit int64, itemType string) ([]*entity.UserFavorite, error)
	CountByUserID(ctx context.Context, userID valueobject.FavoriteUserID, itemType string) (int64, error)
	ListWithConditions(ctx context.Context, conditions map[string]interface{}, offset, limit int64) ([]*entity.UserFavorite, error)
	CountWithConditions(ctx context.Context, conditions map[string]interface{}) (int64, error)
	CheckFavorite(ctx context.Context, userID valueobject.FavoriteUserID, itemID valueobject.FavoriteItemID, itemType valueobject.FavoriteItemType) (bool, error)

	// 统计操作
	CountUsers(ctx context.Context) (int64, error)
	CountFavorites(ctx context.Context) (int64, error)
	CountTodayFavorites(ctx context.Context) (int64, error)
	CountMonthFavorites(ctx context.Context) (int64, error)
	GroupByType(ctx context.Context) ([]*TypeDistributionItem, error)
	GetHotContent(ctx context.Context, limit int) ([]*HotContentItem, error)
	GetTrend(ctx context.Context, period string) ([]*TrendItem, error)
}
