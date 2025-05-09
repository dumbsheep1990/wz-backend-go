package recommend

import (
	"context"

	"google.golang.org/grpc"
)

// 推荐服务接口定义
type RecommendService interface {
	// 获取热门内容
	GetHotContent(ctx context.Context, in *GetHotContentReq) (*GetHotContentResp, error)
	// 获取推荐内容
	GetRecommendContent(ctx context.Context, in *GetRecommendContentReq) (*GetRecommendContentResp, error)
	// 获取相关内容
	GetRelatedContent(ctx context.Context, in *GetRelatedContentReq) (*GetRelatedContentResp, error)
	// 获取热门搜索关键词
	GetHotKeywords(ctx context.Context, in *GetHotKeywordsReq) (*GetHotKeywordsResp, error)
	// 获取热门分类
	GetHotCategories(ctx context.Context, in *GetHotCategoriesReq) (*GetHotCategoriesResp, error)
	// 设置热门内容权重
	SetContentWeight(ctx context.Context, in *SetContentWeightReq) (*SetContentWeightResp, error)
	// 设置推荐规则
	SetRecommendRule(ctx context.Context, in *SetRecommendRuleReq) (*SetRecommendRuleResp, error)
	// 获取推荐规则
	GetRecommendRules(ctx context.Context, in *GetRecommendRulesReq) (*GetRecommendRulesResp, error)
}

// 推荐服务RPC客户端
type recommendServiceClient struct {
	conn *grpc.ClientConn
}

// 创建推荐服务客户端
func NewRecommendService(conn *grpc.ClientConn) RecommendService {
	return &recommendServiceClient{conn: conn}
}

// 以下是请求和响应结构体定义

// 热门内容请求
type GetHotContentReq struct {
	SectionId  int64  `json:"section_id,omitempty"`  // 板块ID
	CategoryId int64  `json:"category_id,omitempty"` // 分类ID
	Type       string `json:"type,omitempty"`        // 内容类型
	Limit      int32  `json:"limit,omitempty"`       // 返回数量
	Period     string `json:"period,omitempty"`      // 时间段：day, week, month, all
}

// 内容简要信息
type ContentBrief struct {
	Id           int64    `json:"id"`
	Type         string   `json:"type"`
	Title        string   `json:"title,omitempty"`
	Summary      string   `json:"summary,omitempty"`
	CoverUrl     string   `json:"cover_url,omitempty"`
	UserId       int64    `json:"user_id"`
	CategoryId   int64    `json:"category_id,omitempty"`
	ViewCount    int64    `json:"view_count,omitempty"`
	LikeCount    int64    `json:"like_count,omitempty"`
	CommentCount int64    `json:"comment_count,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	CreatedAt    string   `json:"created_at"`
	Weight       float64  `json:"weight,omitempty"`    // 权重
	HotScore     float64  `json:"hot_score,omitempty"` // 热度分数
}

// 热门内容响应
type GetHotContentResp struct {
	List  []*ContentBrief `json:"list"`
	Total int64           `json:"total"`
}

// 推荐内容请求
type GetRecommendContentReq struct {
	UserId     int64  `json:"user_id,omitempty"` // 用户ID，用于个性化推荐
	SectionId  int64  `json:"section_id,omitempty"`
	CategoryId int64  `json:"category_id,omitempty"`
	Type       string `json:"type,omitempty"`
	Limit      int32  `json:"limit,omitempty"`
	Offset     int32  `json:"offset,omitempty"`
}

// 推荐内容响应
type GetRecommendContentResp struct {
	List  []*ContentBrief `json:"list"`
	Total int64           `json:"total"`
}

// 相关内容请求
type GetRelatedContentReq struct {
	ContentId int64  `json:"content_id"` // 内容ID
	Type      string `json:"type,omitempty"`
	Limit     int32  `json:"limit,omitempty"`
}

// 相关内容响应
type GetRelatedContentResp struct {
	List []*ContentBrief `json:"list"`
}

// 热门搜索关键词请求
type GetHotKeywordsReq struct {
	Limit  int32  `json:"limit,omitempty"`
	Period string `json:"period,omitempty"` // 时间段：day, week, month
}

// 关键词项
type KeywordItem struct {
	Keyword     string `json:"keyword"`
	SearchCount int64  `json:"search_count"`
	Trend       int32  `json:"trend"` // 趋势: 1-上升, 0-不变, -1-下降
}

// 热门搜索关键词响应
type GetHotKeywordsResp struct {
	Keywords []*KeywordItem `json:"keywords"`
}

// 热门分类请求
type GetHotCategoriesReq struct {
	SectionId int64  `json:"section_id,omitempty"`
	Limit     int32  `json:"limit,omitempty"`
	Period    string `json:"period,omitempty"` // 时间段：day, week, month
}

// 分类项
type CategoryItem struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	ViewCount    int64  `json:"view_count"`
	ContentCount int64  `json:"content_count"`
	Trend        int32  `json:"trend"` // 趋势: 1-上升, 0-不变, -1-下降
}

// 热门分类响应
type GetHotCategoriesResp struct {
	Categories []*CategoryItem `json:"categories"`
}

// 设置内容权重请求
type SetContentWeightReq struct {
	ContentId int64   `json:"content_id"`
	Weight    float64 `json:"weight"`
}

// 设置内容权重响应
type SetContentWeightResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 推荐规则
type RecommendRule struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`             // hot-热门, latest-最新, personalized-个性化
	Params      string `json:"params,omitempty"` // JSON格式的参数
	Status      int32  `json:"status"`
	SortOrder   int32  `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 设置推荐规则请求
type SetRecommendRuleReq struct {
	Id          int64  `json:"id,omitempty"` // 为0时创建新规则
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Params      string `json:"params,omitempty"`
	Status      int32  `json:"status,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
}

