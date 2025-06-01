package valueobject

import (
	"errors"
	"fmt"
)

// AdminID represents a unique identifier for an admin user
type AdminID int64

// NewAdminID creates a new AdminID value object
func NewAdminID(id int64) (AdminID, error) {
	if id <= 0 {
		return 0, errors.New("admin ID must be a positive number")
	}
	return AdminID(id), nil
}

// MustNewAdminID creates a new AdminID and panics if invalid
func MustNewAdminID(id int64) AdminID {
	adminID, err := NewAdminID(id)
	if err != nil {
		panic(err)
	}
	return adminID
}

// Value returns the underlying value
func (id AdminID) Value() int64 {
	return int64(id)
}

// String returns the string representation
func (id AdminID) String() string {
	return fmt.Sprintf("%d", id)
}

// Equals checks if two AdminIDs are equal
func (id AdminID) Equals(other AdminID) bool {
	return id == other
}
