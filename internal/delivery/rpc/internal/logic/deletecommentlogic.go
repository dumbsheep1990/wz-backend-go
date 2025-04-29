package logic

import (
	"context"

	"wz-backend-go/api/rpc/interaction"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除评论
func (l *DeleteCommentLogic) DeleteComment(in *interaction.DeleteCommentRequest) (*interaction.CommentResponse, error) {
	// todo: add your logic here and delete this line

	return &interaction.CommentResponse{}, nil
}
