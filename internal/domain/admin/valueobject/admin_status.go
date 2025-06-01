package valueobject

import (
	"errors"
)

// AdminStatus represents the status of an admin user
type AdminStatus int

const (
	// AdminStatusDisabled represents a disabled admin account
	AdminStatusDisabled AdminStatus = 0
	
	// AdminStatusActive represents an active admin account
	AdminStatusActive AdminStatus = 1
	
	// AdminStatusLocked represents a locked admin account
	AdminStatusLocked AdminStatus = 2
)

// NewAdminStatus creates a new AdminStatus value object
func NewAdminStatus(status int) (AdminStatus, error) {
	switch status {
	case int(AdminStatusDisabled), int(AdminStatusActive), int(AdminStatusLocked):
		return AdminStatus(status), nil
	default:
		return AdminStatusDisabled, errors.New("invalid admin status")
	}
}

// MustNewAdminStatus creates a new AdminStatus and panics if invalid
func MustNewAdminStatus(status int) AdminStatus {
	s, err := NewAdminStatus(status)
	if err != nil {
		panic(err)
	}
	return s
}

// Value returns the underlying integer value
func (s AdminStatus) Value() int {
	return int(s)
}

// String returns a string representation of the status
func (s AdminStatus) String() string {
	switch s {
	case AdminStatusDisabled:
		return "Disabled"
	case AdminStatusActive:
		return "Active"
	case AdminStatusLocked:
		return "Locked"
	default:
		return "Unknown"
	}
}

// IsActive checks if the admin status is active
func (s AdminStatus) IsActive() bool {
	return s == AdminStatusActive
}

// CanActivate determines if this status can be activated
func (s AdminStatus) CanActivate() bool {
	return s == AdminStatusDisabled
}

// CanDisable determines if this status can be disabled
func (s AdminStatus) CanDisable() bool {
	return s == AdminStatusActive
}

// CanLock determines if this status can be locked
func (s AdminStatus) CanLock() bool {
	return s == AdminStatusActive
}

// CanUnlock determines if this status can be unlocked
func (s AdminStatus) CanUnlock() bool {
	return s == AdminStatusLocked
}
