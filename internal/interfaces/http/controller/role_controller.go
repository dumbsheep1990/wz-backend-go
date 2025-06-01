package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/admin/dto"
	"wz-backend-go/internal/application/admin/service"
)

// RoleController handles HTTP requests related to role management
type RoleController struct {
	roleAppService *service.RoleApplicationService
}

// NewRoleController creates a new RoleController
func NewRoleController(roleAppService *service.RoleApplicationService) *RoleController {
	return &RoleController{
		roleAppService: roleAppService,
	}
}

// RegisterRoutes registers all routes for the role controller
func (c *RoleController) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	roleGroup := router.Group("/api/v1/admin/roles")
	roleGroup.Use(authMiddleware)
	{
		roleGroup.POST("", c.CreateRole)
		roleGroup.GET("", c.ListRoles)
		roleGroup.GET("/:id", c.GetRole)
		roleGroup.PUT("/:id", c.UpdateRole)
		roleGroup.DELETE("/:id", c.DeleteRole)
		
		// Permission management
		roleGroup.GET("/:id/permissions", c.GetRolePermissions)
		roleGroup.PUT("/:id/permissions", c.SetRolePermissions)
		roleGroup.POST("/:id/permissions", c.AddRolePermission)
		roleGroup.DELETE("/:id/permissions/:permission", c.RemoveRolePermission)
		
		// Available permissions
		roleGroup.GET("/permissions/available", c.GetAvailablePermissions)
		roleGroup.GET("/permissions/default", c.GetDefaultPermissions)
	}
}

// CreateRole handles role creation
// @Summary Create new role
// @Description Create a new role with permissions
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.RoleCreateRequest true "Role creation request"
// @Success 201 {object} dto.RoleDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/roles [post]
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req dto.RoleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	role, err := c.roleAppService.CreateRole(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create role: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, role)
}

// ListRoles handles listing roles with pagination
// @Summary List roles
// @Description Get a paginated list of roles
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number, default is 1"
// @Param pageSize query int false "Page size, default is 20"
// @Success 200 {object} dto.RoleListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/roles [get]
func (c *RoleController) ListRoles(ctx *gin.Context) {
	// Parse pagination parameters
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	response, err := c.roleAppService.ListRoles(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list roles: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetRole handles retrieving a specific role by ID
// @Summary Get role by ID
// @Description Get detailed information about a role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Success 200 {object} dto.RoleDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/roles/{id} [get]
func (c *RoleController) GetRole(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing role ID",
		})
		return
	}

	response, err := c.roleAppService.GetRoleDetail(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve role: " + err.Error(),
		})
		return
	}

	if response == nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Role not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateRole handles updating a specific role by ID
// @Summary Update role
// @Description Update a role's information and permissions
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Param request body dto.RoleUpdateRequest true "Role update request"
// @Success 200 {object} dto.RoleDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/roles/{id} [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing role ID",
		})
		return
	}

	var req dto.RoleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	response, err := c.roleAppService.UpdateRole(ctx, id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update role: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteRole handles deleting a specific role by ID
// @Summary Delete role
// @Description Delete a role
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/roles/{id} [delete]
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing role ID",
		})
		return
	}

	if err := c.roleAppService.DeleteRole(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete role: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "Role deleted successfully",
	})
}

// GetRolePermissions retrieves all permissions for a role
// @Summary Get role permissions
// @Description Get all permissions assigned to a role
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Success 200 {array} dto.PermissionDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/roles/{id}/permissions [get]
func (c *RoleController) GetRolePermissions(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing role ID",
		})
		return
	}

	permissions, err := c.roleAppService.GetPermissions(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve permissions: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, permissions)
}

// SetRolePermissions sets all permissions for a role
// @Summary Set role permissions
// @Description Replace all permissions for a role with a new set
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Param permissions body []string true "List of permission strings"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/roles/{id}/permissions [put]
func (c *RoleController) SetRolePermissions(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing role ID",
		})
		return
	}

	var permissions []string
	if err := ctx.ShouldBindJSON(&permissions); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	if err := c.roleAppService.SetPermissions(ctx, id, permissions); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to set permissions: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "Permissions updated successfully",
	})
}

// AddRolePermission adds a permission to a role
// @Summary Add role permission
// @Description Add a new permission to a role
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Param permission body string true "Permission string"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/roles/{id}/permissions [post]
func (c *RoleController) AddRolePermission(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing role ID",
		})
		return
	}

	var permission string
	if err := ctx.ShouldBindJSON(&permission); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	if err := c.roleAppService.AddPermission(ctx, id, permission); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to add permission: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "Permission added successfully",
	})
}

// RemoveRolePermission removes a permission from a role
// @Summary Remove role permission
// @Description Remove a permission from a role
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Param permission path string true "Permission string"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/roles/{id}/permissions/{permission} [delete]
func (c *RoleController) RemoveRolePermission(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing role ID",
		})
		return
	}

	permission := ctx.Param("permission")
	if permission == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing permission",
		})
		return
	}

	if err := c.roleAppService.RemovePermission(ctx, id, permission); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to remove permission: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "Permission removed successfully",
	})
}

// GetAvailablePermissions retrieves all available permissions in the system
// @Summary Get available permissions
// @Description Get all available permissions that can be assigned to roles
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} dto.PermissionDTO
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/admin/roles/permissions/available [get]
func (c *RoleController) GetAvailablePermissions(ctx *gin.Context) {
	permissions := c.roleAppService.GetAvailablePermissions()
	ctx.JSON(http.StatusOK, permissions)
}

// GetDefaultPermissions retrieves the default permissions for new roles
// @Summary Get default permissions
// @Description Get the default set of permissions for new roles
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} dto.PermissionDTO
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/admin/roles/permissions/default [get]
func (c *RoleController) GetDefaultPermissions(ctx *gin.Context) {
	permissions := c.roleAppService.GetDefaultPermissions()
	ctx.JSON(http.StatusOK, permissions)
}
