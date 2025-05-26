package users

import (
	"context"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTemplatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTemplatesLogic {
	return &GetTemplatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTemplatesLogic) GetTemplates(userID int64, page, pageSize int) (resp *types.GetTemplatesResp, err error) {
	// 调用服务层获取模板列表
	templates, total, err := l.svcCtx.TemplateService.GetTemplates(l.ctx, userID, page, pageSize)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	// 转换为响应格式
	items := make([]types.TemplateItem, 0, len(templates))
	for _, template := range templates {
		items = append(items, types.TemplateItem{
			ID:          template.ID,
			Name:        template.Name,
			Type:        types.TemplateType(template.Type),
			Preview:     template.Preview,
			Enabled:     template.Enabled,
			IsNew:       template.IsNew,
			PublicShare: template.PublicShare,
			CreatedAt:   template.CreatedAt,
		})
	}

	resp = &types.GetTemplatesResp{
		Total: total,
		Page:  page,
		Size:  pageSize,
		Items: items,
	}

	return resp, nil
}
