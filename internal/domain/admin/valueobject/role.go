package valueobject

import (
	"errors"
	"fmt"
	"strings"
)

// RoleID represents a unique identifier for a role
type RoleID string

// RoleName represents the name of a role
type RoleName string

// NewRoleID creates a new RoleID value object
func NewRoleID(id string) (RoleID, error) {
	if strings.TrimSpace(id) == "" {
		return "", errors.New("role ID cannot be empty")
	}
	return RoleID(id), nil
}

// MustNewRoleID creates a new RoleID and panics if invalid
func MustNewRoleID(id string) RoleID {
	roleID, err := NewRoleID(id)
	if err != nil {
		panic(err)
	}
	return roleID
}

// Value returns the underlying value
func (id RoleID) Value() string {
	return string(id)
}

// String returns the string representation
func (id RoleID) String() string {
	return string(id)
}

// Equals checks if two RoleIDs are equal
func (id RoleID) Equals(other RoleID) bool {
	return id == other
}

// NewRoleName creates a new RoleName value object
func NewRoleName(name string) (RoleName, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("role name cannot be empty")
	}
	if len(name) > 50 {
		return "", fmt.Errorf("role name cannot exceed 50 characters")
	}
	return RoleName(name), nil
}

// MustNewRoleName creates a new RoleName and panics if invalid
func MustNewRoleName(name string) RoleName {
	roleName, err := NewRoleName(name)
	if err != nil {
		panic(err)
	}
	return roleName
}

// Value returns the underlying value
func (n RoleName) Value() string {
	return string(n)
}

// String returns the string representation
func (n RoleName) String() string {
	return string(n)
}

// Equals checks if two RoleNames are equal
func (n RoleName) Equals(other RoleName) bool {
	return strings.EqualFold(string(n), string(other))
}
