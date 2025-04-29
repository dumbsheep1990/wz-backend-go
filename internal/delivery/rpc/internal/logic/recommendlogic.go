package logic

import (
	"context"

	"wz-backend-go/api/rpc/ai"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecommendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendLogic {
	return &RecommendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取推荐内容
func (l *RecommendLogic) Recommend(in *ai.RecommendRequest) (*ai.RecommendResponse, error) {
	// todo: add your logic here and delete this line

	return &ai.RecommendResponse{}, nil
}
