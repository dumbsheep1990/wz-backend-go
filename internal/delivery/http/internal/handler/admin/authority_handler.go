package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/internal/application/admin/dto"
	"wz-backend-go/internal/application/admin/service"
	"wz-backend-go/internal/delivery/http/internal/types"
)

type AuthorityHandler struct {
	roleService *service.RoleApplicationService
}

func NewAuthorityHandler(roleService *service.RoleApplicationService) *AuthorityHandler {
	return &AuthorityHandler{
		roleService: roleService,
	}
}

// GetAuthorityList 获取角色列表
func (h *AuthorityHandler) GetAuthorityList(c *gin.Context) {
	// 调用应用服务获取角色列表
	result, err := h.roleService.ListRoles(c.Request.Context(), 1, 100) // 暂时使用固定分页大小，后续可优化
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取角色列表失败: "+err.Error())
		return
	}

	// 将DTO转换为API响应
	authorities := make([]types.SysAuthority, 0, len(result.Roles))
	for _, role := range result.Roles {
		authority := types.SysAuthority{
			AuthorityId:   role.ID,
			AuthorityName: role.Name,
			ParentId:      role.ParentID,
			DefaultRouter: "", // 可能需要从角色详情获取
			CreatedAt:     role.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     role.UpdatedAt.Format(time.RFC3339),
		}

		authorities = append(authorities, authority)
	}

	c.JSON(http.StatusOK, types.GetAuthorityListResponse{
		Code:    http.StatusOK,
		Message: "获取角色列表成功",
		Data:    authorities,
	})
}

// CreateAuthority 创建角色
func (h *AuthorityHandler) CreateAuthority(c *gin.Context) {
	var req types.CreateAuthorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务DTO
	createReq := dto.RoleCreateRequest{
		ID:          req.AuthorityId,
		Name:        req.AuthorityName,
		ParentID:    req.ParentId,
		Permissions: []string{}, // 默认无权限，后续设置
	}

	// 调用应用服务创建角色
	role, err := h.roleService.CreateRole(c.Request.Context(), createReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "创建角色失败: "+err.Error())
		return
	}

	// 构造响应
	authority := types.SysAuthority{
		AuthorityId:   role.ID,
		AuthorityName: role.Name,
		ParentId:      role.ParentID,
		CreatedAt:     role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     role.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, types.CreateAuthorityResponse{
		Code:    http.StatusOK,
		Message: "创建角色成功",
		Data:    authority,
	})
}

// UpdateAuthority 更新角色
func (h *AuthorityHandler) UpdateAuthority(c *gin.Context) {
	var req types.UpdateAuthorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务DTO
	updateReq := dto.RoleUpdateRequest{
		Name:     req.AuthorityName,
		ParentID: req.ParentId,
	}

	// 调用应用服务更新角色
	role, err := h.roleService.UpdateRole(c.Request.Context(), req.AuthorityId, updateReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "更新角色失败: "+err.Error())
		return
	}

	// 构造响应
	authority := types.SysAuthority{
		AuthorityId:   role.ID,
		AuthorityName: role.Name,
		ParentId:      role.ParentID,
		CreatedAt:     role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     role.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, types.UpdateAuthorityResponse{
		Code:    http.StatusOK,
		Message: "更新角色成功",
		Data:    authority,
	})
}

// DeleteAuthority 删除角色
func (h *AuthorityHandler) DeleteAuthority(c *gin.Context) {
	var req types.DeleteAuthorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用服务删除角色
	err := h.roleService.DeleteRole(c.Request.Context(), req.AuthorityId)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "删除角色失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "删除角色成功",
	})
}

// SetDataAuthority 设置数据权限
func (h *AuthorityHandler) SetDataAuthority(c *gin.Context) {
	var req types.SetDataAuthorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 数据权限可能需要在域层和应用层添加支持
	// 这里假设roleService有处理数据权限的方法
	err := h.roleService.SetDataAuthorities(c.Request.Context(), req.AuthorityId, req.DataAuthorityId)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "设置数据权限失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "设置数据权限成功",
	})
}

// GetCasbinPolicy 获取Casbin策略
func (h *AuthorityHandler) GetCasbinPolicy(c *gin.Context) {
	var req types.GetPolicyPathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取角色的权限
	permissions, err := h.roleService.GetPermissions(c.Request.Context(), req.AuthorityId)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取策略失败: "+err.Error())
		return
	}

	// 将权限转换为策略路径
	policies := make([]types.CasbinPolicyPath, 0, len(permissions))
	for _, perm := range permissions {
		parts := parsePermission(perm.Value)
		if len(parts) >= 2 {
			policies = append(policies, types.CasbinPolicyPath{
				Path:   parts[1],
				Method: parts[0],
			})
		}
	}

	c.JSON(http.StatusOK, types.GetPolicyPathResponse{
		Code:    http.StatusOK,
		Message: "获取策略成功",
		Data:    policies,
	})
}

// UpdateCasbinPolicy 更新Casbin策略
func (h *AuthorityHandler) UpdateCasbinPolicy(c *gin.Context) {
	var req types.UpdateCasbinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 将策略路径转换为权限字符串
	permissionStrs := make([]string, 0, len(req.CasbinInfo))
	for _, policy := range req.CasbinInfo {
		permissionStrs = append(permissionStrs, formatPermission(policy.Method, policy.Path))
	}

	// 设置角色的权限
	err := h.roleService.SetPermissions(c.Request.Context(), req.AuthorityId, permissionStrs)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "更新策略失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "更新策略成功",
	})
}

// 以下是辅助函数
// 在实际实现中，这些函数可能位于专门的工具包或领域服务中

// parsePermission 解析权限字符串为方法和路径
func parsePermission(permission string) []string {
	// 简单实现，实际中可能需要更复杂的解析逻辑
	return []string{"GET", "/api/v1/" + permission}
}

// formatPermission 格式化方法和路径为权限字符串
func formatPermission(method, path string) string {
	// 简单实现，实际中可能需要更复杂的格式化逻辑
	return method + ":" + path
}

// 以下是可能需要在RoleApplicationService中添加的方法

// SetDataAuthorities 设置数据权限范围
func (s *service.RoleApplicationService) SetDataAuthorities(ctx context.Context, roleID string, dataAuthorityIDs []string) error {
	// 这只是一个示例，实际实现应该在应用服务层完成
	panic("需要在RoleApplicationService中实现SetDataAuthorities方法")
}
