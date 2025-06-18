package types

import "time"

// Response 通用API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// UserPointsResponse 用户积分响应
type UserPointsResponse struct {
	ID          int64     `json:"id"`           // ID
	UserID      int64     `json:"user_id"`      // 用户ID
	Username    string    `json:"username"`     // 用户名
	Points      int       `json:"points"`       // 积分变动值
	TotalPoints int       `json:"total_points"` // 总积分
	Type        int       `json:"type"`         // 类型（1增加，2减少）
	Source      string    `json:"source"`       // 来源
	Description string    `json:"description"`  // 描述
	RelatedID   int64     `json:"related_id"`   // 关联ID
	RelatedType string    `json:"related_type"` // 关联类型
	Operator    string    `json:"operator"`     // 操作员（管理员调整时）
	OperatorID  int64     `json:"operator_id"`  // 操作员ID
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
}

// PagedUserPointsResponse 分页用户积分响应
type PagedUserPointsResponse struct {
	List  []*UserPointsResponse `json:"list"`  // 列表数据
	Total int64                 `json:"total"` // 总数
	Page  int                   `json:"page"`  // 当前页
	Size  int                   `json:"size"`  // 每页大小
}

// UserFavoriteResponse 用户收藏响应
type UserFavoriteResponse struct {
	ID        int64     `json:"id"`         // ID
	UserID    int64     `json:"user_id"`    // 用户ID
	Username  string    `json:"username"`   // 用户名
	ItemID    int64     `json:"item_id"`    // 内容ID
	ItemType  string    `json:"item_type"`  // 内容类型
	Title     string    `json:"title"`      // 标题
	Cover     string    `json:"cover"`      // 封面图
	Summary   string    `json:"summary"`    // 内容摘要
	URL       string    `json:"url"`        // 链接地址
	Remark    string    `json:"remark"`     // 备注
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// PagedUserFavoriteResponse 分页用户收藏响应
type PagedUserFavoriteResponse struct {
	List  []*UserFavoriteResponse `json:"list"`  // 列表数据
	Total int64                   `json:"total"` // 总数
	Page  int                     `json:"page"`  // 当前页
	Size  int                     `json:"size"`  // 每页大小
}

// PointsStatisticsResponse 积分统计响应
type PointsStatisticsResponse struct {
	TotalUsers         int64          `json:"total_users"`         // 总用户数
	TotalPoints        int64          `json:"total_points"`        // 总积分数
	AvgPoints          int64          `json:"avg_points"`          // 平均积分
	MaxPoints          int64          `json:"max_points"`          // 最高积分
	TodayIncrease      int64          `json:"today_increase"`      // 今日增加
	TodayDecrease      int64          `json:"today_decrease"`      // 今日减少
	MonthIncrease      int64          `json:"month_increase"`      // 本月增加
	MonthDecrease      int64          `json:"month_decrease"`      // 本月减少
	SourceDistribution []*SourceStats `json:"source_distribution"` // 来源分布
}

// SourceStats 来源统计
type SourceStats struct {
	Source string `json:"source"` // 来源
	Count  int64  `json:"count"`  // 数量
}

// FavoritesStatisticsResponse 收藏统计响应
type FavoritesStatisticsResponse struct {
	TotalUsers       int64        `json:"total_users"`       // 总收藏用户数
	TotalFavorites   int64        `json:"total_favorites"`   // 总收藏数
	TodayFavorites   int64        `json:"today_favorites"`   // 今日收藏
	MonthFavorites   int64        `json:"month_favorites"`   // 本月收藏
	TypeDistribution []*TypeStats `json:"type_distribution"` // 类型分布
}

// TypeStats 类型统计
type TypeStats struct {
	Type  string `json:"type"`  // 类型
	Count int64  `json:"count"` // 数量
}

// HotContentResponse 热门内容响应
type HotContentResponse struct {
	ItemID     int64  `json:"item_id"`     // 内容ID
	ItemType   string `json:"item_type"`   // 内容类型
	Title      string `json:"title"`       // 标题
	Cover      string `json:"cover"`       // 封面图
	Count      int64  `json:"count"`       // 收藏次数
	CreateDate string `json:"create_date"` // 创建日期
}

// TrendDataResponse 趋势数据响应
type TrendDataResponse struct {
	Date  string `json:"date"`  // 日期
	Count int64  `json:"count"` // 数量
}

// PointsRulesResponse 积分规则响应
type PointsRulesResponse struct {
	SignInPoints      int       `json:"sign_in_points"`      // 签到积分
	CommentPoints     int       `json:"comment_points"`      // 评论积分
	SharePoints       int       `json:"share_points"`        // 分享积分
	ArticlePoints     int       `json:"article_points"`      // 发布文章积分
	InvitePoints      int       `json:"invite_points"`       // 邀请积分
	PurchaseRate      int       `json:"purchase_rate"`       // 购买积分比例
	MaxDailyPoints    int       `json:"max_daily_points"`    // 每日最大获取积分
	EnableExchange    bool      `json:"enable_exchange"`     // 是否可兑换商品
	ExchangeRate      int       `json:"exchange_rate"`       // 兑换比例
	MinExchangePoints int       `json:"min_exchange_points"` // 最小兑换积分
	UpdatedAt         time.Time `json:"updated_at"`          // 更新时间
}

// Page management responses
type Page struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Path        string    `json:"path"`
	Type        string    `json:"type"`
	Content     string    `json:"content"`
	SeoTitle    string    `json:"seo_title"`
	SeoDesc     string    `json:"seo_desc"`
	SeoKeywords string    `json:"seo_keywords"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetPageListResponse struct {
	Pages []*Page `json:"pages"`
	Total int     `json:"total"`
}

type PagePreview struct {
	HTML string `json:"html"`
	CSS  string `json:"css"`
	JS   string `json:"js"`
}

// Link management responses
type Link struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Category    string    `json:"category"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	Sort        int       `json:"sort"`
	NewWindow   bool      `json:"new_window"`
	IsActive    bool      `json:"is_active"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetLinkListResponse struct {
	Links []*Link `json:"links"`
	Total int     `json:"total"`
}

type LinkVerifyResult struct {
	IsValid      bool   `json:"is_valid"`
	ResponseCode int    `json:"response_code"`
	ResponseTime int64  `json:"response_time"`
	ErrorMessage string `json:"error_message"`
}

type BatchVerifyResult struct {
	SuccessCount int `json:"success_count"`
	FailCount    int `json:"fail_count"`
	TotalCount   int `json:"total_count"`
}

type LinkCategory struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Component management responses
type Component struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Code        string    `json:"code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetComponentListResponse struct {
	Components []*Component `json:"components"`
	Total      int          `json:"total"`
}

