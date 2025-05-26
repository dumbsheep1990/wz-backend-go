package users

import (
	"context"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTemplateLogic {
	return &GetTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTemplateLogic) GetTemplate(userID, templateID int64) (resp *types.GetTemplateResp, err error) {
	// 调用服务层获取模板详情
	template, err := l.svcCtx.TemplateService.GetTemplate(l.ctx, templateID)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	// 检查权限
	if template.UserID != userID {
		return nil, types.NewForbiddenError("无权访问此模板")
	}

	// 转换为响应格式
	resp = &types.GetTemplateResp{
		ID:          template.ID,
		Name:        template.Name,
		Type:        types.TemplateType(template.Type),
		Preview:     template.Preview,
		Content:     template.Content,
		Enabled:     template.Enabled,
		IsNew:       template.IsNew,
		PublicShare: template.PublicShare,
		CreatedAt:   template.CreatedAt,
		UpdatedAt:   template.UpdatedAt,
	}

	return resp, nil
}
