package public

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"wz-backend-go/internal/domain/model"
	"wz-backend-go/internal/delivery/http/internal/svc"
)

type GetCategoryDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryDetailLogic {
	return &GetCategoryDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetCategoryDetail 获取分类详情
func (l *GetCategoryDetailLogic) GetCategoryDetail(tenantID int64, categoryID int) (*model.CategoryDetail, error) {
	// 获取租户信息
	tenant, err := l.svcCtx.TenantService.GetTenantByID(l.ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 根据租户类型和分类ID获取分类详情
	var categoryDetail model.CategoryDetail
	var found bool

	// 获取租户对应的导航分类
	navigation, err := NewGetTenantNavigationLogic(l.ctx, l.svcCtx).GetTenantNavigation(tenantID, "zh-CN")
	if err != nil {
		return nil, err
	}

	// 查找对应的分类
	for _, category := range navigation.Categories {
		if category.ID == categoryID {
			categoryDetail = model.CategoryDetail{
				ID:          category.ID,
				Name:        category.Name,
				Description: category.Description,
				Icon:        category.Icon,
				Items:       []model.SearchResult{},
			}
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("未找到分类ID: %d", categoryID)
	}

	// 根据租户类型和分类ID获取分类下的项目
	switch tenant.TenantType {
	case model.TenantTypeEnterprise:
		// 企业类型，根据分类ID获取商品列表
		switch categoryID {
		case 1: // 万知车船
			categoryDetail.Items = []model.SearchResult{
				{
					Type:  "vehicle",
					ID:    101,
					Title: "奔驰C级二手车",
					URL:   "/vehicleships/101",
				},
				{
					Type:  "vehicle",
					ID:    102,
					Title: "宝马3系二手车",
					URL:   "/vehicleships/102",
				},
				{
					Type:  "ship",
					ID:    103,
					Title: "私人游艇出售",
					URL:   "/vehicleships/103",
				},
			}
		case 2: // 万知生活
			categoryDetail.Items = []model.SearchResult{
				{
					Type:  "service",
					ID:    201,
					Title: "家政服务",
					URL:   "/life/201",
				},
				{
					Type:  "service",
					ID:    202,
					Title: "送餐服务",
					URL:   "/life/202",
				},
			}
		case 3: // 万知旅游
			categoryDetail.Items = []model.SearchResult{
				{
					Type:  "travel",
					ID:    301,
					Title: "三亚旅游套餐",
					URL:   "/travel/301",
				},
				{
					Type:  "travel",
					ID:    302,
					Title: "欧洲十日游",
					URL:   "/travel/302",
				},
			}
		case 4: // 交易中心
			categoryDetail.Items = []model.SearchResult{
				{
					Type:  "item",
					ID:    401,
					Title: "二手iPhone",
					URL:   "/trade/401",
				},
				{
					Type:  "item",
					ID:    402,
					Title: "家具出售",
					URL:   "/trade/402",
				},
			}
		}
	case model.TenantTypePersonal:
		// 个人类型，根据分类ID获取博客或作品集
		switch categoryID {
		case 1: // 博客文章
			categoryDetail.Items = []model.SearchResult{
				{
					Type:  "blog",
					ID:    501,
					Title: "我的旅行日记",
					URL:   "/blog/501",
				},
				{
					Type:  "blog",
					ID:    502,
					Title: "技术分享：Go语言实践",
					URL:   "/blog/502",
				},
			}
		case 2: // 作品展示
			categoryDetail.Items = []model.SearchResult{
				{
					Type:  "portfolio",
					ID:    601,
					Title: "摄影作品集",
					URL:   "/portfolio/601",
				},
				{
					Type:  "portfolio",
					ID:    602,
					Title: "设计作品展示",
					URL:   "/portfolio/602",
				},
			}
		}
	case model.TenantTypeEducational:
		// 教育机构类型，根据分类ID获取课程或学术资料
		switch categoryID {
		case 1: // 课程中心
			categoryDetail.Items = []model.SearchResult{
				{
					Type:  "course",
					ID:    701,
					Title: "Python编程入门",
					URL:   "/courses/701",
				},
				{
					Type:  "course",
					ID:    702,
					Title: "数据结构与算法",
					URL:   "/courses/702",
				},
			}
		case 2: // 学术资料
			categoryDetail.Items = []model.SearchResult{
				{
					Type:  "academic",
					ID:    801,
					Title: "机器学习研究报告",
					URL:   "/academic/801",
				},
				{
					Type:  "academic",
					ID:    802,
					Title: "人工智能论文集",
					URL:   "/academic/802",
				},
			}
		}
	}

	return &categoryDetail, nil
}
