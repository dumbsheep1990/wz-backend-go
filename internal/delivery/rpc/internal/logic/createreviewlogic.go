package logic

import (
	"context"

	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReviewLogic {
	return &CreateReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 评论管理
func (l *CreateReviewLogic) CreateReview(in *content.CreateReviewRequest) (*content.CreateReviewResponse, error) {
	// todo: add your logic here and delete this line

	return &content.CreateReviewResponse{}, nil
}
