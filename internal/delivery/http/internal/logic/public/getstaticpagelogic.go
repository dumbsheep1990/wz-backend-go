package public

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"wz-backend-go/internal/domain/model"
	"wz-backend-go/internal/delivery/http/internal/svc"
)

type GetStaticPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStaticPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStaticPageLogic {
	return &GetStaticPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetStaticPage 获取静态页面内容
func (l *GetStaticPageLogic) GetStaticPage(tenantID int64, pageName string) (*model.StaticPageResponse, error) {
	// 获取租户信息
	tenant, err := l.svcCtx.TenantService.GetTenantByID(l.ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 根据租户ID和页面名称获取静态页面内容
	// 在实际实现中，应该从数据库或其他存储中获取
	var title, content string
	
	// 常见静态页面类型处理
	switch pageName {
	case "privacy":
		title = "隐私政策"
		content = fmt.Sprintf("<h1>%s的隐私政策</h1><p>本隐私政策描述了我们如何收集、使用和保护您的个人信息。</p>", tenant.Name)
	case "terms":
		title = "服务条款"
		content = fmt.Sprintf("<h1>%s的服务条款</h1><p>使用我们的服务前，请仔细阅读以下条款和条件。</p>", tenant.Name)
	case "about":
		title = "关于我们"
		content = fmt.Sprintf("<h1>关于%s</h1><p>%s</p>", tenant.Name, tenant.Description)
	case "contact":
		title = "联系我们"
		content = fmt.Sprintf("<h1>联系%s</h1><p>请通过以下方式联系我们。</p>", tenant.Name)
	default:
		return nil, fmt.Errorf("页面不存在: %s", pageName)
	}

	// 根据租户类型可能会有不同的内容
	switch tenant.TenantType {
	case model.TenantTypeEnterprise:
		// 企业类型可能有一些额外的页面内容
		if pageName == "about" {
			content += "<p>作为一家专业的企业，我们提供各种商品和服务。</p>"
		}
	case model.TenantTypePersonal:
		// 个人类型可能有一些个性化的内容
		if pageName == "about" {
			content += "<p>这是我的个人博客和作品集展示平台。</p>"
		}
	case model.TenantTypeEducational:
		// 教育机构类型可能有一些教育相关的内容
		if pageName == "about" {
			content += "<p>作为一家教育机构，我们致力于提供高质量的教育课程和学术资源。</p>"
		}
	}

	return &model.StaticPageResponse{
		Title:   title,
		Content: content,
	}, nil
}
