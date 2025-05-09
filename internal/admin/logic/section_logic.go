package logic

import (
	"context"
	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/admin/svc"
	"wz-backend-go/internal/admin/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetSectionListLogic 获取板块列表的逻辑处理器
type GetSectionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetSectionListLogic 创建获取板块列表的逻辑处理器
func NewGetSectionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSectionListLogic {
	return &GetSectionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetSectionList 获取板块列表
func (l *GetSectionListLogic) GetSectionList(req *types.SectionListReq) (*types.SectionListResp, error) {
	// 调用内容服务获取板块列表
	resp, err := l.svcCtx.ContentClient.GetSectionList(l.ctx, &content.GetSectionListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	list := make([]types.SectionDetail, 0, len(resp.List))
	for _, item := range resp.List {
		list = append(list, types.SectionDetail{
			ID:          item.Id,
			Name:        item.Name,
			Code:        item.Code,
			Description: item.Description,
			IconUrl:     item.IconUrl,
			BannerUrl:   item.BannerUrl,
			SortOrder:   item.SortOrder,
			Status:      item.Status,
			ShowInHome:  item.ShowInHome,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	return &types.SectionListResp{
		Total: resp.Total,
		List:  list,
	}, nil
}

// GetSectionDetailLogic 获取板块详情的逻辑处理器
type GetSectionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetSectionDetailLogic 创建获取板块详情的逻辑处理器
func NewGetSectionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSectionDetailLogic {
	return &GetSectionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetSectionDetail 获取板块详情
func (l *GetSectionDetailLogic) GetSectionDetail(req *types.SectionDetailReq) (*types.SectionDetail, error) {
	// 调用内容服务获取板块详情
	resp, err := l.svcCtx.ContentClient.GetSectionDetail(l.ctx, &content.GetSectionDetailReq{
		Id: req.ID,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.SectionDetail{
		ID:          resp.Id,
		Name:        resp.Name,
		Code:        resp.Code,
		Description: resp.Description,
		IconUrl:     resp.IconUrl,
		BannerUrl:   resp.BannerUrl,
		SortOrder:   resp.SortOrder,
		Status:      resp.Status,
		ShowInHome:  resp.ShowInHome,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}

// CreateSectionLogic 创建板块的逻辑处理器
type CreateSectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateSectionLogic 创建创建板块的逻辑处理器
func NewCreateSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSectionLogic {
	return &CreateSectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateSection 创建板块
func (l *CreateSectionLogic) CreateSection(req *types.CreateSectionReq) (*types.SectionDetail, error) {
	// 调用内容服务创建板块
	resp, err := l.svcCtx.ContentClient.CreateSection(l.ctx, &content.CreateSectionReq{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IconUrl:     req.IconUrl,
		BannerUrl:   req.BannerUrl,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		ShowInHome:  req.ShowInHome,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.SectionDetail{
		ID:          resp.Id,
		Name:        resp.Name,
		Code:        resp.Code,
		Description: resp.Description,
		IconUrl:     resp.IconUrl,
		BannerUrl:   resp.BannerUrl,
		SortOrder:   resp.SortOrder,
		Status:      resp.Status,
		ShowInHome:  resp.ShowInHome,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}

// UpdateSectionLogic 更新板块的逻辑处理器
type UpdateSectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateSectionLogic 创建更新板块的逻辑处理器
func NewUpdateSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSectionLogic {
	return &UpdateSectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateSection 更新板块
func (l *UpdateSectionLogic) UpdateSection(req *types.UpdateSectionReq) (*types.SuccessResp, error) {
	// 调用内容服务更新板块
	resp, err := l.svcCtx.ContentClient.UpdateSection(l.ctx, &content.UpdateSectionReq{
		Id:          req.ID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IconUrl:     req.IconUrl,
		BannerUrl:   req.BannerUrl,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		ShowInHome:  req.ShowInHome,
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

// DeleteSectionLogic 删除板块的逻辑处理器
type DeleteSectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteSectionLogic 创建删除板块的逻辑处理器
func NewDeleteSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSectionLogic {
	return &DeleteSectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteSection 删除板块
func (l *DeleteSectionLogic) DeleteSection(req *types.DeleteSectionReq) (*types.SuccessResp, error) {
	// 调用内容服务删除板块
	resp, err := l.svcCtx.ContentClient.DeleteSection(l.ctx, &content.DeleteSectionReq{
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
