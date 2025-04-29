package logic

import (
	"context"

	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteReviewLogic {
	return &DeleteReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteReviewLogic) DeleteReview(in *content.DeleteReviewRequest) (*content.DeleteReviewResponse, error) {
	// todo: add your logic here and delete this line

	return &content.DeleteReviewResponse{}, nil
}
