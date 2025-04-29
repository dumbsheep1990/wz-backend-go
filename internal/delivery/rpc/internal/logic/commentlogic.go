package logic

import (
	"context"

	"wz-backend-go/api/rpc/interaction"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLogic {
	return &CommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发表评论
func (l *CommentLogic) Comment(in *interaction.CommentRequest) (*interaction.CommentResponse, error) {
	// todo: add your logic here and delete this line

	return &interaction.CommentResponse{}, nil
}
