package logic

import (
	"context"

	"wz-backend-go/api/rpc/file"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileLogic {
	return &GetFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取文件信息
func (l *GetFileLogic) GetFile(in *file.GetFileRequest) (*file.GetFileResponse, error) {
	// todo: add your logic here and delete this line

	return &file.GetFileResponse{}, nil
}
