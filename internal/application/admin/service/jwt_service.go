package service

import "context"

// JWTService defines the interface for JWT token operations
type JWTService interface {
	// GenerateToken creates a new JWT token for an admin user
	GenerateToken(adminID int64, username string) (token string, expiresAt int64, err error)
	
	// ValidateToken validates a JWT token and returns the claims
	ValidateToken(token string) (adminID int64, username string, err error)
	
	// RevokeToken blacklists a token to prevent its further use
	RevokeToken(ctx context.Context, token string) error
	
	// IsTokenRevoked checks if a token has been revoked
	IsTokenRevoked(ctx context.Context, token string) (bool, error)
}
