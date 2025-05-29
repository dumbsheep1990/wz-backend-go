package dto

import "time"

// CreateFavoriteRequest 创建收藏请求DTO
type CreateFavoriteRequest struct {
	UserID   int64  `json:"user_id"`   // 用户ID
	ItemID   int64  `json:"item_id"`   // 内容项ID
	ItemType string `json:"item_type"` // 内容项类型
	Title    string `json:"title"`     // 标题
	Cover    string `json:"cover"`     // 封面图片URL
	Summary  string `json:"summary"`   // 摘要
	URL      string `json:"url"`       // 内容URL
	Remark   string `json:"remark"`    // 备注
	TenantID int64  `json:"tenant_id"` // 租户ID
}

// FavoriteDTO 收藏记录DTO
type FavoriteDTO struct {
	ID        int64     `json:"id"`         // 收藏ID
	UserID    int64     `json:"user_id"`    // 用户ID
	Username  string    `json:"username"`   // 用户名
	ItemID    int64     `json:"item_id"`    // 内容项ID
	ItemType  string    `json:"item_type"`  // 内容项类型
	Title     string    `json:"title"`      // 标题
	Cover     string    `json:"cover"`      // 封面图片URL
	Summary   string    `json:"summary"`    // 摘要
	URL       string    `json:"url"`        // 内容URL
	Remark    string    `json:"remark"`     // 备注
	TenantID  int64     `json:"tenant_id"`  // 租户ID
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// UpdateFavoriteRequest 更新收藏请求DTO
type UpdateFavoriteRequest struct {
	ID      int64  `json:"id"`      // 收藏ID
	Cover   string `json:"cover"`   // 封面图片URL
	Summary string `json:"summary"` // 摘要
	URL     string `json:"url"`     // 内容URL
	Remark  string `json:"remark"`  // 备注
}

// ListFavoritesRequest 获取收藏列表请求DTO
type ListFavoritesRequest struct {
	Page      int64  `json:"page"`       // 页码
	PageSize  int64  `json:"page_size"`  // 每页大小
	UserID    int64  `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户名
	Title     string `json:"title"`      // 标题
	ItemType  string `json:"item_type"`  // 内容类型
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
	TenantID  int64  `json:"tenant_id"`  // 租户ID
}

// ListFavoritesResponse 获取收藏列表响应DTO
type ListFavoritesResponse struct {
	List  []*FavoriteDTO `json:"list"`  // 收藏列表
	Total int64          `json:"total"` // 总数
	Page  int64          `json:"page"`  // 当前页码
	Size  int64          `json:"size"`  // 每页大小
}

// BatchDeleteFavoritesRequest 批量删除收藏请求DTO
type BatchDeleteFavoritesRequest struct {
	IDs []int64 `json:"ids"` // 收藏ID列表
}

// CheckFavoriteRequest 检查是否已收藏请求DTO
type CheckFavoriteRequest struct {
	UserID   int64  `json:"user_id"`   // 用户ID
	ItemID   int64  `json:"item_id"`   // 内容项ID
	ItemType string `json:"item_type"` // 内容项类型
}

// CheckFavoriteResponse 检查是否已收藏响应DTO
type CheckFavoriteResponse struct {
	Exists bool `json:"exists"` // 是否已收藏
}

// TypeStatsDTO 类型统计DTO
type TypeStatsDTO struct {
	Type  string `json:"type"`  // 类型
	Count int64  `json:"count"` // 数量
}

// FavoritesStatisticsResponse 收藏统计响应DTO
type FavoritesStatisticsResponse struct {
	TotalUsers       int64           `json:"total_users"`       // 总收藏用户数
	TotalFavorites   int64           `json:"total_favorites"`   // 总收藏数
	TodayFavorites   int64           `json:"today_favorites"`   // 今日收藏
	MonthFavorites   int64           `json:"month_favorites"`   // 本月收藏
	TypeDistribution []*TypeStatsDTO `json:"type_distribution"` // 类型分布
}

// HotContentResponseDTO 热门内容响应DTO
type HotContentResponseDTO struct {
	ItemID     int64  `json:"item_id"`     // 内容ID
	ItemType   string `json:"item_type"`   // 内容类型
	Title      string `json:"title"`       // 标题
	Cover      string `json:"cover"`       // 封面图
	Count      int64  `json:"count"`       // 收藏次数
	CreateDate string `json:"create_date"` // 创建日期
}

// TrendDataResponseDTO 趋势数据响应DTO
type TrendDataResponseDTO struct {
	Date  string `json:"date"`  // 日期
	Count int64  `json:"count"` // 数量
}

// ExportFavoritesDataRequest 导出收藏数据请求DTO
type ExportFavoritesDataRequest struct {
	UserID    int64  `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户名
	Title     string `json:"title"`      // 标题
	ItemType  string `json:"item_type"`  // 内容类型
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
	TenantID  int64  `json:"tenant_id"`  // 租户ID
}

// GetFavoritesTrendRequest 获取收藏趋势请求DTO
type GetFavoritesTrendRequest struct {
	Period string `json:"period"` // 周期，如week/month/year
}
