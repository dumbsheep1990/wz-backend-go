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

// MenuHandler 菜单处理程序
type MenuHandler struct {
	menuService *service.MenuApplicationService
}

// NewMenuHandler 创建菜单处理程序
func NewMenuHandler(menuService *service.MenuApplicationService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

// GetMenuList 获取菜单列表
func (h *MenuHandler) GetMenuList(c *gin.Context) {
	// 调用应用服务获取菜单列表
	result, err := h.menuService.ListMenus(c.Request.Context())
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取菜单列表失败: "+err.Error())
		return
	}

	// 转换为API响应
	menuList := make([]types.MenuInfo, 0, len(result.Menus))
	for _, menu := range result.Menus {
		menuList = append(menuList, convertMenuDTOToMenuInfo(menu))
	}

	c.JSON(http.StatusOK, types.GetMenuListResponse{
		Code:    http.StatusOK,
		Message: "获取菜单列表成功",
		Data:    menuList,
	})
}

// GetMenusByAuthority 获取指定角色的菜单
func (h *MenuHandler) GetMenusByAuthority(c *gin.Context) {
	var req types.GetMenuAuthorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用服务获取角色菜单
	result, err := h.menuService.GetMenusByAuthority(c.Request.Context(), req.AuthorityId)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取角色菜单失败: "+err.Error())
		return
	}

	// 转换为API响应
	menuList := make([]types.MenuInfo, 0, len(result.Menus))
	for _, menu := range result.Menus {
		menuList = append(menuList, convertMenuDTOToMenuInfo(menu))
	}

	c.JSON(http.StatusOK, types.GetMenuAuthorityResponse{
		Code:    http.StatusOK,
		Message: "获取角色菜单成功",
		Data:    menuList,
	})
}

// AddMenu 添加菜单
func (h *MenuHandler) AddMenu(c *gin.Context) {
	var req types.MenuInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务DTO
	createReq := dto.MenuCreateRequest{
		ParentID:  req.ParentId,
		Path:      req.Path,
		Name:      req.Name,
		Hidden:    req.Hidden,
		Component: req.Component,
		Sort:      req.Sort,
		Meta: dto.MenuMeta{
			Title:            req.Meta.Title,
			Icon:             req.Meta.Icon,
			KeepAlive:        req.Meta.KeepAlive,
			DefaultMenu:      req.Meta.DefaultMenu,
			CloseTab:         req.Meta.CloseTab,
			CollapsibleWidth: req.Meta.CollapsibleWidth,
		},
	}

	// 调用应用服务创建菜单
	menu, err := h.menuService.CreateMenu(c.Request.Context(), createReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "添加菜单失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "添加菜单成功",
		Data:    menu.ID,
	})
}

// UpdateMenu 更新菜单
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	var req types.MenuInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务DTO
	updateReq := dto.MenuUpdateRequest{
		ID:        req.Id,
		ParentID:  req.ParentId,
		Path:      req.Path,
		Name:      req.Name,
		Hidden:    req.Hidden,
		Component: req.Component,
		Sort:      req.Sort,
		Meta: dto.MenuMeta{
			Title:            req.Meta.Title,
			Icon:             req.Meta.Icon,
			KeepAlive:        req.Meta.KeepAlive,
			DefaultMenu:      req.Meta.DefaultMenu,
			CloseTab:         req.Meta.CloseTab,
			CollapsibleWidth: req.Meta.CollapsibleWidth,
		},
	}

	// 调用应用服务更新菜单
	_, err := h.menuService.UpdateMenu(c.Request.Context(), updateReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "更新菜单失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "更新菜单成功",
	})
}

// DeleteMenu 删除菜单
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	var req types.DeleteMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用服务删除菜单
	err := h.menuService.DeleteMenu(c.Request.Context(), req.Id)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "删除菜单失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "删除菜单成功",
	})
}

// AddMenuAuthority 添加菜单权限
func (h *MenuHandler) AddMenuAuthority(c *gin.Context) {
	var req types.AddMenuAuthorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 提取菜单IDs
	menuIDs := make([]int64, 0, len(req.Menus))
	for _, menu := range req.Menus {
		menuIDs = append(menuIDs, menu.Id)
	}

	// 调用应用服务设置角色菜单
	menuAuthReq := dto.MenuAuthorityRequest{
		AuthorityID: req.AuthorityId,
		MenuIDs:     menuIDs,
	}

	err := h.menuService.AddMenuAuthority(c.Request.Context(), menuAuthReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "设置菜单权限失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "设置菜单权限成功",
	})
}

// 辅助函数：将DTO转换为API响应类型
func convertMenuDTOToMenuInfo(menu dto.MenuDTO) types.MenuInfo {
	menuInfo := types.MenuInfo{
		Id:        menu.ID,
		CreatedAt: menu.CreatedAt.Format(time.RFC3339),
		UpdatedAt: menu.UpdatedAt.Format(time.RFC3339),
		ParentId:  menu.ParentID,
		Path:      menu.Path,
		Name:      menu.Name,
		Hidden:    menu.Hidden,
		Component: menu.Component,
		Sort:      menu.Sort,
		Meta: types.MenuMeta{
			Title:            menu.Meta.Title,
			Icon:             menu.Meta.Icon,
			KeepAlive:        menu.Meta.KeepAlive,
			DefaultMenu:      menu.Meta.DefaultMenu,
			CloseTab:         menu.Meta.CloseTab,
			CollapsibleWidth: menu.Meta.CollapsibleWidth,
		},
	}

	// 递归转换子菜单
	if len(menu.Children) > 0 {
		menuInfo.Children = make([]types.MenuInfo, 0, len(menu.Children))
		for _, childMenu := range menu.Children {
			menuInfo.Children = append(menuInfo.Children, convertMenuDTOToMenuInfo(childMenu))
		}
	}

	return menuInfo
}
