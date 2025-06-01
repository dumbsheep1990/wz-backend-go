package valueobject

import (
	"errors"
	"strings"
)

// Permission represents a permission for an admin role
type Permission string

// ResourceType represents the type of resource a permission applies to
type ResourceType string

// ActionType represents the type of action a permission allows
type ActionType string

// Common resource types
const (
	ResourceUser        ResourceType = "user"
	ResourceTenant      ResourceType = "tenant"
	ResourceContent     ResourceType = "content"
	ResourceTrade       ResourceType = "trade"
	ResourceSystem      ResourceType = "system"
	ResourceCategory    ResourceType = "category" 
	ResourceMenu        ResourceType = "menu"
	ResourceApi         ResourceType = "api"
	ResourceDictionary  ResourceType = "dictionary"
)

// Common action types
const (
	ActionCreate ActionType = "create"
	ActionRead   ActionType = "read"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
	ActionList   ActionType = "list"
	ActionManage ActionType = "manage" // Full access
	ActionExecute ActionType = "execute" // For APIs and operations
)

// NewPermission creates a new Permission value object
func NewPermission(resource ResourceType, action ActionType) (Permission, error) {
	if strings.TrimSpace(string(resource)) == "" {
		return "", errors.New("resource type cannot be empty")
	}
	if strings.TrimSpace(string(action)) == "" {
		return "", errors.New("action type cannot be empty")
	}
	
	permission := Permission(string(resource) + ":" + string(action))
	return permission, nil
}

// MustNewPermission creates a new Permission and panics if invalid
func MustNewPermission(resource ResourceType, action ActionType) Permission {
	p, err := NewPermission(resource, action)
	if err != nil {
		panic(err)
	}
	return p
}

// ParsePermission parses a permission string into a Permission object
func ParsePermission(permString string) (Permission, error) {
	if strings.TrimSpace(permString) == "" {
		return "", errors.New("permission string cannot be empty")
	}
	
	parts := strings.Split(permString, ":")
	if len(parts) != 2 {
		return "", errors.New("invalid permission format, expected 'resource:action'")
	}
	
	resource := ResourceType(parts[0])
	action := ActionType(parts[1])
	
	return NewPermission(resource, action)
}

// GetResourceType returns the resource type component of the permission
func (p Permission) GetResourceType() ResourceType {
	parts := strings.Split(string(p), ":")
	if len(parts) != 2 {
		return ""
	}
	return ResourceType(parts[0])
}

// GetActionType returns the action type component of the permission
func (p Permission) GetActionType() ActionType {
	parts := strings.Split(string(p), ":")
	if len(parts) != 2 {
		return ""
	}
	return ActionType(parts[1])
}

// Value returns the underlying string value
func (p Permission) Value() string {
	return string(p)
}

// String returns the string representation
func (p Permission) String() string {
	return string(p)
}

// Equals checks if two Permissions are equal
func (p Permission) Equals(other Permission) bool {
	return p == other
}
