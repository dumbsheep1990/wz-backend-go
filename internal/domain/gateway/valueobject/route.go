package valueobject

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// RouteID represents a unique identifier for a route
type RouteID struct {
	value string
}

// NewRouteID creates a new RouteID
func NewRouteID(id string) (RouteID, error) {
	if id == "" {
		return RouteID{}, errors.New("route ID cannot be empty")
	}
	return RouteID{value: id}, nil
}

// String returns the string representation of the RouteID
func (id RouteID) String() string {
	return id.value
}

// Path represents a URL path pattern for routing
type Path struct {
	value string
}

// NewPath creates a new Path value object
func NewPath(path string) (Path, error) {
	if path == "" {
		return Path{}, errors.New("path cannot be empty")
	}
	
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	
	return Path{value: path}, nil
}

// String returns the string representation of the Path
func (p Path) String() string {
	return p.value
}

// Matches checks if this path matches the given request path
func (p Path) Matches(requestPath string) bool {
	// Simple implementation, can be extended for path parameters and wildcards
	pathPattern := p.value
	
	// Handle trailing wildcard
	if strings.HasSuffix(pathPattern, "/*") {
		basePattern := strings.TrimSuffix(pathPattern, "/*")
		return strings.HasPrefix(requestPath, basePattern)
	}
	
	return pathPattern == requestPath
}

// TargetURL represents a URL to forward requests to
type TargetURL struct {
	value *url.URL
}

// NewTargetURL creates a new TargetURL value object
func NewTargetURL(rawURL string) (TargetURL, error) {
	if rawURL == "" {
		return TargetURL{}, errors.New("target URL cannot be empty")
	}
	
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return TargetURL{}, fmt.Errorf("invalid URL format: %w", err)
	}
	
	return TargetURL{value: parsedURL}, nil
}

// URL returns the underlying url.URL
func (t TargetURL) URL() *url.URL {
	return t.value
}

// String returns the string representation of the TargetURL
func (t TargetURL) String() string {
	return t.value.String()
}

// ServiceName represents the name of a service
type ServiceName struct {
	value string
}

// NewServiceName creates a new ServiceName value object
func NewServiceName(name string) (ServiceName, error) {
	if name == "" {
		return ServiceName{}, errors.New("service name cannot be empty")
	}
	
	return ServiceName{value: name}, nil
}

// String returns the string representation of the ServiceName
func (n ServiceName) String() string {
	return n.value
}
