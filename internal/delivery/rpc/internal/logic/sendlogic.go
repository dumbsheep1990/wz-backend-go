package logic

import (
	"context"

	"wz-backend-go/api/rpc/notification"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogic {
	return &SendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送通知
func (l *SendLogic) Send(in *notification.SendRequest) (*notification.SendResponse, error) {
	// todo: add your logic here and delete this line

	return &notification.SendResponse{}, nil
}
