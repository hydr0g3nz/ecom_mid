package user

import (
	"github.com/hydr0g3nz/ecom_mid/domain/common"
)

// Permission represents a permission in the system
type Permission struct {
	common.Entity
	ID          uint   `json:"permission_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Role represents a role in the system
type Role struct {
	common.Entity
	ID          uint         `json:"role_id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}

// AddPermission adds a permission to the role
func (r *Role) AddPermission(permission Permission) {
	// Check if permission already exists
	for _, p := range r.Permissions {
		if p.ID == permission.ID {
			return
		}
	}
	r.Permissions = append(r.Permissions, permission)
}

// RemovePermission removes a permission from the role
func (r *Role) RemovePermission(permissionID uint) {
	updatedPermissions := []Permission{}
	for _, perm := range r.Permissions {
		if perm.ID != permissionID {
			updatedPermissions = append(updatedPermissions, perm)
		}
	}
	r.Permissions = updatedPermissions
}

// HasPermission checks if the role has a specific permission
func (r *Role) HasPermission(permissionName string) bool {
	for _, perm := range r.Permissions {
		if perm.Name == permissionName {
			return true
		}
	}
	return false
}
