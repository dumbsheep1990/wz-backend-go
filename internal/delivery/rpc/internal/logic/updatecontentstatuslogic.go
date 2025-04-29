package logic

import (
	"context"

	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateContentStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateContentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateContentStatusLogic {
	return &UpdateContentStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 内容状态管理
func (l *UpdateContentStatusLogic) UpdateContentStatus(in *content.UpdateContentStatusRequest) (*content.UpdateContentStatusResponse, error) {
	// todo: add your logic here and delete this line

	return &content.UpdateContentStatusResponse{}, nil
}
