package inventory

import (
	"github.com/hydr0g3nz/ecom_mid/domain/common"
	"github.com/hydr0g3nz/ecom_mid/domain/common/vo"
)

// SupplierStatus represents the status of a supplier
type SupplierStatus string

const (
	SupplierStatusActive   SupplierStatus = "active"
	SupplierStatusInactive SupplierStatus = "inactive"
)

// Supplier represents a vendor that provides products
type Supplier struct {
	common.Entity
	SupplierID     uint           `json:"supplier_id"`
	Name           string         `json:"name"`
	ContactPerson  string         `json:"contact_person"`
	Email          string         `json:"email"`
	Phone          string         `json:"phone"`
	Address        string         `json:"address"`
	Status         SupplierStatus `json:"status"`
	PurchaseOrders []PurchaseOrder `json:"purchase_orders,omitempty"`
}

// IsActive checks if the supplier is active
func (s *Supplier) IsActive() bool {
	return s.Status == SupplierStatusActive
}

// Activate sets the supplier status to active
func (s *Supplier) Activate() {
	s.Status = SupplierStatusActive
}

// Deactivate sets the supplier status to inactive
func (s *Supplier) Deactivate() {
	s.Status = SupplierStatusInactive
}
