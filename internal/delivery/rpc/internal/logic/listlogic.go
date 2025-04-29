package logic

import (
	"context"

	"wz-backend-go/api/rpc/notification"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取通知列表
func (l *ListLogic) List(in *notification.ListRequest) (*notification.ListResponse, error) {
	// todo: add your logic here and delete this line

	return &notification.ListResponse{}, nil
}
