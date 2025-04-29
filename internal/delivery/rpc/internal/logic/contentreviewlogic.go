package logic

import (
	"context"

	"wz-backend-go/api/rpc/ai"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ContentReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewContentReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentReviewLogic {
	return &ContentReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 内容审核
func (l *ContentReviewLogic) ContentReview(in *ai.ContentReviewRequest) (*ai.ContentReviewResponse, error) {
	// todo: add your logic here and delete this line

	return &ai.ContentReviewResponse{}, nil
}
