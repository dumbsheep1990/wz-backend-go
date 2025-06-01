package gateway

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/gateway/dto"
	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

// RateLimiter is a simple in-memory rate limiter
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time // key -> list of request timestamps
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
	}
}

// Clean removes expired timestamps
func (r *RateLimiter) Clean() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for key, timestamps := range r.requests {
		var validTimestamps []time.Time
		for _, ts := range timestamps {
			if now.Sub(ts) < time.Hour {
				validTimestamps = append(validTimestamps, ts)
			}
		}
		if len(validTimestamps) == 0 {
			delete(r.requests, key)
		} else {
			r.requests[key] = validTimestamps
		}
	}
}

// Check checks if a request is allowed
func (r *RateLimiter) Check(key string, limit int, duration time.Duration) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	timestamps, exists := r.requests[key]
	if !exists {
		r.requests[key] = []time.Time{now}
		return true
	}

	// Filter timestamps that are within the duration
	var validTimestamps []time.Time
	for _, ts := range timestamps {
		if now.Sub(ts) < duration {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	// Check if request count is within limit
	if len(validTimestamps) < limit {
		r.requests[key] = append(validTimestamps, now)
		return true
	}

	// Update timestamps
	r.requests[key] = validTimestamps
	return false
}

// RateLimitMiddleware handles rate limiting for gateway routes
type RateLimitMiddleware struct {
	limiter *RateLimiter
}

// NewRateLimitMiddleware creates a new RateLimitMiddleware
func NewRateLimitMiddleware() *RateLimitMiddleware {
	limiter := NewRateLimiter()

	// Start a cleanup goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiter.Clean()
		}
	}()

	return &RateLimitMiddleware{
		limiter: limiter,
	}
}

// HandleRateLimit returns a handler function to check rate limits
func (m *RateLimitMiddleware) HandleRateLimit() gin.HandlerFunc {
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

		// Check if rate limiting is enabled
		if !route.RateLimit().IsEnabled() {
			c.Next()
			return
		}

		// Get rate limit details
		rateLimitType := route.RateLimit().Type()
		requests := route.RateLimit().Requests()
		duration := route.RateLimit().Duration()

		// Generate key based on rate limit type
		var key string
		switch rateLimitType {
		case valueobject.IPRateLimit:
			// Get client IP
			clientIP := c.ClientIP()
			key = fmt.Sprintf("ip:%s:%s", clientIP, route.ID().String())
		case valueobject.TokenRateLimit:
			// Get token from context (set by auth middleware)
			userID, exists := c.Get("user_id")
			if !exists {
				// If no user ID, fallback to IP
				clientIP := c.ClientIP()
				key = fmt.Sprintf("ip:%s:%s", clientIP, route.ID().String())
			} else {
				key = fmt.Sprintf("token:%s:%s", userID, route.ID().String())
			}
		case valueobject.ServiceRateLimit:
			// Limit based on service
			key = fmt.Sprintf("service:%s", route.ServiceName().String())
		default:
			// No rate limiting
			c.Next()
			return
		}

		// Check rate limit
		if !m.limiter.Check(key, requests, duration) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Code:    "RATE_LIMIT_EXCEEDED",
				Message: "Rate limit exceeded",
				Details: fmt.Sprintf("Limit: %d requests per %s", requests, duration),
			})
			return
		}

		c.Next()
	}
}
