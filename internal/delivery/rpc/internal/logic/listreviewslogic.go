package logic

import (
	"context"

	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListReviewsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListReviewsLogic {
	return &ListReviewsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListReviewsLogic) ListReviews(in *content.ListReviewsRequest) (*content.ListReviewsResponse, error) {
	// todo: add your logic here and delete this line

	return &content.ListReviewsResponse{}, nil
}
