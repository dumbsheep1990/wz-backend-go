package users

import (
	"context"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTemplateLogic {
	return &DeleteTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTemplateLogic) DeleteTemplate(userID, templateID int64) (resp *types.DeleteTemplateResp, err error) {
	// 先获取模板信息
	template, err := l.svcCtx.TemplateService.GetTemplate(l.ctx, templateID)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	// 检查权限
	if template.UserID != userID {
		return nil, types.NewForbiddenError("无权删除此模板")
	}

	// 调用服务层删除模板
	err = l.svcCtx.TemplateService.DeleteTemplate(l.ctx, templateID)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	resp = &types.DeleteTemplateResp{
		Success: true,
	}

	return resp, nil
}
