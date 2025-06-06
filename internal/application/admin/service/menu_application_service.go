package service

import (
	"context"
	"errors"

	"wz-backend-go/internal/application/admin/dto"
	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/repository"
	domainService "wz-backend-go/internal/domain/admin/service"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// MenuApplicationService 菜单应用服务
type MenuApplicationService struct {
	menuRepo      repository.MenuRepository
	roleRepo      repository.RoleRepository
	menuDomainSvc *domainService.MenuDomainService
}

// NewMenuApplicationService 创建菜单应用服务
func NewMenuApplicationService(
	menuRepo repository.MenuRepository,
	roleRepo repository.RoleRepository,
	menuDomainSvc *domainService.MenuDomainService,
) *MenuApplicationService {
	return &MenuApplicationService{
		menuRepo:      menuRepo,
		roleRepo:      roleRepo,
		menuDomainSvc: menuDomainSvc,
	}
}

// CreateMenu 创建菜单
func (s *MenuApplicationService) CreateMenu(
	ctx context.Context,
	req dto.MenuCreateRequest,
) (*dto.MenuDTO, error) {
	// 创建菜单元数据值对象
	meta := valueobject.NewMenuMeta(
		req.Meta.Title,
		req.Meta.Icon,
		req.Meta.KeepAlive,
		req.Meta.DefaultMenu,
		req.Meta.CloseTab,
		req.Meta.CollapsibleWidth,
	)

	// 创建菜单实体
	menu, err := s.menuDomainSvc.CreateMenu(
		ctx,
		req.ParentID,
		req.Path,
		req.Name,
		req.Hidden,
		req.Component,
		req.Sort,
		meta,
	)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	menuDTO := dto.MapMenuToDTO(menu)
	return &menuDTO, nil
}

// UpdateMenu 更新菜单
func (s *MenuApplicationService) UpdateMenu(
	ctx context.Context,
	req dto.MenuUpdateRequest,
) (*dto.MenuDTO, error) {
	// 查找菜单
	menu, err := s.menuRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, errors.New("菜单不存在")
	}

	// 创建菜单元数据值对象
	meta := valueobject.NewMenuMeta(
		req.Meta.Title,
		req.Meta.Icon,
		req.Meta.KeepAlive,
		req.Meta.DefaultMenu,
		req.Meta.CloseTab,
		req.Meta.CollapsibleWidth,
	)

	// 更新菜单
	err = s.menuDomainSvc.UpdateMenu(
		ctx,
		menu.ID(),
		req.ParentID,
		req.Path,
		req.Name,
		req.Hidden,
		req.Component,
		req.Sort,
		meta,
	)
	if err != nil {
		return nil, err
	}

	// 获取更新后的菜单
	updatedMenu, err := s.menuRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	menuDTO := dto.MapMenuToDTO(updatedMenu)
	return &menuDTO, nil
}

// DeleteMenu 删除菜单
func (s *MenuApplicationService) DeleteMenu(ctx context.Context, id int64) error {
	return s.menuDomainSvc.DeleteMenu(ctx, id)
}

// GetMenu 获取菜单
func (s *MenuApplicationService) GetMenu(ctx context.Context, id int64) (*dto.MenuDTO, error) {
	// 查找菜单
	menu, err := s.menuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, errors.New("菜单不存在")
	}

	// 转换为DTO
	menuDTO := dto.MapMenuToDTO(menu)
	return &menuDTO, nil
}

// ListMenus 获取菜单列表
func (s *MenuApplicationService) ListMenus(ctx context.Context) (*dto.MenuListResponse, error) {
	// 获取所有菜单
	menus, err := s.menuRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	menuTree, err := s.buildMenuTree(menus)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	menuDTOs := make([]dto.MenuDTO, 0, len(menuTree))
	for _, menu := range menuTree {
		menuDTOs = append(menuDTOs, dto.MapMenuToDTO(menu))
	}

	return &dto.MenuListResponse{
		Menus: menuDTOs,
	}, nil
}

// GetMenusByAuthority 获取指定角色的菜单
func (s *MenuApplicationService) GetMenusByAuthority(ctx context.Context, authorityID string) (*dto.MenuAuthorityResponse, error) {
	// 验证角色ID
	roleID, err := valueobject.NewRoleID(authorityID)
	if err != nil {
		return nil, err
	}

	// 检查角色是否存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("角色不存在")
	}

	// 获取角色的菜单
	menus, err := s.menuRepo.FindByAuthorityID(ctx, authorityID)
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	menuTree, err := s.buildMenuTree(menus)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	menuDTOs := make([]dto.MenuDTO, 0, len(menuTree))
	for _, menu := range menuTree {
		menuDTOs = append(menuDTOs, dto.MapMenuToDTO(menu))
	}

	return &dto.MenuAuthorityResponse{
		Menus: menuDTOs,
	}, nil
}

// AddMenuAuthority 设置角色菜单权限
func (s *MenuApplicationService) AddMenuAuthority(ctx context.Context, req dto.MenuAuthorityRequest) error {
	// 验证角色ID
	roleID, err := valueobject.NewRoleID(req.AuthorityID)
	if err != nil {
		return err
	}

	// 检查角色是否存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("角色不存在")
	}

	// 调用领域服务设置角色菜单
	return s.menuDomainSvc.SetRoleMenus(ctx, req.AuthorityID, req.MenuIDs)
}

// buildMenuTree 构建菜单树
func (s *MenuApplicationService) buildMenuTree(menus []*entity.Menu) ([]*entity.Menu, error) {
	// 实现菜单树构建逻辑
	// 这里简单实现，实际中可能需要更复杂的树构建算法
	menuMap := make(map[int64]*entity.Menu)
	var rootMenus []*entity.Menu

	// 第一遍遍历，建立ID到菜单的映射
	for _, menu := range menus {
		menuMap[menu.ID()] = menu
	}

	// 第二遍遍历，建立父子关系
	for _, menu := range menus {
		if menu.ParentID() == 0 {
			// 根菜单
			rootMenus = append(rootMenus, menu)
		} else {
			// 子菜单
			parent, exists := menuMap[menu.ParentID()]
			if !exists {
				continue // 父菜单不存在，跳过
			}
			parent.AddChild(menu)
		}
	}

	return rootMenus, nil
}
