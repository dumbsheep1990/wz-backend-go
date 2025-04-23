package auth

import (
	"context"
	"fmt"
	"time"

	"wz-backend-go/internal/delivery/http/internal/middleware"
	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
	"wz-backend-go/internal/domain/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 验证用户名和密码
	// 这里应该调用用户服务验证用户登录
	// 模拟用户验证过程
	user := &model.User{
		ID:       1001,
		Username: req.Username,
		Role:     model.RolePersonalUser,
	}

	// 如果指定了租户ID，需要验证用户是否属于该租户
	var tenantID int64
	if req.TenantID > 0 {
		// 验证租户存在性
		tenant, err := l.svcCtx.TenantService.GetTenantByID(l.ctx, req.TenantID)
		if err != nil {
			return nil, fmt.Errorf("租户不存在: %v", err)
		}

		// 验证用户是否属于该租户
		hasAccess, role, err := l.svcCtx.TenantService.CheckUserInTenant(l.ctx, tenant.ID, user.ID)
		if err != nil {
			return nil, fmt.Errorf("验证用户租户关系失败: %v", err)
		}

		if !hasAccess {
			return nil, fmt.Errorf("用户不属于该租户")
		}

		// 更新用户角色为租户中的角色
		user.Role = model.UserRole(role)
		tenantID = tenant.ID
	}

	// 生成JWT令牌
	tokenPair, err := middleware.GenerateTokenPair(
		user.ID,
		user.Username,
		tenantID,
		user.Role,
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
	)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %v", err)
	}

	return &types.LoginResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt.Unix(),
		TokenType:    tokenPair.TokenType,
	}, nil
}
