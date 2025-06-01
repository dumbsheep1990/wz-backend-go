package valueobject

import (
	"errors"
	"time"
)

// RateLimitType represents the type of rate limiting
type RateLimitType string

const (
	// NoRateLimit represents no rate limiting
	NoRateLimit RateLimitType = "none"
	
	// IPRateLimit represents rate limiting by IP address
	IPRateLimit RateLimitType = "ip"
	
	// TokenRateLimit represents rate limiting by authentication token
	TokenRateLimit RateLimitType = "token"
	
	// ServiceRateLimit represents rate limiting across all requests to a service
	ServiceRateLimit RateLimitType = "service"
)

// IsValid checks if the rate limit type is valid
func (r RateLimitType) IsValid() bool {
	switch r {
	case NoRateLimit, IPRateLimit, TokenRateLimit, ServiceRateLimit:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (r RateLimitType) String() string {
	return string(r)
}

// RateLimit represents a rate limit configuration
type RateLimit struct {
	limitType RateLimitType
	requests  int
	duration  time.Duration
}

// NewRateLimit creates a new RateLimit
func NewRateLimit(limitType RateLimitType, requests int, durationSeconds int) (RateLimit, error) {
	if !limitType.IsValid() {
		return RateLimit{}, errors.New("invalid rate limit type")
	}
	
	if limitType != NoRateLimit {
		if requests <= 0 {
			return RateLimit{}, errors.New("requests must be positive")
		}
		
		if durationSeconds <= 0 {
			return RateLimit{}, errors.New("duration must be positive")
		}
	}
	
	return RateLimit{
		limitType: limitType,
		requests:  requests,
		duration:  time.Duration(durationSeconds) * time.Second,
	}, nil
}

// Type returns the rate limit type
func (r RateLimit) Type() RateLimitType {
	return r.limitType
}

// Requests returns the maximum number of requests allowed
func (r RateLimit) Requests() int {
	return r.requests
}

// Duration returns the time window for the rate limit
func (r RateLimit) Duration() time.Duration {
	return r.duration
}

// IsEnabled checks if rate limiting is enabled
func (r RateLimit) IsEnabled() bool {
	return r.limitType != NoRateLimit
}

// RequestsPerSecond returns the rate limit as requests per second
func (r RateLimit) RequestsPerSecond() float64 {
	if !r.IsEnabled() {
		return 0
	}
	return float64(r.requests) / r.duration.Seconds()
}

// String returns a human-readable representation of the rate limit
func (r RateLimit) String() string {
	if !r.IsEnabled() {
		return "No rate limit"
	}
	return string(r.limitType) + " rate limit: " + 
		string(r.requests) + " requests per " + 
		r.duration.String()
}
