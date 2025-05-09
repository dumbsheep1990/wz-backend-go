package logic

import (
	"context"
	"wz-backend-go/api/rpc/user"
	"wz-backend-go/internal/admin/svc"
	"wz-backend-go/internal/admin/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetUserListLogic 获取用户列表的逻辑处理器
type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetUserListLogic 创建获取用户列表的逻辑处理器
func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetUserList 获取用户列表
func (l *GetUserListLogic) GetUserList(req *types.UserListReq) (*types.UserListResp, error) {
	// 调用用户服务获取用户列表
	resp, err := l.svcCtx.UserClient.GetUserList(l.ctx, &user.GetUserListReq{
		Page:      req.Page,
		PageSize:  req.PageSize,
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    req.Status,
		Role:      req.Role,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	list := make([]types.UserDetail, 0, len(resp.List))
	for _, item := range resp.List {
		list = append(list, types.UserDetail{
			ID:                item.Id,
			Username:          item.Username,
			Email:             item.Email,
			Phone:             item.Phone,
			Role:              item.Role,
			Status:            item.Status,
			IsVerified:        item.IsVerified,
			IsCompanyVerified: item.IsCompanyVerified,
			DefaultTenantID:   item.DefaultTenantID,
			CreatedAt:         item.CreatedAt,
			UpdatedAt:         item.UpdatedAt,
		})
	}

	return &types.UserListResp{
		Total: resp.Total,
		List:  list,
	}, nil
}

// GetUserDetailLogic 获取用户详情的逻辑处理器
type GetUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetUserDetailLogic 创建获取用户详情的逻辑处理器
func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailLogic {
	return &GetUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetUserDetail 获取用户详情
func (l *GetUserDetailLogic) GetUserDetail(req *types.UserDetailReq) (*types.UserDetail, error) {
	// 调用用户服务获取用户详情
	resp, err := l.svcCtx.UserClient.GetUserDetail(l.ctx, &user.GetUserDetailReq{
		Id: req.ID,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.UserDetail{
		ID:                resp.Id,
		Username:          resp.Username,
		Email:             resp.Email,
		Phone:             resp.Phone,
		Role:              resp.Role,
		Status:            resp.Status,
		IsVerified:        resp.IsVerified,
		IsCompanyVerified: resp.IsCompanyVerified,
		DefaultTenantID:   resp.DefaultTenantID,
		CreatedAt:         resp.CreatedAt,
		UpdatedAt:         resp.UpdatedAt,
	}, nil
}

// CreateUserLogic 创建用户的逻辑处理器
type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateUserLogic 创建创建用户的逻辑处理器
func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateUser 创建用户
func (l *CreateUserLogic) CreateUser(req *types.AdminCreateUserReq) (*types.AdminCreateUserResp, error) {
	// 调用用户服务创建用户
	resp, err := l.svcCtx.UserClient.CreateUser(l.ctx, &user.CreateUserReq{
		Username:        req.Username,
		Password:        req.Password,
		Email:           req.Email,
		Phone:           req.Phone,
		Role:            req.Role,
		Status:          req.Status,
		DefaultTenantID: req.DefaultTenantID,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.AdminCreateUserResp{
		ID:       resp.Id,
		Username: resp.Username,
	}, nil
}

// UpdateUserLogic 更新用户的逻辑处理器
type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateUserLogic 创建更新用户的逻辑处理器
func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateUser 更新用户
func (l *UpdateUserLogic) UpdateUser(req *types.AdminUpdateUserReq) (*types.SuccessResp, error) {
	// 调用用户服务更新用户
	resp, err := l.svcCtx.UserClient.UpdateUser(l.ctx, &user.UpdateUserReq{
		Id:              req.ID,
		Username:        req.Username,
		Password:        req.Password,
		Email:           req.Email,
		Phone:           req.Phone,
		Role:            req.Role,
		Status:          req.Status,
		DefaultTenantID: req.DefaultTenantID,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.SuccessResp{
		Success: resp.Success,
	}, nil
}

// DeleteUserLogic 删除用户的逻辑处理器
type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteUserLogic 创建删除用户的逻辑处理器
func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteUser 删除用户
func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (*types.SuccessResp, error) {
	// 调用用户服务删除用户
	resp, err := l.svcCtx.UserClient.DeleteUser(l.ctx, &user.DeleteUserReq{
		Id: req.ID,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	return &types.SuccessResp{
		Success: resp.Success,
	}, nil
}
