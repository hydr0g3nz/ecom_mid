package user

import (
	"errors"
	"time"

	"github.com/hydr0g3nz/ecom_mid/domain/common"
	"github.com/hydr0g3nz/ecom_mid/domain/common/vo"
)

// CustomerStatus represents the status of a customer
type CustomerStatus string

const (
	CustomerStatusActive    CustomerStatus = "active"
	CustomerStatusInactive  CustomerStatus = "inactive"
	CustomerStatusSuspended CustomerStatus = "suspended"
)

// Gender represents the gender of a customer
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// AddressType represents the type of an address
type AddressType string

const (
	AddressTypeShipping AddressType = "shipping"
	AddressTypeBilling  AddressType = "billing"
)

// CustomerAddress represents a customer's address
type CustomerAddress struct {
	common.Entity
	AddressID   uint        `json:"address_id"`
	CustomerID  uint        `json:"customer_id"`
	AddressType AddressType `json:"address_type"`
	Address     vo.Address  `json:"address"`
}

// Customer represents a customer in the system
type Customer struct {
	common.Entity
	CustomerID   uint            `json:"customer_id"`
	Username     string          `json:"username"`
	Email        string          `json:"email"`
	FirstName    string          `json:"first_name"`
	LastName     string          `json:"last_name"`
	Phone        string          `json:"phone"`
	BirthDate    *time.Time      `json:"birth_date"`
	Gender       Gender          `json:"gender"`
	ProfileImage string          `json:"profile_image"`
	LastLogin    *time.Time      `json:"last_login"`
	Status       CustomerStatus  `json:"status"`
	Points       int             `json:"points"`
	Addresses    []CustomerAddress `json:"addresses"`
}

// NewCustomer creates a new customer with validation
func NewCustomer(username, email, firstName, lastName, phone string) (*Customer, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	
	if firstName == "" || lastName == "" {
		return nil, errors.New("name cannot be empty")
	}
	
	return &Customer{
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Status:    CustomerStatusActive,
		Points:    0,
		Addresses: []CustomerAddress{},
	}, nil
}

// FullName returns the customer's full name
func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}

// IsActive checks if the customer is active
func (c *Customer) IsActive() bool {
	return c.Status == CustomerStatusActive
}

// Suspend sets the customer's status to suspended
func (c *Customer) Suspend() {
	c.Status = CustomerStatusSuspended
}

// Activate sets the customer's status to active
func (c *Customer) Activate() {
	c.Status = CustomerStatusActive
}

// Deactivate sets the customer's status to inactive
func (c *Customer) Deactivate() {
	c.Status = CustomerStatusInactive
}

// UpdateLastLogin updates the last login time to now
func (c *Customer) UpdateLastLogin() {
	now := time.Now()
	c.LastLogin = &now
}

// AddPoints adds points to the customer
func (c *Customer) AddPoints(points int) {
	if points > 0 {
		c.Points += points
	}
}

// UsePoints uses customer points and returns true if successful
func (c *Customer) UsePoints(points int) bool {
	if points <= 0 {
		return false
	}
	
	if c.Points >= points {
		c.Points -= points
		return true
	}
	
	return false
}

// AddAddress adds an address to the customer
func (c *Customer) AddAddress(addr CustomerAddress) {
	// If this is the first address, make it default
	if len(c.Addresses) == 0 {
		addr.Address.IsDefault = true
	}
	
	// If this is marked as default, remove default from others of same type
	if addr.Address.IsDefault {
		for i := range c.Addresses {
			if c.Addresses[i].AddressType == addr.AddressType {
				c.Addresses[i].Address.IsDefault = false
			}
		}
	}
	
	c.Addresses = append(c.Addresses, addr)
}

// UpdateAddress updates an existing address
func (c *Customer) UpdateAddress(addressID uint, newAddress vo.Address) error {
	for i := range c.Addresses {
		if c.Addresses[i].AddressID == addressID {
			// If setting as default, update others of same type
			if newAddress.IsDefault && !c.Addresses[i].Address.IsDefault {
				for j := range c.Addresses {
					if c.Addresses[j].AddressType == c.Addresses[i].AddressType {
						c.Addresses[j].Address.IsDefault = false
					}
				}
			}
			
			c.Addresses[i].Address = newAddress
			return nil
		}
	}
	
	return errors.New("address not found")
}

// RemoveAddress removes an address
func (c *Customer) RemoveAddress(addressID uint) error {
	for i, addr := range c.Addresses {
		if addr.AddressID == addressID {
			// Remove address from the slice
			c.Addresses = append(c.Addresses[:i], c.Addresses[i+1:]...)
			
			// If the removed address was default, set a new default if possible
			if addr.Address.IsDefault {
				for j := range c.Addresses {
					if c.Addresses[j].AddressType == addr.AddressType {
						c.Addresses[j].Address.IsDefault = true
						break
					}
				}
			}
			
			return nil
		}
	}
	
	return errors.New("address not found")
}

// GetDefaultAddress returns the default address of a specific type
func (c *Customer) GetDefaultAddress(addrType AddressType) (*CustomerAddress, error) {
	for _, addr := range c.Addresses {
		if addr.AddressType == addrType && addr.Address.IsDefault {
			return &addr, nil
		}
	}
	
	return nil, errors.New("no default address found")
}
