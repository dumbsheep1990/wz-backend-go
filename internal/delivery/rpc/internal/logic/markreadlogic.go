package logic

import (
	"context"

	"wz-backend-go/api/rpc/notification"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 标记通知为已读
func (l *MarkReadLogic) MarkRead(in *notification.MarkReadRequest) (*notification.MarkReadResponse, error) {
	// todo: add your logic here and delete this line

	return &notification.MarkReadResponse{}, nil
}
