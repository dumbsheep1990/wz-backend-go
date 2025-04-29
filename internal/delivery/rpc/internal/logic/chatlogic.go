package logic

import (
	"context"

	"wz-backend-go/api/rpc/ai"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 客服对话
func (l *ChatLogic) Chat(in *ai.ChatRequest) (*ai.ChatResponse, error) {
	// todo: add your logic here and delete this line

	return &ai.ChatResponse{}, nil
}
