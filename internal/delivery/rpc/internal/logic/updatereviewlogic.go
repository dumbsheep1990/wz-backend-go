package logic

import (
	"context"

	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReviewLogic {
	return &UpdateReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateReviewLogic) UpdateReview(in *content.UpdateReviewRequest) (*content.UpdateReviewResponse, error) {
	// todo: add your logic here and delete this line

	return &content.UpdateReviewResponse{}, nil
}
