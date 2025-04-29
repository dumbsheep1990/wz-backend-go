package model

// TenantListResponse 租户列表响应
type TenantListResponse struct {
	Tenants []TenantInfo `json:"tenants"`
}

// TenantInfo 租户信息
type TenantInfo struct {
	ID          string `json:"id"`          // 租户ID
	Name        string `json:"name"`        // 租户名称
	Description string `json:"description"` // 租户描述
	Subdomain   string `json:"subdomain"`   // 租户子域名
	TenantType  string `json:"tenant_type"` // 租户类型
	Logo        string `json:"logo"`        // 租户logo
}

// NavigationCategory 导航分类
type NavigationCategory struct {
	ID          int    `json:"id"`                    // 分类ID
	Name        string `json:"name"`                  // 分类名称
	URL         string `json:"url"`                   // 分类URL
	Icon        string `json:"icon,omitempty"`        // 分类图标
	Description string `json:"description,omitempty"` // 分类描述
}

// NavigationResponse 导航响应
type NavigationResponse struct {
	Categories []NavigationCategory `json:"categories"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Results []SearchResult `json:"results"` // 搜索结果列表
	Total   int            `json:"total"`   // 结果总数
	Page    int            `json:"page"`    // 当前页码
	Limit   int            `json:"limit"`   // 每页条数
}

// RecommendationResponse 推荐内容响应
type RecommendationResponse struct {
	Recommendations []SearchResult `json:"recommendations"` // 推荐内容列表
}

// CategoryDetail 分类详情
type CategoryDetail struct {
	ID          int            `json:"id"`          // 分类ID
	Name        string         `json:"name"`        // 分类名称
	Description string         `json:"description"` // 分类描述
	Icon        string         `json:"icon"`        // 分类图标
	Items       []SearchResult `json:"items"`       // 分类下的项目
}

// StaticPageResponse 静态页面响应
type StaticPageResponse struct {
	Title   string `json:"title"`   // 页面标题
	Content string `json:"content"` // 页面HTML内容
}
