package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/admin/service"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// AuthMiddleware is a middleware for JWT authentication
type AuthMiddleware struct {
	jwtService      service.JWTService
	adminAppService *service.AdminApplicationService
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(jwtService service.JWTService, adminAppService *service.AdminApplicationService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:      jwtService,
		adminAppService: adminAppService,
	}
}

// AdminAuthHandler returns a Gin middleware for admin authentication
func (m *AuthMiddleware) AdminAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid authorization format. Expected Bearer {token}",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		adminID, username, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid or expired token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Set admin info in context
		c.Set("adminID", adminID)
		c.Set("username", username)

		// Continue to the next handler
		c.Next()
	}
}

// PermissionMiddleware returns a Gin middleware for permission-based authorization
func (m *AuthMiddleware) PermissionMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get admin ID from context (set by AdminAuthHandler)
		adminIDValue, exists := c.Get("adminID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized: admin not authenticated",
			})
			c.Abort()
			return
		}

		adminID, ok := adminIDValue.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Invalid admin ID type",
			})
			c.Abort()
			return
		}

		// Create admin ID value object
		adminIDObj, err := valueobject.NewAdminID(adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Invalid admin ID: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Convert permission string to Permission value object
		permissionObj, err := valueobject.NewPermission(permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Invalid permission: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Check if admin has permission
		hasPermission, err := m.adminAppService.HasPermission(c, adminIDObj, permissionObj)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Error checking permission: " + err.Error(),
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "Forbidden: insufficient permissions",
			})
			c.Abort()
			return
		}

		// Continue to the next handler
		c.Next()
	}
}

// RequireRoleMiddleware returns a Gin middleware that restricts access to specific roles
func (m *AuthMiddleware) RequireRoleMiddleware(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get admin ID from context (set by AdminAuthHandler)
		adminIDValue, exists := c.Get("adminID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized: admin not authenticated",
			})
			c.Abort()
			return
		}

		adminID, ok := adminIDValue.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Invalid admin ID type",
			})
			c.Abort()
			return
		}

		// Create admin ID value object
		adminIDObj, err := valueobject.NewAdminID(adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Invalid admin ID: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Get admin with role details
		admin, err := m.adminAppService.GetAdminDetail(c, adminIDObj)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Error retrieving admin details: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Check if admin has the required role
		if admin.Role.Name != roleName {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "Forbidden: requires role " + roleName,
			})
			c.Abort()
			return
		}

		// Continue to the next handler
		c.Next()
	}
}
