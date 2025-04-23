package public

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"wz-backend-go/internal/domain/model"
	"wz-backend-go/internal/delivery/http/internal/svc"
)

type SearchTenantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchTenantLogic {
	return &SearchTenantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SearchTenant 搜索租户内容
func (l *SearchTenantLogic) SearchTenant(tenantID int64, keyword string, page, limit int) (*model.SearchResponse, error) {
	// 获取租户信息，用于确定搜索的内容类型
	tenant, err := l.svcCtx.TenantService.GetTenantByID(l.ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 根据租户类型进行不同类型的搜索
	var results []model.SearchResult
	var total int

	switch tenant.TenantType {
	case model.TenantTypeEnterprise:
		// 企业类型，主要搜索商品交易相关内容
		// 这里应该调用具体的商品搜索服务
		// 示例数据
		results = []model.SearchResult{
			{
				Type:  "item",
				ID:    101,
				Title: "二手iPhone 13",
				URL:   "/idle/101",
			},
			{
				Type:  "item",
				ID:    102,
				Title: "奔驰C级二手车",
				URL:   "/vehicleships/102",
			},
		}
		total = 2
	case model.TenantTypePersonal:
		// 个人类型，主要搜索博客和作品集
		// 这里应该调用具体的博客、作品搜索服务
		// 示例数据
		results = []model.SearchResult{
			{
				Type:  "blog",
				ID:    201,
				Title: "我的旅行日记",
				URL:   "/blog/201",
			},
			{
				Type:  "portfolio",
				ID:    202,
				Title: "摄影作品集",
				URL:   "/portfolio/202",
			},
		}
		total = 2
	case model.TenantTypeEducational:
		// 教育机构类型，主要搜索课程和学术资料
		// 这里应该调用具体的课程、学术资料搜索服务
		// 示例数据
		results = []model.SearchResult{
			{
				Type:  "course",
				ID:    301,
				Title: "Python编程入门",
				URL:   "/courses/301",
			},
			{
				Type:  "academic",
				ID:    302,
				Title: "机器学习研究报告",
				URL:   "/academic/302",
			},
		}
		total = 2
	default:
		// 默认搜索
		results = []model.SearchResult{}
		total = 0
	}

	return &model.SearchResponse{
		Results: results,
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}
