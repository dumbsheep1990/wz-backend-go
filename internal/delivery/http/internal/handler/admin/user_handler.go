package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/internal/application/admin/dto"
	"wz-backend-go/internal/application/admin/service"
	"wz-backend-go/internal/delivery/http/internal/types"
)

type UserHandler struct {
	adminService *service.AdminApplicationService
	roleService  *service.RoleApplicationService
}

func NewUserHandler(adminService *service.AdminApplicationService, roleService *service.RoleApplicationService) *UserHandler {
	return &UserHandler{
		adminService: adminService,
		roleService:  roleService,
	}
}

// GetUserList 获取用户列表
func (h *UserHandler) GetUserList(c *gin.Context) {
	var req types.GetUserListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 构造过滤条件
	filters := map[string]interface{}{}
	if req.Username != "" {
		filters["username"] = req.Username
	}
	if req.Status != 0 {
		filters["status"] = req.Status
	}
	if req.AuthorityId != "" {
		filters["role_id"] = req.AuthorityId
	}
	if req.CreatedAfter != "" {
		filters["created_after"] = req.CreatedAfter
	}
	if req.CreatedBefore != "" {
		filters["created_before"] = req.CreatedBefore
	}

	// 调用应用服务获取用户列表
	result, err := h.adminService.ListAdmins(c.Request.Context(), req.Page, req.PageSize, filters)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取用户列表失败: "+err.Error())
		return
	}

	// 将DTO转换为API响应
	userList := make([]types.UserInfo, 0, len(result.Admins))
	for _, admin := range result.Admins {
		// 获取用户的角色信息
		roleInfo, err := h.roleService.GetRole(c.Request.Context(), admin.RoleID)
		if err != nil {
			// 角色信息获取失败，使用默认空值
			roleInfo = &dto.RoleDTO{
				ID:   admin.RoleID,
				Name: "未知角色",
			}
		}

		userList = append(userList, types.UserInfo{
			Id:          admin.ID,
			UUID:        admin.UUID,
			Username:    admin.Username,
			Avatar:      admin.Avatar,
			Status:      int(admin.Status),
			AuthorityId: admin.RoleID,
			CreatedAt:   admin.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   admin.UpdatedAt.Format(time.RFC3339),
			LastLoginAt: admin.LastLoginAt.Format(time.RFC3339),
			Authorities: []types.Authority{
				{
					AuthorityId:   roleInfo.ID,
					AuthorityName: roleInfo.Name,
				},
			},
		})
	}

	response := types.GetUserListResponse{
		Code:    http.StatusOK,
		Message: "获取用户列表成功",
		Data: types.UserPageData{
			List:     userList,
			Total:    result.Total,
			Page:     result.Page,
			PageSize: result.PageSize,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetUserInfo 获取用户详情
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		// 如果没有指定ID，则获取当前登录用户信息
		// 从JWT中获取用户ID
		adminID, exists := c.Get("admin_id")
		if !exists {
			httpx.ErrorCtx(c, http.StatusUnauthorized, "未登录")
			return
		}
		idStr = adminID.(string)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 获取用户详细信息
	adminDetail, err := h.adminService.GetAdminDetail(c.Request.Context(), id)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取用户详情失败: "+err.Error())
		return
	}

	// 将角色信息转换为前端需要的格式
	authorities := make([]types.Authority, 0, len(adminDetail.Roles))
	for _, role := range adminDetail.Roles {
		authorities = append(authorities, types.Authority{
			AuthorityId:   role.ID,
			AuthorityName: role.Name,
			ParentId:      role.ParentID,
			DefaultRouter: role.DefaultRouter,
			CreatedAt:     role.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     role.UpdatedAt.Format(time.RFC3339),
		})
	}

	// 构造响应
	userInfo := types.UserInfo{
		Id:          adminDetail.Admin.ID,
		UUID:        adminDetail.Admin.UUID,
		Username:    adminDetail.Admin.Username,
		Avatar:      adminDetail.Admin.Avatar,
		Status:      int(adminDetail.Admin.Status),
		AuthorityId: adminDetail.Admin.RoleID,
		CreatedAt:   adminDetail.Admin.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   adminDetail.Admin.UpdatedAt.Format(time.RFC3339),
		LastLoginAt: adminDetail.Admin.LastLoginAt.Format(time.RFC3339),
		Authorities: authorities,
	}

	c.JSON(http.StatusOK, types.GetUserInfoResponse{
		Code:    http.StatusOK,
		Message: "获取用户详情成功",
		Data:    userInfo,
	})
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req types.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务DTO
	createReq := dto.AdminCreateRequest{
		Username: req.Username,
		Password: req.Password,
		RoleID:   req.AuthorityId,
	}

	// 调用应用服务创建用户
	admin, err := h.adminService.CreateAdmin(c.Request.Context(), createReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "创建用户失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "创建用户成功",
		Data:    admin.ID,
	})
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req types.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务DTO
	updateReq := dto.AdminUpdateRequest{
		Username: req.Username,
		RoleID:   req.AuthorityId,
	}

	if req.Status != 0 {
		status := int32(req.Status)
		updateReq.Status = &status
	}

	// 调用应用服务更新用户
	_, err := h.adminService.UpdateAdmin(c.Request.Context(), req.Id, updateReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "更新用户失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "更新用户成功",
	})
}

