package valueobject

import (
	"errors"
	"regexp"
)

// AuthType represents the type of authentication for a route
type AuthType string

const (
	// NoAuth represents routes that don't require authentication
	NoAuth AuthType = "none"
	
	// JWTAuth represents routes that require JWT authentication
	JWTAuth AuthType = "jwt"
	
	// APIKeyAuth represents routes that require API key authentication
	APIKeyAuth AuthType = "apikey"
	
	// OAuthAuth represents routes that require OAuth authentication
	OAuthAuth AuthType = "oauth"
)

// IsValid checks if the auth type is valid
func (a AuthType) IsValid() bool {
	switch a {
	case NoAuth, JWTAuth, APIKeyAuth, OAuthAuth:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (a AuthType) String() string {
	return string(a)
}

// TokenString represents a JWT or API key token
type TokenString struct {
	value string
}

// NewTokenString creates a new TokenString
func NewTokenString(token string) (TokenString, error) {
	if token == "" {
		return TokenString{}, errors.New("token cannot be empty")
	}
	return TokenString{value: token}, nil
}

// String returns the string representation of the token
func (t TokenString) String() string {
	return t.value
}

// APIKey represents an API key for authentication
type APIKey struct {
	key    string
	secret string
}

// NewAPIKey creates a new APIKey
func NewAPIKey(key, secret string) (APIKey, error) {
	if key == "" {
		return APIKey{}, errors.New("API key cannot be empty")
	}
	
	// API key format validation
	keyPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]{8,32}$`)
	if !keyPattern.MatchString(key) {
		return APIKey{}, errors.New("invalid API key format")
	}
	
	if secret != "" {
		// Secret format validation (if provided)
		secretPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]{16,64}$`)
		if !secretPattern.MatchString(secret) {
			return APIKey{}, errors.New("invalid API secret format")
		}
	}
	
	return APIKey{key: key, secret: secret}, nil
}

// Key returns the API key
func (a APIKey) Key() string {
	return a.key
}

// Secret returns the API secret
func (a APIKey) Secret() string {
	return a.secret
}

// HasSecret checks if the API key has a secret
func (a APIKey) HasSecret() bool {
	return a.secret != ""
}

// AuthScope represents a scope for authentication
type AuthScope struct {
	value string
}

// NewAuthScope creates a new AuthScope
func NewAuthScope(scope string) (AuthScope, error) {
	if scope == "" {
		return AuthScope{}, errors.New("auth scope cannot be empty")
	}
	return AuthScope{value: scope}, nil
}

// String returns the string representation of the scope
func (s AuthScope) String() string {
	return s.value
}
