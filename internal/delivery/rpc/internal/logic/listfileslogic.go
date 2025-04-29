package logic

import (
	"context"

	"wz-backend-go/api/rpc/file"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFilesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFilesLogic {
	return &ListFilesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取文件列表
func (l *ListFilesLogic) ListFiles(in *file.ListFilesRequest) (*file.ListFilesResponse, error) {
	// todo: add your logic here and delete this line

	return &file.ListFilesResponse{}, nil
}
