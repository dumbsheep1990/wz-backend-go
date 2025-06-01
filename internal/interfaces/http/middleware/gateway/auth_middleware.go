package gateway

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"wz-backend-go/internal/domain/gateway/valueobject"
	"wz-backend-go/internal/application/gateway/dto"
	"wz-backend-go/internal/domain/gateway/entity"
)

// AuthMiddleware handles authentication for gateway routes
type AuthMiddleware struct {
	jwtSecret string
	apiKeys   map[string]string // key -> scope mapping
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		apiKeys:   make(map[string]string),
	}
}

// AddAPIKey adds an API key to the middleware
func (m *AuthMiddleware) AddAPIKey(key, scope string) {
	m.apiKeys[key] = scope
}

// GetAuthType returns a handler function to check auth requirements
func (m *AuthMiddleware) GetAuthType() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the route from context (set by route finder middleware)
		routeInterface, exists := c.Get("route")
		if !exists {
			c.Next()
			return
		}
		
		route, ok := routeInterface.(*entity.Route)
		if !ok {
			c.Next()
			return
		}
		
		// Store auth type in context
		c.Set("auth_type", route.AuthType().String())
		c.Next()
	}
}

// HandleAuth returns a handler function to check authentication
func (m *AuthMiddleware) HandleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if auth check is needed
		authTypeInterface, exists := c.Get("auth_type")
		if !exists {
			c.Next()
			return
		}
		
		authTypeStr, ok := authTypeInterface.(string)
		if !ok {
			c.Next()
			return
		}
		
		authType := valueobject.AuthTypeFromString(authTypeStr)
		
		// No authentication required
		if authType == valueobject.NoAuth {
			c.Next()
			return
		}
		
		// JWT authentication
		if authType == valueobject.JWTAuth {
			if err := m.validateJWT(c); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "Invalid or missing JWT token",
					Details: err.Error(),
				})
				return
			}
			c.Next()
			return
		}
		
		// API Key authentication
		if authType == valueobject.APIKeyAuth {
			if err := m.validateAPIKey(c); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "Invalid or missing API key",
					Details: err.Error(),
				})
				return
			}
			c.Next()
			return
		}
		
		// OAuth authentication
		if authType == valueobject.OAuthAuth {
			if err := m.validateOAuth(c); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "Invalid or missing OAuth token",
					Details: err.Error(),
				})
				return
			}
			c.Next()
			return
		}
		
		// Unsupported auth type
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Unsupported authentication type",
		})
	}
}

// validateJWT validates a JWT token
func (m *AuthMiddleware) validateJWT(c *gin.Context) error {
	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return fmt.Errorf("authorization header is required")
	}
	
	// Check if header has Bearer prefix
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return fmt.Errorf("authorization header must be in the format 'Bearer {token}'")
	}
	
	// Extract token
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	
	// Parse token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		return []byte(m.jwtSecret), nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}
	
	// Validate token
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	
	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}
	
	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return fmt.Errorf("token expired")
		}
	}
	
	// Store user ID in context
	if userID, ok := claims["sub"].(string); ok {
		c.Set("user_id", userID)
	}
	
	// Store all claims in context
	c.Set("claims", claims)
	
	return nil
}

// validateAPIKey validates an API key
func (m *AuthMiddleware) validateAPIKey(c *gin.Context) error {
	// Get API key from header
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		return fmt.Errorf("X-API-Key header is required")
	}
	
	// Check if API key exists
	scope, exists := m.apiKeys[apiKey]
	if !exists {
		return fmt.Errorf("invalid API key")
	}
	
	// Store scope in context
	c.Set("api_key_scope", scope)
	
	return nil
}

// validateOAuth validates an OAuth token
func (m *AuthMiddleware) validateOAuth(c *gin.Context) error {
	// This is a placeholder for OAuth validation
	// In a real implementation, this would validate with an OAuth provider
	
	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return fmt.Errorf("authorization header is required")
	}
	
	// Check if header has Bearer prefix
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return fmt.Errorf("authorization header must be in the format 'Bearer {token}'")
	}
	
	// Extract token
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	
	// For now, just check if token is not empty
	if tokenStr == "" {
		return fmt.Errorf("token cannot be empty")
	}
	
	// In a real implementation, this would validate the token with the OAuth provider
	// and extract claims like scope, user ID, etc.
	
	return nil
}
