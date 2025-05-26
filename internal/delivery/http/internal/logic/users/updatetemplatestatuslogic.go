package users

import (
	"context"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTemplateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTemplateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTemplateStatusLogic {
	return &UpdateTemplateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTemplateStatusLogic) UpdateTemplateStatus(req *types.UpdateTemplateStatusReq) (resp *types.UpdateTemplateStatusResp, err error) {
	// 先获取模板信息
	template, err := l.svcCtx.TemplateService.GetTemplate(l.ctx, req.TemplateID)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	// 检查权限
	if template.UserID != req.UserID {
		return nil, types.NewForbiddenError("无权操作此模板")
	}

	// 调用服务层更新模板状态
	err = l.svcCtx.TemplateService.UpdateTemplateStatus(l.ctx, req.TemplateID, req.Enabled)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	resp = &types.UpdateTemplateStatusResp{
		Success: true,
	}

	return resp, nil
}
