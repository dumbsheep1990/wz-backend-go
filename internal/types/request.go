package types

// LinkListRequest 友情链接列表请求参数
type LinkListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Name     string `form:"name"`
	Status   *int   `form:"status"`
}

// LinkCreateRequest 创建友情链接请求参数
type LinkCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Logo        string `json:"logo"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status" binding:"required"`
	Description string `json:"description"`
}

// LinkUpdateRequest 更新友情链接请求参数
type LinkUpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Logo        string `json:"logo"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status" binding:"required"`
	Description string `json:"description"`
}

// SiteConfigUpdateRequest 更新站点配置请求参数
type SiteConfigUpdateRequest struct {
	SiteName       string `json:"site_name" binding:"required"`
	SiteLogo       string `json:"site_logo"`
	SeoTitle       string `json:"seo_title"`
	SeoKeywords    string `json:"seo_keywords"`
	SeoDescription string `json:"seo_description"`
	IcpNumber      string `json:"icp_number"`
	Copyright      string `json:"copyright"`
	ThemeID        int64  `json:"theme_id"`
	ContactEmail   string `json:"contact_email"`
	ContactPhone   string `json:"contact_phone"`
	Address        string `json:"address"`
}

// ThemeListRequest 主题列表请求参数
type ThemeListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Name     string `form:"name"`
	Status   *int   `form:"status"`
}

// ThemeCreateRequest 创建主题请求参数
type ThemeCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Preview     string `json:"preview"`
	Description string `json:"description"`
	Status      int    `json:"status" binding:"required"`
	IsDefault   int    `json:"is_default"`
	Config      string `json:"config"`
}

// ThemeUpdateRequest 更新主题请求参数
type ThemeUpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Preview     string `json:"preview"`
	Description string `json:"description"`
	Status      int    `json:"status" binding:"required"`
	IsDefault   int    `json:"is_default"`
	Config      string `json:"config"`
}

// UserMessageListRequest 用户消息列表请求参数
type UserMessageListRequest struct {
	Page     int  `form:"page"`
	PageSize int  `form:"pageSize"`
	Type     *int `form:"type"`
	Status   *int `form:"status"`
}

// UserMessageCreateRequest 创建用户消息请求参数
type UserMessageCreateRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	IsImportant int    `json:"is_important"`
	RelatedID   int64  `json:"related_id"`
	RelatedType string `json:"related_type"`
}

// UserPointsListRequest 用户积分列表请求参数
type UserPointsListRequest struct {
	Page     int   `form:"page"`
	PageSize int   `form:"pageSize"`
	UserID   int64 `form:"userId" binding:"required"`
}

// UserPointsCreateRequest 创建用户积分请求参数
type UserPointsCreateRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	Points      int    `json:"points" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	Source      string `json:"source" binding:"required"`
	Description string `json:"description"`
	RelatedID   int64  `json:"related_id"`
	RelatedType string `json:"related_type"`
}

// UserFavoriteListRequest 用户收藏列表请求参数
type UserFavoriteListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	UserID   int64  `form:"userId" binding:"required"`
	ItemType string `form:"itemType"`
}

// UserFavoriteCreateRequest 创建用户收藏请求参数
type UserFavoriteCreateRequest struct {
	UserID   int64  `json:"user_id" binding:"required"`
	ItemID   int64  `json:"item_id" binding:"required"`
	ItemType string `json:"item_type" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Cover    string `json:"cover"`
	Summary  string `json:"summary"`
	URL      string `json:"url"`
	Remark   string `json:"remark"`
}

// StatisticsRequest 统计数据请求参数
type StatisticsRequest struct {
	Type      string `form:"type" binding:"required"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

// 用户积分相关请求

// CreateUserPointsRequest 创建用户积分请求
type CreateUserPointsRequest struct {
	UserID      int64  `json:"user_id"`      // 用户ID
	Points      int    `json:"points"`       // 积分值
	Type        int    `json:"type"`         // 类型（1增加，2减少）
	Source      string `json:"source"`       // 来源
	Description string `json:"description"`  // 描述
	RelatedID   int64  `json:"related_id"`   // 关联ID
	RelatedType string `json:"related_type"` // 关联类型
}

// ListUserPointsRequest 获取用户积分列表请求
type ListUserPointsRequest struct {
	UserID   int64 `json:"user_id"`   // 用户ID
	Page     int   `json:"page"`      // 页码
	PageSize int   `json:"page_size"` // 每页数量
	Type     int   `json:"type"`      // 类型过滤（1增加，2减少）
	TenantID int64 `json:"tenant_id"` // 租户ID
}

// 用户收藏相关请求

// CreateUserFavoriteRequest 创建用户收藏请求
type CreateUserFavoriteRequest struct {
	UserID   int64  `json:"user_id"`   // 用户ID
	ItemID   int64  `json:"item_id"`   // 内容ID
	ItemType string `json:"item_type"` // 内容类型
	Title    string `json:"title"`     // 标题
	Cover    string `json:"cover"`     // 封面图
	Summary  string `json:"summary"`   // 内容摘要
	URL      string `json:"url"`       // 链接地址
	Remark   string `json:"remark"`    // 备注
}

// ListUserFavoritesRequest 获取用户收藏列表请求
type ListUserFavoritesRequest struct {
	UserID   int64  `json:"user_id"`   // 用户ID
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
	ItemType string `json:"item_type"` // 项目类型过滤
	TenantID int64  `json:"tenant_id"` // 租户ID
}

// ListPointsRequest 获取积分列表请求
type ListPointsRequest struct {
	Page      int64  `json:"page"`       // 当前页码
	PageSize  int64  `json:"page_size"`  // 每页大小
	UserID    int64  `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户名
	Type      int    `json:"type"`       // 积分类型
	Source    string `json:"source"`     // 积分来源
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
}

