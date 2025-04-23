package public

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"wz-backend-go/internal/domain/model"
	"wz-backend-go/internal/delivery/http/internal/svc"
)

type GetRecommendationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRecommendationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendationsLogic {
	return &GetRecommendationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetRecommendations 获取租户推荐内容
func (l *GetRecommendationsLogic) GetRecommendations(tenantID int64, limit int, contentType string) (*model.RecommendationResponse, error) {
	// 获取租户信息，用于确定推荐内容类型
	tenant, err := l.svcCtx.TenantService.GetTenantByID(l.ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 根据租户类型和可选的内容类型参数获取推荐内容
	var recommendations []model.SearchResult

	switch tenant.TenantType {
	case model.TenantTypeEnterprise:
		// 企业类型推荐，主要是商品交易相关
		if contentType == "" || contentType == "item" {
			// 商品推荐
			recommendations = append(recommendations, model.SearchResult{
				Type:  "item",
				ID:    101,
				Title: "热门二手iPhone",
				URL:   "/idle/101",
			})
			recommendations = append(recommendations, model.SearchResult{
				Type:  "item",
				ID:    102,
				Title: "精选二手车",
				URL:   "/vehicleships/102",
			})
		}
	case model.TenantTypePersonal:
		// 个人类型推荐，主要是博客作品集相关
		if contentType == "" || contentType == "blog" {
			// 博客推荐
			recommendations = append(recommendations, model.SearchResult{
				Type:  "blog",
				ID:    201,
				Title: "最新博客文章",
				URL:   "/blog/201",
			})
		}
		if contentType == "" || contentType == "portfolio" {
			// 作品集推荐
			recommendations = append(recommendations, model.SearchResult{
				Type:  "portfolio",
				ID:    202,
				Title: "精选作品展示",
				URL:   "/portfolio/202",
			})
		}
	case model.TenantTypeEducational:
		// 教育机构类型推荐，主要是课程和学术资料相关
		if contentType == "" || contentType == "course" {
			// 课程推荐
			recommendations = append(recommendations, model.SearchResult{
				Type:  "course",
				ID:    301,
				Title: "热门课程：Python编程入门",
				URL:   "/courses/301",
			})
		}
		if contentType == "" || contentType == "academic" {
			// 学术资料推荐
			recommendations = append(recommendations, model.SearchResult{
				Type:  "academic",
				ID:    302,
				Title: "最新研究报告：机器学习应用",
				URL:   "/academic/302",
			})
		}
	}

	// 限制返回数量
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return &model.RecommendationResponse{
		Recommendations: recommendations,
	}, nil
}
