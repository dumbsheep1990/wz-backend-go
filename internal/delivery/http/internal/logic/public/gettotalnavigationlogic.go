package public

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"wz-backend-go/internal/domain/model"
	"wz-backend-go/internal/delivery/http/internal/svc"
)

type GetTotalNavigationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTotalNavigationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTotalNavigationLogic {
	return &GetTotalNavigationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetTotalNavigation 获取总站全局导航
func (l *GetTotalNavigationLogic) GetTotalNavigation(lang string) (*model.NavigationResponse, error) {
	// 这里可以根据语言参数获取不同语言的导航
	// 简单示例，实际应该从数据库或配置中获取
	categories := []model.NavigationCategory{
		{
			ID:   1,
			Name: "所有市场",
			URL:  "/market",
		},
		{
			ID:   2,
			Name: "所有房产",
			URL:  "/realestate",
		},
		{
			ID:   3,
			Name: "所有职位",
			URL:  "/jobs",
		},
		{
			ID:   4,
			Name: "所有课程",
			URL:  "/courses",
		},
	}

	// 如果有其他语言的支持，可以在这里添加语言判断
	if lang != "zh-CN" {
		// 处理其他语言...
	}

	return &model.NavigationResponse{
		Categories: categories,
	}, nil
}
