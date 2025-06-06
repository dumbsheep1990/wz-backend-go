package dto

import (
	"time"

	"wz-backend-go/internal/domain/admin/entity"
)

// MenuDTO 菜单DTO
type MenuDTO struct {
	ID        int64     `json:"id"`
	ParentID  int64     `json:"parentId"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Hidden    bool      `json:"hidden"`
	Component string    `json:"component"`
	Sort      int       `json:"sort"`
	Meta      MenuMeta  `json:"meta"`
	Children  []MenuDTO `json:"children,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MenuMeta 菜单元数据
type MenuMeta struct {
	Title            string `json:"title"`
	Icon             string `json:"icon"`
	KeepAlive        bool   `json:"keepAlive"`
	DefaultMenu      bool   `json:"defaultMenu"`
	CloseTab         bool   `json:"closeTab"`
	CollapsibleWidth int    `json:"collapsibleWidth,omitempty"`
}

// MenuCreateRequest 创建菜单请求
type MenuCreateRequest struct {
	ParentID  int64    `json:"parentId"`
	Path      string   `json:"path"`
	Name      string   `json:"name"`
	Hidden    bool     `json:"hidden"`
	Component string   `json:"component"`
	Sort      int      `json:"sort"`
	Meta      MenuMeta `json:"meta"`
}

// MenuUpdateRequest 更新菜单请求
type MenuUpdateRequest struct {
	ID        int64    `json:"id"`
	ParentID  int64    `json:"parentId"`
	Path      string   `json:"path"`
	Name      string   `json:"name"`
	Hidden    bool     `json:"hidden"`
	Component string   `json:"component"`
	Sort      int      `json:"sort"`
	Meta      MenuMeta `json:"meta"`
}

// MenuListResponse 菜单列表响应
type MenuListResponse struct {
	Menus []MenuDTO `json:"menus"`
}

// MenuAuthorityRequest 获取或设置角色菜单请求
type MenuAuthorityRequest struct {
	AuthorityID string  `json:"authorityId"`
	MenuIDs     []int64 `json:"menuIds,omitempty"`
}

// MenuAuthorityResponse 获取角色菜单响应
type MenuAuthorityResponse struct {
	Menus []MenuDTO `json:"menus"`
}

// MapMenuToDTO 将领域实体映射为DTO
func MapMenuToDTO(menu *entity.Menu) MenuDTO {
	dto := MenuDTO{
		ID:        menu.ID(),
		ParentID:  menu.ParentID(),
		Path:      menu.Path(),
		Name:      menu.Name(),
		Hidden:    menu.Hidden(),
		Component: menu.Component(),
		Sort:      menu.Sort(),
		Meta: MenuMeta{
			Title:            menu.Meta().Title(),
			Icon:             menu.Meta().Icon(),
			KeepAlive:        menu.Meta().KeepAlive(),
			DefaultMenu:      menu.Meta().DefaultMenu(),
			CloseTab:         menu.Meta().CloseTab(),
			CollapsibleWidth: menu.Meta().CollapsibleWidth(),
		},
		CreatedAt: menu.CreatedAt(),
		UpdatedAt: menu.UpdatedAt(),
	}
	
	if len(menu.Children()) > 0 {
		dto.Children = make([]MenuDTO, 0, len(menu.Children()))
		for _, child := range menu.Children() {
			dto.Children = append(dto.Children, MapMenuToDTO(child))
		}
	}
	
	return dto
}
