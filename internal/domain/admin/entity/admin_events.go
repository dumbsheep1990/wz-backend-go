package entity

import (
	"time"

	"wz-backend-go/internal/domain/admin/valueobject"
)

// DomainEvent is the base interface for all domain events
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
	EntityID() string
}

// BaseDomainEvent provides common functionality for all domain events
type BaseDomainEvent struct {
	occurredAt time.Time
	entityID   string
}

// OccurredAt returns when the event occurred
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EntityID returns the ID of the entity that triggered the event
func (e BaseDomainEvent) EntityID() string {
	return e.entityID
}

// AdminCreatedEvent is triggered when a new admin is created
type AdminCreatedEvent struct {
	BaseDomainEvent
	AdminID   valueobject.AdminID
	Username  valueobject.Username
	Role      valueobject.RoleID
}

// EventType returns the type of this event
func (e AdminCreatedEvent) EventType() string {
	return "admin.created"
}

// NewAdminCreatedEvent creates a new AdminCreatedEvent
func NewAdminCreatedEvent(admin *Admin) AdminCreatedEvent {
	return AdminCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			occurredAt: time.Now(),
			entityID:   admin.ID.String(),
		},
		AdminID:   admin.ID,
		Username:  admin.Username,
		Role:      admin.Role,
	}
}

// AdminRoleChangedEvent is triggered when an admin's role is changed
type AdminRoleChangedEvent struct {
	BaseDomainEvent
	AdminID   valueobject.AdminID
	OldRole   valueobject.RoleID
	NewRole   valueobject.RoleID
}

// EventType returns the type of this event
func (e AdminRoleChangedEvent) EventType() string {
	return "admin.role_changed"
}

// NewAdminRoleChangedEvent creates a new AdminRoleChangedEvent
func NewAdminRoleChangedEvent(admin *Admin, oldRole valueobject.RoleID) AdminRoleChangedEvent {
	return AdminRoleChangedEvent{
		BaseDomainEvent: BaseDomainEvent{
			occurredAt: time.Now(),
			entityID:   admin.ID.String(),
		},
		AdminID:   admin.ID,
		OldRole:   oldRole,
		NewRole:   admin.Role,
	}
}

// AdminStatusChangedEvent is triggered when an admin's status is changed
type AdminStatusChangedEvent struct {
	BaseDomainEvent
	AdminID    valueobject.AdminID
	OldStatus  valueobject.AdminStatus
	NewStatus  valueobject.AdminStatus
}

// EventType returns the type of this event
func (e AdminStatusChangedEvent) EventType() string {
	return "admin.status_changed"
}

// NewAdminStatusChangedEvent creates a new AdminStatusChangedEvent
func NewAdminStatusChangedEvent(admin *Admin, oldStatus valueobject.AdminStatus) AdminStatusChangedEvent {
	return AdminStatusChangedEvent{
		BaseDomainEvent: BaseDomainEvent{
			occurredAt: time.Now(),
			entityID:   admin.ID.String(),
		},
		AdminID:    admin.ID,
		OldStatus:  oldStatus,
		NewStatus:  admin.Status,
	}
}

// AdminLoggedInEvent is triggered when an admin logs in
type AdminLoggedInEvent struct {
	BaseDomainEvent
	AdminID   valueobject.AdminID
	Username  valueobject.Username
	IP        string
	UserAgent string
}

// EventType returns the type of this event
func (e AdminLoggedInEvent) EventType() string {
	return "admin.logged_in"
}

// NewAdminLoggedInEvent creates a new AdminLoggedInEvent
func NewAdminLoggedInEvent(admin *Admin, ip, userAgent string) AdminLoggedInEvent {
	return AdminLoggedInEvent{
		BaseDomainEvent: BaseDomainEvent{
			occurredAt: time.Now(),
			entityID:   admin.ID.String(),
		},
		AdminID:   admin.ID,
		Username:  admin.Username,
		IP:        ip,
		UserAgent: userAgent,
	}
}

// AdminPasswordChangedEvent is triggered when an admin changes their password
type AdminPasswordChangedEvent struct {
	BaseDomainEvent
	AdminID  valueobject.AdminID
	Username valueobject.Username
}

// EventType returns the type of this event
func (e AdminPasswordChangedEvent) EventType() string {
	return "admin.password_changed"
}

// NewAdminPasswordChangedEvent creates a new AdminPasswordChangedEvent
func NewAdminPasswordChangedEvent(admin *Admin) AdminPasswordChangedEvent {
	return AdminPasswordChangedEvent{
		BaseDomainEvent: BaseDomainEvent{
			occurredAt: time.Now(),
			entityID:   admin.ID.String(),
		},
		AdminID:  admin.ID,
		Username: admin.Username,
	}
}

// RoleCreatedEvent is triggered when a new role is created
type RoleCreatedEvent struct {
	BaseDomainEvent
	RoleID      valueobject.RoleID
	Name        valueobject.RoleName
	Description string
}

// EventType returns the type of this event
func (e RoleCreatedEvent) EventType() string {
	return "role.created"
}

// NewRoleCreatedEvent creates a new RoleCreatedEvent
func NewRoleCreatedEvent(role *Role) RoleCreatedEvent {
	return RoleCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			occurredAt: time.Now(),
			entityID:   role.ID.String(),
		},
		RoleID:      role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}

// RoleUpdatedEvent is triggered when a role is updated
type RoleUpdatedEvent struct {
	BaseDomainEvent
	RoleID      valueobject.RoleID
	Name        valueobject.RoleName
	Description string
}

// EventType returns the type of this event
func (e RoleUpdatedEvent) EventType() string {
	return "role.updated"
}

// NewRoleUpdatedEvent creates a new RoleUpdatedEvent
func NewRoleUpdatedEvent(role *Role) RoleUpdatedEvent {
	return RoleUpdatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			occurredAt: time.Now(),
			entityID:   role.ID.String(),
		},
		RoleID:      role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}

// RolePermissionsChangedEvent is triggered when a role's permissions are changed
type RolePermissionsChangedEvent struct {
	BaseDomainEvent
	RoleID      valueobject.RoleID
	Permissions []valueobject.Permission
}

// EventType returns the type of this event
func (e RolePermissionsChangedEvent) EventType() string {
	return "role.permissions_changed"
}

// NewRolePermissionsChangedEvent creates a new RolePermissionsChangedEvent
func NewRolePermissionsChangedEvent(role *Role) RolePermissionsChangedEvent {
	return RolePermissionsChangedEvent{
		BaseDomainEvent: BaseDomainEvent{
			occurredAt: time.Now(),
			entityID:   role.ID.String(),
		},
		RoleID:      role.ID,
		Permissions: role.Permissions,
	}
}