// ListFavoritesRequest 获取收藏列表请求
type ListFavoritesRequest struct {
	Page      int64  `json:"page"`       // 当前页码
	PageSize  int64  `json:"page_size"`  // 每页大小
	UserID    int64  `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户名
	Title     string `json:"title"`      // 标题
	ItemType  string `json:"item_type"`  // 内容类型
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
}

// AdminAddPointsRequest 管理员添加/调整积分请求
type AdminAddPointsRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`        // 用户ID
	Points      int    `json:"points" binding:"required,gt=0"`    // 积分值
	Type        int    `json:"type" binding:"required,oneof=1 2"` // 类型（1增加，2减少）
	Description string `json:"description" binding:"required"`    // 描述
}

// BatchDeleteFavoritesRequest 批量删除收藏请求
type BatchDeleteFavoritesRequest struct {
	IDs []int64 `json:"ids" binding:"required"` // 收藏ID列表
}

// PointsRulesRequest 积分规则设置请求
type PointsRulesRequest struct {
	SignInPoints      int  `json:"sign_in_points"`      // 签到积分
	CommentPoints     int  `json:"comment_points"`      // 评论积分
	SharePoints       int  `json:"share_points"`        // 分享积分
	ArticlePoints     int  `json:"article_points"`      // 发布文章积分
	InvitePoints      int  `json:"invite_points"`       // 邀请积分
	PurchaseRate      int  `json:"purchase_rate"`       // 购买积分比例
	MaxDailyPoints    int  `json:"max_daily_points"`    // 每日最大获取积分
	EnableExchange    bool `json:"enable_exchange"`     // 是否可兑换商品
	ExchangeRate      int  `json:"exchange_rate"`       // 兑换比例
	MinExchangePoints int  `json:"min_exchange_points"` // 最小兑换积分
}

// Page management types
type GetPageListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Type     string `form:"type"`
	Status   string `form:"status"`
	Title    string `form:"title"`
}

type CreatePageRequest struct {
	Title       string `json:"title" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Content     string `json:"content"`
	SeoTitle    string `json:"seo_title"`
	SeoDesc     string `json:"seo_desc"`
	SeoKeywords string `json:"seo_keywords"`
	Status      string `json:"status"`
}

type UpdatePageRequest struct {
	ID          int64  `json:"id"`
	Title       string `json:"title" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Content     string `json:"content"`
	SeoTitle    string `json:"seo_title"`
	SeoDesc     string `json:"seo_desc"`
	SeoKeywords string `json:"seo_keywords"`
	Status      string `json:"status"`
}

type DeletePageRequest struct {
	ID int64 `json:"id"`
}

type GetPageDetailRequest struct {
	ID int64 `json:"id"`
}

type TogglePageStatusRequest struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type PreviewPageRequest struct {
	ID int64 `json:"id"`
}

type BatchUpdatePageRequest struct {
	IDs    []int64 `json:"ids"`
	Status string  `json:"status"`
}

// Link management types
type GetLinkListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Category string `form:"category"`
	Status   string `form:"status"`
	Name     string `form:"name"`
}

type CreateLinkRequest struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	NewWindow   bool   `json:"new_window"`
	IsActive    bool   `json:"is_active"`
}

type UpdateLinkRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	NewWindow   bool   `json:"new_window"`
	IsActive    bool   `json:"is_active"`
}

type DeleteLinkRequest struct {
	ID int64 `json:"id"`
}

type GetLinkDetailRequest struct {
	ID int64 `json:"id"`
}

type VerifyLinkRequest struct {
	ID int64 `json:"id"`
}

type BatchVerifyLinksRequest struct {
	IDs []int64 `json:"ids"`
}

type GetLinkCategoriesRequest struct {
}

type UpdateLinkSortRequest struct {
	SortData []struct {
		ID   int64 `json:"id"`
		Sort int   `json:"sort"`
	} `json:"sort_data"`
}

// Component management types
type GetComponentListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Type     string `form:"type"`
	Name     string `form:"name"`
}

type CreateComponentRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	Code        string `json:"code" binding:"required"`
}

type UpdateComponentRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	Code        string `json:"code" binding:"required"`
}

type DeleteComponentRequest struct {
	ID int64 `json:"id"`
}

type GetComponentDetailRequest struct {
	ID int64 `json:"id"`
}

type PreviewComponentRequest struct {
	ID int64 `json:"id"`
}

type ImportComponentRequest struct {
	Source string `json:"source" binding:"required"`
}

type GetComponentTypesRequest struct {
}

// Theme management types  
type GetThemeListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Name     string `form:"name"`
}

type CreateThemeRequest struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	PrimaryColor      string `json:"primary_color" binding:"required"`
	TextColor         string `json:"text_color"`
	HeaderColor       string `json:"header_color"`
	LogoColor         string `json:"logo_color"`
	MenuTextColor     string `json:"menu_text_color"`
	ContentBgColor    string `json:"content_bg_color"`
	SidebarColor      string `json:"sidebar_color"`
	SidebarTextColor  string `json:"sidebar_text_color"`
	CardColor         string `json:"card_color"`
	LinkColor         string `json:"link_color"`
}

type UpdateThemeRequest struct {
	ID                int64  `json:"id"`
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	PrimaryColor      string `json:"primary_color" binding:"required"`
	TextColor         string `json:"text_color"`
	HeaderColor       string `json:"header_color"`
	LogoColor         string `json:"logo_color"`
	MenuTextColor     string `json:"menu_text_color"`
	ContentBgColor    string `json:"content_bg_color"`
	SidebarColor      string `json:"sidebar_color"`
	SidebarTextColor  string `json:"sidebar_text_color"`
	CardColor         string `json:"card_color"`
	LinkColor         string `json:"link_color"`
}

type DeleteThemeRequest struct {
	ID int64 `json:"id"`
}

type GetThemeDetailRequest struct {
	ID int64 `json:"id"`
}

type ApplyThemeRequest struct {
	ID int64 `json:"id"`
}

type GetCurrentThemeRequest struct {
}

type PreviewThemeRequest struct {
	ID int64 `json:"id"`
}

type ExportThemeRequest struct {
	ID int64 `json:"id"`
}

type ImportThemeRequest struct {
	Data string `json:"data" binding:"required"`
}

// 树形导航管理请求类型

// GetNavigationTreeRequest 获取导航树请求
type GetNavigationTreeRequest struct {
	Type string `json:"type" form:"type"` // main, footer, sidebar
}

// CreateNavigationItemRequest 创建导航项请求
type CreateNavigationItemRequest struct {
	Name      string `json:"name" binding:"required"`
	URL       string `json:"url" binding:"required"`
	Icon      string `json:"icon"`
	Visible   bool   `json:"visible"`
	NewWindow bool   `json:"newWindow"`
	ParentID  *int64 `json:"parentId"`
	SortOrder int32  `json:"sortOrder"`
	Type      string `json:"type"` // main, footer, sidebar
}

// UpdateNavigationItemRequest 更新导航项请求
type UpdateNavigationItemRequest struct {
	ID        int64  `json:"id"`
	Name      string `json:"name" binding:"required"`
	URL       string `json:"url" binding:"required"`
	Icon      string `json:"icon"`
	Visible   bool   `json:"visible"`
	NewWindow bool   `json:"newWindow"`
	ParentID  *int64 `json:"parentId"`
	SortOrder int32  `json:"sortOrder"`
	Type      string `json:"type"`
}

// DeleteNavigationItemRequest 删除导航项请求
type DeleteNavigationItemRequest struct {
	ID int64 `json:"id"`
}

// GetNavigationItemRequest 获取导航项详情请求
type GetNavigationItemRequest struct {
	ID int64 `json:"id"`
}

// UpdateNavigationOrderRequest 更新导航排序请求
type UpdateNavigationOrderRequest struct {
	Type      string                        `json:"type"`
	OrderData []NavigationOrderItemRequest `json:"orderData"`
}

// NavigationOrderItemRequest 导航排序项目
type NavigationOrderItemRequest struct {
	ID       int64 `json:"id"`
	ParentID *int64 `json:"parentId"`
	Order    int32 `json:"order"`
}

// ToggleNavigationVisibilityRequest 切换导航可见性请求
type ToggleNavigationVisibilityRequest struct {
	ID      int64 `json:"id"`
	Visible bool  `json:"visible"`
}

// BatchDeleteNavigationItemsRequest 批量删除导航项请求
type BatchDeleteNavigationItemsRequest struct {
	IDs []int64 `json:"ids"`
}

// ExportNavigationTreeRequest 导出导航树请求
type ExportNavigationTreeRequest struct {
	Type string `json:"type"`
}

// ImportNavigationTreeRequest 导入导航树请求
type ImportNavigationTreeRequest struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
