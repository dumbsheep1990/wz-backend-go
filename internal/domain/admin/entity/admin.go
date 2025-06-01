package entity

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"wz-backend-go/internal/domain/admin/valueobject"
)

// Admin represents an administrator user entity
type Admin struct {
	ID        valueobject.AdminID
	Username  valueobject.Username
	Password  string // Hashed password
	Role      valueobject.RoleID
	Status    valueobject.AdminStatus
	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAdmin creates a new Admin entity
func NewAdmin(
	id valueobject.AdminID,
	username valueobject.Username,
	password string,
	role valueobject.RoleID,
	status valueobject.AdminStatus,
) (*Admin, error) {
	// Hash the password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	
	admin := &Admin{
		ID:        id,
		Username:  username,
		Password:  hashedPassword,
		Role:      role,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	return admin, nil
}

// ChangePassword changes the admin's password
func (a *Admin) ChangePassword(currentPassword, newPassword string) error {
	// Verify current password
	if !a.VerifyPassword(currentPassword) {
		return errors.New("current password is incorrect")
	}
	
	// Hash and set new password
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	
	a.Password = hashedPassword
	a.UpdatedAt = time.Now()
	
	return nil
}

// VerifyPassword checks if the provided password matches the admin's password
func (a *Admin) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

// UpdateRole updates the admin's role
func (a *Admin) UpdateRole(role valueobject.RoleID) {
	a.Role = role
	a.UpdatedAt = time.Now()
}

// UpdateStatus updates the admin's status
func (a *Admin) UpdateStatus(status valueobject.AdminStatus) error {
	switch status {
	case valueobject.AdminStatusActive:
		if !a.Status.CanActivate() {
			return errors.New("admin account cannot be activated from its current status")
		}
	case valueobject.AdminStatusDisabled:
		if !a.Status.CanDisable() {
			return errors.New("admin account cannot be disabled from its current status")
		}
	case valueobject.AdminStatusLocked:
		if !a.Status.CanLock() {
			return errors.New("admin account cannot be locked from its current status")
		}
	default:
		return errors.New("invalid admin status")
	}
	
	a.Status = status
	a.UpdatedAt = time.Now()
	return nil
}

// RecordLogin records a successful login
func (a *Admin) RecordLogin() {
	a.LastLogin = time.Now()
	a.UpdatedAt = time.Now()
}

// IsActive checks if the admin account is active
func (a *Admin) IsActive() bool {
	return a.Status.IsActive()
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	return string(hashedBytes), nil
}