// 设置推荐规则响应
type SetRecommendRuleResp struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Rule    *RecommendRule `json:"rule,omitempty"`
}

// 获取推荐规则请求
type GetRecommendRulesReq struct {
	Type   string `json:"type,omitempty"`
	Status int32  `json:"status,omitempty"`
}

// 获取推荐规则响应
type GetRecommendRulesResp struct {
	Rules []*RecommendRule `json:"rules"`
	Total int64            `json:"total"`
}

// 实现RecommendService接口的方法
func (c *recommendServiceClient) GetHotContent(ctx context.Context, in *GetHotContentReq) (*GetHotContentResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetHotContentResp{
		List:  []*ContentBrief{},
		Total: 0,
	}, nil
}

func (c *recommendServiceClient) GetRecommendContent(ctx context.Context, in *GetRecommendContentReq) (*GetRecommendContentResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetRecommendContentResp{
		List:  []*ContentBrief{},
		Total: 0,
	}, nil
}

func (c *recommendServiceClient) GetRelatedContent(ctx context.Context, in *GetRelatedContentReq) (*GetRelatedContentResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetRelatedContentResp{
		List: []*ContentBrief{},
	}, nil
}

func (c *recommendServiceClient) GetHotKeywords(ctx context.Context, in *GetHotKeywordsReq) (*GetHotKeywordsResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetHotKeywordsResp{
		Keywords: []*KeywordItem{},
	}, nil
}

func (c *recommendServiceClient) GetHotCategories(ctx context.Context, in *GetHotCategoriesReq) (*GetHotCategoriesResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetHotCategoriesResp{
		Categories: []*CategoryItem{},
	}, nil
}

func (c *recommendServiceClient) SetContentWeight(ctx context.Context, in *SetContentWeightReq) (*SetContentWeightResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &SetContentWeightResp{
		Success: true,
	}, nil
}

func (c *recommendServiceClient) SetRecommendRule(ctx context.Context, in *SetRecommendRuleReq) (*SetRecommendRuleResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &SetRecommendRuleResp{
		Success: true,
		Rule:    &RecommendRule{},
	}, nil
}

func (c *recommendServiceClient) GetRecommendRules(ctx context.Context, in *GetRecommendRulesReq) (*GetRecommendRulesResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetRecommendRulesResp{
		Rules: []*RecommendRule{},
		Total: 0,
	}, nil
}
