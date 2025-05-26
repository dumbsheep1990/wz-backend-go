package users

import (
	"context"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"
	"github.com/wz-project/wz-backend-go/internal/domain/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTemplateLogic {
	return &CreateTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTemplateLogic) CreateTemplate(req *types.CreateTemplateReq) (resp *types.CreateTemplateResp, err error) {
	// 验证请求参数
	if req.Name == "" {
		return nil, types.NewBadRequestError("模板名称不能为空")
	}

	// 构建领域模型
	template := &model.Template{
		UserID:      req.UserID,
		Name:        req.Name,
		Type:        model.TemplateType(req.Type),
		Preview:     req.Preview,
		Content:     req.Content,
		PublicShare: req.PublicShare,
		Enabled:     true, // 默认启用
		IsNew:       true, // 新创建默认为新模板
	}

	// 调用服务层创建模板
	err = l.svcCtx.TemplateService.CreateTemplate(l.ctx, template)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	// 构建响应
	resp = &types.CreateTemplateResp{
		ID:        template.ID,
		Name:      template.Name,
		Type:      types.TemplateType(template.Type),
		Preview:   template.Preview,
		Enabled:   template.Enabled,
		IsNew:     template.IsNew,
		CreatedAt: template.CreatedAt,
	}

	return resp, nil
}
