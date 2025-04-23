package public

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"wz-backend-go/internal/domain/model"
	"wz-backend-go/internal/delivery/http/internal/svc"
)

type GetTenantNavigationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTenantNavigationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTenantNavigationLogic {
	return &GetTenantNavigationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetTenantNavigation 获取租户导航分类
func (l *GetTenantNavigationLogic) GetTenantNavigation(tenantID int64, lang string) (*model.NavigationResponse, error) {
	// 获取租户信息，用于确定租户类型
	tenant, err := l.svcCtx.TenantService.GetTenantByID(l.ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 根据租户类型返回不同的导航
	var categories []model.NavigationCategory

	switch tenant.TenantType {
	case model.TenantTypeEnterprise:
		// 企业类型导航，主要是商品交易相关
		categories = []model.NavigationCategory{
			{
				ID:          1,
				Name:        "万知车船",
				URL:         "/vehicleships",
				Icon:        "car",
				Description: "车辆和船舶交易平台",
			},
			{
				ID:          2,
				Name:        "万知生活",
				URL:         "/life",
				Icon:        "home",
				Description: "日常生活服务",
			},
			{
				ID:          3,
				Name:        "万知旅游",
				URL:         "/travel",
				Icon:        "plane",
				Description: "旅游相关服务",
			},
			{
				ID:          4,
				Name:        "交易中心",
				URL:         "/trade",
				Icon:        "shopping-cart",
				Description: "商品交易中心",
			},
		}
	case model.TenantTypePersonal:
		// 个人类型导航，主要是博客作品集相关
		categories = []model.NavigationCategory{
			{
				ID:          1,
				Name:        "博客文章",
				URL:         "/blog",
				Icon:        "book",
				Description: "博客文章列表",
			},
			{
				ID:          2,
				Name:        "作品展示",
				URL:         "/portfolio",
				Icon:        "image",
				Description: "个人作品展示",
			},
			{
				ID:          3,
				Name:        "关于我",
				URL:         "/about",
				Icon:        "user",
				Description: "个人介绍",
			},
			{
				ID:          4,
				Name:        "联系方式",
				URL:         "/contact",
				Icon:        "envelope",
				Description: "联系方式",
			},
		}
	case model.TenantTypeEducational:
		// 教育机构类型导航，主要是课程和学术资料相关
		categories = []model.NavigationCategory{
			{
				ID:          1,
				Name:        "课程中心",
				URL:         "/courses",
				Icon:        "graduation-cap",
				Description: "在线课程列表",
			},
			{
				ID:          2,
				Name:        "学术资料",
				URL:         "/academic",
				Icon:        "file-text",
				Description: "学术研究资料",
			},
			{
				ID:          3,
				Name:        "教师团队",
				URL:         "/teachers",
				Icon:        "users",
				Description: "教师团队介绍",
			},
			{
				ID:          4,
				Name:        "学习路径",
				URL:         "/learning-paths",
				Icon:        "road",
				Description: "推荐学习路径",
			},
		}
	default:
		// 默认导航
		categories = []model.NavigationCategory{
			{
				ID:          1,
				Name:        "首页",
				URL:         "/home",
				Icon:        "home",
				Description: "站点首页",
			},
			{
				ID:          2,
				Name:        "关于我们",
				URL:         "/about",
				Icon:        "info-circle",
				Description: "关于我们",
			},
			{
				ID:          3,
				Name:        "联系我们",
				URL:         "/contact",
				Icon:        "phone",
				Description: "联系我们",
			},
		}
	}

	// 如果有其他语言的支持，可以在这里添加语言判断
	if lang != "zh-CN" {
		// 处理其他语言...
	}

	return &model.NavigationResponse{
		Categories: categories,
	}, nil
}
