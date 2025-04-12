package user

import (
	"errors"
	"time"

	"github.com/hydr0g3nz/ecom_mid/domain/common"
)

// StaffStatus represents the status of a staff member
type StaffStatus string

const (
	StaffStatusActive    StaffStatus = "active"
	StaffStatusInactive  StaffStatus = "inactive"
	StaffStatusSuspended StaffStatus = "suspended"
)

// Staff represents a staff member in the system
type Staff struct {
	common.Entity
	StaffID   uint       `json:"staff_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Phone     string     `json:"phone"`
	Position  string     `json:"position"`
	LastLogin *time.Time `json:"last_login"`
	Status    StaffStatus `json:"status"`
	Roles     []Role     `json:"roles"`
}

// NewStaff creates a new staff with validation
func NewStaff(username, email, firstName, lastName, phone, position string) (*Staff, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	
	if firstName == "" || lastName == "" {
		return nil, errors.New("name cannot be empty")
	}
	
	return &Staff{
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Position:  position,
		Status:    StaffStatusActive,
		Roles:     []Role{},
	}, nil
}

// FullName returns the staff's full name
func (s *Staff) FullName() string {
	return s.FirstName + " " + s.LastName
}

// IsActive checks if the staff is active
func (s *Staff) IsActive() bool {
	return s.Status == StaffStatusActive
}

// Suspend sets the staff's status to suspended
func (s *Staff) Suspend() {
	s.Status = StaffStatusSuspended
}

// Activate sets the staff's status to active
func (s *Staff) Activate() {
	s.Status = StaffStatusActive
}

// Deactivate sets the staff's status to inactive
func (s *Staff) Deactivate() {
	s.Status = StaffStatusInactive
}

// UpdateLastLogin updates the last login time to now
func (s *Staff) UpdateLastLogin() {
	now := time.Now()
	s.LastLogin = &now
}

// HasPermission checks if the staff has a specific permission
func (s *Staff) HasPermission(permissionName string) bool {
	for _, role := range s.Roles {
		for _, permission := range role.Permissions {
			if permission.Name == permissionName {
				return true
			}
		}
	}
	return false
}

// AddRole adds a role to the staff
func (s *Staff) AddRole(role Role) {
	// Check if role already exists
	for _, r := range s.Roles {
		if r.ID == role.ID {
			return
		}
	}
	s.Roles = append(s.Roles, role)
}

// RemoveRole removes a role from the staff
func (s *Staff) RemoveRole(roleID uint) {
	updatedRoles := []Role{}
	for _, role := range s.Roles {
		if role.ID != roleID {
			updatedRoles = append(updatedRoles, role)
		}
	}
	s.Roles = updatedRoles
}
