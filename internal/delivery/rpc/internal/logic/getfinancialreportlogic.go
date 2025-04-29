package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFinancialReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFinancialReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFinancialReportLogic {
	return &GetFinancialReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 报表管理
func (l *GetFinancialReportLogic) GetFinancialReport(in *trade.GetFinancialReportRequest) (*trade.GetFinancialReportResponse, error) {
	// todo: add your logic here and delete this line

	return &trade.GetFinancialReportResponse{}, nil
}
