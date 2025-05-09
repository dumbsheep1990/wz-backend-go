package logic

import (
	"context"
	"wz-backend-go/api/rpc/recommend"
	"wz-backend-go/internal/admin/svc"
	"wz-backend-go/internal/admin/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetRecommendRulesLogic 获取推荐规则列表的逻辑处理器
type GetRecommendRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetRecommendRulesLogic 创建获取推荐规则列表的逻辑处理器
func NewGetRecommendRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendRulesLogic {
	return &GetRecommendRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetRecommendRules 获取推荐规则列表
func (l *GetRecommendRulesLogic) GetRecommendRules(req *types.RecommendRuleListReq) (*types.RecommendRuleListResp, error) {
	// 调用推荐服务获取推荐规则列表
	resp, err := l.svcCtx.RecommendClient.GetRecommendRules(l.ctx, &recommend.GetRecommendRulesReq{
		Type:   req.Type,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	list := make([]types.RecommendRuleDetail, 0, len(resp.Rules))
	for _, item := range resp.Rules {
		list = append(list, types.RecommendRuleDetail{
			ID:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Type:        item.Type,
			Params:      item.Params,
			Status:      item.Status,
			SortOrder:   item.SortOrder,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	return &types.RecommendRuleListResp{
		Total: resp.Total,
		List:  list,
	}, nil
}

// SaveRecommendRuleLogic 保存推荐规则的逻辑处理器
type SaveRecommendRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewSaveRecommendRuleLogic 创建保存推荐规则的逻辑处理器
func NewSaveRecommendRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRecommendRuleLogic {
	return &SaveRecommendRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SaveRecommendRule 保存推荐规则
func (l *SaveRecommendRuleLogic) SaveRecommendRule(req *types.SaveRecommendRuleReq) (*types.SuccessResp, error) {
	// 调用推荐服务保存推荐规则
	resp, err := l.svcCtx.RecommendClient.SetRecommendRule(l.ctx, &recommend.SetRecommendRuleReq{
		Id:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Params:      req.Params,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.SuccessResp{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// SetContentWeightLogic 设置内容权重的逻辑处理器
type SetContentWeightLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewSetContentWeightLogic 创建设置内容权重的逻辑处理器
func NewSetContentWeightLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetContentWeightLogic {
	return &SetContentWeightLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SetContentWeight 设置内容权重
func (l *SetContentWeightLogic) SetContentWeight(req *types.SetContentWeightReq) (*types.SuccessResp, error) {
	// 调用推荐服务设置内容权重
	resp, err := l.svcCtx.RecommendClient.SetContentWeight(l.ctx, &recommend.SetContentWeightReq{
		ContentId: req.ContentID,
		Weight:    req.Weight,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.SuccessResp{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// GetHotContentLogic 获取热门内容的逻辑处理器
type GetHotContentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetHotContentLogic 创建获取热门内容的逻辑处理器
func NewGetHotContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotContentLogic {
	return &GetHotContentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetHotContent 获取热门内容
func (l *GetHotContentLogic) GetHotContent(req *types.HotContentReq) (*types.HotContentResp, error) {
	// 调用推荐服务获取热门内容
	resp, err := l.svcCtx.RecommendClient.GetHotContent(l.ctx, &recommend.GetHotContentReq{
		SectionId:  req.SectionID,
		CategoryId: req.CategoryID,
		Type:       req.Type,
		Limit:      req.Limit,
		Period:     req.Period,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	list := make([]types.ContentBrief, 0, len(resp.List))
	for _, item := range resp.List {
		list = append(list, types.ContentBrief{
			ID:           item.Id,
			Type:         item.Type,
			Title:        item.Title,
			Summary:      item.Summary,
			CoverUrl:     item.CoverUrl,
			UserID:       item.UserId,
			CategoryID:   item.CategoryId,
			ViewCount:    item.ViewCount,
			LikeCount:    item.LikeCount,
			CommentCount: item.CommentCount,
			Tags:         item.Tags,
			CreatedAt:    item.CreatedAt,
			Weight:       item.Weight,
			HotScore:     item.HotScore,
		})
	}

	return &types.HotContentResp{
		Total: resp.Total,
		List:  list,
	}, nil
}

// GetHotKeywordsLogic 获取热门关键词的逻辑处理器
type GetHotKeywordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetHotKeywordsLogic 创建获取热门关键词的逻辑处理器
func NewGetHotKeywordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotKeywordsLogic {
	return &GetHotKeywordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetHotKeywords 获取热门关键词
func (l *GetHotKeywordsLogic) GetHotKeywords(req *types.HotKeywordsReq) (*types.HotKeywordsResp, error) {
	// 调用推荐服务获取热门关键词
	resp, err := l.svcCtx.RecommendClient.GetHotKeywords(l.ctx, &recommend.GetHotKeywordsReq{
		Limit:  req.Limit,
		Period: req.Period,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	keywords := make([]types.KeywordItem, 0, len(resp.Keywords))
	for _, item := range resp.Keywords {
		keywords = append(keywords, types.KeywordItem{
			Keyword:     item.Keyword,
			SearchCount: item.SearchCount,
			Trend:       item.Trend,
		})
	}

	return &types.HotKeywordsResp{
		Keywords: keywords,
	}, nil
}

// GetHotCategoriesLogic 获取热门分类的逻辑处理器
type GetHotCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetHotCategoriesLogic 创建获取热门分类的逻辑处理器
func NewGetHotCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotCategoriesLogic {
	return &GetHotCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetHotCategories 获取热门分类
func (l *GetHotCategoriesLogic) GetHotCategories(req *types.HotCategoriesReq) (*types.HotCategoriesResp, error) {
	// 调用推荐服务获取热门分类
	resp, err := l.svcCtx.RecommendClient.GetHotCategories(l.ctx, &recommend.GetHotCategoriesReq{
		SectionId: req.SectionID,
		Limit:     req.Limit,
		Period:    req.Period,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	categories := make([]types.CategoryItem, 0, len(resp.Categories))
	for _, item := range resp.Categories {
		categories = append(categories, types.CategoryItem{
			ID:           item.Id,
			Name:         item.Name,
			ViewCount:    item.ViewCount,
			ContentCount: item.ContentCount,
			Trend:        item.Trend,
		})
	}

	return &types.HotCategoriesResp{
		Categories: categories,
	}, nil
}
