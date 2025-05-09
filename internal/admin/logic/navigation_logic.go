package logic

import (
	"context"
	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/admin/svc"
	"wz-backend-go/internal/admin/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetNavigationListLogic 获取导航列表的逻辑处理器
type GetNavigationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetNavigationListLogic 创建获取导航列表的逻辑处理器
func NewGetNavigationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNavigationListLogic {
	return &GetNavigationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetNavigationList 获取导航列表
func (l *GetNavigationListLogic) GetNavigationList(req *types.NavigationListReq) (*types.NavigationListResp, error) {
	// 调用内容服务获取导航列表
	resp, err := l.svcCtx.ContentClient.GetNavigationList(l.ctx, &content.GetNavigationListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Type:     req.Type,
		Status:   req.Status,
		ParentId: req.ParentId,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	list := make([]types.NavigationDetail, 0, len(resp.List))
	for _, item := range resp.List {
		list = append(list, types.NavigationDetail{
			ID:         item.Id,
			Name:       item.Name,
			Type:       item.Type,
			Url:        item.Url,
			Target:     item.Target,
			IconUrl:    item.IconUrl,
			ParentId:   item.ParentId,
			SectionId:  item.SectionId,
			CategoryId: item.CategoryId,
			SortOrder:  item.SortOrder,
			Status:     item.Status,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		})
	}

	return &types.NavigationListResp{
		Total: resp.Total,
		List:  list,
	}, nil
}

// GetNavigationDetailLogic 获取导航详情的逻辑处理器
type GetNavigationDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetNavigationDetailLogic 创建获取导航详情的逻辑处理器
func NewGetNavigationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNavigationDetailLogic {
	return &GetNavigationDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetNavigationDetail 获取导航详情
func (l *GetNavigationDetailLogic) GetNavigationDetail(req *types.NavigationDetailReq) (*types.NavigationDetail, error) {
	// 调用内容服务获取导航详情
	resp, err := l.svcCtx.ContentClient.GetNavigationDetail(l.ctx, &content.GetNavigationDetailReq{
		Id: req.ID,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.NavigationDetail{
		ID:         resp.Id,
		Name:       resp.Name,
		Type:       resp.Type,
		Url:        resp.Url,
		Target:     resp.Target,
		IconUrl:    resp.IconUrl,
		ParentId:   resp.ParentId,
		SectionId:  resp.SectionId,
		CategoryId: resp.CategoryId,
		SortOrder:  resp.SortOrder,
		Status:     resp.Status,
		CreatedAt:  resp.CreatedAt,
		UpdatedAt:  resp.UpdatedAt,
	}, nil
}

// CreateNavigationLogic 创建导航的逻辑处理器
type CreateNavigationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateNavigationLogic 创建创建导航的逻辑处理器
func NewCreateNavigationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNavigationLogic {
	return &CreateNavigationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateNavigation 创建导航
func (l *CreateNavigationLogic) CreateNavigation(req *types.CreateNavigationReq) (*types.NavigationDetail, error) {
	// 调用内容服务创建导航
	resp, err := l.svcCtx.ContentClient.CreateNavigation(l.ctx, &content.CreateNavigationReq{
		Name:       req.Name,
		Type:       req.Type,
		Url:        req.Url,
		Target:     req.Target,
		IconUrl:    req.IconUrl,
		ParentId:   req.ParentId,
		SectionId:  req.SectionId,
		CategoryId: req.CategoryId,
		SortOrder:  req.SortOrder,
		Status:     req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.NavigationDetail{
		ID:         resp.Id,
		Name:       resp.Name,
		Type:       resp.Type,
		Url:        resp.Url,
		Target:     resp.Target,
		IconUrl:    resp.IconUrl,
		ParentId:   resp.ParentId,
		SectionId:  resp.SectionId,
		CategoryId: resp.CategoryId,
		SortOrder:  resp.SortOrder,
		Status:     resp.Status,
		CreatedAt:  resp.CreatedAt,
		UpdatedAt:  resp.UpdatedAt,
	}, nil
}

// UpdateNavigationLogic 更新导航的逻辑处理器
type UpdateNavigationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateNavigationLogic 创建更新导航的逻辑处理器
func NewUpdateNavigationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNavigationLogic {
	return &UpdateNavigationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateNavigation 更新导航
func (l *UpdateNavigationLogic) UpdateNavigation(req *types.UpdateNavigationReq) (*types.SuccessResp, error) {
	// 调用内容服务更新导航
	resp, err := l.svcCtx.ContentClient.UpdateNavigation(l.ctx, &content.UpdateNavigationReq{
		Id:         req.ID,
		Name:       req.Name,
		Type:       req.Type,
		Url:        req.Url,
		Target:     req.Target,
		IconUrl:    req.IconUrl,
		ParentId:   req.ParentId,
		SectionId:  req.SectionId,
		CategoryId: req.CategoryId,
		SortOrder:  req.SortOrder,
		Status:     req.Status,
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

// DeleteNavigationLogic 删除导航的逻辑处理器
type DeleteNavigationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteNavigationLogic 创建删除导航的逻辑处理器
func NewDeleteNavigationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNavigationLogic {
	return &DeleteNavigationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteNavigation 删除导航
func (l *DeleteNavigationLogic) DeleteNavigation(req *types.DeleteNavigationReq) (*types.SuccessResp, error) {
	// 调用内容服务删除导航
	resp, err := l.svcCtx.ContentClient.DeleteNavigation(l.ctx, &content.DeleteNavigationReq{
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
