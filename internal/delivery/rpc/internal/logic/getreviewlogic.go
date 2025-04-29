package logic

import (
	"context"

	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReviewLogic {
	return &GetReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetReviewLogic) GetReview(in *content.GetReviewRequest) (*content.GetReviewResponse, error) {
	// todo: add your logic here and delete this line

	return &content.GetReviewResponse{}, nil
}
