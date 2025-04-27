package logic

import (
	"context"

	"wz-backend-go/internal/delivery/rpc/interaction/interaction"
	"wz-backend-go/internal/delivery/rpc/interaction/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportLogic {
	return &ReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 举报
func (l *ReportLogic) Report(in *interaction.ReportRequest) (*interaction.ReportResponse, error) {
	// todo: add your logic here and delete this line

	return &interaction.ReportResponse{}, nil
}