// UpdateSelfInfo 更新自身信息
func (h *UserHandler) UpdateSelfInfo(c *gin.Context) {
	var req types.UpdateSelfInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 从JWT中获取用户ID
	adminID, exists := c.Get("admin_id")
	if !exists {
		httpx.ErrorCtx(c, http.StatusUnauthorized, "未登录")
		return
	}

	id, _ := strconv.ParseInt(adminID.(string), 10, 64)

	// 转换为应用服务DTO
	updateReq := dto.AdminUpdateRequest{
		// 只更新允许用户自己修改的字段
		Avatar: req.Avatar,
		// 其他个人设置可以保存到额外的用户设置表或通过其他服务处理
	}

	// 调用应用服务更新用户
	_, err := h.adminService.UpdateAdmin(c.Request.Context(), id, updateReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "更新个人信息失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "更新个人信息成功",
	})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	var req types.DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用服务删除用户
	err := h.adminService.DeleteAdmin(c.Request.Context(), req.Id)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "删除用户失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "删除用户成功",
	})
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req types.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 从JWT中获取用户ID
	adminID, exists := c.Get("admin_id")
	if !exists {
		httpx.ErrorCtx(c, http.StatusUnauthorized, "未登录")
		return
	}

	id, _ := strconv.ParseInt(adminID.(string), 10, 64)

	// 转换为应用服务DTO
	passwordReq := dto.AdminPasswordChangeRequest{
		OldPassword: req.Password,
		NewPassword: req.NewPassword,
	}

	// 调用应用服务修改密码
	err := h.adminService.ChangePassword(c.Request.Context(), id, passwordReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "修改密码失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "修改密码成功",
	})
}

// ResetPassword 重置密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req types.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务DTO
	passwordReq := dto.AdminPasswordResetRequest{
		ID:          req.Id,
		NewPassword: req.Password,
	}

	// 这里假设adminService有一个ResetPassword方法，如果没有，需要在应用层添加
	err := h.adminService.ResetPassword(c.Request.Context(), passwordReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "重置密码失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "重置密码成功",
	})
}

// SetUserAuthority 设置用户权限
func (h *UserHandler) SetUserAuthority(c *gin.Context) {
	var req types.SetUserAuthorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务DTO
	updateReq := dto.AdminUpdateRequest{
		RoleID: req.AuthorityId,
	}

	// 调用应用服务更新用户角色
	_, err := h.adminService.UpdateAdmin(c.Request.Context(), req.UserId, updateReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "设置用户角色失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "设置用户角色成功",
	})
}

// SetUserAuthorities 设置用户权限组
func (h *UserHandler) SetUserAuthorities(c *gin.Context) {
	var req types.SetUserAuthoritiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 多角色功能可能需要在域层和应用层添加支持
	// 这里暂时使用第一个角色作为主角色
	if len(req.AuthorityIds) == 0 {
		httpx.ErrorCtx(c, http.StatusBadRequest, "权限组不能为空")
		return
	}

	// 转换为应用服务DTO
	updateReq := dto.AdminUpdateRequest{
		RoleID: req.AuthorityIds[0], // 使用第一个角色作为主角色
	}

	// 调用应用服务更新用户角色
	_, err := h.adminService.UpdateAdmin(c.Request.Context(), req.UserId, updateReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "设置用户权限组失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "设置用户权限组成功",
	})
}

// 在实际实现中，可能需要添加以下辅助方法或扩展adminService
// ResetPassword方法未在AdminApplicationService中看到，需要添加
func (s *service.AdminApplicationService) ResetPassword(ctx context.Context, req dto.AdminPasswordResetRequest) error {
	// 这只是一个示例，实际实现应该在应用服务层完成
	panic("需要在AdminApplicationService中实现ResetPassword方法")
}
