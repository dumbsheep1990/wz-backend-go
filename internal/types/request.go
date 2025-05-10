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
