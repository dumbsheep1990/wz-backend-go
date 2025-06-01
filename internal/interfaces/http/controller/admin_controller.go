package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/admin/dto"
	"wz-backend-go/internal/application/admin/service"
)

// AdminController handles HTTP requests related to admin management
type AdminController struct {
	adminAppService *service.AdminApplicationService
}

// NewAdminController creates a new AdminController
func NewAdminController(adminAppService *service.AdminApplicationService) *AdminController {
	return &AdminController{
		adminAppService: adminAppService,
	}
}

// RegisterRoutes registers all routes for the admin controller
func (c *AdminController) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	adminGroup := router.Group("/api/v1/admin")
	{
		// Public routes
		adminGroup.POST("/login", c.Login)

		// Protected routes
		protected := adminGroup.Group("/")
		protected.Use(authMiddleware)
		{
			// Admin management
			protected.GET("/profile", c.GetCurrentAdmin)
			protected.PUT("/profile/password", c.ChangePassword)
			
			// Admin CRUD operations
			protected.POST("/users", c.CreateAdmin)
			protected.GET("/users", c.ListAdmins)
			protected.GET("/users/:id", c.GetAdmin)
			protected.PUT("/users/:id", c.UpdateAdmin)
			protected.DELETE("/users/:id", c.DeleteAdmin)
		}
	}
}

// Login handles admin authentication
// @Summary Admin login
// @Description Authenticate admin and return JWT token
// @Tags admin
// @Accept json
// @Produce json
// @Param request body dto.AdminLoginRequest true "Login credentials"
// @Success 200 {object} dto.AdminLoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/admin/login [post]
func (c *AdminController) Login(ctx *gin.Context) {
	var req dto.AdminLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Get client information for login tracking
	ip := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	response, err := c.adminAppService.Login(ctx, req, ip, userAgent)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Authentication failed: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetCurrentAdmin returns the currently authenticated admin
// @Summary Get current admin profile
// @Description Get detailed information about the currently authenticated admin
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.AdminDetailResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/admin/profile [get]
func (c *AdminController) GetCurrentAdmin(ctx *gin.Context) {
	// Get admin ID from context (set by auth middleware)
	adminID, exists := ctx.Get("AdminID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized: Admin ID not found in context",
		})
		return
	}

	id, ok := adminID.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error: Invalid admin ID type",
		})
		return
	}

	response, err := c.adminAppService.GetCurrentAdmin(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve admin profile: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ChangePassword handles password change for the current admin
// @Summary Change admin password
// @Description Change the password for the currently authenticated admin
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.AdminPasswordChangeRequest true "Password change request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/profile/password [put]
func (c *AdminController) ChangePassword(ctx *gin.Context) {
	var req dto.AdminPasswordChangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Get admin ID from context
	adminID, exists := ctx.Get("AdminID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized: Admin ID not found in context",
		})
		return
	}

	id, ok := adminID.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error: Invalid admin ID type",
		})
		return
	}

	if err := c.adminAppService.ChangePassword(ctx, id, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to change password: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "Password changed successfully",
	})
}

// CreateAdmin handles admin creation
// @Summary Create new admin
// @Description Create a new admin account
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.AdminCreateRequest true "Admin creation request"
// @Success 201 {object} dto.AdminDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/users [post]
func (c *AdminController) CreateAdmin(ctx *gin.Context) {
	var req dto.AdminCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	admin, err := c.adminAppService.CreateAdmin(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create admin: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, admin)
}

// ListAdmins handles listing admins with pagination and filtering
// @Summary List admins
// @Description Get a paginated list of admins with optional filtering
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number, default is 1"
// @Param pageSize query int false "Page size, default is 20"
// @Param username query string false "Filter by username"
// @Param status query int false "Filter by status"
// @Param roleId query string false "Filter by role ID"
// @Success 200 {object} dto.AdminListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/users [get]
func (c *AdminController) ListAdmins(ctx *gin.Context) {
	// Parse pagination parameters
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Parse filter parameters
	filters := make(map[string]interface{})

	if username := ctx.Query("username"); username != "" {
		filters["username"] = username
	}

	if statusStr := ctx.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			filters["status"] = status
		}
	}

	if roleID := ctx.Query("roleId"); roleID != "" {
		filters["role_id"] = roleID
	}

	response, err := c.adminAppService.ListAdmins(ctx, page, pageSize, filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list admins: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetAdmin handles retrieving a specific admin by ID
// @Summary Get admin by ID
// @Description Get detailed information about an admin by their ID
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Admin ID"
// @Success 200 {object} dto.AdminDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/users/{id} [get]
func (c *AdminController) GetAdmin(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid admin ID format",
		})
		return
	}

	response, err := c.adminAppService.GetAdminDetail(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve admin: " + err.Error(),
		})
		return
	}

	if response == nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Admin not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateAdmin handles updating a specific admin by ID
// @Summary Update admin
// @Description Update an admin's information
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Admin ID"
// @Param request body dto.AdminUpdateRequest true "Admin update request"
// @Success 200 {object} dto.AdminDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/users/{id} [put]
func (c *AdminController) UpdateAdmin(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid admin ID format",
		})
		return
	}

	var req dto.AdminUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	response, err := c.adminAppService.UpdateAdmin(ctx, id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update admin: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteAdmin handles deleting a specific admin by ID
// @Summary Delete admin
// @Description Delete an admin account
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Admin ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/users/{id} [delete]
func (c *AdminController) DeleteAdmin(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid admin ID format",
		})
		return
	}

	// Get current admin ID to prevent self-deletion
	currentAdminID, exists := ctx.Get("AdminID")
	if exists {
		if currentID, ok := currentAdminID.(int64); ok && currentID == id {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot delete your own account",
			})
			return
		}
	}

	if err := c.adminAppService.DeleteAdmin(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete admin: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "Admin deleted successfully",
	})
}
