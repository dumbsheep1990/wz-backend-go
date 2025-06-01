package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

// Username represents an admin username value object
type Username string

// UsernameMinLength is the minimum length for a username
const UsernameMinLength = 3

// UsernameMaxLength is the maximum length for a username
const UsernameMaxLength = 32

// Username validation regex: alphanumeric, underscore, and hyphen
var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// NewUsername creates a new Username value object
func NewUsername(username string) (Username, error) {
	// Trim spaces
	username = strings.TrimSpace(username)
	
	// Check length
	if len(username) < UsernameMinLength {
		return "", errors.New("username must be at least 3 characters")
	}
	if len(username) > UsernameMaxLength {
		return "", errors.New("username must be no more than 32 characters")
	}
	
	// Check pattern
	if !usernamePattern.MatchString(username) {
		return "", errors.New("username must contain only letters, numbers, underscores, and hyphens")
	}
	
	return Username(username), nil
}

// MustNewUsername creates a new Username and panics if invalid
func MustNewUsername(username string) Username {
	u, err := NewUsername(username)
	if err != nil {
		panic(err)
	}
	return u
}

// Value returns the underlying string value
func (u Username) Value() string {
	return string(u)
}

// String returns the string representation
func (u Username) String() string {
	return string(u)
}

// Equals checks if two Usernames are equal
func (u Username) Equals(other Username) bool {
	return strings.EqualFold(string(u), string(other))
}

// IsEmpty checks if the username is empty
func (u Username) IsEmpty() bool {
	return len(strings.TrimSpace(string(u))) == 0
}
