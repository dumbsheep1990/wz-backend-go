package users

import (
	"context"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"
	"github.com/wz-project/wz-backend-go/internal/domain/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTemplateLogic {
	return &UpdateTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTemplateLogic) UpdateTemplate(req *types.UpdateTemplateReq) (resp *types.UpdateTemplateResp, err error) {
	// 验证请求参数
	if req.Name == "" {
		return nil, types.NewBadRequestError("模板名称不能为空")
	}

	// 先获取原模板信息
	existingTemplate, err := l.svcCtx.TemplateService.GetTemplate(l.ctx, req.TemplateID)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	// 检查权限
	if existingTemplate.UserID != req.UserID {
		return nil, types.NewForbiddenError("无权操作此模板")
	}

	// 更新模板信息
	template := &model.Template{
		ID:          req.TemplateID,
		UserID:      req.UserID,
		Name:        req.Name,
		Type:        model.TemplateType(req.Type),
		Preview:     req.Preview,
		Content:     req.Content,
		PublicShare: req.PublicShare,
		Enabled:     existingTemplate.Enabled, // 保持原有启用状态
		IsNew:       existingTemplate.IsNew,   // 保持原有新模板标记
	}

	// 调用服务层更新模板
	err = l.svcCtx.TemplateService.UpdateTemplate(l.ctx, template)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	resp = &types.UpdateTemplateResp{
		Success: true,
	}

	return resp, nil
}
