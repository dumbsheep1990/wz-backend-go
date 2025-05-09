package logic

import (
	"context"
	"wz-backend-go/api/rpc/ad"
	"wz-backend-go/internal/admin/svc"
	"wz-backend-go/internal/admin/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetAdSpaceListLogic 获取广告位列表的逻辑处理器
type GetAdSpaceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetAdSpaceListLogic 创建获取广告位列表的逻辑处理器
func NewGetAdSpaceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdSpaceListLogic {
	return &GetAdSpaceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetAdSpaceList 获取广告位列表
func (l *GetAdSpaceListLogic) GetAdSpaceList(req *types.AdSpaceListReq) (*types.AdSpaceListResp, error) {
	// 调用广告服务获取广告位列表
	resp, err := l.svcCtx.AdClient.GetAdSpaceList(l.ctx, &ad.GetAdSpaceListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Position: req.Position,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	list := make([]types.AdSpaceDetail, 0, len(resp.List))
	for _, item := range resp.List {
		list = append(list, types.AdSpaceDetail{
			ID:          item.Id,
			Name:        item.Name,
			Position:    item.Position,
			Description: item.Description,
			Width:       item.Width,
			Height:      item.Height,
			Type:        item.Type,
			MaxAds:      item.MaxAds,
			Status:      item.Status,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	return &types.AdSpaceListResp{
		Total: resp.Total,
		List:  list,
	}, nil
}

// GetAdSpaceDetailLogic 获取广告位详情的逻辑处理器
type GetAdSpaceDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetAdSpaceDetailLogic 创建获取广告位详情的逻辑处理器
func NewGetAdSpaceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdSpaceDetailLogic {
	return &GetAdSpaceDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetAdSpaceDetail 获取广告位详情
func (l *GetAdSpaceDetailLogic) GetAdSpaceDetail(req *types.AdSpaceDetailReq) (*types.AdSpaceDetail, error) {
	// 调用广告服务获取广告位详情
	resp, err := l.svcCtx.AdClient.GetAdSpaceDetail(l.ctx, &ad.GetAdSpaceDetailReq{
		Id: req.ID,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.AdSpaceDetail{
		ID:          resp.Id,
		Name:        resp.Name,
		Position:    resp.Position,
		Description: resp.Description,
		Width:       resp.Width,
		Height:      resp.Height,
		Type:        resp.Type,
		MaxAds:      resp.MaxAds,
		Status:      resp.Status,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}

// CreateAdSpaceLogic 创建广告位的逻辑处理器
type CreateAdSpaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateAdSpaceLogic 创建创建广告位的逻辑处理器
func NewCreateAdSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdSpaceLogic {
	return &CreateAdSpaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateAdSpace 创建广告位
func (l *CreateAdSpaceLogic) CreateAdSpace(req *types.CreateAdSpaceReq) (*types.AdSpaceDetail, error) {
	// 调用广告服务创建广告位
	resp, err := l.svcCtx.AdClient.CreateAdSpace(l.ctx, &ad.CreateAdSpaceReq{
		Name:        req.Name,
		Position:    req.Position,
		Description: req.Description,
		Width:       req.Width,
		Height:      req.Height,
		Type:        req.Type,
		MaxAds:      req.MaxAds,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.AdSpaceDetail{
		ID:          resp.Id,
		Name:        resp.Name,
		Position:    resp.Position,
		Description: resp.Description,
		Width:       resp.Width,
		Height:      resp.Height,
		Type:        resp.Type,
		MaxAds:      resp.MaxAds,
		Status:      resp.Status,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}

// UpdateAdSpaceLogic 更新广告位的逻辑处理器
type UpdateAdSpaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateAdSpaceLogic 创建更新广告位的逻辑处理器
func NewUpdateAdSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdSpaceLogic {
	return &UpdateAdSpaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateAdSpace 更新广告位
func (l *UpdateAdSpaceLogic) UpdateAdSpace(req *types.UpdateAdSpaceReq) (*types.SuccessResp, error) {
	// 调用广告服务更新广告位
	resp, err := l.svcCtx.AdClient.UpdateAdSpace(l.ctx, &ad.UpdateAdSpaceReq{
		Id:          req.ID,
		Name:        req.Name,
		Position:    req.Position,
		Description: req.Description,
		Width:       req.Width,
		Height:      req.Height,
		Type:        req.Type,
		MaxAds:      req.MaxAds,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.SuccessResp{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// DeleteAdSpaceLogic 删除广告位的逻辑处理器
type DeleteAdSpaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteAdSpaceLogic 创建删除广告位的逻辑处理器
func NewDeleteAdSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAdSpaceLogic {
	return &DeleteAdSpaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteAdSpace 删除广告位
func (l *DeleteAdSpaceLogic) DeleteAdSpace(req *types.DeleteAdSpaceReq) (*types.SuccessResp, error) {
	// 调用广告服务删除广告位
	resp, err := l.svcCtx.AdClient.DeleteAdSpace(l.ctx, &ad.DeleteAdSpaceReq{
		Id: req.ID,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.SuccessResp{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}