type ComponentPreview struct {
	HTML string `json:"html"`
	CSS  string `json:"css"`
	JS   string `json:"js"`
}

type ComponentType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Theme management responses
type Theme struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	PrimaryColor      string    `json:"primary_color"`
	TextColor         string    `json:"text_color"`
	HeaderColor       string    `json:"header_color"`
	LogoColor         string    `json:"logo_color"`
	MenuTextColor     string    `json:"menu_text_color"`
	ContentBgColor    string    `json:"content_bg_color"`
	SidebarColor      string    `json:"sidebar_color"`
	SidebarTextColor  string    `json:"sidebar_text_color"`
	CardColor         string    `json:"card_color"`
	LinkColor         string    `json:"link_color"`
	IsDefault         bool      `json:"is_default"`
	IsCurrent         bool      `json:"is_current"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type GetThemeListResponse struct {
	Themes []*Theme `json:"themes"`
	Total  int      `json:"total"`
}

type ApplyThemeResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ThemePreview struct {
	CSS string `json:"css"`
}

type ThemeExport struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

// 树形导航管理响应类型

// NavigationItem 导航项
type NavigationItem struct {
	ID        int64             `json:"id"`
	Name      string            `json:"name"`
	URL       string            `json:"url"`
	Icon      string            `json:"icon"`
	Visible   bool              `json:"visible"`
	NewWindow bool              `json:"newWindow"`
	ParentID  *int64            `json:"parentId"`
	SortOrder int32             `json:"sortOrder"`
	Type      string            `json:"type"`
	Children  []*NavigationItem `json:"children"`
	CreatedAt string            `json:"createdAt,omitempty"`
	UpdatedAt string            `json:"updatedAt,omitempty"`
}

// GetNavigationTreeResponse 获取导航树响应
type GetNavigationTreeResponse struct {
	NavigationTree []*NavigationItem `json:"navigationTree"`
}

// NavigationTreeExport 导航树导出数据
type NavigationTreeExport struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

// NavigationTreeImportResult 导航树导入结果
type NavigationTreeImportResult struct {
	Success       bool   `json:"success"`
	ImportedCount int    `json:"importedCount"`
	Message       string `json:"message"`
}
