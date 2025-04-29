package logic

import (
	"context"

	"wz-backend-go/api/rpc/file"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 上传文件
func (l *UploadLogic) Upload(in *file.UploadRequest) (*file.UploadResponse, error) {
	// todo: add your logic here and delete this line

	return &file.UploadResponse{}, nil
}
